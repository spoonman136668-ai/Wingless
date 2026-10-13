"""W-R25: counterfeit calibration evidence gives wrong confident model."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS=Path(__file__).resolve().parent.parent/"20261013-w-r24-noisy-calibration"/"execute.py"
spec=importlib.util.spec_from_file_location("w24_frozen",PREVIOUS)
assert spec and spec.loader
w24=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w24)
t=w24.w23.w22.transition
CANDIDATES=w24.CANDIDATES


def true_parameters(seed):
    assert 104<=seed<=111
    return (2 if seed%2==0 else 1),(seed-104)%4


def alternative_parameters(seed):
    a,b=true_parameters(seed)
    return (1 if a==2 else 2),(b+1)%4


def sequence(seed,episode):
    assert 0<=episode<16
    a,b=true_parameters(seed)
    values=[(episode+seed)%4,((episode//4)+3*seed)%4]
    while len(values)<32:
        values.append(t(values[-2],values[-1],a,b))
    return tuple(values)


def evidence(seed,condition):
    assert condition in ("HONEST_FOUR","SPOOF_FOUR","SPOOF_FOUR_PLUS_TRUSTED_FOUR")
    spoof=condition!="HONEST_FOUR"
    n=8 if condition=="SPOOF_FOUR_PLUS_TRUSTED_FOUR" else 4
    rows=[]
    for episode in range(n):
        seq=sequence(seed,episode)
        if spoof and episode<4:
            a,b=alternative_parameters(seed)
            received_label=t(seq[0],seq[1],a,b)
        else:
            received_label=seq[2]
        rows.append((seq[0],seq[1],received_label))
    assert len(rows)==n
    assert len({(a,b) for a,b,_ in rows})==n
    return rows


def select(exemplars):
    matching=[
        (a,b) for a,b in CANDIDATES
        if all(t(before,last,a,b)==label for before,last,label in exemplars)
    ]
    return matching[0] if len(matching)==1 else None,matching


def fixture():
    p=[(x,0,t(x,0,2,1)) for x in range(4)]
    c=[(x,0,t(x,0,1,2)) for x in range(4)]
    assert select(p)[0]==(2,1)
    assert select(c)[0]==(1,2)
    # Conflicting four examples for a new context cannot be satisfied together.
    assert select(p+c)[0] is None
    return "PASS"


def evaluate(seed,model):
    if model is None:
        return None
    cm=[[0]*4 for _ in range(4)]
    digest=hashlib.sha256()
    for episode in range(8,16):
        seq=sequence(seed,episode)
        for tick in range(2,32):
            prediction=t(seq[tick-2],seq[tick-1],model[0],model[1])
            actual=seq[tick]
            cm[actual][prediction]+=1
            digest.update(bytes((episode,tick,actual,prediction)))
    assert sum(map(sum,cm))==240
    correct=sum(cm[i][i] for i in range(4))
    return {"correct":correct,"total":240,
            "accuracy":correct/240,"confusion":cm,"trace_sha256":digest.hexdigest()}


def run():
    assert fixture()=="PASS"
    conditions=("HONEST_FOUR","SPOOF_FOUR","SPOOF_FOUR_PLUS_TRUSTED_FOUR")
    rows=[]
    for seed in range(104,112):
        for condition in conditions:
            source=evidence(seed,condition)
            selected,matches=select(source)
            rows.append({
                "seed":seed,"condition":condition,"labels_available":len(source),
                "accepted":selected is not None,"selected":selected,
                "matching_count":len(matches),
                "test":evaluate(seed,selected),
            })
    def subset(name):
        return [r for r in rows if r["condition"]==name]
    honest=subset("HONEST_FOUR")
    poisoned=subset("SPOOF_FOUR")
    combined=subset("SPOOF_FOUR_PLUS_TRUSTED_FOUR")
    honest_accept=sum(r["accepted"] for r in honest)
    spoof_accept=sum(r["accepted"] for r in poisoned)
    later_reject=sum(not r["accepted"] for r in combined)
    honest_accuracy=sum(r["test"]["correct"] for r in honest if r["accepted"])/(max(1,honest_accept)*240)
    spoof_accuracy=sum(r["test"]["correct"] for r in poisoned if r["accepted"])/(max(1,spoof_accept)*240)
    if honest_accept==8 and honest_accuracy==1 and spoof_accept==8 and spoof_accuracy<1 and later_reject==8:
        label="SPOOF_ACCEPTANCE_FAILURE_AND_LATE_REJECTION"
    elif any(r["accepted"] for r in combined):
        label="FALSE_REASSURANCE"
    else:
        label="OTHER"
    return {
        "schema":"wingless.w-r25-poisoned-calibration.v1",
        "predecessor_head":"a235914611ff2a6c5280ad44d7f902f89b67f43a",
        "classification":label,"fixture":"PASS",
        "honest_accepted":honest_accept,"honest_test_accuracy":honest_accuracy,
        "spoof_accepted":spoof_accept,"spoof_test_accuracy":spoof_accuracy,
        "mixed_evidence_rejected":later_reject,
        "calibration_examples_initial":4,"calibration_examples_extended":8,
        "rows":rows,"model_calls":0,
        "claim_limit":"Counterfeit fixed-family synthetic teacher labels, not an adversarial robustness claim on real systems.",
    }


if __name__=="__main__":
    data=run()
    assert data==run()
    print(json.dumps(data,sort_keys=True,indent=2))
