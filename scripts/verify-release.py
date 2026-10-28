#!/usr/bin/env python3
"""Validate release archives, Debian payloads, source isolation, and checksums."""

import argparse
import hashlib
import io
from pathlib import Path
import struct
import tarfile
import zipfile


def binary_architecture(content, system):
    """Read an executable header and return the architecture encoded by its toolchain."""
    if system == "linux":
        if content[:4] != b"\x7fELF":
            raise ValueError("Linux executable is not ELF")
        return {62: "amd64", 183: "arm64"}[struct.unpack_from("<H", content, 18)[0]]
    if system == "darwin":
        if content[:4] != b"\xcf\xfa\xed\xfe":
            raise ValueError("macOS executable is not Mach-O")
        return {0x01000007: "amd64", 0x0100000C: "arm64"}[struct.unpack_from("<I", content, 4)[0]]
    if system == "windows":
        if content[:2] != b"MZ":
            raise ValueError("Windows executable is not PE")
        offset = struct.unpack_from("<I", content, 0x3C)[0]
        if content[offset:offset+4] != b"PE\0\0":
            raise ValueError("Windows PE signature is invalid")
        return {0x8664: "amd64", 0xAA64: "arm64"}[struct.unpack_from("<H", content, offset+4)[0]]
    raise ValueError(f"unsupported target: {system}")


def read_ar(path):
    """Read the standard ar members used by Debian binary packages."""
    data = path.read_bytes()
    if data[:8] != b"!<arch>\n":
        raise ValueError("Debian package is not an ar archive")
    members, offset = {}, 8
    while offset < len(data):
        header = data[offset:offset+60]
        if header[58:60] != b"`\n":
            raise ValueError("invalid ar member")
        name = header[:16].decode().strip().rstrip("/")
        size = int(header[48:58])
        members[name] = data[offset+60:offset+60+size]
        offset += 60 + size + size % 2
    return members


def verify(dist):
    """Check every promised target, its packaged executable, and all checksum entries."""
    for system in ["linux", "darwin", "windows"]:
        for arch in ["amd64", "arm64"]:
            extension = "zip" if system == "windows" else "tar.gz"
            paths = list(dist.glob(f"lagerator_*_{system}_{arch}.{extension}"))
            if len(paths) != 1:
                raise ValueError(f"expected one archive for {system}/{arch}: {paths}")
            if extension == "zip":
                with zipfile.ZipFile(paths[0]) as archive:
                    names = archive.namelist()
                    content = archive.read("lgrt.exe")
            else:
                with tarfile.open(paths[0]) as archive:
                    names = archive.getnames()
                    content = archive.extractfile("lgrt").read()
            if binary_architecture(content, system) != arch:
                raise ValueError(f"incorrect executable architecture: {paths[0]}")
            if "LICENSE" not in names or not any("licenses/" in name for name in names):
                raise ValueError(f"missing licenses: {paths[0]}")
            if "lgrtdata.json" in names or "bin/ginv" in names:
                raise ValueError("inventory or development binary included in release")
    for arch in ["amd64", "arm64"]:
        packages = list(dist.glob(f"lagerator_*_{arch}.deb"))
        if len(packages) != 1:
            raise ValueError(f"expected one Debian package for {arch}")
        members = read_ar(packages[0])
        data_name = next(name for name in members if name.startswith("data.tar"))
        with tarfile.open(fileobj=io.BytesIO(members[data_name])) as archive:
            names = {member.name.removeprefix("./"): member for member in archive}
            content = archive.extractfile(names["usr/bin/lgrt"]).read()
            if binary_architecture(content, "linux") != arch:
                raise ValueError("Debian executable architecture does not match package")
            if not any(name.startswith("usr/share/doc/lagerator/licenses/") for name in names):
                raise ValueError("Debian package is missing dependency licenses")
        control_name = next(name for name in members if name.startswith("control.tar"))
        with tarfile.open(fileobj=io.BytesIO(members[control_name])) as archive:
            member = next(item for item in archive if item.name.removeprefix("./") == "control")
            control = archive.extractfile(member).read().decode()
            if f"Architecture: {arch}\n" not in control:
                raise ValueError("incorrect Debian control architecture")
    support = Path(__file__).resolve().parent.parent / "packaging/generated"
    source, = support.glob("*_source.tar.gz")
    with tarfile.open(source) as archive:
        names = archive.getnames()
        root = names[0].split("/")[0]
        for required in ["go.mod", "vendor/modules.txt", "main.go"]:
            if f"{root}/{required}" not in names:
                raise ValueError(f"source archive missing {required}")
        if f"{root}/lgrtdata.json" in names or any(name.startswith(f"{root}/bin/") for name in names):
            raise ValueError("source archive includes private inventory or binaries")
    assets = {path.name: path for tree in [dist, support] for path in tree.rglob("*") if path.is_file()}
    lines = (dist / "checksums.txt").read_text().splitlines()
    for line in lines:
        digest, name = line.split(maxsplit=1)
        path = assets[Path(name.strip().lstrip("*")).name]
        if hashlib.sha256(path.read_bytes()).hexdigest() != digest:
            raise ValueError(f"checksum mismatch: {path}")
    for path in [source, *support.glob("*.ebuild"), support / "Manifest"]:
        if not any(path.name == Path(line.split(maxsplit=1)[1].strip()).name for line in lines):
            raise ValueError(f"support artifact missing from checksums: {path}")
    print("Verified six platform archives, two Debian packages, vendored source, licenses, and checksums.")


def main():
    """Parse the release distribution path and validate its contents."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("dist", type=Path)
    verify(parser.parse_args().dist)


if __name__ == "__main__":
    main()
