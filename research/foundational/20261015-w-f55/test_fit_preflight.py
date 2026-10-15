"""W-F55 FIT-only source-constructive witness regression."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
s=importlib.util.spec_from_file_location("w55_public_fit",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Tests(unittest.TestCase):
 def test_source_boundaries(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  with self.assertRaises(ValueError):m.source(602000,1108,8)
  with self.assertRaises(ValueError):m.source(601100,1100,8)
 def test_exact_public_maps_and_inverse(self):
  case=m.source(601000,1100,6)
  ans,maps=m.constructive_public_only(case)
  self.assertEqual(ans,case["evaluator_truth"])
  self.assertEqual(len(maps),4)
  self.assertEqual(len(case["public_calibration"]),4)
  self.assertTrue(all(len(r["pairs"])==2 for r in case["public_calibration"]))
 def test_fit_only_frozen_witness(self):
  result=m.fit_only_report()
  self.assertEqual(result["train_fit_cases"],128)
  self.assertGreaterEqual(result["public_constructive_accuracy"],.85)
  self.assertTrue(result["disjoint_diagnostic_unused"])
if __name__=="__main__":unittest.main()
