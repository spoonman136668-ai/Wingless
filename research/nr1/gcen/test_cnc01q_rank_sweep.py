"""Focused controls for CNC-01Q-G6; no original checkpoint data."""
import json
import unittest

from cnc01q_rank_sweep import BATCHES, RANKS, compare


class RankSweepTests(unittest.TestCase):
    def test_frozen_rank_batch_crossover_and_quality_controls(self):
        receipt = compare(repeats=2)
        print("CNC01Q_RANK_SWEEP_JSON " + json.dumps(receipt, sort_keys=True), flush=True)
        self.assertEqual(receipt["status"], "DEVELOPMENT_ONLY_SYNTHETIC_RANK_SWEEP")
        self.assertEqual(receipt["scope"], "NO_REAL_QWEN_NO_GPU_NO_INFERENCE_NO_ACCEPTANCE")
        self.assertEqual(len(receipt["results"]), 2 * len(RANKS) * len(BATCHES))
        for row in receipt["results"]:
            self.assertGreater(row["direct_median_cpu_seconds"], 0)
            self.assertGreater(row["reconstructed_median_cpu_seconds"], 0)
            self.assertGreater(row["factor_fp16_bytes"], 0)
            self.assertGreaterEqual(row["relative_output_squared_error"], 0)
            self.assertGreaterEqual(row["relative_weight_squared_error"], 0)
        random_rows = [x for x in receipt["results"] if x["family"] == "full_rank_random" and x["batch"] == 32]
        random_rows.sort(key=lambda x: x["rank"])
        for before, after in zip(random_rows, random_rows[1:]):
            self.assertGreaterEqual(before["relative_weight_squared_error"] + 1e-6, after["relative_weight_squared_error"])
        # A genuinely low-rank matrix must be recoverable in our positive control.
        positive = [x for x in receipt["results"] if x["family"] == "rank8_positive_control" and x["rank"] == 8]
        self.assertEqual(len(positive), len(BATCHES))
        for row in positive:
            self.assertLess(row["relative_weight_squared_error"], 1e-8)
            self.assertLess(row["relative_output_squared_error"], 1e-8)
        self.assertLess(random_rows[0]["factor_fp16_bytes"], random_rows[0]["dense_fp16_bytes"])
        self.assertGreater(random_rows[-1]["factor_fp16_bytes"], random_rows[-1]["dense_fp16_bytes"])

    def test_invalid_repeat_budget(self):
        for value in (0, 31, -1, True, 1.5):
            with self.subTest(value=value), self.assertRaises(ValueError):
                compare(repeats=value)


if __name__ == "__main__":
    unittest.main()
