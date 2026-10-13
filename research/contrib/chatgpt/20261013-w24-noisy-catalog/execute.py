"""W-R24: single corrupted calibration label and out-of-catalogue rejection."""
import importlib.util
import json
from pathlib import Path

PRIOR=Path(__file__).resolve().parent.parent/"20261013-w23-catalog-abstention"/"execute.py"
_spec=importlib.util.spec_from_file_location("w23_frozen",PRIOR)
assert _spec and _spec.loader
w23=importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(w23)


def family(seed):
    assert 96<=seed<=107
    if seed<=103:
        return (2 if seed%2==0 else 1, (seed-(96 if seed<100 else 100))%4)
    return (3,(seed-104)%4)


def population(seed):
    if seed<100:return "CLEAN"
    if seed<104:return "ONE_CORRUPTED_LABEL"
    return "OUT_OF_CATALOGUE"


def episode(seed,e):
    assert 0<=e<16
    a,b=family(seed)
    seq=[(e+seed)%4,((e//4)+3*seed)%4]
    while len(seq)<32:
        seq.append((seq[-1]+a*seq[-2]+b)%4)
    return tuple(seq)


def observed_calibration(seed):
    rows=[]
    for ep in range(4):
        seq=episode(seed,ep)
        observed=seq[2]
        if population(seed)=="ONE_CORRUPTED_LABEL" and ep==seed%4:
            observed=(observed+1)%4
        rows.append((seq[0],seq[1],observed))
    return rows


def select_models(rows):
    rankings=[]
    for a,b in w23.w22.CANDIDATES:
        mistakes=sum((last+a*before+b)%4!=out
                     for before,last,out in rows)
        rankings.append((mistakes,a,b))
    rankings.sort()
    best=rankings[0]
    tied=sum(item[0]==best[0] for item in rankings)
    strict=((best[1],best[2]) if best[0]==0 and tied==1 else None)
    tolerant=((best[1],best[2]) if best[0]<=1 and tied==1 else None)
    forced=(best[1],best[2])
    return strict,tolerant,forced,best[0],tied,rankings[1][0]


def accuracy(seed,model):
    if model is None:return None
    a,b=model
    matched=0
    for ep in range(4,16):
        seq=episode(seed,ep)
        for tick in range(2,32):
            matched+=int((seq[tick-1]+a*seq[tick-2]+b)%4==seq[tick])
    return matched/360


def run():
    rows=[]
    for seed in range(96,108):
        observations=observed_calibration(seed)
        strict,tolerant,forced,errors,ties,second=select_models(observations)
        rows.append({
            "seed":seed,"cohort":population(seed),
            "calibration_labels":4,
            "calibration_corruptions":int(population(seed)=="ONE_CORRUPTED_LABEL"),
            "min_errors":errors,"best_ties":ties,"runnerup_errors":second,
            "strict_accepted":strict is not None,
            "tolerant_accepted":tolerant is not None,
            "strict_accuracy":accuracy(seed,strict),
            "tolerant_accuracy":accuracy(seed,tolerant),
            "forced_accuracy":accuracy(seed,forced),
        })
    clean=[r for r in rows if r["cohort"]=="CLEAN"]
    noisy=[r for r in rows if r["cohort"]=="ONE_CORRUPTED_LABEL"]
    out=[r for r in rows if r["cohort"]=="OUT_OF_CATALOGUE"]
    clean_ok=sum(r["tolerant_accepted"] and r["tolerant_accuracy"]>=.95 for r in clean)
    noisy_ok=sum(r["tolerant_accepted"] and r["tolerant_accuracy"]>=.95 for r in noisy)
    out_rej=sum(not r["tolerant_accepted"] for r in out)
    if clean_ok==4 and noisy_ok==4 and out_rej==4:
        label="ROBUST_CATALOG_SELECTIVITY"
    elif noisy_ok>0 and out_rej>0:
        label="PARTIAL"
    else:
        label="FAILURE"
    return {
        "schema":"wingless.w-r24-noisy-catalog-selection.v1",
        "predecessor_head":"51304a86f4291036e87f22716397881f24ea652e",
        "classification":label,
        "catalog_count":len(w23.w22.CANDIDATES),
        "clean_accepted_correct":clean_ok,
        "noisy_accepted_correct":noisy_ok,
        "out_of_catalogue_rejected":out_rej,
        "strict_noisy_rejections":sum(not r["strict_accepted"] for r in noisy),
        "rows":rows,
        "model_calls":0,
        "claim_limit":"Exactly one synthetic mislabeled exemplar and known eight-template family; no generic noise robustness."
    }


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
