package db

type HostInventoryRel struct {
	ID          int `db:"id" json:"id"`
	ProjectId   int `db:"project_id" json:"project_id"`
	HostId      int `db:"host_id" json:"host_id"`
	InventoryId int `db:"inventory_id" json:"inventory_id"`
}
