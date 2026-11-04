"""G9I original-weight independent dev seeds and sparse correction guards."""
import unittest
from pathlib import Path
from unittest.mock import patch
import numpy as np
from expert_scaling_comparison import ROWS,COLS,COUNTS,SEEDS as OLD_SEEDS
from g9i_original_sparse_comparison import (
    SEEDS, BUDGET_K, compare_expert, compare_windows, run
)

class G9IOriginalTests(unittest.TestCase):
    def test_separate_frozen_seed_family_and_numeric_operator(self):
        self.assertEqual(SEEDS,(107,163,229,271,347))
        self.assertFalse(set(SEEDS)&set(OLD_SEEDS))
        rng=np.random.default_rng(511)
        w=rng.standard_normal((ROWS,COLS)).astype(np.float32)
        w[:,7]*=16
        x=rng.standard_normal((32,COLS)).astype(np.float32)
        result=compare_expert(w,x)
        comp=result["companded"]["256"]
        weight_errors=[]
        for k in BUDGET_K:
            case=result["sparse"][str(k)]
            self.assertGreater(case["bytes"],comp["bytes"])
            self.assertLess(case["output_squared_error"],comp["output_squared_error"])
            self.assertLess(case["weight_squared_error"],comp["weight_squared_error"])
            weight_errors.append(case["weight_squared_error"])
        # Squared WEIGHT residual is monotone for top-absolute-error selection.
        # Output error need not be monotone because output terms can cancel.
        self.assertEqual(weight_errors,sorted(weight_errors,reverse=True))
    def test_frozen_original_custody_and_nonheldout(self):
        w=np.random.default_rng(507).standard_normal((16,ROWS,COLS)).astype(np.float32)
        sample=compare_expert(w[0],np.ones((32,COLS),dtype=np.float32))
        with patch("g9i_original_sparse_comparison.compare_expert",return_value=sample) as called:
            rows=compare_windows(w,SEEDS[0])
        self.assertEqual(called.call_count,16)
        self.assertEqual([r["experts"] for r in rows],list(COUNTS))
        self.assertTrue(all(len(r["sparse_cases_block256"])==3 for r in rows))
        with self.assertRaises(ValueError):compare_windows(w,11)
        with self.assertRaises(ValueError):compare_windows(w[:4],SEEDS[0])
        src=Path(__file__).with_name("g9i_original_sparse_comparison.py").read_text()
        for needed in (
            "check_upstream(manifest, build_manifest())",
            'stage_verified(response, staged, manifest["files"][SHARD])',
            'inspect_header(staged, index)',
            '"scientific_acceptance": "NOT_GRANTED"',
            '"authentic_model_activations": "NOT_CAPTURED"'):
            self.assertIn(needed,src)
        with self.assertRaisesRegex(ValueError,"source SHA"):
            run(Path("missing"),Path("."),Path("receipt"),"main")
