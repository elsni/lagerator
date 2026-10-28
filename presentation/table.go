package presentation

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/elsni/lagerator/data"
)

type sortEntry struct {
	Name string
	Text string
}

// PrintList prints the table, optionally sorted by name.
func PrintList[T Entity](st *data.DataTable[T], sortname bool) {
	PrintListFiltered(st, sortname, nil)
}

// PrintListByTag prints entries filtered by tag name.
func PrintListByTag[T Entity](st *data.DataTable[T], tagname string, sortname bool) {
	tagid, found := data.Db.Tags.GetFirstOccurance(tagname)
	if !found {
		fmt.Printf("Unknown tag \"%s\"\n", tagname)
		return
	}
	PrintListFiltered(st, sortname, func(set data.Dataset[T]) bool {
		return slices.Contains(set.Tags, tagid)
	})
}

// PrintListByTagId prints entries filtered by tag id.
func PrintListByTagId[T Entity](st *data.DataTable[T], tagid uint32, sortname bool) {
	PrintListFiltered(st, sortname, func(set data.Dataset[T]) bool {
		return slices.Contains(set.Tags, tagid)
	})
}

// PrintListFiltered prints entries that match a filter, optionally sorted by name.
func PrintListFiltered[T Entity](st *data.DataTable[T], sortname bool, filter func(data.Dataset[T]) bool) {
	var t []sortEntry
	if len(*st) == 0 {
		return
	}
	fmt.Println(GetTableHeader((*st)[0]))
	for _, set := range *st {
		if filter != nil && !filter(set) {
			continue
		}
		line := GetTableRow(set)
		if line != "" {
			t = append(t, sortEntry{Name: strings.ToUpper(set.Name), Text: line})
		}
	}
	if sortname {
		sort.Slice(t, func(i, j int) bool {
			return t[i].Name < t[j].Name
		})
	}
	for _, entry := range t {
		fmt.Println(entry.Text)
	}
}
