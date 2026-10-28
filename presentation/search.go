package presentation

import (
	"fmt"
	"strings"

	"github.com/elsni/lagerator/data"
)

// FindItem prints matching items for a search string.
func FindItem(searchstring string, sortname bool) {
	if len(data.Db.Items) == 0 {
		fmt.Println("no data")
		return
	}
	s := strings.ToLower(searchstring)
	PrintListFiltered(&data.Db.Items, sortname, func(item data.Dataset[data.Item]) bool {
		return strings.Contains(strings.ToLower(item.Name), s) ||
			strings.Contains(strings.ToLower(item.Description), s) ||
			strings.Contains(strings.ToLower(item.Data.Location), s)
	})
}
