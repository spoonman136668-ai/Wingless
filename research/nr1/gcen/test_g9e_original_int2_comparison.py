"""G9E original-window comparison frozen offline tests. No network access."""
import unittest
from pathlib import Path
from unittest.mock import patch

import numpy as np

from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS
from g9e_original_int2_comparison import compare_expert, compare_windows, run


class G9EOriginalInt2Tests(unittest.TestCase):
    def test_frozen_numerical_error_and_storage(self):
        rng = np.random.default_rng(31)
        w = rng.standard_normal((ROWS, COLS)).astype(np.float32)
        x = rng.standard_normal((32, COLS)).astype(np.float32)
        result = compare_expert(w, x)
        self.assertEqual(result["original_bf16_bytes"], 2*ROWS*COLS)
        self.assertEqual(result["int4_storage_bytes"], 18456)
        for block in (128, 256):
            q = result["int2"][str(block)]
            self.assertEqual(q["bytes"], ROWS*COLS//4 + ROWS*(COLS//block)*2 + 24)
            self.assertLess(q["bytes"], result["int4_storage_bytes"])
            self.assertGreater(q["output_squared_error"], result["int4_output_squared_error"])
            self.assertTrue(np.isfinite(q["weight_squared_error"]))

    def test_frozen_group_seeds_counts_and_negative_cases(self):
        w = np.random.default_rng(7).standard_normal((16, ROWS, COLS)).astype(np.float32)
        fake = compare_expert(w[0], np.ones((32, COLS), dtype=np.float32))
        with patch("g9e_original_int2_comparison.compare_expert", return_value=fake) as called:
            results = compare_windows(w, SEEDS[0])
        self.assertEqual(called.call_count, 16)
        self.assertEqual([v["experts"] for v in results], list(COUNTS))
        self.assertEqual([v["seed"] for v in results], [SEEDS[0]]*3)
        self.assertTrue(all(len(v["int2_cases"]) == 2 for v in results))
        with self.assertRaises(ValueError):
            compare_windows(w, 999)
        with self.assertRaises(ValueError):
            compare_windows(w[:8], SEEDS[0])

    def test_original_custody_boundary_and_no_model_forward(self):
        path = Path(__file__).with_name("g9e_original_int2_comparison.py")
        src = path.read_text(encoding="utf-8")
        for present in ("stage_verified(response, staged, manifest[\"files\"][SHARD])",
                        "inspect_header(staged, index)",
                        "select_experts(header, index[\"weight_map\"])",
                        "read_bf16_window(stream, header_bytes, off, selected[\"shape\"], row, col)",
                        "check_upstream(manifest, build_manifest())",
                        '"scientific_acceptance": "NOT_GRANTED"',
                        '"authentic_model_activations": "NOT_CAPTURED"',
                        '"protected_heldout_used": False'):
            self.assertIn(present, src)
        for forbidden in ("from_pretrained(", ".generate(", "load_checkpoint_and_dispatch("):
            self.assertNotIn(forbidden, src)
        with self.assertRaisesRegex(ValueError, "source SHA"):
            run(Path("missing"), Path("."), Path("missing-receipt"), "main")
