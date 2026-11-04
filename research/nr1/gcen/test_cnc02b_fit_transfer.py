"""Frozen CNC-02B developmental split and adaptation negative controls."""
import unittest
import numpy as np
from cnc02b_fit_transfer import compare_window, STEPS, ROWS, COLS


class CNC02BFitTransferTests(unittest.TestCase):
    def test_bounded_training_and_unseen_transfer_are_distinct(self):
        rng = np.random.default_rng(23)
        weight = rng.standard_normal((ROWS, COLS)).astype(np.float32)
        report = compare_window(weight)
        self.assertEqual(report["updates_allowed"], STEPS)
        self.assertEqual(report["fit_seed"], 6011)
        self.assertEqual(report["transfer_seed"], 17201)
        self.assertNotEqual(report["fit_seed"], report["transfer_seed"])
        self.assertEqual(report["fit_size"], 32)
        self.assertEqual(report["transfer_size"], 32)
        self.assertFalse(report["fp32_master_copy_retained"])
        self.assertFalse(report["full_model_training"])
        self.assertEqual(report["scientific_acceptance"], "NOT_GRANTED")
        self.assertLessEqual(report["compressed_direct_fit_error"],
                             report["frozen_int4_fit_error"] + 1e-6)
        self.assertTrue(np.isfinite(report["compressed_direct_transfer_error"]))
        self.assertTrue(np.isfinite(report["frozen_int4_transfer_error"]))
        self.assertLessEqual(report["code_edits"], STEPS)
        self.assertGreater(report["persistent_compressed_bytes"], 0)

    def test_rejects_invalid_source_or_update_budget(self):
        valid = np.ones((ROWS, COLS), dtype=np.float32)
        for weight in (valid[:1], np.full_like(valid, np.nan),
                       np.zeros_like(valid)):
            with self.assertRaises(ValueError):
                compare_window(weight)
        for steps in (0, STEPS+1, True, 1.5):
            with self.assertRaises(ValueError):
                compare_window(valid, steps=steps)

    def test_reproducible_frozen_sampling(self):
        weight = np.random.default_rng(91).standard_normal((ROWS, COLS)).astype(np.float32)
        a = compare_window(weight, steps=1)
        b = compare_window(weight, steps=1)
        self.assertEqual(a, b)


if __name__ == "__main__":
    unittest.main()
