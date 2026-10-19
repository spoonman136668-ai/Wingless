"""T16 Wingless mathematically deterministic checks; no original cohort birth."""
import importlib.util
import sys
import unittest
from pathlib import Path

p=Path(__file__).with_name("20261019-wingless-t16b-persistent-bindings.py")
spec=importlib.util.spec_from_file_location("wing_t16b",p)
m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m)

class PersistentBindingTests(unittest.TestCase):
    def test_unseen_role_with_inverse_cancels_exactly(self):
        rows=(((2,5),),)
        self.assertEqual(m.shared_role_certificate(0,rows,[(0,False),(0,True)],7,budget=1)[0],7)

    def test_repeat_binding_affine_not_independent_approximation(self):
        for slope in range(1,11):
            op=(slope,5)
            rows=(((3,(slope*3+5)%11),),)
            result,receipt=m.shared_role_certificate(0,rows,[(0,False),(0,True),
                (0,False),(0,True)],7,budget=1)
            self.assertEqual(result,7)
            self.assertTrue(receipt["certified"])

    def test_affine_single_role_exact_quantifier(self):
        rows=(((2,7),),)
        word=[(0,False),(0,False),(0,True)]
        got,diag=m.shared_role_certificate(0,rows,word,2,32)
        self.assertEqual(got,7)
        self.assertEqual(diag["assignments"],10)

    def test_s3_repeated_primitive_binding(self):
        rows=(((0,1),),)
        result,info=m.shared_role_certificate(1,rows,[(0,True),(0,False),
                            (0,True),(0,False)],2,1)
        self.assertEqual(result,2)
        self.assertTrue(info["certified"])

    def test_ambiguity_never_certified(self):
        rows=(((0,1),),)
        got,info=m.shared_role_certificate(1,rows,[(0,False)],2,budget=2)
        self.assertIsNone(got)
        self.assertEqual(info["reason"],"MULTIPLE_VALID_OUTCOMES")

    def test_budget_exhaustion_never_claims_answer(self):
        rows=(((2,7),),((3,9),))
        got,info=m.shared_role_certificate(0,rows,[(0,False),(1,False)],3,1)
        self.assertIsNone(got)
        self.assertEqual(info["reason"],"BOUNDED_UNKNOWN")

    def test_full_public_exact_can_fallback(self):
        a=(2,3);b=(4,5)
        rows=(((0,3),(1,5)),((0,5),(1,9)))
        word=[(0,False),(1,False),(0,True)]
        v,receipt=m.bound_query_solve(0,rows,word,2,(0,0,0,0),cap=1)
        exact=m.public_exact_fallback(0,rows,word,2)
        self.assertEqual(v,exact)
        self.assertTrue(receipt["sound"])
        self.assertLessEqual(receipt["queries"],2)

    def test_inconsistent_public_data_cannot_be_certified(self):
        got,info=m.shared_role_certificate(0,(((0,1),(0,2)),),[(0,False)],3)
        self.assertIsNone(got)
        self.assertEqual(info["reason"],"INCONSISTENT_PUBLIC_EVIDENCE")

    def test_public_only_training_scores_are_stable(self):
        training=[(0,(((0,1),(1,2)),((0,2),(1,3))),[(0,False),(1,False)],2)]
        a=m.score_role_requests(training)
        self.assertEqual(a,m.score_role_requests(training))
        self.assertEqual(len(a),4)

if __name__=="__main__":
    unittest.main()
