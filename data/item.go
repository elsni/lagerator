package data

type Item struct {
	Location   string `json:"location"`
	Condition  string `json:"condition"`
	Amount     int    `json:"amount"`
	BoxId      uint32 `json:"boxId"`
	CategoryId uint32 `json:"categoryId"`
}

// GetShelfIdforItem returns the shelf id for an item id.
func GetShelfIdforItem(itemid uint32) uint32 {
	item, ok := Db.Items.GetPtr(itemid)
	if !ok {
		return 0
	}
	return GetShelfIdforBox(item.Data.BoxId)
}

type ItemTable = DataTable[Item]
