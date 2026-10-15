import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("w_f49",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class TestF49(unittest.TestCase):
 def test_preregistered_seeds_only(self):
  self.assertFalse(set(m.SEEDS)&set(m.RESERVED))
  for seed,family in ((462000,704),(411100,700),(461100,704)):
   with self.assertRaises(ValueError):m.source(seed,family,5)
 def test_two_hidden_worlds_consistent_public_observations(self):
  w,p,start,chain=m.source(461100,700,3)
  bank=m.p.worlds(p)
  self.assertEqual(len(bank),1296)
  for world in (bank[0],bank[-1]):
   self.assertTrue(all(world[r][x]==y for r,x,y in p))
 def test_four_candidates_and_public_prediction_only(self):
  _,pub,start,chain=m.source(461100,701,3)
  self.assertEqual(len(sorted(m.p.candidate_actions(pub))[:m.BEAM]),4)
  self.assertIn(m.outcome(m.p.worlds(pub)[0],pub,start,chain,None),range(6))
if __name__=="__main__":unittest.main()
