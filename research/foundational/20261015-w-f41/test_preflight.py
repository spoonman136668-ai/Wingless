"""Source preflight; no prospective 72-question TRAIN-only diagnostic."""
import importlib.util,pathlib,unittest
P=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("wf41source",P)
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
class Gate(unittest.TestCase):
  def test_reserved_cannot_be_generated(self):
    self.assertFalse(set(m.SEEDS)&set(m.RESERVED))
    with self.assertRaises(ValueError):m.source(312000,354,3)
  def test_truth_and_posteriors(self):
    private,public,start,path=m.source(311000,350,3)
    self.assertEqual(m.follow(private,start,path),m.independent_follow(private,start,path))
    self.assertEqual(len(public),24)
    a,base,H,bank,answers=m.score_actions(public,start,path,12,17)
    self.assertEqual(len(bank),12)
    self.assertTrue(all(len(a[q])==3 for q in a))
  def test_private_port_rejects_third_and_duplicate(self):
    secret,public,start,path=m.source(311000,350,5)
    port=m.Port(secret,public)
    found=next((r,x) for r in range(8) for x in range(8) if (r,x) not in port.known)
    port.buy(*found)
    with self.assertRaises(ValueError):port.buy(*found)
    another=next((r,x) for r in range(8) for x in range(8) if (r,x) not in port.known)
    port.buy(*another)
    with self.assertRaises(ValueError):port.buy(7,7)
  def test_entropy_not_identical_to_decision_value(self):
    private,public,start,path=m.source(311001,351,3)
    scores,*rest=m.score_actions(public,start,path,48,99)
    self.assertTrue(all(score[0]>=-1e-12 and score[1]>=-1e-12 for score in scores.values()))
    self.assertEqual(m.SAMPLE_BUDGETS,(12,48))
if __name__=="__main__":unittest.main()
