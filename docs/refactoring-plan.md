# Refactoring and improvement plan

Assessment date: 2026-10-28. The assessment below describes the original baseline.
Entity, dataset, table, and search rendering have since moved into `presentation/`;
data containers no longer require presentation methods. CLI output is checked
against 45 pre-refactoring snapshots, and `logic.go` has been split to satisfy
the source-file size limit. The remaining proposals are not yet implemented.

## Assessment

| Question | Finding | Recommendation |
| --- | --- | --- |
| Can storage be simpler? | JSON persistence is already small; complexity comes from global state, presentation methods, repeated table traversal, and scattered saves. | Keep JSON and the existing file format initially. Extract one file store and keep the model independent of storage and rendering. |
| Is Go current? | `go.mod` declares Go 1.21; the local compiler is Go 1.27.1. Latest stable is Go 1.27.3. | Adopt Go 1.27 as the supported baseline and validate with 1.27.3 in CI and locally. |
| Is architecture clean? | Packages have recognizable responsibilities, but dependencies mix models, terminal rendering, persistence, business operations, and interactive UI. | Use a small model/application/storage/presentation split with explicit dependencies. |
| Is the README complete? | The original examples omitted editing, deletion, detail commands, tag removal, and many list variants. Built-in help also omits commands and contains misleading descriptions. | Maintain a complete command reference now; later derive help and reference documentation from command metadata. |

Go versions were checked against the [official release history](https://go.dev/doc/devel/release).
The `go` directive sets the minimum Go version and language semantics; it does
not imply that the installed compiler is Go 1.21.

## Storage decision

Keep one JSON snapshot, generic typed tables, numeric references, and the current
path. Change `Dataset[T CustomData]` to an unconstrained data container once
rendering methods move out of the model. Do not flatten the hierarchy into an
untyped entity map or introduce an ORM.

Introduce a concrete JSON store with an explicit path and error-returning
load/save operations. It should decode into a fresh value and return errors
without partially replacing the active database. Only a missing file means a
new inventory; permission failures and malformed JSON must stop the operation.
Save through a temporary file in the same directory, sync and close it, then
rename it over the destination. Clean up temporary files on failure. Never
silently reset data or substitute another storage backend.

The application owns the database and ID allocation. Initially derive the next
ID from all loaded records, including deleted records, preserving the current
JSON schema. Reject ID exhaustion rather than wrapping. Commit a successful
mutation once, and only report success after persistence succeeds. Stage edits
so cancellation or failed persistence does not change active application state.

Atomic replacement prevents partially written snapshots; it does not prevent
two processes from overwriting each other's updates. Before promising concurrent
use, add an exclusive lock covering load, mutation, and save, with an explicit
busy error. Avoid holding a lock throughout an interactive form; use a revision
check before committing edited data, or explicitly reject concurrent sessions.

| Option | Benefits | Costs / decision |
| --- | --- | --- |
| JSON snapshot | No database dependency; easy to inspect and back up; retains existing data. | Whole-file writes and linear lookups. Recommended for the current small CLI. |
| SQLite | Transactions, constraints, indexes, and concurrency support. | Driver, schema, migrations, and import verification. Reconsider when concurrent writers or measured query costs justify it. |
| Key/value engine | Transactional key/value access. | Hierarchical joins and relationship validation still require application code. No demonstrated simplification here. |

## Architecture target

Keep `main.go` as the composition root. Four internal packages are sufficient:

- `internal/inventory`: plain models, ID allocation, queries, and validation.
- `internal/app`: inventory use cases; accepts database/store dependencies and returns results/errors.
- `internal/storage`: JSON file access, atomic persistence, and later locking if required.
- `internal/presentation`: CLI dispatch, terminal rendering, and explicit TUI forms in separate files.

Presentation invokes application operations. Application depends on inventory;
storage serializes inventory. Neither inventory nor application imports the TUI.
Define a small storage interface where the application consumes it only if
needed for tests; avoid per-entity repositories and a generic backend framework.

Replace reflection-based form construction with explicit typed field bindings.
Currently every `uint32` is assumed to be a reference dropdown and its options
depend on field order. Keep shared widgets, but validate typed edits before
applying them. Remove unused demo code and commented experiments during this step.

## Ordered implementation

### 1. Establish regression coverage and fix correctness issues

- Add tests before correcting category duplicate detection (`idx > 0` misses the first entry).
- Add tests for deleted-room/shelf name checks, which currently include deleted entries.
- Test that tagging an unknown target creates no tag; the current test explicitly expects an orphan tag and must change with this intentional behavior fix.
- Verify unchanged tags survive form submission: `taglist` currently starts empty and depends on the widget's change callback to initialize it.
- Define parent deletion behavior. Recommended default: reject deletion while live children or category references exist; tag deletion explicitly detaches references. Preserve existing deleted records on load and report inconsistencies without silently rewriting them.
- Specify warehouse scope explicitly: current warehouse validation does not scope generic selection or global lists. Preserve global lists; document and test any intended restrictions on creation and moves.

Acceptance: focused regression tests pass; all intentional behavior changes are documented.

### 2. Modernize the build baseline

- Update the Go directive/README together; use Go 1.27.3 for validation.
- Add CI for tests, vet, and build. Run race checks where supported.
- Review dependency upgrades individually, especially tview/tcell, and test UI behavior after upgrades.

Acceptance: reproducible build and passing checks on the supported toolchain.

### 3. Isolate and harden persistence

- Extract the JSON store, explicit path, and error handling described above.
- Test all-table round trips, existing-format compatibility, deleted records, ID restoration, malformed/truncated JSON, read/write failures, failed replacement, and preservation of the previous file.
- Replace process-wide HOME changes in tests with per-test temporary storage paths.
- Add a sanitized fixture; do not use the checked-in inventory snapshot as test data.

Acceptance: no ignored persistence errors; a failed save preserves the prior snapshot and never produces a success message.

### 4. Separate model, use cases, and presentation

- Remove `data.Db` and `id.IdSource` globals through explicit application ownership.
- Move `Show`, table formatting, colors, stdout writes, and selection dialogs out of models/use cases.
- Split the 586-line `logic.go` by responsibility; keep every source file below 500 lines and retain English function doc comments.
- Replace duplicated edit dispatch and repeated traversal where a small typed helper improves clarity.
- Replace reflection-based forms with typed edits and test save/cancel, reference selection, numeric validation, tag retention, and keyboard navigation using simulated terminal input.

Acceptance: independent application instances can be tested without changing globals, HOME, or stdout; inventory/application packages have no TUI dependency.

### 5. Align CLI, documentation, and repository hygiene

- Use command metadata for names, arguments, descriptions, and generated help/reference output; keep short command names compatible.
- Return nonzero exit status for invalid commands, arguments, and failed operations.
- Correct built-in help: missing `ac`, `ei`, `di`, `si`, and version aliases; `at`/`rt` accept tag names, not tag IDs; tag deletion only marks the tag deleted today.
- Add CLI integration tests for argument validation, aliases, command registration/help coverage, and an inventory workflow in a temporary directory.
- Document search fields, global listing scope, name/ID resolution, form controls, deletion rules, storage guarantees, and backup/restore.
- Review the tracked `bin/ginv` binary and `lgrtdata.json` inventory snapshot for removal from source control; replace examples with sanitized fixtures if needed.

Acceptance: every registered command appears in help and the README; documented behavior matches tests; `go test ./...`, `go vet ./...`, build, and diff checks pass. Commit each implementation session with a meaningful message.

## Verified baseline

`go test ./... -cover` and `go vet ./...` pass on local Go 1.27.1.
Statement coverage is 59.7% for `data` and 59.0% for `logic`; other packages have
no direct tests. Passing checks do not cover persistence failures, interactive
workflows, or the correctness issues listed above.
