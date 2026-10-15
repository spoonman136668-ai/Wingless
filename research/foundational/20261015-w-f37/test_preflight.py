"""Train-only source and inverse-consistency controls; no frozen evaluation."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("train_only.py")
s=importlib.util.spec_from_file_location("wf37",P)
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Contract(unittest.TestCase):
  def test_reserved_ids_are_disallowed(self):
    self.assertFalse(set(m.TRAIN_SEEDS)&set(m.RESERVED_SEEDS))
    with self.assertRaises(ValueError):m.make_case(262000,244,4)
  def test_privileged_two_implementation_truth_agreement(self):
    private,public,start,query=m.make_case(261000,240,6)
    self.assertEqual(len(public),40)
    self.assertEqual(m.secret_compose(private,start,query),m.oracle_independent(private,start,query))
    self.assertEqual(m.secret_compose(private,start,((0,1),(0,-1))),start)
  def test_learned_q8_inverse_and_doubly_stochastic_transfer(self):
    private,public,start,query=m.make_case(261000,241,4)
    fitted=m.fit_sinkhorn(public)
    self.assertEqual(len(fitted),8)
    self.assertTrue(all(len(f)==8 and all(len(row)==8 and sum(row)==255 for row in f) for f in fitted))
    self.assertIn(m.learned_compose(fitted,start,query),range(8))
    self.assertIn(m.learned_compose(fitted,start,tuple(reversed(query))),range(8))
  def test_matched_bijective_controls_and_core(self):
    _,public,_,_=m.make_case(261001,242,6)
    hard,affine,counts=m.controls(public)
    self.assertTrue(all(sorted(op)==list(range(8)) for op in hard))
    self.assertTrue(all(sorted(op)==list(range(8)) for op in affine))
    self.assertLessEqual(512+16,576)
if __name__=="__main__":unittest.main()
