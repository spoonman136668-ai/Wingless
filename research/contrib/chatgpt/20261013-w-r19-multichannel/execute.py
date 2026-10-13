"""W-R19: local swap counterfactual with channel-ablation controls."""
import importlib.util
import json
from pathlib import Path

OLD=Path(__file__).resolve().parent.parent/"20261013-w-r18-adjacent-order"/"execute.py"
spec=importlib.util.spec_from_file_location("w18_frozen",OLD)
assert spec and spec.loader
w18=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w18)
original=w18.old.original
shift=w18.old.nuisance_shift

def pair(seed,entity,positive,phase):
    base=original(seed,entity)
    development=phase=="development"
    modulus=3 if development else 5
    delta=17 if development else 61
    jitter=[((5 if development else 7)*t+seed+entity)%modulus-modulus//2 for t in range(16)]
    assert len(set(jitter))==modulus
    shifted=shift(base,delta)
    observed=tuple(tuple((v+jitter[t])%256 for v in row) for t,row in enumerate(shifted))
    if positive:
        test=observed
    else:
        test=list(observed)
        i=(4 if entity%2==0 else 10) if development else (7 if entity%2==0 else 12)
        test[i],test[i+1]=test[i+1],test[i]
        test=tuple(test)
    assert len(test)==16 and test[0]==observed[0]
    assert sorted(test)==sorted(observed)
    assert (test==observed)==positive
    return base,test


def backward_steps(rows,channel):
    result=0
    for a,b in zip(rows,rows[1:]):
        delta=(b[channel]-a[channel]+128)%256-128
        if delta<0:
            result+=1
    return result


def feature(a,b,kind):
    if kind=="unordered":
        def bins(rows):
            anchor=rows[0][0]
            vals=[0]*16
            for r in rows:vals[((r[0]-anchor)%256)%16]+=1
            return vals
        u,v=bins(a),bins(b)
        return sum(abs(x-y) for x,y in zip(u,v))/16
    if kind=="first_backward":
        return backward_steps(b,0)
    if kind=="fourth_backward":
        return backward_steps(b,3)
    assert kind=="all_channel_backward_mean"
    return sum(backward_steps(b,c) for c in range(4))/4


def evaluate(kind):
    pos=[]
    neg=[]
    for seed in range(56,60):
        for entity in range(9):
            for positive in (True,False):
                a,b=pair(seed,entity,positive,"development")
                (pos if positive else neg).append(feature(a,b,kind))
    assert len(pos)==len(neg)==36
    maxp=max(pos)
    minn=min(neg)
    threshold=(maxp+minn)/2
    trace=[]
    for seed in range(60,64):
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=pair(seed,entity,positive,"heldout")
                pred=1 if feature(a,b,kind)<=threshold else -1
                true=1 if positive else -1
                if pred==1 and true==1:tp+=1
                elif pred==-1 and true==-1:tn+=1
                elif pred==1:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        trace.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,"accuracy":(tp+tn)/18})
    return {"arm":kind,"train_positive_max":maxp,"train_negative_min":minn,
            "training_separable":maxp<minn,"train_only_threshold":threshold,
            "learned_parameters":1,"heldout_accuracy":sum(r["tp"]+r["tn"] for r in trace)/72,
            "heldout_by_seed":trace}


def run():
    arms=[evaluate(kind) for kind in ("unordered","first_backward","fourth_backward","all_channel_backward_mean")]
    acc={a["arm"]:a["heldout_accuracy"] for a in arms}
    if acc["fourth_backward"]>=.85 and acc["unordered"]<=.55:
        label="FOURTH_CHANNEL_TRANSFER"
    elif acc["fourth_backward"]>acc["first_backward"]:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r19-multichannel-local.v1",
            "predecessor_head":"a1c11d76e670a4838a7b05d2b6905d2318183394",
            "classification":label,"arms":arms,"model_calls":0,
            "claim_limit":"The fourth-channel monotonic signal is synthetic, not language comprehension"}


if __name__=="__main__":
    r=run()
    assert r==run()
    print(json.dumps(r,sort_keys=True,indent=2))
