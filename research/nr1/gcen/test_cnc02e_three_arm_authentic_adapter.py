"""CNC-02E alternate original-weight runner: frozen three-arm admission tests."""
import unittest
from pathlib import Path
import numpy as np

from cnc02c_real_weight_fit_transfer import (
    evaluate_verified_windows, EXPERT_COUNT, ROWS, COLS, run,
)
from cnc02d_matched_controls import compare_controls
from expert_scaling_comparison import SEEDS


class ThreeArmAdapterTests(unittest.TestCase):
    def test_evaluates_four_experts_without_receipt_aliasing(self):
        w = np.random.default_rng(71).standard_normal(
            (EXPERT_COUNT, ROWS, COLS)).astype(np.float32)
        results = evaluate_verified_windows(w, seed=SEEDS[0],
                                             comparator=compare_controls)
        self.assertEqual(len(results), EXPERT_COUNT)
        self.assertEqual([r["original_expert_id"] for r in results],
                         list(range(EXPERT_COUNT)))
        self.assertEqual([r["source_window_seed"] for r in results],
                         [SEEDS[0]]*EXPERT_COUNT)
        for row in results:
            self.assertEqual(set(("frozen_int4", "direct_int4", "fp32_sgd", "qat_ste"))-set(row), set())
            self.assertEqual(row["gradient_budget_matching"],
                             "UPDATE_COUNT_ONLY_NOT_MATCHED_FLOPS_OR_PARAMETERS")
            self.assertEqual(row["transient_peak_rss_bytes"], "UNMEASURED")
            self.assertEqual(row["scientific_acceptance"], "NOT_GRANTED")
        results[0]["status"] = "MUTATED"
        self.assertNotEqual(results[1]["status"], "MUTATED")

    def test_rejects_unknown_adapter_and_mode_before_source_access(self):
        w = np.ones((EXPERT_COUNT, ROWS, COLS), dtype=np.float32)
        with self.assertRaises(ValueError):
            evaluate_verified_windows(w, seed=SEEDS[0], comparator=lambda w: None)
        with self.assertRaisesRegex(ValueError, "unfrozen comparison mode"):
            run(Path("nonexistent"), Path("."), Path("absent"), "a"*40,
                comparison_mode="promote")


if __name__ == "__main__":
    unittest.main()
