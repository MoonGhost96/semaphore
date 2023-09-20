package model

import "github.com/ansible-semaphore/semaphore/db"

// Inventory 在db.Inventory的基础上添加了Host信息
type Inventory struct {
	ID        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	ProjectID int    `json:"project_id"`
	Inventory string `json:"inventory"`

	SSHKeyID *int         `json:"ssh_key_id"`
	SSHKey   db.AccessKey `json:"-"`

	BecomeKeyID *int         `json:"become_key_id"`
	BecomeKey   db.AccessKey `json:"-"`

	Type string `json:"type"`

	HostInvRels []db.HostInventoryRel `json:"host_inv_rels"`
	Hosts       []db.Host             `json:"hosts"`
}

func ConvertInvModel2InvDB(inventoryModel Inventory) db.Inventory {
	return db.Inventory{
		ID:          inventoryModel.ID,
		Name:        inventoryModel.Name,
		ProjectID:   inventoryModel.ProjectID,
		Inventory:   inventoryModel.Inventory,
		SSHKeyID:    inventoryModel.SSHKeyID,
		SSHKey:      inventoryModel.SSHKey,
		BecomeKeyID: inventoryModel.BecomeKeyID,
		BecomeKey:   inventoryModel.BecomeKey,
		Type:        inventoryModel.Type,
	}
}

func ConvertInvDB2InvModel(inventoryDB db.Inventory, hostInvRels []db.HostInventoryRel, hosts []db.Host) Inventory {
	inventoryModel := Inventory{
		ID:          inventoryDB.ID,
		Name:        inventoryDB.Name,
		ProjectID:   inventoryDB.ProjectID,
		Inventory:   inventoryDB.Inventory,
		SSHKeyID:    inventoryDB.SSHKeyID,
		SSHKey:      inventoryDB.SSHKey,
		BecomeKeyID: inventoryDB.BecomeKeyID,
		BecomeKey:   inventoryDB.BecomeKey,
		Type:        inventoryDB.Type,
	}
	if hostInvRels != nil {
		inventoryModel.HostInvRels = hostInvRels
	}
	if hosts != nil {
		inventoryModel.Hosts = hosts
	}
	return inventoryModel
}
