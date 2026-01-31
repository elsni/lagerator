package data

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	Db = NewDatabase()
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

// dbFilePath returns the path to the database file under the current HOME.
func dbFilePath(t *testing.T) string {
	t.Helper()
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".lgrt", "lgrtdata.json")
}

// TestDatabaseSaveLoad verifies Save/Load and id restoration using a temp HOME.
func TestDatabaseSaveLoad(t *testing.T) {
	resetDb()

	wh := NewDataset[Warehouse]("WH1", Warehouse{})
	Db.Warehouses.Add(wh)
	Db.CurrentWarehouse = wh.ID

	tagID := Db.Tags.AddSimple("fragile")
	cat := NewDataset[Category]("Tools", Category{})
	cat.Tags = []uint32{tagID}
	Db.Categories.Add(cat)

	item := NewDataset[Item]("Hammer", Item{Amount: 2, Location: "A"})
	item.Tags = []uint32{tagID}
	Db.Items.Add(item)

	Db.Save()

	dbPath := dbFilePath(t)
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected db file at %s: %v", dbPath, err)
	}

	db2 := NewDatabase()
	db2.Load()

	if len(db2.Warehouses) != 1 || db2.Warehouses[0].Name != "WH1" {
		t.Fatalf("warehouse not loaded correctly: %+v", db2.Warehouses)
	}
	if db2.CurrentWarehouse != wh.ID {
		t.Fatalf("expected current warehouse %d, got %d", wh.ID, db2.CurrentWarehouse)
	}
	if len(db2.Tags) != 1 || db2.Tags[0].Name != "fragile" {
		t.Fatalf("tags not loaded correctly: %+v", db2.Tags)
	}
	if id.IdSource.LastId != db2.FindLastId() {
		t.Fatalf("id source not restored: got %d, want %d", id.IdSource.LastId, db2.FindLastId())
	}
}

// TestDatabaseLoadMissingFile ensures Load is a no-op when no file exists.
func TestDatabaseLoadMissingFile(t *testing.T) {
	resetDb()
	_ = os.Remove(dbFilePath(t))
	Db.Load()
	if len(Db.Warehouses) != 0 || len(Db.Rooms) != 0 || len(Db.Items) != 0 {
		t.Fatalf("expected empty database after Load without file")
	}
}

// TestCountTagOccurance verifies counting tag usage across tables.
func TestCountTagOccurance(t *testing.T) {
	resetDb()
	tagID := Db.Tags.AddSimple("fragile")

	cat := NewDataset[Category]("Tools", Category{})
	cat.Tags = []uint32{tagID}
	Db.Categories.Add(cat)

	wh := NewDataset[Warehouse]("WH", Warehouse{})
	wh.Tags = []uint32{tagID}
	Db.Warehouses.Add(wh)

	room := NewDataset[Room]("R1", Room{WarehouseId: wh.ID})
	room.Tags = []uint32{tagID}
	Db.Rooms.Add(room)

	shelf := NewDataset[Shelf]("S1", Shelf{RoomId: room.ID})
	shelf.Tags = []uint32{tagID}
	Db.Shelves.Add(shelf)

	box := NewDataset[Box]("B1", Box{ShelfId: shelf.ID})
	box.Tags = []uint32{tagID}
	Db.Boxes.Add(box)

	item := NewDataset[Item]("I1", Item{BoxId: box.ID})
	item.Tags = []uint32{tagID}
	Db.Items.Add(item)

	if got := Db.CountTagOccurance(tagID); got != 6 {
		t.Fatalf("expected 6 tag uses, got %d", got)
	}
}

// TestCountTagOccuranceMissing verifies that unknown tags return zero.
func TestCountTagOccuranceMissing(t *testing.T) {
	resetDb()
	if got := Db.CountTagOccurance(999); got != 0 {
		t.Fatalf("expected zero tag uses for unknown tag, got %d", got)
	}
}

// TestGetTagListAndIds verifies tag list formatting and id resolution.
func TestGetTagListAndIds(t *testing.T) {
	resetDb()
	tagOne := Db.Tags.AddSimple("one")
	tagTwo := Db.Tags.AddSimple("two")

	list := GetTagList([]uint32{tagOne, tagTwo})
	if list != "one, two" {
		t.Fatalf("unexpected tag list: %q", list)
	}

	ids := GetTagIds("one, three")
	if len(ids) != 2 {
		t.Fatalf("expected 2 tag ids, got %d", len(ids))
	}
	if ids[0] != tagOne {
		t.Fatalf("expected existing tag id %d, got %d", tagOne, ids[0])
	}
	if _, ok := Db.Tags.GetFirstOccurance("three"); !ok {
		t.Fatalf("expected tag \"three\" to be created")
	}
}

// TestGetTagListEmpty verifies empty inputs return empty strings.
func TestGetTagListEmpty(t *testing.T) {
	resetDb()
	if list := GetTagList(nil); list != "" {
		t.Fatalf("expected empty list, got %q", list)
	}
	if list := GetTagList([]uint32{}); list != "" {
		t.Fatalf("expected empty list, got %q", list)
	}
}

// TestGetTagIdsEmpty verifies blank tag lists return no ids.
func TestGetTagIdsEmpty(t *testing.T) {
	resetDb()
	ids := GetTagIds("  , , ")
	if len(ids) != 0 {
		t.Fatalf("expected empty ids, got %v", ids)
	}
}

// TestFindItem verifies search output for items.
func TestFindItem(t *testing.T) {
	resetDb()
	item := NewDataset[Item]("Hammer", Item{Amount: 1, Location: "Shelf A"})
	item.Description = "Steel tool"
	Db.Items.Add(item)
	Db.Items.Add(NewDataset[Item]("Nails", Item{Amount: 50, Location: "Shelf B"}))

	out := captureOutput(t, func() {
		Db.FindItem("hammer", false)
	})
	if !strings.Contains(out, "Hammer") {
		t.Fatalf("expected output to include Hammer, got: %s", out)
	}
	if strings.Contains(out, "Nails") {
		t.Fatalf("did not expect Nails in output: %s", out)
	}
}

// TestFindItemNoData verifies the no data output.
func TestFindItemNoData(t *testing.T) {
	resetDb()
	out := captureOutput(t, func() {
		Db.FindItem("anything", false)
	})
	if !strings.Contains(out, "no data") {
		t.Fatalf("expected no data output, got: %s", out)
	}
}

// TestFindItemNoMatch verifies no matches are printed.
func TestFindItemNoMatch(t *testing.T) {
	resetDb()
	Db.Items.Add(NewDataset[Item]("Hammer", Item{Amount: 1, Location: "Shelf A"}))
	out := captureOutput(t, func() {
		Db.FindItem("wrench", false)
	})
	if strings.Contains(out, "Hammer") {
		t.Fatalf("did not expect Hammer in output: %s", out)
	}
}

// TestDataTableOperations verifies core DataTable methods.
func TestDataTableOperations(t *testing.T) {
	id.IdSource.SetLastId(0)
	dt := NewDataTable[Category]()

	setA := NewDataset[Category]("Alpha", Category{})
	dt.Add(setA)
	setB := NewDataset[Category]("Bravo", Category{})
	dt.Add(setB)

	names := dt.GetNames()
	if len(names) != 2 || names[0].Name != "Alpha" || names[1].Name != "Bravo" {
		t.Fatalf("unexpected names: %+v", names)
	}
	raw := dt.GetNamesRaw()
	if len(raw) != 2 || raw[0] != "Alpha" || raw[1] != "Bravo" {
		t.Fatalf("unexpected raw names: %+v", raw)
	}
	if got, ok := dt.GetFirstOccurance("alpha"); !ok || got != setA.ID {
		t.Fatalf("expected GetFirstOccurance to return %d, got %d (ok=%v)", setA.ID, got, ok)
	}
	if sets := dt.GetSetsByName("Alpha"); len(sets) != 1 || sets[0].ID != setA.ID {
		t.Fatalf("unexpected GetSetsByName result: %+v", sets)
	}
	if idx, got := dt.GetDataByName("Bravo"); idx != 1 || got != setB.ID {
		t.Fatalf("unexpected GetDataByName: idx=%d id=%d", idx, got)
	}
	if idx := dt.GetIdx(setB.ID); idx != 1 {
		t.Fatalf("expected index 1 for setB, got %d", idx)
	}
	if set, ok := dt.GetPtr(setA.ID); !ok || set.Name != "Alpha" {
		t.Fatalf("GetPtr failed: %+v ok=%v", set, ok)
	}
	if _, ok := dt.GetPtr(999); ok {
		t.Fatalf("expected GetPtr to fail for missing id")
	}

	out := captureOutput(t, func() {
		dt.PrintListFiltered(true, nil)
	})
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "Bravo") {
		t.Fatalf("expected list output to include names, got: %s", out)
	}

	dt.Delete(setA.ID)
	if idx := dt.GetIdx(setA.ID); idx != -1 {
		t.Fatalf("expected deleted item to be missing, got idx=%d", idx)
	}
}

// TestDataTableAddEmptyName verifies default naming on Add.
func TestDataTableAddEmptyName(t *testing.T) {
	dt := NewDataTable[Category]()
	dt.Add(Dataset[Category]{ID: 1, Name: "", Data: Category{}})
	if dt[0].Name != "Unnamed" {
		t.Fatalf("expected default name, got %q", dt[0].Name)
	}
}
