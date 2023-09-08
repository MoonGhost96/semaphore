package bolt

import "github.com/ansible-semaphore/semaphore/db"

func (d *BoltDb) GetHostInvRel(projectID int, hostInvRelID int) (hostInvRel db.HostInventoryRel, err error) {
	err = d.getObject(projectID, db.HostInvRelProps, intObjectID(hostInvRelID), &hostInvRel)

	if err != nil {
		return
	}

	return
}

func (d *BoltDb) GetHostInvRels(projectID int, params db.RetrieveQueryParams) ([]db.HostInventoryRel, error) {
	var hostInvRels []db.HostInventoryRel
	err := d.getObjects(projectID, db.HostInvRelProps, params, nil, &hostInvRels)
	return hostInvRels, err
}

func (d *BoltDb) GetHostInvRelRefs(projectID int, hostInvRelID int) (db.ObjectReferrers, error) {
	return d.getObjectRefs(projectID, db.HostInvRelProps, hostInvRelID)
}

func (d *BoltDb) DeleteHostInvRel(projectID int, hostInvRelID int) error {
	return d.deleteObject(projectID, db.HostInvRelProps, intObjectID(hostInvRelID), nil)
}

func (d *BoltDb) UpdateHostInvRel(hostInvRel db.HostInventoryRel) error {
	return d.updateObject(hostInvRel.ProjectId, db.InventoryProps, hostInvRel)
}

func (d *BoltDb) CreateHostInvRel(hostInvRel db.HostInventoryRel) (newHostInvRel db.HostInventoryRel, err error) {
	newInventory, err := d.createObject(hostInvRel.ProjectId, db.HostInvRelProps, hostInvRel)
	return newInventory.(db.HostInventoryRel), err
}
