import importlib.util,json,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
sp=importlib.util.spec_from_file_location("w_f53",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
pr=json.loads(pathlib.Path(__file__).with_name("prereg.json").read_text())
class Qualification(unittest.TestCase):
 def test_source_cohorts_original_separation(self):
  self.assertEqual((m.FIT.start,m.FIT.stop),tuple(pr["fit_seeds"]))
  self.assertEqual((m.DIAG.start,m.DIAG.stop),tuple(pr["diagnostic_seeds"]))
  self.assertFalse(set(m.FIT)&set(m.DIAG));self.assertFalse(set(m.DIAG)&set(m.RESERVED))
  for args in ((582000,1108,3),(561100,1104,3),(581100,1100,3)):
   with self.assertRaises(ValueError):m.source(*args)
 def test_frozen_exact_fair_controls_and_full_law_bank(self):
  self.assertEqual(m.ARMS,(pr["candidate"],*pr["controls"],pr["unfair_control"]))
  self.assertEqual(len(m.LAWS),8)
  self.assertEqual(len(m.PATHS),8)
 def test_fit_only_minimax_calibration_and_public_actions(self):
  weights,n,correct=m.fit()
  self.assertEqual(n,64)
  params,families=m.fit_group_robust(weights)
  self.assertIn(params[0],(0.,.08,.16,.24,.32,.40))
  self.assertIn(params[1],("nearest","support","majority"))
  self.assertEqual(set(families),set(m.FIT_FAM))
  d=m.source(581000,1100,3)
  self.assertEqual(len(d["calibration"]),5)
  for arm in m.ARMS[:-1]:
   self.assertIn(m.predict(arm,d,weights,params),(-1,1))
 def test_audit_script_compiles_without_consuming_DIAG(self):
  q=pathlib.Path(__file__).with_name("audit_result.py")
  compile(q.read_text(),str(q),"exec")
if __name__=="__main__":unittest.main()
