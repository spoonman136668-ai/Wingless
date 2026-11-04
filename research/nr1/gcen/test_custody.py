"""Focused custody tests using small synthetic safetensors; no Qwen weights."""
import hashlib
import json
import tempfile
import unittest
from pathlib import Path
from custody import REQUIRED, verify_manifest, extract_tensor

class CustodyTests(unittest.TestCase):
    def test_incomplete_manifest_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with self.assertRaises(ValueError): verify_manifest(d,{})
    def test_invalid_hash_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with self.assertRaises(ValueError): verify_manifest(d,{name:"0"*63 for name in REQUIRED})
    def test_missing_files_rejected(self):
        with tempfile.TemporaryDirectory() as d:
            with self.assertRaises(ValueError): verify_manifest(d,{name:"0"*64 for name in REQUIRED})
    def test_extract_tensor(self):
        import numpy as np
        from safetensors.numpy import save_file
        with tempfile.TemporaryDirectory() as d:
            root=Path(d)
            for name in REQUIRED[:16]: save_file({"probe":np.arange(6,dtype=np.float32).reshape(2,3)},str(root/name))
            (root/"config.json").write_text("{}",encoding="utf-8")
            (root/"model.safetensors.index.json").write_text(json.dumps({"weight_map":{"probe":REQUIRED[0]}}),encoding="utf-8")
            manifest={name:hashlib.sha256((root/name).read_bytes()).hexdigest() for name in REQUIRED}
            self.assertTrue(verify_manifest(root,manifest))
            np.testing.assert_array_equal(extract_tensor(root,manifest,"probe"),np.arange(6,dtype=np.float32).reshape(2,3))
            (root/REQUIRED[0]).write_bytes(b"corrupt")
            with self.assertRaises(ValueError): verify_manifest(root,manifest)

if __name__=="__main__": unittest.main()
