"""G9U source/format negative controls; routine CI never invokes CUDA."""
import unittest
import numpy as np
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
if __name__=="__main__":unittest.main()
