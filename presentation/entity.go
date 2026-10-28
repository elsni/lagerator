package presentation

import "github.com/elsni/lagerator/data"

// Entity lists the inventory payloads supported by terminal rendering.
type Entity interface {
	data.Warehouse | data.Room | data.Shelf | data.Box | data.Item | data.Category | data.Tag
}

// entityTableRow dispatches rendering for the supported inventory payloads.
func entityTableRow[T Entity](value T, ownid uint32) string {
	switch d := any(value).(type) {
	case data.Warehouse:
		return tableRowWarehouse(d, ownid)
	case data.Room:
		return tableRowRoom(d, ownid)
	case data.Shelf:
		return tableRowShelf(d, ownid)
	case data.Box:
		return tableRowBox(d, ownid)
	case data.Item:
		return tableRowItem(d, ownid)
	case data.Category:
		return tableRowCategory(d, ownid)
	case data.Tag:
		return tableRowTag(d, ownid)
	}
	panic("unsupported inventory presentation type")
}

// entityTableHeader dispatches rendering for the supported inventory payloads.
func entityTableHeader[T Entity](value T) string {
	switch d := any(value).(type) {
	case data.Warehouse:
		return tableHeaderWarehouse(d)
	case data.Room:
		return tableHeaderRoom(d)
	case data.Shelf:
		return tableHeaderShelf(d)
	case data.Box:
		return tableHeaderBox(d)
	case data.Item:
		return tableHeaderItem(d)
	case data.Category:
		return tableHeaderCategory(d)
	case data.Tag:
		return tableHeaderTag(d)
	}
	panic("unsupported inventory presentation type")
}

// showEntity dispatches rendering for the supported inventory payloads.
func showEntity[T Entity](value T) {
	switch d := any(value).(type) {
	case data.Warehouse:
		showWarehouse(d)
		return
	case data.Room:
		showRoom(d)
		return
	case data.Shelf:
		showShelf(d)
		return
	case data.Box:
		showBox(d)
		return
	case data.Item:
		showItem(d)
		return
	case data.Category:
		showCategory(d)
		return
	case data.Tag:
		showTag(d)
		return
	}
	panic("unsupported inventory presentation type")
}
