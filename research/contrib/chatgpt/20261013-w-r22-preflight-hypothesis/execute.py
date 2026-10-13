"""W-R22: fresh cohort, preflightable finite hypothesis search after INVALID W-R21."""
import hashlib
import json

CANDIDATES=tuple((a,b) for a in (1,2) for b in range(4))

def sequence(seed,e):
    assert 80<=seed<=83 and 0<=e<16
    a=2 if seed%2==0 else 1
    b=(seed-80)%4
    x=[(e+seed)%4,((e//4)+3*seed)%4]
    while len(x)<32: x.append((x[-1]+a*x[-2]+b)%4)
    return tuple(x)

def select(rows):
    return [(a,b) for a,b in CANDIDATES if all(
        (last+a*before+b)%4==label for before,last,label in rows)]

def evaluate(seed,arm,fit,counts):
    cm=[[0]*4 for _ in range(4)]
    trace=hashlib.sha256()
    for e in range(4,16):
        seq=sequence(seed,e)
        for t in range(2,32):
            past=seq[:t]
            if arm=="fixed": guess=(past[-1]+past[-2]+2)%4
            elif arm=="rote":
                row=counts[4*past[-2]+past[-1]]
                guess=max(range(4),key=lambda v:(row[v],-v))
            else:
                if arm=="shuffled":
                    past=past[::2]+past[1::2]
                a,b=fit
                guess=(past[-1]+a*past[-2]+b)%4
            actual=seq[t]
            cm[actual][guess]+=1
            trace.update(bytes((e,t,actual,guess)))
    return {"seed":seed,"arm":arm,"correct":sum(cm[i][i] for i in range(4)),
            "n":360,"confusion":cm,"sha256":trace.hexdigest()}

def run():
    names=("fixed","rote","selected","shuffled")
    rows=[]
    ambiguity=[]
    selected={}
    for seed in range(80,84):
        calibrations=[]
        for e in range(4):
            x=sequence(seed,e)
            calibrations.append((x[0],x[1],x[2]))
        fits=select(calibrations)
        if len(fits)!=1:
            ambiguity.append({"seed":seed,"candidate_count":len(fits)})
            continue
        selected[str(seed)]=fits[0]
        counts=[[0]*4 for _ in range(16)]
        for before,last,label in calibrations:
            counts[4*before+last][label]+=1
        assert sum(map(sum,counts))==4
        for arm in names: rows.append(evaluate(seed,arm,fits[0],counts))
    accuracies={arm:sum(r["correct"] for r in rows if r["arm"]==arm)/1440
                for arm in names} if not ambiguity else {}
    if ambiguity: outcome="HYPOTHESIS_AMBIGUITY"
    elif accuracies["selected"]>=.95 and all(accuracies[k]<=.65 for k in names if k!="selected"):
        outcome="BOUNDED_HYPOTHESIS_TRANSFER"
    elif accuracies["selected"]>max(accuracies[k] for k in names if k!="selected"):
        outcome="PARTIAL"
    else: outcome="NEGATIVE"
    return {"id":"W-R22","classification":outcome,"ambiguity":ambiguity,
            "selected":selected,"accuracy":accuracies,"rows":rows,
            "fresh_seeds":[80,81,82,83],"calibration_labels_per_seed":4,
            "model_calls":0,"prior_W_R21":"INVALID_SOURCE_BUG_NOT_SCIENTIFIC_NEGATIVE",
            "claim_limit":"Eight predesigned hypotheses and four labeled synthetic demonstrations."}

if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
