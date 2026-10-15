"""W-F33 bounded TRAIN-only scientific source contract."""
import importlib.util,sys,unittest
from pathlib import Path
p=Path(__file__).with_name("execute.py")
spec=importlib.util.spec_from_file_location("wf33",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class TestSource(unittest.TestCase):
    def test_disjoint_source(self):
        self.assertFalse(set(m.TRAIN_SEEDS)&set(m.RESERVED_SEEDS))
        self.assertFalse(set(m.TRAIN_FAMILIES)&set(m.RESERVED_FAMILIES))
    def test_role_order_counterfactual_has_same_topology(self):
        a=m.make(191000,160,3,0);b=m.make(191000,160,3,2)
        self.assertEqual(a["edges"],b["edges"])
        self.assertEqual(a["source"],b["source"])
        self.assertEqual(a["target"],b["target"])
        self.assertNotEqual(a["query"],b["query"])
        self.assertTrue(m.coded_oracle(a))
        self.assertFalse(m.coded_oracle(b))
        self.assertEqual(m.no_role(a),m.no_role(b))
    def test_directed_edge_and_full_path(self):
        base=m.make(191004,162,4,0)
        reverse=m.make(191004,162,4,5)
        self.assertTrue(m.coded_oracle(base))
        self.assertFalse(m.coded_oracle(reverse))
    def test_bounded_no_truth_interface(self):
        w=m.initial_weights()
        self.assertEqual(len(w)*8,288)
        a=m.make(191005,163,5,0)
        score,path,reads=m.propagate(a,w,return_path=True)
        self.assertGreater(score,0)
        self.assertEqual(len(path),5)
        self.assertEqual(reads,len(a["query"])*len(a["edges"]))
if __name__=="__main__":unittest.main()
