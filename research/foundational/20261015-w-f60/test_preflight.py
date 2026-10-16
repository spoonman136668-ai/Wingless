"""W-F60 pre-firstuse D4 multi-query source/custody tests."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
t=importlib.util.spec_from_file_location("w60_source",p)
m=importlib.util.module_from_spec(t);sys.modules[t.name]=m;t.loader.exec_module(m)
class Tests(unittest.TestCase):
 def test_sources_disjoint(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  with self.assertRaises(ValueError):m.source(721100,1384)
  with self.assertRaises(ValueError):m.source(722000,1388)
 def test_public_bridges_and_pair_dependence(self):
  c=m.source(721000,1380,8)
  self.assertEqual(len(c["public_bridges"]),2)
  self.assertEqual(len(c["second_chain"]),6)
  self.assertTrue(all(len(x)==2 for x in c["public"]))
  self.assertNotEqual(tuple((r,-d) for r,d in reversed(c["chain"])),c["second_chain"])
  guesses,cost,n=m.evaluate(c)
  self.assertTrue(all(len(v)==2 for v in guesses.values()))
  self.assertEqual(tuple(guesses[m.ARMS[-1]]),(c["truth"],c["second_truth"]))
  self.assertTrue(1<=n<=16)
 def test_stable_noncommutative_structured_law(self):
  self.assertEqual(len(m.MAPS),512)
  self.assertTrue(any(m.compose(x,y)!=m.compose(y,x) for x in m.PERMS for y in m.PERMS))
if __name__=="__main__":unittest.main()
