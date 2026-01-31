package data

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/elsni/lagerator/id"
	"github.com/elsni/lagerator/terminal"
)

type CustomData interface {
	GetTableHeader() string
	GetTableRow(ownid uint32) string
	Show()
}

type Dataset[T CustomData] struct {
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
type Sortentry struct {
	Name string
	Text string
}

// NewDataset creates a dataset with a new id and timestamps.
func NewDataset[T CustomData](name string, data T) Dataset[T] {
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

// GetPrintName returns the display name with markers and a max width.
func (d Dataset[T]) GetPrintName(abb int) string {
	deleted := ""
	selected := ""
	if d.Deleted {
		deleted = "(deleted)"
	}
	if d.ID == Db.CurrentWarehouse {
		selected = "*"
	}
	s := fmt.Sprintf("%s%s%s", d.Name, selected, deleted)
	if len([]rune(s)) > abb {
		s = string([]rune(s)[:abb-1]) + "…"
	}
	return s
}

// GetTableRow returns the formatted row or an empty string when deleted.
func (d Dataset[T]) GetTableRow() string {
	if d.Deleted {
		return ""
	}
	return fmt.Sprintf("%5d %-30s ", d.ID, d.GetPrintName(30)) + d.Data.GetTableRow(d.ID)
}

// GetTableHeader prints the headline and returns the header row.
func (d Dataset[T]) GetTableHeader() string {
	t := strings.Split(fmt.Sprintf("%T", *new(T)), ".")
	fmt.Println(terminal.GetHeadlineText(fmt.Sprintf(" %s ", t[1])))
	return fmt.Sprintf("%s%5s %-30s %s", terminal.SetBgColor(terminal.COLORBLUE), "ID", "Name", d.Data.GetTableHeader())
}

// Show prints dataset details to stdout.
func (d Dataset[T]) Show() {
	t := strings.Split(fmt.Sprintf("%T", *new(T)), ".")
	fmt.Println(terminal.GetHeadlineText(fmt.Sprintf("%-12s %-30s", "Type", t[1])))
	fmt.Printf("%s %d\n", terminal.GetLabelText("Id"), d.ID)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Name"), d.GetPrintName(256))
	desclines := strings.Split(d.Description, "\n")
	dlabel := "Description"
	for idx, line := range desclines {
		if idx > 0 {
			dlabel = ""
		}
		fmt.Printf("%s %s\n", terminal.GetLabelText(dlabel), line)
	}

	fmt.Printf("%s %s\n", terminal.GetLabelText("Created"), terminal.GetTimeString(d.Created))
	fmt.Printf("%s %s\n", terminal.GetLabelText("Updated"), terminal.GetTimeString(d.Updated))
	d.Data.Show()
	fmt.Printf("%s %s\n", terminal.GetLabelText("Tags"), GetTagList(d.Tags))
}

type DataTable[T CustomData] []Dataset[T]

// NewDataTable creates an empty DataTable with default capacity.
func NewDataTable[T CustomData]() DataTable[T] {
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

// PrintList prints the table, optionally sorted by name.
func (st *DataTable[T]) PrintList(sortname bool) {
	st.PrintListFiltered(sortname, nil)
}

// PrintListByTag prints entries filtered by tag name.
func (st *DataTable[T]) PrintListByTag(tagname string, sortname bool) {
	tagid, found := Db.Tags.GetFirstOccurance(tagname)
	if !found {
		fmt.Printf("Unknown tag \"%s\"\n", tagname)
		return
	}
	st.PrintListFiltered(sortname, func(set Dataset[T]) bool {
		return slices.Contains(set.Tags, tagid)
	})
}

// PrintListByTagId prints entries filtered by tag id.
func (st *DataTable[T]) PrintListByTagId(tagid uint32, sortname bool) {
	st.PrintListFiltered(sortname, func(set Dataset[T]) bool {
		return slices.Contains(set.Tags, tagid)
	})
}

// PrintListFiltered prints entries that match a filter, optionally sorted by name.
func (st *DataTable[T]) PrintListFiltered(sortname bool, filter func(Dataset[T]) bool) {
	var t []Sortentry
	if len(*st) == 0 {
		return
	}
	fmt.Println((*st)[0].GetTableHeader())
	for _, set := range *st {
		if filter != nil && !filter(set) {
			continue
		}
		line := set.GetTableRow()
		if line != "" {
			t = append(t, Sortentry{Name: strings.ToUpper(set.Name), Text: line})
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
