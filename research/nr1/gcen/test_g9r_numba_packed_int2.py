"""G9R frozen-input negative controls; no Numba required by routine CI."""
import unittest
import numpy as np
from g9r_numba_packed_int2 import validate,relative_squared_error,BATCH
class G9RTests(unittest.TestCase):
    def setUp(self):
        self.x=np.zeros((BATCH,256),np.float32)
        self.d=np.zeros((2,64),np.uint8)
        self.s=np.ones((2,1),np.float16)
    def test_valid(self):
        self.assertEqual(validate(self.x,self.d,self.s).shape,(8,256))
    def test_wrong_batch(self):
        with self.assertRaises(ValueError):validate(self.x[:4],self.d,self.s)
    def test_nonfinite(self):
        self.x[0,0]=np.nan
        with self.assertRaises(ValueError):validate(self.x,self.d,self.s)
    def test_shape_mismatch(self):
        with self.assertRaises(ValueError):validate(self.x,self.d[:1],self.s)
    def test_dense_reference_error_exact(self):
        ref=np.array([[2.,-1.]],dtype=np.float32)
        self.assertEqual(relative_squared_error(ref,ref),0.)
        other=np.array([[1.,-1.]],dtype=np.float32)
        self.assertAlmostEqual(relative_squared_error(other,ref),0.2)

    def test_dense_reference_invalid_shapes(self):
        a=np.ones((1,2),dtype=np.float32)
        with self.assertRaisesRegex(ValueError,"matched"):
            relative_squared_error(a,np.ones((2,2),dtype=np.float32))

    def test_dense_reference_zero_energy_rejected(self):
        z=np.zeros((1,2),dtype=np.float32)
        with self.assertRaisesRegex(ValueError,"energy"):
            relative_squared_error(z,z)

    def test_nonfinite_scales(self):
        self.s[0,0]=np.inf
        with self.assertRaises(ValueError):validate(self.x,self.d,self.s)
if __name__=="__main__":unittest.main()
