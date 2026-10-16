"""W-F58 source-only role reuse, Bayesian normalization and identity."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
sp=importlib.util.spec_from_file_location("w58_source",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
class Tests(unittest.TestCase):
 def test_protected_identity(self):
  self.assertTrue(set(m.FIT).isdisjoint(m.DIAG))
  with self.assertRaises(ValueError):m.source(681100,1304,10)
  with self.assertRaises(ValueError):m.source(682000,1308,8)
 def test_public_posterior_and_reuse(self):
  c=m.source(681000,1300,8)
  self.assertEqual(len(c["public"]),4)
  self.assertTrue(all(len(pairs)==3 and len(set(x for x,y in pairs))==3 for role,pairs in c["public"]))
  self.assertTrue(all(sum(role==r for role,d in c["chain"])>=2 for r in range(4)))
  probs=m.posterior(c)
  for p in probs:self.assertAlmostEqual(sum(p),1.,places=10)
 def test_computation_controls(self):
  c=m.source(681000,1300,10);post=m.posterior(c)
  self.assertEqual(m.exact(c,post)[1],14**4*10)
  self.assertEqual(m.beam(c,post)[1],16*10)
  self.assertEqual(m.factorized(c,post)[1],10*7*14)
if __name__=="__main__":unittest.main()
