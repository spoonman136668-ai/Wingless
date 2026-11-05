"""Frozen synthetic-only G9F controls: no authentic original weights accessed."""
import unittest
import numpy as np
from packed_companded_int2 import (
    BLOCKS,LEVELS,CUTOFF,pack_companded_int2,apply_companded_int2,stored_bytes)
from packed_int2 import pack_int2,apply_packed_int2,stored_bytes as oldbytes
from packed_int4 import pack_int4,apply_packed_int4,stored_bytes as int4bytes

class G9FFrozenTests(unittest.TestCase):
    def test_static_codebook_and_correct_storage(self):
        np.testing.assert_array_equal(LEVELS,np.array([-1.5104,-0.4528,0.4528,1.5104],dtype=np.float32))
        self.assertAlmostEqual(CUTOFF,float((LEVELS[2]+LEVELS[3])/2))
        w=np.random.default_rng(8).standard_normal((3,256)).astype(np.float32)
        x=np.random.default_rng(9).standard_normal((4,256)).astype(np.float32)
        for block in BLOCKS:
            p,s=pack_companded_int2(w,block=block)
            self.assertEqual(p.dtype,np.uint8)
            self.assertEqual(s.dtype,np.float16)
            self.assertEqual(p.shape,(3,64))
            self.assertEqual(s.shape,(3,256//block))
            self.assertEqual(stored_bytes(p,s,256,block=block),p.nbytes+s.nbytes+24)
            self.assertTrue(np.isfinite(apply_companded_int2(x,p,s,block=block)).all())

    def test_frozen_synthetic_comparison_vs_uniform_and_int4(self):
        rng=np.random.default_rng(23)
        w=rng.standard_normal((128,256)).astype(np.float32)
        x=rng.standard_normal((32,256)).astype(np.float32)
        ref=x@w.T
        den=float(np.sum(ref.astype(np.float64)**2))
        p4,s4=pack_int4(w,block=64)
        y4=apply_packed_int4(x,p4,s4,block=64)
        err4=float(np.sum((y4.astype(np.float64)-ref.astype(np.float64))**2)/den)
        for block in BLOCKS:
            p,s=pack_companded_int2(w,block=block)
            out=apply_companded_int2(x,p,s,block=block)
            q,t=pack_int2(w,block=block)
            old=apply_packed_int2(x,q,t,block=block)
            err=float(np.sum((out.astype(np.float64)-ref.astype(np.float64))**2)/den)
            older=float(np.sum((old.astype(np.float64)-ref.astype(np.float64))**2)/den)
            self.assertLess(err,older)
            self.assertLess(err,0.25)
            self.assertGreater(err,err4)
            self.assertEqual(stored_bytes(p,s,256,block=block),oldbytes(q,t,256,block=block))
            self.assertLess(stored_bytes(p,s,256,block=block),int4bytes(p4,s4,256,block=64))

    def test_zero_and_fail_closed(self):
        w=np.zeros((2,256),dtype=np.float32)
        p,s=pack_companded_int2(w)
        self.assertTrue(np.all(s==0))
        self.assertTrue(np.all(apply_companded_int2(np.ones((3,256),dtype=np.float32),p,s)==0))
        for block in (0,64,129,True):
            with self.subTest(block=block),self.assertRaises(ValueError):pack_companded_int2(w,block=block)
        for value in (np.nan,np.inf,1e35):
            with self.assertRaises(ValueError):pack_companded_int2(np.full_like(w,value))
        with self.assertRaises(ValueError):apply_companded_int2(np.ones((2,255)),p,s)
        with self.assertRaises(ValueError):apply_companded_int2(np.ones((2,256)),p,-np.ones_like(s))
        with self.assertRaises(ValueError):stored_bytes(p,s,256,block=128)
