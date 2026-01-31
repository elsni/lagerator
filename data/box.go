package data

import (
	"fmt"
	"sort"

	"github.com/elsni/lagerator/terminal"
)

type Box struct {
	Location string `json:"location"`
	Type     string `json:"type"`
	ShelfId  uint32 `json:"shelfId"`
}

// GetTableRow returns the formatted row for a box.
func (d Box) GetTableRow(ownid uint32) string {
	shelfName := GetPrintNameById(&Db.Shelves, d.ShelfId, 20)
	roomId := GetRoomIdforShelf(d.ShelfId)
	roomName := GetPrintNameById(&Db.Rooms, roomId, 20)
	whName := GetPrintNameById(&Db.Warehouses, GetWarehouseIdforRoom(roomId), 20)
	return fmt.Sprintf("%-20s %-20s %-20s", shelfName, roomName, whName)
}

// GetTableHeader returns the box table header.
func (d Box) GetTableHeader() string {
	return fmt.Sprintf("%-20s %-20s %-20s%s", "Shelf", "Room", "Warehouse", terminal.ResetColor())
}

// Show prints box details.
func (d Box) Show() {
	fmt.Printf("%s %s\n", terminal.GetLabelText("Location"), d.Location)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Type"), d.Type)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Shelf"), GetPrintNameById(&Db.Shelves, d.ShelfId, 20))
}

// GetRoomIdforBox returns the room id for a box id.
func GetRoomIdforBox(boxid uint32) uint32 {
	return GetRoomIdforShelf(GetShelfIdforBox(boxid))
}

// GetBoxNamesforShelf returns boxes belonging to a shelf.
func GetBoxNamesforShelf(st *BoxTable, sid uint32) []Listentry {
	var names []Listentry
	for _, set := range *st {
		if set.Data.ShelfId == sid && !set.Deleted {
			names = append(names, Listentry{Id: set.ID, Name: set.Name})
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].Name < names[j].Name
	})
	return names
}

type BoxTable = DataTable[Box]
