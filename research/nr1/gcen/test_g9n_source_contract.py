"""G9N original-weight-free contract against Wingless G9B/G28 imports."""
import math
import unittest
from g9b_tensor_budget import FROZEN_CONFIG, MODEL, REVISION, expected_tensors
from g9n_source_contract import freeze_original_assignment

SHARDS = tuple(f"model-{i:05d}-of-00016.safetensors" for i in range(1, 17))


class OriginalSourceContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config = dict(FROZEN_CONFIG)
        shapes = expected_tensors(cls.config)
        cls.index = {
            "weight_map": {name: SHARDS[i % 16]
                           for i, name in enumerate(shapes)},
            "metadata": {"total_size": sum(2 * math.prod(shape)
                                           for shape, _ in shapes.values())}
        }
        cls.manifest = {"model": MODEL, "revision": REVISION,
                        "files": {k: "a" * 64 for k in (*SHARDS,
                            "config.json", "model.safetensors.index.json")}}
        cls.devices = {"": "disk", "model.layers.0": "cpu"}
        cls.kwargs = dict(manifest=cls.manifest, config=cls.config,
                          index=cls.index, device_map=cls.devices,
                          cpu_budget_bytes=4 * 1024**3,
                          disk_budget_bytes=256 * 1024**3)

    def test_pinned_original_source_and_budget(self):
        r = freeze_original_assignment(**self.kwargs)
        self.assertEqual(r["indexed_tensors"], 18867)
        self.assertEqual(len(r["target_devices"]), 18867)
        self.assertEqual(r["original_tensor_bytes"], 61064245248)
        self.assertEqual(r["target_devices"]["model.layers.0.input_layernorm.weight"], "cpu")
        self.assertEqual(r["target_devices"]["model.layers.1.input_layernorm.weight"], "disk")
        self.assertFalse(r["full_model_forward_executed"])
        self.assertEqual(r["scientific_acceptance"], "NOT_GRANTED")

    def test_rejects_unrecognized_configuration(self):
        with self.assertRaises(ValueError):
            freeze_original_assignment(**(self.kwargs |
                {"config": dict(self.config, num_experts=64)}))

    def test_rejects_index_missing_last_tensor(self):
        index = dict(self.index, weight_map=dict(self.index["weight_map"]))
        index["weight_map"].pop(next(reversed(index["weight_map"])))
        with self.assertRaises(ValueError):
            freeze_original_assignment(**(self.kwargs | {"index": index}))

    def test_rejects_missing_shard_custody(self):
        manifest = dict(self.manifest, files=dict(self.manifest["files"]))
        del manifest["files"][SHARDS[-1]]
        with self.assertRaises(ValueError):
            freeze_original_assignment(**(self.kwargs | {"manifest": manifest}))

    def test_rejects_non_cpu_disk(self):
        with self.assertRaises(ValueError):
            freeze_original_assignment(**(self.kwargs |
                {"device_map": {"": "cuda:0"}}))

    def test_rejects_too_small_cpu_budget(self):
        with self.assertRaises(ValueError):
            freeze_original_assignment(**(self.kwargs |
                {"cpu_budget_bytes": 1}))

    def test_rejects_no_disk_placement(self):
        with self.assertRaises(ValueError):
            freeze_original_assignment(**(self.kwargs |
                {"device_map": {"": "cpu"}, "cpu_budget_bytes": 128 * 1024**3}))


if __name__ == "__main__":
    unittest.main()
