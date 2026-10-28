package presentation

import (
	"fmt"
	"strings"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/terminal"
)

// GetPrintName returns the display name with markers and a max width.
func GetPrintName[T any](d data.Dataset[T], abb int) string {
	deleted := ""
	selected := ""
	if d.Deleted {
		deleted = "(deleted)"
	}
	if d.ID == data.Db.CurrentWarehouse {
		selected = "*"
	}
	s := fmt.Sprintf("%s%s%s", d.Name, selected, deleted)
	if len([]rune(s)) > abb {
		s = string([]rune(s)[:abb-1]) + "…"
	}
	return s
}

// GetTableRow returns the formatted row or an empty string when deleted.
func GetTableRow[T Entity](d data.Dataset[T]) string {
	if d.Deleted {
		return ""
	}
	return fmt.Sprintf("%5d %-30s ", d.ID, GetPrintName(d, 30)) + entityTableRow(d.Data, d.ID)
}

// GetTableHeader prints the headline and returns the header row.
func GetTableHeader[T Entity](d data.Dataset[T]) string {
	t := strings.Split(fmt.Sprintf("%T", *new(T)), ".")
	fmt.Println(terminal.GetHeadlineText(fmt.Sprintf(" %s ", t[1])))
	return fmt.Sprintf("%s%5s %-30s %s", terminal.SetBgColor(terminal.COLORBLUE), "ID", "Name", entityTableHeader(d.Data))
}

// Show prints dataset details to stdout.
func Show[T Entity](d data.Dataset[T]) {
	t := strings.Split(fmt.Sprintf("%T", *new(T)), ".")
	fmt.Println(terminal.GetHeadlineText(fmt.Sprintf("%-12s %-30s", "Type", t[1])))
	fmt.Printf("%s %d\n", terminal.GetLabelText("Id"), d.ID)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Name"), GetPrintName(d, 256))
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
	showEntity(d.Data)
	fmt.Printf("%s %s\n", terminal.GetLabelText("Tags"), data.GetTagList(d.Tags))
}

// GetPrintNameById returns the printable name for an id or a fallback.
func GetPrintNameById[T any](tbl *data.DataTable[T], id uint32, width int) string {
	if set, ok := tbl.GetPtr(id); ok {
		return GetPrintName(*set, width)
	}
	return "not found"
}
