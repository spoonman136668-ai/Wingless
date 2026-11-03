"""Negative controls for pinned-source custody preflight; synthetic hashes only."""
import copy
import unittest
from custody import REQUIRED
from source_manifest import MODEL, REVISION, admit_source_manifest


class SourceManifestTests(unittest.TestCase):
    def setUp(self):
        self.doc = {
            "model": MODEL, "revision": REVISION,
            "source": {"uri": f"https://huggingface.co/{MODEL}/tree/{REVISION}", "revision": REVISION},
            "files": {name: "a" * 64 for name in REQUIRED},
        }

    def test_complete_shape_admitted_only_as_preflight(self):
        self.assertEqual(admit_source_manifest(self.doc), self.doc["files"])

    def test_wrong_revision_rejected(self):
        self.doc["revision"] = "0" * 40
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)

    def test_source_revision_mismatch_rejected(self):
        self.doc["source"]["revision"] = "0" * 40
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)

    def test_unrecognized_source_rejected(self):
        self.doc["source"]["uri"] = "https://example.invalid/weights"
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)

    def test_missing_file_rejected(self):
        del self.doc["files"][REQUIRED[0]]
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)

    def test_extra_file_rejected(self):
        self.doc["files"]["extra"] = "a" * 64
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)

    def test_invalid_hash_rejected(self):
        self.doc["files"][REQUIRED[0]] = "A" * 64
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)

    def test_extra_metadata_rejected(self):
        self.doc["unreviewed"] = True
        with self.assertRaises(ValueError): admit_source_manifest(self.doc)


if __name__ == "__main__":
    unittest.main()
