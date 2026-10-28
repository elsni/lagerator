package logic

import (
	"fmt"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/presentation"
	"github.com/elsni/lagerator/ui"
)

// ShowAny prints details for any object by id.
func ShowAny(id uint32) {
	kind, idx := findTableById(id)
	switch kind {
	case kindCategory:
		presentation.Show(data.Db.Categories[idx])
	case kindWarehouse:
		presentation.Show(data.Db.Warehouses[idx])
	case kindRoom:
		presentation.Show(data.Db.Rooms[idx])
	case kindShelf:
		presentation.Show(data.Db.Shelves[idx])
	case kindBox:
		presentation.Show(data.Db.Boxes[idx])
	case kindItem:
		presentation.Show(data.Db.Items[idx])
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
	presentation.PrintListFiltered(&data.Db.Items, sortname, func(set data.Dataset[data.Item]) bool {
		return set.Data.CategoryId == data.Db.Categories[idx].ID
	})
}

// PrintItemsOfBox lists items in a box.
func PrintItemsOfBox(boxnameOrId string, sortname bool) {
	idx := SelectSet(&data.Db.Boxes, boxnameOrId, "Box", "show")
	if idx < 0 {
		return
	}
	presentation.PrintListFiltered(&data.Db.Items, sortname, func(set data.Dataset[data.Item]) bool {
		return set.Data.BoxId == data.Db.Boxes[idx].ID
	})
}
