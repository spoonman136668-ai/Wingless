"""G9P frozen operator development controls without external checkpoint access."""
import unittest
import numpy as np
from g9p_original_expert_comparison import compare_original_expert, SEEDS


class G9PFullExpertTests(unittest.TestCase):
    def setUp(self):
        self.w=np.random.default_rng(59).normal(size=(8,256)).astype(np.float32)

    def test_matched_arms_and_exact_bf16_reference(self):
        r=compare_original_expert(self.w)
        self.assertEqual(r["comparison_count"],8)
        self.assertEqual(r["frozen_seeds"],list(SEEDS))
        self.assertEqual(r["probe_inputs"],"SYNTHETIC_GAUSSIAN_NOT_QWEN_ACTIVATIONS")
        for row in r["comparisons"]:
            self.assertGreater(row["encoded_bytes_including_accounted_metadata"],0)
            self.assertGreaterEqual(row["relative_output_squared_error"],0)
            self.assertGreaterEqual(row["warm_cpu_operator_median_seconds"],0)
            if row["arm"]=="original_bf16_float32_compute":
                self.assertEqual(row["relative_output_squared_error"],0.0)
                self.assertEqual(row["compression_ratio_vs_original_bf16"],1.0)

    def test_nonfrozen_seeds_rejected(self):
        with self.assertRaisesRegex(ValueError,"seeds"):
            compare_original_expert(self.w,seeds=(1,2))

    def test_nonfinite_original_rejected(self):
        w=self.w.copy()
        w[0,0]=np.nan
        with self.assertRaisesRegex(ValueError,"invalid"):
            compare_original_expert(w)

    def test_nonmultiple_256_rejected(self):
        with self.assertRaisesRegex(ValueError,"invalid"):
            compare_original_expert(self.w[:,:255])

    def test_zero_original_rejected(self):
        with self.assertRaisesRegex(ValueError,"invalid"):
            compare_original_expert(np.zeros_like(self.w))
if __name__=="__main__":
    unittest.main()
