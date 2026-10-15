"""W-F56 source-only FIT calibration/noisy posterior proof."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
sp=importlib.util.spec_from_file_location("w56_fitonly",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
class Preflight(unittest.TestCase):
 def test_frozen_source_isolation(self):
  self.assertTrue(set(m.FIT).isdisjoint(m.DIAG))
  self.assertTrue(set(m.FIT_FAMILIES).isdisjoint(m.DIAG_FAMILIES))
  with self.assertRaises(ValueError):m.source(621100,1144,8)
  with self.assertRaises(ValueError):m.source(622000,1148,8)
 def test_normalized_public_role_posteriors(self):
  case=m.source(621000,1140,6)
  self.assertEqual(len(case["public"]),4)
  self.assertTrue(all(len(obs)==2 for role,obs in case["public"]))
  posterior=m.role_posterior(case["public"])
  self.assertEqual(len(posterior),4)
  for vals in posterior.values():
   self.assertEqual(len(vals),14)
   self.assertAlmostEqual(sum(vals),1,places=10)
 def test_exact_control_honest_work(self):
  case=m.source(621000,1140,6)
  exact,work=m.predictor(case,"full_exact_public_bayes")
  beam,partial=m.predictor(case,"budget_matched_fixed_public_posterior_beam")
  self.assertIn(exact,range(7))
  self.assertIn(beam,range(7))
  self.assertEqual(work,14**4*6)
  self.assertEqual(partial,16*6)
if __name__=="__main__":unittest.main()
