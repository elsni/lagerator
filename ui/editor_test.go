package ui

import (
	"reflect"
	"testing"

	"github.com/elsni/lagerator/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// editorFixture provides an item and reference choices for editor interaction tests.
func editorFixture() (data.Dataset[data.Item], [][]DropdownOptions) {
	return data.Dataset[data.Item]{ID: 1, Name: "Item", Tags: []uint32{7}, Data: data.Item{
		Amount: 1, BoxId: 4, CategoryId: 6,
	}}, [][]DropdownOptions{
		{{Id: 4, Name: "Box A"}, {Id: 5, Name: "Box [red]"}},
		{{Id: 6, Name: "Tools"}},
	}
}

// TestEditorSave exercises updated widgets, reflection writes, and unchanged tag retention.
func TestEditorSave(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	var result data.Dataset[data.Item]
	var saved bool
	form := newEditForm(item, options, " Edit ", func(value data.Dataset[data.Item], accepted bool) {
		result, saved = value, accepted
	})
	form.GetFormItemByLabel("Name").(*tview.InputField).SetText("Updated")
	form.GetFormItemByLabel("Location").(*tview.InputField).SetText("Shelf B")
	amount := form.GetFormItemByLabel("Amount").(*tview.InputField)
	amount.SetText("2")
	pressKey(amount, tcell.KeyRune, 'x')
	if amount.GetText() != "2" {
		t.Fatal("amount accepted non-numeric input")
	}
	form.GetFormItemByLabel("Box").(*referenceInput).SetText(referenceLabel(options[0][1]))
	form.Draw(simulationScreen(t))
	pressKey(form.GetButton(form.GetButtonIndex("Ok")), tcell.KeyEnter, 0)
	if !saved || result.Name != "Updated" || result.Data.Location != "Shelf B" ||
		result.Data.Amount != 2 || result.Data.BoxId != 5 || result.Data.CategoryId != 6 ||
		!reflect.DeepEqual(result.Tags, item.Tags) || result.Updated == 0 {
		t.Fatalf("saved item is incorrect: %+v, saved=%v", result, saved)
	}
}

// TestEditorCancelDoesNotCreateTags verifies Quit and Escape do not apply tag edits.
func TestEditorCancelDoesNotCreateTags(t *testing.T) {
	for _, mode := range []string{"Quit", "Escape"} {
		t.Run(mode, func(t *testing.T) {
			setupInventory(t)
			item, options := editorFixture()
			called, saved := false, true
			form := newEditForm(item, options, " Edit ", func(value data.Dataset[data.Item], accepted bool) {
				called, saved = true, accepted
			})
			form.GetFormItemByLabel("Tags").(*tview.InputField).SetText("new tag")
			if mode == "Quit" {
				pressKey(form.GetButton(form.GetButtonIndex("Quit")), tcell.KeyEnter, 0)
			} else {
				pressKey(form, tcell.KeyEscape, 0)
			}
			if !called || saved || len(data.Db.Tags) != 3 {
				t.Fatalf("cancel mutated inventory: called=%v saved=%v tags=%+v", called, saved, data.Db.Tags)
			}
		})
	}
}

// TestEditorSaveCompletedTags verifies literal tag completion is resolved to its existing ID.
func TestEditorSaveCompletedTags(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	var result data.Dataset[data.Item]
	form := newEditForm(item, options, " Edit ", func(value data.Dataset[data.Item], accepted bool) {
		result = value
	})
	field := form.GetFormItemByLabel("Tags").(*tview.InputField)
	field.SetText("fragile, [")
	field.Autocomplete()
	pressKey(field, tcell.KeyEnter, 0)
	pressKey(form.GetButton(form.GetButtonIndex("Ok")), tcell.KeyEnter, 0)
	if !reflect.DeepEqual(result.Tags, []uint32{7, 8}) || len(data.Db.Tags) != 3 {
		t.Fatalf("completed tags were not preserved: %+v", result.Tags)
	}
}
