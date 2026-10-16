"""W-F61B source-only preflight: no new FIT, DIAG, or original generation."""
import importlib.util,json,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_feasibility.py")
sp=importlib.util.spec_from_file_location("w61b_fit_source",p)
m=importlib.util.module_from_spec(sp);sys.modules[sp.name]=m;sp.loader.exec_module(m)
v=json.loads(p.with_name("prereg-amendment.json").read_text())
class Tests(unittest.TestCase):
 def test_sealed_source_denial(self):
  self.assertEqual((m.FIT.start,m.FIT.stop),tuple(v["fit"]["seeds"]))
  self.assertEqual(tuple(m.FIT_FAMILIES),tuple(v["fit"]["families"]))
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  for seed,family in ((751000,1440),(751100,1444),(752000,1448),
                      (771100,1484),(772000,1488)):
   with self.assertRaises(ValueError):m.source(seed,family,8)
 def test_historical_control_identity(self):
  self.assertEqual(tuple(v["new_controls"]),m.ARMS)
  self.assertEqual(m.ARMS[2],"public_top2_role_beam")
  self.assertEqual(m.ARMS[3],"F60_MAP_algorithm_port_F61_observations")
  self.assertEqual(m.historical_F60().MAPS,m.MAPS)
  self.assertTrue(v["prior_preflight_generated_one_original_FIT_case"])
  self.assertFalse(v["future_fit_results_generated"])
 def test_public_synthetic_evaluator_without_fit_identity(self):
  # Hand-built public observations; no call to source() with an FIT seed.
  synthetic=dict(family=1480,public=tuple((((0,0),(1,1))) for _ in range(4)),
    public_bridges=((0,2,0,0),(1,3,1,1),(2,0,0,0)),
    start=0,chain=((0,1),(1,-1)),second_chain=((2,1),),
    counterfactual_start=1,counterfactual_chain=((3,-1),(0,1)))
  outcomes,cost,n=m.evaluate(synthetic)
  self.assertEqual(set(outcomes),set(m.ARMS))
  self.assertTrue(all(len(x)==3 for x in outcomes.values()))
  self.assertEqual(outcomes["strong_structured_MAP"],outcomes[m.ARMS[3]])
  self.assertGreater(cost["public_top2_role_beam"],0)
  self.assertTrue(1<=n<=64)
 def test_noncommutative_group(self):
  self.assertEqual(len(m.MAPS),512)
  self.assertTrue(any(m.compose(x,y)!=m.compose(y,x) for x in m.PERMS for y in m.PERMS))
if __name__=="__main__":unittest.main()
