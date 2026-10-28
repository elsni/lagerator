package data

type Shelf struct {
	Location string `json:"location"`
	RoomId   uint32 `json:"roomId"`
}

type ShelfTable = DataTable[Shelf]

// Checks if a shelf exists in the current warehouse with a given room
// IsShelfExistent checks if a shelf exists in a room within a warehouse.
func IsShelfExistent(st *ShelfTable, shelfname string, roomname string, warehouseid uint32) (bool, uint32, uint32) {
	roomexist, roomid := IsRoomExistent(&Db.Rooms, roomname, warehouseid)
	if roomexist {
		for _, set := range *st {
			if set.Name == shelfname && set.Data.RoomId == roomid {
				return true, set.ID, roomid
			}
		}
	}
	return false, 0, 0
}

// GetWarehouseIdforShelf returns the warehouse id for a shelf id.
func GetWarehouseIdforShelf(shelfid uint32) uint32 {
	return GetWarehouseIdforRoom(GetRoomIdforShelf(shelfid))
}

// GetShelfNamesforRoom returns shelves belonging to a room.
func GetShelfNamesforRoom(st *ShelfTable, rid uint32) []Listentry {
	var names []Listentry
	for _, set := range *st {
		if set.Data.RoomId == rid && !set.Deleted {
			names = append(names, Listentry{Id: set.ID, Name: set.Name})
		}
	}
	return names
}
