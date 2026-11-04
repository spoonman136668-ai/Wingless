"""Offline G32 source and negative controls; no actual source weights transferred."""
import io
from pathlib import Path
import unittest

import numpy as np
from g32_real_expert_scaling import (
    COLS, ROWS, EXPERT_IDS, SHARD, RESULT, read_bf16_window, select_experts
)


class G32RealExpertScalingTests(unittest.TestCase):
    def make_inventory(self, count=16):
        names = [f"model.layers.0.mlp.experts.{i}.down_proj.weight"
                 for i in range(count)]
        desc = {name: {"dtype": "BF16", "shape": [ROWS, COLS],
                       "data_offsets": [i * ROWS * COLS * 2,
                                        (i + 1) * ROWS * COLS * 2]}
                for i, name in enumerate(names)}
        weights = {name: SHARD for name in names}
        return desc, weights

    def test_deterministic_first_complete_layer(self):
        header, index = self.make_inventory()
        selected = select_experts(header, index)
        self.assertEqual(selected["layer"], 0)
        self.assertEqual(len(selected["names"]), len(EXPERT_IDS))
        self.assertEqual(selected["shape"], [ROWS, COLS])
        self.assertEqual(selected["offsets"][0], 0)

    def test_incomplete_original_group_rejected(self):
        header, index = self.make_inventory(15)
        with self.assertRaisesRegex(ValueError, "no complete"):
            select_experts(header, index)

    def test_wrong_source_shard_rejected(self):
        header, index = self.make_inventory()
        index[next(iter(index))] = "model-00002-of-00016.safetensors"
        with self.assertRaisesRegex(ValueError, "no complete"):
            select_experts(header, index)

    def test_duplicate_or_wrong_dtype_rejected(self):
        for changes in ({"dtype": "F16"}, {"shape": [ROWS, COLS - 1]},
                        {"data_offsets": [0, 3]}):
            with self.subTest(changes=changes):
                header, index = self.make_inventory()
                header[next(iter(header))] = {**header[next(iter(header))], **changes}
                with self.assertRaises(ValueError):
                    select_experts(header, index)

    def test_bf16_window_numeric_decode(self):
        raw = (0x3f80).to_bytes(2, "little") * (4 * 8)
        source = io.BytesIO(b"\x00" * 8 + b"{}\n" + raw)
        # 3 is the exact placeholder header byte count.
        values = read_bf16_window(source, 3, 0, [4, 8], 1, 2, rows=2, cols=3)
        self.assertTrue(np.array_equal(values, np.ones((2, 3), dtype=np.float32)))

    def test_bf16_window_rejects_truncation_and_nan(self):
        for raw, error in ((b"", "short"), ((0x7fc0).to_bytes(2, "little") * 8, "nonfinite")):
            with self.subTest(error=error):
                source = io.BytesIO(b"\x00" * 11 + raw)
                with self.assertRaisesRegex(ValueError, error):
                    read_bf16_window(source, 3, 0, [1, 8], 0, 0, rows=1, cols=8)

    def test_bounds_reject_noninteger_and_out_of_range(self):
        source = io.BytesIO(b"\x00" * 50)
        for kw in ({"row0": True}, {"row0": -1}, {"cols": 9}, {"rows": 2}):
            args = {"header_bytes": 3, "start": 0, "shape": [1, 8],
                    "row0": 0, "col0": 0, "rows": 1, "cols": 8}
            args.update(kw)
            with self.subTest(kw=kw), self.assertRaises(ValueError):
                read_bf16_window(source, **args)

    def test_exact_existing_comparison_and_no_model_forward(self):
        source = Path(__file__).with_name("g32_real_expert_scaling.py").read_text(encoding="utf-8")
        self.assertIn("from expert_scaling_comparison import COUNTS, COLS, ROWS, SEEDS, evaluate_group", source)
        self.assertIn("check_upstream(manifest, upstream)", source)
        self.assertIn("stage_verified(response, stage, manifest[\"files\"][SHARD])", source)
        self.assertIn("inspect_header(stage, index)", source)
        self.assertIn("entry = evaluate_group(windows[:experts])", source)
        for forbidden in ("from_pretrained(", ".generate(", "Qwen3MoeForCausalLM(",
                          "load_checkpoint_and_dispatch("):
            self.assertNotIn(forbidden, source)
        self.assertIn('"protected_heldout_used": False', source)
        self.assertIn('"scientific_acceptance": "NOT_GRANTED"', source)
        self.assertEqual(RESULT, "G32_REAL_BF16_EXPERT_WINDOW_SCALING_SCREEN_COMPLETE")

    def test_preexisting_legacy_step_replaced_not_duplicated(self):
        workflow = (Path(__file__).resolve().parents[3] / ".github" / "workflows" /
                    "native.yml").read_text(encoding="utf-8")
        self.assertEqual(workflow.count("  g32_real_expert_sharing:"), 1)
        self.assertIn("    needs: native", workflow.split("  g32_real_expert_sharing:", 1)[1])
        self.assertNotIn("      - name: GCEN one-shot 4-8-16 expert sharing scaling pilot", workflow)
        section = workflow.split("  g32_real_expert_sharing:", 1)[1]
        self.assertIn("[gcen-sharing-scale]", section)
        self.assertNotIn("continue-on-error:", section)
        self.assertIn("gcen-g32-real-expert-scaling.json", section)
        self.assertIn("if: always()", section)


if __name__ == "__main__":
    unittest.main()
