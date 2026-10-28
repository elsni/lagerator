package data

type Warehouse struct {
	Location string `json:"location"`
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
