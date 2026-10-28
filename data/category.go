package data

import "sort"

type Category struct {
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
