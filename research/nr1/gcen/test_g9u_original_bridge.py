"""CPU-only negative controls for original G9U/G9N source bridge."""
import unittest
from unittest.mock import patch
import numpy as np
import g9n_one_original_tensor as g9n
import g9u_verified_loader as g9u

NAME="model.layers.0.mlp.experts.0.down_proj.weight"

class Tensor:
    def __init__(self, data): self.data=data
    def to(self,*,dtype): return self
    def numpy(self): return self.data

class BridgeTests(unittest.TestCase):
    def setUp(self):
        self.w=np.ones((2048,768),dtype=np.float32)
        self.receipt={
            "tensor_name":NAME,
            "tensor_original_dtype":"BF16",
            "original_shard_authenticated":True,
            "shard_sha256":"a"*64,
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

    def test_rejects_modified_tensor(self):
        original=self.w.copy()
        original[0,0]=2.0
        with patch.object(g9u,"load_verified_original_expert",
                          return_value=(self.w,self.receipt)):
            with self.assertRaisesRegex(RuntimeError,"tensor bytes"):
                g9n.verify_g9u_original_custody("shard","manifest",NAME,Tensor(original))

if __name__=="__main__":unittest.main()
