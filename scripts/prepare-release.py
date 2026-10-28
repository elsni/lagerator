#!/usr/bin/env python3
"""Prepare vendored release sources, dependency licenses, and Gentoo metadata."""

import argparse
import gzip
import hashlib
from pathlib import Path
import re
import shutil
import subprocess
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parent.parent


def command(*args, cwd=ROOT):
    """Run a required build command and return its textual output."""
    return subprocess.check_output(args, cwd=cwd, text=True).strip()


def source_files():
    """List tracked source files while excluding inventory snapshots and binaries."""
    paths = command("git", "ls-files", "-z").split("\0")
    return [Path(path) for path in paths if path and path != "lgrtdata.json"
            and Path(path).parts[0] not in {"bin", "dist", "vendor"}]


def license_files(vendor):
    """Find the license and notice files retained by Go's vendor operation."""
    return sorted(path for path in vendor.rglob("*") if path.is_file()
                  and path.name.lower().startswith(("license", "copying", "notice")))


def archive_tree(source, destination, epoch):
    """Create a reproducible source tarball with stable ownership and timestamps."""
    with destination.open("wb") as raw:
        with gzip.GzipFile(filename="", mode="wb", fileobj=raw, mtime=epoch) as stream:
            with tarfile.open(fileobj=stream, mode="w") as archive:
                for path in [source, *sorted(source.rglob("*"))]:
                    info = archive.gettarinfo(str(path), str(path.relative_to(source.parent)))
                    info.uid = info.gid = 0
                    info.uname = info.gname = "root"
                    info.mtime = epoch
                    if path.is_file():
                        with path.open("rb") as content:
                            archive.addfile(info, content)
                    else:
                        archive.addfile(info)


def manifest_entry(path):
    """Describe a Gentoo distfile using its size, BLAKE2B, and SHA512 digests."""
    content = path.read_bytes()
    return (f"DIST {path.name} {len(content)} BLAKE2B "
            f"{hashlib.blake2b(content).hexdigest()} SHA512 "
            f"{hashlib.sha512(content).hexdigest()}\n")


def prepare(version, output, clean=False):
    """Generate release support artifacts for a stable or local snapshot version."""
    if not re.fullmatch(r"\d+\.\d+\.\d+(?:-dev)?", version):
        raise ValueError("expected a stable X.Y.Z version or an X.Y.Z-dev snapshot")
    if clean:
        if output != ROOT / "packaging/generated":
            raise ValueError("--clean only supports the generated packaging directory")
        if output.exists():
            shutil.rmtree(output)
    output.mkdir(parents=True, exist_ok=True)
    if any(output.iterdir()):
        raise ValueError(f"output directory must be empty: {output}")
    epoch = int(command("git", "show", "-s", "--format=%ct", "HEAD"))
    with tempfile.TemporaryDirectory(prefix="lagerator-source-") as temporary:
        source = Path(temporary) / f"lagerator-{version}"
        source.mkdir()
        for relative in source_files():
            target = source / relative
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(ROOT / relative, target)
        command("go", "mod", "vendor", "-o", str(source / "vendor"))
        licenses = source / "THIRD_PARTY_LICENSES"
        for path in license_files(source / "vendor"):
            target = licenses / path.relative_to(source / "vendor")
            target.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(path, target)
        go_license = licenses / "Go" / "LICENSE"
        go_license.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(ROOT / "packaging/licenses/Go-LICENSE", go_license)
        shutil.copytree(licenses, output / "licenses")
        tarball = output / f"lagerator_{version}_source.tar.gz"
        archive_tree(source, tarball, epoch)
    template = (ROOT / "packaging/gentoo/lagerator.ebuild.in").read_text()
    ebuild_version = version.removesuffix("-dev")
    (output / f"lagerator-{ebuild_version}.ebuild").write_text(
        template.replace("@VERSION@", version))
    (output / "Manifest").write_text(manifest_entry(tarball))
    print(f"Prepared release sources and Gentoo files in {output}")


def main():
    """Parse release parameters and generate the requested support artifacts."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", required=True)
    parser.add_argument("--output", type=Path, default=ROOT / "packaging/generated")
    parser.add_argument("--clean", action="store_true")
    args = parser.parse_args()
    prepare(args.version, args.output.resolve(), args.clean)


if __name__ == "__main__":
    main()
