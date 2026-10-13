"""W-R53: frozen 16-D ridge normal equations versus full-batch paired perceptron."""
import hashlib
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r52-batch-diversity"/"execute.py"
spec=importlib.util.spec_from_file_location("w52_pinned",SOURCE)
assert spec and spec.loader
w52=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w52)

w37=w52.w37
w42=w52.w42
w45=w52.w45
GENERATORS=w52.GEN
TEST_SITES=w52.SITES
DEV=(tuple(range(3680,3696)),tuple(range(3712,3728)),tuple(range(3744,3760)))
HELDOUT=tuple(range(3808,3856))
ENCODINGS=("identity","xor55")
ARMS=("affine_batch","affine_ridge","mixed_batch","mixed_ridge")
LAMBDA=1.0

def solve(matrix,vector):
    n=len(vector)
    a=[list(row)+[vector[i]] for i,row in enumerate(matrix)]
    for c in range(n):
        pivot=max(range(c,n),key=lambda i:abs(a[i][c]))
        assert abs(a[pivot][c])>1e-12,"SINGULAR_RIDGE_MATRIX"
        a[pivot],a[c]=a[c],a[pivot]
        v=a[c][c]
        for k in range(c,n+1):a[c][k]/=v
        for i in range(n):
            if i==c:continue
            factor=a[i][c]
            for k in range(c,n+1):a[i][k]-=factor*a[c][k]
    x=[row[-1] for row in a]
    assert max(abs(sum(matrix[i][j]*x[j] for j in range(n))-vector[i])
               for i in range(n))<=1e-7,"RIDGE_NUMERIC_RESIDUAL_FAILED"
    return x

def ridge_fit(examples):
    n=len(examples)
    gram=[[0]*16 for _ in range(16)]
    rhs=[0]*16
    for vec,label in examples:
        for i in range(16):
            rhs[i]+=label*vec[i]
            for j in range(16):
                gram[i][j]+=vec[i]*vec[j]
    # Data matrix is exact integer accumulation; regularizer is positive.
    system=[[gram[i][j]/n + (LAMBDA if i==j else 0.0)
             for j in range(16)] for i in range(16)]
    target=[v/n for v in rhs]
    weights=solve(system,target)
    return weights

def batch_fit(examples):
    weights=[0.0]*16
    updates=0
    for _ in range(32):
        delta=[0]*16
        for x,y in examples:
            if y*sum(w*v for w,v in zip(weights,x))<=0:
                for j in range(16):delta[j]+=y*x[j]
                updates+=1
        weights=[w+.1*d/len(examples) for w,d in zip(weights,delta)]
    return weights,updates

def train():
    models={}
    for index,block in enumerate(DEV,1):
        for source in ("affine","mixed"):
            samples,counts,fingerprint=w52.training_examples(
                block,"affine_only_full_batch" if source=="affine"
                else "three_source_balanced_full_batch")
            assert len(samples)==1728
            bw,updates=batch_fit(samples)
            rw=ridge_fit(samples)
            for solver,weights in (("batch",bw),("ridge",rw)):
                models[f"{index}:{source}_{solver}"]={
                    "weights":weights,
                    "weights_sha256":hashlib.sha256(json.dumps(
                        weights,separators=(",",":")).encode()).hexdigest(),
                    "development_block":index,"source_distribution":source,
                    "solver":solver,"training_examples":1728,
                    "training_corpus_sha256":fingerprint,
                    "source_counts":counts,"batch_updates":updates if solver=="batch" else None}
    assert len(models)==12
    return models

def run():
    models=train()
    results={}
    for encoding in ENCODINGS:
        per={}
        for key,model in models.items():
            families={g:{"correct":0.0,"presentations":0} for g in GENERATORS}
            by_seed=[]
            for seed in HELDOUT:
                own={"seed":seed,"correct":0.0,"presentations":0,
                     "ties":0,"source_correct":{}}
                for generator in GENERATORS:
                    correct=0.0
                    for site in TEST_SITES:
                        for entity in range(9):
                            base,pos,neg=w37.recoded_example(
                                seed,entity,generator,encoding,site,3,False)
                            rawbase,rawpos,rawneg=w37.recoded_example(
                                seed,entity,generator,"identity",site,3,False)
                            a=w42.feature(base,pos);b=w42.feature(base,neg)
                            assert a==w42.feature(rawbase,rawpos)
                            assert b==w42.feature(rawbase,rawneg)
                            margin=sum(w*(x-y) for w,x,y in zip(model["weights"],a,b))
                            s=w45.credit(margin,True)+w45.credit(-margin,False)
                            correct+=s
                            own["correct"]+=s
                            own["presentations"]+=2
                            own["ties"]+=2*int(abs(margin)<=1e-10)
                            families[generator]["correct"]+=s
                            families[generator]["presentations"]+=2
                    own["source_correct"][generator]=correct
                assert own["presentations"]==108
                by_seed.append(own)
            per[key]={"accuracy":sum(x["correct"] for x in by_seed)/(48*108),
                      "source_accuracy":{g:v["correct"]/v["presentations"] for g,v in families.items()},
                      "weight_sha256":model["weights_sha256"],
                      "training_corpus_sha256":model["training_corpus_sha256"],
                      "source_training_example_counts":model["source_counts"],
                      "batch_updates":model["batch_updates"],
                      "by_seed":by_seed}
        results[encoding]=per
    assert results["identity"].keys()==results["xor55"].keys()
    for key in results["identity"]:
        assert results["identity"][key]["accuracy"]==results["xor55"][key]["accuracy"]
        assert results["identity"][key]["source_accuracy"]==results["xor55"][key]["source_accuracy"]
    comparisons=[]
    for i in (1,2,3):
        batch=results["identity"][f"{i}:mixed_batch"]
        ridge=results["identity"][f"{i}:mixed_ridge"]
        floor_batch=min(batch["source_accuracy"].values())
        floor_ridge=min(ridge["source_accuracy"].values())
        deltas={g:ridge["source_accuracy"][g]-batch["source_accuracy"][g] for g in GENERATORS}
        gain=floor_ridge-floor_batch
        comparisons.append({"dev_cohort":i,"mixed_batch_accuracy":batch["accuracy"],
                            "mixed_ridge_accuracy":ridge["accuracy"],
                            "worst_source_gain":gain,
                            "source_deltas":deltas,
                            "strict":gain>=.03 and ridge["accuracy"]>=.55 and min(deltas.values())>=-.05,
                            "partial":gain>=.02})
    if all(c["strict"] for c in comparisons):
        label="RIDGE_WORST_FAMILY_GAIN"
    elif sum(c["partial"] for c in comparisons)>=2:
        label="PARTIAL_RIDGE_GAIN"
    else:label="NO_GAIN"
    return {"schema":"wingless.w-r53-ridge-relative-pair.v1",
            "parent_head":"b1e4972f04142f1f35f76f9483b12491dec4846c",
            "classification":label,"development_blocks":DEV,
            "sealed_heldout_seeds":HELDOUT,"lambda_normalized":LAMBDA,
            "models":results,"mixed_source_comparisons":comparisons,
            "models_trained_on_new_sealed":False,
            "model_calls":0,
            "claim_limit":"Synthetic paired-comparison ridge readout versus 32-epoch batch perceptron, with unequal computational objectives/cost; no autonomous cognition."}

if __name__=="__main__":
    print(json.dumps(run(),indent=2,sort_keys=True))
