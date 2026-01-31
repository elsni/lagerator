package data

import (
	"sort"

	"github.com/elsni/lagerator/terminal"
)

type Category struct {
}

// GetTableRow returns the formatted row for a category.
func (d Category) GetTableRow(ownid uint32) string {
	return ""
}

// GetTableHeader returns the category table header.
func (d Category) GetTableHeader() string {
	return terminal.ResetColor()
}

// Show prints category details.
func (d Category) Show() {
}

type CategoryTable = DataTable[Category]

// GetCategoriesSorted returns categories sorted by name.
func GetCategoriesSorted(ct *CategoryTable) []Listentry {
	var names []Listentry
	for _, set := range *ct {
		if !set.Deleted {
			names = append(names, Listentry{Id: set.ID, Name: set.Name})
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].Name < names[j].Name
	})
	return names
}
