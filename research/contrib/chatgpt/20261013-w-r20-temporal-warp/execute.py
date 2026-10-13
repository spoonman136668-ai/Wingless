"""W-R20: bounded monotone one-repeat/skip timing recognition."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r18-local-order"/"execute.py"
spec=importlib.util.spec_from_file_location("wr18_frozen",PARENT)
assert spec and spec.loader
prior=importlib.util.module_from_spec(spec)
spec.loader.exec_module(prior)

ARMS=("unordered","position","transition","monotone_warp")


def observed_pair(seed,entity,positive,phase):
    dev=phase=="development"
    base=prior.original(seed,entity)
    shift=17 if dev else 61
    m=3 if dev else 5
    multiplier=5 if dev else 7
    noises=[(multiplier*t+entity+seed)%m-(m//2) for t in range(16)]
    assert len(set(noises))==m
    transformed=tuple(
        tuple((v+shift+noises[t])%256 for v in row)
        for t,row in enumerate(base)
    )
    warp_at=6 if dev else 9
    warped=tuple(transformed[t-1] if t==warp_at else transformed[t]
                 for t in range(16))
    if positive:
        candidate=warped
    elif dev:
        tail=warped[1:]
        candidate=warped[:1]+(tail[::-1] if entity%2==0 else tail[3:]+tail[:3])
    else:
        changed=list(warped)
        changed[11],changed[12]=changed[12],changed[11]
        candidate=tuple(changed)
    assert candidate[0]==warped[0]
    assert sorted(candidate)==sorted(warped)
    assert (candidate==warped)==positive
    assert len(candidate)==16
    return base,candidate


def score(a,b,kind):
    x=prior.anchor_relative(a)
    y=prior.anchor_relative(b)
    if kind=="unordered":
        bins_a=[0]*16
        bins_b=[0]*16
        for u,v in zip(x,y):
            bins_a[u%16]+=1
            bins_b[v%16]+=1
        return sum(abs(u-v) for u,v in zip(bins_a,bins_b))/32
    if kind=="position":
        return sum(prior.circular_distance(u,v) for u,v in zip(x,y))/16
    if kind=="transition":
        return sum(
            prior.circular_distance((x[t+1]-x[t])%256,(y[t+1]-y[t])%256)
            for t in range(15)
        )/15
    assert kind=="monotone_warp"
    # Every possible positive single-repeat/skip alignment is frozen
    # before first use. No dynamic model weights or oracle-provided index.
    return min(
        sum(prior.circular_distance(
            x[t-1] if t==warp_at else x[t],y[t])
            for t in range(16)
        )/16
        for warp_at in range(1,16)
    )


def develop(kind):
    yes=[]
    no=[]
    for seed in range(64,68):
        for entity in range(9):
            for positive in (True,False):
                a,b=observed_pair(seed,entity,positive,"development")
                (yes if positive else no).append(score(a,b,kind))
    assert len(yes)==len(no)==36
    high=max(yes)
    low=min(no)
    return {"threshold":(high+low)/2,"train_max_positive":high,
            "train_min_negative":low,"train_separable":high<low}


def holdout(kind,fit):
    rows=[]
    for seed in range(68,72):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=observed_pair(seed,entity,positive,"heldout")
                pred=score(a,b,kind)<=fit["threshold"]
                if positive and pred:tp+=1
                elif not positive and not pred:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "accuracy":(tp+tn)/18})
    return {
        "arm":kind,"learned_parameters":1,
        "fit":fit,"accuracy":sum(r["tp"]+r["tn"] for r in rows)/72,
        "heldout":rows
    }


def run():
    arms=[holdout(kind,develop(kind)) for kind in ARMS]
    accuracies={row["arm"]:row["accuracy"] for row in arms}
    aligned=accuracies["monotone_warp"]
    if (aligned>=.85 and accuracies["unordered"]<=.55
            and aligned>accuracies["position"]):
        classification="WARP_TRANSFER"
    elif aligned>max(accuracies["unordered"],accuracies["position"]):
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {
        "schema":"wingless.w-r20-monotone-warp.v1",
        "predecessor_head":"3e96d6f45819abf75ed3d99dbf68b85f8792f098",
        "classification":classification,
        "alignment_candidates":15,
        "arms":arms,"model_calls":0,
        "scope":"Synthetic one-repeat/skip alignment with jitter, not language reasoning"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
