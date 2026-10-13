"""W-R18 frozen one-swap local temporal sensitivity falsification."""
import importlib.util
import json
from pathlib import Path

PRIOR = Path(__file__).resolve().parent.parent / "20261012-w-r17-true-tick-jitter" / "execute.py"
spec = importlib.util.spec_from_file_location("w17_frozen", PRIOR)
assert spec and spec.loader
prior = importlib.util.module_from_spec(spec)
spec.loader.exec_module(prior)
old = prior.w12


def pairs(seed, entity, positive, phase):
    reference = old.original(seed, entity)
    development = phase == "development"
    modulus = 3 if development else 5
    delta = 17 if development else 61
    jitter = tuple(((5 if development else 7)*t+seed+entity) % modulus
                   - modulus//2 for t in range(16))
    assert len(set(jitter)) == modulus
    shifted = old.nuisance_shift(reference, delta)
    observed = tuple(tuple((value+jitter[t]) % 256 for value in row)
                     for t,row in enumerate(shifted))
    if positive:
        candidate = observed
    elif development:
        tail = observed[1:]
        candidate = observed[:1]+(tail[::-1] if entity % 2 == 0
                                 else tail[3:]+tail[:3])
    else:
        arr = list(observed)
        arr[7], arr[8] = arr[8], arr[7]
        candidate = tuple(arr)
    assert observed[0] == candidate[0]
    assert sorted(observed) == sorted(candidate)
    assert (observed == candidate) == positive
    return reference, candidate


def distance(x,y):
    d=(x-y)%256
    return min(d,256-d)/128.


def values(reference,candidate,kind):
    anchor1,anchor2=reference[0][0],candidate[0][0]
    a=[(r[0]-anchor1)%256 for r in reference]
    b=[(r[0]-anchor2)%256 for r in candidate]
    if kind=="unordered":
        x=[0]*16
        y=[0]*16
        for p,q in zip(a,b):
            x[p%16]+=1
            y[q%16]+=1
        return sum(abs(p-q) for p,q in zip(x,y))/16.
    if kind=="pooled":
        return sum(distance(p,q) for p,q in zip(a,b))/16.
    if kind=="max_position":
        return max(distance(p,q) for p,q in zip(a,b))
    assert kind=="max_transition"
    t1=[(a[(i+1)%16]-a[i])%256 for i in range(16)]
    t2=[(b[(i+1)%16]-b[i])%256 for i in range(16)]
    return max(distance(p,q) for p,q in zip(t1,t2))


def example(seed,entity,positive,phase,kind):
    a,b=pairs(seed,entity,positive,phase)
    return values(a,b,kind),1 if positive else -1


def arm(kind):
    pos=[]
    neg=[]
    for seed in range(48,52):
        for entity in range(9):
            for positive in (True,False):
                feature,label=example(seed,entity,positive,"development",kind)
                (pos if positive else neg).append(feature)
    assert len(pos)==len(neg)==36
    maxpos=max(pos)
    minneg=min(neg)
    threshold=(maxpos+minneg)/2
    trace=[]
    for seed in range(52,56):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                feature,label=example(seed,entity,positive,"heldout",kind)
                prediction=1 if feature<=threshold else -1
                if label==1 and prediction==1: tp+=1
                elif label==-1 and prediction==-1: tn+=1
                elif label==-1: fp+=1
                else: fn+=1
        assert tp+tn+fp+fn==18
        trace.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                      "accuracy":(tp+tn)/18})
    return {"arm":kind,"parameters":1,"train_max_positive":maxpos,
            "train_min_negative":minneg,"train_range_separable":maxpos<minneg,
            "train_only_threshold":threshold,"heldout_by_seed":trace,
            "heldout_accuracy":sum(r["tp"]+r["tn"] for r in trace)/72}


def run():
    arms=[arm(x) for x in ("unordered","pooled","max_position","max_transition")]
    score={a["arm"]:a["heldout_accuracy"] for a in arms}
    if score["max_transition"]>=.85 and score["unordered"]<=.55:
        classification="LOCAL_ORDER_TRANSFER"
    elif score["max_transition"]>max(score["pooled"],score["unordered"]):
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {"schema":"wingless.w-r18-local-order.v1",
            "predecessor_head":"21b3c3fee08e5af3ed6c22218d67f59ed40d60a9",
            "classification":classification,"arms":arms,"model_calls":0,
            "claim_limit":"Toy adjacent-swap discrimination, not language prediction or general intelligence"}


if __name__=="__main__":
    outcome=run()
    assert outcome==run()
    print(json.dumps(outcome,sort_keys=True,indent=2))
