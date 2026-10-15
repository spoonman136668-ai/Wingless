"""Source preflight only: never run complete W-F45 TRAIN cohort."""
import importlib.util,pathlib,sys,unittest
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("w_f45_preflight",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class Preflight(unittest.TestCase):
  def test_registered_source_and_original_isolation(self):
      self.assertFalse(set(m.SEEDS)&set(m.RESERVED))
      self.assertFalse(set(m.FAMILIES)&set(m.RESERVED_FAMILIES))
      with self.assertRaises(ValueError):m.source(372000,464,5)
  def test_public_typed_structure_and_exact_support(self):
      private,public,start,chain=m.source(371000,460,5)
      self.assertEqual(len(public),12)
      self.assertEqual(chain[:3],((0,1),(1,-1),(0,-1)))
      self.assertIn(private,m.prior.prior.worlds(public))
      self.assertEqual(len(m.prior.prior.worlds(public)),1296)
  def test_one_step_public_only_and_cost_accounted(self):
      private,public,start,chain=m.source(371000,462,3)
      meter=[0]
      selected,margin=m.one_step(list(public),start,chain,meter)
      self.assertTrue(selected is None or selected in m.prior.prior.candidate_actions(public))
      self.assertGreaterEqual(margin,0)
      self.assertEqual(meter[0],1296*12)
  def test_purchase_cap(self):
      private,public,start,chain=m.source(371000,461,5)
      port=m.prior.prior.PurchasePort(private,public)
      a=m.prior.prior.candidate_actions(public)
      port.buy(*a[0]);port.buy(*a[1])
      with self.assertRaises(ValueError):port.buy(*a[2])
  def test_deep_gate_is_predeclared(self):
      self.assertEqual(m.MAX_DEEP_ABSOLUTE,3)
      self.assertLessEqual(m.MAX_DEEP_ABSOLUTE/(len(m.SEEDS)*len(m.FAMILIES)*len(m.LENGTHS)),m.MAX_DEEP_FRAC)
if __name__=="__main__":unittest.main()
