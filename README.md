# Lagerator

A terminal-based personal inventory manager written in Go. Track what you own,
where it is stored, and how much you have, using short CLI commands and an
interactive editor.

```text
Warehouse → Room → Shelf → Box → Item
```

Items have a quantity, condition, location, and optional category. Tags and
case-insensitive search help you find things across your inventory. Data is
stored locally in a single JSON file.

[Install](#installation) · [Quick start](#quick-start) ·
[Editor controls](#interactive-editor) · [Commands](#command-reference) ·
[Data and backups](#data-and-backups) · [Releases](#releases) ·
[Development](#development)

## Installation

You need **Go 1.26.0 or newer**, Git, and Make. Interactive forms require a terminal.

```sh
git clone https://github.com/elsni/lagerator.git
cd lagerator
make build
./bin/lgrt --version
```

The executable is `bin/lgrt`. To use `lgrt` in the examples below, add the build
directory to your current shell's PATH:

```sh
export PATH="$PWD/bin:$PATH"
```

Alternatively, install it to `/usr/local/bin`:

```sh
make install
```

This builds the release executable and uses `sudo` to copy it into place.
Installation does not import inventory data.

## Quick start

For a new inventory, create the hierarchy from the top down:

```sh
lgrt aw Home                       # Create a warehouse
lgrt sww Home                      # Make it active
lgrt ar Basement                   # Add a room
lgrt as "Shelf A" Basement         # Add a shelf to that room
lgrt ab "Box 1" "Shelf A"           # Add a box to that shelf
lgrt ac Tools                      # Create an item category
lgrt ai "Box 1"                    # Open the item editor
```

Enter a name, choose the category, and press **Ctrl+S** to save. The add-item
editor opens another entry after each save, reusing the previous category and
location. Press **Ctrl+Q** when you have finished adding items; earlier saves
remain stored.

List your items and search for one:

```sh
lgrt li
lgrt f camera
lgrt lib "Box 1"
lgrt lic Tools
```

Use IDs from list output for further operations. For example, if an item's ID
is `42`:

```sh
lgrt s 42                          # Show details
lgrt e 42                          # Edit the item
lgrt at fragile 42                 # Attach a tag, creating it if necessary
lgrt lit fragile                   # Find items with that tag
lgrt rt fragile 42                 # Remove the attachment
```

Quote names or search text that contain spaces. The example ID `42` must be
replaced with an ID from your own inventory.

## Interactive editor

The editor starts in **Name**. Common fields are name, description, and tags;
item fields also include location, condition, amount, box, and category.

| Control | Action |
| --- | --- |
| **Tab** / **Shift+Tab** | Move forward/backward between fields, including the multiline Description |
| **Page Down** / **Page Up** | Move forward/backward between fields |
| **Ctrl+S** or **Ok** | Save the entry |
| **Ctrl+Q**, **Esc**, or **Quit** | Cancel the current entry |
| **Ctrl+U** | Clear the text in the focused input field |
| **↑** / **↓** | Navigate autocomplete suggestions |
| **Enter** / **Tab** in suggestions | Accept the selected completion |

Reference fields such as **Box** and **Category** support case-insensitive
substring search. Clear the existing value, type part of a name, then accept a
suggestion. Choices include IDs to distinguish duplicate names. Incomplete or
invalid choices block saving and focus the field with an explanation; other
entered values stay intact. Untouched references are preserved.

**Tags** accepts a comma-separated list, for example `fragile, electronics`.
Start typing a tag prefix to complete an existing tag, including after a comma.
You can also enter a new tag; it is created when you save. Tab accepts an open
completion before moving to the next field on a subsequent press.

Bracketed names such as `Box [red]` display literally. Dialogs support bracketed
paste, including multiline descriptions, and set the terminal title to
`Lagerator` where supported.

## Command reference

```text
lgrt <operation> [arguments]
```

Run `lgrt` without arguments for built-in help. This reference also includes
commands currently omitted from that help.

### Names, IDs, and scope

- `name-or-id` tries a case-insensitive name match first, then a numeric ID.
- Ambiguous names open a selection dialog when an operation needs one record.
  Type-specific detail commands show all matching records.
- `sww` selects a warehouse by its exact, case-sensitive name.
- Lists and search span all warehouses. Selecting a warehouse does not filter
  those results; it determines the parent of newly created rooms.
- Search matches substrings in item **name**, **description**, and **location**,
  ignoring case. It does not search category or tag names.
- Unsorted lists preserve stored order. For paired list/search commands, the
  second variant sorts by name. `lt` always sorts tags by name.

| Command | Arguments | Function |
| --- | --- | --- |
| `version`, `--version`, `-v` | none | Show version, build metadata, and author |
| `aw` | `name` | Create warehouse |
| `sww` | `name` | Select active warehouse by exact name |
| `ar` | `name` | Create room in active warehouse |
| `as` | `name room-name-or-id` | Create shelf in selected room |
| `ab` | `name shelf-name-or-id` | Create box on selected shelf |
| `ai` | `box-name-or-id` | Add items repeatedly through the editor |
| `ac` | `name` | Create category |
| `lw`, `lws` | none | List warehouses; second command sorts by name |
| `lr`, `lrs` | none | List rooms; second sorts by name |
| `ls`, `lss` | none | List shelves; second sorts by name |
| `lb`, `lbs` | none | List boxes; second sorts by name |
| `li`, `lis` | none | List items; second sorts by name |
| `lc`, `lcs` | none | List categories; second sorts by name |
| `lt` | none | List tags sorted by name, with usage counts |
| `lic`, `lics` | `category-name-or-id` | List category items; second sorts by name |
| `lib`, `libs` | `box-name-or-id` | List box items; second sorts by name |
| `lit`, `lits` | `tag-name-or-id` | List tagged items; second sorts by name |
| `f`, `fs` | `text` | Search items; second sorts by name |
| `e` | `id` | Edit category, warehouse, room, shelf, box, or item |
| `ec`, `ew`, `er`, `es`, `eb`, `ei` | `name-or-id` | Edit the corresponding record type |
| `d` | `id` | Confirm deletion of category, warehouse, room, shelf, box, or item |
| `dc`, `dw`, `dr`, `ds`, `db`, `di`, `dt` | `name-or-id` | Confirm deletion of the corresponding type, including tags |
| `s` | `id` | Show category, warehouse, room, shelf, box, or item details |
| `sc`, `sw`, `sr`, `ss`, `sb`, `si` | `name-or-id` | Show details for the corresponding type |
| `at` | `tag-name object-id` | Attach a tag, creating it on first use |
| `rt` | `tag-name object-id` | Remove a tag from one object |
| `mi` | `item-id box-name-or-id` | Move item to box |
| `mb` | `box-id shelf-name-or-id` | Move box to shelf |

### Moving records

Use the source record's ID and the destination's name or ID:

```text
lgrt mi <item-id> <box-name-or-id>
lgrt mb <box-id> <shelf-name-or-id>
```

### Deleting records

Deletion asks for confirmation and marks the selected record as deleted. There
is no restore command or automatic deletion of children. Deleting a parent can
leave unresolved paths for its children.

Deleting a tag hides it but leaves its IDs on tagged records. Use `rt` to remove
an attachment while the tag is still active.

## Data and backups

Lagerator loads `~/.lgrt/lgrtdata.json` at startup and rewrites it after mutations.
A missing file starts an empty inventory; the file is created on the first save.
The selected warehouse is stored with the inventory.

Close Lagerator before copying or restoring the file:

```sh
# Back up an existing inventory; -i asks before replacing an existing backup.
cp -i ~/.lgrt/lgrtdata.json ~/lgrtdata-backup.json

# Restore a backup; this replaces the current inventory after confirmation.
cp -i ~/lgrtdata-backup.json ~/.lgrt/lgrtdata.json
```

Current writes are not atomic, write errors are not reported, and simultaneous
writers are not coordinated. Use one Lagerator process at a time and keep backups.

### Importing from ginventory

The ginventory JSON format is compatible with Lagerator. To import from a sibling
checkout, close both applications, back up any existing Lagerator inventory, and
copy the data file:

```sh
mkdir -p ~/.lgrt
cp -i ../ginventory/ginvdata.json ~/.lgrt/lgrtdata.json
lgrt lw
lgrt li
```

This is a one-time copy. The original ginventory file remains unchanged; the two
files are not synchronized.

## Screenshots

### Item editor

![Item edit form](screenshots/lgrt_edit.png)

### Search results

![Search results](screenshots/lgrt_find.png)

### Item details

![Item details](screenshots/lgrt_show.png)

## Releases

Packaged downloads are listed on [GitHub Releases](https://github.com/elsni/lagerator/releases).

Release automation builds Debian packages and Linux, macOS, and Windows binaries
for amd64 and arm64, plus vendored source and a Gentoo ebuild. Stable version tags
create a **draft GitHub release** after native tests and package checks pass.

Existing tags can also be released with **Actions → Draft release → Run workflow**.
See [Creating a release](docs/releases.md) for local snapshot builds, tagging,
installation, Gentoo overlay setup, and publishing. Artifacts are generated by
the workflow; setting up the configuration does not publish a release.

## Development

### Build and checks

```sh
make build                        # Build bin/lgrt with version metadata
make release                      # Build a stripped executable
make test                         # Run all tests
go test -count=1 -race ./...       # Run uncached tests with race checks
go vet ./...                      # Run static diagnostics
```

Override the version when building a release:

```sh
make build VERSION=1.0.13
```

Builds also embed the Git commit and build date, shown by `lgrt --version`.
`make run` starts the program through Go and prints help without arguments.

### Project structure

| Path | Responsibility |
| --- | --- |
| `main.go` | Load the database and dispatch commands |
| `args/` | Parse CLI arguments, route commands, and print help |
| `data/` | Plain entities, generic datasets/tables, relationships, and JSON persistence |
| `logic/` | Inventory operations, selection, tagging, and editor workflows |
| `presentation/` | Render names, tables, details, filtered lists, and search results |
| `ui/` | Interactive forms, searchable references, completion, and confirmation dialogs |
| `id/`, `terminal/`, `loggi/` | ID allocation, terminal formatting, and debug logging |
| `testdata/` | Sanitized inventory and CLI output snapshots |

The application currently shares a process-wide database and ID source.
Records use numeric parent references and include timestamps, tags, and a
soft-deletion flag. Rendering is separate from the model types.

### Tests

Data and logic tests isolate file writes with temporary home directories.
Presentation tests cover markers, Unicode truncation, ordering, and filtering.
CLI integration tests compare 45 cases against pre-refactoring output snapshots
and verify that read-only commands leave inventory files unchanged. UI tests use
a simulated terminal to exercise navigation, literal names, completion,
save/cancel, tag retention, and multiline paste.

The maintainer targets `make getdata` and `make putdata` copy inventory between
the repository and `~/.lgrt`. `putdata` replaces the active inventory after moving
the previous file to `~/.lgrt/lgrtdata.json.bak`.

Further details:

- [Refactoring plan](docs/refactoring-plan.md): architecture assessment and remaining work.
- [Dependency review](docs/dependency-audit.md): library versions, incorporated features, and validation.

## License

[GPL-3.0-or-later](LICENSE). Contributors are listed in [AUTHORS](AUTHORS).
