"""W-F54 FIT-only hypothesis/preflight; diagnostic data never touched."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
spec=importlib.util.spec_from_file_location("w54_fitonly",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class Tests(unittest.TestCase):
 def test_source_folds(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  with self.assertRaises(ValueError):m.source(582000,1048,8)
  with self.assertRaises(ValueError):m.source(581100,1040,8)
 def test_public_transition_likelihood(self):
  d=m.source(581000,1040,6)
  for sign in (-1,1):
   p=m.likelihood(d,sign)
   self.assertEqual(len(p),4)
   self.assertTrue(all(abs(sum(z)-1)<1e-10 for z in p.values()))
   state,ops,trace=m.propagate(d,sign,"fixed_mass_two_beam")
   self.assertEqual(len(trace),6)
   self.assertGreater(ops,0)
 def test_true_fit_and_budget(self):
  params,codes,n=m.learner_fit()
  self.assertEqual(n,128)
  self.assertEqual(len(params),8)
  self.assertEqual(len(codes),8)
  w,_=m.prev.trained_model()
  d=m.source(581000,1040,8)
  a,work,memory=m.infer(d,m.ARMS[0],params,w)
  b,exact,size=m.infer(d,"exact_full_public_bayes",params,w)
  self.assertIn(a,range(7))
  self.assertIn(b,range(7))
  self.assertLessEqual(work/exact,.5)
  self.assertLessEqual(memory+8+2,64)
if __name__=="__main__":unittest.main()
