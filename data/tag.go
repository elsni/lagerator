package data

type Tag struct {
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
