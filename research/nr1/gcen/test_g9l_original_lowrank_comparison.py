"""G9L frozen original source controls; offline synthetic fixtures."""
import unittest
from pathlib import Path
from unittest.mock import patch
import numpy as np
from expert_scaling_comparison import ROWS,COLS,COUNTS,SEEDS as PRIOR_SEEDS
from g9l_original_lowrank_comparison import SEEDS,RANKS,compare_expert,compare_windows,run

class G9LOriginalTests(unittest.TestCase):
    def test_frozen_independent_original_selection(self):
        self.assertEqual(SEEDS,(541,563,587,613,641))
        self.assertFalse(set(SEEDS)&set(PRIOR_SEEDS))
        self.assertFalse(set(SEEDS)&{107,163,229,271,347})
    def test_direct_corrected_synthetic_fixture(self):
        rng=np.random.default_rng(809)
        w=rng.standard_normal((ROWS,COLS)).astype(np.float32)
        x=rng.standard_normal((32,COLS)).astype(np.float32)
        r=compare_expert(w,x)
        for rank in RANKS:
            c=r["lowrank"][str(rank)]
            self.assertGreater(c["bytes"],r["companded_int2_bytes"])
            self.assertLess(c["output_squared_error"],r["companded_squared_error"])
        self.assertEqual(len(r["lowrank"]),3)
    def test_16_experts_and_custody(self):
        rng=np.random.default_rng(901)
        w=rng.standard_normal((16,ROWS,COLS)).astype(np.float32)
        fixture=compare_expert(w[0],np.ones((32,COLS),dtype=np.float32))
        with patch("g9l_original_lowrank_comparison.compare_expert",return_value=fixture) as called:
            records=compare_windows(w,SEEDS[0])
        self.assertEqual(called.call_count,16)
        self.assertEqual([x["experts"] for x in records],list(COUNTS))
        self.assertTrue(all(len(x["lowrank_cases"])==3 for x in records))
        with self.assertRaises(ValueError):compare_windows(w,11)
        src=Path(__file__).with_name("g9l_original_lowrank_comparison.py").read_text()
        for required in ("check_upstream(manifest, build_manifest())",
           'stage_verified(response, staged, manifest["files"][SHARD])',
           'inspect_header(staged, index)',
           '"scientific_acceptance": "NOT_GRANTED"',
           '"authentic_model_activations": "NOT_CAPTURED"'):
            self.assertIn(required,src)
        with self.assertRaisesRegex(ValueError,"source SHA"):
            run(Path("missing"),Path("."),Path("missing"),"main")
