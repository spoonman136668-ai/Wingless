"""W-F32 pre-original train-source controls. No evaluation seeds read."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("execute.py")
spec=importlib.util.spec_from_file_location("wf32",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class Check(unittest.TestCase):
    def test_disjoint_seed_family(self):
        self.assertFalse(set(m.TRAIN_SEEDS)&set(m.TEST_SEEDS))
        self.assertFalse(set(m.TRAIN_FAMILIES)&set(m.TEST_FAMILIES))
    def test_counterfactual_truth_train_only(self):
        for v in range(8):
            c=m.construct(171000,140,3,v)
            self.assertEqual(m.exact(c),v in (0,1,4,6),v)
    def test_role_order_matters(self):
        a=m.construct(171001,141,4,0)
        b=m.construct(171001,141,4,2)
        self.assertTrue(m.exact(a));self.assertFalse(m.exact(b))
        self.assertEqual(a["source"],b["source"])
        self.assertEqual(a["target"],b["target"])
        self.assertEqual(set(a["edges"]),set(b["edges"]), "paired counterfactual must preserve same edges")
    def test_state_cap_and_pure_inference(self):
        a=m.construct(171001,141,4,0)
        self.assertEqual(m.construct(171001,141,4,0),a)
        self.assertLessEqual(4*8+128,m.STATE_CEILING)
        self.assertEqual(len(m.local_features(1,1)),4)
if __name__=="__main__":unittest.main()
