"""R2 real-file I/O tests, strictly synthetic."""
import hashlib
from pathlib import Path
import tempfile
import unittest

from route_io_r2 import (
    BYTES_PER_EXPERT, CORPUS, POLICIES, create_corpus, experiment, measure, payload
)


class R2Tests(unittest.TestCase):
    def test_corpus_exact_content_and_distinct_hashes(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            refs = create_corpus(root)
            self.assertEqual(len(refs), CORPUS)
            self.assertEqual(len(set(refs.values())), CORPUS)
            for expert in range(CORPUS):
                self.assertEqual(hashlib.sha256((root / f"expert-{expert:03d}.bin").read_bytes()).hexdigest(), refs[expert])

    def test_no_cache_misses_every_time(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            refs = create_corpus(root)
            r = measure(root, refs, [0, 1], [0, 0, 1, 0], "no_cache")
            self.assertEqual((r["hits"], r["misses"]), (0, 4))
            self.assertEqual(r["read_syscall_payload_bytes"], 4 * BYTES_PER_EXPERT)

    def test_lru_hits_on_repeated_expert(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            refs = create_corpus(root)
            r = measure(root, refs, [], [0, 0, 0, 1], "lru")
            self.assertEqual((r["hits"], r["misses"]), (2, 2))

    def test_static_hotset_bypasses_non_hot_experts(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            refs = create_corpus(root)
            r = measure(root, refs, [0] * 8, [20, 20, 0, 0], "training_hotset")
            self.assertEqual((r["hits"], r["misses"]), (1, 3))

    def test_bad_sha_rejected_without_acceptance(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            refs = create_corpus(root)
            (root / "expert-000.bin").write_bytes(b"corrupted")
            with self.assertRaisesRegex(ValueError, "tampered"):
                measure(root, refs, [], [0], "lru")

    def test_capacity_and_unknown_policy_fail(self):
        with tempfile.TemporaryDirectory() as td:
            root = Path(td)
            refs = create_corpus(root)
            for cap in (0, CORPUS, -1, True):
                with self.subTest(cap=cap), self.assertRaises(ValueError):
                    measure(root, refs, [], [0], "lru", cap)
            with self.assertRaises(ValueError):
                measure(root, refs, [], [0], "hindsight_oracle")

    def test_frozen_experiment_and_no_gpu_claim(self):
        r = experiment()
        self.assertEqual(len(r["observations"]), 27)
        self.assertEqual(r["scientific_acceptance"], "NOT_GRANTED")
        self.assertFalse(r["gpu_transfers_measured"])
        self.assertFalse(r["original_weights_used"])
        self.assertFalse(r["model_inference_executed"])
        self.assertEqual(r["os_page_cache_control"], "NOT_COLD_OR_CONTROLLED")
        self.assertTrue(all(x["hits"] + x["misses"] == x["requests"] for x in r["observations"]))


if __name__ == "__main__":
    unittest.main()
