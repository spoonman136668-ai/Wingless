"""Source-free boundary tests; never open MiniGrid or launch a native model."""
import unittest
from ipc_source_runner import BridgeError
from exact_exposure_v2 import exact_train_v2, TRAIN_STEPS

class ExactExposureV2Tests(unittest.TestCase):
    def test_budget_cut_is_allowed_only_at_final_boundary(self):
        calls=[]
        def fake(_, policy, seed, limit):
            calls.append((seed,limit))
            return {"steps":limit, "completed":False}
        rows=exact_train_v2(object(),fake)
        self.assertEqual(sum(r["steps"] for r in rows),TRAIN_STEPS)
        self.assertEqual(len(calls),1)
        self.assertEqual(calls[0][1],TRAIN_STEPS)

    def test_completed_episode_allows_next_seed(self):
        calls=[]
        def fake(_,policy,seed,limit):
            calls.append(limit)
            if len(calls)==1:
                return {"steps":60,"completed":True}
            return {"steps":limit,"completed":False}
        rows=exact_train_v2(object(),fake)
        self.assertEqual([r["steps"] for r in rows],[60,196])
        self.assertEqual(calls,[256,196])

    def test_early_incomplete_episode_never_resets(self):
        calls=[]
        def fake(_,policy,seed,limit):
            calls.append(limit)
            return {"steps":20,"completed":False}
        with self.assertRaisesRegex(BridgeError,"no source reset"):
            exact_train_v2(object(),fake)
        self.assertEqual(calls,[256])

    def test_invalid_step_counts_fail_closed(self):
        for invalid in (0,-1,257,True,1.5):
            def fake(_,policy,seed,limit):
                return {"steps":invalid,"completed":True}
            with self.subTest(invalid=invalid):
                with self.assertRaises(BridgeError):
                    exact_train_v2(object(),fake)

if __name__=="__main__":
    unittest.main()
