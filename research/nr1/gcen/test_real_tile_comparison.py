"""Synthetic negative controls for the frozen real-tile reconstruction screen."""
import unittest
import numpy as np
from real_tile_comparison import compare_tile, EXPERTS, ROWS, COLS


class RealTileComparisonTests(unittest.TestCase):
    def test_finite_four_experts_and_matched_target(self):
        weights = np.random.default_rng(11).standard_normal((EXPERTS, ROWS, COLS)).astype(np.float32)
        result = compare_tile(weights)
        self.assertGreaterEqual(result["minimum_compression_ratio"], 4)
        self.assertEqual(result["original_fp16_bytes"], 16384)
        self.assertEqual(result["svd_factor_fp16_bytes"], 3880)
        self.assertEqual(result["shared_factor_fp16_bytes"], 3840)
        self.assertGreaterEqual(result["svd_relative_squared_error"], 0)
        self.assertGreaterEqual(result["shared_relative_squared_error"], 0)

    def test_bad_shape_rejected(self):
        with self.assertRaises(ValueError): compare_tile(np.zeros((1, ROWS, COLS), dtype=np.float32))

    def test_nonfinite_rejected(self):
        weights = np.ones((EXPERTS, ROWS, COLS), dtype=np.float32)
        weights[0, 0, 0] = np.nan
        with self.assertRaises(ValueError): compare_tile(weights)

    def test_zero_energy_rejected(self):
        with self.assertRaises(ValueError): compare_tile(np.zeros((EXPERTS, ROWS, COLS), dtype=np.float32))

    def test_deterministic(self):
        weights = np.random.default_rng(23).standard_normal((EXPERTS, ROWS, COLS)).astype(np.float32)
        self.assertEqual(compare_tile(weights), compare_tile(weights))


if __name__ == "__main__":
    unittest.main()
