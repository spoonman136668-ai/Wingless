"""CNC-02D frozen three-arm control contracts (no original Qwen download)."""
import unittest
import numpy as np

from cnc02d_matched_controls import (
    compare_controls, reconstruct, _fixed_scale_requantized,
    ROWS, COLS, FIT_SEED, TRANSFER_SEED, STEPS,
)
from packed_int4 import pack_int4, apply_packed_int4


class CNC02DControls(unittest.TestCase):
    def test_fixed_codebook_roundtrip_and_master_count(self):
        w = np.random.default_rng(91).standard_normal((ROWS, COLS)).astype(np.float32)
        codes, scales = pack_int4(w)
        recon = reconstruct(codes, scales)
        np.testing.assert_allclose(recon, apply_packed_int4(
            np.eye(COLS, dtype=np.float32), codes, scales).T, rtol=0, atol=0)
        np.testing.assert_array_equal(_fixed_scale_requantized(recon, scales), codes)

    def test_independent_fit_transfer_and_three_conventional_controls(self):
        w = np.random.default_rng(23).standard_normal((ROWS, COLS)).astype(np.float32)
        record = compare_controls(w)
        self.assertEqual(record["status"], "SYNTHETIC_ACTIVATION_DEVELOPMENT_COMPARISON_ONLY")
        self.assertEqual(record["steps_each"], STEPS)
        self.assertNotEqual(record["fit_seed"], record["transfer_seed"])
        self.assertEqual(record["fit_seed"], FIT_SEED)
        self.assertEqual(record["transfer_seed"], TRANSFER_SEED)
        self.assertEqual(record["gradient_budget_matching"], "UPDATE_COUNT_ONLY_NOT_MATCHED_FLOPS_OR_PARAMETERS")
        self.assertEqual(set(("frozen_int4","direct_int4","fp32_sgd","qat_ste")) -
                         set(record), set())
        for key in ("frozen_int4", "direct_int4", "fp32_sgd", "qat_ste"):
            for error in record[key].values():
                self.assertTrue(np.isfinite(error))
                self.assertGreaterEqual(error, 0)
        self.assertLessEqual(record["direct_int4"]["fit_relative_squared_error"],
                             record["frozen_int4"]["fit_relative_squared_error"] + 1e-7)
        self.assertEqual(record["persistent_weight_bytes"]["fp32_sgd"], ROWS*COLS*4)
        self.assertLess(record["persistent_weight_bytes"]["direct_int4"],
                        record["persistent_weight_bytes"]["fp32_sgd"])
        self.assertEqual(record["transient_peak_rss_bytes"], "UNMEASURED")
        self.assertFalse(record["full_model_training"])

    def test_fixed_controls_reproducible_and_fail_closed(self):
        w = np.random.default_rng(37).standard_normal((ROWS, COLS)).astype(np.float32)
        self.assertEqual(compare_controls(w), compare_controls(w))
        for source, steps in ((w[:1],STEPS), (np.zeros_like(w),STEPS),
                              (np.full_like(w, np.nan),STEPS), (w,0),
                              (w,1),(w,True),(w,3)):
            with self.subTest(shape=source.shape, steps=steps):
                with self.assertRaises(ValueError):
                    compare_controls(source, steps=steps)


if __name__ == "__main__":
    unittest.main()
