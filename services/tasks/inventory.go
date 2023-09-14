package tasks

import (
	"fmt"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/db/factory"
	"io/ioutil"
	"strconv"

	"github.com/ansible-semaphore/semaphore/util"
)

func (t *TaskRunner) installInventory() (err error) {
	if t.inventory.SSHKeyID != nil {
		err = t.inventory.SSHKey.Install(db.AccessKeyRoleAnsibleUser)
		if err != nil {
			return
		}
	}

	if t.inventory.BecomeKeyID != nil {
		err = t.inventory.BecomeKey.Install(db.AccessKeyRoleAnsibleBecomeUser)
		if err != nil {
			return
		}
	}

	if t.inventory.Type == db.InventoryStatic || t.inventory.Type == db.InventoryStaticYaml || t.inventory.Type == db.InventoryHost {
		err = t.installStaticInventory()
	}

	return
}

func (t *TaskRunner) installStaticInventory() error {
	t.Log("installing static inventory")

	path := util.Config.TmpPath + "/inventory_" + strconv.Itoa(t.task.ID)
	if t.inventory.Type == db.InventoryStaticYaml {
		path += ".yml"
	}

	var content string

	if t.inventory.Type == db.InventoryHost {
		content = getTmpInventoryFileContent(t)
	} else {
		content = t.inventory.Inventory
	}

	// create inventory file
	// 对于记录host账密的临时inventory文件，会在运行完任务后删除该临时文件
	return ioutil.WriteFile(path, []byte(content), 0664)
}

func getTmpInventoryFileContent(t *TaskRunner) (content string) {
	// 将各个host对应的账密信息写成inventory文件
	store := factory.CreateStore()
	store.Connect("root")

	hostInvRels, err := store.GetHostInvRels(t.inventory.ProjectID, db.RetrieveQueryParams{
		QueryIdName:  "inventory_id",
		QueryIdValue: t.inventory.ID,
	})
	if err != nil {
		t.Log("[ERROR] get db host_inventory_rels info error")
		return
	}

	var hostIds []int
	var hosts []db.Host

	for _, v := range hostInvRels {
		hostId := v.HostId
		hostIds = append(hostIds, hostId)
	}

	if hostIds != nil {
		hosts, err = store.GetHosts(t.inventory.ProjectID, db.RetrieveQueryParams{
			QueryIdName:   "id",
			QueryIdValues: hostIds,
		})
		if err != nil {
			t.Log("[ERROR] get db hosts info error")
			return
		}
	}
	for i, host := range hosts {
		userName := host.UserName
		password := host.Password
		ip := host.HostIP
		content += fmt.Sprintf("%s ansible_user=%s ansible_password=%s", ip, userName, password)
		if i != (len(hosts) - 1) {
			content += "\n"
		}
	}
	return
}
