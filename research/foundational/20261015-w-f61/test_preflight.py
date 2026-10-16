"""W-F61 frozen counterfactual public source focused-only qualification."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
sp=importlib.util.spec_from_file_location("w61_fit_source",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
class Tests(unittest.TestCase):
 def test_sealed_source_denial(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  with self.assertRaises(ValueError):m.source(751100,1444,8)
  with self.assertRaises(ValueError):m.source(752000,1448,8)
 def test_public_counterfactual_roles(self):
  c=m.source(751000,1440,8)
  self.assertEqual(len(c["public_bridges"]),3)
  self.assertEqual(len(c["secondary_chain"]),6)
  self.assertEqual(len(c["counterfactual_chain"]),6)
  self.assertTrue(all(len(x)==2 for x in c["public"]))
  self.assertNotEqual(c["secondary_chain"],c["counterfactual_chain"])
  outcomes,cost,n=m.evaluate(c)
  self.assertEqual(set(outcomes),set(m.ARMS))
  self.assertEqual(len(outcomes[m.ARMS[0]]),3)
  self.assertTrue(all(len(x)==3 for x in outcomes.values()))
  self.assertEqual(outcomes[m.ARMS[-1]],(c["truth"],c["second_truth"],c["counterfactual_truth"]))
  self.assertEqual(outcomes["strong_structured_MAP"],outcomes["frozen_F60_public_coupled_MAP"])
  self.assertTrue(1<=n<=4**3)
 def test_noncommutative_group(self):
  self.assertEqual(len(m.MAPS),512)
  self.assertTrue(any(m.compose(x,y)!=m.compose(y,x) for x in m.PERMS for y in m.PERMS))
if __name__=="__main__":unittest.main()
