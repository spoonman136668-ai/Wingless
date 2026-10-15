"""W-F51 controls, source identities, and public operator inference."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("wf51",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Focused(unittest.TestCase):
 def test_fresh_family_seed_identity_and_originals_excluded(self):
  self.assertFalse(set(m.FIT)&set(m.DIAG))
  self.assertFalse(set(m.FIT_FAM)&set(m.DIAG_FAM))
  self.assertFalse(set(m.DIAG)&set(m.RESERVED))
  for seed,fam in ((502000,908),(522000,908),(501100,904),(521100,900)):
   with self.assertRaises(ValueError):m.source(seed,fam,3)
 def test_full_public_control_is_registered(self):
  self.assertEqual(m.ARMS,("learned_operator_posterior","always_shift","always_reflect",
   "full_five_bit_public_nonlinear_rule","three_bit_public_parity",
   "public_hash","private_operator_oracle_unfair"))
  for ctx in ((0,0,0,0,0),(1,1,0,0,1),(1,1,1,1,1)):
   expected=-1 if (ctx[0]&ctx[1])^(ctx[2]&ctx[3])^ctx[4] else 1
   self.assertEqual(m.guess("full_five_bit_public_nonlinear_rule",ctx,0,1,1,3),expected)
 def test_fit_only_and_inverse_recomposition(self):
  posterior,n=m.fit()
  self.assertEqual(n,64)
  self.assertEqual(len(posterior),32)
  d=m.source(521000,900,5)
  self.assertEqual(len(d["context"]),5)
  for sign in (-1,1):
   self.assertIn(m.compose(d["context"],d["public"],d["chain"],d["start"],sign),range(7))
if __name__=="__main__":unittest.main()
