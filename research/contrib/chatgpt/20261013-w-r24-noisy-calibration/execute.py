"""W-R24: frozen noisy four-symbol hypothesis selection with calibrated abstention."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r23-open-set" / "execute.py"
spec = importlib.util.spec_from_file_location("w23_frozen", PREVIOUS)
assert spec and spec.loader
w23 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w23)
CANDIDATES = w23.CANDIDATES


def dynamics(seed, episode):
    assert 96 <= seed <= 103 and 0 <= episode < 16
    a = 2 if seed % 2 == 0 else 1
    b = (seed-96)%4
    seq = [(episode+seed)%4, ((episode//4)+3*seed)%4]
    while len(seq)<32:
        seq.append(w23.w22.transition(seq[-2],seq[-1],a,b))
    return tuple(seq)


def calibration(seed, corrupt_episodes):
    out=[]
    for episode in range(8):
        seq=dynamics(seed,episode)
        label=(seq[2]+(1 if episode in corrupt_episodes else 0))%4
        out.append((seq[0],seq[1],label))
    assert len(out)==8
    assert len({(x,y) for x,y,_ in out})==8
    return out


def mismatch(model,cases):
    a,b=model
    return sum(w23.w22.transition(x,y,a,b)!=target
               for x,y,target in cases)


def choose(cases, allowed_errors):
    counts=[(mismatch(candidate,cases),candidate) for candidate in CANDIDATES]
    minimum=min(count for count,model in counts)
    best=sorted(candidate for count,candidate in counts if count==minimum)
    return (best[0] if len(best)==1 and minimum<=allowed_errors else None),minimum,len(best)


def fixture_only():
    # Artificial 8-context fixture, no science generator seed.
    contexts=[(x,y) for y in (0,1) for x in range(4)]
    base=[(x,y,w23.w22.transition(x,y,2,3)) for x,y in contexts]
    one=[(x,y,(z+int(i==0))%4) for i,(x,y,z) in enumerate(base)]
    two=[(x,y,(z+int(i in (0,4)))%4) for i,(x,y,z) in enumerate(base)]
    assert choose(one,1)[0]==(2,3)
    assert choose(two,1)[0] is None
    return "PASS"


def test(seed,model):
    confusion=[[0]*4 for _ in range(4)]
    trace=hashlib.sha256()
    for episode in range(8,16):
        seq=dynamics(seed,episode)
        for tick in range(2,32):
            pred=w23.w22.transition(seq[tick-2],seq[tick-1],model[0],model[1])
            true=seq[tick]
            confusion[true][pred]+=1
            trace.update(bytes((episode,tick,true,pred)))
    assert sum(map(sum,confusion))==240
    correct=sum(confusion[i][i] for i in range(4))
    return {"correct":correct,"total":240,"accuracy":correct/240,
            "confusion":confusion,"trace_sha256":trace.hexdigest()}


def run():
    assert fixture_only()=="PASS"
    conditions=(("ONE_CORRUPT",(0,)),("TWO_CORRUPT",(0,4)))
    rows=[]
    for condition,corrupt in conditions:
        for seed in range(96,104):
            triples=calibration(seed,corrupt)
            robust,err,ties=choose(triples,1)
            strict,_,_=choose(triples,0)
            forced=min(CANDIDATES,key=lambda c:(mismatch(c,triples),c))
            rows.append({
                "seed":seed,"condition":condition,"best_errors":err,
                "best_tie_count":ties,"robust_accepted":robust is not None,
                "strict_accepted":strict is not None,
                "robust_model":robust,
                "robust_test":test(seed,robust) if robust else None,
                "forced_model":forced,"forced_test":test(seed,forced),
            })
    one=[r for r in rows if r["condition"]=="ONE_CORRUPT"]
    two=[r for r in rows if r["condition"]=="TWO_CORRUPT"]
    accepted=sum(r["robust_accepted"] for r in one)
    rejected_two=sum(not r["robust_accepted"] for r in two)
    mistakes=sum(r["robust_test"]["total"]-r["robust_test"]["correct"]
                 for r in rows if r["robust_accepted"])
    if accepted==8 and rejected_two==8 and mistakes==0 and all(
            not r["strict_accepted"] for r in rows):
        label="ONE_ERROR_RECOVERY_TWO_ERROR_REJECTION"
    elif accepted>0 and mistakes==0:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {
        "schema":"wingless.w-r24-noisy-calibration.v1",
        "predecessor_head":"0c899bceb80f3bb54e3f98e797ed5e547c0d62bc",
        "classification":label,"fixture_only":"PASS",
        "one_error_accepted_seeds":accepted,
        "two_error_rejected_seeds":rejected_two,
        "strict_zero_error_accepted":sum(r["strict_accepted"] for r in rows),
        "wrong_accepted_predictions":mistakes,
        "forced_one_accuracy":sum(r["forced_test"]["correct"] for r in one)/(8*240),
        "forced_two_accuracy":sum(r["forced_test"]["correct"] for r in two)/(8*240),
        "rows":rows,"model_calls":0,
        "claim_limit":"Synthetic fixed-candidate mismatch tolerance on two preset corruptions; not general noisy learning.",
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
