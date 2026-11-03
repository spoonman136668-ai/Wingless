"""G24 metadata and out-of-core scratch diagnostics: synthetic offline controls."""
import hashlib
import json
import unittest
from pathlib import Path

from out_of_core_preflight import (
    BASE, MAX_BYTES, PINNED_ORIGINAL_SHARD_BYTES, SHARDS, analyze, pinned_metadata,
)


class G24PreflightTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.pinned = json.loads((Path(__file__).parent / "upstream_full_manifest_unverified_weights.json").read_text())

    def setUp(self):
        self.manifest = dict(self.pinned, files=dict(self.pinned["files"]))
        self.config = {"model_type": "qwen3_moe", "num_hidden_layers": 48,
                       "num_experts": 128, "num_experts_per_tok": 8}
        self.index = {"metadata": {"total_size": PINNED_ORIGINAL_SHARD_BYTES - 1024},
                      "weight_map": {
                          f"model.layers.{i}.mlp.experts.0.down_proj.weight": name
                          for i, name in enumerate(SHARDS)
                      }}
        self.index["weight_map"]["model.layers.47.mlp.experts.127.down_proj.weight"] = SHARDS[-1]

    def test_valid_index_and_no_runtime_claim(self):
        result = analyze(self.config, self.index, self.manifest, 8 * 1024**3)
        self.assertEqual(result["expected_shards"], 16)
        self.assertTrue(result["single_shard_staging_capacity_present"])
        self.assertFalse(result["full_checkpoint_staging_capacity_present"])
        self.assertFalse(result["execution_authorized_by_this_probe"])
        self.assertEqual(result["full_model_forward"], "NOT_EXECUTED")

    def test_high_scratch_does_not_grant_runtime(self):
        result = analyze(self.config, self.index, self.manifest, 160 * 1024**3)
        self.assertTrue(result["full_checkpoint_plus_separate_offload_planning_capacity_present"])
        self.assertEqual(result["compatible_offload_engine"], "NOT_VERIFIED")
        self.assertEqual(result["expert_forward_hook"], "NOT_VERIFIED")

    def test_wrong_architecture_fails(self):
        self.config["model_type"] = "qwen3"
        with self.assertRaisesRegex(ValueError, "architecture"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_wrong_topology_fails(self):
        self.config["num_experts_per_tok"] = 4
        with self.assertRaisesRegex(ValueError, "topology"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_unknown_shard_fails(self):
        self.index["weight_map"]["new.weight"] = "model-00017-of-00016.safetensors"
        with self.assertRaisesRegex(ValueError, "unknown"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_missing_shard_coverage_fails(self):
        self.index["weight_map"]["model.layers.0.mlp.experts.0.down_proj.weight"] = SHARDS[1]
        with self.assertRaisesRegex(ValueError, "16 original"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_missing_last_expert_fails(self):
        del self.index["weight_map"]["model.layers.47.mlp.experts.127.down_proj.weight"]
        with self.assertRaisesRegex(ValueError, "tensor mapping"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_invalid_index_total_fails(self):
        self.index["metadata"]["total_size"] = -1
        with self.assertRaisesRegex(ValueError, "total"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_unrealistically_small_index_total_fails(self):
        self.index["metadata"]["total_size"] = 100
        with self.assertRaisesRegex(ValueError, "unexpectedly small"):
            analyze(self.config, self.index, self.manifest, 8 * 1024**3)

    def test_invalid_free_space_fails(self):
        with self.assertRaises(ValueError):
            analyze(self.config, self.index, self.manifest, -1)
        with self.assertRaises(ValueError):
            analyze(self.config, self.index, self.manifest, True)

    def test_pinned_metadata_bytes_hash_match(self):
        encoded = json.dumps(self.config, sort_keys=True).encode("utf-8")
        self.manifest["files"]["config.json"] = hashlib.sha256(encoded).hexdigest()
        def fetch(url, limit):
            self.assertEqual(url, f"{BASE}/config.json")
            self.assertEqual(limit, MAX_BYTES["config.json"])
            return encoded
        self.assertEqual(pinned_metadata("config.json", self.manifest, fetch), self.config)

    def test_pinned_metadata_tampering_fails(self):
        def fetch(url, limit):
            return b"{}"
        with self.assertRaisesRegex(ValueError, "checksum"):
            pinned_metadata("config.json", self.manifest, fetch)

    def test_malformed_json_rejected_even_with_matching_sha(self):
        raw = b"{invalid JSON"
        self.manifest["files"]["config.json"] = hashlib.sha256(raw).hexdigest()
        with self.assertRaises(json.JSONDecodeError):
            pinned_metadata("config.json", self.manifest, lambda url, limit: raw)

    def test_metadata_empty_bytes_rejected(self):
        with self.assertRaisesRegex(ValueError, "invalid or oversized"):
            pinned_metadata("config.json", self.manifest, lambda url, limit: b"")

    def test_weight_shard_metadata_fetch_rejected(self):
        with self.assertRaisesRegex(ValueError, "allowlisted"):
            pinned_metadata(SHARDS[0], self.manifest, lambda url, limit: b"x")

    def test_oversized_metadata_rejected(self):
        def fetch(url, limit):
            return b"x" * (limit + 1)
        with self.assertRaisesRegex(ValueError, "oversized"):
            pinned_metadata("config.json", self.manifest, fetch)


if __name__ == "__main__":
    unittest.main()
