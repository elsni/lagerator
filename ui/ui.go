package ui

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/loggi"
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
	app := tview.NewApplication()
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
func EditItem[T data.CustomData](r data.Dataset[T], idoptions [][]DropdownOptions, caption string) (data.Dataset[T], bool) {
	saved := false
	app := tview.NewApplication()
	var fields []string
	fields, _ = reflections.Fields(r.Data)
	ididx := 0
	fheight := 17
	taglist := ""

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
			var opt []string
			initialopt := -1
			for i, ido := range idoptions[ididx] {
				opt = append(opt, " "+ido.Name+" ")
				if ido.Id == fieldvalue.(uint32) {
					initialopt = i
				}
			}
			//form.AddDropDown(fieldName[:len(fieldName)-2], opt, initialopt, nil)
			dd := tview.NewDropDown()
			dd.SetLabel(fieldName[:len(fieldName)-2])
			dd.SetOptions(opt, nil)
			dd.SetCurrentOption(initialopt)
			dd.SetTextOptions("", "", "", "", "please select")
			form.AddFormItem(dd)
			ididx += 1
			fheight += 2
		/*case "[]uint32":
		taglist := tview.NewList().ShowSecondaryText(false)
		alltags := data.Db.Tags.GetNames()
		for _, tn := range alltags {
			taglist.AddItem(tn.Name, "", 0, func() {

			})
		}
		//form.AddFormItem(taglist)

		ipf := tview.NewInputField().
			SetLabel("Add Tag: ").
			SetFieldWidth(40)

		ipf.SetAutocompleteFunc(func(currentText string) (entries []string) {
			if len(currentText) == 0 {
				return
			}
			for _, tag := range alltags {
				if strings.HasPrefix(strings.ToLower(tag.Name), strings.ToLower(currentText)) {
					entries = append(entries, tag.Name)
				}
			}
			if len(entries) <= 1 {
				entries = nil
			}
			return
		})
		ipf.SetAutocompletedFunc(func(text string, index, source int) bool {
			if source != tview.AutocompletedNavigate {
				ipf.SetText(text)
			}
			return source == tview.AutocompletedEnter || source == tview.AutocompletedClick
		})
		*/
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
	tagfield := tview.NewInputField().SetLabel("Tags").SetText(data.GetTagList(r.Tags)).SetFieldWidth(40).SetChangedFunc(func(text string) { taglist = text })
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
		r.Tags = data.GetTagIds(taglist)
		r.Updated = time.Now().Unix()
		saved = true
		app.Stop()
	})
	form.AddButton("Quit", func() {
		saved = false
		app.Stop()
	})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			loggi.Log.Log("ESC")
			saved = false
			app.Stop()
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
	if err := app.SetRoot(form, false).SetFocus(form).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return r, saved
}

// Alert shows a confirmation dialog and returns true on acceptance.
func Alert(message string) bool {
	app := tview.NewApplication()
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
func SelectItem[T data.CustomData](items []data.Dataset[T], action string) int {
	numitems := len(items)
	lastidx := numitems - 1
	if numitems == 0 {
		return 0
	}

	app := tview.NewApplication()
	menu := tview.NewList()

	result := -1

	for i, set := range items {
		menu.AddItem(set.GetTableRow(), "", rune('a'+i), nil)
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
	menu.SetRect(0, 0, len(items[0].GetTableRow())+6, len(items)+3)

	if err := app.SetRoot(menu, false).Run(); err != nil {
		panic(err)
	}
	return result
}
