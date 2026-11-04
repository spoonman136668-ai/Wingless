"""CNC-G9D frozen mixed-precision storage screen, no weights or inference."""
import math
import unittest
from g9b_tensor_budget import FROZEN_CONFIG, expected_tensors
from g9d_mixed_precision_budget import mixed_report, _tensor_bytes, BASE, PROTECTED

SOURCE_SHA = "c" * 40


def original_index(config):
    shapes = expected_tensors(config)
    return {"weight_map": {
                name: f"model-{i%16+1:05d}-of-00016.safetensors"
                for i, name in enumerate(sorted(shapes))
            },
            "metadata": {"total_size": sum(2*math.prod(shape) for shape, _ in shapes.values())}}


class MixedPrecisionTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config = dict(FROZEN_CONFIG)
        cls.index = original_index(cls.config)
        cls.report = mixed_report(cls.config, cls.index, source_sha=SOURCE_SHA)

    def test_exact_pinned_budget_and_preservation_of_g9b(self):
        r = self.report
        self.assertEqual(r["source_index_tensor_count"], 18867)
        self.assertEqual(r["source_index_bf16_bytes"], 61064245248)
        self.assertEqual(r["vram_ceiling_bytes"], 12*2**30)
        cases = {v["name"]: v for v in r["scenarios"]}
        self.assertEqual(len(cases), 6)
        base = cases["all_int2_baseline"]
        self.assertEqual(base["remaining_bytes"], 448974848)
        self.assertEqual(base["protected_expert_count"], 0)
        self.assertEqual(base["extra_expert_bytes_vs_full_int2"], 0)
        self.assertTrue(base["analytical_fit_under_assumptions"])
        self.assertFalse(r["full_model_forward"])
        self.assertFalse(r["actual_gpu_measurement"])
        self.assertEqual(r["scientific_acceptance"], "NOT_GRANTED")

    def test_frozen_mask_cost_conservation_and_boundaries(self):
        cases = {v["name"]: v for v in self.report["scenarios"]}
        one = cases["first_layer_int4"]["extra_expert_bytes_vs_full_int2"]
        self.assertEqual(one, 128 * 1_437_696)
        self.assertEqual(one, cases["last_layer_int4"]["extra_expert_bytes_vs_full_int2"])
        self.assertEqual(cases["both_boundary_layers_int4"]["extra_expert_bytes_vs_full_int2"], 2*one)
        self.assertEqual(cases["first_expert_each_layer_int4"]["extra_expert_bytes_vs_full_int2"], 48 * 1_437_696)
        self.assertEqual(cases["first_two_experts_each_layer_int4"]["extra_expert_bytes_vs_full_int2"], 96 * 1_437_696)
        for name, case in cases.items():
            self.assertEqual(case["remaining_bytes"], cases["all_int2_baseline"]["remaining_bytes"] - case["extra_expert_bytes_vs_full_int2"])
            self.assertEqual(case["total_expert_count"], 6144)
            self.assertEqual(case["quality_or_router_sensitivity"], "UNVERIFIED")
            self.assertEqual(case["gpu_runtime_or_residency"], "UNVERIFIED")
        self.assertTrue(cases["both_boundary_layers_int4"]["analytical_fit_under_assumptions"])
        self.assertLess(cases["both_boundary_layers_int4"]["remaining_bytes"], 100 * (1<<20))

    def test_fails_closed_on_wrong_source_and_config(self):
        with self.assertRaises(ValueError):
            mixed_report(self.config, self.index, source_sha="latest")
        with self.assertRaises(ValueError):
            mixed_report(dict(self.config, hidden_size=1024), self.index, source_sha=SOURCE_SHA)
        with self.assertRaises(ValueError):
            _tensor_bytes((2,255), BASE)
        self.assertGreater(_tensor_bytes((768,2048), PROTECTED), _tensor_bytes((768,2048), BASE))


if __name__ == "__main__":
    unittest.main()
