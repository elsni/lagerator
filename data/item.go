package data

import (
	"fmt"

	"github.com/elsni/lagerator/terminal"
)

type Item struct {
	Location   string `json:"location"`
	Condition  string `json:"condition"`
	Amount     int    `json:"amount"`
	BoxId      uint32 `json:"boxId"`
	CategoryId uint32 `json:"categoryId"`
}

// GetTableRow returns the formatted row for an item.
func (d Item) GetTableRow(ownid uint32) string {
	boxName := GetPrintNameById(&Db.Boxes, d.BoxId, 15)
	shelfId := GetShelfIdforBox(d.BoxId)
	shelfName := GetPrintNameById(&Db.Shelves, shelfId, 20)
	roomId := GetRoomIdforShelf(shelfId)
	roomName := GetPrintNameById(&Db.Rooms, roomId, 20)
	whName := GetPrintNameById(&Db.Warehouses, GetWarehouseIdforRoom(roomId), 20)
	catName := GetPrintNameById(&Db.Categories, d.CategoryId, 15)
	return fmt.Sprintf("%-5d %-15s %-15s %-20s %-20s %-20s", d.Amount, catName, boxName, shelfName, roomName, whName)
}

// GetTableHeader returns the item table header.
func (d Item) GetTableHeader() string {
	return fmt.Sprintf("%-5s %-15s %-15s %-20s %-20s %-20s%s", "Amnt", "Category", "Box", "Shelf", "Room", "Warehouse", terminal.ResetColor())
}

// GetShelfIdforItem returns the shelf id for an item id.
func GetShelfIdforItem(itemid uint32) uint32 {
	item, ok := Db.Items.GetPtr(itemid)
	if !ok {
		return 0
	}
	return GetShelfIdforBox(item.Data.BoxId)
}

type ItemTable = DataTable[Item]

// Show prints item details.
func (d Item) Show() {
	fmt.Printf("%s %s\n", terminal.GetLabelText("Location"), d.Location)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Condition"), d.Condition)
	fmt.Printf("%s %d\n", terminal.GetLabelText("Amount"), d.Amount)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Box"), GetPrintNameById(&Db.Boxes, d.BoxId, 999))
	fmt.Printf("%s %s\n", terminal.GetLabelText("Category"), GetPrintNameById(&Db.Categories, d.CategoryId, 999))
}
