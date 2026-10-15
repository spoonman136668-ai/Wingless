import importlib.util,json,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
sp=importlib.util.spec_from_file_location("w_f54",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
pr=json.loads(pathlib.Path(__file__).with_name("prereg.json").read_text())
class Tests(unittest.TestCase):
 def test_frozen_identity_and_exclusions(self):
  self.assertEqual((m.FIT.start,m.FIT.stop),tuple(pr["fit_seeds"]))
  self.assertEqual((m.DIAG.start,m.DIAG.stop),tuple(pr["diagnostic_seeds"]))
  self.assertEqual(list(m.ARMS),pr["decision_arms"])
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  for v in ((602000,1208,0),(601100,1200,0),(581100,1204,0)):
   with self.assertRaises(ValueError):m.source(*v)
 def test_public_actions_never_access_private_label(self):
  e=m.source(601000,1200,0)
  self.assertEqual(len(e["initial"]),3)
  fit,n=m.fit();self.assertEqual(n,64)
  for arm in m.ARMS[:-1]:
   chosen=m.select(arm,e["target"],e["initial"],fit)
   self.assertTrue(chosen is None or (chosen!=e["target"] and chosen not in
                                     [c for c,v in e["initial"]]))
  self.assertEqual(m.public_predict(e["target"],e["initial"],fit) in (0,1),True)
 def test_independent_audit_compiles(self):
  compile(pathlib.Path(__file__).with_name("audit_result.py").read_text(),"audit","exec")
if __name__=="__main__":unittest.main()
