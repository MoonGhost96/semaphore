package bolt

import "github.com/ansible-semaphore/semaphore/db"

func (d *BoltDb) GetHost(projectID int, hostID int) (host db.Host, err error) {
	err = d.getObject(projectID, db.HostProps, intObjectID(hostID), &host)

	if err != nil {
		return
	}

	return
}

func (d *BoltDb) GetHosts(projectID int, params db.RetrieveQueryParams) (hosts []db.Host, err error) {
	err = d.getObjects(projectID, db.HostProps, params, nil, &hosts)
	return
}

func (d *BoltDb) GetHostRefs(projectID int, hostID int) (db.ObjectReferrers, error) {
	return d.getObjectRefs(projectID, db.HostProps, hostID)
}

func (d *BoltDb) DeleteHost(projectID int, hostID int) error {
	return d.deleteObject(projectID, db.HostProps, intObjectID(hostID), nil)
}

func (d *BoltDb) UpdateHost(host db.Host) error {
	return d.updateObject(host.ProjectID, db.HostProps, host)
}

func (d *BoltDb) CreateHost(host db.Host) (db.Host, error) {
	newHost, err := d.createObject(host.ProjectID, db.HostProps, host)
	return newHost.(db.Host), err
}
