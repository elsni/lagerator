package data

// Typ
type Room struct {
	Location    string `json:"location"`
	WarehouseId uint32 `json:"warehouseId"`
}

// Liste
type RoomTable = DataTable[Room]

// GetRoomNamesforWarehouse returns rooms belonging to a warehouse.
func GetRoomNamesforWarehouse(rt *RoomTable, wid uint32) []Listentry {
	var names []Listentry
	for _, set := range *rt {
		if set.Data.WarehouseId == wid && !set.Deleted {
			names = append(names, Listentry{Id: set.ID, Name: set.Name})
		}
	}
	return names
}

// Checks if a room exists in the given warehouse
// IsRoomExistent checks whether a room exists in the warehouse.
func IsRoomExistent(rt *RoomTable, name string, wid uint32) (bool, uint32) {
	for _, set := range *rt {
		if set.Data.WarehouseId == wid && set.Name == name {
			return true, set.ID
		}
	}
	return false, 0
}
