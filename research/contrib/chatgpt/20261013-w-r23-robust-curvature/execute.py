"""W-R23: bounded robust online curvature under unseen post-tick8 regime shift."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r22-online-prediction"/"execute.py"
spec=importlib.util.spec_from_file_location("w22_frozen",PARENT)
assert spec and spec.loader
w22=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w22)

ARMS=("momentum","frozen","naive_online","median3_online")


def learn():
    differences=[]
    for seed in range(88,92):
        for entity in range(9):
            raw=w22.series(seed,entity,"development")
            for t in range(2,16):
                differences.append(w22.circular_signed(raw[t]-2*raw[t-1]+raw[t-2]))
    differences.sort()
    assert len(differences)==504
    count=len(differences)
    return (differences[count//2-1]+differences[count//2])/2


def median(items):
    ordered=sorted(items)
    n=len(ordered)
    assert 1<=n<=3
    if n%2:return ordered[n//2]
    return (ordered[n//2-1]+ordered[n//2])/2


def predict(raw,t,arm,frozen):
    momentum=2*raw[t-1]-raw[t-2]
    if arm=="momentum":return momentum
    if arm=="frozen":return momentum+frozen
    if arm=="naive_online":
        curvature=(frozen if t<3 else
                   w22.circular_signed(raw[t-1]-2*raw[t-2]+raw[t-3]))
        return momentum+curvature
    assert arm=="median3_online"
    curvatures=[
        w22.circular_signed(raw[j]-2*raw[j-1]+raw[j-2])
        for j in range(max(2,t-3),t)
    ]
    return momentum+(median(curvatures) if curvatures else frozen)


def evaluate(arm,frozen):
    by_seed=[]
    all_errors=[]
    all_pre=[]
    all_post=[]
    for seed in range(92,96):
        rows=[]
        for entity in range(9):
            raw=w22.series(seed,entity,"heldout")
            for t in range(2,16):
                prediction=predict(raw,t,arm,frozen)
                error=w22.circular_abs(prediction,raw[t])
                rows.append((t,error))
        assert len(rows)==126
        pre=[err for t,err in rows if t<=8]
        post=[err for t,err in rows if t>=9]
        assert len(pre)==len(post)==63
        errors=[error for t,error in rows]
        by_seed.append({"seed":seed,"full_mae":sum(errors)/len(errors),
                        "prechange_mae":sum(pre)/len(pre),
                        "postchange_mae":sum(post)/len(post)})
        all_errors.extend(errors)
        all_pre.extend(pre)
        all_post.extend(post)
    return {"arm":arm,"n":len(all_errors),
            "full_mae":sum(all_errors)/len(all_errors),
            "prechange_mae":sum(all_pre)/len(all_pre),
            "postchange_mae":sum(all_post)/len(all_post),
            "by_seed":by_seed}


def run():
    frozen=learn()
    results=[evaluate(arm,frozen) for arm in ARMS]
    r={entry["arm"]:entry for entry in results}
    robust=r["median3_online"]["postchange_mae"]
    static=r["frozen"]["postchange_mae"]
    naive=r["naive_online"]["postchange_mae"]
    by_seed_benefit=sum(a["postchange_mae"]<b["postchange_mae"]
                        for a,b in zip(r["median3_online"]["by_seed"],
                                       r["frozen"]["by_seed"]))
    if robust<=.90*static and robust<=naive and by_seed_benefit>=3:
        classification="ROBUST_ONLINE_TRANSFER"
    elif robust<static:
        classification="PARTIAL"
    else:
        classification="NEGATIVE"
    return {
        "schema":"wingless.w-r23-median-online.v1",
        "predecessor_head":"25d608c3e3f4b18e3922e82741048ad53d57190c",
        "classification":classification,
        "development_curvature":frozen,
        "heldout_seeds_improved":by_seed_benefit,
        "models":results,
        "model_calls":0,
        "scope":"Synthetic causal fixed-five-tick context next-step prediction"
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
