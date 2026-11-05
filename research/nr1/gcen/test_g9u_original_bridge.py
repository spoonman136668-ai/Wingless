"""CPU-only negative controls for original G9U/G9N source bridge."""
import unittest
import sys
import types
from unittest.mock import patch
import numpy as np
import g9n_one_original_tensor as g9n
import g9u_verified_loader as g9u
from runtime_feasibility import MODEL, REVISION
from g31_original_shard_stage import SHARD

NAME="model.layers.0.mlp.experts.0.down_proj.weight"

class Tensor:
    def __init__(self, data): self.data=data
    def to(self,*,dtype): return self
    def numpy(self): return self.data

class BridgeTests(unittest.TestCase):
    def setUp(self):
        # Routine native CI has NumPy but intentionally does not install Torch.
        # Inject only the dtype marker used by the bridge; real BF16 custody
        # remains exclusively exercised in the isolated original-weight job.
        mock_torch=types.SimpleNamespace(float32="mock-fp32-dtype")
        patcher=patch.dict(sys.modules,{"torch":mock_torch})
        patcher.start()
        self.addCleanup(patcher.stop)
        self.w=np.ones((2048,768),dtype=np.float32)
        self.receipt={
            "tensor_name":NAME,
            "tensor_original_dtype":"BF16",
            "original_shard_authenticated":True,
            "shard_sha256":g9u.PINNED_SHARD_SHA256,
            "shard":SHARD,
            "model":MODEL,
            "revision":REVISION,
            "tensor_shape":[2048,768],
            "source_authentication_method":"IN_PROCESS_SHA256_PINNED_INDEX_BF16",
        }

    def run_bridge(self):
        with patch.object(g9u,"load_verified_original_expert",
                          return_value=(self.w,self.receipt)):
            return g9n.verify_g9u_original_custody("shard","manifest",NAME,Tensor(self.w.copy()))

    def test_positive_mocked_parity(self):
        report=self.run_bridge()
        self.assertEqual(report["result"],"G9U_ORIGINAL_BF16_CPU_LOADER_PARITY_VERIFIED")
        self.assertFalse(report["cuda_executed"])

    def test_rejects_wrong_identity(self):
        self.receipt["tensor_name"]="different"
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_missing_authentication(self):
        self.receipt["original_shard_authenticated"]=False
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_wrong_dtype_receipt(self):
        self.receipt["tensor_original_dtype"]="FP32"
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_wrong_pinned_shard_digest(self):
        self.receipt["shard_sha256"]="b"*64
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_wrong_model_revision(self):
        self.receipt["revision"]="f"*40
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_wrong_shard_name(self):
        self.receipt["shard"]="model-00002-of-00016.safetensors"
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_wrong_tensor_shape_receipt(self):
        self.receipt["tensor_shape"]=[768,2048]
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_unverifiable_method(self):
        self.receipt["source_authentication_method"]="SELF_DECLARED"
        with self.assertRaisesRegex(RuntimeError,"identity"): self.run_bridge()

    def test_rejects_non_mapping_custody(self):
        with patch.object(g9u,"load_verified_original_expert",
                          return_value=(self.w,True)):
            with self.assertRaisesRegex(RuntimeError,"identity"):
                g9n.verify_g9u_original_custody("shard","manifest",NAME,Tensor(self.w.copy()))

    def test_rejects_modified_tensor(self):
        original=self.w.copy()
        original[0,0]=2.0
        with patch.object(g9u,"load_verified_original_expert",
                          return_value=(self.w,self.receipt)):
            with self.assertRaisesRegex(RuntimeError,"tensor bytes"):
                g9n.verify_g9u_original_custody("shard","manifest",NAME,Tensor(original))

if __name__=="__main__":unittest.main()
