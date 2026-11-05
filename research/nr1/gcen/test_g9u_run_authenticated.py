"""Offline admission/receipt tests; no original weights, download, or CUDA."""
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

import g9u_run_authenticated as run

SHA="f"*40

class RunnerTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root=Path(self.temp.name)
        self.shard=self.root/"shard"
        self.manifest=self.root/"manifest"
        self.output=self.root/"result.json"
        self.shard.write_bytes(b"test")
        self.manifest.write_text("{}")

    def test_rejects_invalid_sha_before_git(self):
        with self.assertRaisesRegex(ValueError,"source SHA"):
            run.check_checkout(self.root,"invalid")

    def test_rejects_wrong_head(self):
        with patch.object(run.subprocess,"run") as proc:
            proc.return_value.stdout=str(self.root)+"\n" if not hasattr(self,"step") else ""
            proc.side_effect=[
                type("R",(),{"stdout":str(self.root)+"\n"})(),
                type("R",(),{"stdout":"e"*40+"\n"})(),
            ]
            with self.assertRaisesRegex(ValueError,"source SHA mismatch"):
                run.check_checkout(self.root,SHA)

    def test_rejects_dirty_checkout(self):
        with patch.object(run.subprocess,"run") as proc:
            proc.side_effect=[
                type("R",(),{"stdout":str(self.root)+"\n"})(),
                type("R",(),{"stdout":SHA+"\n"})(),
                type("R",(),{"stdout":"?? stray_file\n"})(),
            ]
            with self.assertRaisesRegex(ValueError,"dirty checkout"):
                run.check_checkout(self.root,SHA)

    def test_rejects_existing_output(self):
        self.output.write_text("existing")
        with patch.object(run,"check_checkout",return_value=self.root):
            with self.assertRaises(FileExistsError):
                self.invoke()

    def invoke(self):
        return run.execute(repository_root=self.root,source_sha=SHA,
                           shard_path=self.shard,manifest_path=self.manifest,
                           output=self.output)

    def test_rejects_missing_authentication(self):
        with patch.object(run,"check_checkout",return_value=self.root),patch.object(
            run,"run_authenticated_expert",return_value={"exact_source_sha":SHA}):
            with self.assertRaisesRegex(RuntimeError,"incomplete"):
                self.invoke()
        self.assertFalse(self.output.exists())

    def test_positive_stubbed_receipt(self):
        fake={"exact_source_sha":SHA,"original_shard_authenticated_in_process":True,
              "measurements":[]}
        with patch.object(run,"check_checkout",return_value=self.root),patch.object(
            run,"run_authenticated_expert",return_value=fake) as benchmark:
            result=self.invoke()
        self.assertEqual(result["benchmark"],fake)
        self.assertEqual(json.loads(self.output.read_text())["source_sha"],SHA)
        self.assertFalse(result["full_model_forward"])
        benchmark.assert_called_once()

    def test_rejects_receipt_sha_drift(self):
        with patch.object(run,"check_checkout",return_value=self.root),patch.object(
            run,"run_authenticated_expert",return_value={
                "exact_source_sha":"e"*40,
                "original_shard_authenticated_in_process":True}):
            with self.assertRaisesRegex(RuntimeError,"incomplete"):
                self.invoke()

if __name__=="__main__":
    unittest.main()
