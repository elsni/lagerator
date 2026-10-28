package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TestApplicationTitleAndPaste exercises a complete simulated terminal event loop.
func TestApplicationTitleAndPaste(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	app := configureApplication(tview.NewApplication().SetScreen(screen))
	var title string
	app.SetAfterDrawFunc(func(tcell.Screen) { title = screen.GetTitle() })
	text := tview.NewTextArea()
	app.SetRoot(text, true).SetFocus(text)
	app.QueueEvent(tcell.NewEventPaste(true))
	app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'A', tcell.ModNone))
	app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))
	app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'B', tcell.ModNone))
	app.QueueEvent(tcell.NewEventPaste(false))
	app.QueueEvent(tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone))
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	if title != "Lagerator" || text.GetText() != "A\nB" {
		t.Fatalf("application configuration failed: title=%q text=%q", title, text.GetText())
	}
}
