"""Source-free PG8 preregistration and invalid-input tests; no environment run."""
from __future__ import annotations

import unittest

import prospective_geometry_pg8 as pg8


class PG8ProtocolTests(unittest.TestCase):
    def test_frozen_families_and_seed_custody(self):
        pg8.locked_contract()
        self.assertEqual(pg8.TARGET_ENV, "MiniGrid-Empty-8x8-v0")
        self.assertEqual(pg8.TEST_STEPS, 100)
        self.assertEqual(pg8.TARGET_SEEDS, (18041, 18047, 18061, 18077,
                                            18089, 18101, 18119, 18133))

    def test_predeclared_positive_is_never_science(self):
        counts = {name: 0 for name in pg8.ARM_NAMES}
        counts["progress_native"] = 1
        counts["public_reactive"] = 1
        self.assertEqual(pg8.classify(counts), "PG8_DEVELOPMENT_POSITIVE_NOT_SCIENCE")

    def test_zero_goal_native_is_no_go(self):
        counts = {name: 0 for name in pg8.ARM_NAMES}
        counts["public_reactive"] = 0
        self.assertEqual(pg8.classify(counts), "PG8_NATIVE_TRANSFER_NO_GO")

    def test_stronger_control_requires_no_go(self):
        counts = {name: 1 for name in pg8.ARM_NAMES}
        counts["public_reactive"] = 2
        self.assertEqual(pg8.classify(counts), "PG8_NATIVE_TRANSFER_NO_GO")

    def test_extra_or_missing_control_rejected(self):
        counts = {name: 0 for name in pg8.ARM_NAMES}
        with self.assertRaisesRegex(Exception, "PG8_ARM_SET_DRIFT"):
            pg8.classify(counts)
        counts["public_reactive"] = 0
        counts["unpreregistered"] = 0
        with self.assertRaisesRegex(Exception, "PG8_ARM_SET_DRIFT"):
            pg8.classify(counts)

    def test_noninteger_count_rejected(self):
        counts = {name: 0 for name in pg8.ARM_NAMES}
        counts["public_reactive"] = True
        with self.assertRaisesRegex(Exception, "PG8_GOAL_COUNT_INVALID"):
            pg8.classify(counts)


if __name__ == "__main__":
    unittest.main()
