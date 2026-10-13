"""W-R35: frozen comparison of 1/2/3-error triggers under noisy-feedback bursts."""
import hashlib, importlib.util, json
from pathlib import Path

prior=Path(__file__).resolve().parent.parent/"20261013-w-r26-latent-alphabet"/"execute.py"
spec=importlib.util.spec_from_file_location("w26",prior)
w26=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)
assert w26.fixture()=="PASS"
ARMS=("FROZEN_ADD","K1","K2","K3")
PROFILES=("ISOLATED","BURST_2")


def sequence(seed,ep):
    assert 216<=seed<=223 and 0<=ep<16
    order=sorted(range(16),key=lambda i:(
        hashlib.sha256(f"W35-order/{seed}/{i}".encode()).digest(),i))
    o=order[ep]
    xs=[(o+seed)%4,((o//4)+3*seed)%4]
    while len(xs)<32:
        t=len(xs);a,b=xs[-2:]
        xs.append((a^b) if seed%2==0 and t>=16 else (a+b)%4)
    p=w26.PERMUTATIONS[hashlib.sha256(
        f"W35-perm/{seed}".encode()).digest()[0]%24]
    return tuple(p[x] for x in xs)


def xor_table(p):
    inv=[0]*4
    for i,label in enumerate(p):inv[label]=i
    return tuple(p[inv[a]^inv[b]] for a,b in w26.ALL_CONTEXTS)


def calibration(seed):
    labels=[sequence(seed,e)[:3] for e in range(16)]
    if len({(a,b) for a,b,_ in labels})!=16:return {"invalid":"MISSING_CONTEXT"}
    fits=[p for p in w26.PERMUTATIONS if w26.fits(p,labels)]
    aa={w26.table_for(p) for p in fits}
    xx={xor_table(p) for p in fits}
    if len(aa)!=1 or len(xx)!=1:
        return {"invalid":"AMBIGUOUS_HYPOTHESIS","add":len(aa),"xor":len(xx)}
    a,x=next(iter(aa)),next(iter(xx))
    return {"a":a,"x":x,"candidate_count":len(fits),
            "add_sha256":hashlib.sha256(bytes(a)).hexdigest(),
            "xor_sha256":hashlib.sha256(bytes(x)).hexdigest()}


def evaluate(seed,profile,arm,a,x):
    ticks=(7,24) if profile=="ISOLATED" else (7,8,24,25)
    metrics={"pre":{"correct":0,"wrong":0},"post":{"correct":0,"wrong":0}}
    digest=hashlib.sha256()
    switches=0;corrupt=0
    for ep in range(16):
        xs=sequence(seed,ep);mode=0;streak=0
        for t in range(3,32):
            pred=(x if mode else a)[4*xs[t-2]+xs[t-1]]
            actual=xs[t]
            score=metrics["pre" if t<16 else "post"]
            score["correct" if pred==actual else "wrong"]+=1
            is_corrupt=t in ticks
            report=(actual+1)%4 if is_corrupt else actual
            corrupt+=int(is_corrupt)
            if arm!="FROZEN_ADD":
                mismatch=report!=pred
                streak=streak+1 if mismatch else 0
                need={"K1":1,"K2":2,"K3":3}[arm]
                if streak==need:
                    mode=1-mode;switches+=1;streak=0
            digest.update(bytes((ep,t,actual,pred,report,mode)))
    assert corrupt==16*len(ticks)
    assert sum(metrics["pre"].values())==208
    assert sum(metrics["post"].values())==256
    return {"seed":seed,"cohort":"development" if seed<220 else "heldout",
            "condition":"SHIFT" if seed%2==0 else "STABLE",
            "noise_profile":profile,"arm":arm,"metrics":metrics,
            "switches":switches,"feedback_corruptions":corrupt,
            "trace_sha256":digest.hexdigest()}


def run():
    fingerprints={str(s):hashlib.sha256(bytes(
        val for e in range(16) for val in sequence(s,e))).hexdigest()
        for s in range(216,224)}
    if len(set(fingerprints.values()))!=8:
        return {"classification":"INVALID_IDENTITIES","fingerprints":fingerprints,
                "rows":[],"model_calls":0}
    rows=[];fits={}
    for seed in range(216,224):
        f=calibration(seed)
        if "invalid" in f:
            return {"classification":"AMBIGUOUS_HYPOTHESIS","seed":seed,
                    "reason":f,"rows":[],"model_calls":0}
        a,x=f.pop("a"),f.pop("x")
        fits[str(seed)]=f
        for profile in PROFILES:
            for arm in ARMS:rows.append(evaluate(seed,profile,arm,a,x))
    held=[r for r in rows if r["seed"]>=220]
    agg={p:{c:{a:{
        "pre":{k:sum(r["metrics"]["pre"][k] for r in held
                        if (r["noise_profile"],r["condition"],r["arm"])==(p,c,a))
               for k in ("correct","wrong")},
        "post":{k:sum(r["metrics"]["post"][k] for r in held
                         if (r["noise_profile"],r["condition"],r["arm"])==(p,c,a))
                for k in ("correct","wrong")},
        "switches":sum(r["switches"] for r in held
                       if (r["noise_profile"],r["condition"],r["arm"])==(p,c,a))}
        for a in ARMS} for c in ("STABLE","SHIFT")} for p in PROFILES}
    burst=agg["BURST_2"]["STABLE"]
    shift=agg["ISOLATED"]["SHIFT"]
    classification=("BURST_CORROBORATION_CEILING"
        if (burst["K3"]["post"]["wrong"]<burst["K2"]["post"]["wrong"]
            and burst["K3"]["switches"]<burst["K2"]["switches"]
            and shift["K2"]["post"]["correct"]>=shift["K3"]["post"]["correct"])
        else "MIXED")
    return {"schema":"wingless.w-r35-bursty-noise.v1",
            "classification":classification,
            "parent":"4eb34e2ea945c1df98eb4ca79c2c9db6b4cb2930",
            "fingerprints":fingerprints,"calibrations":fits,
            "heldout":agg,"rows":rows,"model_calls":0,
            "claim_limit":"Finite ADD/XOR source with separate perfect state observation and corrupted supervisory reports."}


if __name__=="__main__":
    result=run()
    assert result==run()
    print(json.dumps(result,sort_keys=True,indent=2))
