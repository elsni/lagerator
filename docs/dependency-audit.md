# Dependency updates and feature review

Updated on 2026-10-28. Versions were resolved with `go get <module>@latest`,
`go get -u ./...`, and `go mod tidy`. All nine external modules imported by the
application and its tests now use the latest stable versions on their existing
module paths. The Go minimum is now 1.26.0, required by the updated x packages.

## Installed dependencies

| Module | Before | Installed | Module minimum Go |
| --- | --- | --- | --- |
| `github.com/gdamore/tcell/v2` | v2.7.4 | v2.13.10 | 1.24.0 |
| `github.com/rivo/tview` | April 2024 snapshot | v0.42.0 | 1.18 |
| `github.com/oleiade/reflections` | v1.0.1 | v1.1.0 | 1.22.6 |
| `github.com/gdamore/encoding` | v1.0.1 | v1.0.1 | 1.9 |
| `github.com/lucasb-eyer/go-colorful` | v1.2.0 | v1.4.1 | 1.12 |
| `github.com/rivo/uniseg` | v0.4.7 | v0.4.7 | 1.18 |
| `golang.org/x/sys` | v0.19.0 | v0.49.0 | 1.26.0 |
| `golang.org/x/term` | v0.19.0 | v0.47.0 | 1.26.0 |
| `golang.org/x/text` | v0.14.0 | v0.43.0 | 1.26.0 |

`go-runewidth` is no longer imported by the resolved application and was removed
from explicit requirements by `go mod tidy`. The full module graph still contains
upstream test/tool requirements that the application does not import; these were
not promoted into application dependencies merely to override upstream versions.

## Features incorporated

- **Literal inventory names:** tview's new `DropDown.SetUseStyleTags` and
  `List.SetUseStyleTags` preserve bracketed names such as `Box [red]` instead of
  interpreting them as styling. This is covered by simulated terminal rendering
  and dropdown keyboard-selection tests.
- **Terminal title:** tview's new `Application.SetTitle` identifies inventory
  dialogs as `Lagerator` in terminals that support window titles.
- **Tag completion:** the existing autocomplete API now completes the final
  comma-separated tag with case-insensitive matching. Earlier tags are preserved,
  deleted tags are excluded, arrows navigate suggestions, and Enter/Tab accepts.
  Suggestions are escaped for literal rendering; the original selected string
  is recovered by index because the callback receives text with style tags
  stripped. No tag is created until the editor is saved.
- **Bracketed paste:** enabled for dialogs so pasted multiline descriptions are
  delivered as text rather than interpreted as navigation or submission keys.
  This API already existed before the upgrade and is newly enabled in Lagerator.
- **Tag retention:** saving the editor reads the tag widget's current value,
  preserving unchanged tags without depending on initialization callbacks.

## Features reviewed but not added

- Reflections v1.1.0 adds JSON-tag-to-field lookup and assignable-value support.
  The current editor already has Go field names and writes matching concrete
  types. Neither feature requires an additional application feature now.
- tview adds richer field/dropdown styling. Existing colors remain appropriate;
  adding a theme system would extend this task without a demonstrated need.
- tcell image/graphics and clipboard-request APIs do not address the present
  inventory workflows. Existing terminal and Unicode fixes are received through
  the library upgrade without adding new application controls.

## Validation

The upgrade is checked with tests, race checks, vet, and build. The 45 CLI output
snapshots must remain unchanged. New tests run real widget handlers and a full
application event loop against tcell's simulated terminal: literal names,
completion acceptance, reference selection, numeric input rejection, editor
save/cancel, tag retention, terminal title, and multiline paste.

This is a dependency compatibility and feature review, not a security advisory audit.

## Upstream references

- [tcell releases](https://github.com/gdamore/tcell/releases)
- [tview dropdown API](https://github.com/rivo/tview/blob/v0.42.0/dropdown.go)
- [tview list API](https://github.com/rivo/tview/blob/v0.42.0/list.go)
- [tview application API](https://github.com/rivo/tview/blob/v0.42.0/application.go)
- [tview input-field API](https://github.com/rivo/tview/blob/v0.42.0/inputfield.go)
- [reflections v1.1.0 changes](https://github.com/oleiade/reflections/releases/tag/v1.1.0)
