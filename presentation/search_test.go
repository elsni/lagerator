package presentation

import (
	"strings"
	"testing"

	"github.com/elsni/lagerator/data"
)

// TestFindItem verifies search output for items.
func TestFindItem(t *testing.T) {
	resetDb(t)
	item := data.NewDataset[data.Item]("Hammer", data.Item{Amount: 1, Location: "Shelf A"})
	item.Description = "Steel tool"
	data.Db.Items.Add(item)
	data.Db.Items.Add(data.NewDataset[data.Item]("Nails", data.Item{Amount: 50, Location: "Shelf B"}))

	out := captureOutput(t, func() {
		FindItem("hammer", false)
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
	resetDb(t)
	out := captureOutput(t, func() {
		FindItem("anything", false)
	})
	if !strings.Contains(out, "no data") {
		t.Fatalf("expected no data output, got: %s", out)
	}
}

// TestFindItemNoMatch verifies no matches are printed.
func TestFindItemNoMatch(t *testing.T) {
	resetDb(t)
	data.Db.Items.Add(data.NewDataset[data.Item]("Hammer", data.Item{Amount: 1, Location: "Shelf A"}))
	out := captureOutput(t, func() {
		FindItem("wrench", false)
	})
	if strings.Contains(out, "Hammer") {
		t.Fatalf("did not expect Hammer in output: %s", out)
	}
}
