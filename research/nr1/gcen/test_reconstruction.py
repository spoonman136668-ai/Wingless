"""Development-only focused synthetic tests."""
import unittest
import numpy as np
from reconstruction import original, svd, shared_basis, storage_bytes

class ReconstructionTests(unittest.TestCase):
    def test_original(self):
        w=np.arange(12,dtype=np.float32).reshape(3,4)
        np.testing.assert_array_equal(original(w),w)
    def test_svd_full_rank(self):
        w=np.arange(12,dtype=np.float32).reshape(3,4)
        np.testing.assert_allclose(svd(w,3),w,atol=1e-4)
    def test_shared_full_rank(self):
        w=np.arange(24,dtype=np.float32).reshape(2,3,4)
        np.testing.assert_allclose(shared_basis(w,4),w,atol=1e-4)
    def test_storage(self):
        self.assertEqual(storage_bytes((4,6,8),2,'original'),384)
        self.assertEqual(storage_bytes((4,6,8),2,'svd'),240)
        self.assertEqual(storage_bytes((4,6,8),2,'shared'),128)
    def test_invalid_rank(self):
        with self.assertRaises(ValueError): svd(np.ones((3,4)),0)

if __name__ == '__main__': unittest.main()
