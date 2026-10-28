# Creating a release

Releases are built by GitHub Actions and remain **drafts** until you publish them
on GitHub. No GitHub token is needed for local snapshots. CI uses its scoped
`GITHUB_TOKEN` to create the draft; no personal access token is required.

## Release assets

For a tag such as `v1.0.13`, the workflow produces:

| Target | Architectures | Files |
| --- | --- | --- |
| Debian/Ubuntu | amd64, arm64 | `lagerator_1.0.13_<arch>.deb` |
| Linux | amd64, arm64 | `lagerator_1.0.13_linux_<arch>.tar.gz` |
| macOS | amd64 (Intel), arm64 (Apple Silicon) | `lagerator_1.0.13_darwin_<arch>.tar.gz` |
| Windows | amd64, arm64 | `lagerator_1.0.13_windows_<arch>.zip` |
| Source | all | `lagerator_1.0.13_source.tar.gz`, including vendored dependencies |
| Gentoo | amd64, arm64 | `lagerator-1.0.13.ebuild` and `Manifest` |
| Integrity | all | `checksums.txt` covering every downloadable artifact above |

The executable is `lgrt` (`lgrt.exe` on Windows). Archives contain documentation
and project/dependency licenses. The Debian package installs `/usr/bin/lgrt` and
licenses under `/usr/share/doc/lagerator/`. It never creates or replaces inventory
files. macOS and Windows downloads are command-line executables, not GUI installers;
code signing and notarization are not configured.

## Local validation

Install Go 1.26.0+, Python 3.9+, Git, Make, and **GoReleaser 2.19.0**. The
configuration is validated against that pinned GoReleaser version. Python helpers
use only the standard library. On macOS, GoReleaser is available through Homebrew:

```sh
brew install goreleaser
goreleaser --version
make release-check
make release-snapshot
```

`release-snapshot` cross-compiles all six executables, creates the archives and
two Debian packages, and verifies executable architecture, packaged license
files, source isolation, and checksums. It does not upload anything.

Binary packages appear in `dist/`; source, ebuild, Manifest, and dependency license
files are generated in `packaging/generated/`. Both directories are ignored by
Git. Each snapshot replaces the generated files. With no existing Git tags,
GoReleaser uses a `0.0.1-dev` snapshot version; this is not the real release version.
Snapshot ebuilds are for inspection and local tests, not published Gentoo installs.

Run application checks before creating a tag:

```sh
go test -count=1 -race ./...
go vet ./...
go mod verify
```

To check Debian installation on a disposable amd64 Debian/Ubuntu system, transfer
the matching `.deb` and run:

```sh
sudo apt install ./lagerator_1.0.13_amd64.deb
lgrt --version
```

## Create a draft on GitHub

1. Commit and push the release configuration to `main`. GitHub Actions must be
   enabled for the repository.
2. Choose an unused stable `vX.Y.Z` tag. The tag determines the release binary's
   version; commit and date are also embedded automatically. If you update the
   Makefile's `VERSION` and `args/args.go` default version, update the version
   expectations in `testdata/cli-output.json` too, run tests, and commit/push that
   version bump before tagging.
3. Create and push the annotated tag:

   ```sh
   git tag -a v1.0.13 -m "Release 1.0.13"
   git push origin v1.0.13
   ```

4. Watch **Actions → Draft release**. It runs native tests on Linux, macOS, and
   Windows using Go 1.27.x, then packages the release on Linux. The amd64 Debian
   package is installed and its version command is run before draft creation.
5. Open **Releases**, inspect the draft assets, test the downloads and Gentoo
   ebuild on their target systems, edit release notes, and choose **Publish release**.
   Gentoo's public source URL becomes usable after publication.

Only stable `vX.Y.Z` tags are accepted. The workflow has no prerelease or nightly
publishing path. If it fails before draft creation, fix the cause and rerun it.
If a draft already exists, review or delete that draft before retrying: the
workflow deliberately fails rather than overwriting an existing release.

## Gentoo overlay

The generated ebuild inherits `go-module` and builds the vendored source archive
with `-mod=vendor`. It requires Go 1.26.0+ and includes dependency license metadata.
It is intended for a personal overlay; it is not an official Gentoo tree package.

An overlay skeleton lives in `packaging/gentoo/overlay/`. On a Gentoo test system,
copy it to `/var/db/repos/lagerator`, then create
`/etc/portage/repos.conf/lagerator.conf` with:

```ini
[lagerator]
location = /var/db/repos/lagerator
masters = gentoo
auto-sync = no
```

For a published release, place its downloaded `.ebuild` and `Manifest` in:

```text
/var/db/repos/lagerator/app-misc/lagerator/
```

Accept the testing keyword for your architecture, then install:

```sh
emerge --ask app-misc/lagerator
```

For Portage validation, run `pkgcheck scan` in the overlay and use `ebuild` to
exercise unpack, test, compile, and install on a Gentoo machine. The Manifest
contains the release source archive's BLAKE2B and SHA512 digests. Do not substitute
GitHub's automatically generated source archive: it does not include vendored
dependencies and has different checksums.

A snapshot can be tested before publication by copying its generated source
archive into Gentoo's configured distfiles directory and its ebuild/Manifest into
the overlay. Portage then uses that local distfile without fetching the snapshot
URL. Native Gentoo/Portage validation still requires a Gentoo environment.

## Implementation

- `.goreleaser.yaml`: cross-compilation, archives, Debian metadata, checksums.
- `.github/workflows/ci.yml`: native checks on three operating systems.
- `.github/workflows/release.yml`: stable tags, tests, packaging, Debian installation,
  and draft creation after checks pass.
- `scripts/prepare-release.py`: vendored source, dependency licenses, ebuild, Manifest.
- `scripts/verify-release.py`: executable headers, Debian payloads, source safety,
  bundled licenses, and checksums.
- `packaging/gentoo/lagerator.ebuild.in`: versioned Gentoo template.

The custom source archive excludes the checked-in `lgrtdata.json` inventory and
`bin/` executables. It includes sanitized test fixtures. GitHub's automatic source
archives follow the repository contents and are separate from this curated asset.

References: [GoReleaser Actions](https://goreleaser.com/customization/ci/actions/),
[Debian packaging](https://goreleaser.com/customization/package/nfpm/),
[Gentoo go-module eclass](https://devmanual.gentoo.org/eclass-reference/go-module.eclass/index.html).
