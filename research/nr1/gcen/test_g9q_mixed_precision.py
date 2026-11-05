"""G9Q frozen precision assignment / equal-budget negative controls."""
import unittest
import numpy as np
from g9q_mixed_precision import pack_mixed, apply_mixed, stored_bytes


class G9QMixedTests(unittest.TestCase):
    def setUp(self):
        self.w=np.random.default_rng(33).normal(size=(16,256)).astype(np.float32)
        self.x=np.random.default_rng(77).normal(size=(3,256)).astype(np.float32)

    def test_ranked_and_uniform_are_complete_and_equal_budget(self):
        ranked=pack_mixed(self.w, denominator=8, policy="ranked")
        control=pack_mixed(self.w, denominator=8, policy="uniform")
        self.assertEqual(len(ranked[0]), 2)
        self.assertEqual(len(control[0]), 2)
        self.assertEqual(stored_bytes(ranked,16,256), stored_bytes(control,16,256))
        for state in (ranked, control):
            y=apply_mixed(self.x,state)
            self.assertEqual(y.shape,(3,16))
            self.assertTrue(np.isfinite(y).all())
            self.assertGreater(stored_bytes(state,16,256),0)
            self.assertLess(stored_bytes(state,16,256),self.w.size*2)

    def test_ranking_uses_only_weight_not_probe_inputs(self):
        a=pack_mixed(self.w,denominator=4,policy="ranked")
        b=pack_mixed(self.w.copy(),denominator=4,policy="ranked")
        self.assertTrue(np.array_equal(a[0], b[0]))
        self.assertTrue(np.array_equal(a[1], b[1]))

    def test_duplicate_partition_rejected(self):
        state=list(pack_mixed(self.w))
        state[1]=state[1].copy()
        state[1][0]=state[0][0]
        with self.assertRaisesRegex(ValueError,"partition"):
            apply_mixed(self.x,tuple(state))

    def test_nonfrozen_policy_and_fraction_rejected(self):
        with self.assertRaisesRegex(ValueError,"nonfrozen"):
            pack_mixed(self.w, denominator=3)
        with self.assertRaisesRegex(ValueError,"nonfrozen"):
            pack_mixed(self.w, policy="oracle")

    def test_nonfinite_source_or_input_rejected(self):
        w=self.w.copy()
        w[0,0]=np.nan
        with self.assertRaisesRegex(ValueError,"source"):
            pack_mixed(w)
        state=pack_mixed(self.w)
        x=self.x.copy()
        x[0,0]=np.inf
        with self.assertRaisesRegex(ValueError,"activation"):
            apply_mixed(x,state)


if __name__=="__main__":
    unittest.main()
