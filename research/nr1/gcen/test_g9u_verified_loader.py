"""Focused source-custody negative controls; never downloads model weights."""
import json
from pathlib import Path
import tempfile
import unittest

from g9u_verified_loader import load_verified_original_expert, PINNED_SHARD_SHA256
from g31_original_shard_stage import SHARD
from runtime_feasibility import MODEL, REVISION


class FailClosedCustodyTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        root = Path(self.temp.name)
        self.shard = root / SHARD
        self.shard.write_bytes(b"forged source shard")
        self.manifest = root / "manifest.json"
        files = {f"model-{i:05d}-of-00016.safetensors": "a"*64 for i in range(1, 17)}
        files.update({"config.json":"a"*64,"model.safetensors.index.json":"a"*64})
        files[SHARD] = PINNED_SHARD_SHA256
        self.data = {"model":MODEL,"revision":REVISION,"files":files}
        self.write_manifest()

    def write_manifest(self):
        self.manifest.write_text(json.dumps(self.data),encoding="utf-8")

    def test_rejects_forged_original_bytes(self):
        with self.assertRaisesRegex(ValueError, "SHA256 mismatch"):
            load_verified_original_expert(self.shard,self.manifest)

    def test_rejects_unpinned_manifest_digest(self):
        self.data["files"][SHARD] = "b"*64
        self.write_manifest()
        with self.assertRaisesRegex(ValueError, "authenticated pin"):
            load_verified_original_expert(self.shard,self.manifest)

    def test_rejects_wrong_model_revision(self):
        self.data["revision"]="f"*40
        self.write_manifest()
        with self.assertRaisesRegex(ValueError,"original-model identity"):
            load_verified_original_expert(self.shard,self.manifest)

    def test_rejects_missing_shard(self):
        self.shard.unlink()
        with self.assertRaisesRegex(ValueError,"path identity"):
            load_verified_original_expert(self.shard,self.manifest)

    def test_rejects_wrong_filename(self):
        other=self.shard.with_name("other.safetensors")
        self.shard.rename(other)
        with self.assertRaisesRegex(ValueError,"path identity"):
            load_verified_original_expert(other,self.manifest)

    def test_rejects_symlink(self):
        real=self.shard.with_name("real")
        self.shard.rename(real)
        try: self.shard.symlink_to(real)
        except (OSError,NotImplementedError): self.skipTest("symlink unavailable")
        with self.assertRaisesRegex(ValueError,"path identity"):
            load_verified_original_expert(self.shard,self.manifest)

if __name__ == "__main__":
    unittest.main()
