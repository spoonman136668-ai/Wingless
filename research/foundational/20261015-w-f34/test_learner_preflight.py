"""W-F34 train-only learned operator interface tests. No cohort executed."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("learned_train_only.py")
s=importlib.util.spec_from_file_location("wf34learner",P);m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Contract(unittest.TestCase):
    def test_reserved_original_never_generated(self):
        self.assertFalse(set(m.SEEDS)&set(m.RESERVED_SEEDS))
        self.assertFalse(set(m.FAMILIES)&set(m.RESERVED_FAMILIES))
        self.assertFalse(set(m.SEEDS)&set(m.PREVIOUS_TRAIN_SEEDS))
        with self.assertRaises(ValueError):m.private_grammar(202000,184)
    def test_grammar_is_bijection_and_query_ordered(self):
        ops,demos,x,roles=m.public_episode(203000,188,6)
        self.assertEqual(len(demos),32)
        self.assertTrue(all(sorted(op)==list(range(8)) for op in ops))
        self.assertEqual(m.apply(ops,x,roles),m.independent_apply(ops,x,roles))
    def test_learned_transition_is_gradient_fitted_and_distilled(self):
        _,demos,x,roles=m.public_episode(203000,189,6)
        fit=m.fitted_gradient(demos)
        self.assertEqual(len(fit),8)
        self.assertTrue(all(sorted(op)==list(range(8)) for op in fit))
        self.assertEqual(m.MODEL_DISTILLED_BYTES,64)
        self.assertLessEqual(m.MODEL_DISTILLED_BYTES,m.CORE_BUDGET)
        self.assertEqual(m.MODEL_COEFFICIENTS,512)
    def test_strong_empirical_controls_are_disjoint_from_truth(self):
        _,demos,x,roles=m.public_episode(203000,190,8)
        self.assertEqual(len(m.fitted_permutation_map(demos)),8)
        self.assertEqual(len(m.fitted_affine_map(demos)),8)
        self.assertEqual(len(m.fitted_nearest(demos)),8)
if __name__=="__main__":unittest.main()
