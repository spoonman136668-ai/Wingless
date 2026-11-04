"""G9H fixed selection, encoded storage, outlier and source gate tests."""
import unittest
import numpy as np
from sparse_residual_int2 import (
    BUDGET_K, pack_sparse_residual, apply_sparse_residual,
    stored_bytes, correction_bytes_for_shape
)
from packed_companded_int2 import (
    pack_companded_int2, apply_companded_int2, stored_bytes as int2_bytes
)
from g9h_sparse_budget import sparse_budget


class G9HSparseTests(unittest.TestCase):
    def test_synthetic_sparse_outlier_correction_and_exact_byte_count(self):
        rng = np.random.default_rng(811)
        w = rng.standard_normal((16, 256)).astype(np.float32)
        w[:, 15] *= 20
        x = rng.standard_normal((7, 256)).astype(np.float32)
        ref = x @ w.T
        base, scales = pack_companded_int2(w, block=256)
        referr = float(np.sum((apply_companded_int2(x, base, scales,
                              block=256) - ref).astype(np.float64)**2))
        prior = referr
        for k in BUDGET_K:
            encoded = pack_sparse_residual(w, k=k, block=256)
            y = apply_sparse_residual(x, *encoded, k=k, block=256)
            err = float(np.sum((y - ref).astype(np.float64)**2))
            self.assertLess(err, prior)
            prior = err
            self.assertEqual(stored_bytes(*encoded, 256, k=k, block=256),
                int2_bytes(encoded[0], encoded[1], 256, block=256) +
                4*17 + 4*16*k)
            self.assertEqual(encoded[3].dtype, np.uint16)
            self.assertEqual(encoded[4].dtype, np.float16)
            self.assertEqual(encoded[2].dtype, np.uint32)

    def test_all_zero_exact_and_fail_closed(self):
        w = np.zeros((4, 256), dtype=np.float32)
        enc = pack_sparse_residual(w, k=1)
        x = np.ones((2, 256), dtype=np.float32)
        self.assertTrue(np.array_equal(
            apply_sparse_residual(x, *enc, k=1),
            np.zeros((2,4),dtype=np.float32)))
        for bad in (0, 4, 5, True):
            with self.assertRaises(ValueError):
                pack_sparse_residual(w, k=bad)
        with self.assertRaises(ValueError):
            pack_sparse_residual(np.full_like(w, np.inf))
        with self.assertRaises(ValueError):
            apply_sparse_residual(x, enc[0], enc[1], enc[2], enc[3]+256, enc[4], k=1)
        with self.assertRaises(ValueError):
            correction_bytes_for_shape((128,65536),k=1)

    def test_exact_model_census_and_budget_ceiling(self):
        report = sparse_budget()
        self.assertEqual(report["parameter_total"],30532122624)
        self.assertEqual(report["indexed_tensor_count"],18867)
        self.assertEqual(report["baseline_remaining_bytes"],448974848)
        self.assertEqual(len(report["cases"]),3)
        cases = report["cases"]
        self.assertTrue(cases[0]["analytical_fit"])
        self.assertTrue(cases[1]["analytical_fit"])
        self.assertFalse(cases[2]["model_behavior_verified"])
        for c in cases:
            self.assertEqual(c["model_total_planned_bytes"] +
                c["remaining_after_reserve_bytes"],12*2**30)
        self.assertFalse(report["actual_original_model_index_fetched"])
        self.assertFalse(report["original_qwen_activation_probes"])


if __name__=="__main__":
    unittest.main()
