"""Offline negative controls for one-shot original-shard staging. No network."""
import hashlib
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

from g31_original_shard_stage import SHARD
from runtime_feasibility import MODEL, REVISION
import g9u_stage_pinned_shard as stage


class Response:
    status = 200
    def __init__(self, payload):
        self.buf = io.BytesIO(payload)
        self.headers = {"Content-Length": str(len(payload)),
                        "Content-Range": None, "Content-Encoding": None}
    def geturl(self):
        return "https://cdn.example/pinned"
    def read(self, n):
        return self.buf.read(n)
    def __enter__(self):
        return self
    def __exit__(self, *args):
        self.buf.close()


class StageTests(unittest.TestCase):
    def setUp(self):
        d = tempfile.TemporaryDirectory()
        self.addCleanup(d.cleanup)
        self.root = Path(d.name)
        self.tools = self.root / "tools"
        self.tools.mkdir()
        self.evidence = self.root / "evidence"
        self.evidence.mkdir()
        self.output = self.evidence / "stage.json"
        self.manifest = self.root / "manifest.json"
        self.payload = b"tiny-test-source-only"
        self.digest = hashlib.sha256(self.payload).hexdigest()
        files = {f"model-{n:05d}-of-00016.safetensors": "a"*64
                 for n in range(1, 17)}
        files[SHARD] = self.digest
        files["config.json"] = "b"*64
        files["model.safetensors.index.json"] = "c"*64
        self.manifest.write_text(json.dumps({"model": MODEL, "revision": REVISION,
                                             "files": files}), encoding="utf-8")
        for name, value in (("PINNED_SHARD_SHA256", self.digest),
                            ("EXPECTED_BYTES", len(self.payload))):
            p = patch.object(stage, name, value)
            p.start()
            self.addCleanup(p.stop)
    def run_stage(self):
        return stage.stage(toolcache_root=self.tools, manifest_path=self.manifest,
                           output=self.output, source_sha="f"*40)

    def test_authenticated_stage_then_reuse_without_second_request(self):
        with patch.object(stage, "urlopen", return_value=Response(self.payload)) as opener:
            result = self.run_stage()
        self.assertEqual(opener.call_count, 1)
        self.assertEqual(result["status"], "DOWNLOADED_VERIFIED_ORIGINAL_SHARD")
        self.assertFalse(result["gpu_kernel_executed"])
        target = self.tools / stage.DIRNAME / SHARD
        self.assertEqual(target.read_bytes(), self.payload)
        self.assertEqual(list(target.parent.iterdir()), [target])
        self.output.unlink()
        with patch.object(stage, "urlopen", side_effect=AssertionError("network denied")):
            self.assertEqual(self.run_stage()["status"], "REUSED_VERIFIED_ORIGINAL_SHARD")

    def test_invalid_stream_hash_fails_closed_and_cleans(self):
        with patch.object(stage, "urlopen",
                          return_value=Response(b"x"*len(self.payload))):
            with self.assertRaisesRegex(ValueError, "SHA256"):
                self.run_stage()
        self.assertFalse(self.output.exists())
        self.assertEqual(list((self.tools / stage.DIRNAME).iterdir()), [])

    def test_existing_corruption_never_redownloads_or_overwrites(self):
        directory = self.tools / stage.DIRNAME
        directory.mkdir()
        (directory / SHARD).write_bytes(b"x"*len(self.payload))
        with patch.object(stage, "urlopen", side_effect=AssertionError("network denied")):
            with self.assertRaisesRegex(ValueError, "SHA256"):
                self.run_stage()
        self.assertFalse(self.output.exists())

    def test_wrong_manifest_never_stages(self):
        manifest = json.loads(self.manifest.read_text())
        manifest["revision"] = "0"*40
        self.manifest.write_text(json.dumps(manifest))
        with patch.object(stage, "urlopen", side_effect=AssertionError("network denied")):
            with self.assertRaisesRegex(ValueError, "identity"):
                self.run_stage()
        self.assertFalse(self.output.exists())

    def test_disk_floor_denies_before_request(self):
        class Disk:
            free = 0
        with patch.object(stage.shutil, "disk_usage", return_value=Disk()), patch.object(
                stage, "urlopen", side_effect=AssertionError("network denied")):
            with self.assertRaisesRegex(RuntimeError, "storage"):
                self.run_stage()

    def test_existing_receipt_and_bad_sha_rejected(self):
        self.output.write_text("keep")
        with self.assertRaisesRegex(ValueError, "overwrite"):
            self.run_stage()
        self.output.unlink()
        with self.assertRaisesRegex(ValueError, "source SHA"):
            stage.stage(toolcache_root=self.tools, manifest_path=self.manifest,
                        output=self.output, source_sha="invalid")


if __name__ == "__main__":
    unittest.main()
