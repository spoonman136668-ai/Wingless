"""Offline fail-closed G9B architecture/budget tests (no network or weights)."""
import unittest
import math

from g9b_tensor_budget import (
    CONTEXT, EXPECTED_TENSORS, FROZEN_CONFIG, VRAM_CEILING,
    expected_tensors, make_report, quantized_expert_bytes, reconcile_index,
)

SOURCE_SHA = "b" * 40


def fixture():
    config = dict(FROZEN_CONFIG)
    shapes = expected_tensors(config)
    names = sorted(shapes)
    index = {"metadata": {"total_size": sum(
        2 * math.prod(shape) for shape, _ in shapes.values())},
        "weight_map": {key: f"model-{(i % 16)+1:05d}-of-00016.safetensors"
                       for i, key in enumerate(names)}}
    return config, shapes, index


class G9BTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config, cls.shapes, cls.index = fixture()

    def test_exact_18867_keys_source_total(self):
        self.assertEqual(len(self.shapes), EXPECTED_TENSORS)
        self.assertEqual(reconcile_index(self.shapes, self.index), 61_064_245_248)
        self.assertEqual(self.shapes["model.layers.47.mlp.experts.127.down_proj.weight"],
                         ((2048, 768), "experts"))

    def test_reject_config_drift_even_if_same_tensor_count(self):
        for key, other in (("vocab_size", 123), ("torch_dtype", "float32"),
                           ("tie_word_embeddings", True), ("head_dim", 64),
                           ("attention_bias", True), ("mlp_only_layers", [0])):
            with self.subTest(key=key), self.assertRaisesRegex(ValueError, "config mismatch"):
                expected_tensors(dict(self.config, **{key: other}))

    def test_reject_index_missing_extra_duplicate_shards_and_bytes(self):
        mismatches = []
        wrong = {"metadata": dict(self.index["metadata"]),
                 "weight_map": dict(self.index["weight_map"])}
        del wrong["weight_map"]["model.norm.weight"]
        mismatches.append(wrong)
        wrong = {"metadata": {"total_size": self.index["metadata"]["total_size"] - 2},
                 "weight_map": dict(self.index["weight_map"])}
        mismatches.append(wrong)
        wrong = {"metadata": dict(self.index["metadata"]),
                 "weight_map": dict(self.index["weight_map"])}
        wrong["weight_map"] = {k: "model-00001-of-00016.safetensors" for k in wrong["weight_map"]}
        mismatches.append(wrong)
        for altered in mismatches:
            with self.subTest(test=mismatches.index(altered)), self.assertRaises(ValueError):
                reconcile_index(self.shapes, altered)

    def test_2bit_block64_f32_expert_cost_and_headroom(self):
        report = make_report(self.config, self.index, source_sha=SOURCE_SHA)
        self.assertEqual(report["parameter_total"], 30_532_122_624)
        self.assertEqual(report["indexed_tensor_count"], 18_867)
        self.assertEqual(report["kv_cache_bytes_per_token"], 98_304)
        self.assertEqual(report["context_tokens"], CONTEXT)
        self.assertEqual(report["vram_target_bytes"], VRAM_CEILING)
        self.assertEqual(report["parameter_classes"]["experts"]["parameters"], 28_991_029_248)
        scenarios = {x["name"]: x for x in report["scenarios"]}
        self.assertEqual(scenarios["int2_block64_fp32"]["scales_bytes"], 1_811_939_328)
        self.assertLess(scenarios["int2_block64_fp32"]["remaining_bytes"], 0)
        self.assertLess(scenarios["int3_block64_fp32"]["remaining_bytes"], 0)
        self.assertLess(scenarios["int4_block64_fp32"]["remaining_bytes"], 0)
        self.assertGreater(scenarios["int2_block128_fp16"]["remaining_bytes"], 0)
        self.assertFalse(scenarios["int2_block128_fp16"]["real_gpu_residency_verified"])
        self.assertEqual(scenarios["int2_block128_fp16"]["operator_scope"], "NOT_IMPLEMENTED")
        self.assertFalse(report["full_model_forward"])
        self.assertFalse(report["deployment_or_activation"])

    def test_fail_closed_invalid_source_identity_and_storage(self):
        with self.assertRaisesRegex(ValueError, "exact source"):
            make_report(self.config, self.index, source_sha="main")
        with self.assertRaises(ValueError):
            quantized_expert_bytes(self.shapes, 5, 64, 4)
        with self.assertRaises(ValueError):
            quantized_expert_bytes(self.shapes, 3, 62, 4)


if __name__ == "__main__":
    unittest.main()
