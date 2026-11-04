"""Frozen negative controls for GCEN 4/8/16 expert-sharing screen."""
import unittest
import numpy as np
from expert_scaling_comparison import COUNTS, ROWS, COLS, rank_for_matched_bytes, evaluate_group


class ExpertScalingTests(unittest.TestCase):
    def test_discrete_factor_budgets(self):
        for experts in COUNTS:
            rank=rank_for_matched_bytes(experts)
            independent=2*experts*14*(ROWS+COLS+1)
            shared=2*(experts*ROWS+COLS)*rank
            self.assertLessEqual(abs(independent-shared)/max(independent,shared),.015)
            self.assertLessEqual(max(independent,shared),experts*ROWS*COLS*2/4)

    def test_finite_errors(self):
        w=np.random.default_rng(11).standard_normal((4,ROWS,COLS)).astype(np.float32)
        result=evaluate_group(w)
        self.assertEqual(result["experts"],4)
        self.assertGreaterEqual(result["independent_normalized_squared_error"],0)
        self.assertGreaterEqual(result["shared_normalized_squared_error"],0)
        self.assertEqual(result["peak_memory"],"UNMEASURED")

    def test_wrong_shape_rejected(self):
        with self.assertRaises(ValueError):
            evaluate_group(np.ones((5,ROWS,COLS),dtype=np.float32))

    def test_nonfinite_rejected(self):
        w=np.ones((4,ROWS,COLS),dtype=np.float32)
        w[0,0,0]=np.nan
        with self.assertRaises(ValueError): evaluate_group(w)

    def test_zero_energy_rejected(self):
        with self.assertRaises(ValueError): evaluate_group(np.zeros((4,ROWS,COLS),dtype=np.float32))

    def test_rank_invalid_count_rejected(self):
        with self.assertRaises(ValueError): rank_for_matched_bytes(7)


if __name__=="__main__": unittest.main()
