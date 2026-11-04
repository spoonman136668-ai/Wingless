"""CNC-02G original BF16 weight-window consensus adapter controls."""
import unittest
from pathlib import Path
from unittest.mock import patch
import numpy as np

from cnc02c_real_weight_fit_transfer import (
    evaluate_verified_windows, EXPERT_COUNT, ROWS, COLS, run
)
from cnc02f_consensus_learning import compare_consensus
from expert_scaling_comparison import SEEDS


class ConsensusAuthenticAdapterTests(unittest.TestCase):
    def test_four_authentic_compatible_expert_windows(self):
        w = np.random.default_rng(67).standard_normal(
            (EXPERT_COUNT, ROWS, COLS)).astype(np.float32)
        records = evaluate_verified_windows(w, seed=SEEDS[0],
                                             comparator=compare_consensus)
        self.assertEqual(len(records), EXPERT_COUNT)
        self.assertEqual([r["original_expert_id"] for r in records],
                         list(range(EXPERT_COUNT)))
        for r in records:
            self.assertEqual(r["source_window_seed"], SEEDS[0])
            self.assertEqual(r["status"],
                             "FROZEN_FRESH_COHORT_SYNTHETIC_DEVELOPMENT_ONLY")
            self.assertEqual(r["scientific_acceptance"], "NOT_GRANTED")
            self.assertEqual(r["original_model_activations"], "NOT_CAPTURED")
        records[0]["status"] = "tampered"
        self.assertNotEqual(records[1]["status"], "tampered")

    def test_unknown_mode_and_comparator_fail_closed_before_network(self):
        w = np.ones((EXPERT_COUNT, ROWS, COLS), dtype=np.float32)
        with self.assertRaisesRegex(ValueError,"unfrozen experiment comparator"):
            evaluate_verified_windows(w, seed=SEEDS[0], comparator=lambda x: x)
        with self.assertRaisesRegex(ValueError,"unfrozen comparison mode"):
            run(Path("missing"),Path("."),Path("missing-receipt"),"a"*40,
                comparison_mode="unapproved")


if __name__ == "__main__":
    unittest.main()
