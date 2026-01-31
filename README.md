# Lagerator

Terminal-based personal inventory management in Go.

`Lagerator` keeps your stuff organized in a hierarchy:
Warehouse -> Room -> Shelf -> Box -> Item.
Items can also be assigned to categories and tagged.

## Features
- Hierarchical inventory structure (Warehouse/Room/Shelf/Box/Item)
- Categories and tags
- Text search across items
- TUI editor for adding/editing entries

## Installation

### Requirements
- Go 1.21+

### Build from source
```
git clone https://github.com/elsni/lagerator.git
cd Lagerator
make build
```

The binary is created at `bin/lgrt`.

### Install to /usr/local/bin
```
make install
```

## Usage

Run:
```
lgrt <operation> [args]
```

Show version/author:
```
lgrt --version
```

Build with version metadata:
```
make build VERSION=0.1.0
```

### Getting started (empty database)
When starting with an empty database, create the hierarchy top‑down:
1) Create a warehouse
2) Switch to it
3) Add rooms, shelves, boxes
4) Add items to a box (interactive)

Example (minimal workflow):
```
lgrt aw "Home"
lgrt sww "Home"
lgrt ar "Basement"
lgrt as "Shelf A" "Basement"
lgrt ab "Box 1" "Shelf A"
lgrt ai "Box 1"
```

Example (with categories + tags):
```
lgrt ac "Tools"
lgrt ai "Box 1"
lgrt at "fragile" <itemId>
```

Examples:
```
# add and select a warehouse
lgrt aw "Home"
lgrt sww "Home"

# add structure
lgrt ar "Basement"
lgrt as "Shelf A" "Basement"
lgrt ab "Box 1" "Shelf A"

# add items (interactive TUI)
lgrt ai "Box 1"

# list
lgrt lw
lgrt lr
lgrt ls
lgrt lb
lgrt li

# find items
lgrt f "camera"
lgrt fs "camera"

# move
lgrt mi <itemId> <box name|id>
lgrt mb <boxId> <shelf name|id>
```

For the full command list, run `lgrt` without arguments.

## Data storage
Data is stored in:
```
~/.lgrt/lgrtdata.json
```

## Screenshots
![Item edit form](screenshots/lgrt_edit.png)
![Search results](screenshots/lgrt_find.png)
![Item details](screenshots/lgrt_show.png)

## Development
Run tests:
```
make test
```

## License
GPL-3.0-or-later
