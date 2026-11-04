"""Frozen G9J evidence admission tests: source and storage are immutable."""
import unittest
from g9j_pareto_decision import decision_report, G9I_BASELINE, SOURCE_SHA


class G9JDecisionTests(unittest.TestCase):
    def test_original_identity_and_no_science_claim(self):
        r = decision_report()
        self.assertEqual(r["source_sha"], SOURCE_SHA)
        self.assertEqual(r["authentic_original_artifact_id"], 13308383049)
        self.assertEqual(len(r["cases"]), 6)
        self.assertEqual(r["frozen_seeds"], [107,163,229,271,347])
        self.assertFalse(r["full_model_forward"])
        self.assertFalse(r["gpu_residency_measured"])
        self.assertFalse(r["original_qwen_activations_captured"])
        self.assertFalse(r["protected_heldout_used"])
        self.assertEqual(r["scientific_acceptance"], "NOT_GRANTED")

    def test_monotonic_accuracy_vs_storage_and_int4_gap(self):
        r = decision_report()
        arms = {x["method"]: x for x in r["cases"]}
        for name in ("sparse_k1", "sparse_k2", "sparse_k3"):
            self.assertGreater(arms[name]["relative_output_squared_error"],
                8*arms["int4_b64_fp32"]["relative_output_squared_error"])
            self.assertGreater(arms[name]["actual_window_stored_bytes"],
                arms["companded_int2_b256"]["actual_window_stored_bytes"])
            self.assertLess(arms[name]["relative_output_squared_error"],
                arms["companded_int2_b256"]["relative_output_squared_error"])
            self.assertGreater(arms[name]["whole_model_analytical_headroom_bytes"],0)
        self.assertEqual([arms[n]["actual_window_stored_bytes"]
            for n in ("companded_int2_b256","sparse_k1","sparse_k2","sparse_k3")],
            [135552,152000,160192,168384])

    def test_separate_development_decision_not_acceptance(self):
        r=decision_report()
        self.assertEqual(r["decision"],
            "NO_GO_DIRECT_MODEL_TRIAL_FOR_UNIFORM_K1_K2_K3_SPARSE_CORRECTIONS")
        self.assertEqual(r["comparison_scope"],
            "16_REAL_QWEN_BF16_EXPERT_WINDOWS_SYNTHETIC_ACTIVATION_INPUTS")


if __name__ == "__main__":
    unittest.main()
