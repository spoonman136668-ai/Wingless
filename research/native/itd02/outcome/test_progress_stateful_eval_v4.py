"""Prospective V4 five-arm source protocol, no MiniGrid/native source execution."""
import unittest
import progress_stateful_eval_v4 as current
import progress_stateful_eval_v3 as predecessor
from ipc_source_runner import BridgeError

class StagnationV4SourceProtocol(unittest.TestCase):
    def test_fresh_disjoint_training_and_five_arm_test_starts(self):
        self.assertEqual(current.TRAIN_SEEDS,tuple(7101+2*i for i in range(256)))
        self.assertEqual(current.TEST_SEEDS,(8101,8111,8123,8137,8147,8161,8171,8191))
        previous=set(predecessor.TRAIN_SEEDS)|set(predecessor.TEST_SEEDS)|set(current.PREVIOUS_HELDOUT)
        self.assertFalse(set(current.TRAIN_SEEDS)&previous)
        self.assertFalse(set(current.TEST_SEEDS)&(previous|set(current.TRAIN_SEEDS)))
        self.assertEqual(current.TRAIN_STEPS,256)
        self.assertEqual(current.TEST_STEPS,100)

    def test_native_feedback_exact_count_rejects_gaps(self):
        self.assertTrue(current.assert_update_delta(256,356,100,"progress_native",8101))
        for before,after,steps in ((255,356,100),(256,355,100),(256,357,100),
                                   (256,256,0),(256,357,101),(256,257,True)):
            with self.subTest(values=(before,after,steps)):
                with self.assertRaises(BridgeError):
                    current.assert_update_delta(before,after,steps,"progress_native",8101)

    def test_five_arm_equal_public_starts_and_feedback(self):
        names=("progress_native","progress_v3","motion_native","original_native","linear_q")
        def fixture(bad_start=False,bad_count=False):
            result={}
            for name in names:
                eval_rows=[]
                for seed in current.TEST_SEEDS:
                    before={"x":5,"y":7,"direction":1,"goal_x":11,"goal_y":18,"front_blocked":False}
                    if bad_start and name=="progress_v3" and seed==current.TEST_SEEDS[-1]:
                        before={**before,"x":6}
                    ep={"seed":seed,"steps":5,"events":[{"before":before,"after":before,"action":0}]}
                    count=4 if bad_count and name=="progress_native" else 5
                    eval_rows.append({"seed":seed,"episode":ep,"online_updates":count})
                result[name]={"evaluation":eval_rows}
            return result
        self.assertTrue(current.verify_matched_eval(fixture()))
        with self.assertRaisesRegex(BridgeError,"unmatched public initial"):
            current.verify_matched_eval(fixture(bad_start=True))
        with self.assertRaisesRegex(BridgeError,"unmatched evaluation updates"):
            current.verify_matched_eval(fixture(bad_count=True))

if __name__=="__main__":
    unittest.main()
