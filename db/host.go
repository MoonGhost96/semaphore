package db

import (
	"fmt"
	"net"
)

type Host struct {
	ID        int     `db:"id" json:"id"`
	ProjectID int     `db:"project_id" json:"project_id"`
	Name      string  `db:"name" json:"name" binding:"required"`
	HostIP    string  `db:"host_IP" json:"host_IP"`
	UserName  string  `db:"user_name" json:"user_name"`
	Password  string  `db:"password" json:"password"`
	SSHKey    *string `db:"sshKey" json:"sshKey"`
}

func (host Host) Validate() error {
	addr := net.ParseIP(host.HostIP)
	if addr == nil {
		return fmt.Errorf("host ip is not valid")
	}

	return nil
}
