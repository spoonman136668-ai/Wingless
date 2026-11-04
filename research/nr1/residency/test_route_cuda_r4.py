"""R4 source-bound offline negative controls; no GPU access or torch import."""
import copy
import subprocess
import unittest
from unittest.mock import patch

from route_cuda_r4 import SCHEMA, MODES, SIZES_MIB, TRIALS, plan, validate_plan, validate_receipt, verify_source_identity


class R4OfflineTests(unittest.TestCase):
    @staticmethod
    def fixture():
        return {
            "schema": SCHEMA, "status": "R4_SYNTHETIC_CUDA_COPY_MEASURED",
            "cuda_available": True, "cuda_device_name": "synthetic-test-fixture",
            "source_sha": "a" * 40,
            "source_custody": {"source_sha": "a" * 40, "checkout_clean": True,
                               "source_identity_verified": True},
            "original_weights_used": False, "model_inference_executed": False,
            "scientific_acceptance": "NOT_GRANTED",
            "trials": [
                {"size_mib": size, "source_mode": mode, "repetition": rep,
                 "device_bytes": size * 1024**2, "correct_copy": True,
                 "host_synchronized_wall_ms": .1, "cuda_event_elapsed_ms": .01}
                for size, mode, rep in plan()
            ],
        }

    def test_exact_frozen_order(self):
        self.assertTrue(validate_plan(plan()))
        self.assertEqual(len(plan()), TRIALS * len(SIZES_MIB) * len(MODES))
        self.assertEqual({r for _, _, r in plan()}, set(range(TRIALS)))

    def test_missing_or_reordered_trial_rejected(self):
        for rows in (plan()[:-1], tuple(reversed(plan()))):
            with self.assertRaises(ValueError):
                validate_plan(rows)

    def test_valid_synthetic_receipt_fixture(self):
        self.assertTrue(validate_receipt(self.fixture()))

    def test_missing_or_wrong_parity_rejected(self):
        receipt = self.fixture()
        receipt["trials"][0]["correct_copy"] = False
        with self.assertRaisesRegex(ValueError, "parity"):
            validate_receipt(receipt)

    def test_invalid_timing_and_size_rejected(self):
        for key, invalid in (("cuda_event_elapsed_ms", 0),
                             ("host_synchronized_wall_ms", float("nan")),
                             ("device_bytes", 12)):
            receipt = self.fixture()
            receipt["trials"][0][key] = invalid
            with self.assertRaises(ValueError):
                validate_receipt(receipt)

    def test_false_inference_or_weights_claim_rejected(self):
        for key in ("original_weights_used", "model_inference_executed"):
            receipt = self.fixture()
            receipt[key] = True
            with self.assertRaises(ValueError):
                validate_receipt(receipt)

    def test_scope_and_device_fail_closed(self):
        for key, value in (("cuda_available", False),
                           ("cuda_device_name", ""),
                           ("scientific_acceptance", "PASS")):
            receipt = self.fixture()
            receipt[key] = value
            with self.assertRaises(ValueError):
                validate_receipt(receipt)

    def test_invalid_or_mismatched_source_identity(self):
        for value in ("", "g" * 40, "a" * 39):
            with self.assertRaises(ValueError):
                verify_source_identity(value)
        with patch("route_cuda_r4.subprocess.run", side_effect=[
            subprocess.CompletedProcess([], 0, stdout="b" * 40 + "\n"),
        ]):
            with self.assertRaisesRegex(RuntimeError, "SHA_MISMATCH"):
                verify_source_identity("a" * 40, repository_root=".")

    def test_dirty_or_untracked_checkout_rejected(self):
        with patch("route_cuda_r4.subprocess.run", side_effect=[
            subprocess.CompletedProcess([], 0, stdout="a" * 40 + "\n"),
            subprocess.CompletedProcess([], 0, stdout="?? research/other.py\n"),
        ]) as mocked:
            with self.assertRaisesRegex(RuntimeError, "DIRTY"):
                verify_source_identity("a" * 40, repository_root=".")
            self.assertEqual(mocked.call_count, 2)

    def test_clean_identity_accepted(self):
        with patch("route_cuda_r4.subprocess.run", side_effect=[
            subprocess.CompletedProcess([], 0, stdout="a" * 40 + "\n"),
            subprocess.CompletedProcess([], 0, stdout=""),
        ]):
            result = verify_source_identity("a" * 40, repository_root=".")
        self.assertTrue(result["checkout_clean"])
        self.assertTrue(result["source_identity_verified"])

    def test_unverified_receipt_provenance_rejected(self):
        for key, value in (("source_identity_verified", False),
                           ("checkout_clean", False), ("source_sha", "b" * 40)):
            receipt = self.fixture()
            receipt["source_custody"][key] = value
            with self.assertRaisesRegex(ValueError, "provenance"):
                validate_receipt(receipt)

    def test_source_import_requires_no_torch(self):
        import sys
        self.assertNotIn("torch", sys.modules)
        from route_cuda_r4 import run
        with self.assertRaisesRegex(ValueError, "immutable"):
            run("bad")


if __name__ == "__main__":
    unittest.main()
