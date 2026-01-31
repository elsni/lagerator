package logic

import (
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/id"
)

// TestMain sets a temporary HOME so tests don't touch the real database file.
func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "lgrttest-*")
	if err != nil {
		panic(err)
	}
	originalHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempDir)
	code := m.Run()
	if originalHome == "" {
		_ = os.Unsetenv("HOME")
	} else {
		_ = os.Setenv("HOME", originalHome)
	}
	_ = os.RemoveAll(tempDir)
	os.Exit(code)
}

// resetDb resets the global database and id source.
func resetDb() {
	id.IdSource.SetLastId(0)
	data.Db = data.NewDatabase()
}

// captureOutput captures stdout for the duration of fn.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	_ = r.Close()
	return string(out)
}

// TestCurrentWarehouseExists verifies the current warehouse validation.
func TestCurrentWarehouseExists(t *testing.T) {
	resetDb()
	if CurrentWarehouseExists() {
		t.Fatalf("expected CurrentWarehouseExists to be false for empty db")
	}

	wh := data.NewDataset[data.Warehouse]("WH", data.Warehouse{})
	data.Db.Warehouses.Add(wh)
	data.Db.CurrentWarehouse = wh.ID
	if !CurrentWarehouseExists() {
		t.Fatalf("expected CurrentWarehouseExists to be true for valid warehouse")
	}
}

// TestSwitchWarehouseMissing verifies switching to an unknown warehouse.
func TestSwitchWarehouseMissing(t *testing.T) {
	resetDb()
	out := captureOutput(t, func() {
		SwitchWarehouse("unknown")
	})
	if !strings.Contains(out, "does not exist") {
		t.Fatalf("expected missing warehouse message, got: %s", out)
	}
}

// TestSwitchWarehousePositive verifies switching to an existing warehouse.
func TestSwitchWarehousePositive(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	if data.Db.CurrentWarehouse == 0 {
		t.Fatalf("expected current warehouse to be set")
	}
}

// TestAddWarehouseDuplicate verifies duplicate warehouses are rejected.
func TestAddWarehouseDuplicate(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	AddWarehouse("WH1")
	if len(data.Db.Warehouses) != 1 {
		t.Fatalf("expected 1 warehouse, got %d", len(data.Db.Warehouses))
	}
}

// TestAddCategoryPositive verifies categories are created.
func TestAddCategoryPositive(t *testing.T) {
	resetDb()
	AddCategory("Tools")
	if len(data.Db.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(data.Db.Categories))
	}
}

// TestAddRoomWithoutWarehouse verifies adding rooms without active warehouse.
func TestAddRoomWithoutWarehouse(t *testing.T) {
	resetDb()
	out := captureOutput(t, func() {
		AddRoomToCurrentWarehouse("room")
	})
	if !strings.Contains(out, "Switch to valid warehouse first") {
		t.Fatalf("expected warning for missing warehouse, got: %s", out)
	}
}

// TestAddRoomShelfBoxFlow verifies room/shelf/box creation path.
func TestAddRoomShelfBoxFlow(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")

	AddRoomToCurrentWarehouse("R1")
	if len(data.Db.Rooms) != 1 {
		t.Fatalf("expected 1 room, got %d", len(data.Db.Rooms))
	}

	AddShelfToRoom("S1", "R1")
	if len(data.Db.Shelves) != 1 {
		t.Fatalf("expected 1 shelf, got %d", len(data.Db.Shelves))
	}

	AddBoxToShelf("B1", "S1")
	if len(data.Db.Boxes) != 1 {
		t.Fatalf("expected 1 box, got %d", len(data.Db.Boxes))
	}
}

// TestMoveItemPositive verifies item moves between boxes.
func TestMoveItemPositive(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	AddRoomToCurrentWarehouse("R1")
	AddShelfToRoom("S1", "R1")
	AddBoxToShelf("B1", "S1")
	AddBoxToShelf("B2", "S1")

	box1 := data.Db.Boxes[0]
	box2 := data.Db.Boxes[1]
	item := data.NewDataset[data.Item]("I1", data.Item{BoxId: box1.ID})
	data.Db.Items.Add(item)

	MoveItem(item.ID, strconv.FormatUint(uint64(box2.ID), 10))
	if data.Db.Items[0].Data.BoxId != box2.ID {
		t.Fatalf("expected item to move to box %d", box2.ID)
	}
}

// TestMoveItemInvalid verifies invalid item id is rejected.
func TestMoveItemInvalid(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	AddRoomToCurrentWarehouse("R1")
	AddShelfToRoom("S1", "R1")
	AddBoxToShelf("B1", "S1")

	out := captureOutput(t, func() {
		MoveItem(999, "B1")
	})
	if !strings.Contains(out, "No item with id") {
		t.Fatalf("expected missing item message, got: %s", out)
	}
}

// TestMoveBoxPositive verifies box moves between shelves.
func TestMoveBoxPositive(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	AddRoomToCurrentWarehouse("R1")
	AddShelfToRoom("S1", "R1")
	AddShelfToRoom("S2", "R1")
	AddBoxToShelf("B1", "S1")

	box := data.Db.Boxes[0]
	shelf2 := data.Db.Shelves[1]
	MoveBox(box.ID, strconv.FormatUint(uint64(shelf2.ID), 10))
	if data.Db.Boxes[0].Data.ShelfId != shelf2.ID {
		t.Fatalf("expected box to move to shelf %d", shelf2.ID)
	}
}

// TestMoveBoxInvalid verifies invalid box id is rejected.
func TestMoveBoxInvalid(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	AddRoomToCurrentWarehouse("R1")
	AddShelfToRoom("S1", "R1")

	out := captureOutput(t, func() {
		MoveBox(999, "S1")
	})
	if !strings.Contains(out, "No box with id") {
		t.Fatalf("expected missing box message, got: %s", out)
	}
}

// TestDropDownOpts verifies dropdown options for each object type.
func TestDropDownOpts(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	AddRoomToCurrentWarehouse("R1")
	AddShelfToRoom("S1", "R1")
	AddBoxToShelf("B1", "S1")
	AddCategory("C1")
	item := data.NewDataset[data.Item]("I1", data.Item{BoxId: data.Db.Boxes[0].ID})
	data.Db.Items.Add(item)

	roomID := data.Db.Rooms[0].ID
	shelfID := data.Db.Shelves[0].ID
	boxID := data.Db.Boxes[0].ID
	itemID := data.Db.Items[0].ID

	opts := GetDropDownOpts(roomID)
	if len(opts) != 2 || len(opts[0]) != 1 || opts[0][0].Name != "WH1" {
		t.Fatalf("unexpected room dropdown opts: %+v", opts)
	}

	opts = GetDropDownOpts(shelfID)
	if len(opts[0]) != 1 || opts[0][0].Name != "R1" {
		t.Fatalf("unexpected shelf dropdown opts: %+v", opts)
	}

	opts = GetDropDownOpts(boxID)
	if len(opts[0]) != 1 || opts[0][0].Name != "S1" {
		t.Fatalf("unexpected box dropdown opts: %+v", opts)
	}

	opts = GetDropDownOpts(itemID)
	if len(opts[0]) != 1 || opts[0][0].Name != "B1" {
		t.Fatalf("unexpected item dropdown opts: %+v", opts)
	}
	if len(opts[1]) != 1 || opts[1][0].Name != "C1" {
		t.Fatalf("unexpected category dropdown opts: %+v", opts)
	}
}

// TestSelectSet verifies selecting by name and id, and error cases.
func TestSelectSet(t *testing.T) {
	resetDb()
	cat := data.NewDataset[data.Category]("Tools", data.Category{})
	data.Db.Categories.Add(cat)

	if idx := SelectSet(&data.Db.Categories, "Tools", "Category", "show"); idx != 0 {
		t.Fatalf("expected index 0, got %d", idx)
	}
	if idx := SelectSet(&data.Db.Categories, strconv.FormatUint(uint64(cat.ID), 10), "Category", "show"); idx != 0 {
		t.Fatalf("expected index 0 by id, got %d", idx)
	}

	out := captureOutput(t, func() {
		SelectSet(&data.Db.Categories, "missing", "Category", "show")
	})
	if !strings.Contains(out, "No category with name") {
		t.Fatalf("expected missing name message, got: %s", out)
	}

	out = captureOutput(t, func() {
		SelectSet(&data.Db.Categories, "999", "Category", "show")
	})
	if !strings.Contains(out, "No category with ID") {
		t.Fatalf("expected missing id message, got: %s", out)
	}
}

// TestShowSet verifies output for name and id.
func TestShowSet(t *testing.T) {
	resetDb()
	cat := data.NewDataset[data.Category]("Tools", data.Category{})
	data.Db.Categories.Add(cat)

	out := captureOutput(t, func() {
		ShowSet(&data.Db.Categories, "Tools", "Category")
	})
	if !strings.Contains(out, "Tools") {
		t.Fatalf("expected output to include name, got: %s", out)
	}

	out = captureOutput(t, func() {
		ShowSet(&data.Db.Categories, strconv.FormatUint(uint64(cat.ID), 10), "Category")
	})
	if !strings.Contains(out, "Tools") {
		t.Fatalf("expected output to include name by id, got: %s", out)
	}

	out = captureOutput(t, func() {
		ShowSet(&data.Db.Categories, "missing", "Category")
	})
	if !strings.Contains(out, "No category with name") {
		t.Fatalf("expected missing name message, got: %s", out)
	}
}

// TestDeleteEditSetInvalid verifies EditSet/DeleteSet return on invalid names.
func TestDeleteEditSetInvalid(t *testing.T) {
	resetDb()
	out := captureOutput(t, func() {
		DeleteSet(&data.Db.Categories, "missing", "Category")
	})
	if !strings.Contains(out, "No category with name") {
		t.Fatalf("expected missing name message, got: %s", out)
	}

	out = captureOutput(t, func() {
		EditSet(&data.Db.Categories, "missing", "Category")
	})
	if !strings.Contains(out, "No category with name") {
		t.Fatalf("expected missing name message, got: %s", out)
	}
}

// TestAddRemoveTag verifies tag add/remove behavior.
func TestAddRemoveTag(t *testing.T) {
	resetDb()
	cat := data.NewDataset[data.Category]("Tools", data.Category{})
	data.Db.Categories.Add(cat)

	AddTag("fragile", cat.ID)
	tagID, ok := data.Db.Tags.GetFirstOccurance("fragile")
	if !ok {
		t.Fatalf("expected tag to be created")
	}
	if !strings.Contains(data.GetTagList(data.Db.Categories[0].Tags), "fragile") {
		t.Fatalf("expected category to have fragile tag")
	}

	RemoveTag("fragile", cat.ID)
	if strings.Contains(data.GetTagList(data.Db.Categories[0].Tags), "fragile") {
		t.Fatalf("expected tag to be removed")
	}

	// remove unknown tag
	out := captureOutput(t, func() {
		RemoveTag("missing", cat.ID)
	})
	if !strings.Contains(out, "not found") {
		t.Fatalf("expected missing tag message, got: %s", out)
	}

	// add tag to unknown id still creates tag, but reports no record found
	out = captureOutput(t, func() {
		AddTag("orphan", 999)
	})
	if !strings.Contains(out, "No record found") {
		t.Fatalf("expected no record found, got: %s", out)
	}
	if _, ok := data.Db.Tags.GetFirstOccurance("orphan"); !ok {
		t.Fatalf("expected orphan tag to be created")
	}

	// remove tag from unknown id
	out = captureOutput(t, func() {
		RemoveTag("orphan", 999)
	})
	if !strings.Contains(out, "No record found") {
		t.Fatalf("expected no record found on remove, got: %s", out)
	}
	if tagID == 0 {
		t.Fatalf("unexpected tag id state")
	}
}

// TestPrintItemsOfCategory verifies positive and negative listings.
func TestPrintItemsOfCategory(t *testing.T) {
	resetDb()
	cat := data.NewDataset[data.Category]("Tools", data.Category{})
	data.Db.Categories.Add(cat)
	item := data.NewDataset[data.Item]("I1", data.Item{CategoryId: cat.ID})
	data.Db.Items.Add(item)

	out := captureOutput(t, func() {
		PrintItemsOfCategory("Tools", false)
	})
	if !strings.Contains(out, "I1") {
		t.Fatalf("expected item in output, got: %s", out)
	}

	out = captureOutput(t, func() {
		PrintItemsOfCategory("missing", false)
	})
	if !strings.Contains(out, "No category with name") {
		t.Fatalf("expected missing category message, got: %s", out)
	}
}

// TestPrintItemsOfBox verifies positive and negative listings.
func TestPrintItemsOfBox(t *testing.T) {
	resetDb()
	AddWarehouse("WH1")
	SwitchWarehouse("WH1")
	AddRoomToCurrentWarehouse("R1")
	AddShelfToRoom("S1", "R1")
	AddBoxToShelf("B1", "S1")
	box := data.Db.Boxes[0]
	item := data.NewDataset[data.Item]("I1", data.Item{BoxId: box.ID})
	data.Db.Items.Add(item)

	out := captureOutput(t, func() {
		PrintItemsOfBox("B1", false)
	})
	if !strings.Contains(out, "I1") {
		t.Fatalf("expected item in output, got: %s", out)
	}

	out = captureOutput(t, func() {
		PrintItemsOfBox("missing", false)
	})
	if !strings.Contains(out, "No box with name") {
		t.Fatalf("expected missing box message, got: %s", out)
	}
}

// TestShowAnyUnknown verifies message for missing id.
func TestShowAnyUnknown(t *testing.T) {
	resetDb()
	out := captureOutput(t, func() {
		ShowAny(12345)
	})
	if !strings.Contains(out, "No record found") {
		t.Fatalf("expected no record found message, got: %s", out)
	}
}

// TestDeleteAnyUnknown verifies delete on unknown id.
func TestDeleteAnyUnknown(t *testing.T) {
	resetDb()
	out := captureOutput(t, func() {
		DeleteAny(12345)
	})
	if !strings.Contains(out, "No record found") {
		t.Fatalf("expected no record found message, got: %s", out)
	}
}
