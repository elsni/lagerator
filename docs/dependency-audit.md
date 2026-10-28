# Dependency version audit

Checked on 2026-10-28 with `go list -m -u -json all` and
`go list -m -f '{{.Path}} {{.Version}} go={{.GoVersion}}' <module>@latest`.
These are available versions for the existing module paths; compatibility with
Lagerator has not been tested. No dependency versions were changed.

## Direct dependencies

| Module | Current | Latest | Latest module's minimum Go |
| --- | --- | --- | --- |
| `github.com/gdamore/tcell/v2` | v2.7.4 | v2.13.10 | 1.24.0 |
| `github.com/rivo/tview` | v0.0.0-20240406141410-79d4cc321256 | v0.42.0 | 1.18 |
| `github.com/oleiade/reflections` | v1.0.1 | v1.1.0 | 1.22.6 |

The Go minimums above belong to each module; the resolved dependency graph may
require a newer Go version. Upstream references:
[tcell releases](https://github.com/gdamore/tcell/releases),
[tview tags](https://github.com/rivo/tview/tags), and
[reflections releases](https://github.com/oleiade/reflections/releases).

## Indirect dependencies explicitly required in go.mod

| Module | Current | Latest | Latest module's minimum Go |
| --- | --- | --- | --- |
| `github.com/gdamore/encoding` | v1.0.1 | v1.0.1 | — |
| `github.com/lucasb-eyer/go-colorful` | v1.2.0 | v1.4.1 | 1.12 |
| `github.com/mattn/go-runewidth` | v0.0.15 | v0.0.31 | 1.23 |
| `github.com/rivo/uniseg` | v0.4.7 | v0.4.7 | — |
| `golang.org/x/sys` | v0.19.0 | v0.49.0 | 1.26.0 |
| `golang.org/x/term` | v0.19.0 | v0.47.0 | 1.26.0 |
| `golang.org/x/text` | v0.14.0 | v0.43.0 | 1.26.0 |

Eight of the ten explicitly required modules have newer versions. The full
module graph also contains older upstream test/tool requirements; do not add
all of them as explicit application dependencies merely to update them.

## Proposed upgrade sequence

1. Raise the declared Go baseline from 1.21 before adopting libraries that
   require newer Go. Updating all modules listed here requires at least Go 1.26.
2. Update tview and tcell together, then test input handling, dropdowns, selection
   lists, editor save/cancel, and terminal layout.
3. Update reflections independently and verify reflection-based field editing.
4. Resolve indirect requirements and run `go mod tidy`; inspect the resulting
   module graph rather than retaining unnecessary explicit pins.
5. Run tests (including CLI snapshots), race checks, vet, and build. Perform an
   interactive TUI check because CLI snapshots do not exercise the forms.

This audit checks version availability and declared Go requirements, not security
advisories or upgrade compatibility.
