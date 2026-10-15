"""F42 prospective source-only tests; NEVER run whole F42 evaluation here."""
import importlib.util,pathlib,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("wf42",p);m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Gate(unittest.TestCase):
  def test_identity_isolation(self):
    self.assertFalse(set(m.SEEDS)&set(m.RESERVED))
    self.assertFalse(set(m.FAMILIES)&set(m.RESERVED_FAMILIES))
    with self.assertRaises(ValueError):m.case(322000,374,3)
  def test_private_map_truth_dual_oracle(self):
    secret,public,start,path=m.case(321000,370,3)
    self.assertEqual(len(public),24)
    self.assertEqual(m.prior.follow(secret,start,path),m.prior.independent_follow(secret,start,path))
  def test_public_sample_ranking_no_private_oracle(self):
    secret,public,start,path=m.case(321000,371,3)
    small,_,_,_,_=m.prior.score_actions(public,start,path,12,711)
    self.assertTrue(len(small)>0)
    self.assertTrue(all(len(a)==2 and a[0] in range(8) and a[1] in range(8) for a in small))
  def test_sample_rank_banks_are_independent(self):
    tokens=("b12","b48","b192","b192_independent")
    rank_seeds=[m.independent_rank_seed(101,t) for t in tokens]
    self.assertEqual(len(set(rank_seeds)),len(tokens))
    _,public,_,_=m.case(321000,370,3)
    a=m.prior.completions(public,12,rank_seeds[1])
    b=m.prior.completions(public,12,rank_seeds[2])
    self.assertNotEqual(a,b)
  def test_max_two_private_port(self):
    secret,public,start,path=m.case(321001,372,5)
    port=m.prior.Port(secret,public)
    choices=list(m.prior.score_actions(public,start,path,12,91)[0])
    port.buy(*choices[0]);port.buy(*choices[1])
    with self.assertRaises(ValueError):port.buy(*choices[2])
if __name__=="__main__":unittest.main()
