"""W-F54 source-free row-level evaluation audit against frozen F54 prereg.

This module never imports the experiment's code, train/test source generator,
or private model objects, and cannot dispatch any experiment. Original first-use
and independent F3 authority stay false.
"""
import json,math,sys
from pathlib import Path

def verify(data,pr):
    assert pr["schema"]=="wingless.f54.learned-relation-transition-invariants.prereg.v1"
    assert data["schema"]=="wingless.f54.learned-role-transition.trainonly.v1"
    assert data["stage"]=="TRAIN_ONLY_DISJOINT_FIT_DIAG_LEARNED_TRANSITIONS"
    assert data["fit_seeds"]==pr["fit"]["seeds"]
    assert data["fit_families"]==pr["fit"]["families"]
    assert data["diagnostic_seeds"]==pr["diagnostic"]["seeds"]
    assert data["diagnostic_families"]==pr["diagnostic"]["families"]
    assert data["original_seeds_NOT_RESERVED"]==pr["future_original_NOT_RESERVED"]["seeds"]
    assert data["independent_F3"] is False and data["original_executed"] is False
    assert data["accepted_ref_mutation"] is False and data["paid_api_calls"]==0
    assert data["physical_native_576B_unverified"] is True
    arms=pr["controls"]
    assert len(arms)==len(set(arms)) and len(arms)==9
    assert set(data["scores"])==set(arms)
    assert data["transition_parameter_bytes"]==8
    assert len(data["learned_transition_temperature_codes"])==8
    assert all(type(x)==int and 0<=x<=5 for x in data["learned_transition_temperature_codes"])
    assert data["train_examples"]==128
    assert len(data["F53_frozen_public_fitted_weights"])==8
    assert data["recurrent_belief_bytes"]==7
    assert data["selected_state_index_bytes"]==2
    assert data["cases"]==128 and len(data["rows"])==128
    seen=set();hits={a:0 for a in arms};by={str(f):{a:0 for a in arms} for f in pr["diagnostic"]["families"]}
    works={a:0 for a in arms if a!="private_truth_oracle_unfair"}
    candidate=arms[0];oracle="private_truth_oracle_unfair"
    for row in data["rows"]:
        identity=(row["seed"],row["family"],row["length"])
        assert identity not in seen;seen.add(identity)
        assert pr["diagnostic"]["seeds"][0]<=row["seed"]<pr["diagnostic"]["seeds"][1]
        assert row["family"] in pr["diagnostic"]["families"]
        assert row["length"] in pr["diagnostic"]["path_lengths"]
        assert len(row["chain"])==row["length"] and len(row["public_observations"])==4
        assert len(row["public_calibration"])==4 and len(row["public_context"])==3
        assert set(row["answers"])==set(arms)
        assert set(row["abstract_operations"])==set(works)
        assert set(row["state_payload_bytes"])==set(works)
        assert row["answers"][oracle]==row["truth"]
        assert type(row["truth"])==int and 0<=row["truth"]<=6
        for arm in arms:
            assert type(row["answers"][arm])==int and 0<=row["answers"][arm]<=6
            hit=int(row["answers"][arm]==row["truth"])
            hits[arm]+=hit;by[str(row["family"])][arm]+=hit
            if arm!=oracle:
                ops=row["abstract_operations"][arm];memory=row["state_payload_bytes"][arm]
                assert type(ops)==int and ops>0
                assert type(memory)==int and memory>0 and memory<=14
                works[arm]+=ops
        assert row["state_payload_bytes"][candidate]+8+2<=pr["resource_constraints"]["total_model_parameter_and_state_payload_bytes_max"]
    assert len(seen)==128 and by==data["by_family"]
    assert works==data["posterior_abstract_operations"]
    for arm,metric in data["scores"].items():
        assert metric["correct"]==hits[arm] and metric["total"]==128
        assert math.isclose(metric["accuracy"],hits[arm]/128,abs_tol=1e-12)
    fair=arms[1:7]
    # Do NOT compare equal-work candidates to the extra-work full exact reference.
    strongest=max(fair,key=lambda a:hits[a])
    gain=(hits[candidate]-hits[strongest])/128
    ratio=works[candidate]/works["exact_full_public_bayes"]
    no_family_regression=all(by[str(f)][candidate]>=max(by[str(f)][a] for a in fair)
        for f in pr["diagnostic"]["families"])
    passed=(gain>=pr["thresholds"]["gain_vs_strongest_equal_work_fair_min"] and
       ratio<=pr["thresholds"]["posterior_operation_ratio_vs_exact_max"] and
       no_family_regression)
    assert data["strongest_fair_control"]==strongest
    assert math.isclose(data["candidate_gain"],gain,abs_tol=1e-12)
    assert math.isclose(data["posterior_work_ratio_vs_exact"],ratio,abs_tol=1e-12)
    assert data["per_family_nondegradation"]==no_family_regression
    assert data["frozen_gate_met"]==passed
    assert data["verdict"]==("TRAIN_ONLY_TRANSITION_GENERALIZATION_SUPPORTED" if passed
        else "TRAIN_ONLY_TRANSITION_GENERALIZATION_NEGATIVE")
    return dict(schema="wingless.f54.independent-train-metric-recount.v1",rows=len(seen),
        candidate_correct=hits[candidate],strongest_fair=strongest,
        strongest_correct=hits[strongest],gain_percentage_points=round(gain*100,6),
        abstract_work_ratio_to_exact=round(ratio,6),
        per_family_nondegradation=no_family_regression,
        source_frozen_train_gate=passed,
        train_only=True,original_executed=False,independent_F3=False)

if __name__=="__main__":
    assert len(sys.argv)==3
    doc=json.loads(Path(sys.argv[1]).read_text())
    pre=json.loads(Path(sys.argv[2]).read_text())
    print(json.dumps(verify(doc,pre),sort_keys=True))
