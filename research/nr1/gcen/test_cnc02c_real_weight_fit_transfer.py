"""CNC-02C original-shard adapter tests, offline only."""
from pathlib import Path
import unittest
from unittest.mock import patch
import numpy as np

from cnc02c_real_weight_fit_transfer import (
    evaluate_verified_windows, EXPERT_COUNT, ROWS, COLS, SCHEMA, run,
)
from expert_scaling_comparison import SEEDS


class CNC02CControls(unittest.TestCase):
    def test_frozen_expert_count_and_grouped_probe(self):
        w = np.random.default_rng(42).standard_normal(
            (EXPERT_COUNT, ROWS, COLS)).astype(np.float32)
        fake = {"scope": "SYNTHETIC_FIT_TRANSFER_ONLY",
                "compressed_direct_transfer_error": .2}
        with patch("cnc02c_real_weight_fit_transfer.compare_window",
                   return_value=fake) as probe:
            rows = evaluate_verified_windows(w, seed=SEEDS[0])
        self.assertEqual(probe.call_count, EXPERT_COUNT)
        self.assertEqual([r["original_expert_id"] for r in rows],
                         list(range(EXPERT_COUNT)))
        self.assertEqual([r["source_window_seed"] for r in rows],
                         [SEEDS[0]]*EXPERT_COUNT)
        self.assertEqual(rows[0]["scope"], "SYNTHETIC_FIT_TRANSFER_ONLY")
        for wrong in (w[:3], w[:, :, :-1],
                      np.full(w.shape, np.nan, dtype=np.float32)):
            with self.assertRaises(ValueError):
                evaluate_verified_windows(wrong, seed=SEEDS[0])
        with self.assertRaises(ValueError):
            evaluate_verified_windows(w, seed=99999)

    def test_source_custody_and_scientific_scope_not_weakened(self):
        src = Path(__file__).with_name(
            "cnc02c_real_weight_fit_transfer.py").read_text(encoding="utf-8")
        for token in (
            "check_upstream(manifest, build_manifest())",
            "stage_verified(response, stage, manifest[\"files\"][SHARD])",
            "inspect_header(stage, index)",
            "select_experts(header, index[\"weight_map\"])",
            "read_bf16_window(stream, header_bytes, offset,",
            "len(rows) != EXPERT_COUNT * len(SEEDS)",
            '"original_model_activations": "NOT_CAPTURED"',
            '"fp32_baseline_matched": False',
            '"scientific_acceptance": "NOT_GRANTED"',
            '"production_activation": False',
        ):
            self.assertIn(token, src)
        self.assertEqual(EXPERT_COUNT, 4)
        self.assertEqual(SCHEMA, "wingless.cnc02c.real-original-compressed-fit-transfer.v1")
        with self.assertRaisesRegex(ValueError, "exact source SHA"):
            run(Path("missing"), Path("."), Path("output"), "main")


if __name__ == "__main__":
    unittest.main()
