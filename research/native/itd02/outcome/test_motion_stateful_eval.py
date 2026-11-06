"""Source-free negative and identity tests. No MiniGrid or native processes."""
import unittest

import motion_stateful_eval as m
from ipc_source_runner import BridgeError

class StatefulEvalContractTests(unittest.TestCase):
    def test_prospective_heldout_is_disjoint(self):
        self.assertEqual(m.TEST_SEEDS,(2801,2803,2819,2833,2843,2851,2861,2879))
        self.assertEqual(m.TRAIN_SEEDS,tuple(1701+2*i for i in range(256)))
        self.assertFalse(set(m.TEST_SEEDS)&(set(m.TRAIN_SEEDS)|m.PREVIOUS_HELDOUT))
        self.assertEqual(m.TRAIN_STEPS,256)
        self.assertEqual(m.TEST_STEPS,100)

    def test_native_feedback_update_conservation(self):
        self.assertTrue(m.assert_update_delta(256,356,100,"native",2801))
        for before,after,steps in ((255,356,100),(256,355,100),(256,357,100),
                                   (256,256,0),(256,357,101),(256,257,True),
                                   (256.0,257,1)):
            with self.subTest(values=(before,after,steps)):
                with self.assertRaises(BridgeError):
                    m.assert_update_delta(before,after,steps,"native",2801)

    def test_public_start_and_episode_update_mismatch_fail_closed(self):
        def arms(other_start=False,missing_update=False):
            out={}
            for name in ("motion_native","original_native","linear_q"):
                entries=[]
                for seed in m.TEST_SEEDS:
                    prior={"x":5,"y":4,"direction":2}
                    if other_start and seed==m.TEST_SEEDS[-1] and name=="linear_q":
                        prior={"x":6,"y":4,"direction":2}
                    episode={"seed":seed,"steps":5,
                             "events":[{"before":prior,"after":prior,"action":0}]}
                    entries.append({"seed":seed,"episode":episode,
                                    "online_updates":4 if missing_update and name=="motion_native" else 5})
                out[name]={"evaluation":entries}
            return out
        self.assertTrue(m.verify_matched_eval(arms()))
        with self.assertRaisesRegex(BridgeError,"unmatched public initial"):
            m.verify_matched_eval(arms(other_start=True))
        with self.assertRaisesRegex(BridgeError,"unmatched evaluation updates"):
            m.verify_matched_eval(arms(missing_update=True))

if __name__=="__main__":
    unittest.main()
