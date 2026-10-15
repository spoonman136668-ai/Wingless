import importlib.util,json,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
spec=importlib.util.spec_from_file_location("w_f52",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
pr=json.loads(pathlib.Path(__file__).with_name("prereg.json").read_text())
class Qualification(unittest.TestCase):
 def test_exact_registered_controls_and_hypothesis_bank(self):
  self.assertEqual(m.ARMS,(pr["candidate"],*pr["baseline_controls"],pr["unfair_control"]))
  self.assertEqual(tuple(m.LAWS),tuple(pr["public_law_hypothesis_bank"]))
  self.assertEqual(len(m.PATHS),8)
 def test_distinct_identity_and_rejected_originals(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.DIAG)&set(m.RESERVED))
  self.assertEqual([m.DIAG.start,m.DIAG.stop],pr["diagnostic_seeds"])
  for v in ((562000,1008,5),(521100,904,3),(561100,1000,3)):
   with self.assertRaises(ValueError):m.source(*v)
 def test_public_law_only_with_five_anchors(self):
  d=m.source(561000,1000,5)
  self.assertEqual(len(d["calibration"]),5)
  p=m.public_law_probability(d["context"],d["calibration"])
  self.assertTrue(0<p<1)
  for calibration in d["calibration"]:self.assertIn(m.infer_sign(calibration),(-1,1))
  for name in m.ARMS[:-1]:
   self.assertIn(m.predict(name,d,(0.,0.,0.)),(-1,1))
 def test_disjoint_fit_only_and_truth_composition(self):
  weights,n,correct=m.fit()
  self.assertEqual(n,8*4*2)
  self.assertEqual(len(weights),3)
  d=m.source(561000,1000,3)
  for sign in (-1,1):
   self.assertIn(m.compose(d["public"],d["chain"],d["start"],sign),range(m.N))
if __name__=="__main__":unittest.main()
