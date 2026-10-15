"""W-F53 FIT-only preflight. No novel diagnostic source is consumed."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("w53_preflight",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Preflight(unittest.TestCase):
 def test_ancestry_disjoint(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  with self.assertRaises(ValueError):m.source(562000,988,8)
  with self.assertRaises(ValueError):m.source(561100,980,8)
 def test_real_fitted_value_model(self):
  weights,n=m.trained_model()
  self.assertEqual(len(weights),8)
  self.assertEqual(n,16*4*(6+8)*7)
  self.assertTrue(all(-127<=v<=127 for v in weights))
  self.assertTrue(any(v!=0 for v in weights))
 def test_bounded_compute_fit_case(self):
  d=m.source(561000,980,6);w,_=m.trained_model()
  for arm in m.ARMS[:-1]:
   prediction,cost,state=m.run_arm(d,arm,w)
   self.assertIn(prediction,range(7))
   self.assertGreater(cost,0)
  a,c,b=m.run_arm(d,m.ARMS[0],w)
  _,exact,eb=m.run_arm(d,"full_exact_public_bayes",w)
  self.assertLessEqual(c/exact,.5)
  self.assertLessEqual(b+len(w)+2,64)
if __name__=="__main__":unittest.main()
