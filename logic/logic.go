package logic

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/ui"
)

type tableKind int

const (
	kindUnknown tableKind = iota
	kindCategory
	kindWarehouse
	kindRoom
	kindShelf
	kindBox
	kindItem
)

// findTableById resolves an id to a table kind and index.
func findTableById(id uint32) (tableKind, int) {
	if idx := data.Db.Categories.GetIdx(id); idx > -1 {
		return kindCategory, idx
	}
	if idx := data.Db.Warehouses.GetIdx(id); idx > -1 {
		return kindWarehouse, idx
	}
	if idx := data.Db.Rooms.GetIdx(id); idx > -1 {
		return kindRoom, idx
	}
	if idx := data.Db.Shelves.GetIdx(id); idx > -1 {
		return kindShelf, idx
	}
	if idx := data.Db.Boxes.GetIdx(id); idx > -1 {
		return kindBox, idx
	}
	if idx := data.Db.Items.GetIdx(id); idx > -1 {
		return kindItem, idx
	}
	return kindUnknown, -1
}

// tableName returns a display name for a table kind.
func tableName(kind tableKind) string {
	switch kind {
	case kindCategory:
		return "category"
	case kindWarehouse:
		return "warehouse"
	case kindRoom:
		return "room"
	case kindShelf:
		return "shelf"
	case kindBox:
		return "box"
	case kindItem:
		return "item"
	default:
		return "record"
	}
}

// CurrentWarehouseExists reports whether the current warehouse id is valid.
func CurrentWarehouseExists() bool {
	if data.Db.CurrentWarehouse == 0 {
		return false
	}
	_, ok := data.Db.Warehouses.GetPtr(data.Db.CurrentWarehouse)
	return ok
}

// SwitchWarehouse selects the active warehouse by name.
func SwitchWarehouse(whname string) {
	index, id := data.Db.Warehouses.GetDataByName(whname)
	if index == -1 {
		fmt.Printf("warehouse \"%s\" does not exist\n", whname)
		return
	}
	data.Db.CurrentWarehouse = id
	data.Db.Save()
	fmt.Printf("Warehouse \"%s\" is now active\n", whname)
}

// AddWarehouse creates a new warehouse.
func AddWarehouse(whname string) {
	whexists, _ := data.IsWarehouseExistent(&data.Db.Warehouses, whname)
	if whexists {
		fmt.Printf("The warehouse \"%s\" aldready exists\n", whname)
		return
	}
	nwh := data.NewDataset[data.Warehouse](whname, data.Warehouse{})
	data.Db.Warehouses.Add(nwh)
	data.Db.Save()

	fmt.Printf("Added warehouse \"%s\" with ID %d\n", whname, nwh.ID)
}

// AddCategory creates a new category.
func AddCategory(cname string) {
	idx, _ := data.Db.Categories.GetDataByName(cname)
	if idx > 0 {
		fmt.Printf("The Categroy \"%s\" aldready exists\n", cname)
		return
	}
	nc := data.NewDataset[data.Category](cname, data.Category{})
	data.Db.Categories.Add(nc)
	data.Db.Save()

	fmt.Printf("Added Category \"%s\" with ID %d\n", cname, nc.ID)
}

// AddRoomToCurrentWarehouse creates a room in the active warehouse.
func AddRoomToCurrentWarehouse(roomname string) {
	if !CurrentWarehouseExists() {
		fmt.Println("Switch to valid warehouse first")
		return
	}
	wh, ok := data.Db.Warehouses.GetPtr(data.Db.CurrentWarehouse)
	if !ok {
		fmt.Println("Switch to valid warehouse first")
		return
	}
	roomexists, _ := data.IsRoomExistent(&data.Db.Rooms, roomname, data.Db.CurrentWarehouse)
	if roomexists {
		fmt.Printf("The room \"%s\" already exists in warehouse \"%s\"\n", roomname, wh.Name)
		return
	}
	nroom := data.NewDataset[data.Room](roomname, data.Room{WarehouseId: data.Db.CurrentWarehouse})
	data.Db.Rooms.Add(nroom)
	data.Db.Save()

	fmt.Printf("Added room \"%s\" with id %d to warehouse \"%s\"\n", roomname, nroom.ID, wh.Name)
}

// AddShelfToRoom creates a shelf in a room.
func AddShelfToRoom(shelfname string, roomname string) {
	if !CurrentWarehouseExists() {
		fmt.Println("Switch to valid warehouse first")
		return
	}
	idx := SelectSet(&data.Db.Rooms, roomname, "Room", "add")
	if idx < 0 {
		return
	}
	newshelf := data.NewDataset[data.Shelf](shelfname, data.Shelf{RoomId: data.Db.Rooms[idx].ID})
	data.Db.Shelves.Add(newshelf)
	data.Db.Save()
	fmt.Printf("Added shelf \"%s\" with id %d to room \"%s\"\n", shelfname, newshelf.ID, data.Db.Rooms[idx].Name)
}

// AddBoxToShelf creates a box on a shelf.
func AddBoxToShelf(boxname string, shelfname string) {
	if !CurrentWarehouseExists() {
		fmt.Println("Switch to valid warehouse first")
		return
	}
	idx := SelectSet(&data.Db.Shelves, shelfname, "Shelf", "add")
	if idx < 0 {
		return
	}

	newbox := data.NewDataset[data.Box](boxname, data.Box{ShelfId: data.Db.Shelves[idx].ID})
	data.Db.Boxes.Add(newbox)
	data.Db.Save()
	fmt.Printf("Added Box \"%s\" with id %d to shelf \"%s\"\n", boxname, newbox.ID, data.Db.Shelves[idx].Name)
}

// AddItems opens the item editor and appends items to a box.
func AddItems(boxnameorid string) {
	index := SelectSet(&data.Db.Boxes, boxnameorid, "Box", "add")
	oldcategory := uint32(0)
	oldlocation := ""
	if index == -1 {
		return
	}
	for {
		item := data.NewDataset[data.Item]("", data.Item{BoxId: data.Db.Boxes[index].ID, CategoryId: oldcategory, Location: oldlocation, Amount: 1})

		// build [][]IdOptions for populating the reference dropdowns
		idopts := ToDropDownOpts(data.GetBoxNamesforShelf(&data.Db.Boxes, data.Db.Boxes[index].Data.ShelfId))
		idopts = AppendToDropDownOpts(idopts, data.GetCategoriesSorted(&data.Db.Categories))
		// open form
		item, saved := ui.EditItem(item, idopts, " Add ")
		oldcategory = item.Data.CategoryId
		oldlocation = item.Data.Location
		if saved {
			data.Db.Items.Add(item)
			data.Db.Save()
		} else {
			return
		}
	}
}

// MoveItem moves an item to a different box.
func MoveItem(itemid uint32, boxnameorid string) {
	itemidx := data.Db.Items.GetIdx(itemid)
	if itemidx == -1 {
		fmt.Printf("No item with id %d found\n", itemid)
		return
	}
	boxidx := SelectSet(&data.Db.Boxes, boxnameorid, "Box", "move")
	if boxidx == -1 {
		return
	}
	data.Db.Items[itemidx].Data.BoxId = data.Db.Boxes[boxidx].ID
	data.Db.Items[itemidx].Updated = time.Now().Unix()
	data.Db.Save()
	fmt.Printf("Moved Item \"%s\" to \"%s\"\n", data.Db.Items[itemidx].Name, data.Db.Boxes[boxidx].Name)
}

// MoveBox moves a box to a different shelf.
func MoveBox(boxid uint32, shelfnameorid string) {
	boxidx := data.Db.Boxes.GetIdx(boxid)
	if boxidx == -1 {
		fmt.Printf("No box with id %d found\n", boxid)
		return
	}
	shelfidx := SelectSet(&data.Db.Shelves, shelfnameorid, "Shelf", "move")
	if shelfidx == -1 {
		return
	}
	data.Db.Boxes[boxidx].Data.ShelfId = data.Db.Shelves[shelfidx].ID
	data.Db.Boxes[boxidx].Updated = time.Now().Unix()
	data.Db.Save()
	fmt.Printf("Moved box \"%s\" to \"%s\"\n", data.Db.Boxes[boxidx].Name, data.Db.Shelves[shelfidx].Name)
}

// Convert Listentry slice to slice of slice of IdOptions (which is basically the same as Listentry) to prevent circular imports
// Used for populating the dropdown (reference-) fields of the edit form
// ToDropDownOpts converts list entries to dropdown options.
func ToDropDownOpts(list []data.Listentry) [][]ui.DropdownOptions {
	var out1 []ui.DropdownOptions
	var outlist [][]ui.DropdownOptions
	for _, e := range list {
		out1 = append(out1, ui.DropdownOptions{Name: e.Name, Id: e.Id})
	}
	return append(outlist, out1)
}

// Append slice of Listentry for Items-Form. Items contains a second reference field (categories)
// AppendToDropDownOpts appends list entries to dropdown options.
func AppendToDropDownOpts(idopts [][]ui.DropdownOptions, list []data.Listentry) [][]ui.DropdownOptions {
	var out1 []ui.DropdownOptions
	for _, e := range list {
		out1 = append(out1, ui.DropdownOptions{Name: e.Name, Id: e.Id})
	}
	return append(idopts, out1)
}

// gets the ID options to fill the parent and category dropdown fields in the edit mask
// GetDropDownOpts builds dropdown options for edit forms.
func GetDropDownOpts(id uint32) [][]ui.DropdownOptions {
	var outlist [][]ui.DropdownOptions
	var parentopts []ui.DropdownOptions
	var categoryopts []ui.DropdownOptions

	// check if id belongs to a room
	if _, ok := data.Db.Rooms.GetPtr(id); ok {
		// if so, parent dropdown shows warehouses
		for _, entry := range data.Db.Warehouses.GetNames() {
			parentopts = append(parentopts, ui.DropdownOptions{Id: entry.Id, Name: entry.Name})
		}
	}

	// check if id belongs to a shelf
	if _, ok := data.Db.Shelves.GetPtr(id); ok {
		// if so, parent dropdown shows the rooms of current warehouse
		for _, entry := range data.GetRoomNamesforWarehouse(&data.Db.Rooms, data.GetWarehouseIdforShelf(id)) {
			parentopts = append(parentopts, ui.DropdownOptions{Id: entry.Id, Name: entry.Name})
		}
	}

	// check if id belongs to a box
	if _, ok := data.Db.Boxes.GetPtr(id); ok {
		// if so, parent dropdown shows the shelves of the current room
		for _, entry := range data.GetShelfNamesforRoom(&data.Db.Shelves, data.GetRoomIdforBox(id)) {
			parentopts = append(parentopts, ui.DropdownOptions{Id: entry.Id, Name: entry.Name})
		}
	}

	// check if id belongs to an item
	if _, ok := data.Db.Items.GetPtr(id); ok {
		// if so, parent dropdown shows the boxes of the current shelf
		for _, entry := range data.GetBoxNamesforShelf(&data.Db.Boxes, data.GetShelfIdforItem(id)) {
			parentopts = append(parentopts, ui.DropdownOptions{Id: entry.Id, Name: entry.Name})
		}
	}

	// for category selection dropdown (edit item)
	for _, entry := range data.Db.Categories.GetNames() {
		categoryopts = append(categoryopts, ui.DropdownOptions{Id: entry.Id, Name: entry.Name})
	}
	outlist = append(outlist, parentopts)
	outlist = append(outlist, categoryopts)
	return outlist
}

// get the index of a set from its name or id. let the user select if name is ambigous
// returns the index of the found set or -1 if not found
// SelectSet resolves a name or id and may prompt when ambiguous.
func SelectSet[T data.CustomData](tbl *data.DataTable[T], setname string, tablename string, action string) int {
	sets := tbl.GetSetsByName(setname)
	set := data.Dataset[T]{}
	if len(sets) == 0 {
		id, err := strconv.ParseUint(setname, 10, 32)
		if err != nil {
			fmt.Printf("No %s with name \"%s\" found.\n", strings.ToLower(tablename), setname)
			return -1
		}
		idset, ok := tbl.GetPtr(uint32(id))
		if !ok {
			fmt.Printf("No %s with ID %s found.\n", strings.ToLower(tablename), setname)
			return -1
		}
		set = *idset
	} else if len(sets) > 1 {
		resultindex := ui.SelectItem(sets, action)
		if resultindex > -1 {
			set = sets[resultindex]
		} else {
			return -1
		}
	} else {
		set = sets[0]
	}
	return tbl.GetIdx(set.ID)
}

// EditSet opens the edit UI for a selected set.
func EditSet[T data.CustomData](tbl *data.DataTable[T], setname string, tablename string) {
	idx := SelectSet(tbl, setname, tablename, "edit")
	if idx < 0 {
		return
	}
	set, saved := ui.EditItem((*tbl)[idx], GetDropDownOpts((*tbl)[idx].ID), " Edit ")
	if saved {
		(*tbl)[idx] = set
		data.Db.Save()
	}
}

// DeleteSet deletes a selected set after confirmation.
func DeleteSet[T data.CustomData](tbl *data.DataTable[T], setname string, tablename string) {
	idx := SelectSet(tbl, setname, tablename, "delete")
	if idx < 0 {
		return
	}

	DeleteSetId(tbl, idx, tablename)
}

// DeleteSetId deletes a set by index after confirmation.
func DeleteSetId[T data.CustomData](tbl *data.DataTable[T], idx int, tablename string) {
	if ui.Alert(fmt.Sprintf("Do you really want to delete %s \"%s\"?", strings.ToLower(tablename), (*tbl)[idx].Name)) {
		tbl.Delete((*tbl)[idx].ID)
		data.Db.Save()
		fmt.Printf("%s \"%s\" deleted\n", tablename, (*tbl)[idx].Name)
	}
}

// ShowSet prints details for the selected set or id.
func ShowSet[T data.CustomData](tbl *data.DataTable[T], setname string, tablename string) {
	sets := tbl.GetSetsByName(setname)
	if len(sets) == 0 {
		id, err := strconv.ParseUint(setname, 10, 32)
		if err != nil {
			fmt.Printf("No %s with name \"%s\" found.\n", strings.ToLower(tablename), setname)
			return
		}
		idset, ok := tbl.GetPtr(uint32(id))
		if !ok {
			fmt.Printf("No %s with ID %s found.\n", strings.ToLower(tablename), setname)
			return
		}
		(*idset).Show()
	} else {

		for _, set := range sets {
			set.Show()
		}
	}
}

// AddTagById attaches a tag to a dataset at the given index.
func AddTagById[T data.CustomData](tbl *data.DataTable[T], idx int, tagid uint32, tagname string) bool {
	objname := strings.Split(fmt.Sprintf("%T", *new(T)), ".")
	if slices.Contains((*tbl)[idx].Tags, tagid) {
		fmt.Printf("%s \"%s\" is already tagged with \"%s\"", objname[1], (*tbl)[idx].Name, tagname)
		return false
	}
	(*tbl)[idx].Tags = append((*tbl)[idx].Tags, tagid)
	(*tbl)[idx].Updated = time.Now().Unix()
	fmt.Printf("Added Tag \"%s\" to %s \"%s\"", tagname, objname[1], (*tbl)[idx].Name)
	data.Db.Save()
	return true
}

// RemoveTagById removes a tag from a dataset at the given index.
func RemoveTagById[T data.CustomData](tbl *data.DataTable[T], idx int, tagid uint32, tagname string) bool {
	objname := strings.Split(fmt.Sprintf("%T", *new(T)), ".")
	if !slices.Contains((*tbl)[idx].Tags, tagid) {
		fmt.Printf("%s \"%s\" is not tagged with \"%s\"", objname[1], (*tbl)[idx].Name, tagname)
		return false
	}
	// copy only tags that are not of the id to remove, preserving order
	nt := []uint32{}
	for _, tid := range (*tbl)[idx].Tags {
		if tid != tagid {
			nt = append(nt, tid)
		}
	}

	(*tbl)[idx].Tags = nt
	(*tbl)[idx].Updated = time.Now().Unix()
	fmt.Printf("Removed Tag \"%s\" from %s \"%s\"", tagname, objname[1], (*tbl)[idx].Name)
	data.Db.Save()
	return true
}

// AddTag adds a tag to the object identified by id.
func AddTag(tagname string, id uint32) {
	tagid, found := data.Db.Tags.GetFirstOccurance(tagname)
	if !found {
		tagid = data.Db.Tags.AddSimple(tagname)
	}

	kind, idx := findTableById(id)
	switch kind {
	case kindCategory:
		AddTagById(&data.Db.Categories, idx, tagid, tagname)
	case kindWarehouse:
		AddTagById(&data.Db.Warehouses, idx, tagid, tagname)
	case kindRoom:
		AddTagById(&data.Db.Rooms, idx, tagid, tagname)
	case kindShelf:
		AddTagById(&data.Db.Shelves, idx, tagid, tagname)
	case kindBox:
		AddTagById(&data.Db.Boxes, idx, tagid, tagname)
	case kindItem:
		AddTagById(&data.Db.Items, idx, tagid, tagname)
	default:
		fmt.Println("No record found.")
	}
}

// RemoveTag removes a tag from the object identified by id.
func RemoveTag(tagname string, id uint32) {
	tagid, found := data.Db.Tags.GetFirstOccurance(tagname)
	if !found {
		fmt.Printf("Tag \"%s\" not found\n", tagname)
		return
	}
	kind, idx := findTableById(id)
	switch kind {
	case kindCategory:
		RemoveTagById(&data.Db.Categories, idx, tagid, tagname)
	case kindWarehouse:
		RemoveTagById(&data.Db.Warehouses, idx, tagid, tagname)
	case kindRoom:
		RemoveTagById(&data.Db.Rooms, idx, tagid, tagname)
	case kindShelf:
		RemoveTagById(&data.Db.Shelves, idx, tagid, tagname)
	case kindBox:
		RemoveTagById(&data.Db.Boxes, idx, tagid, tagname)
	case kindItem:
		RemoveTagById(&data.Db.Items, idx, tagid, tagname)
	default:
		fmt.Println("No record found.")
	}
}

// ShowAny prints details for any object by id.
func ShowAny(id uint32) {
	kind, idx := findTableById(id)
	switch kind {
	case kindCategory:
		data.Db.Categories[idx].Show()
	case kindWarehouse:
		data.Db.Warehouses[idx].Show()
	case kindRoom:
		data.Db.Rooms[idx].Show()
	case kindShelf:
		data.Db.Shelves[idx].Show()
	case kindBox:
		data.Db.Boxes[idx].Show()
	case kindItem:
		data.Db.Items[idx].Show()
	default:
		fmt.Println("No record found.")
	}
}

// DeleteAny deletes any object by id after confirmation.
func DeleteAny(id uint32) {
	kind, idx := findTableById(id)
	switch kind {
	case kindCategory:
		DeleteSetId(&data.Db.Categories, idx, tableName(kind))
	case kindWarehouse:
		DeleteSetId(&data.Db.Warehouses, idx, tableName(kind))
	case kindRoom:
		DeleteSetId(&data.Db.Rooms, idx, tableName(kind))
	case kindShelf:
		DeleteSetId(&data.Db.Shelves, idx, tableName(kind))
	case kindBox:
		DeleteSetId(&data.Db.Boxes, idx, tableName(kind))
	case kindItem:
		DeleteSetId(&data.Db.Items, idx, tableName(kind))
	default:
		fmt.Println("No record found.")
	}
}

// EditAny opens the edit UI for any object by id.
func EditAny(id uint32) {
	kind, idx := findTableById(id)
	switch kind {
	case kindCategory:
		set, saved := ui.EditItem(data.Db.Categories[idx], GetDropDownOpts(data.Db.Categories[idx].ID), " Edit ")
		if saved {
			data.Db.Categories[idx] = set
			data.Db.Save()
		}
	case kindWarehouse:
		set, saved := ui.EditItem(data.Db.Warehouses[idx], GetDropDownOpts(data.Db.Warehouses[idx].ID), " Edit ")
		if saved {
			data.Db.Warehouses[idx] = set
			data.Db.Save()
		}
	case kindRoom:
		set, saved := ui.EditItem(data.Db.Rooms[idx], GetDropDownOpts(data.Db.Rooms[idx].ID), " Edit ")
		if saved {
			data.Db.Rooms[idx] = set
			data.Db.Save()
		}
	case kindShelf:
		set, saved := ui.EditItem(data.Db.Shelves[idx], GetDropDownOpts(data.Db.Shelves[idx].ID), " Edit ")
		if saved {
			data.Db.Shelves[idx] = set
			data.Db.Save()
		}
	case kindBox:
		set, saved := ui.EditItem(data.Db.Boxes[idx], GetDropDownOpts(data.Db.Boxes[idx].ID), " Edit ")
		if saved {
			data.Db.Boxes[idx] = set
			data.Db.Save()
		}
	case kindItem:
		set, saved := ui.EditItem(data.Db.Items[idx], GetDropDownOpts(data.Db.Items[idx].ID), " Edit ")
		if saved {
			data.Db.Items[idx] = set
			data.Db.Save()
		}
	default:
		fmt.Println("No record found.")
	}
}

// PrintItemsOfCategory lists items in a category.
func PrintItemsOfCategory(catnameOrId string, sortname bool) {
	idx := SelectSet(&data.Db.Categories, catnameOrId, "Category", "show")
	if idx < 0 {
		return
	}
	data.Db.Items.PrintListFiltered(sortname, func(set data.Dataset[data.Item]) bool {
		return set.Data.CategoryId == data.Db.Categories[idx].ID
	})
}

// PrintItemsOfBox lists items in a box.
func PrintItemsOfBox(boxnameOrId string, sortname bool) {
	idx := SelectSet(&data.Db.Boxes, boxnameOrId, "Box", "show")
	if idx < 0 {
		return
	}
	data.Db.Items.PrintListFiltered(sortname, func(set data.Dataset[data.Item]) bool {
		return set.Data.BoxId == data.Db.Boxes[idx].ID
	})
}
