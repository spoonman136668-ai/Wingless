"""G28 bounded CPU/disk inference mapping negative controls; no real model loading."""
from types import SimpleNamespace
import unittest

from g28_original_device_plan import (
    GIB, DISK_BUDGET_BYTES, CPU_BUDGET_GIB,
    resolve_assignment, summarize_plan,
)
from test_g27_original_meta_mapping import synthetic_reference_keys
from out_of_core_preflight import SHARDS


class FakeMetaParameter:
    def __init__(self, count=1, device="meta"):
        self.count = count
        self.device = SimpleNamespace(type=device)
    def numel(self):
        return self.count


class G28DevicePlanTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.keys = synthetic_reference_keys()
        cls.index = {k: SHARDS[i % 16] for i, k in enumerate(cls.keys)}
        cls.tensors = [(name, FakeMetaParameter()) for name in cls.keys]

    def test_frozen_budget_bands(self):
        self.assertEqual(CPU_BUDGET_GIB, (4, 8, 12))
        self.assertEqual(DISK_BUDGET_BYTES, 256 * GIB)

    def test_longest_parent_mapping_not_arbitrary_parent(self):
        mapping = {"": "disk", "model.layers.0": "cpu", "model.layers.0.mlp": "disk"}
        self.assertEqual(resolve_assignment("model.layers.0.self_attn.q_proj.weight", mapping), "cpu")
        self.assertEqual(resolve_assignment("model.layers.0.mlp.experts.0.down_proj.weight", mapping), "disk")
        self.assertEqual(resolve_assignment("lm_head.weight", mapping), "disk")

    def test_cpu_and_disk_plan_reconciles_every_tensor(self):
        result = summarize_plan(self.tensors,
            {"": "disk", "model.layers.0": "cpu"},
            self.index, 4 * GIB, 80 * GIB, 2 * GIB)
        self.assertEqual(result["cpu_tensors"] + result["disk_tensors"], 18867)
        self.assertEqual(result["cpu_weight_bytes"] + result["disk_weight_bytes"], 18867 * 2)
        self.assertTrue(result["disk_tensors"] > 0)
        self.assertFalse(result["forward_and_disk_loading_verified"])

    def test_wrong_source_index_rejected(self):
        index = dict(self.index)
        index.pop("lm_head.weight")
        with self.assertRaises(ValueError):
            summarize_plan(self.tensors, {"": "disk"}, index, 4 * GIB, 80 * GIB, 2 * GIB)

    def test_excess_cpu_budget_fails_closed(self):
        with self.assertRaisesRegex(ValueError, "exceed proposed budget"):
            summarize_plan(self.tensors,
                {"": "disk", "model.layers.0": "cpu"},
                self.index, 1, 80 * GIB, 2 * GIB)

    def test_all_cpu_not_misreported_as_offload(self):
        with self.assertRaisesRegex(ValueError, "no disk offload"):
            summarize_plan(self.tensors, {"": "cpu"},
                self.index, 4 * GIB, 80 * GIB, 2 * GIB)

    def test_unmapped_checkpoint_parameter_rejected(self):
        with self.assertRaisesRegex(ValueError, "unmapped checkpoint tensor"):
            resolve_assignment("model.layers.11.mlp.experts.4.up_proj.weight", {"lm_head": "cpu"})

    def test_gpu_target_forbidden_on_cpu_host(self):
        with self.assertRaisesRegex(ValueError, "unexpectedly targets GPU"):
            resolve_assignment("model.layers.0.mlp.experts.0.down_proj.weight", {"": 0})

    def test_nested_gpu_target_forbidden_on_cpu_host(self):
        with self.assertRaisesRegex(ValueError, "unexpectedly targets GPU"):
            resolve_assignment("model.layers.0.mlp.experts.0.down_proj.weight",
                               {"": "disk", "model.layers.0": "cuda:0"})

    def test_non_meta_source_rejected(self):
        tensors = list(self.tensors)
        tensors[0] = (tensors[0][0], FakeMetaParameter(device="cpu"))
        with self.assertRaisesRegex(ValueError, "not meta"):
            summarize_plan(tensors, {"": "disk"}, self.index, 4 * GIB, 80 * GIB, 2 * GIB)

    def test_invalid_disk_free_values_rejected(self):
        for free in (-1, True, 1.5):
            with self.subTest(free=free), self.assertRaises(ValueError):
                summarize_plan(self.tensors, {"": "disk"}, self.index, 4 * GIB, free, 2 * GIB)

    def test_disk_stage_is_not_runtime_admission(self):
        result = summarize_plan(self.tensors, {"": "disk"},
                self.index, 4 * GIB, 1 * GIB, 2 * GIB)
        self.assertFalse(result["can_stage_full_checkpoint_plus_separate_offload_copy_and_workspace"])
        self.assertFalse(result["forward_and_disk_loading_verified"])


if __name__ == "__main__":
    unittest.main()
