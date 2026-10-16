"""W-F55 independent, source-free row/transition/score auditor.

Does not import experiment module. Checks all public affine maps from actual
calibration, every inferred role-map and every reported compositional trace.
No source materialization of protected or unobserved cohorts.
"""
import argparse,json,math
from pathlib import Path
ARMS=("learned_shared_compositional_transition","F54_frozen_transition_model",
      "fixed_equal_work_mass_beam","fixed_equal_work_entropy_beam",
      "public_full_algebraic_solver","full_exact_public_posterior",
      "public_parity_rule","public_support_nearest","private_truth_oracle_unfair")
SLOPES=(8,9,10,1,2,3);F54_CODES=(2,1,1,2,1,2,1,2)
def verify(doc,pr):
    assert pr["schema"]=="wingless.f55.public-constructive-binding-efficiency.prereg.v1"
    assert doc["schema"]=="wingless.f55.learned-public-role-bindings.trainonly.v1"
    assert doc["stage"]=="TRAIN_ONLY_NEW_PUBLIC_IDENTIFIABLE_ROLE_TRANSFER"
    assert pr["controls"]==list(ARMS)
    assert doc["fit_seeds"]==pr["fit"]["seeds"]
    assert doc["fit_families"]==pr["fit"]["families"]
    assert doc["diagnostic_seeds"]==pr["diagnostic"]["seeds"]
    assert doc["diagnostic_families"]==pr["diagnostic"]["families"]
    assert doc["original_seeds_NOT_RESERVED"]==pr["future_original_NOT_RESERVED"]["seeds"]
    assert doc["learned_model_bytes"]==12 and doc["recurrent_state_payload_bytes"]==11
    assert doc["learned_model_bytes"]+doc["recurrent_state_payload_bytes"]<=64
    q=doc["learned_int8_parameters"]
    assert len(q)==12 and all(type(v)==int and -127<=v<=127 for v in q)
    assert doc["fit_labeled_role_examples"]==512 and doc["cases"]==128
    assert len(doc["rows"])==128
    assert doc["physical_native_576B_unverified"] is True
    assert doc["F54_frozen_control_ADAPTED_NOT_IDENTICAL"] is True
    for name in ("independent_F3","primary_original_executed","accepted_ref_mutation",
                 "live_optimizer_promotion"):
        assert doc[name] is False
    assert doc["paid_api_calls"]==0
    counts={a:0 for a in ARMS}
    family={str(f):{a:0 for a in ARMS} for f in doc["diagnostic_families"]}
    works={a:0 for a in ARMS[:-1]}
    seen=set()
    for row in doc["rows"]:
        seed,fam,L=row["seed"],row["family"],row["length"]
        assert (seed,fam,L) not in seen;seen.add((seed,fam,L))
        assert pr["diagnostic"]["seeds"][0]<=seed<pr["diagnostic"]["seeds"][1]
        assert fam in pr["diagnostic"]["families"] and L in (8,10)
        assert 0<=row["start"]<11 and len(row["chain"])==L
        assert len(row["public_calibration"])==4
        assert set(row["answers"])==set(ARMS)
        assert set(row["abstract_work"])==set(ARMS[:-1])
        assert set(row["public_role_map_receipts"])==set(ARMS[:-1])
        assert set(row["trace"])==set(ARMS[:-1])
        observations={}
        true={}
        for entry in row["public_calibration"]:
            role=entry["role"];(x,y),(z,t)=entry["pairs"]
            assert 0<=role<4 and role not in observations and (z-x)%11!=0
            slope=((t-y)*pow((z-x)%11,-1,11))%11
            shift=(y-slope*x)%11
            assert slope in SLOPES
            observations[role]=entry;true[role]=(slope,shift)
        assert set(observations)==set(range(4))
        def run(maps):
            state=row["start"];trace=[]
            for role,direction in row["chain"]:
                assert role in range(4) and direction in (-1,1)
                a,b=maps[role]
                state=((a*state+b) if direction==1 else
                       ((state-b)*pow(a,-1,11)))%11
                trace.append(state)
            return state,trace
        truth,_=run(true)
        assert row["evaluator_truth"]==truth
        assert row["answers"]["private_truth_oracle_unfair"]==truth
        for arm in ARMS[:-1]:
            mm=row["public_role_map_receipts"][arm]
            assert len(mm)==4
            maps={}
            for role,(a,b) in enumerate(mm):
                assert a in SLOPES and type(b)==int and 0<=b<11
                x,y=observations[role]["pairs"][0]
                assert (a*x+b)%11==y
                maps[role]=(a,b)
                if arm in ("public_full_algebraic_solver","full_exact_public_posterior"):
                    assert (a,b)==true[role]
                if arm=="learned_shared_compositional_transition":
                    (x,y),(z,t)=observations[role]["pairs"]
                    dx=(z-x)%11;dy=(t-y)%11
                    feat=[int(s*dx%11==dy) for s in SLOPES]
                    sc=[(q[k+6]+q[k]*feat[k])/16 for k in range(6)]
                    expected=SLOPES[max(range(6),key=lambda k:(sc[k],-k))]
                    assert a==expected
            result,trace=run(maps)
            assert row["answers"][arm]==result
            assert row["trace"][arm]==trace
            correct=int(result==truth);counts[arm]+=correct;family[str(fam)][arm]+=correct
            work=row["abstract_work"][arm]
            expect=(4*5+L) if arm=="public_full_algebraic_solver" else (
                 4*6*11*2+L if arm=="full_exact_public_posterior" else 4*6*2+L)
            assert type(work)==int and work==expect
            works[arm]+=work
        counts[ARMS[-1]]+=1;family[str(fam)][ARMS[-1]]+=1
    assert len(seen)==128 and doc["by_family"]==family
    for arm in ARMS:
        metric=doc["scores"][arm]
        assert metric["correct"]==counts[arm] and metric["total"]==128
        assert math.isclose(metric["accuracy"],counts[arm]/128,abs_tol=1e-12)
    fair=ARMS[1:5]+ARMS[6:8]
    best=max(fair,key=lambda a:counts[a])
    gain=(counts[ARMS[0]]-counts[best])/128
    loss=(counts["full_exact_public_posterior"]-counts[ARMS[0]])/128
    ratio=works[ARMS[0]]/works["full_exact_public_posterior"]
    noregress=all(family[str(f)][ARMS[0]]>=max(family[str(f)][a] for a in fair)
                  for f in pr["diagnostic"]["families"])
    witness=counts["public_full_algebraic_solver"]/128
    passed=(witness>=.85 and gain>=.05 and loss<=.02 and ratio<=.5 and noregress)
    assert doc["strongest_fair_control"]==best
    assert math.isclose(doc["candidate_gain"],gain,abs_tol=1e-12)
    assert math.isclose(doc["loss_vs_exact"],loss,abs_tol=1e-12)
    assert math.isclose(doc["work_ratio_vs_exact"],ratio,abs_tol=1e-12)
    assert math.isclose(doc["public_constructive_witness_accuracy"],witness,abs_tol=1e-12)
    assert doc["no_family_regression"]==noregress and doc["frozen_gate_met"]==passed
    assert doc["verdict"]==("TRAIN_ONLY_LEARNED_ROLE_BINDING_SUPPORTED" if passed
                            else "TRAIN_ONLY_LEARNED_ROLE_BINDING_NEGATIVE")
    return dict(schema="wingless.f55.independent-train-recount.v1",
        rows=len(seen),candidate_correct=counts[ARMS[0]],
        strongest_fair=best,strongest_correct=counts[best],
        public_constructive_correct=counts["public_full_algebraic_solver"],
        gain_pp=round(100*gain,6),work_ratio=ratio,nondegradation=noregress,
        train_gate_passed=passed,
        no_original_or_F3=True)
if __name__=="__main__":
    parser=argparse.ArgumentParser()
    parser.add_argument("result");parser.add_argument("prereg")
    args=parser.parse_args()
    print(json.dumps(verify(json.loads(Path(args.result).read_text()),
            json.loads(Path(args.prereg).read_text())),sort_keys=True))
