"""W-R25: preregistered reference-signature alignment without missingness oracle."""
import importlib.util
import json
from pathlib import Path

BASE=Path(__file__).resolve().parent.parent/"20261013-w-r24-missingness"/"execute.py"
spec=importlib.util.spec_from_file_location("w24_frozen",BASE)
assert spec and spec.loader
w24=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w24)
w16=w24.w16

DEV=tuple(range(320,336))
HELDOUT=tuple(range(336,352))
ARMS=("zero_sentinel","structural_fill","reference_signature_alignment")
UNSEEN_FILLS=("constant73","linear_neighbour_interpolation")
GENS=("affine_prior","quadratic_irregular","xorshift_independent")
SWAPS=("swap10_11","swap12_13")


def aligned_signature(a,b):
    assert len(a)==len(b)==16
    offsets=[(b[0][channel]-a[0][channel])%256 for channel in range(4)]
    matched=[]
    for tick,row in enumerate(b):
        best_index=0
        best_error=999.0
        for ref in range(16):
            delta=sum(w16.circular_distance((a[ref][c]+offsets[c])%256,row[c])
                      for c in range(4))/4.0
            if delta<best_error:
                best_error=delta
                best_index=ref
        # circular_distance is normalized by 128 in the frozen W-R16B source.
        if best_error <= 12.0/128.0:
            displacement=min((tick-best_index)%16,(best_index-tick)%16)/8.0
            matched.append(displacement)
        else:
            matched.append(0.0)
    assert len(matched)==16
    return matched


def feature(a,b,kind):
    if kind=="reference_signature_alignment":
        return aligned_signature(a,b)
    assert kind in ("zero_sentinel","structural_fill")
    return w24.feature(a,b,kind)


def train(kind):
    examples=[]
    for seed in DEV:
        for entity in range(9):
            for swap in w24.TRAIN_SWAPS:
                for amp in (1,2):
                    fill=w24.TRAIN_FILLS[(seed+entity+swap+amp)%3]
                    count=(2,4)[(seed+entity+amp)%2]
                    for positive in (True,False):
                        a,b=w24.pair(seed,entity,amp,"development",swap,
                                     fill,count,positive,"affine_prior")
                        examples.append((feature(a,b,kind),1 if positive else -1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            signed=label*(bias+sum(w*x for w,x in zip(weights,vec)))
            if signed<=0:
                weights=[w+.1*label*x for w,x in zip(weights,vec)]
                bias+=.1*label
                updates+=1
    return weights,bias,updates


def test(kind,model,fill,generator,count,swap):
    weights,bias,_=model
    seeds=[]
    for seed in HELDOUT:
        tp=tn=fp=fn=0
        for entity in range(9):
            for positive in (True,False):
                a,b=w24.pair(seed,entity,3,"heldout",swap,fill,count,
                             positive,generator)
                vec=feature(a,b,kind)
                predicted=bias+sum(w*x for w,x in zip(weights,vec))>=0
                if positive and predicted:tp+=1
                elif not positive and not predicted:tn+=1
                elif not positive:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        seeds.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn})
    correct=sum(r["tp"]+r["tn"] for r in seeds)
    return {"fill":fill,"generator":generator,"masked":count,"swap":swap,
            "correct":correct,"total":288,"accuracy":correct/288,"seed_rows":seeds}


def run():
    arms=[]
    for kind in ARMS:
        model=train(kind)
        cells=[test(kind,model,fill,gen,count,swap)
               for fill in UNSEEN_FILLS for gen in GENS
               for count in (2,4,6) for swap in SWAPS]
        assert len(cells)==36
        accuracy=sum(c["correct"] for c in cells)/sum(c["total"] for c in cells)
        by_fill={fill:sum(c["correct"] for c in cells if c["fill"]==fill)/
                       sum(c["total"] for c in cells if c["fill"]==fill)
                 for fill in UNSEEN_FILLS}
        arms.append({"arm":kind,"width":len(model[0]),"examples":1728,
                     "epochs":32,"updates":model[2],
                     "accuracy":accuracy,"by_fill":by_fill,
                     "cells":cells})
    by={arm["arm"]:arm for arm in arms}
    aligned=by["reference_signature_alignment"]
    structure=by["structural_fill"]
    gain=aligned["accuracy"]-structure["accuracy"]
    if all(c["accuracy"]>=.85 for c in aligned["cells"]) and gain>=.10:
        label="ALIGNMENT_GENERALIZES"
    elif gain>=.05:
        label="PARTIAL"
    else:
        label="NEGATIVE"
    return {"schema":"wingless.w-r25-reference-alignment.v1",
            "parent_head":"bf81aafa6a43b040b00ab5f5b22ceda489c217c2",
            "classification":label,"gain_vs_structural":gain,
            "development_seeds":DEV,"heldout_seeds":HELDOUT,
            "arms":arms,"model_calls":0,
            "claim_limit":"Synthetic observed-only reference alignment, no general reasoning"}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
