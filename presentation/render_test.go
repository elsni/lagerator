package presentation

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/id"
	"github.com/elsni/lagerator/terminal"
)

// resetDb isolates the database and ID source for a presentation test.
func resetDb(t *testing.T) {
	t.Helper()
	previousDB, previousID := data.Db, id.IdSource
	t.Cleanup(func() {
		data.Db, id.IdSource = previousDB, previousID
	})
	data.Db = data.NewDatabase()
	id.IdSource.SetLastId(0)
}

// captureOutput captures terminal output without depending on pipe buffer capacity.
func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "stdout")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	previous := os.Stdout
	os.Stdout = file
	defer func() { os.Stdout = previous }()
	fn()
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	output, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	return string(output)
}

// TestPrintNamePreservesMarkersAndUnicode checks the exact legacy display rules.
func TestPrintNamePreservesMarkersAndUnicode(t *testing.T) {
	resetDb(t)
	data.Db.CurrentWarehouse = 3
	set := data.Dataset[data.Room]{ID: 3, Name: "Ä長Room", Deleted: true}
	if got := GetPrintName(set, 256); got != "Ä長Room*(deleted)" {
		t.Fatalf("unexpected markers: %q", got)
	}
	if got := GetPrintName(set, 4); got != "Ä長R…" {
		t.Fatalf("unexpected truncation: %q", got)
	}
	if got := GetTableRow(set); got != "" {
		t.Fatalf("deleted row must remain hidden: %q", got)
	}
	if got := GetPrintNameById(&data.Db.Rooms, 999, 20); got != "not found" {
		t.Fatalf("missing reference text changed: %q", got)
	}
}

// TestListFilteringPreservesHeadersAndOrder checks filtering, sorting, and empty tables.
func TestListFilteringPreservesHeadersAndOrder(t *testing.T) {
	resetDb(t)
	table := data.DataTable[data.Category]{
		{ID: 1, Name: "Zulu", Tags: []uint32{10}},
		{ID: 2, Name: "Alpha", Tags: []uint32{10}},
		{ID: 3, Name: "Hidden", Tags: []uint32{10}, Deleted: true},
	}
	data.Db.Tags.Add(data.Dataset[data.Tag]{ID: 10, Name: "fragile"})
	unsorted := captureOutput(t, func() { PrintList(&table, false) })
	sorted := captureOutput(t, func() { PrintListByTag(&table, "FRAGILE", true) })
	if strings.Index(unsorted, "Zulu") > strings.Index(unsorted, "Alpha") ||
		strings.Index(sorted, "Alpha") > strings.Index(sorted, "Zulu") ||
		strings.Contains(sorted, "Hidden") {
		t.Fatalf("list sorting/filtering changed: %q / %q", unsorted, sorted)
	}
	unknown := captureOutput(t, func() { PrintListByTag(&table, "unknown", false) })
	if unknown != "Unknown tag \"unknown\"\n" {
		t.Fatalf("unknown tag message changed: %q", unknown)
	}
	empty := data.NewDataTable[data.Category]()
	if got := captureOutput(t, func() { PrintList(&empty, false) }); got != "" {
		t.Fatalf("empty list emitted output: %q", got)
	}
	filtered := captureOutput(t, func() { PrintListByTagId(&table, 999, false) })
	want := terminal.GetHeadlineText(" Category ") + "\n" +
		terminal.SetBgColor(terminal.COLORBLUE) + "   ID Name                           " + terminal.ResetColor() + "\n"
	if filtered != want {
		t.Fatalf("empty results must preserve header: %q, want %q", filtered, want)
	}
}
