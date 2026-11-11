"""No actual source environment: prospective split and accounting contracts."""
import unittest
import motion_source_compare as m
import exact_exposure_gate as v1
import exact_exposure_v3 as v3

class TestMotionComparisonSourceFree(unittest.TestCase):
 def test_disjoint_preregistered_splits(self):
  self.assertEqual(m.TRAIN_SEEDS,tuple(1701+2*i for i in range(256)))
  self.assertEqual(m.TEST_SEEDS,(2609,2617,2621,2633,2647,2657,2663,2671))
  prior=set(v3.TRAIN_SEEDS_V3)|set(v1.TEST_SEEDS)
  self.assertFalse(set(m.TRAIN_SEEDS)&prior)
  self.assertFalse(set(m.TEST_SEEDS)&(set(m.TRAIN_SEEDS)|prior))
  self.assertEqual(m.TRAIN_STEPS,256)
  self.assertEqual(m.TEST_STEPS,100)

 def test_only_actual_position_change_counts(self):
  origin={"x":1,"y":2}
  moved={"x":2,"y":2}
  events=[{"action":0,"before":origin,"after":origin},
          {"action":2,"before":origin,"after":moved}]
  row={"events":events,"return":0.0,"goal_reached":False}
  d=m.diagnostics([row])
  self.assertEqual(d["position_advances"],1)
  self.assertEqual(d["selected_action_counts"],{"0":1,"1":0,"2":1})
  self.assertEqual(d["authentic_steps"],2)
  self.assertEqual(d["goal_count"],0)

if __name__=="__main__":unittest.main()
