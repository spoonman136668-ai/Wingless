"""W-F57 source-bound FIT-only focused checks."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
s=importlib.util.spec_from_file_location("wf57_source",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Checks(unittest.TestCase):
 def test_no_diagnostic_or_original(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  with self.assertRaises(ValueError):m.source(641100,1184,8)
  with self.assertRaises(ValueError):m.source(642000,1188,8)
 def test_public_coupling_and_exact(self):
  self.assertEqual(len(m.MAPS),1372)
  d=m.source(641000,1180,6)
  self.assertEqual(len(d["public"]),4)
  for arm in m.ARMS:
   guess,work=m.predict(d,arm)
   self.assertIn(guess,range(7))
   self.assertGreater(work,0)
  self.assertEqual(m.predict(d,m.ARMS[0])[1],1372*14)
if __name__=="__main__":unittest.main()
