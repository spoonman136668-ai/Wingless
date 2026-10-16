"""W-F55 prospective learned public-role operator inference and variable binding.

Cohorts/controls frozen in prereg-candidate.json, blob
47515bea8590e3a5c849ce058587fa1852d991a7.
FIT-only supervision learns class-specific coefficients for predicting a role's
affine slope from two public input/output calibration observations. No hidden
operator or evaluator truth crosses the diagnostic action boundary. The
recurrent consumer composes recovered affine maps for unseen role order and
inverse binding. The exact public algebraic control is honestly costed and may
make the frozen +5pp gain impossible. Only TRAIN, not F3 or native residency.
"""
from __future__ import annotations
import importlib.util,json,math,sys
from pathlib import Path
HERE=Path(__file__).resolve().parent
s=importlib.util.spec_from_file_location("w55_public_fit_feasibility",HERE/"source_feasibility.py")
base=importlib.util.module_from_spec(s);sys.modules[s.name]=base;s.loader.exec_module(base)
ARMS=(
 "learned_shared_compositional_transition","F54_frozen_transition_model",
 "fixed_equal_work_mass_beam","fixed_equal_work_entropy_beam",
 "public_full_algebraic_solver","full_exact_public_posterior",
 "public_parity_rule","public_support_nearest","private_truth_oracle_unfair")
SLOPES=tuple(x%base.PRIME for x in base.SLOPES)
# Original F54 source-free result witness: 8 role/direction quantized
# temperatures. An explicitly adapted historical control, NOT an exact
# replay of the F54 seven-state scientific system on a changed F55 source.
F54_FROZEN_TEMPERATURE_CODES=(2,1,1,2,1,2,1,2)
def features(o):
    (x,y),(z,t)=o["pairs"]
    dx=(z-x)%base.PRIME;dy=(t-y)%base.PRIME
    assert dx!=0
    return tuple(int((a*dx)%base.PRIME==dy) for a in SLOPES)
def train():
    # Cross-entropy on FIT-only true operator class, not intermediate
    # diagnostic labels. Six supervised per-class logit slopes plus six
    # class biases: 12 quantized signed bytes, frozen before diagnosis.
    obs=[]
    for seed in base.FIT:
        for fam in base.FIT_FAM:
            for length in (6,8):
                case=base.source(seed,fam,length)
                for entry in case["public_calibration"]:
                    target=case["evaluator_private_maps"][entry["role"]][0]%base.PRIME
                    obs.append((features(entry),SLOPES.index(target)))
    w=[0.]*6;b=[0.]*6
    for epoch in range(6):
        for feat,y in obs:
            logit=[b[k]+w[k]*feat[k] for k in range(6)]
            top=max(logit);exp=[math.exp(z-top) for z in logit]
            probs=[z/sum(exp) for z in exp]
            for k in range(6):
                delta=(probs[k]-int(k==y))
                b[k]-=.035*(delta+.00015*b[k])
                w[k]-=.035*(delta*feat[k]+.00015*w[k])
    # int8 quantization 1/16: report true compressed *payload* bytes.
    encoded=tuple(max(-127,min(127,round(z*16))) for z in w+b)
    need=len(obs)
    return encoded,need
def select(o,arm,quant):
    evidence=features(o)
    if arm=="public_full_algebraic_solver" or arm=="full_exact_public_posterior":
        (x,y),(z,t)=o["pairs"]
        a=((t-y)*pow((z-x)%base.PRIME,-1,base.PRIME))%base.PRIME
    elif arm=="public_parity_rule":
        a=SLOPES[(o["role"]+sum(sum(p) for p in o["pairs"]))%6]
    elif arm=="public_support_nearest":
        # Deliberately weak one-pair public role prior (no hidden map).
        a=SLOPES[(o["pairs"][0][0]+o["pairs"][0][1]+o["role"])%6]
    elif arm=="learned_shared_compositional_transition":
        score=[(quant[k+6]+quant[k]*evidence[k])/16 for k in range(6)]
        a=SLOPES[max(range(6),key=lambda k:(score[k],-k))]
    elif arm=="F54_frozen_transition_model":
        # Historical temperature parameters cannot be refitted. Transfer
        # these as a fixed prior only; explicitly not equivalent to old F54.
        role=o["role"]
        score=[12*evidence[k]+(F54_FROZEN_TEMPERATURE_CODES[role*2+(k%2)]-2)*.1
               for k in range(6)]
        a=SLOPES[max(range(6),key=lambda k:(score[k],-k))]
    elif arm=="fixed_equal_work_mass_beam":
        a=SLOPES[max(range(6),key=lambda k:(evidence[k],-k))]
    elif arm=="fixed_equal_work_entropy_beam":
        a=SLOPES[max(range(6),key=lambda k:(evidence[k],int(k%2==0),-k))]
    else:raise ValueError("unregistered arm")
    x,y=o["pairs"][0]
    shift=(y-a*x)%base.PRIME
    return a,shift
def compose(start,chain,maps):
    state=start;trace=[]
    for role,direction in chain:
        a,b=maps[role]
        if direction==1:state=(a*state+b)%base.PRIME
        elif direction==-1:state=((state-b)*pow(a,-1,base.PRIME))%base.PRIME
        else:raise ValueError("invalid binding direction")
        trace.append(state)
    return state,trace
def infer(public,start,chain,arm,quant):
    if arm=="private_truth_oracle_unfair":
        raise ValueError("oracle has no public runtime API")
    need={v["role"]:v for v in public}
    assert set(need)==set(range(4))
    maps={r:select(o,arm,quant) for r,o in need.items()}
    answer,trace=compose(start,chain,maps)
    # Abstract operations: equal candidate/beam matching budget, explicit
    # reference comparator more work. Not walltime/native byte measurement.
    cost=(4*6*2+len(chain)) if arm not in (
        "public_full_algebraic_solver","full_exact_public_posterior") else (
        (4*5+len(chain)) if arm=="public_full_algebraic_solver" else
        (4*6*11*2+len(chain)))
    return answer,cost,maps,trace
def fit_only_report():
    q,n=train();assert n==512
    nfit=0;correct=0
    for seed in base.FIT:
        for fam in base.FIT_FAM:
            for length in (6,8):
                case=base.source(seed,fam,length)
                answer,_,_,_=infer(case["public_calibration"],case["start"],
                                 case["chain"],ARMS[0],q)
                correct+=int(answer==case["evaluator_truth"]);nfit+=1
    return dict(schema="wingless.f55.learned-operator-fit-preflight.v1",
        learned_quantized_signed_int8_parameters=list(q),
        fit_labeled_role_examples=n,fit_cases=nfit,
        fit_learned_correct=correct,precode_fit_seeds=[601000,601016],
        diagnostic_cohort_NOT_EVALUATED=[601100,601116],
        future_original_NOT_RESERVED=[602000,602016],
        source_independent_F3=False,physical_native_576B_unmeasured=True)
def run_diagnostic():
    q,n=train()
    rows=[];arms=ARMS
    for seed in base.DIAG:
        for fam in base.DIAG_FAM:
            for length in (8,10):
                case=base.source(seed,fam,length)
                obs=case["public_calibration"]
                answers={};costs={};maps={};trace={}
                for arm in arms[:-1]:
                    answer,work,m,t=infer(obs,case["start"],case["chain"],arm,q)
                    answers[arm]=answer;costs[arm]=work
                    maps[arm]=[[m[k][0],m[k][1]] for k in sorted(m)]
                    trace[arm]=t
                answers[arms[-1]]=case["evaluator_truth"]
                rows.append(dict(seed=seed,family=fam,length=length,start=case["start"],
                    public_calibration=obs,chain=[list(v) for v in case["chain"]],
                    evaluator_truth=case["evaluator_truth"],answers=answers,
                    abstract_work=costs,public_role_map_receipts=maps,trace=trace))
    scores={a:dict(correct=sum(r["answers"][a]==r["evaluator_truth"] for r in rows),
                   total=len(rows)) for a in arms}
    for s in scores.values():s["accuracy"]=s["correct"]/s["total"]
    by={str(f):{a:sum(r["answers"][a]==r["evaluator_truth"] for r in rows if r["family"]==f)
                for a in arms} for f in base.DIAG_FAM}
    fair=arms[1:5]+arms[6:8] # include cheap exact algebraic solver: ceiling!
    strongest=max(fair,key=lambda a:scores[a]["correct"])
    candidate=arms[0];delta=(scores[candidate]["correct"]-scores[strongest]["correct"])/len(rows)
    exact=scores["full_exact_public_posterior"]["correct"]
    loss=(exact-scores[candidate]["correct"])/len(rows)
    work=sum(r["abstract_work"][candidate] for r in rows)/sum(
        r["abstract_work"]["full_exact_public_posterior"] for r in rows)
    no_family=all(by[str(f)][candidate]>=max(by[str(f)][a] for a in fair)
                  for f in base.DIAG_FAM)
    witness=scores["public_full_algebraic_solver"]["accuracy"]
    passed=witness>=.85 and delta>=.05 and loss<=.02 and work<=.5 and no_family
    return dict(schema="wingless.f55.learned-public-role-bindings.trainonly.v1",
        stage="TRAIN_ONLY_NEW_PUBLIC_IDENTIFIABLE_ROLE_TRANSFER",
        fit_seeds=[base.FIT.start,base.FIT.stop],fit_families=list(base.FIT_FAM),
        diagnostic_seeds=[base.DIAG.start,base.DIAG.stop],
        diagnostic_families=list(base.DIAG_FAM),
        original_seeds_NOT_RESERVED=[base.ORIGINAL.start,base.ORIGINAL.stop],
        fit_labeled_role_examples=n,learned_int8_parameters=list(q),
        learned_model_bytes=12,recurrent_state_payload_bytes=11,
        physical_native_576B_unverified=True,rows=rows,cases=len(rows),
        scores=scores,by_family=by,strongest_fair_control=strongest,
        candidate_gain=delta,loss_vs_exact=loss,
        work_ratio_vs_exact=work,public_constructive_witness_accuracy=witness,
        no_family_regression=no_family,frozen_gate_met=passed,
        F54_frozen_control_ADAPTED_NOT_IDENTICAL=True,
        verdict=("TRAIN_ONLY_LEARNED_ROLE_BINDING_SUPPORTED" if passed
                 else "TRAIN_ONLY_LEARNED_ROLE_BINDING_NEGATIVE"),
        independent_F3=False,primary_original_executed=False,
        accepted_ref_mutation=False,live_optimizer_promotion=False,paid_api_calls=0)
if __name__=="__main__":
    if len(sys.argv)!=2:raise SystemExit("explicit --fit-only or --diagnostic-train flag required")
    if sys.argv[1]=="--fit-only":d=fit_only_report()
    elif sys.argv[1]=="--diagnostic-train":d=run_diagnostic()
    else:raise SystemExit("no F3 or primary original execution")
    print(json.dumps(d,sort_keys=True))
