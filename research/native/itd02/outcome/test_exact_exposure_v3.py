"""V3 preregistration checks that require no source-world execution."""
import unittest
import exact_exposure_gate as frozen
from exact_exposure_v3 import TRAIN_SEEDS_V3
from exact_exposure_v2 import exact_train_v2

class ProspectiveV3Tests(unittest.TestCase):
    def test_disjoint_immutable_start_schedule(self):
        self.assertEqual(TRAIN_SEEDS_V3,tuple(1009+2*i for i in range(256)))
        self.assertEqual(len(set(TRAIN_SEEDS_V3)),256)
        self.assertFalse(set(TRAIN_SEEDS_V3).intersection(frozen.TEST_SEEDS))
        self.assertEqual(frozen.TRAIN_STEPS,256)
    def test_exact_accounting_without_fabricated_steps(self):
        import exact_exposure_v2 as v2
        prior=v2.TRAIN_SEEDS
        v2.TRAIN_SEEDS=TRAIN_SEEDS_V3
        seen=[]
        try:
            def fake(env,policy,seed,remaining):
                seen.append(seed)
                return {"steps":1,"completed":True}
            rows=exact_train_v2(None,fake)
            self.assertEqual(len(rows),256)
            self.assertEqual(sum(x["steps"] for x in rows),256)
            self.assertEqual(tuple(seen),TRAIN_SEEDS_V3)
        finally:
            v2.TRAIN_SEEDS=prior

if __name__=="__main__":
    unittest.main()
