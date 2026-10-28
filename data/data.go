package data

import (
	"sort"
	"strings"
	"time"

	"github.com/elsni/lagerator/id"
)

type Dataset[T any] struct {
	ID          uint32   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Created     int64    `json:"created"`
	Updated     int64    `json:"updated"`
	Deleted     bool     `json:"deleted"`
	Tags        []uint32 `json:"tags"`
	Data        T        `json:"data"`
}

type Listentry struct {
	Id   uint32
	Name string
}

// NewDataset creates a dataset with a new id and timestamps.
func NewDataset[T any](name string, data T) Dataset[T] {
	d := Dataset[T]{
		ID:      id.IdSource.GetNewId(),
		Name:    name,
		Data:    data,
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
		Deleted: false,
		Tags:    make([]uint32, 0, 128),
	}
	return d
}

type DataTable[T any] []Dataset[T]

// NewDataTable creates an empty DataTable with default capacity.
func NewDataTable[T any]() DataTable[T] {
	return make(DataTable[T], 0, 128)
}

// GetNames returns all non-deleted names and their ids.
func (dt *DataTable[T]) GetNames() []Listentry {
	var names []Listentry
	for _, set := range *dt {
		if !set.Deleted {
			names = append(names, Listentry{Id: set.ID, Name: set.Name})
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].Name < names[j].Name
	})
	return names
}

// GetFirstOccurance returns the first case-insensitive name match and its id.
func (dt *DataTable[T]) GetFirstOccurance(name string) (uint32, bool) {
	for _, set := range *dt {
		if !set.Deleted && strings.EqualFold(set.Name, name) {
			return set.ID, true
		}
	}
	return 0, false
}

// GetNamesRaw returns sorted names of all non-deleted entries.
func (dt *DataTable[T]) GetNamesRaw() []string {
	var names []string
	for _, set := range *dt {
		if !set.Deleted {
			names = append(names, set.Name)
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names
}

// GetSetsByName returns all non-deleted datasets with the given name.
func (dt *DataTable[T]) GetSetsByName(name string) []Dataset[T] {
	list := make([]Dataset[T], 0, 128)
	for _, set := range *dt {
		if strings.EqualFold(set.Name, name) && !set.Deleted {
			list = append(list, set)
		}
	}
	return list
}

// Add appends a dataset, ensuring a default name.
func (dt *DataTable[T]) Add(set Dataset[T]) {
	if set.Name == "" {
		set.Name = "Unnamed"
	}
	*dt = append(*dt, set)
}

// AddSimple adds a dataset with a name, new id, and default data.
func (dt *DataTable[T]) AddSimple(name string) uint32 {
	data := new(T)
	set := Dataset[T]{
		ID:      id.IdSource.GetNewId(),
		Name:    name,
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
		Data:    *data,
		Deleted: false,
	}
	*dt = append(*dt, set)
	return set.ID
}

// GetPtr returns a pointer to the dataset inside the table.
func (dt *DataTable[T]) GetPtr(id uint32) (*Dataset[T], bool) {
	idx := dt.GetIdx(id)
	if idx < 0 {
		return nil, false
	}
	return &(*dt)[idx], true
}

// GetDataByName returns index and id for an exact name match.
func (dt *DataTable[T]) GetDataByName(name string) (int, uint32) {
	for i, set := range *dt {
		if set.Name == name && !set.Deleted {
			return i, set.ID
		}
	}
	return -1, 0
}

// GetIdx returns the index for an id, or -1 if not found.
func (st *DataTable[T]) GetIdx(id uint32) int {
	for i, set := range *st {
		if set.ID == id && !set.Deleted {
			return i
		}
	}
	return -1
}

// Delete marks an item as deleted and updates the timestamp.
func (st *DataTable[T]) Delete(id uint32) {
	idx := st.GetIdx(id)
	(*st)[idx].Deleted = true
	(*st)[idx].Updated = time.Now().Unix()
}
