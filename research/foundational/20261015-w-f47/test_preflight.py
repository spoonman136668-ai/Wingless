"""W-F47 source-only checks; no complete prospective fit/diagnostic execution."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("w47_preflight",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Preflight(unittest.TestCase):
  def test_disjoint_families_seeds_and_future_originals(self):
    self.assertFalse(set(m.FIT)&set(m.DIAG))
    self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
    self.assertFalse(set(m.DIAG)&set(m.RESERVED))
    self.assertFalse(set(m.DIAG_FAM)&set(m.RESERVED_FAM))
    with self.assertRaises(ValueError):m.source(422000,568,3)
    with self.assertRaises(ValueError):m.source(421100,560,5)
  def test_public_feature_basis_only(self):
    private,public,start,chain=m.source(421000,560,3)
    f=m.public_features(public,start,chain)
    self.assertEqual(len(f),m.FEATURE_DIM)
    self.assertTrue(all(isinstance(x,float) or isinstance(x,int) for x in f))
    self.assertTrue(all(abs(x)<=2. for x in f))
  def test_deep_budget_predeclared(self):
    self.assertLessEqual(m.LEARNED_DEEP_BUDGET,m.PREREG_DEEP_CEILING)
    self.assertEqual(m.LEARNED_DEEP_BUDGET,4)
  def test_one_source_terminal_not_holdout(self):
    private,public,start,chain=m.source(421000,561,3)
    answer,receipts,cost=m.run_arm("one_step_value",private,public,start,chain)
    self.assertIn(answer,range(6))
    self.assertLessEqual(len(receipts),2)
    self.assertGreater(cost,0)
    self.assertEqual(len({(r,x) for r,x,y in receipts}),len(receipts))
if __name__=="__main__":unittest.main()
