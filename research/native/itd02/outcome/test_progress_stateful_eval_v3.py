"""V2 successor preregistration only. Never execute MiniGrid in this test."""
import unittest
import progress_stateful_eval_v3 as current
import progress_stateful_eval as failed_v1
import progress_stateful_eval_v2 as failed_v2
from ipc_source_runner import BridgeError

class TurnEvidenceSuccessorV3Protocol(unittest.TestCase):
 def test_predeclared_new_cohorts_never_reuse_observed_v1(self):
  self.assertEqual(current.TRAIN_SEEDS,tuple(5101+2*i for i in range(256)))
  self.assertEqual(current.TEST_SEEDS,(6101,6113,6121,6133,6151,6163,6173,6197))
  previous=set(failed_v1.TRAIN_SEEDS)|set(failed_v1.TEST_SEEDS)|set(failed_v1.PREVIOUS_HELDOUT)|set(failed_v2.TRAIN_SEEDS)|set(failed_v2.TEST_SEEDS)
  self.assertFalse(set(current.TRAIN_SEEDS)&previous)
  self.assertFalse(set(current.TEST_SEEDS)&(previous|set(current.TRAIN_SEEDS)))
  self.assertEqual(current.TRAIN_STEPS,256)
  self.assertEqual(current.TEST_STEPS,100)
 def test_fail_closed_update_count(self):
  self.assertTrue(current.assert_update_delta(256,356,100,"progress_native",6101))
  for a,b,n in ((256,355,100),(255,356,100),(256,257,0),(256,257,False)):
   with self.subTest(c=(a,b,n)):
    with self.assertRaises(BridgeError):
     current.assert_update_delta(a,b,n,"progress_native",6101)
if __name__=="__main__":
 unittest.main()
