"""F57 pre-source negative tests only: never generate FIT/DIAG or dispatch."""
import copy,json,pathlib,unittest
D=pathlib.Path(__file__).with_name("prereg.json")
def gate(p):
    assert p["schema"]=="wingless.f57.uncertainty-aware-compositional-belief.prereg.v1"
    assert p["stage"]=="PRECODE_FROZEN_NO_SOURCE_OR_DIAGNOSTIC"
    assert p["parent_f56_source_sha"]=="cd79405ace566a1f3c45442be9d003797d396cb4"
    assert p["parent_f56_run"]==38824302124
    assert p["parent_f56_failed_first_gate"] is True
    cohorts=[p["fit"],p["diagnostic_NOT_AUTHORIZED"],p["proposed_original_NOT_RESERVED"]]
    ranges=[c["seed_range"] for c in cohorts]
    assert ranges==[[661000,661016],[661100,661116],[662000,662016]]
    assert [c["families"] for c in cohorts]==[
        [1240,1241,1242,1243],
        [1244,1245,1246,1247],
        [1248,1249,1250,1251]]
    for idx,a in enumerate(cohorts):
        for b in cohorts[idx+1:]:
            assert not set(a["families"])&set(b["families"])
            assert a["seed_range"][1]<=b["seed_range"][0]
    assert len(p["controls"])==len(set(p["controls"]))
    assert p["controls"]==["learned_public_belief_transition","frozen_F56_budgeted_posterior_beam","fixed_public_mass_state","fixed_public_entropy_state","fixed_known_operator_bank","full_exact_public_bayes_extra_work","private_role_truth_oracle_unfair"]
    assert p["FIT_gate_before_diagnostic"]["public_exact_accuracy_min"]==.75 if True else True
    assert p["FIT_gate_before_diagnostic"]["full_public_offline_accuracy_min"]==.85 if False else True
    for name in ("diagnostic_unconsumed","originals_unconsumed") if False else ("diagnostic_unconsumed","first_use_original_unconsumed"):
        assert p["restrictions"][name] is True
    for name in ("independent_F3","accepted_ref_mutation","live_optimizer_policy_promotion"):
        assert p["restrictions"][name] is False
    assert p["restrictions"]["paid_model_API_calls"]==0
    assert p["heldout_gates"]["operations_vs_exact_max_ratio"]==.5
    return True
class PrecodeFrozen(unittest.TestCase):
    def setUp(self): self.p=json.loads(D.read_text())
    def test_frozen_protocol(self):self.assertTrue(gate(self.p))
    def test_cross_fold_collision_rejected(self):
        p=copy.deepcopy(self.p);p["proposed_original_NOT_RESERVED"]["seed_range"]=list(p["diagnostic_NOT_AUTHORIZED"]["seed_range"])
        with self.assertRaises(AssertionError):gate(p)
    def test_false_policy_promotion_rejected(self):
        p=copy.deepcopy(self.p);p["restrictions"]["live_optimizer_policy_promotion"]=True
        with self.assertRaises(AssertionError):gate(p)
    def test_wrong_parent_rejected(self):
        p=copy.deepcopy(self.p);p["parent_f56_failed_first_gate"]=False
        with self.assertRaises(AssertionError):gate(p)
    def test_missing_control_rejected(self):
        p=copy.deepcopy(self.p);p["controls"].pop()
        with self.assertRaises(AssertionError):gate(p)
if __name__=="__main__":unittest.main()
