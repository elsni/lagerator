package ui

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/loggi"
	"github.com/elsni/lagerator/presentation"
	"github.com/gdamore/tcell/v2"
	"github.com/oleiade/reflections"
	"github.com/rivo/tview"
)

type DropdownOptions struct {
	Id   uint32
	Name string
}

// TestForm shows a demo form for UI testing.
func TestForm() {
	app := configureApplication(tview.NewApplication())
	form := tview.NewForm().
		AddInputField("Monday", "Test", 40, nil, nil).
		AddInputField("Tuesday", "Test", 40, nil, nil).
		AddInputField("Wednesday", "Test", 40, nil, nil).
		AddInputField("Thursday", "Test", 40, nil, nil)
	form.AddButton("Quit", func() {
		app.Stop()
	})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			app.Stop()
			return nil
		case tcell.KeyPgDn:
			fallthrough
		case tcell.KeyPgUp:
			maxidx := form.GetFormItemCount() - 1
			fidx, _ := form.GetFocusedItemIndex()
			// focused itemindex is -1 when a button is focused
			if fidx > -1 {
				if event.Key() == tcell.KeyPgDn {
					fidx += 1
					if fidx > maxidx {
						fidx = 0
					}
				}
				if event.Key() == tcell.KeyPgUp {
					fidx -= 1
					if fidx < 0 {
						fidx = maxidx
					}
				}

				form.SetFocus(fidx)
			}
			return nil
		}

		return event
	})
	form.SetRect(0, 0, 56, 16)
	if err := app.SetRoot(form, false).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// EditItem opens the edit form for a dataset and returns the updated data.
// idoptions are the reference dropdown entries, grouped per field.
func EditItem[T any](r data.Dataset[T], idoptions [][]DropdownOptions, caption string) (data.Dataset[T], bool) {
	saved := false
	app := configureApplication(tview.NewApplication())
	form := newEditForm(r, idoptions, caption, func(result data.Dataset[T], accepted bool) {
		r, saved = result, accepted
		app.Stop()
	})
	if err := app.SetRoot(form, false).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return r, saved
}

// newEditForm builds an editor whose completion callback handles saving or cancellation.
func newEditForm[T any](r data.Dataset[T], idoptions [][]DropdownOptions, caption string, done func(data.Dataset[T], bool)) *tview.Form {
	var fields []string
	fields, _ = reflections.Fields(r.Data)
	ididx := 0
	fheight := 17

	form := tview.NewForm().
		AddTextView("Id", fmt.Sprint(r.ID), 5, 1, false, false).
		AddInputField("Name", r.Name, 40, nil, func(text string) { r.Name = text }).
		AddTextArea("Description", r.Description, 40, 0, 0, func(text string) { r.Description = text })

	for _, fieldName := range fields {
		fieldtype, _ := reflections.GetFieldType(r.Data, fieldName)
		fieldvalue, _ := reflections.GetField(r.Data, fieldName)
		switch fieldtype {
		case "string":
			form.AddInputField(fieldName, fieldvalue.(string), 40, nil, nil)
			fheight += 2
		case "uint32":
			form.AddFormItem(newReferenceDropdown(fieldName[:len(fieldName)-2], idoptions[ididx], fieldvalue.(uint32)))
			ididx += 1
			fheight += 2
		case "int":
			form.AddInputField(fieldName, fmt.Sprint(fieldvalue), 6, func(text string, last rune) bool {
				_, err := strconv.Atoi(text)
				return err == nil
			}, nil)
			fheight += 2
		default:
			loggi.Log.Log(fieldtype)
		}

	}
	tagfield := newTagInput(data.GetTagList(r.Tags))
	form.AddFormItem(tagfield)
	form.AddButton("Ok", func() {
		ididx = 0
		for i, fieldName := range fields {
			fieldtype, _ := reflections.GetFieldType(r.Data, fieldName)
			f := form.GetFormItem(i + 3)
			if f != nil {
				fmt.Println(fieldtype)
				switch fieldtype {
				case "string":
					text := f.(*tview.InputField).GetText()
					reflections.SetField(&r.Data, fieldName, text)

				case "uint32":
					index, _ := f.(*tview.DropDown).GetCurrentOption()
					if index > -1 {
						reflections.SetField(&r.Data, fieldName, idoptions[ididx][index].Id)
					}
					ididx += 1
				case "int":
					text := f.(*tview.InputField).GetText()
					val, _ := strconv.Atoi(text)
					reflections.SetField(&r.Data, fieldName, val)
				default:
				}
			}
		}
		r.Tags = data.GetTagIds(tagfield.GetText())
		r.Updated = time.Now().Unix()
		done(r, true)
	})
	form.AddButton("Quit", func() {
		done(r, false)
	})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			loggi.Log.Log("ESC")
			done(r, false)
			return nil
		case tcell.KeyPgDn:
			fallthrough
		case tcell.KeyPgUp:
			loggi.Log.Log("Page Key")
			maxidx := form.GetFormItemCount() - 1
			fidx, _ := form.GetFocusedItemIndex()
			loggi.Log.Log(fmt.Sprintf("Current Focus: %d", fidx))
			// focused itemindex is -1 when a button is focused
			if fidx > -1 {
				if event.Key() == tcell.KeyPgDn {
					fidx += 1
					if fidx > maxidx {
						fidx = 1
					}
				}
				if event.Key() == tcell.KeyPgUp {
					fidx -= 1
					if fidx < 1 {
						fidx = maxidx
					}
				}

				loggi.Log.Log(fmt.Sprintf("New Focus set: %d", fidx))
				form.SetFocus(fidx)
				nfidx, _ := form.GetFocusedItemIndex()
				loggi.Log.Log(fmt.Sprintf("New Focus read: %d", nfidx))
			}
			return nil
		}

		return event
	})
	form.SetFieldBackgroundColor(tcell.NewRGBColor(20, 20, 20))
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetButtonTextColor(tcell.ColorRed)

	form.SetBorder(true).SetTitle(caption + reflect.TypeOf(r.Data).String()[5:] + " ").SetTitleAlign(tview.AlignLeft)
	form.SetRect(0, 0, 56, fheight)
	return form
}

// Alert shows a confirmation dialog and returns true on acceptance.
func Alert(message string) bool {
	app := configureApplication(tview.NewApplication())
	form := tview.NewForm()
	result := false
	form.AddTextView("", message, 30, 3, true, false)

	form.AddButton("No", func() {
		result = false
		app.Stop()
	})

	form.AddButton("Yes", func() {
		result = true
		app.Stop()
	})

	form.SetBorder(true).SetTitle("Please confirm").SetTitleAlign(tview.AlignLeft)
	form.SetFieldBackgroundColor(tcell.NewRGBColor(20, 20, 20))
	form.SetFieldTextColor(tcell.ColorWhite)
	form.SetButtonTextColor(tcell.ColorRed)
	form.SetRect(0, 0, 32, 10)
	if err := app.SetRoot(form, false).Run(); err != nil {
		panic(err)
	}
	return result
}

// SelectItem lets the user choose an item and returns its index or -1.
func SelectItem[T presentation.Entity](items []data.Dataset[T], action string) int {
	numitems := len(items)
	lastidx := numitems - 1
	if numitems == 0 {
		return 0
	}

	app := configureApplication(tview.NewApplication())
	menu := newSelectionList()

	result := -1

	for i, set := range items {
		menu.AddItem(presentation.GetTableRow(set), "", rune('a'+i), nil)
	}
	menu.AddItem("Quit", "Press to exit", 'q', nil)
	menu.SetTitle(fmt.Sprintf("Select %s to %s", reflect.TypeOf(items[0].Data).String()[5:], action))
	menu.SetTitleColor(tcell.ColorYellow)
	menu.SetBorder(true)
	menu.SetBorderColor(tcell.ColorDarkCyan)
	//menu.Set
	menu.ShowSecondaryText(false)
	menu.SetSelectedFunc(func(idx int, maintext, secondarytext string, shortcut rune) {
		if idx <= lastidx {
			result = idx
		}
		app.Stop()
	})
	menu.SetRect(0, 0, len(presentation.GetTableRow(items[0]))+6, len(items)+3)

	if err := app.SetRoot(menu, false).Run(); err != nil {
		panic(err)
	}
	return result
}
