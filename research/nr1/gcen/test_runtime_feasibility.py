"""G23 no-download static-memory capacity and custody-boundary tests."""
import json
from pathlib import Path
import unittest

from runtime_feasibility import (
    PINNED_ORIGINAL_SHARD_BYTES,
    capacity_assessment,
    check_manifest,
)


class RuntimeFeasibilityTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.manifest = json.loads((Path(__file__).parent / "upstream_full_manifest_unverified_weights.json").read_text())

    def test_pinned_manifest_exact_18_files(self):
        self.assertTrue(check_manifest(self.manifest))

    def test_missing_shard_fails(self):
        doc = dict(self.manifest, files=dict(self.manifest["files"]))
        del doc["files"]["model-00009-of-00016.safetensors"]
        with self.assertRaises(ValueError):
            check_manifest(doc)

    def test_wrong_revision_fails(self):
        with self.assertRaises(ValueError):
            check_manifest(dict(self.manifest, revision="0" * 40))

    def test_short_sha_fails(self):
        doc = dict(self.manifest, files=dict(self.manifest["files"]))
        doc["files"]["model-00001-of-00016.safetensors"] = "a"
        with self.assertRaises(ValueError):
            check_manifest(doc)

    def test_32gib_ram_plus_12gib_vram_insufficient(self):
        result = capacity_assessment(32 * 1024**3, 12 * 1024**3)
        self.assertEqual(result["status"], "INSUFFICIENT_FOR_SIMULTANEOUS_ORIGINAL_WEIGHT_RESIDENCY")
        self.assertGreater(PINNED_ORIGINAL_SHARD_BYTES, result["optimistic_combined_bytes"])
        self.assertEqual(result["offload_or_paging_feasibility"], "UNVERIFIED")
        self.assertFalse(result["execution_authorized_by_this_probe"])

    def test_unknown_gpu_does_not_claim_absence(self):
        result = capacity_assessment(16 * 1024**3, None)
        self.assertEqual(result["status"], "GPU_CAPACITY_UNKNOWN_NO_RESIDENCY_ADMISSION")
        self.assertIsNone(result["optimistic_combined_bytes"])

    def test_sufficient_optimistic_total_still_unqualified(self):
        result = capacity_assessment(128 * 1024**3, 12 * 1024**3)
        self.assertEqual(result["status"], "CAPACITY_LOWER_BOUND_MET_RUNTIME_UNVERIFIED")
        self.assertFalse(result["original_model_forward_verified"])
        self.assertFalse(result["runner_qualification_independently_verified"])

    def test_exact_floor_not_sufficient_for_model_inference(self):
        result = capacity_assessment(PINNED_ORIGINAL_SHARD_BYTES, 0)
        self.assertEqual(result["status"], "CAPACITY_LOWER_BOUND_MET_RUNTIME_UNVERIFIED")

    def test_reject_invalid_capacity_and_boolean(self):
        for ram, gpu in [(0, 100), (-1, None), (True, 1), (1, -1), (1, True), (1, 1.2)]:
            with self.subTest(ram=ram, gpu=gpu):
                with self.assertRaises(ValueError):
                    capacity_assessment(ram, gpu)


if __name__ == "__main__":
    unittest.main()
