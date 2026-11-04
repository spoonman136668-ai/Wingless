"""CNC-02F independent FIT cohort, deterministic consensus and FP32 controls."""
import unittest
import numpy as np

from cnc02f_consensus_learning import (
    compare_consensus, ROWS, COLS, STEPS,
    FIT_A_SEED, FIT_B_SEED, TRANSFER_SEED,
)


class ConsensusLearningTests(unittest.TestCase):
    def test_fresh_disjoint_cohorts_and_no_transfer_selection(self):
        w = np.random.default_rng(29).standard_normal((ROWS, COLS)).astype(np.float32)
        record = compare_consensus(w)
        self.assertEqual(len({FIT_A_SEED, FIT_B_SEED, TRANSFER_SEED}), 3)
        self.assertEqual(record["step_count_each"], STEPS)
        self.assertEqual(record["examples_per_cohort"], 32)
        self.assertLessEqual(
            record["compressed_candidate_updates_accepted"] +
            record["compressed_candidate_updates_rejected"], STEPS)
        self.assertEqual(record["persistent_weight_bytes"]["direct_compressed"], 18456)
        self.assertEqual(record["persistent_weight_bytes"]["fp32_sgd"], ROWS*COLS*4)
        self.assertFalse(record["equal_compute_or_parameter_updates"])
        self.assertEqual(record["transient_peak_memory"], "UNMEASURED")
        self.assertEqual(record["original_model_activations"], "NOT_CAPTURED")
        self.assertFalse(record["full_model_training"])
        self.assertEqual(record["scientific_acceptance"], "NOT_GRANTED")
        for name in ("frozen_int4_transfer_error", "consensus_int4_transfer_error",
                     "fp32_sgd_transfer_error", "consensus_fit_a_error",
                     "consensus_fit_b_error"):
            self.assertTrue(np.isfinite(record[name]))
            self.assertGreaterEqual(record[name], 0)

    def test_reproducibility_and_fail_closed(self):
        w = np.random.default_rng(31).standard_normal((ROWS, COLS)).astype(np.float32)
        self.assertEqual(compare_consensus(w), compare_consensus(w))
        for source,steps in ((w[:1],STEPS), (np.zeros_like(w),STEPS),
                             (np.full_like(w,np.nan),STEPS),
                             (w,0),(w,1),(w,3),(w,True)):
            with self.subTest(shape=source.shape,steps=steps), self.assertRaises(ValueError):
                compare_consensus(source,steps=steps)

    def test_transfer_target_cannot_affect_update_choice(self):
        import inspect
        from cnc02f_consensus_learning import compare_consensus
        source = inspect.getsource(compare_consensus)
        proposal_line = source.index("proposal = direct.train_step(x_a, y_a)")
        validation_line = source.index("after_b = _error(direct.predict(x_b), y_b)")
        self.assertLess(proposal_line, validation_line)
        # The only transfer use occurs in post-update reporting.
        training_region = source[source.index("for _ in range(steps):"):
                                 source.index("if not np.array_equal(direct.scales")]
        self.assertNotIn("x_t", training_region)
        self.assertNotIn("y_t", training_region)


if __name__ == "__main__":
    unittest.main()
