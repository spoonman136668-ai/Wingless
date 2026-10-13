"""W-R34: causal model toggling with precommitted corrupted supervisor reports."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS = Path(__file__).resolve().parent.parent / "20261013-w-r26-latent-alphabet" / "execute.py"
spec = importlib.util.spec_from_file_location("w26_pinned", PREVIOUS)
assert spec and spec.loader
w26 = importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)
assert w26.fixture() == "PASS" and len(w26.PERMUTATIONS) == 24
ARMS = ("FROZEN_ADD", "ONE_ERROR_TOGGLE", "TWO_CONSECUTIVE_ERRORS")


def order(seed):
    assert 208 <= seed <= 215
    return tuple(sorted(range(16), key=lambda i: (
        hashlib.sha256(("W34-order/"+str(seed)+"/"+str(i)).encode()).digest(),i)))


def sequence(seed, episode):
    original = order(seed)[episode]
    xs = [(original+seed)%4,((original//4)+3*seed)%4]
    while len(xs) < 32:
        t=len(xs)
        before,previous=xs[-2:]
        new = (before ^ previous) if seed%2==0 and t>=16 else (before+previous)%4
        xs.append(new)
    p = w26.PERMUTATIONS[hashlib.sha256(
        ("W34-perm/"+str(seed)).encode()).digest()[0]%24]
    return tuple(p[k] for k in xs)


def xor_table(p):
    inverse=[0]*4
    for latent,observed in enumerate(p):
        inverse[observed]=latent
    return tuple(p[inverse[a]^inverse[b]] for a,b in w26.ALL_CONTEXTS)


def calibrate(seed):
    cases=[sequence(seed,e)[:3] for e in range(16)]
    if len({(a,b) for a,b,_ in cases})!=16:
        return {"invalid":"MISSING_CONTEXTS"}
    p=[perm for perm in w26.PERMUTATIONS if w26.fits(perm,cases)]
    additions={w26.table_for(perm) for perm in p}
    xors={xor_table(perm) for perm in p}
    if len(additions)!=1 or len(xors)!=1:
        return {"invalid":"AMBIGUOUS_HYPOTHESIS","fits":len(p),
                "add_tables":len(additions),"xor_tables":len(xors)}
    add,xor=next(iter(additions)),next(iter(xors))
    return {"add":add,"xor":xor,"candidates":len(p),"calibration_labels":16,
            "add_sha256":hashlib.sha256(bytes(add)).hexdigest(),
            "xor_sha256":hashlib.sha256(bytes(xor)).hexdigest()}


def evaluate(seed,arm,add,xor):
    before={"correct":0,"wrong":0}
    after={"correct":0,"wrong":0}
    trace=hashlib.sha256()
    switches=[]
    corrupt_count=0
    for episode in range(16):
        xs=sequence(seed,episode)
        mode=0
        streak=0
        switched=[]
        for t in range(3,32):
            # A separate reliable state channel supplies the true past context.
            # The label report below is a DISTINCT noisy supervision channel.
            model = xor if mode else add
            pred=model[4*xs[t-2]+xs[t-1]]
            truth=xs[t]  # inaccessible to learner; evaluator-only scoring
            bucket = before if t<16 else after
            bucket["correct" if pred==truth else "wrong"]+=1
            noisy=(t==7 or t==24)
            reported=(truth+1)%4 if noisy else truth
            corrupt_count+=int(noisy)
            disagreement=reported!=pred
            if arm=="ONE_ERROR_TOGGLE" and disagreement:
                mode=1-mode
                switched.append(t)
            elif arm=="TWO_CONSECUTIVE_ERRORS":
                streak=streak+1 if disagreement else 0
                if streak==2:
                    mode=1-mode
                    switched.append(t)
                    streak=0
            trace.update(bytes((episode,t,truth,pred,reported,mode)))
        switches.append(switched)
    assert corrupt_count==32
    assert before["correct"]+before["wrong"]==16*13
    assert after["correct"]+after["wrong"]==16*16
    return {"seed":seed,"cohort":"development" if seed<212 else "heldout",
            "condition":"SHIFT" if seed%2==0 else "STABLE","arm":arm,
            "pre":before,"post":after,"noisy_supervision_reports":corrupt_count,
            "switch_count":sum(map(len,switches)),
            "switch_events":switches,"trace_sha256":trace.hexdigest()}


def run():
    fingerprints={str(seed):hashlib.sha256(bytes(
        token for episode in range(16) for token in sequence(seed,episode)
    )).hexdigest() for seed in range(208,216)}
    if len(set(fingerprints.values()))!=8:
        return {"schema":"wingless.w-r34-noisy-feedback.v1",
                "classification":"INVALID_IDENTITIES","fingerprints":fingerprints,
                "rows":[],"model_calls":0}
    calibrations={}
    rows=[]
    for seed in range(208,216):
        fit=calibrate(seed)
        if "invalid" in fit:
            return {"schema":"wingless.w-r34-noisy-feedback.v1",
                    "classification":"AMBIGUOUS_HYPOTHESIS","seed":seed,
                    "details":fit,"rows":[],"model_calls":0}
        add,xor=fit.pop("add"),fit.pop("xor")
        calibrations[str(seed)]=fit
        rows.extend(evaluate(seed,arm,add,xor) for arm in ARMS)
    heldout=[r for r in rows if r["seed"]>=212]
    summary={condition:{arm:{
        "pre":{key:sum(r["pre"][key] for r in heldout
                        if r["condition"]==condition and r["arm"]==arm)
               for key in ("correct","wrong")},
        "post":{key:sum(r["post"][key] for r in heldout
                         if r["condition"]==condition and r["arm"]==arm)
                for key in ("correct","wrong")},
        "switches":sum(r["switch_count"] for r in heldout
                       if r["condition"]==condition and r["arm"]==arm)}
        for arm in ARMS} for condition in ("STABLE","SHIFT")}
    stable=summary["STABLE"]
    shift=summary["SHIFT"]
    tradeoff=(stable["TWO_CONSECUTIVE_ERRORS"]["post"]["wrong"]<
              stable["ONE_ERROR_TOGGLE"]["post"]["wrong"]
              and stable["TWO_CONSECUTIVE_ERRORS"]["switches"]<
              stable["ONE_ERROR_TOGGLE"]["switches"]
              and shift["TWO_CONSECUTIVE_ERRORS"]["post"]["correct"]<=
              shift["ONE_ERROR_TOGGLE"]["post"]["correct"])
    return {"schema":"wingless.w-r34-noisy-feedback.v1",
            "classification":"NOISE_HYSTERESIS_TRADEOFF" if tradeoff else "MIXED",
            "parent":"18742dba3d6a2266d592bb740b1df17a0ca5482f",
            "source_fingerprints":fingerprints,"calibrations":calibrations,
            "heldout":summary,"rows":rows,"model_calls":0,
            "claim_limit":"Two frozen algebraic models and separate perfect context/noisy-label channels; no unrestricted noise robustness."}


if __name__=="__main__":
    original=run()
    assert original==run()
    print(json.dumps(original,indent=2,sort_keys=True))
