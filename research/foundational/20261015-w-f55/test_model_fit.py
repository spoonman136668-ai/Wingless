"""W-F55 source FIT-only preflight; diagnostic source is NOT visited."""
from __future__ import annotations
import importlib.util,pathlib,sys,unittest,copy
p=pathlib.Path(__file__).with_name("source_pilot.py")
s=importlib.util.spec_from_file_location("w55_fit_learned",p)
m=importlib.util.module_from_spec(s);sys.modules[s.name]=m;s.loader.exec_module(m)
class F55Fit(unittest.TestCase):
    def test_frozen_cohorts_and_controls(self):
        self.assertEqual((m.base.FIT.start,m.base.FIT.stop),(601000,601016))
        self.assertEqual((m.base.DIAG.start,m.base.DIAG.stop),(601100,601116))
        self.assertEqual((m.base.ORIGINAL.start,m.base.ORIGINAL.stop),(602000,602016))
        self.assertEqual(len(m.ARMS),9)
        self.assertEqual(len(set(m.ARMS)),9)
        with self.assertRaises(ValueError):m.base.source(602000,1108,8)
        with self.assertRaises(ValueError):m.base.source(601000,1104,8)
    def test_fit_label_training_is_deterministic_and_quantized(self):
        a,n=m.train();b,n2=m.train()
        self.assertEqual(a,b);self.assertEqual(n,n2)
        self.assertEqual(n,512);self.assertEqual(len(a),12)
        self.assertTrue(all(type(x) is int and -127<=x<=127 for x in a))
        self.assertNotEqual(a,(0,)*12)
    def test_public_only_composition_and_inverse_transfer(self):
        q,n=m.train()
        for seed,fam,L in ((601000,1100,6),(601002,1102,8),(601008,1101,8)):
            case=m.base.source(seed,fam,L)
            public=copy.deepcopy(case["public_calibration"])
            for arm in m.ARMS[:-1]:
                a,cost,maps,trace=m.infer(public,case["start"],case["chain"],arm,q)
                self.assertEqual(len(maps),4)
                self.assertEqual(len(trace),L)
                self.assertIn(a,range(11));self.assertGreater(cost,0)
                self.assertLessEqual(len(q)+11,64)
            reference,_,_,_=m.infer(public,case["start"],case["chain"],
                "public_full_algebraic_solver",q)
            self.assertEqual(reference,case["evaluator_truth"])
    def test_fit_only_cli_has_no_diagnostic(self):
        result=m.fit_only_report()
        self.assertEqual(result["fit_cases"],128)
        self.assertEqual(result["fit_labeled_role_examples"],512)
        self.assertTrue(result["diagnostic_cohort_NOT_EVALUATED"])
        self.assertFalse(result["source_independent_F3"])
if __name__=="__main__":unittest.main()
