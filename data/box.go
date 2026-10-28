package data

import "sort"

type Box struct {
	Location string `json:"location"`
	Type     string `json:"type"`
	ShelfId  uint32 `json:"shelfId"`
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
