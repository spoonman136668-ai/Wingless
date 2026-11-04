"""Offline collector tests: synthetic metadata; no network or model weights."""
import hashlib
import json
import unittest

from custody import REQUIRED
from source_manifest import REVISION
from collect_upstream_manifest import API, BASE, MODEL, build_manifest


class UpstreamCollectorTests(unittest.TestCase):
    def setUp(self):
        self.revision = REVISION
        self.shards = REQUIRED[:16]
        self.info = {
            "sha": REVISION,
            "siblings": [
                {"rfilename": name, "lfs": {"sha256": "a" * 64}}
                for name in self.shards
            ] + [{"rfilename": name} for name in REQUIRED[16:]],
        }

    def fetch(self, url, limit):
        if url == API:
            return json.dumps(self.info).encode()
        if url.endswith("/config.json"):
            return b"{}"
        if url.endswith("/model.safetensors.index.json"):
            return b'{"weight_map":{}}'
        raise AssertionError("unexpected URL: " + url)

    def test_complete_pinned_inventory(self):
        doc = build_manifest(self.fetch)
        self.assertEqual(len(doc["files"]), 18)
        self.assertEqual(doc["revision"], REVISION)
        self.assertEqual(doc["files"][self.shards[0]], "a" * 64)
        self.assertEqual(doc["files"]["config.json"], hashlib.sha256(b"{}").hexdigest())

    def test_revision_mismatch_fails(self):
        self.info["sha"] = "0" * 40
        with self.assertRaisesRegex(ValueError, "revision mismatch"):
            build_manifest(self.fetch)

    def test_missing_shard_fails(self):
        self.info["siblings"].pop(0)
        with self.assertRaisesRegex(ValueError, "incomplete"):
            build_manifest(self.fetch)

    def test_missing_lfs_digest_fails(self):
        self.info["siblings"][0].pop("lfs")
        with self.assertRaisesRegex(ValueError, "LFS"):
            build_manifest(self.fetch)

    def test_duplicate_sibling_fails(self):
        self.info["siblings"].append(self.info["siblings"][0].copy())
        with self.assertRaisesRegex(ValueError, "duplicate"):
            build_manifest(self.fetch)


if __name__ == "__main__":
    unittest.main()
