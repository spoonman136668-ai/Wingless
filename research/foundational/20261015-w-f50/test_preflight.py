"""W-F50 focused tests: never execute prospective diagnostic cohort."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
sp=importlib.util.spec_from_file_location("w50_preflight",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
class Preflight(unittest.TestCase):
 def test_identity(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  with self.assertRaises(ValueError):m.source(502000,808,3)
  with self.assertRaises(ValueError):m.source(501100,800,3)
 def test_inverse_consistency(self):
  d=m.source(501000,800,5)
  for sign in (1,-1):
   for r,x,y in d["public"]:
    shift=(y-sign*x)%m.N
    self.assertEqual((sign*(y-shift))%m.N,x)
 def test_fit_labels_only(self):
  p,n=m.fit()
  self.assertEqual(n,16*4*2)
  self.assertEqual(len(p),8)
  self.assertTrue(all(0<val<1 for val in p.values()))
 def test_public_only_compose(self):
  d=m.source(501000,800,5)
  for sign in (1,-1):
   predicted=m.compose(d["context"],d["public"],d["chain"],d["start"],sign)
   self.assertIn(predicted,range(m.N))
if __name__=="__main__":unittest.main()
