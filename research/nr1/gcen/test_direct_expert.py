"""Synthetic CNC-01Q algebraic equivalence and negative controls."""
import unittest
import numpy as np
from direct_expert import factorize, apply_direct, factor_storage_bytes


class DirectExpertTests(unittest.TestCase):
    def test_factorized_operator_matches_dense_reconstruction(self):
        rng = np.random.default_rng(17)
        for rows, cols, rank, batch in ((32, 64, 8, 13), (64, 32, 8, 1), (16, 16, 16, 9)):
            w = rng.standard_normal((rows, cols)).astype(np.float32)
            x = rng.standard_normal((batch, cols)).astype(np.float32)
            a, b = factorize(w, rank)
            direct = apply_direct(x, a, b)
            reconstructed = x @ (a @ b).T
            np.testing.assert_allclose(direct, reconstructed, rtol=2e-5, atol=2e-5)
            self.assertEqual(direct.shape, (batch, rows))

    def test_reference_distortion_and_storage(self):
        rng = np.random.default_rng(23)
        w = rng.standard_normal((128, 256)).astype(np.float32)
        x = rng.standard_normal((32, 256)).astype(np.float32)
        a, b = factorize(w, 14)
        direct = apply_direct(x, a, b)
        true_output = x @ w.T
        relative_error = np.linalg.norm(direct - true_output) / np.linalg.norm(true_output)
        self.assertTrue(np.isfinite(relative_error))
        self.assertGreaterEqual(relative_error, 0)
        self.assertEqual(factor_storage_bytes(a, b), 14 * (128 + 256) * 2)
        self.assertLess(factor_storage_bytes(a, b), w.size * 2 / 4)

    def test_input_validation(self):
        w = np.ones((4, 8), dtype=np.float32)
        for rank in (0, 5, True, 1.5):
            with self.assertRaises(ValueError): factorize(w, rank)
        with self.assertRaises(ValueError): factorize(np.array([[np.nan]]), 1)
        a, b = factorize(w, 2)
        with self.assertRaises(ValueError): apply_direct(np.ones((3, 7)), a, b)
        with self.assertRaises(ValueError): apply_direct(np.ones((3, 8)), a, np.full_like(b, np.inf))
        with self.assertRaises(ValueError): factor_storage_bytes(a, b, 0)


if __name__ == "__main__":
    unittest.main()
