"""W-R42: frozen decision-tree versus linear readout on identical XOR-invariant lag2 state."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r41-feature-audit"/"execute.py"
spec=importlib.util.spec_from_file_location("w41_pinned",SOURCE)
assert spec and spec.loader
w41=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w41)
w40=w41.w40
w37=w41.w37

TRAIN=tuple(range(2032,2048))
HELDOUT=tuple(range(2048,2064))
CONDITIONS=("identity","xor55")
GENERATORS=w40.GENERATORS
TEST_SITES=w40.TEST_SITES
TRAIN_SITES=w40.TRAIN_SITES
ARMS=("linear16","depth2_tree")


def feature(base,candidate):
    v=w40.feature(base,candidate,"ordered_lag2")
    assert len(v)==16 and all(isinstance(x,int) and 0<=x<=32 for x in v)
    return tuple(v)


def development_examples():
    examples=[]
    same=0
    for seed in TRAIN:
        for entity in range(9):
            for site in TRAIN_SITES:
                for amplitude in (1,2):
                    base,pos,neg=w37.recoded_example(
                        seed,entity,"affine_prior","identity",site,amplitude,True)
                    same+=int(pos==neg)
                    examples.extend(((feature(base,pos),1),(feature(base,neg),-1)))
    assert len(examples)==1728
    return examples,same


def fit_linear(examples):
    weights=[0.0]*16
    bias=0.0
    updates=0
    for _ in range(32):
        for vec,label in examples:
            if label*(bias+sum(w*x for w,x in zip(weights,vec)))<=0:
                weights=[w+.1*label*x for w,x in zip(weights,vec)]
                bias+=.1*label
                updates+=1
    return {"weights":weights,"bias":bias,"updates":updates,"model_scalars":17}


def majority(pos,neg):
    return 1 if pos>neg else -1


def fit_stump(examples):
    # A development-only fixed exhaustive split grid, not heldout search.
    if not examples:
        return {"idx":0,"threshold":0,"left":-1,"right":-1,
                "train_correct":0,"candidates":16*33}
    best=None
    total_pos=sum(label>0 for _,label in examples)
    total_neg=len(examples)-total_pos
    for idx in range(16):
        pos_hist=[0]*33
        neg_hist=[0]*33
        for vec,label in examples:
            if label>0:pos_hist[vec[idx]]+=1
            else:neg_hist[vec[idx]]+=1
        left_pos=left_neg=0
        for threshold in range(33):
            left_pos+=pos_hist[threshold]
            left_neg+=neg_hist[threshold]
            right_pos=total_pos-left_pos
            right_neg=total_neg-left_neg
            correct=max(left_pos,left_neg)+max(right_pos,right_neg)
            if best is None or correct>best["train_correct"]:
                best={"idx":idx,"threshold":threshold,
                      "left":majority(left_pos,left_neg),
                      "right":majority(right_pos,right_neg),
                      "train_correct":correct,"candidates":16*33}
    return best


def fit_tree(examples):
    root=fit_stump(examples)
    left=[(v,y) for v,y in examples if v[root["idx"]]<=root["threshold"]]
    right=[(v,y) for v,y in examples if v[root["idx"]]>root["threshold"]]
    child_left=fit_stump(left)
    child_right=fit_stump(right)
    return {"root":root,"child_left":child_left,"child_right":child_right,
            "development_left":len(left),"development_right":len(right),
            "train_correct":sum(predict_stump(child_left,v)==y for v,y in left)+
                            sum(predict_stump(child_right,v)==y for v,y in right),
            "candidate_split_tests":3*16*33,
            "model_scalars":10}


def predict_stump(node,vec):
    return node["left"] if vec[node["idx"]]<=node["threshold"] else node["right"]


def predict(arm,model,vec):
    if arm=="linear16":
        return 1 if model["bias"]+sum(w*x for w,x in zip(model["weights"],vec))>=0 else -1
    root=model["root"]
    chosen=model["child_left"] if vec[root["idx"]]<=root["threshold"] else model["child_right"]
    return predict_stump(chosen,vec)


def eval_cell(arm,model,condition,generator,site):
    rows=[]
    correct=total=identity_pairs=checks=0
    for seed in HELDOUT:
        tp=tn=fp=fn=same=verified=0
        for entity in range(9):
            base,pos,neg=w37.recoded_example(
                seed,entity,generator,condition,site,3,False)
            origbase,origpos,origneg=w37.recoded_example(
                seed,entity,generator,"identity",site,3,False)
            same+=int(pos==neg)
            for truth,candidate,original in ((True,pos,origpos),(False,neg,origneg)):
                v=feature(base,candidate)
                assert v==feature(origbase,original)
                verified+=1
                predicted=predict(arm,model,v)>0
                if truth and predicted:tp+=1
                elif not truth and not predicted:tn+=1
                elif not truth:fp+=1
                else:fn+=1
        assert tp+tn+fp+fn==18
        correct+=tp+tn
        total+=18
        identity_pairs+=same
        checks+=verified
        rows.append({"seed":seed,"tp":tp,"tn":tn,"fp":fp,"fn":fn,
                     "identical_observation_pairs":same,
                     "xor_feature_checks":verified})
    return {"generator":generator,"site":site,
            "correct":correct,"total":total,"accuracy":correct/total,
            "identical_pairs":identity_pairs,"xor_feature_checks":checks,
            "by_seed":rows}


def run():
    examples,same=development_examples()
    linear=fit_linear(examples)
    tree=fit_tree(examples)
    models={"linear16":linear,"depth2_tree":tree}
    rows=[]
    for arm in ARMS:
        for condition in CONDITIONS:
            cells=[eval_cell(arm,models[arm],condition,g,site)
                   for g in GENERATORS for site in TEST_SITES]
            rows.append({"arm":arm,"condition":condition,
                         "development_example_count":1728,
                         "heldout_seeds":HELDOUT,"feature_width":16,
                         "model_scalars":models[arm]["model_scalars"],
                         "accuracy":sum(c["correct"] for c in cells)/
                                    sum(c["total"] for c in cells),
                         "cells":cells})
    by={(r["arm"],r["condition"]):r["accuracy"] for r in rows}
    assert all(by[a,"identity"]==by[a,"xor55"] for a in ARMS)
    gain=by["depth2_tree","identity"]-by["linear16","identity"]
    if by["depth2_tree","identity"]>=.65 and by["depth2_tree","xor55"]>=.65 and gain>=.10:
        label="NONLINEAR_READOUT_GAIN"
    elif gain>=.05:
        label="PARTIAL_GAIN"
    else:
        label="NO_GAIN"
    return {"schema":"wingless.w-r42-nonlinear-readout.v1",
            "parent_head":"260dc136709f553e8369518605d5010fcfc4f34c",
            "classification":label,"development_seeds":TRAIN,
            "sealed_heldout_seeds":HELDOUT,
            "identical_training_pairs":same,
            "linear_training_updates":linear["updates"],
            "tree_train_correct":tree["train_correct"],
            "tree_candidate_split_tests":tree["candidate_split_tests"],
            "tree_nodes":{k:tree[k] for k in ("root","child_left","child_right")},
            "tree_identity_minus_linear_accuracy":gain,
            "rows":rows,"model_calls":0,
            "claim_limit":"Fixed synthetic lag2 bit-invariant temporal features with readout class and training compute deliberately unequal; no language cognition."}


if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
