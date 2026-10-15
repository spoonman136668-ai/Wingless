import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).parent/"source_pilot.py"
s=importlib.util.spec_from_file_location("wf46",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class F46(unittest.TestCase):
 def test_future_originals_inaccessible(self):
  for seed in (391000,392000,372000):
   with self.assertRaises(ValueError):m.source(seed,500,3)
  with self.assertRaises(ValueError):m.source(391100,504,3)
 def test_public_only_deterministic_choice(self):
  private,public,start,query=m.source(391100,500,3)
  self.assertEqual(len(public),12)
  a=m.select(list(public),start,query,[0],True)
  b=m.select(list(public),start,query,[0],True)
  self.assertEqual(a,b)
 def test_new_query_roles(self):
  self.assertNotEqual(m.source(391100,500,5)[3],m.source(391100,503,5)[3])
if __name__=="__main__":unittest.main()
