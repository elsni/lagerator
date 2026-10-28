package ui

import (
	"testing"

	"github.com/elsni/lagerator/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TestEditorTabNavigation exercises focus transitions through the application event loop.
func TestEditorTabNavigation(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	screen := tcell.NewSimulationScreen("UTF-8")
	app := configureApplication(tview.NewApplication().SetScreen(screen))
	form := newEditForm(item, options, " Edit ", func(data.Dataset[data.Item], bool) { app.Stop() })
	bindEditorInput(app, form)
	app.SetRoot(form, false).SetFocus(form)
	for _, char := range []rune{'n', 'd', 'l', 'c', '2'} {
		app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, char, tcell.ModNone))
		app.QueueEvent(tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone))
	}
	app.QueueEvent(tcell.NewEventKey(tcell.KeyCtrlQ, 0, tcell.ModNone))
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	for label, want := range map[string]string{"Name": "Itemn", "Location": "l", "Condition": "c", "Amount": "12"} {
		got := form.GetFormItemByLabel(label).(*tview.InputField).GetText()
		if got != want {
			t.Errorf("%s=%q, want %q", label, got, want)
		}
	}
}

// TestEditorBacktabNavigation checks reverse navigation and returning through Description.
func TestEditorBacktabNavigation(t *testing.T) {
	setupInventory(t)
	item, options := editorFixture()
	screen := tcell.NewSimulationScreen("UTF-8")
	app := configureApplication(tview.NewApplication().SetScreen(screen))
	form := newEditForm(item, options, " Edit ", func(data.Dataset[data.Item], bool) { app.Stop() })
	bindEditorInput(app, form)
	app.SetRoot(form, false).SetFocus(form)
	for _, event := range []struct {
		key  tcell.Key
		char rune
	}{
		{tcell.KeyTab, 0}, {tcell.KeyTab, 0}, {tcell.KeyTab, 0},
		{tcell.KeyRune, 'c'}, {tcell.KeyBacktab, 0}, {tcell.KeyRune, 'l'},
		{tcell.KeyBacktab, 0}, {tcell.KeyRune, 'd'}, {tcell.KeyTab, 0},
		{tcell.KeyRune, 'x'}, {tcell.KeyTab, 0}, {tcell.KeyRune, 'y'},
		{tcell.KeyCtrlQ, 0},
	} {
		app.QueueEvent(tcell.NewEventKey(event.key, event.char, tcell.ModNone))
	}
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	for label, want := range map[string]string{"Location": "lx", "Condition": "cy"} {
		if got := form.GetFormItemByLabel(label).(*tview.InputField).GetText(); got != want {
			t.Errorf("%s=%q, want %q", label, got, want)
		}
	}
	if got := form.GetFormItemByLabel("Description").(*tview.TextArea).GetText(); got != "d" {
		t.Errorf("Description=%q, want d", got)
	}
}
