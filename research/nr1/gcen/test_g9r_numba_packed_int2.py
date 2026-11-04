"""G9R frozen-input negative controls; no Numba required by routine CI."""
import unittest
import numpy as np
from g9r_numba_packed_int2 import validate,BATCH
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
    def test_nonfinite_scales(self):
        self.s[0,0]=np.inf
        with self.assertRaises(ValueError):validate(self.x,self.d,self.s)
if __name__=="__main__":unittest.main()
