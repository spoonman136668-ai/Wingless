"""G9K synthetic mechanism, storage and negative custody controls."""
import unittest
import numpy as np
from lowrank_residual_int2 import (
    RANKS, pack_lowrank_residual, apply_lowrank_residual, stored_bytes,
    correction_bytes_for_shape,
)
from packed_companded_int2 import pack_companded_int2, apply_companded_int2
from g9k_lowrank_budget import report


class G9KTests(unittest.TestCase):
    def test_synthetic_dense_lowrank_residual_direct_operator(self):
        rng=np.random.default_rng(401)
        w=rng.standard_normal((32,256)).astype(np.float32)
        x=rng.standard_normal((6,256)).astype(np.float32)
        y=x@w.T
        p,s=pack_companded_int2(w)
        baseline=np.sum((apply_companded_int2(x,p,s)-y).astype(np.float64)**2)
        for rank in RANKS:
            enc=pack_lowrank_residual(w,rank=rank)
            pred=apply_lowrank_residual(x,*enc,rank=rank)
            error=np.sum((pred-y).astype(np.float64)**2)
            self.assertLess(error,baseline)
            self.assertEqual(stored_bytes(*enc,256,rank=rank),
                p.nbytes+s.nbytes+24+2*rank*(32+256)+24)
    def test_zero_and_fail_closed(self):
        w=np.zeros((4,256),dtype=np.float32)
        enc=pack_lowrank_residual(w,rank=1)
        pred=apply_lowrank_residual(np.ones((2,256)),*enc,rank=1)
        self.assertTrue(np.all(pred==0))
        for bad in (0,4,True):
            with self.assertRaises(ValueError):pack_lowrank_residual(w,rank=bad)
        with self.assertRaises(ValueError):pack_lowrank_residual(np.full_like(w,np.inf))
        with self.assertRaises(ValueError):correction_bytes_for_shape((10,256),rank=4)
        with self.assertRaises(ValueError):apply_lowrank_residual(np.ones((2,256)),enc[0],enc[1],enc[2],enc[3][:,:0],rank=1)
    def test_full_model_frozen_budget_and_scope(self):
        r=report()
        self.assertEqual(r["parameter_total"],30532122624)
        self.assertEqual(r["indexed_tensor_count"],18867)
        self.assertEqual(r["baseline_remaining_bytes"],448974848)
        self.assertEqual([c["rank"] for c in r["cases"]],[1,2,3])
        self.assertFalse(r["model_forward"])
        self.assertFalse(r["original_qwen_weight_windows_used"])
        self.assertFalse(r["protected_heldout_used"])
        self.assertEqual(r["scientific_acceptance"],"NOT_GRANTED")
        for case in r["cases"]:
            self.assertTrue(case["analytical_12_gib_fit"])
            self.assertEqual(case["actual_gpu_memory"],"UNVERIFIED")
