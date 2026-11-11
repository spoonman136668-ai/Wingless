"""Source-free preregistration and test-time custody: no MiniGrid trials."""
import unittest
import progress_stateful_eval as p
import motion_stateful_eval as prior
from ipc_source_runner import BridgeError

class ProspectiveGoalProgressContracts(unittest.TestCase):
    def test_fresh_prospectively_disjoint_source_starts(self):
        self.assertEqual(p.TRAIN_SEEDS,tuple(3101+2*i for i in range(256)))
        self.assertEqual(p.TEST_SEEDS,(3701,3719,3733,3743,3761,3779,3793,3803))
        old=set(prior.TRAIN_SEEDS)|set(prior.TEST_SEEDS)|prior.PREVIOUS_HELDOUT
        self.assertFalse(set(p.TRAIN_SEEDS)&old)
        self.assertFalse(set(p.TEST_SEEDS)&(old|set(p.TRAIN_SEEDS)))
        self.assertEqual(p.TRAIN_STEPS,256)
        self.assertEqual(p.TEST_STEPS,100)

    def test_exact_online_causal_feedback_accounting(self):
        self.assertTrue(p.assert_update_delta(256,356,100,"progress_native",3701))
        for a,b,n in ((256,256,0),(256,356,99),(256,355,100),
                      (255,356,100),(256,357,100),(256,357,101),
                      (256,257,True)):
            with self.subTest(case=(a,b,n)):
                with self.assertRaises(BridgeError):
                    p.assert_update_delta(a,b,n,"progress_native",3701)

    def test_four_arm_public_start_and_feedback_guard(self):
        def make(bad_start=False,bad_updates=False):
            arms={}
            for arm in ("progress_native","motion_native","original_native","linear_q"):
                rows=[]
                for seed in p.TEST_SEEDS:
                    before={"x":10,"y":20,"direction":1}
                    if bad_start and arm=="linear_q" and seed==p.TEST_SEEDS[-1]:
                        before={"x":11,"y":20,"direction":1}
                    rows.append({"seed":seed,
                        "episode":{"seed":seed,"steps":5,
                                   "events":[{"before":before,"action":0,"after":before}]},
                        "online_updates":4 if bad_updates and arm=="progress_native" else 5})
                arms[arm]={"evaluation":rows}
            return arms
        self.assertTrue(p.verify_matched_eval(make()))
        with self.assertRaisesRegex(BridgeError,"unmatched public initial"):
            p.verify_matched_eval(make(bad_start=True))
        with self.assertRaisesRegex(BridgeError,"unmatched evaluation updates"):
            p.verify_matched_eval(make(bad_updates=True))

if __name__=="__main__":
    unittest.main()
