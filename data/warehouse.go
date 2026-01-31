package data

import "github.com/elsni/lagerator/terminal"

type Warehouse struct {
	Location string `json:"location"`
}

// GetTableRow returns the warehouse row content.
func (d Warehouse) GetTableRow(ownid uint32) string {
	return ""
}

// GetTableHeader returns the warehouse table header.
func (d Warehouse) GetTableHeader() string {
	return terminal.ResetColor()
}

// Show prints warehouse details.
func (d Warehouse) Show() {
}

type WarehouseTable = DataTable[Warehouse]

// Checks if a room exists in the given warehouse
// IsWarehouseExistent checks if a warehouse name exists.
func IsWarehouseExistent(wt *WarehouseTable, name string) (bool, uint32) {
	for _, set := range *wt {
		if set.Name == name && !set.Deleted {
			return true, set.ID
		}
	}
	return false, 0
}
