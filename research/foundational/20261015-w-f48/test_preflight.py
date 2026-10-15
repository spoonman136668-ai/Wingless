import importlib.util, pathlib, sys, unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("w_f48_preflight",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Preflight(unittest.TestCase):
 def test_no_original_or_train_diagnostic_overlap(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.DIAG)&set(m.ORIGINAL))
  self.assertFalse(set(m.FAMILIES)&set(m.ORIGINAL_FAMILIES))
  for seed,family in ((412000,604),(391100,600),(411100,604)):
   with self.assertRaises(ValueError):m.source(seed,family,5)
 def test_learned_public_rank_is_metered(self):
  private,public,start,chain=m.source(411000,600,3)
  meter=[0]
  a=m.learned_action(public,start,chain,0.5,meter)
  self.assertTrue(a is None or a in m.old.f45.prior.prior.candidate_actions(public))
  self.assertGreater(meter[0],0)
 def test_paid_cap_and_original_isolation(self):
  private,public,start,chain=m.source(411000,601,5)
  p=m.old.f45.prior.prior.PurchasePort(private,public)
  a=m.old.f45.prior.prior.candidate_actions(public)
  p.buy(*a[0]);p.buy(*a[1])
  with self.assertRaises(ValueError):p.buy(*a[2])
 def test_only_train_grid_can_fit(self):
  self.assertEqual(len(m.BETAS),6)
  self.assertEqual(len(list(m.cases(m.FIT))),32)
  self.assertEqual(len(list(m.cases(m.DIAG))),32)
if __name__=="__main__":unittest.main()
