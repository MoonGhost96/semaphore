package sql

import "github.com/ansible-semaphore/semaphore/db"

func (d *SqlDb) GetHost(projectID int, hostID int) (host db.Host, err error) {
	err = d.getObject(projectID, db.HostProps, hostID, &host)

	if err != nil {
		return
	}

	return
}

func (d *SqlDb) GetHosts(projectID int, params db.RetrieveQueryParams) ([]db.Host, error) {
	var hosts []db.Host
	err := d.getObjects(projectID, db.HostProps, params, &hosts)
	return hosts, err
}

func (d *SqlDb) GetHostRefs(projectID int, hostID int) (db.ObjectReferrers, error) {
	return d.getObjectRefs(projectID, db.HostProps, hostID)
}

func (d *SqlDb) DeleteHost(projectID int, hostID int) error {
	return d.deleteObject(projectID, db.HostProps, hostID)
}

func (d *SqlDb) UpdateHost(host db.Host) error {
	_, err := d.exec(
		"update project__host set name=?, host_ip=?, user_name=?, password=?, sshKey=? where id=?",
		host.Name,
		host.HostIP,
		host.UserName,
		host.Password,
		host.SSHKey,
		host.ID)

	return err
}

func (d *SqlDb) CreateHost(host db.Host) (newHost db.Host, err error) {
	insertID, err := d.insert(
		"id",
		"insert into project__host (project_id, name, host_ip, user_name, password, sshKey) values (?, ?, ?, ?, ?, ?)",
		host.ProjectID,
		host.Name,
		host.HostIP,
		host.UserName,
		host.Password,
		host.SSHKey)

	if err != nil {
		return
	}

	newHost = host
	newHost.ID = insertID
	return
}
