package args

import (
	"fmt"
	"os"
	"strconv"

	"github.com/elsni/lagerator/data"
	"github.com/elsni/lagerator/logic"
	"github.com/elsni/lagerator/terminal"
)

const appName = "Lagerator"
const appAuthor = "elsni"

var appVersion = "1.0.12"
var buildCommit = "unknown"
var buildDate = "unknown"

// requireArgs validates the minimum argument count.
func requireArgs(should int, args []string) bool {
	if len(args) < should {
		fmt.Println("Too few arguments")
		return false
	}
	return true
}

// ConvId parses a uint32 id from a string.
func ConvId(arg string) (uint32, error) {
	id, err := strconv.ParseUint(arg, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// parseID parses an id and prints a custom error on failure.
func parseID(arg string, errMsg string) (uint32, bool) {
	id, err := ConvId(arg)
	if err != nil {
		fmt.Println(errMsg)
		return 0, false
	}
	return id, true
}

// versionString returns the formatted version string.
func versionString() string {
	return fmt.Sprintf("%s %s (%s, %s) by %s", appName, appVersion, buildCommit, buildDate, appAuthor)
}

// ProcessArgs routes CLI arguments to command handlers.
func ProcessArgs() {
	args := os.Args[1:]
	if len(args) < 1 {
		PrintUsage()
		return
	}
	cmd := args[0]
	rest := args[1:]

	handlers := map[string]func([]string){
		"version":   func(_ []string) { fmt.Println(versionString()) },
		"--version": func(_ []string) { fmt.Println(versionString()) },
		"-v":        func(_ []string) { fmt.Println(versionString()) },
		"sww": func(a []string) {
			if requireArgs(1, a) {
				logic.SwitchWarehouse(a[0])
			}
		},
		"ac": func(a []string) {
			if requireArgs(1, a) {
				logic.AddCategory(a[0])
			}
		},
		"aw": func(a []string) {
			if requireArgs(1, a) {
				logic.AddWarehouse(a[0])
			}
		},
		"ar": func(a []string) {
			if requireArgs(1, a) {
				logic.AddRoomToCurrentWarehouse(a[0])
			}
		},
		"as": func(a []string) {
			if requireArgs(2, a) {
				logic.AddShelfToRoom(a[0], a[1])
			}
		},
		"ab": func(a []string) {
			if requireArgs(2, a) {
				logic.AddBoxToShelf(a[0], a[1])
			}
		},
		"ai": func(a []string) {
			if requireArgs(1, a) {
				logic.AddItems(a[0])
			}
		},
		"at": func(a []string) {
			if requireArgs(2, a) {
				id, ok := parseID(a[1], "Error: not an ID")
				if !ok {
					return
				}
				logic.AddTag(a[0], id)
			}
		},
		"rt": func(a []string) {
			if requireArgs(2, a) {
				id, ok := parseID(a[1], "Error: not an ID")
				if !ok {
					return
				}
				logic.RemoveTag(a[0], id)
			}
		},
		"mi": func(a []string) {
			if requireArgs(2, a) {
				id, ok := parseID(a[0], "Error: not an ID")
				if !ok {
					return
				}
				logic.MoveItem(id, a[1])
			}
		},
		"mb": func(a []string) {
			if requireArgs(2, a) {
				id, ok := parseID(a[0], "Error: not an ID")
				if !ok {
					return
				}
				logic.MoveBox(id, a[1])
			}
		},
		"lc":  func(_ []string) { data.Db.Categories.PrintList(false) },
		"lcs": func(_ []string) { data.Db.Categories.PrintList(true) },
		"lw":  func(_ []string) { data.Db.Warehouses.PrintList(false) },
		"lws": func(_ []string) { data.Db.Warehouses.PrintList(true) },
		"lr":  func(_ []string) { data.Db.Rooms.PrintList(false) },
		"lrs": func(_ []string) { data.Db.Rooms.PrintList(true) },
		"ls":  func(_ []string) { data.Db.Shelves.PrintList(false) },
		"lss": func(_ []string) { data.Db.Shelves.PrintList(true) },
		"lb":  func(_ []string) { data.Db.Boxes.PrintList(false) },
		"lbs": func(_ []string) { data.Db.Boxes.PrintList(true) },
		"li":  func(_ []string) { data.Db.Items.PrintList(false) },
		"lis": func(_ []string) { data.Db.Items.PrintList(true) },
		"lt":  func(_ []string) { data.Db.Tags.PrintList(true) },
		"lic": func(a []string) {
			if requireArgs(1, a) {
				logic.PrintItemsOfCategory(a[0], false)
			}
		},
		"lics": func(a []string) {
			if requireArgs(1, a) {
				logic.PrintItemsOfCategory(a[0], true)
			}
		},
		"lib": func(a []string) {
			if requireArgs(1, a) {
				logic.PrintItemsOfBox(a[0], false)
			}
		},
		"libs": func(a []string) {
			if requireArgs(1, a) {
				logic.PrintItemsOfBox(a[0], true)
			}
		},
		"lit": func(a []string) {
			if requireArgs(1, a) {
				idx := logic.SelectSet[data.Tag](&data.Db.Tags, a[0], "Tags", "list")
				if idx == -1 {
					fmt.Println("Unknown Tag " + a[0])
					return
				}
				data.Db.Items.PrintListByTagId(data.Db.Tags[idx].ID, false)
			}
		},
		"lits": func(a []string) {
			if requireArgs(1, a) {
				idx := logic.SelectSet[data.Tag](&data.Db.Tags, a[0], "Tags", "list")
				if idx == -1 {
					fmt.Println("Unknown Tag " + a[0])
					return
				}
				data.Db.Items.PrintListByTagId(data.Db.Tags[idx].ID, true)
			}
		},
		"e": func(a []string) {
			if requireArgs(1, a) {
				id, ok := parseID(a[0], "Error: not an ID")
				if !ok {
					return
				}
				logic.EditAny(id)
			}
		},
		"ec": func(a []string) {
			if requireArgs(1, a) {
				logic.EditSet[data.Category](&data.Db.Categories, a[0], "Category")
			}
		},
		"ew": func(a []string) {
			if requireArgs(1, a) {
				logic.EditSet[data.Warehouse](&data.Db.Warehouses, a[0], "Warehouse")
			}
		},
		"er": func(a []string) {
			if requireArgs(1, a) {
				logic.EditSet[data.Room](&data.Db.Rooms, a[0], "Room")
			}
		},
		"es": func(a []string) {
			if requireArgs(1, a) {
				logic.EditSet[data.Shelf](&data.Db.Shelves, a[0], "Shelf")
			}
		},
		"eb": func(a []string) {
			if requireArgs(1, a) {
				logic.EditSet[data.Box](&data.Db.Boxes, a[0], "Box")
			}
		},
		"ei": func(a []string) {
			if requireArgs(1, a) {
				logic.EditSet[data.Item](&data.Db.Items, a[0], "Item")
			}
		},
		"d": func(a []string) {
			if requireArgs(1, a) {
				id, ok := parseID(a[0], "Error: not an ID")
				if !ok {
					return
				}
				logic.DeleteAny(id)
			}
		},
		"dc": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Category](&data.Db.Categories, a[0], "Category")
			}
		},
		"dw": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Warehouse](&data.Db.Warehouses, a[0], "Warehouse")
			}
		},
		"dr": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Room](&data.Db.Rooms, a[0], "Room")
			}
		},
		"ds": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Shelf](&data.Db.Shelves, a[0], "Shelf")
			}
		},
		"db": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Box](&data.Db.Boxes, a[0], "Box")
			}
		},
		"di": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Item](&data.Db.Items, a[0], "Item")
			}
		},
		"dt": func(a []string) {
			if requireArgs(1, a) {
				logic.DeleteSet[data.Tag](&data.Db.Tags, a[0], "Tag")
			}
		},
		"s": func(a []string) {
			if requireArgs(1, a) {
				id, ok := parseID(a[0], "Error: Not an Id")
				if !ok {
					return
				}
				logic.ShowAny(id)
			}
		},
		"sc": func(a []string) {
			if requireArgs(1, a) {
				logic.ShowSet[data.Category](&data.Db.Categories, a[0], "Category")
			}
		},
		"sw": func(a []string) {
			if requireArgs(1, a) {
				logic.ShowSet[data.Warehouse](&data.Db.Warehouses, a[0], "Warehouse")
			}
		},
		"ss": func(a []string) {
			if requireArgs(1, a) {
				logic.ShowSet[data.Shelf](&data.Db.Shelves, a[0], "Shelf")
			}
		},
		"sr": func(a []string) {
			if requireArgs(1, a) {
				logic.ShowSet[data.Room](&data.Db.Rooms, a[0], "Room")
			}
		},
		"sb": func(a []string) {
			if requireArgs(1, a) {
				logic.ShowSet[data.Box](&data.Db.Boxes, a[0], "Box")
			}
		},
		"si": func(a []string) {
			if requireArgs(1, a) {
				logic.ShowSet[data.Item](&data.Db.Items, a[0], "Item")
			}
		},
		"f": func(a []string) {
			if requireArgs(1, a) {
				data.Db.FindItem(a[0], false)
			}
		},
		"fs": func(a []string) {
			if requireArgs(1, a) {
				data.Db.FindItem(a[0], true)
			}
		},
	}

	if handler, ok := handlers[cmd]; ok {
		handler(rest)
		return
	}
	fmt.Printf("Unknown operation \"%s\"\n", cmd)
}

// PrintUsage prints CLI usage text.
func PrintUsage() {
	fmt.Println(terminal.GetHeadlineText(appName + " - console inventory management"))
	fmt.Println("usage: lgrt <operation> [object|list|id|name|searchstring]")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Operations: "))
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Add objects:"))
	fmt.Println("aw  <name>                       Add a warehouse")
	fmt.Println("ar  <name>                       Add a room to the current warehouse")
	fmt.Println("as  <shelfname> <roomname or ID> Add a shelf to a room")
	fmt.Println("ab  <boxname> <shelfname or ID>  Add a box to a shelf")
	fmt.Println("ai  <boxname or ID>              add items interactively to a specific box")
	fmt.Println("sww <name>                       switch to warehouse")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("List objects:"))
	fmt.Println("lc     list categories")
	fmt.Println("lcs    list categories sorted by name")
	fmt.Println("lw     list warehouses")
	fmt.Println("lws    list warehouses sorted by name")
	fmt.Println("lr     list rooms")
	fmt.Println("lrs    list rooms sorted by name")
	fmt.Println("ls     list shelves")
	fmt.Println("lss    list shelves sorted by name")
	fmt.Println("lb     list boxes")
	fmt.Println("lbs    list boxes sorted by name")
	fmt.Println("li     list items")
	fmt.Println("lis    list items sorted by name")
	fmt.Println("lic    <category name or id> list items of a category")
	fmt.Println("lics   <category name or id> list items of a category sorted by name")
	fmt.Println("lib    <box name or id>      list items in a box")
	fmt.Println("libs   <box name or id>      list items in a box sorted by name")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Edit objects:"))
	fmt.Println("e <id>     edit any object")
	fmt.Println("ec <name>  edit category")
	fmt.Println("ew <name>  edit warehouse")
	fmt.Println("er <name>  edit room")
	fmt.Println("es <name>  edit shelf")
	fmt.Println("eb <name>  edit box")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Reorganize:"))
	fmt.Println("mi <itemid> <box name or id>   move item to another box")
	fmt.Println("mb <boxid>  <shelf name or id> move box to another shelf")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Delete objects:"))
	fmt.Println("d <id>     delete object")
	fmt.Println("dc <name>  delete category")
	fmt.Println("dw <name>  delete warehouse")
	fmt.Println("dr <name>  delete room")
	fmt.Println("ds <name>  delete shelf")
	fmt.Println("db <name>  delete box")
	fmt.Println("deleting objects won't break integrity, since they are only marked as deleted.")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Tagging:"))
	fmt.Println("Tags are automatically added on first use.")
	fmt.Println("lt                              list tags and number of uses")
	fmt.Println("at   <tagname or id> <objectId> add tag to any object")
	fmt.Println("rt   <tagname or id> <objectId> remove Tag from object")
	fmt.Println("dt   <tagname or id>            delete tag from all tagged objects")
	fmt.Println("lit  <tagname or id>            list items by tag")
	fmt.Println("lits <tagname or id>            list items by tag sorted by name")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Show object details:"))
	fmt.Println("s <id>     show object")
	fmt.Println("sc <name>  show category")
	fmt.Println("sw <name>  show warehouse")
	fmt.Println("sr <name>  show room")
	fmt.Println("ss <name>  show shelf")
	fmt.Println("sb <name>  show box")
	fmt.Println()
	fmt.Println(terminal.GetHeadlineText("Find items:"))
	fmt.Println("f  <searchstring>  list sorted by Id")
	fmt.Println("fs <searchstring>  list sorted by name")

}
