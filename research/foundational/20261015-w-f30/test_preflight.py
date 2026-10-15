"""Pre-original source and train-only interface tests. Never evaluate heldout."""
import importlib.util
import pathlib
import unittest
P=pathlib.Path(__file__).with_name("execute.py")
spec=importlib.util.spec_from_file_location("wf30",P)
mod=importlib.util.module_from_spec(spec)
spec.loader.exec_module(mod)

class Preflight(unittest.TestCase):
    def test_disjoint_seeds_and_structures(self):
        self.assertFalse(set(mod.TRAIN_SEEDS)&set(mod.TEST_SEEDS))
        self.assertFalse(set(mod.TRAIN_FAMILIES)&set(mod.TEST_FAMILIES))
    def test_novel_roles_and_candidate_count(self):
        ref,choices,idx=mod.make_case(161000,120,0,False)
        self.assertEqual(len(choices),8)
        self.assertIn(idx,range(8))
        self.assertTrue(all(max((r for a,b,r in g[1]),default=0)<4 for g in choices))
        self.assertEqual(len(mod.features(ref,choices[0])),mod.FEATURE_COUNT)
    def test_no_truth_in_scorer_and_deterministic_source(self):
        a=mod.make_case(161007,121,1,False)
        b=mod.make_case(161007,121,1,False)
        self.assertEqual(a,b)
        self.assertEqual(mod.classify(a[0],a[1],[0.]*mod.FEATURE_COUNT,"no_link"),-1)
    def test_graph_family_shift_is_real(self):
        self.assertEqual(mod.graph(161000,120,7,False)[0],8)
        self.assertEqual(mod.graph(161000,124,7,True)[0],9)
if __name__=="__main__":unittest.main()
