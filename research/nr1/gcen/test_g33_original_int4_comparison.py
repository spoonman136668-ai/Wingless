"""CNC-01Q-G8 offline source, numeric, and custody-boundary controls.

No real Qwen weights are downloaded by this test file.
"""
from pathlib import Path
import unittest
from unittest.mock import patch

import numpy as np

from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS
from g33_original_int4_comparison import (
    BATCH, BLOCK, RESULT, compare_expert, compare_windows, run,
)


class G8RealInt4Tests(unittest.TestCase):
    def test_independent_svd_and_int4_synthetic_comparison(self):
        rng = np.random.default_rng(23)
        w = rng.standard_normal((ROWS, COLS)).astype(np.float32)
        x = rng.standard_normal((BATCH, COLS)).astype(np.float32)
        row = compare_expert(w, x)
        self.assertEqual(row["original_bf16_bytes"], ROWS * COLS * 2)
        self.assertEqual(row["svd_rank14_fp16_factor_bytes"], 14 * (ROWS + COLS + 1) * 2)
        self.assertEqual(row["packed_int4_bytes_including_metadata"], 18456)
        self.assertGreater(row["original_weight_energy"], 0)
        self.assertGreater(row["original_output_energy"], 0)
        self.assertLess(
            row["int4_output_squared_error"],
            row["svd_output_squared_error"]
        )
        self.assertLess(row["int4_weight_squared_error"] / row["original_weight_energy"], 0.08)

    def test_rejects_wrong_shapes_nonfinite_and_zero_energy(self):
        valid_x = np.ones((BATCH, COLS), dtype=np.float32)
        valid_w = np.ones((ROWS, COLS), dtype=np.float32)
        for w, x in [
            (valid_w[:, :-1], valid_x),
            (valid_w, valid_x[:, :-1]),
            (np.full_like(valid_w, np.nan), valid_x),
            (valid_w, np.full_like(valid_x, np.inf)),
            (np.zeros_like(valid_w), valid_x),
            (valid_w, np.zeros_like(valid_x)),
        ]:
            with self.subTest(wshape=w.shape, xshape=x.shape):
                with self.assertRaises(ValueError):
                    compare_expert(w, x)

    def test_frozen_4_8_16_group_aggregation_without_repeated_svd(self):
        fake = {
            "original_weight_energy": 2.0,
            "original_output_energy": 3.0,
            "int4_weight_squared_error": .04,
            "int4_output_squared_error": .06,
            "svd_output_squared_error": 1.5,
            "original_bf16_bytes": ROWS * COLS * 2,
            "packed_int4_bytes_including_metadata": 18456,
            "svd_rank14_fp16_factor_bytes": 10780,
        }
        windows = np.ones((max(COUNTS), ROWS, COLS), dtype=np.float32)
        with patch("g33_original_int4_comparison.compare_expert", return_value=fake) as probe:
            rows = compare_windows(windows, SEEDS[0])
        self.assertEqual(probe.call_count, max(COUNTS))
        self.assertEqual([r["experts"] for r in rows], list(COUNTS))
        self.assertEqual([r["seed"] for r in rows], [SEEDS[0]] * len(COUNTS))
        for row in rows:
            self.assertEqual(len(row["per_expert"]), row["experts"])
            self.assertAlmostEqual(row["int4_relative_weight_squared_error"], .02)
            self.assertAlmostEqual(row["int4_relative_output_squared_error"], .02)
            self.assertAlmostEqual(row["svd_rank14_relative_output_squared_error"], .5)
            self.assertEqual(row["packed_int4_bytes_including_metadata"], 18456 * row["experts"])
            self.assertEqual(row["sample_input_kind"], "SYNTHETIC_GAUSSIAN_NOT_NATIVE_ACTIVATIONS")
        with self.assertRaises(ValueError):
            compare_windows(windows, 9999)
        with self.assertRaises(ValueError):
            compare_windows(windows[:-1], SEEDS[0])

    def test_runner_source_reuses_strict_original_shard_custody(self):
        source = Path(__file__).with_name("g33_original_int4_comparison.py").read_text(encoding="utf-8")
        for requirement in (
            "check_upstream(manifest, upstream)",
            "stage_verified(response, stage, manifest[\"files\"][SHARD])",
            "inspect_header(stage, index)",
            "select_experts(header, index[\"weight_map\"])",
            "read_bf16_window(stream, width, offset, selected[\"shape\"], row0, col0)",
            "if len(comparisons) != len(SEEDS) * len(COUNTS):",
            '"scientific_acceptance": "NOT_GRANTED"',
            '"production_activation": False',
            '"protected_heldout_used": False',
            '"authentic_model_activations": "NOT_CAPTURED"',
        ):
            self.assertIn(requirement, source)
        for forbidden in (
            "from_pretrained(",
            ".generate(",
            "Qwen3MoeForCausalLM(",
            "load_checkpoint_and_dispatch(",
        ):
            self.assertNotIn(forbidden, source)
        self.assertEqual(BLOCK, 64)
        self.assertEqual(RESULT, "G8_ORIGINAL_BF16_WINDOW_INT4_DEVELOPMENT_COMPLETED")

    def test_runner_fails_closed_before_network_without_exact_sha(self):
        with self.assertRaisesRegex(ValueError, "exact SHA"):
            run(Path("nonexistent"), Path("."), Path("nonexistent-result"), "main")


if __name__ == "__main__":
    unittest.main()
