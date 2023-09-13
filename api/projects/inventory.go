package projects

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/model"
	"net/http"

	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/context"
)

// InventoryMiddleware ensures an inventory exists and loads it to the context
func InventoryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := context.Get(r, "project").(db.Project)
		inventoryID, err := helpers.GetIntParam("inventory_id", w, r)
		if err != nil {
			return
		}

		inventoryDB, err := helpers.Store(r).GetInventory(project.ID, inventoryID)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		var hostInvRels []db.HostInventoryRel
		hostInvRels, err = helpers.Store(r).GetHostInvRels(project.ID, db.RetrieveQueryParams{
			QueryIdName:  "inventory_id",
			QueryIdValue: inventoryID,
		})
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		var hostIds []int
		var hosts []db.Host

		for _, v := range hostInvRels {
			hostId := v.HostId
			hostIds = append(hostIds, hostId)
		}

		// 目前如果hostIds传空，会查host表的所有host
		if hostIds != nil {
			hosts, err = helpers.Store(r).GetHosts(project.ID, db.RetrieveQueryParams{
				QueryIdName:   "id",
				QueryIdValues: hostIds,
			})
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
		}

		inventoryModel := model.ConvertInvDB2InvModel(inventoryDB, hostInvRels, hosts)

		context.Set(r, "inventory", inventoryModel)
		next.ServeHTTP(w, r)
	})
}

func GetInventoryRefs(w http.ResponseWriter, r *http.Request) {
	inventory := context.Get(r, "inventory").(model.Inventory)
	refs, err := helpers.Store(r).GetInventoryRefs(inventory.ProjectID, inventory.ID)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, refs)
}

// GetInventory returns inventories from the database
func GetInventory(w http.ResponseWriter, r *http.Request) {
	if inventory := context.Get(r, "inventory"); inventory != nil {
		helpers.WriteJSON(w, http.StatusOK, inventory.(model.Inventory))
		return
	}

	project := context.Get(r, "project").(db.Project)

	var inventories []model.Inventory

	inventoriesDB, err := helpers.Store(r).GetInventories(project.ID, helpers.QueryParams(r.URL))
	for _, v := range inventoriesDB {
		if v.Type != db.InventoryHost {
			inventories = append(inventories, model.ConvertInvDB2InvModel(v, nil, nil))
			continue
		}
		var hostInvRels []db.HostInventoryRel
		var hosts []db.Host
		var hostIds []int
		hostInvRels, err = helpers.Store(r).GetHostInvRels(v.ProjectID, db.RetrieveQueryParams{
			QueryIdName:  "inventory_id",
			QueryIdValue: v.ID,
		})
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		for _, hostInvRel := range hostInvRels {
			hostId := hostInvRel.HostId
			hostIds = append(hostIds, hostId)
		}

		if hostIds != nil {
			hosts, err = helpers.Store(r).GetHosts(project.ID, db.RetrieveQueryParams{
				// 注意下面QueryIdName要填成"id"，host表里没有host_id这一列
				QueryIdName:   "id",
				QueryIdValues: hostIds,
			})
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
		}

		inventories = append(inventories, model.ConvertInvDB2InvModel(v, hostInvRels, hosts))
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, inventories)
}

// AddInventory creates an inventory in the database
func AddInventory(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)

	var inventoryModel model.Inventory

	if !helpers.Bind(w, r, &inventoryModel) {
		return
	}

	if inventoryModel.ProjectID != project.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	switch inventoryModel.Type {
	case db.InventoryStatic, db.InventoryStaticYaml, db.InventoryFile, db.InventoryHost:
		break
	default:
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Not supported inventory type",
		})
		return
	}

	inventoryDB := model.ConvertInvModel2InvDB(inventoryModel)
	newInventoryDB, err := helpers.Store(r).CreateInventory(inventoryDB)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}
	// 如果inventoryModel.Type是绑定主机的新类型，添加关系表数据
	if inventoryModel.Type == db.InventoryHost {
		//todo 改成批量
		for i, _ := range inventoryModel.HostInvRels {
			inventoryModel.HostInvRels[i].InventoryId = newInventoryDB.ID
			_, err = helpers.Store(r).CreateHostInvRel(inventoryModel.HostInvRels[i])
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
		}
	}

	user := context.Get(r, "user").(*db.User)

	objType := db.EventInventory
	desc := "Inventory " + inventoryModel.Name + " created"
	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &project.ID,
		ObjectType:  &objType,
		ObjectID:    &newInventoryDB.ID,
		Description: &desc,
	})

	if err != nil {
		// Write error to log but return ok to user, because inventory created
		log.Error(err)
	}

	// 此处没有修改成返回Model结构，感觉没必要返回对象的
	helpers.WriteJSON(w, http.StatusCreated, newInventoryDB)
}

// IsValidInventoryPath tests a path to ensure it is below the cwd
func IsValidInventoryPath(path string) bool {

	currentPath, err := os.Getwd()
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	relPath, err := filepath.Rel(currentPath, absPath)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(relPath, "..")
}

// UpdateInventory writes updated values to an existing inventory item in the database
func UpdateInventory(w http.ResponseWriter, r *http.Request) {
	oldInventoryModel := context.Get(r, "inventory").(model.Inventory)

	var inventoryModel model.Inventory

	if !helpers.Bind(w, r, &inventoryModel) {
		return
	}

	if inventoryModel.ID != oldInventoryModel.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Inventory ID in body and URL must be the same",
		})
		return
	}

	if inventoryModel.ProjectID != oldInventoryModel.ProjectID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	switch inventoryModel.Type {
	case db.InventoryStatic, db.InventoryStaticYaml, db.InventoryHost:
		break
	case db.InventoryFile:
		if !IsValidInventoryPath(inventoryModel.Inventory) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inventoryDB := model.ConvertInvModel2InvDB(inventoryModel)
	err := helpers.Store(r).UpdateInventory(inventoryDB)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	if inventoryModel.Type == db.InventoryHost {
		// Todo 目前采取全量更新的形式，后续考虑改为增量
		for _, v := range inventoryModel.HostInvRels {
			if v.ID != 0 {
				err = helpers.Store(r).UpdateHostInvRel(v)
				if err != nil {
					helpers.WriteError(w, err)
					return
				}
			} else {
				_, err = helpers.Store(r).CreateHostInvRel(v)
				if err != nil {
					helpers.WriteError(w, err)
					return
				}
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveInventory deletes an inventory from the database
func RemoveInventory(w http.ResponseWriter, r *http.Request) {
	inventoryModel := context.Get(r, "inventory").(model.Inventory)
	var err error

	err = helpers.Store(r).DeleteInventory(inventoryModel.ProjectID, inventoryModel.ID)
	if err == db.ErrInvalidOperation {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Inventory is in use by one or more templates",
			"inUse": true,
		})
		return
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	// 1 删除host_inv_rel关系表中该inventory_id对应的项
	// 1.1 先查询该inventoryId与哪些host关联
	hostInvRels, err := helpers.Store(r).GetHostInvRels(inventoryModel.ProjectID, db.RetrieveQueryParams{
		QueryIdName:  "inventory_id",
		QueryIdValue: inventoryModel.ID,
	})
	if err != nil {
		helpers.WriteError(w, err)
		return
	}
	// 1.2 遍历删除Host_Inv_Rel表中当前project下所有关联该inventory的信息
	for _, v := range hostInvRels {
		err = helpers.Store(r).DeleteHostInvRel(inventoryModel.ProjectID, v.ID)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
	}

	desc := "Inventory " + inventoryModel.Name + " deleted"

	user := context.Get(r, "user").(*db.User)

	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &inventoryModel.ProjectID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
