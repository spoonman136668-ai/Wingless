"""W-F59 safe D4 typed-binding public-source structural preflight."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
spec=importlib.util.spec_from_file_location("w59_source",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class Tests(unittest.TestCase):
 def test_protected_sources(self):
  with self.assertRaises(ValueError):m.source(701100,1344)
  with self.assertRaises(ValueError):m.source(702000,1348)
  self.assertFalse(set(m.FIT)&set(m.DIAG))
 def test_group_noncommutativity_and_public_closure(self):
  self.assertEqual(len(m.MAPS),512)
  self.assertTrue(any(m.compose(a,b)!=m.compose(b,a) for a in m.PERMS for b in m.PERMS))
  case=m.source(701000,1340)
  self.assertEqual(len(case["public"]),4)
  self.assertEqual(m.compose(case["private_maps"][3],m.compose(case["private_maps"][2],m.compose(case["private_maps"][1],case["private_maps"][0]))),m.IDENTITY)
 def test_all_public_controls_and_no_primary_use(self):
  case=m.source(701000,1340)
  guesses,work=m.evaluate(case)
  self.assertEqual(set(guesses),set(m.ARMS))
  self.assertEqual(work[m.ARMS[0]],512*(12+case["length"]))
  self.assertTrue(all(guesses[x] in range(4) for x in m.ARMS))
if __name__=="__main__":unittest.main()
