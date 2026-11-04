"""R4 source-bound offline negative controls; no GPU access or torch import."""
import copy
import unittest

from route_cuda_r4 import SCHEMA, MODES, SIZES_MIB, TRIALS, plan, validate_plan, validate_receipt


class R4OfflineTests(unittest.TestCase):
    @staticmethod
    def fixture():
        return {
            "schema": SCHEMA, "status": "R4_SYNTHETIC_CUDA_COPY_MEASURED",
            "cuda_available": True, "cuda_device_name": "synthetic-test-fixture",
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

    def test_source_import_requires_no_torch(self):
        import sys
        self.assertNotIn("torch", sys.modules)
        from route_cuda_r4 import run
        with self.assertRaisesRegex(ValueError, "immutable"):
            run("bad")


if __name__ == "__main__":
    unittest.main()
