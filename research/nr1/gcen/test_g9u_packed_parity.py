"""CPU-only packed-CUDA parity receipt negative controls; no GPU or weights."""
import unittest
import numpy as np
from g9u_gpu_operator_benchmark import (
    verify_packed_operator_parity, BATCH, ROWS
)

class PackedCudaParityTests(unittest.TestCase):
    def setUp(self):
        self.reference=np.ones((BATCH,ROWS),dtype=np.float32)

    def test_exact_parity(self):
        report=verify_packed_operator_parity(self.reference,self.reference.copy())
        self.assertEqual(report["relative_squared_error"],0.0)

    def test_small_rounding_difference_allowed(self):
        observed=self.reference.copy()
        observed[0,0]+=1e-5
        self.assertLess(verify_packed_operator_parity(
            self.reference,observed)["maximum_absolute_difference"],2e-4)

    def test_rejects_corrupt_decode(self):
        observed=self.reference.copy()
        observed[0,0]=5.0
        with self.assertRaisesRegex(RuntimeError,"disagrees"):
            verify_packed_operator_parity(self.reference,observed)

    def test_rejects_global_scale_drift(self):
        with self.assertRaisesRegex(RuntimeError,"disagrees"):
            verify_packed_operator_parity(self.reference,self.reference*1.001)

    def test_rejects_wrong_shape(self):
        with self.assertRaisesRegex(ValueError,"inputs"):
            verify_packed_operator_parity(self.reference,self.reference[:1])

    def test_rejects_nonfinite(self):
        observed=self.reference.copy()
        observed[0,0]=np.nan
        with self.assertRaisesRegex(ValueError,"inputs"):
            verify_packed_operator_parity(self.reference,observed)

if __name__=="__main__":
    unittest.main()
