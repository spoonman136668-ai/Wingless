"""W-R34 frozen symmetric-vs-asymmetric sensor codec representations."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r33-quantized-retraining"/"execute.py"
spec=importlib.util.spec_from_file_location("w33_pinned",SOURCE)
assert spec and spec.loader
w33=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w33)
w32=w33.w32
w25=w33.w25
w20=w33.w20
w16=w33.w16

TRAIN=tuple(range(1184,1200))
HELDOUT=tuple(range(1200,1216))
BITS=(8,4,2,1)
ARMS=("asymmetric_refit","symmetric_refit")
GENERATORS=w33.GENERATORS
SITES=w33.SITES
INDICES=w33.INDICES


def representation(base,candidate,arm,bits):
    if arm=="symmetric_refit":
        base=tuple(w32.quantize(row,bits) for row in base)
    else:
        assert arm=="asymmetric_refit"
    v=w25.new_feature(base,candidate)
    assert len(v)==16
    return v


def train(arm,bits):
    examples=[]
    identical=0
    for seed in TRAIN:
        for entity in range(9):
            base=w20.base_sequence(seed,entity,"affine_prior")
            for site in INDICES:
                for amplitude in (1,2):
                    original=w16.observed(base,seed,entity,amplitude,17)
                    positive=tuple(w32.quantize(row,bits) for row in original)
                    tail=list(positive[1:])
                    tail[site],tail[site+1]=tail[site+1],tail[site]
                    negative=positive[:1]+tuple(tail)
                    assert sorted(positive)==sorted(negative)
                    identical+=int(positive==negative)
                    examples.append((representation(base,positive,arm,bits),1))
                    examples.append((representation(base,negative,arm,bits),-1))
    assert len(examples)==1728
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            if label*(bias+sum(w*x for w,x in zip(weights,vec)))<=0:
                weights=[w+0.1*label*x for w,x in zip(weights,vec)]
                bias+=0.1*label
                updates+=1
    return weights,bias,updates,identical


def evaluate(arm,bits,model,generator,site):
    weights,bias,updates,identical=model
    by_seed=[]
    correct=total=observationally_identical=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=0
        for entity in range(9):
            base,positive,negative=w32.observed_pair(seed,entity,generator,bits,site)
            same+=int(positive==negative)
            for truth,candidate in ((True,positive),(False,negative)):
                feat=representation(base,candidate,arm,bits)
                pred=bias+sum(w*x for w,x in zip(weights,feat))>=0
                if truth and pred: tp+=1
                elif not truth and not pred: tn+=1
                elif not truth: fp+=1
                else: fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        observationally_identical+=same
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "identical_positive_negative_pairs":same})
    return {"generator":generator,"swap_site":site,
            "correct":correct,"total":total,"accuracy":correct/total,
            "identical_pairs":observationally_identical,"by_seed":by_seed}


def run():
    arms=[]
    for bits in BITS:
        for arm in ARMS:
            model=train(arm,bits)
            cells=[evaluate(arm,bits,model,g,s)
                   for g in GENERATORS for s in SITES]
            arms.append({"arm":arm,"bits":bits,"width":16,
                         "training_examples":1728,"epochs":32,
                         "training_updates":model[2],
                         "identical_training_pairs":model[3],
                         "aggregate_accuracy":sum(c["correct"] for c in cells)/
                                              sum(c["total"] for c in cells),
                         "cells":cells})
    by={(a["bits"],a["arm"]):a for a in arms}
    gain={str(bits):by[(bits,"symmetric_refit")]["aggregate_accuracy"]-
                   by[(bits,"asymmetric_refit")]["aggregate_accuracy"] for bits in BITS}
    sym4=by[(4,"symmetric_refit")]["aggregate_accuracy"]
    if gain["4"]>=.10 and sym4>=.70 and gain["8"]>=-.05:
        label="CODEC_ALIGNMENT_RECOVERY"
    elif gain["4"]>=.05 or gain["2"]>=.05:
        label="PARTIAL_GAIN"
    else:
        label="NO_GAIN"
    return {"schema":"wingless.w-r34-symmetric-quantizer-feature.v1",
            "parent_head":"f7113cf52145cd5f6cb2be43a88fedd38fbd8d94",
            "classification":label,
            "development_seeds":TRAIN,"sealed_heldout_seeds":HELDOUT,
            "accuracy_gain_symmetric_minus_asymmetric":gain,
            "arms":arms,"model_calls":0,
            "claim_limit":"Clean synthetic quantized temporal pair classification, not general language reasoning."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
