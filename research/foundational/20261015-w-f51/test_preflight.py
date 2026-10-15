"""W-F51 focused tests do not consume novel diagnostic cohort."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
spec=importlib.util.spec_from_file_location("w51_preflight",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class Preflight(unittest.TestCase):
 def test_strict_source_and_laws(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  with self.assertRaises(ValueError):m.source(522000,908,3)
  with self.assertRaises(ValueError):m.source(521100,900,3)
  self.assertEqual(len(set(m.LAWS[:4])&set(m.LAWS[4:])),0)
 def test_public_calibration_legit(self):
  d=m.source(521000,900,3)
  self.assertEqual(len(d["calibration"]),4)
  for a in d["calibration"]:
   self.assertIn(m.infer_sign(a),(-1,1))
  self.assertEqual(len(d["public"]),4)
 def test_train_fit_only(self):
  w,n,k=m.fit()
  self.assertEqual(n,64)
  self.assertEqual(len(w),3)
  self.assertLessEqual(k,64)
 def test_inverse(self):
  d=m.source(521000,902,5)
  for sign in (-1,1):
   self.assertIn(m.compose(d["public"],d["chain"],d["start"],sign),range(7))
if __name__=="__main__":unittest.main()
