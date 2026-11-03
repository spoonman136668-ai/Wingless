"""Focused controls for larger original-weight reconstruction screening."""
import unittest
import numpy as np
from real_window_comparison import compare_window, EXPERTS, ROWS, COLS, SVD_RANK, SHARED_RANK


class RealWindowTests(unittest.TestCase):
    def test_matched_budget_and_errors(self):
        w = np.random.default_rng(11).standard_normal((EXPERTS, ROWS, COLS)).astype(np.float32)
        r = compare_window(w)
        self.assertEqual(r["svd_rank"], SVD_RANK)
        self.assertEqual(r["shared_rank"], SHARED_RANK)
        self.assertLess(r["svd_factor_fp16_bytes"], r["original_fp16_bytes"] / 4)
        self.assertLess(r["shared_factor_fp16_bytes"], r["original_fp16_bytes"] / 4)
        self.assertEqual(len(r["per_expert"]), EXPERTS)
        self.assertTrue(0 <= r["svd_normalized_squared_error"] <= 1.1)
        self.assertTrue(0 <= r["shared_normalized_squared_error"] <= 1.1)

    def test_rejects_wrong_shape(self):
        with self.assertRaises(ValueError):
            compare_window(np.ones((EXPERTS, ROWS, COLS - 1), dtype=np.float32))

    def test_rejects_nan(self):
        w = np.ones((EXPERTS, ROWS, COLS), dtype=np.float32)
        w[0, 1, 1] = np.nan
        with self.assertRaises(ValueError):
            compare_window(w)

    def test_rejects_zero_energy(self):
        with self.assertRaises(ValueError):
            compare_window(np.zeros((EXPERTS, ROWS, COLS), dtype=np.float32))


if __name__ == "__main__":
    unittest.main()
