"""Synthetic identity-drift and manifest tests for session-scoped checkpoint custody."""
import hashlib
import json
import tempfile
import unittest
from pathlib import Path

from custody import REQUIRED
from verified_session import VerifiedCheckpoint


class VerifiedSessionTests(unittest.TestCase):
    def fixture(self, root):
        for name in REQUIRED[:16]:
            (root / name).write_bytes(b"synthetic shard")
        (root / "config.json").write_text("{}", encoding="utf-8")
        (root / "model.safetensors.index.json").write_text(
            json.dumps({"weight_map": {"probe": REQUIRED[0]}}), encoding="utf-8"
        )
        return {name: hashlib.sha256((root / name).read_bytes()).hexdigest()
                for name in REQUIRED}

    def test_complete_manifest_accepts_unchanged_files(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            session = VerifiedCheckpoint(root, self.fixture(root))
            self.assertEqual(session._checked(REQUIRED[0]), root / REQUIRED[0])

    def test_missing_manifest_entry_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            manifest = self.fixture(root)
            del manifest[REQUIRED[0]]
            with self.assertRaisesRegex(ValueError, "incomplete"):
                VerifiedCheckpoint(root, manifest)

    def test_wrong_digest_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            manifest = self.fixture(root)
            manifest[REQUIRED[0]] = "0" * 64
            with self.assertRaisesRegex(ValueError, "SHA256 mismatch"):
                VerifiedCheckpoint(root, manifest)

    def test_mutation_after_verification_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            session = VerifiedCheckpoint(root, self.fixture(root))
            (root / REQUIRED[0]).write_bytes(b"changed content and size")
            with self.assertRaisesRegex(ValueError, "identity drift"):
                session._checked(REQUIRED[0])

    def test_unknown_file_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            session = VerifiedCheckpoint(root, self.fixture(root))
            with self.assertRaisesRegex(ValueError, "unknown"):
                session._checked("not-a-checkpoint-file")

    def test_symlink_rejected_when_supported(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            manifest = self.fixture(root)
            path = root / REQUIRED[0]
            target = root / "target.bin"
            target.write_bytes(path.read_bytes())
            path.unlink()
            try:
                path.symlink_to(target)
            except (OSError, NotImplementedError):
                self.skipTest("symlink creation unavailable")
            with self.assertRaisesRegex(ValueError, "symlinked"):
                VerifiedCheckpoint(root, manifest)


if __name__ == "__main__":
    unittest.main()
