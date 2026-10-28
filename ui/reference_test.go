package ui

import (
	"strings"
	"testing"

	"github.com/elsni/lagerator/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TestReferenceSearchKeyboard verifies substring completion, literal names, and duplicate IDs.
func TestReferenceSearchKeyboard(t *testing.T) {
	for _, key := range []tcell.Key{tcell.KeyEnter, tcell.KeyTab} {
		field := newReferenceInput("Box", []DropdownOptions{
			{Id: 4, Name: "Storage [red]"}, {Id: 5, Name: "Storage [red]"},
		}, 4)
		field.SetText("RAGE [R")
		field.Autocomplete()
		field.SetRect(0, 0, 80, 1)
		field.Focus(nil)
		screen := simulationScreen(t)
		field.Draw(screen)
		if !strings.Contains(screenText(screen), "Storage [red] (#5)") {
			t.Fatalf("matching duplicate not rendered literally: %q", screenText(screen))
		}
		pressKey(field, tcell.KeyDown, 0)
		pressKey(field, key, 0)
		if id, ok := field.referenceID(); !ok || id != 5 {
			t.Fatalf("completion lost reference: id=%d valid=%v text=%q", id, ok, field.GetText())
		}
	}
}

// TestReferenceValidation preserves missing original IDs and rejects unfinished searches.
func TestReferenceValidation(t *testing.T) {
	field := newReferenceInput("Box", []DropdownOptions{{Id: 4, Name: "Box"}}, 99)
	if id, ok := field.referenceID(); id != 99 || !ok {
		t.Fatal("untouched missing reference changed")
	}
	field.SetText("unknown")
	if _, ok := field.referenceID(); ok {
		t.Fatal("unknown reference accepted")
	}
}

// TestEditorShortcutsAndInvalidSelection verifies visible controls and validation before saving.
func TestEditorShortcutsAndInvalidSelection(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	calls, saved := 0, false
	form := newEditForm(item, options, " Edit ", func(result data.Dataset[data.Item], accepted bool) {
		calls++
		saved = accepted
	})
	screen := simulationScreen(t)
	form.Draw(screen)
	if !strings.Contains(screenText(screen), "Ctrl+S: save") {
		t.Fatal("save shortcut is hidden")
	}
	box := form.GetFormItemByLabel("Box").(*referenceInput)
	box.SetText("not a box")
	pressKey(form, tcell.KeyCtrlS, 0)
	if calls != 0 || !strings.Contains(form.GetTitle(), "valid Box") {
		t.Fatal("invalid reference was saved")
	}
	box.SetText(referenceLabel(options[0][1]))
	pressKey(form, tcell.KeyCtrlS, 0)
	if calls != 1 || !saved {
		t.Fatal("Ctrl+S did not save corrected form")
	}
	pressKey(form, tcell.KeyCtrlQ, 0)
	if calls != 2 || saved {
		t.Fatal("Ctrl+Q did not cancel")
	}
}

// TestEditorDescriptionTabNavigation verifies multiline fields allow forward and reverse navigation.
func TestEditorDescriptionTabNavigation(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	form := newEditForm(item, options, " Edit ", func(data.Dataset[data.Item], bool) {})
	form.SetFocus(2)
	pressKey(form, tcell.KeyTab, 0)
	if index, _ := form.GetFocusedItemIndex(); index != 3 {
		t.Fatalf("Tab focus=%d", index)
	}
	form.SetFocus(2)
	pressKey(form, tcell.KeyBacktab, 0)
	if index, _ := form.GetFocusedItemIndex(); index != 1 {
		t.Fatalf("Shift+Tab focus=%d", index)
	}
}

// TestEditorApplicationShortcuts verifies save works through the real application focus chain.
func TestEditorApplicationShortcuts(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	screen := tcell.NewSimulationScreen("UTF-8")
	app := configureApplication(tview.NewApplication().SetScreen(screen))
	saved := false
	form := newEditForm(item, options, " Edit ", func(result data.Dataset[data.Item], accepted bool) {
		saved = accepted
		app.Stop()
	})
	app.SetInputCapture(form.GetInputCapture())
	app.SetRoot(form, false).SetFocus(form)
	app.QueueEvent(tcell.NewEventKey(tcell.KeyCtrlS, 0, tcell.ModNone))
	app.QueueEvent(tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone))
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	if !saved {
		t.Fatal("save shortcut did not reach the focused form")
	}
}
