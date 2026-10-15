"""Pre-first-use public-interface/custody tests; no TRAIN-only cohort executed."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("train_only.py")
s=importlib.util.spec_from_file_location("wf36",P)
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class SourceContract(unittest.TestCase):
  def test_disjoint_original_and_past_source(self):
    self.assertFalse(set(m.TRAIN_SEEDS)&set(m.RESERVED_SEEDS))
    self.assertFalse(set(m.TRAIN_FAMILIES)&set(m.RESERVED_FAMILIES))
    with self.assertRaises(ValueError):m.family_graph(242000,214,4)
    with self.assertRaises(ValueError):m.family_graph(202000,210,4)
  def test_public_triples_and_dual_ordered_oracle(self):
    grammar,public,start,query=m.family_graph(241000,210,4)
    self.assertEqual(len(public),m.ROLES*m.DEMOS)
    self.assertEqual(len(query),4)
    self.assertEqual(m.target_oracle(grammar,start,query),m.independent_oracle(grammar,start,query))
    self.assertTrue(all(direction in (-1,1) for role,direction in query))
  def test_genuine_soft_quantization_and_inversion(self):
    _,public,start,query=m.family_graph(241000,211,6)
    q=m.learned_quantized(public)
    self.assertEqual(len(q),m.ROLES)
    self.assertTrue(all(len(q[r])==m.NODES for r in range(m.ROLES)))
    self.assertTrue(all(sum(row)==255 for role in q for row in role))
    self.assertIn(m.uncertain_compose(q,start,query),range(8))
    self.assertIn(m.uncertain_compose(q,start,tuple((r,-d) for r,d in query)),range(8))
  def test_fair_matched_demos_and_native_buffer(self):
    _,public,start,query=m.family_graph(241001,212,4)
    hard,affine,counts=m.hard_controls(public)
    self.assertEqual(len(hard),8)
    self.assertTrue(all(sorted(op)==list(range(8)) for op in hard))
    self.assertTrue(all(sorted(op)==list(range(8)) for op in affine))
    self.assertLessEqual(m.BYTES_Q8+m.BYTES_WORK,m.NATIVE_BYTES)
if __name__=="__main__":unittest.main()
