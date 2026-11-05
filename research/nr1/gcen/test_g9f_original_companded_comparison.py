"""G9F original source comparator frozen offline test; no original weights accessed."""
import unittest
from pathlib import Path
from unittest.mock import patch
import numpy as np
from expert_scaling_comparison import ROWS,COLS,COUNTS,SEEDS
from g9f_original_companded_comparison import compare_expert,compare_windows,run

class G9FOriginalTests(unittest.TestCase):
    def test_synthetic_companded_fidelity_and_equal_storage(self):
        rng=np.random.default_rng(23)
        w=rng.standard_normal((ROWS,COLS)).astype(np.float32)
        x=rng.standard_normal((32,COLS)).astype(np.float32)
        row=compare_expert(w,x)
        for block in (128,256):
            z=row["companded"][str(block)]
            p=row["int2"][str(block)]
            self.assertEqual(z["bytes"],p["bytes"])
            self.assertLess(z["output_squared_error"],p["output_squared_error"])
            self.assertGreater(z["output_squared_error"],row["int4_output_squared_error"])
    def test_frozen_original_windows_and_custody_boundaries(self):
        weights=np.random.default_rng(7).standard_normal((16,ROWS,COLS)).astype(np.float32)
        fake=compare_expert(weights[0],np.ones((32,COLS),dtype=np.float32))
        with patch("g9f_original_companded_comparison.compare_expert",return_value=fake) as probe:
            result=compare_windows(weights,SEEDS[0])
        self.assertEqual(probe.call_count,16)
        self.assertEqual([r["experts"] for r in result],list(COUNTS))
        for r in result:
            self.assertEqual(len(r["int2_cases"]),2)
            for c in r["int2_cases"]:
                self.assertIn("companded_relative_output_squared_error",c)
        with self.assertRaises(ValueError):compare_windows(weights,999)
        with self.assertRaises(ValueError):compare_windows(weights[:8],SEEDS[0])
        code=Path(__file__).with_name("g9f_original_companded_comparison.py").read_text()
        for required in ("check_upstream(manifest, build_manifest())",
           'stage_verified(response, staged, manifest["files"][SHARD])',
           'inspect_header(staged, index)','select_experts(header, index["weight_map"])',
           '"scientific_acceptance": "NOT_GRANTED"',
           '"authentic_model_activations": "NOT_CAPTURED"'):
            self.assertIn(required,code)
        with self.assertRaisesRegex(ValueError,"source SHA"):
            run(Path("nonexistent"),Path("."),Path("nonexistent"),"main")
