package ui

import (
	"reflect"
	"strings"
	"testing"

	"github.com/elsni/lagerator/data"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// setupInventory isolates the database used by widget tests.
func setupInventory(t *testing.T) {
	t.Helper()
	previous := data.Db
	data.Db = data.NewDatabase()
	t.Cleanup(func() { data.Db = previous })
	data.Db.Tags.Add(data.Dataset[data.Tag]{ID: 7, Name: "fragile"})
	data.Db.Tags.Add(data.Dataset[data.Tag]{ID: 8, Name: "[red]"})
	data.Db.Tags.Add(data.Dataset[data.Tag]{ID: 9, Name: "deleted", Deleted: true})
}

// simulationScreen creates an isolated terminal for widget rendering and input tests.
func simulationScreen(t *testing.T) tcell.SimulationScreen {
	t.Helper()
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	screen.SetSize(120, 50)
	t.Cleanup(screen.Fini)
	return screen
}

// screenText reads terminal cells so tests verify visible text rather than widget settings.
func screenText(screen tcell.Screen) string {
	width, height := screen.Size()
	var output strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			main, _, _, _ := screen.GetContent(x, y)
			output.WriteRune(main)
		}
		output.WriteByte('\n')
	}
	return output.String()
}

// pressKey sends a keyboard event through a widget's real input handler.
func pressKey(widget tview.Primitive, key tcell.Key, char rune) {
	widget.InputHandler()(tcell.NewEventKey(key, char, tcell.ModNone), func(tview.Primitive) {})
}

// TestLiteralSelectionNames verifies bracketed names survive list and dropdown rendering.
func TestLiteralSelectionNames(t *testing.T) {
	screen := simulationScreen(t)
	name := "Box [red] 長"
	list := newSelectionList().AddItem(name, "", 0, nil)
	list.SetRect(0, 0, 80, 3)
	list.Draw(screen)
	if !strings.Contains(screenText(screen), name) {
		t.Fatalf("list interpreted inventory text as markup: %q", screenText(screen))
	}
	screen.Clear()
	field := newReferenceDropdown("Box", []DropdownOptions{{Id: 4, Name: "First"}, {Id: 5, Name: name}}, 4)
	field.SetRect(0, 0, 80, 1)
	field.Focus(func(tview.Primitive) {})
	pressKey(field, tcell.KeyEnter, 0)
	pressKey(field, tcell.KeyDown, 0)
	pressKey(field, tcell.KeyEnter, 0)
	field.Draw(screen)
	index, text := field.GetCurrentOption()
	if index != 1 || text != " "+name+" " || !strings.Contains(screenText(screen), name) {
		t.Fatalf("dropdown selection/rendering changed: index=%d text=%q screen=%q", index, text, screenText(screen))
	}
}

// TestTagCompletions verifies case-insensitive final-token completion and deleted-tag exclusion.
func TestTagCompletions(t *testing.T) {
	setupInventory(t)
	for _, tc := range []struct {
		input string
		want  []string
	}{
		{"FR", []string{"fragile"}},
		{"new, fr", []string{"new, fragile"}},
		{"new,[", []string{"new,[red]"}},
		{"new,\t[", []string{"new,\t[red]"}},
		{"", nil}, {"new, ", nil}, {"del", nil}, {"unknown", nil},
	} {
		t.Run(tc.input, func(t *testing.T) {
			if got := tagCompletions(tc.input); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("completions=%q, want %q", got, tc.want)
			}
		})
	}
}

// TestTagCompletionKeyboardAndRendering preserves literal brackets when completing with Enter or Tab.
func TestTagCompletionKeyboardAndRendering(t *testing.T) {
	setupInventory(t)
	for _, key := range []tcell.Key{tcell.KeyEnter, tcell.KeyTab} {
		screen := simulationScreen(t)
		field := newTagInput("new, [")
		field.SetRect(0, 0, 80, 1)
		field.Focus(func(tview.Primitive) {})
		field.Draw(screen)
		if !strings.Contains(screenText(screen), "new, [red]") {
			t.Fatalf("completion is not displayed literally: %q", screenText(screen))
		}
		pressKey(field, key, 0)
		if field.GetText() != "new, [red]" {
			t.Fatalf("completion changed tag text: %q", field.GetText())
		}
	}
}
