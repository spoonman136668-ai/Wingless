"""W-R23: abstain when four observations disprove every frozen hypothesis."""
import hashlib
import importlib.util
import json
from pathlib import Path

PRIOR=Path(__file__).resolve().parent.parent/"20261013-w-r22-hypothesis-prereg"/"execute.py"
_spec=importlib.util.spec_from_file_location("w22_frozen",PRIOR)
assert _spec and _spec.loader
w22=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w22)


def sequence(seed,episode):
    assert 84<=seed<=91 and 0<=episode<16
    if seed<88:
        a=2 if seed%2==0 else 1
        b=(seed-84)%4
    else:
        a=3
        b=(seed-88)%4
    seq=[(episode+seed)%4,((episode//4)+3*seed)%4]
    while len(seq)<32:
        seq.append((seq[-1]+a*seq[-2]+b)%4)
    return tuple(seq)


def calibrate(seed):
    return [tuple(sequence(seed,e)[t] for t in range(3)) for e in range(4)]


def forced_model(rows):
    ranked=[]
    for a,b in w22.CANDIDATES:
        errors=sum((last+a*before+b)%4!=true
                   for before,last,true in rows)
        ranked.append((errors,a,b))
    errors,a,b=min(ranked)
    return (a,b,errors)


def evaluate(seed,model,kind):
    if model is None:
        return {"seed":seed,"kind":kind,"abstained":True,"test_predictions":0}
    a,b=model
    hits=0
    trace=hashlib.sha256()
    for ep in range(4,16):
        seq=sequence(seed,ep)
        for tick in range(2,32):
            guess=(seq[tick-1]+a*seq[tick-2]+b)%4
            actual=seq[tick]
            hits+=int(guess==actual)
            trace.update(bytes((ep,tick,actual,guess)))
    return {
        "seed":seed,"kind":kind,"abstained":False,
        "test_predictions":360,"correct":hits,"accuracy":hits/360,
        "trace_sha256":trace.hexdigest()
    }


def run():
    w22.smoke_fixture_only()
    rows=[]
    for seed in range(84,92):
        cal=calibrate(seed)
        matches=w22.select_hypotheses(cal)
        best_a,best_b,best_errors=forced_model(cal)
        accepted=matches[0] if len(matches)==1 else None
        rows.append({
            "seed":seed,
            "heldout_domain":"in_catalog" if seed<88 else "out_of_catalog",
            "calibration_labels":len(cal),
            "exact_hypotheses":matches,
            "strict":evaluate(seed,accepted,"strict_abstention"),
            "forced":evaluate(seed,(best_a,best_b),"forced_best_fit"),
            "forced_calibration_errors":best_errors
        })
    in_domain=[r for r in rows if r["seed"]<88]
    out=[r for r in rows if r["seed"]>=88]
    correct_in=sum(bool(r["strict"]["test_predictions"]==360 and
                        r["strict"]["accuracy"]>=.95) for r in in_domain)
    reject_out=sum(r["strict"]["abstained"] for r in out)
    if correct_in==4 and reject_out==4:
        classification="CORRECT_CATALOG_ABSTENTION"
    elif reject_out>0:
        classification="PARTIAL"
    else:
        classification="FAILURE"
    return {
        "schema":"wingless.w-r23-bounded-catalog-abstention.v1",
        "predecessor_head":"3838bedecfb3cb2133f36846960fddf878498928",
        "classification":classification,
        "candidate_count":len(w22.CANDIDATES),
        "in_catalog_accepted_correct":correct_in,
        "out_catalog_correctly_rejected":reject_out,
        "out_catalog_forced_mean_accuracy":
            sum(r["forced"]["accuracy"] for r in out)/4,
        "rows":rows,
        "model_calls":0,
        "claim_limit":"Exact fit/reject for predeclared known modular recurrences; not calibrated uncertainty outside toy family."
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
