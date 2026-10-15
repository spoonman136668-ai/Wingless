"""W-F52 focused tests only: do NOT execute new diagnostic cohort."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
spec=importlib.util.spec_from_file_location("wf52_preflight",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)
class Test(unittest.TestCase):
 def test_disjoint_source_and_frozen_family(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  with self.assertRaises(ValueError):m.source(542000,948,6)
  with self.assertRaises(ValueError):m.source(541100,940,6)
 def test_public_posterior_quantization(self):
  d=m.source(541000,940,4)
  for sign in (-1,1):
   ps=m.role_likelihoods(d,sign)
   self.assertEqual(len(ps),4)
   self.assertTrue(all(abs(sum(x)-1)<1e-11 for x in ps))
   val,cost=m.propagate(d,sign,True)
   self.assertAlmostEqual(sum(val),1,places=8)
   self.assertGreater(cost,0)
 def test_train_only_and_inference_budget(self):
  weights,count=m.train()
  self.assertEqual(count,128)
  self.assertEqual(len(weights),7)
  self.assertTrue(all(-127<=w<=127 for w in weights))
  d=m.source(541000,940,6)
  p,work,bytes_used=m.decide("trained_quantized_state",d,weights)
  q,exactwork,exactbytes=m.decide("full_exact_public_posterior",d,weights)
  self.assertLessEqual(work/exactwork,.5)
  self.assertEqual(bytes_used,7)
  self.assertIn(p,range(7))
  self.assertIn(q,range(7))
if __name__=="__main__":unittest.main()
