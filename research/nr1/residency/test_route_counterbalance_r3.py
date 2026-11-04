"""R3 counterbalanced policy-order and negative controls."""
import tempfile
from pathlib import Path
import unittest

from route_cache_r1 import traces
from route_io_r2 import (BYTES_PER_EXPERT, CAPACITY, POLICIES,
                         create_corpus, frozen_requests)
from route_counterbalance_r3 import (PERMUTATIONS, PATTERNS, ROUNDS,
                                      SEEDS, balanced_order, checked_warmup,
                                      experiment, matched_trial)


class R3Tests(unittest.TestCase):
    def test_six_unique_orders(self):
        for seed in SEEDS:
            orders = balanced_order(seed)
            self.assertEqual(len(orders), 6)
            self.assertEqual(set(orders), set(PERMUTATIONS))

    def test_frozen_seed_and_round_guard(self):
        for rounds in (0, 1, 7, True):
            with self.subTest(rounds=rounds), self.assertRaises(ValueError):
                balanced_order(11, rounds)

    def test_missing_policy_or_duplicate_rejected(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            h = create_corpus(root)
            with self.assertRaises(ValueError):
                matched_trial(root, h, [], [0, 1], ("lru", "lru", "no_cache"))

    def test_warmup_rejects_corrupted_file(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            h = create_corpus(root)
            (root / "expert-004.bin").write_bytes(b"wrong")
            with self.assertRaises(ValueError):
                checked_warmup(root, h)

    def test_application_cache_is_cold_for_each_policy(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            h = create_corpus(root)
            checked_warmup(root, h)
            train, ev = traces("phase_shift", 11)
            request = frozen_requests(ev)
            controls = matched_trial(root, h, [e for row in train for e in row],
                                     request, POLICIES)
            self.assertEqual(controls["no_cache"]["hits"], 0)
            self.assertTrue(all(x["application_cache_initial_state"] == "EMPTY"
                                for x in controls.values()))
            self.assertTrue(all(x["application_peak_capacity_bytes"] ==
                                CAPACITY * BYTES_PER_EXPERT for x in controls.values()))

    def test_entire_frozen_matrix_and_scientific_boundary(self):
        receipt = experiment()
        self.assertEqual(len(receipt["records"]), len(PATTERNS) * len(SEEDS) * ROUNDS)
        self.assertEqual(len(receipt["summaries"]), len(PATTERNS) * len(POLICIES))
        for pattern in PATTERNS:
            for seed in SEEDS:
                group = [r for r in receipt["records"]
                         if r["pattern"] == pattern and r["seed"] == seed]
                self.assertEqual(len(group), 6)
                self.assertEqual({tuple(r["policy_order"]) for r in group}, set(PERMUTATIONS))
        self.assertEqual(receipt["os_page_cache_control"], "UNCONTROLLED")
        self.assertFalse(receipt["gpu_transfer_measured"])
        self.assertFalse(receipt["physical_disk_bytes_measured"])
        self.assertFalse(receipt["original_weights_used"])
        self.assertEqual(receipt["scientific_acceptance"], "NOT_GRANTED")


if __name__ == "__main__":
    unittest.main()
