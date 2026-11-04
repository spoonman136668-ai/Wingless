"""CNC-G9G frozen robustness & mixed-capacity controls."""
import unittest
from g9g_robustness_screen import (
    BLOCKS, FAMILIES, SEEDS, analyze, compare_case, synthetic_weights
)


class G9GTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.report = analyze()

    def test_frozen_sources_and_stress_risk(self):
        r=self.report
        self.assertEqual(r["case_count"],18)
        self.assertEqual(len(r["family_metrics"]),6)
        self.assertEqual(len(r["mixed_precision_capacity"]),6)
        self.assertEqual(set(r["families"]),set(FAMILIES))
        self.assertEqual(r["seeds"],list(SEEDS))
        self.assertEqual(r["blocks"],list(BLOCKS))
        for row in r["cases"]:
            self.assertEqual(row["companded_int2_bytes"],row["uniform_int2_bytes"])
            self.assertLess(row["companded_int2_bytes"],row["int4_control_bytes"])
            for key in ("uniform_int2_relative_output_squared_error",
                        "companded_int2_relative_output_squared_error",
                        "int4_relative_output_squared_error"):
                self.assertGreaterEqual(row[key],0)
            self.assertEqual(row["bf16_reference_bytes"],65536)
        self.assertFalse(r["model_forward"])
        self.assertFalse(r["original_qwen_weights"])
        self.assertFalse(r["protected_heldout_used"])
        self.assertEqual(r["gpu_residency"],"UNVERIFIED")

    def test_exact_g9d_index_budget_and_guard(self):
        cases={r["name"]:r for r in self.report["mixed_precision_capacity"]}
        self.assertEqual(cases["all_int2_baseline"]["remaining_bytes"],448974848)
        self.assertEqual(cases["first_layer_int4"]["protected_experts"],128)
        self.assertLess(cases["first_layer_int4"]["remaining_bytes"],
                        cases["all_int2_baseline"]["remaining_bytes"])
        self.assertTrue(cases["both_boundary_layers_int4"]["analytical_fit"])
        self.assertFalse(cases["first_two_experts_each_layer_int4"]["analytical_fit"])
        self.assertEqual(self.report["source_identity"],"SYNTHETIC_FIXTURE_NOT_ORIGINAL_INDEX")

    def test_fail_closed_unfrozen_cases(self):
        for bad in (3,0,100,True):
            with self.assertRaises(ValueError):synthetic_weights(bad,"gaussian")
        with self.assertRaises(ValueError):synthetic_weights(SEEDS[0],"original")
        with self.assertRaises(ValueError):compare_case(SEEDS[0],"gaussian",64)


if __name__=="__main__":
    unittest.main()
