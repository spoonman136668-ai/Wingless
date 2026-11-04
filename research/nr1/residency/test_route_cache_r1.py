"""Offline negative controls for frozen intact-expert residency simulation."""
import unittest
from route_cache_r1 import (
    N, FANOUT, TRAIN, EVAL, EXPERT_BYTES, BANDWIDTH_BPS,
    PATTERNS, SEEDS, CAPACITIES, POLICIES, traces, simulate, evaluate
)


class ResidencyTests(unittest.TestCase):
    def test_trace_determinism_and_nonduplicated_tokens(self):
        for pattern in PATTERNS:
            for seed in SEEDS:
                train, test = traces(pattern, seed)
                self.assertEqual((train, test), traces(pattern, seed))
                self.assertEqual((len(train), len(test)), (TRAIN, EVAL))
                self.assertTrue(all(len(row) == FANOUT and len(set(row)) == FANOUT
                                    and all(0 <= e < N for e in row)
                                    for row in train + test))

    def test_no_cache_control(self):
        training, evaluation = traces("hotset", 11)
        r = simulate(training, evaluation, "no_cache", 16)
        self.assertEqual(r["hit_count"], 0)
        self.assertEqual(r["miss_count"], EVAL * FANOUT)
        self.assertEqual(r["peak_gpu_resident_bytes_assumed"], 0)
        self.assertEqual(r["transferred_bytes"], EVAL * FANOUT * EXPERT_BYTES)

    def test_oracle_not_deployable(self):
        training, evaluation = traces("cyclic_adversarial", 11)
        for cap in CAPACITIES:
            optimal = simulate(training, evaluation, "hindsight_belady", cap)
            lru = simulate(training, evaluation, "lru", cap)
            self.assertLessEqual(optimal["miss_count"], lru["miss_count"])
            self.assertTrue(optimal["future_access_used"])
            self.assertFalse(lru["future_access_used"])

    def test_training_static_hotset_cannot_see_future_phase(self):
        training, evaluation = traces("phase_shift", 23)
        r = simulate(training, evaluation, "training_hotset", 16)
        self.assertGreater(r["miss_count"], 0)
        self.assertLessEqual(r["peak_gpu_resident_bytes_assumed"], 16 * EXPERT_BYTES)

    def test_cache_bounds_and_count_conservation(self):
        training, evaluation = traces("uniform", 37)
        for cap in CAPACITIES:
            for policy in POLICIES:
                with self.subTest(cap=cap, policy=policy):
                    r = simulate(training, evaluation, policy, cap)
                    self.assertLessEqual(r["peak_gpu_resident_bytes_assumed"], cap * EXPERT_BYTES)
                    self.assertEqual(r["hit_count"] + r["miss_count"], EVAL * FANOUT)
                    self.assertEqual(r["blocking_page_faults"], r["miss_count"])
                    self.assertEqual(r["transferred_bytes"], r["miss_count"] * EXPERT_BYTES)
                    self.assertFalse(r["runtime_measured"])

    def test_predicted_latency_includes_faults_and_setup(self):
        train, evaluation = traces("uniform", 11)
        r = simulate(train, evaluation, "no_cache", 8)
        per = 1.0 + FANOUT * (1000 * EXPERT_BYTES / BANDWIDTH_BPS + .15)
        self.assertAlmostEqual(r["predicted_total_ms"], EVAL * per)
        self.assertAlmostEqual(r["predicted_p95_token_ms"], per)

    def test_unrepresentable_capacity_and_bandwidth_fail(self):
        train, evaluation = traces("uniform", 11)
        for cap in (0, -1, True, 128, 1.5):
            with self.subTest(cap=cap), self.assertRaises(ValueError):
                simulate(train, evaluation, "lru", cap)
        for bandwidth in (0, -1, float("inf"), float("nan")):
            with self.subTest(bandwidth=bandwidth), self.assertRaises(ValueError):
                simulate(train, evaluation, "lru", 8, bandwidth=bandwidth)
        with self.assertRaises(ValueError):
            simulate(train, evaluation, "made_up_policy", 8)

    def test_bad_or_oracular_input_rejected(self):
        training, evaluation = traces("uniform", 11)
        for token in ((1,) * FANOUT, (-1,) + evaluation[0][1:],
                      (N,) + evaluation[0][1:], (True,) + evaluation[0][1:]):
            with self.subTest(token=token), self.assertRaises(ValueError):
                simulate(training, [token], "lru", 8)

    def test_all_frozen_scenarios_and_strict_scope(self):
        receipt = evaluate()
        self.assertEqual(len(receipt["scenarios"]), len(PATTERNS) * len(SEEDS) * len(CAPACITIES))
        self.assertEqual(receipt["oracle_role"], "NON_DEPLOYABLE_HINDSIGHT_BOUND")
        self.assertEqual(receipt["scientific_acceptance"], "NOT_GRANTED")
        self.assertFalse(receipt["original_weights_loaded"])
        self.assertFalse(receipt["authentic_router_traces"])
        self.assertFalse(receipt["full_model_inference"])
        self.assertFalse(receipt["protected_heldout_used"])


if __name__ == "__main__":
    unittest.main()
