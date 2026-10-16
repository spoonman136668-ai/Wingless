"""W-F57 source-independent small FIT predictor witnesses."""
import importlib.util
import pathlib
import sys
import unittest
HERE=pathlib.Path(__file__).resolve().parent
def load(name,file):
    x=importlib.util.spec_from_file_location(name,HERE/file)
    y=importlib.util.module_from_spec(x)
    sys.modules[name]=y
    x.loader.exec_module(y)
    return y
a=load("audit_w57","audit_fit.py")
s=load("source_w57","source_fit.py")
class AuditCheck(unittest.TestCase):
    def test_distinct_control_and_belief_predictions(self):
        case=s.source(661001,1241,10)
        row={"seed":case["seed"],"family":case["family"],"length":case["length"],
            "public":[{"role":r,"pairs":[list(v) for v in pairs]} for r,pairs in case["public"]],
            "start":case["start"],"chain":[list(x) for x in case["chain"]]}
        x,w=a.answer(row,(1.,1.,1.,1.))
        self.assertEqual(x["learned_public_belief_state"],s.factorized_transition(case,(1.,1.,1.,1.))[0])
        self.assertEqual(x["frozen_F56_posterior_beam"],s.frozen_F56_budgeted_beam(case)[0])
        self.assertEqual(x["known_operator_MAP"],s.fixed_known_operator_bank(case)[0])
        self.assertTrue(all(w[k]>0 for k in x))
    def test_missing_public_role_blocks(self):
        row={"family":1241,"public":[{"role":0,"pairs":[[0,1],[1,2]]}]}
        with self.assertRaises(AssertionError):
            a.posterior(row)
if __name__=="__main__":unittest.main()
