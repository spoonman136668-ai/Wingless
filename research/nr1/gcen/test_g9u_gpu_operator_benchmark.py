"""G9U source/format negative controls; routine CI never invokes CUDA."""
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest.mock import patch
import numpy as np
import g9u_gpu_operator_benchmark as g9u
from g9u_gpu_operator_benchmark import (
    validate_source,validate_packed,ROWS,COLUMNS,BLOCK
)
SHA="e"*40

class G9UAdmissionTests(unittest.TestCase):
    def setUp(self):
        self.w=np.ones((ROWS,COLUMNS),dtype=np.float32)
        self.data=np.zeros((ROWS,COLUMNS//4),dtype=np.uint8)
        self.scales=np.ones((ROWS,COLUMNS//BLOCK),dtype=np.float16)
    def test_source_shape_and_sha(self):
        self.assertEqual(validate_source(self.w,SHA).shape,(ROWS,COLUMNS))
    def test_source_bad_sha(self):
        with self.assertRaisesRegex(ValueError,"SHA"):validate_source(self.w,"not-a-sha")
    def test_source_nonfinite(self):
        self.w[0,0]=np.nan
        with self.assertRaisesRegex(ValueError,"source"):validate_source(self.w,SHA)
    def test_source_wrong_shape(self):
        with self.assertRaisesRegex(ValueError,"source"):validate_source(self.w[:1],SHA)
    def test_packed_format(self):
        self.assertEqual(validate_packed(self.data,self.scales)[0].shape,self.data.shape)
    def test_packed_wrong_dtype(self):
        with self.assertRaisesRegex(ValueError,"format"):validate_packed(self.data.astype(np.int16),self.scales)
    def test_packed_wrong_inventory(self):
        with self.assertRaisesRegex(ValueError,"format"):validate_packed(self.data[:100],self.scales)
    def test_packed_nonfinite_scale(self):
        self.scales[0,0]=np.inf
        with self.assertRaisesRegex(ValueError,"format"):validate_packed(self.data,self.scales)

class G9UCheckoutIdentityTests(unittest.TestCase):
    def check(self, *outputs, sha=SHA):
        with patch.object(g9u.subprocess,"run",
                          side_effect=[SimpleNamespace(stdout=value) for value in outputs]) as runner:
            result=g9u.verify_executing_checkout(sha)
        return result,runner.call_args_list

    def test_accepts_clean_exact_checkout(self):
        root=str(Path(g9u.__file__).resolve().parents[3])
        observed,calls=self.check(root+"\\n",SHA+"\\n","")
        self.assertEqual(observed,SHA)
        self.assertEqual(len(calls),3)

    def test_rejects_unrelated_repository_root(self):
        with self.assertRaisesRegex(RuntimeError,"root mismatch"):
            self.check("/tmp/unrelated",SHA,"")

    def test_rejects_different_executing_sha(self):
        root=str(Path(g9u.__file__).resolve().parents[3])
        with self.assertRaisesRegex(RuntimeError,"SHA mismatch"):
            self.check(root,"a"*40,"")

    def test_rejects_dirty_tracked_checkout(self):
        root=str(Path(g9u.__file__).resolve().parents[3])
        with self.assertRaisesRegex(RuntimeError,"dirty tracked"):
            self.check(root,SHA," M research/nr1/gcen/g9u_gpu_operator_benchmark.py")

    def test_rejects_invalid_source_sha_before_git(self):
        with patch.object(g9u.subprocess,"run") as runner:
            with self.assertRaisesRegex(ValueError,"SHA"):
                g9u.verify_executing_checkout("invalid")
            runner.assert_not_called()

    def test_rejects_unavailable_git_checkout(self):
        with patch.object(g9u.subprocess,"run",side_effect=FileNotFoundError):
            with self.assertRaisesRegex(RuntimeError,"unavailable"):
                g9u.verify_executing_checkout(SHA)

    def test_rejects_checkout_timeout(self):
        with patch.object(g9u.subprocess,"run",
                          side_effect=g9u.subprocess.TimeoutExpired("git",10)):
            with self.assertRaisesRegex(RuntimeError,"unavailable"):
                g9u.verify_executing_checkout(SHA)

    def test_rejects_checkout_mismatch_before_weight_loading(self):
        with patch.object(g9u,"verify_executing_checkout",
                          side_effect=RuntimeError("SHA mismatch")):
            with patch("g9u_verified_loader.load_verified_original_expert") as loader:
                with self.assertRaisesRegex(RuntimeError,"SHA mismatch"):
                    g9u.run_authenticated_expert(
                        shard_path="not-loaded",manifest_path="not-loaded",source_sha=SHA)
                loader.assert_not_called()

if __name__=="__main__":unittest.main()
