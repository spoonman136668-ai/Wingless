"""Offline R5 causal admission and equality-of-budget controls."""
import unittest

from route_cache_r1 import CAPACITIES, EVAL, EXPERT_BYTES, FANOUT, PATTERNS, SEEDS, traces
from route_admission_r5 import (
    BASELINES, NEW_POLICY, RECENCY_WINDOW, adaptive_simulate, evaluate,
)


class R5Tests(unittest.TestCase):
    def test_all_frozen_scenarios_and_scope(self):
        receipt = evaluate()
        self.assertEqual(receipt["frozen_scenario_count"],
                         len(PATTERNS) * len(SEEDS) * len(CAPACITIES))
        self.assertEqual(len(receipt["cases"]), receipt["frozen_scenario_count"])
        self.assertEqual(receipt["baseline_policies"], list(BASELINES))
        self.assertEqual(receipt["candidate_policy"], NEW_POLICY)
        self.assertFalse(receipt["original_weights_used"])
        self.assertFalse(receipt["actual_router_traces_used"])
        self.assertFalse(receipt["gpu_transfer_measured"])
        self.assertEqual(receipt["scientific_acceptance"], "NOT_GRANTED")

    def test_prefix_invariance_no_future_peeking(self):
        training, full = traces("phase_shift", 11)
        part = full[:3]
        changed = part + tuple(tuple(range(100, 108)) for _ in range(EVAL - len(part)))
        original = adaptive_simulate(training, part, 16)
        alternate = adaptive_simulate(training, changed[:len(part)], 16)
        self.assertEqual(original["token_misses"], alternate["token_misses"])
        self.assertEqual(original["hits"], alternate["hits"])
        self.assertFalse(original["future_access_used"])

    def test_first_request_is_bypassed_before_repeat(self):
        row = tuple(range(FANOUT))
        training, _ = traces("uniform", 11)
        r = adaptive_simulate(training, (row,), 8)
        self.assertEqual(r["hits"], 0)
        self.assertEqual(r["admission_bypasses"], FANOUT)
        self.assertEqual(r["peak_resident_experts"], 0)

    def test_repeated_experts_can_be_admitted(self):
        row = tuple(range(FANOUT))
        training, _ = traces("uniform", 11)
        r = adaptive_simulate(training, (row, row, row), 8)
        self.assertEqual(r["hits"], FANOUT)
        self.assertEqual(r["misses"], FANOUT * 2)
        self.assertEqual(r["admission_bypasses"], FANOUT)

    def test_capacity_and_byte_conservation(self):
        for p in PATTERNS:
            training, evaluation = traces(p, 23)
            for cap in CAPACITIES:
                r = adaptive_simulate(training, evaluation, cap)
                self.assertLessEqual(r["peak_gpu_resident_bytes_assumed"], cap * EXPERT_BYTES)
                self.assertEqual(r["hits"] + r["misses"], EVAL * FANOUT)
                self.assertEqual(r["transferred_bytes"], r["misses"] * EXPERT_BYTES)
                self.assertFalse(r["measured_hardware_latency"])

    def test_invalid_capacity_and_window_fail_closed(self):
        training, evaluation = traces("uniform", 11)
        for cap in (0, 128, True, 1.5, -1):
            with self.subTest(cap=cap), self.assertRaises(ValueError):
                adaptive_simulate(training, evaluation, cap)
        for window in (0, 7, True, 4097):
            with self.subTest(window=window), self.assertRaises(ValueError):
                adaptive_simulate(training, evaluation, 16, window=window)

    def test_corrupt_or_duplicate_routing_rejected(self):
        training, evaluation = traces("uniform", 11)
        invalid = ((0,) * FANOUT, (-1,) + evaluation[0][1:], (128,) + evaluation[0][1:])
        for row in invalid:
            with self.subTest(row=row), self.assertRaises(ValueError):
                adaptive_simulate(training, (row,), 16)


if __name__ == "__main__":
    unittest.main()
