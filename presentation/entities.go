package presentation

import (
	"fmt"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/terminal"
)

// tableRowWarehouse renders the warehouse GetTableRow output.
func tableRowWarehouse(d data.Warehouse, ownid uint32) string {
	return ""
}

// tableHeaderWarehouse renders the warehouse GetTableHeader output.
func tableHeaderWarehouse(d data.Warehouse) string {
	return terminal.ResetColor()
}

// showWarehouse renders the warehouse Show output.
func showWarehouse(d data.Warehouse) {
}

// tableRowRoom renders the room GetTableRow output.
func tableRowRoom(d data.Room, ownid uint32) string {
	return fmt.Sprintf("%-20s", GetPrintNameById(&data.Db.Warehouses, d.WarehouseId, 20))
}

// tableHeaderRoom renders the room GetTableHeader output.
func tableHeaderRoom(d data.Room) string {
	return fmt.Sprintf("%-20s%s", "Warehouse", terminal.ResetColor())
}

// showRoom renders the room Show output.
func showRoom(d data.Room) {
	fmt.Printf("%s %s\n", terminal.GetLabelText("Location"), d.Location)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Warehouse"), GetPrintNameById(&data.Db.Warehouses, d.WarehouseId, 20))
}

// tableRowShelf renders the shelf GetTableRow output.
func tableRowShelf(d data.Shelf, ownid uint32) string {
	roomName := GetPrintNameById(&data.Db.Rooms, d.RoomId, 20)
	whName := GetPrintNameById(&data.Db.Warehouses, data.GetWarehouseIdforRoom(d.RoomId), 20)
	return fmt.Sprintf("%-20s %-20s", roomName, whName)
}

// tableHeaderShelf renders the shelf GetTableHeader output.
func tableHeaderShelf(d data.Shelf) string {
	return fmt.Sprintf("%-20s %-20s%s", "Room", "Warehouse", terminal.ResetColor())
}

// showShelf renders the shelf Show output.
func showShelf(d data.Shelf) {
	fmt.Printf("%s %s\n", terminal.GetLabelText("Location"), d.Location)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Room"), GetPrintNameById(&data.Db.Rooms, d.RoomId, 999))
}

// tableRowBox renders the box GetTableRow output.
func tableRowBox(d data.Box, ownid uint32) string {
	shelfName := GetPrintNameById(&data.Db.Shelves, d.ShelfId, 20)
	roomId := data.GetRoomIdforShelf(d.ShelfId)
	roomName := GetPrintNameById(&data.Db.Rooms, roomId, 20)
	whName := GetPrintNameById(&data.Db.Warehouses, data.GetWarehouseIdforRoom(roomId), 20)
	return fmt.Sprintf("%-20s %-20s %-20s", shelfName, roomName, whName)
}

// tableHeaderBox renders the box GetTableHeader output.
func tableHeaderBox(d data.Box) string {
	return fmt.Sprintf("%-20s %-20s %-20s%s", "Shelf", "Room", "Warehouse", terminal.ResetColor())
}

// showBox renders the box Show output.
func showBox(d data.Box) {
	fmt.Printf("%s %s\n", terminal.GetLabelText("Location"), d.Location)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Type"), d.Type)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Shelf"), GetPrintNameById(&data.Db.Shelves, d.ShelfId, 20))
}

// tableRowItem renders the item GetTableRow output.
func tableRowItem(d data.Item, ownid uint32) string {
	boxName := GetPrintNameById(&data.Db.Boxes, d.BoxId, 15)
	shelfId := data.GetShelfIdforBox(d.BoxId)
	shelfName := GetPrintNameById(&data.Db.Shelves, shelfId, 20)
	roomId := data.GetRoomIdforShelf(shelfId)
	roomName := GetPrintNameById(&data.Db.Rooms, roomId, 20)
	whName := GetPrintNameById(&data.Db.Warehouses, data.GetWarehouseIdforRoom(roomId), 20)
	catName := GetPrintNameById(&data.Db.Categories, d.CategoryId, 15)
	return fmt.Sprintf("%-5d %-15s %-15s %-20s %-20s %-20s", d.Amount, catName, boxName, shelfName, roomName, whName)
}

// tableHeaderItem renders the item GetTableHeader output.
func tableHeaderItem(d data.Item) string {
	return fmt.Sprintf("%-5s %-15s %-15s %-20s %-20s %-20s%s", "Amnt", "Category", "Box", "Shelf", "Room", "Warehouse", terminal.ResetColor())
}

// showItem renders the item Show output.
func showItem(d data.Item) {
	fmt.Printf("%s %s\n", terminal.GetLabelText("Location"), d.Location)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Condition"), d.Condition)
	fmt.Printf("%s %d\n", terminal.GetLabelText("Amount"), d.Amount)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Box"), GetPrintNameById(&data.Db.Boxes, d.BoxId, 999))
	fmt.Printf("%s %s\n", terminal.GetLabelText("Category"), GetPrintNameById(&data.Db.Categories, d.CategoryId, 999))
}

// tableRowCategory renders the category GetTableRow output.
func tableRowCategory(d data.Category, ownid uint32) string {
	return ""
}

// tableHeaderCategory renders the category GetTableHeader output.
func tableHeaderCategory(d data.Category) string {
	return terminal.ResetColor()
}

// showCategory renders the category Show output.
func showCategory(d data.Category) {
}

// tableRowTag renders the tag GetTableRow output.
func tableRowTag(d data.Tag, ownid uint32) string {
	return fmt.Sprintf("%4d", data.Db.CountTagOccurance(ownid))
}

// tableHeaderTag renders the tag GetTableHeader output.
func tableHeaderTag(d data.Tag) string {
	return fmt.Sprintf("%-15s%s", "uses", terminal.ResetColor())
}

// showTag renders the tag Show output.
func showTag(d data.Tag) {
}
