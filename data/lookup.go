package data

// GetPrintNameById returns the printable name for an id or a fallback.
func GetPrintNameById[T CustomData](tbl *DataTable[T], id uint32, width int) string {
	if set, ok := tbl.GetPtr(id); ok {
		return set.GetPrintName(width)
	}
	return "not found"
}

// GetWarehouseIdforRoom returns the warehouse id for a room id.
func GetWarehouseIdforRoom(roomid uint32) uint32 {
	room, ok := Db.Rooms.GetPtr(roomid)
	if !ok {
		return 0
	}
	return room.Data.WarehouseId
}

// GetRoomIdforShelf returns the room id for a shelf id.
func GetRoomIdforShelf(shelfid uint32) uint32 {
	shelf, ok := Db.Shelves.GetPtr(shelfid)
	if !ok {
		return 0
	}
	return shelf.Data.RoomId
}

// GetShelfIdforBox returns the shelf id for a box id.
func GetShelfIdforBox(boxid uint32) uint32 {
	box, ok := Db.Boxes.GetPtr(boxid)
	if !ok {
		return 0
	}
	return box.Data.ShelfId
}
