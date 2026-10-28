"""Check release source filtering, reproducibility, and Gentoo distfile metadata."""

import hashlib
import importlib.util
from pathlib import Path
import tarfile
import tempfile
import unittest

SPEC = importlib.util.spec_from_file_location("prepare_release", Path(__file__).with_name("prepare-release.py"))
RELEASE = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(RELEASE)


class ReleaseTests(unittest.TestCase):
    """Exercise packaging invariants that protect inventory data and source integrity."""

    def test_source_excludes_inventory_and_binaries(self):
        """Ensure private inventory snapshots and checked-in executables are not shipped."""
        paths = RELEASE.source_files()
        self.assertNotIn(Path("lgrtdata.json"), paths)
        self.assertFalse(any(path.parts[0] == "bin" for path in paths))
        self.assertIn(Path("go.mod"), paths)
        self.assertIn(Path("main.go"), paths)

    def test_reproducible_archive_and_manifest(self):
        """Ensure repeated archives match and their Gentoo checksums cover actual content."""
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            source = root / "lagerator-1.0.13"
            source.mkdir()
            (source / "main.go").write_text("package main\n")
            first, second = root / "first.tar.gz", root / "second.tar.gz"
            RELEASE.archive_tree(source, first, 1700000000)
            RELEASE.archive_tree(source, second, 1700000000)
            self.assertEqual(first.read_bytes(), second.read_bytes())
            with tarfile.open(first) as archive:
                self.assertIn("lagerator-1.0.13/main.go", archive.getnames())
            entry = RELEASE.manifest_entry(first)
            self.assertIn(f"DIST {first.name} {first.stat().st_size} BLAKE2B ", entry)
            self.assertIn(hashlib.sha512(first.read_bytes()).hexdigest(), entry)


if __name__ == "__main__":
    unittest.main()
