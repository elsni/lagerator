package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

// referenceInput searches inventory names while retaining unambiguous numeric references.
type referenceInput struct {
	*tview.InputField
	options  []DropdownOptions
	original uint32
	initial  string
}

// referenceLabel distinguishes duplicate names and preserves literal inventory text.
func referenceLabel(option DropdownOptions) string {
	return fmt.Sprintf("%s (#%d)", option.Name, option.Id)
}

// newReferenceInput builds a substring-searchable field with keyboard completion.
func newReferenceInput(label string, options []DropdownOptions, current uint32) *referenceInput {
	field := &referenceInput{InputField: tview.NewInputField(), options: options, original: current}
	field.SetLabel(label).SetFieldWidth(40)
	for _, option := range options {
		if option.Id == current {
			field.initial = referenceLabel(option)
			break
		}
	}
	field.SetText(field.initial)
	var matches []string
	field.SetAutocompleteFunc(func(query string) []string {
		matches = nil
		for _, option := range options {
			text := referenceLabel(option)
			if strings.Contains(strings.ToLower(text), strings.ToLower(strings.TrimSpace(query))) {
				matches = append(matches, text)
			}
		}
		escaped := make([]string, len(matches))
		for i, text := range matches {
			escaped[i] = tview.Escape(text)
		}
		return escaped
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

// referenceID validates a completed selection without changing an untouched missing reference.
func (field *referenceInput) referenceID() (uint32, bool) {
	text := field.GetText()
	if text == field.initial {
		return field.original, true
	}
	for _, option := range field.options {
		if text == referenceLabel(option) {
			return option.Id, true
		}
	}
	return 0, false
}
