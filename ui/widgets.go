package ui

import (
	"strings"

	"github.com/elsni/lagerator/data"
	"github.com/rivo/tview"
)

// configureApplication configures the terminal title and bracketed paste for inventory dialogs.
func configureApplication(app *tview.Application) *tview.Application {
	return app.SetTitle("Lagerator").EnablePaste(true)
}

// newReferenceDropdown displays inventory names literally and preserves numeric references.
func newReferenceDropdown(label string, options []DropdownOptions, current uint32) *tview.DropDown {
	texts := make([]string, len(options))
	selected := -1
	for index, option := range options {
		texts[index] = " " + option.Name + " "
		if option.Id == current {
			selected = index
		}
	}
	field := tview.NewDropDown().SetUseStyleTags(false)
	field.SetLabel(label)
	field.SetOptions(texts, nil)
	field.SetCurrentOption(selected)
	field.SetTextOptions("", "", "", "", "please select")
	return field
}

// newSelectionList prevents bracketed inventory names from becoming style markup.
func newSelectionList() *tview.List {
	return tview.NewList().SetUseStyleTags(false, false)
}

// newTagInput completes the final comma-separated tag without interpreting style markup.
func newTagInput(text string) *tview.InputField {
	field := tview.NewInputField().SetLabel("Tags").SetText(text).SetFieldWidth(40)
	var matches []string
	field.SetAutocompleteFunc(func(text string) []string {
		matches = tagCompletions(text)
		entries := make([]string, len(matches))
		for index, match := range matches {
			entries[index] = tview.Escape(match)
		}
		return entries
	})
	field.SetAutocompletedFunc(func(text string, index, source int) bool {
		if source == tview.AutocompletedNavigate {
			return false
		}
		field.SetText(matches[index])
		return true
	})
	return field
}

// tagCompletions suggests existing tags for the last token, preserving earlier tokens.
func tagCompletions(text string) []string {
	prefix, token := "", text
	if index := strings.LastIndex(text, ","); index >= 0 {
		prefix, token = text[:index+1], text[index+1:]
		prefix += token[:len(token)-len(strings.TrimLeft(token, " \t"))]
	}
	query := strings.ToLower(strings.TrimSpace(token))
	if query == "" {
		return nil
	}
	var matches []string
	for _, name := range data.Db.Tags.GetNamesRaw() {
		if strings.HasPrefix(strings.ToLower(name), query) {
			matches = append(matches, prefix+name)
		}
	}
	return matches
}
