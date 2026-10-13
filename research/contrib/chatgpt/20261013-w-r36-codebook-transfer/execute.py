"""W-R36 frozen zero-shot byte-codebook transfer of W-R34 models."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r35-cohort-replication"/"execute.py"
spec=importlib.util.spec_from_file_location("w35_pinned",SOURCE)
assert spec and spec.loader
w35=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w35)
w34=w35.w34

HELDOUT=tuple(range(1600,1632))
TRANSFORMS=("identity","byte_invert","xor55","channel_rotate","offset37")
ARMS=w34.ARMS
BITS=4
GENERATORS=w34.GENERATORS
SITES=w34.SITES


def recode(rows,fmt):
    result=[]
    for row in rows:
        if fmt=="identity":
            changed=tuple(row)
        elif fmt=="byte_invert":
            changed=tuple(255-v for v in row)
        elif fmt=="xor55":
            changed=tuple(v^0x55 for v in row)
        elif fmt=="channel_rotate":
            changed=tuple(row[1:])+tuple(row[:1])
        else:
            assert fmt=="offset37"
            changed=tuple((v+37)%256 for v in row)
        assert len(changed)==4
        assert all(isinstance(v,int) and 0<=v<=255 for v in changed)
        result.append(changed)
    return tuple(result)


def cell(arm,fmt,model,generator,site):
    weights,bias,updates,_=model
    by_seed=[]
    correct=total=indistinguishable=0
    for seed in HELDOUT:
        tp=tn=fp=fn=identical=0
        for entity in range(9):
            base,positive,negative=w34.w32.observed_pair(seed,entity,generator,BITS,site)
            transformed_base=recode(base,fmt)
            transformed_positive=recode(positive,fmt)
            transformed_negative=recode(negative,fmt)
            assert (positive==negative)==(transformed_positive==transformed_negative)
            assert sorted(transformed_positive)==sorted(transformed_negative)
            identical+=int(transformed_positive==transformed_negative)
            for truth,candidate in ((True,transformed_positive),(False,transformed_negative)):
                v=w34.representation(transformed_base,candidate,arm,BITS)
                guess=bias+sum(w*x for w,x in zip(weights,v))>=0
                if truth and guess:tp+=1
                elif not truth and not guess:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        indistinguishable+=identical
        by_seed.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                        "identical_positive_negative_pairs":identical})
    return {"generator":generator,"site":site,"correct":correct,"total":total,
            "accuracy":correct/total,"observation_identical_pairs":indistinguishable,
            "by_seed":by_seed}


def run():
    models={arm:w34.train(arm,BITS) for arm in ARMS}
    rows=[]
    for fmt in TRANSFORMS:
        for arm in ARMS:
            model=models[arm]
            cells=[cell(arm,fmt,model,generator,site)
                   for generator in GENERATORS for site in SITES]
            rows.append({"format":fmt,"arm":arm,"bits":BITS,
                         "training_examples":1728,"epochs":32,"width":16,
                         "training_updates":model[2],
                         "accuracy":sum(x["correct"] for x in cells)/
                                    sum(x["total"] for x in cells),"cells":cells})
    by={(row["format"],row["arm"]):row for row in rows}
    losses={}
    for fmt in TRANSFORMS:
        if fmt=="identity":continue
        losses[fmt]={arm:by[("identity",arm)]["accuracy"]-by[(fmt,arm)]["accuracy"]
                     for arm in ARMS}
    affected=sum(any(loss>=.10 for loss in losses[fmt].values()) for fmt in losses)
    if affected>=2: label="CODEBOOK_SENSITIVE"
    elif affected==1: label="PARTIAL"
    else: label="TRANSFER_STABLE"
    return {"schema":"wingless.w-r36-bijective-codebook-transfer.v1",
            "parent_head":"7528b49520035347597967e476392ef4a2b97a15",
            "classification":label,"sealed_heldout_seeds":HELDOUT,
            "formats":TRANSFORMS,"four_bit_loss_by_format_and_arm":losses,
            "arms":rows,"model_calls":0,
            "claim_limit":"Synthetic byte-label transform invariance of frozen weights, not semantic transfer."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
