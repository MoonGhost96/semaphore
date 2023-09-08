package projects

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/gorilla/context"
	"net/http"
)

// HostMiddleware ensures an inventory exists and loads it to the context
func HostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := context.Get(r, "project").(db.Project)
		hostID, err := helpers.GetIntParam("host_id", w, r)
		if err != nil {
			return
		}

		host, err := helpers.Store(r).GetHost(project.ID, hostID)

		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		context.Set(r, "host", host)
		next.ServeHTTP(w, r)
	})
}

func GetHostRefs(w http.ResponseWriter, r *http.Request) {
	host := context.Get(r, "host").(db.Host)
	refs, err := helpers.Store(r).GetHostRefs(host.ProjectID, host.ID)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, refs)
}

// GetHosts returns hosts from the database
func GetHosts(w http.ResponseWriter, r *http.Request) {
	if host := context.Get(r, "host"); host != nil {
		helpers.WriteJSON(w, http.StatusOK, host.(db.Host))
		return
	}

	project := context.Get(r, "project").(db.Project)

	hosts, err := helpers.Store(r).GetHosts(project.ID, helpers.QueryParams(r.URL))

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, hosts)
}

// AddHost adds a host to the database
func AddHost(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)

	var host db.Host

	if !helpers.Bind(w, r, &host) {
		return
	}

	if host.ProjectID != project.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	if err := host.Validate(); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "IP地址填写有误",
		})
		return
	}

	newHost, err := helpers.Store(r).CreateHost(host)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)

	objType := db.EventHost
	desc := "Host, Name: " + host.Name + ", IP:" + host.Name + " created"
	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &project.ID,
		ObjectType:  &objType,
		ObjectID:    &newHost.ID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	helpers.WriteJSON(w, http.StatusCreated, newHost)
}

// UpdateHost updates host in database
func UpdateHost(w http.ResponseWriter, r *http.Request) {
	oldHost := context.Get(r, "host").(db.Host)

	var host db.Host

	if !helpers.Bind(w, r, &host) {
		return
	}

	if host.ID != oldHost.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Host ID in body and URL must be the same",
		})
		return
	}

	if host.ProjectID != oldHost.ProjectID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	err := helpers.Store(r).UpdateHost(host)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveHost deletes a host and host_inv_rels related to it from the database
func RemoveHost(w http.ResponseWriter, r *http.Request) {
	host := context.Get(r, "host").(db.Host)
	var err error

	// Todo 应该该为批删除，一个个删太慢
	// Todo 一个逻辑内删除多张表内数据，应改为事务性操作
	// 1. 删除host_inv_rel表中该主机与所有Inventory的关联信息
	// 1.1 先查询该hostId与哪些inventory关联
	hostInvRels, err := helpers.Store(r).GetHostInvRels(host.ProjectID, db.RetrieveQueryParams{
		QueryIdName:  "host_id",
		QueryIdValue: host.ID,
	})
	if err != nil {
		helpers.WriteError(w, err)
		return
	}
	// 1.2 遍历删除Host_Inv_Rel表中当前project下所有关联该Host的信息
	for _, v := range hostInvRels {
		err = helpers.Store(r).DeleteHostInvRel(host.ProjectID, v.ID)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
	}
	// 2 删除host表中该主机信息
	err = helpers.Store(r).DeleteHost(host.ProjectID, host.ID)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	desc := "Host, Name: " + host.Name + ", IP:" + host.Name + " deleted"

	user := context.Get(r, "user").(*db.User)

	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   &host.ProjectID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}
