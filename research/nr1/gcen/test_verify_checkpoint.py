"""Synthetic-only tests for the offline checkpoint verifier."""
import hashlib
import tempfile
import unittest
from pathlib import Path

from custody import REQUIRED
from source_manifest import MODEL, REVISION
from verify_checkpoint import verify_checkpoint


class CheckpointVerificationTests(unittest.TestCase):
    def setUp(self):
        self.workspace = tempfile.TemporaryDirectory()
        self.addCleanup(self.workspace.cleanup)
        self.root = Path(self.workspace.name)
        for name in REQUIRED:
            (self.root / name).write_bytes(name.encode("utf-8"))
        self.manifest = {
            "model": MODEL,
            "revision": REVISION,
            "source": {
                "uri": "https://huggingface.co/Qwen/Qwen3-30B-A3B/tree/" + REVISION,
                "revision": REVISION,
            },
            "files": {
                name: hashlib.sha256((self.root / name).read_bytes()).hexdigest()
                for name in REQUIRED
            },
        }

    def test_receipt_does_not_claim_authoritative_custody(self):
        result = verify_checkpoint(self.root, self.manifest)
        self.assertEqual(result["status"], "LOCAL_BYTES_VERIFIED_SOURCE_AUTHORITY_UNVERIFIED")
        self.assertEqual(len(result["files"]), 18)

    def test_mutated_file_rejected(self):
        (self.root / REQUIRED[0]).write_bytes(b"changed")
        with self.assertRaises(ValueError):
            verify_checkpoint(self.root, self.manifest)

    def test_missing_file_rejected(self):
        (self.root / REQUIRED[0]).unlink()
        with self.assertRaises(ValueError):
            verify_checkpoint(self.root, self.manifest)

    def test_revision_drift_rejected(self):
        self.manifest["revision"] = "0" * 40
        with self.assertRaises(ValueError):
            verify_checkpoint(self.root, self.manifest)


if __name__ == "__main__":
    unittest.main()
