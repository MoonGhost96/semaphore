package sql

import "github.com/ansible-semaphore/semaphore/db"

func (d *SqlDb) GetHostInvRel(projectID int, hostInvRelID int) (hostInvRel db.HostInventoryRel, err error) {
	err = d.getObject(projectID, db.HostInvRelProps, hostInvRelID, &hostInvRel)

	if err != nil {
		return
	}

	return
}

func (d *SqlDb) GetHostInvRels(projectID int, params db.RetrieveQueryParams) ([]db.HostInventoryRel, error) {
	var hostInvRels []db.HostInventoryRel
	err := d.getObjects(projectID, db.HostInvRelProps, params, &hostInvRels)
	return hostInvRels, err
}

func (d *SqlDb) GetHostInvRelRefs(projectID int, hostInvRelID int) (db.ObjectReferrers, error) {
	return d.getObjectRefs(projectID, db.HostInvRelProps, hostInvRelID)
}

func (d *SqlDb) DeleteHostInvRel(projectID int, hostInvRelID int) error {
	return d.deleteObject(projectID, db.HostInvRelProps, hostInvRelID)
}

func (d *SqlDb) UpdateHostInvRel(hostInvRel db.HostInventoryRel) error {
	_, err := d.exec(
		"update project__host__inventory__rel set host_id=?, inventory_id=? where id=?",
		hostInvRel.HostId,
		hostInvRel.InventoryId,
		hostInvRel.ID)

	return err
}

func (d *SqlDb) CreateHostInvRel(hostInvRel db.HostInventoryRel) (newHostInvRel db.HostInventoryRel, err error) {
	insertID, err := d.insert(
		"id",
		"insert into project__host__inventory__rel (project_id, host_id, inventory_id) values (?, ?, ?)",
		hostInvRel.ProjectId,
		hostInvRel.HostId,
		hostInvRel.InventoryId)

	if err != nil {
		return
	}

	newHostInvRel = hostInvRel
	newHostInvRel.ID = insertID
	return
}
