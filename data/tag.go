package data

import (
	"fmt"

	"github.com/elsni/lagerator/terminal"
)

type Tag struct {
}

// GetTableRow returns the formatted row for a tag.
func (d Tag) GetTableRow(ownid uint32) string {
	return fmt.Sprintf("%4d", Db.CountTagOccurance(ownid))
}

// GetTableHeader returns the tag table header.
func (d Tag) GetTableHeader() string {
	return fmt.Sprintf("%-15s%s", "uses", terminal.ResetColor())
}

// Show prints tag details.
func (d Tag) Show() {
}

type TagTable = DataTable[Tag]

// Checks if a room exists in the given warehouse
// IsTagExistant checks if a tag name exists.
func IsTagExistant(wt *TagTable, name string) (bool, uint32) {
	for _, set := range *wt {
		if set.Name == name && !set.Deleted {
			return true, set.ID
		}
	}
	return false, 0
}
