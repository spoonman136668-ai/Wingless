"""F51 PRECODE-only: source identities and public baselines, no experiment."""
import copy
import json
import pathlib
import unittest

P=pathlib.Path(__file__).with_name("prereg-candidate.json")
def validate(p):
    assert p["schema"]=="wingless.f51.operator-generalization-strong-public-controls.prereg.v1"
    assert p["stage"]=="PRECODE_FROZEN_CANDIDATE_NOT_OWNER_ADMITTED"
    assert p["parent_f50_source_sha"]=="646b9db1d48eb03150ad79f1c88b7a44354ab96f"
    assert p["independent_f50_review"]["observed_learned_correct"]==112
    assert p["independent_f50_review"]["observed_three_bit_parity_correct"]==112
    required={"three_bit_public_parity","four_bit_public_parity","FIT_context_truth_table","FIT_public_logistic","private_operator_oracle_unfair","learned_context_operator_posterior"}
    assert required.issubset(p["controls"])
    assert len(p["controls"])==len(set(p["controls"]))
    assert p["evaluator_only_controls"]==["private_operator_oracle_unfair"]
    assert p["thresholds"]["minimum_gain_vs_strongest_fair"]==0.05
    ranges=[(p[k]["seeds_start_inclusive"],p[k]["seeds_end_exclusive"])
            for k in ("fresh_fit","fresh_diagnostic","proposed_original_NOT_RESERVED")]
    ranges += [tuple(r) for r in p["excluded_prior_originals"]]
    assert all(type(a) is int and type(b) is int and a<b for a,b in ranges)
    assert all(a[1]<=b[0] or b[1]<=a[0]
               for i,a in enumerate(ranges) for b in ranges[i+1:])
    families=[set(p[k]["families"]) for k in ("fresh_fit","fresh_diagnostic","proposed_original_NOT_RESERVED")]
    assert all(len(x)==4 for x in families)
    assert not any(a&b for i,a in enumerate(families) for b in families[i+1:])
    for flag in ("source_code_implemented","train_or_diagnostic_executed","owner_admitted",
                 "primary_original_executed","independent_F3","accepted_ref_mutation",
                 "dreamrsi_gain_proven","thetaevolve_training_authorized"):
        assert p[flag] is False
    assert p["paid_api_calls"]==0
    return True

class FrozenF51(unittest.TestCase):
    def setUp(self):
        self.p=json.loads(P.read_text())
    def test_precode_contract(self):self.assertTrue(validate(self.p))
    def test_missing_strong_public_control_rejected(self):
        p=copy.deepcopy(self.p);p["controls"].remove("three_bit_public_parity")
        with self.assertRaises(AssertionError):validate(p)
    def test_original_collision_rejected(self):
        p=copy.deepcopy(self.p);p["proposed_original_NOT_RESERVED"]["seeds_start_inclusive"]=521100
        with self.assertRaises(AssertionError):validate(p)
    def test_claimed_promotion_rejected(self):
        p=copy.deepcopy(self.p);p["owner_admitted"]=True
        with self.assertRaises(AssertionError):validate(p)
if __name__=="__main__":unittest.main()
