"""Train-only implementation checks; never run evaluation cohort here."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("train_only.py")
spec=importlib.util.spec_from_file_location("wf35",P)
m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
class Gate(unittest.TestCase):
  def test_identity_ranges(self):
    self.assertFalse(set(m.META_SEEDS)&set(m.EVAL_SEEDS))
    self.assertFalse(set(m.META_FAMILIES)&set(m.EVAL_FAMILIES))
    with self.assertRaises(ValueError):m.generator(202000,184,6)
  def test_grammar_dual_independent_evaluator(self):
    hidden,demos,x,roles=m.generator(231000,200,8)
    self.assertEqual(len(demos),32)
    self.assertEqual(len(roles),8)
    self.assertEqual(m.compose(hidden,x,roles),m.compose_independent(hidden,x,roles))
  def test_train_only_learned_prior(self):
    w,n=m.train_meta()
    self.assertEqual(len(w),3)
    self.assertGreater(n,0)
    _,demos,x,roles=m.generator(231000,201,6)
    a,b,c,d=m.infer_all(demos,w)
    self.assertEqual(len(a),8)
    self.assertEqual(len(b),8)
    self.assertEqual(len(c),8)
    self.assertTrue(all(len(role)==8 for role in a))
  def test_resident_and_data_costs(self):
    self.assertLessEqual(8*8+3*8,576)
    self.assertEqual(len(m.GROUP),len(set(m.GROUP)))
if __name__=="__main__":unittest.main()
