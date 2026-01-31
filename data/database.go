package data

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/elsni/lagerator/id"
)

type Database struct {
	CurrentWarehouse uint32         `json:"currentWarehouseid"`
	Warehouses       WarehouseTable `json:"warehouses"`
	Rooms            RoomTable      `json:"rooms"`
	Shelves          ShelfTable     `json:"shelves"`
	Boxes            BoxTable       `json:"boxes"`
	Items            ItemTable      `json:"items"`
	Categories       CategoryTable  `json:"categories"`
	Tags             TagTable       `json:"tags"`
}

// NewDatabase creates a Database with empty tables.
func NewDatabase() *Database {
	return &Database{
		Warehouses: NewDataTable[Warehouse](),
		Rooms:      NewDataTable[Room](),
		Shelves:    NewDataTable[Shelf](),
		Boxes:      NewDataTable[Box](),
		Items:      NewDataTable[Item](),
		Categories: NewDataTable[Category](),
		Tags:       NewDataTable[Tag](),
	}
}

// Save writes the database to the user config directory.
func (db *Database) Save() {
	dirname, _ := os.UserHomeDir()
	os.Mkdir(dirname+"/.lgrt", os.ModePerm)
	json, _ := json.Marshal(db)
	os.WriteFile(dirname+"/.lgrt/lgrtdata.json", json, 0644)
}

// Load reads the database file and updates the id source.
func (db *Database) Load() {
	dirname, _ := os.UserHomeDir()
	bytes, err := os.ReadFile(dirname + "/.lgrt/lgrtdata.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, db)
	if err != nil {
		panic("Data corrupted!")
	}
	id.IdSource.SetLastId(db.FindLastId())
}

// FindItem prints matching items for a search string.
func (db *Database) FindItem(searchstring string, sortname bool) {
	if len(Db.Items) == 0 {
		fmt.Println("no data")
		return
	}
	s := strings.ToLower(searchstring)
	Db.Items.PrintListFiltered(sortname, func(item Dataset[Item]) bool {
		return strings.Contains(strings.ToLower(item.Name), s) ||
			strings.Contains(strings.ToLower(item.Description), s) ||
			strings.Contains(strings.ToLower(item.Data.Location), s)
	})
}

// FindLastId returns the highest id across all tables.
func (db *Database) FindLastId() uint32 {
	var id uint32 = 0
	for _, cat := range db.Tags {
		if (cat.ID) > id {
			id = cat.ID
		}
	}
	for _, cat := range db.Categories {
		if (cat.ID) > id {
			id = cat.ID
		}
	}
	for _, wh := range db.Warehouses {
		if (wh.ID) > id {
			id = wh.ID
		}
	}
	for _, room := range db.Rooms {
		if (room.ID) > id {
			id = room.ID
		}
	}
	for _, shelf := range db.Shelves {
		if (shelf.ID) > id {
			id = shelf.ID
		}
	}
	for _, box := range db.Boxes {
		if (box.ID) > id {
			id = box.ID
		}
	}
	for _, item := range db.Items {
		if (item.ID) > id {
			id = item.ID
		}
	}
	return id
}

// CountTagOccurance counts how many entries use a tag id.
func (db *Database) CountTagOccurance(tagid uint32) int {
	count := 0
	for _, cat := range db.Categories {
		if slices.Contains(cat.Tags, tagid) {
			count++
		}
	}
	for _, wh := range db.Warehouses {
		if slices.Contains(wh.Tags, tagid) {
			count++
		}
	}
	for _, room := range db.Rooms {
		if slices.Contains(room.Tags, tagid) {
			count++
		}
	}
	for _, shelf := range db.Shelves {
		if slices.Contains(shelf.Tags, tagid) {
			count++
		}
	}
	for _, box := range db.Boxes {
		if slices.Contains(box.Tags, tagid) {
			count++
		}
	}
	for _, item := range db.Items {
		if slices.Contains(item.Tags, tagid) {
			count++
		}
	}
	return count
}

var Db = NewDatabase()

// GetTagList returns a comma-separated list of tag names.
func GetTagList(tagidlist []uint32) string {
	s := ""
	for _, id := range tagidlist {
		tag, ok := Db.Tags.GetPtr(id)
		if ok {
			s += tag.Name + ", "
		}
	}
	if len(s) == 0 {
		return ""
	}
	return s[:len(s)-2]
}

// GetTagIds resolves a comma-separated list of tag names to ids.
func GetTagIds(taglist string) []uint32 {
	result := []uint32{}
	tagnames := strings.Split(taglist, ",")
	for _, tag := range tagnames {
		cleantag := strings.TrimSpace(tag)
		if cleantag != "" {
			tagid, found := Db.Tags.GetFirstOccurance(strings.TrimSpace(tag))
			if !found {
				tagid = Db.Tags.AddSimple(cleantag)
			}
			result = append(result, tagid)
		}
	}
	return result
}
