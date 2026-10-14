"""W68 development-only feature identity test; no heldout first-use."""
import importlib.util,json
from pathlib import Path
p=Path(__file__).resolve().parent/"execute.py"
spec=importlib.util.spec_from_file_location("w67_frozen",p)
assert spec and spec.loader
w67=importlib.util.module_from_spec(spec);spec.loader.exec_module(w67)
m=w67.w66.w65.w64.w63.w61.m.m
reports={"cases":0,"minus1_positive_equal":0,"minus1_negative_equal":0,
         "plus1_positive_equal":0,"plus1_negative_equal":0,
         "minus1_positive_difference_l1":0.0,"plus1_positive_difference_l1":0.0}
for seed in range(9008,9016):
    for entity in (9,10,11):
        for site in w67.w66.w65.w64.SITES:
            base,pos,neg=m.w37.recoded_example(seed,entity,"xorshift_independent","identity",site,1,False)
            observations={}
            for phase in ("intact","minus1","plus1"):
                x0,xp,xn=(w67.acquire(x,(0,1,2,3),phase) for x in (base,pos,neg))
                observations[phase]=(m.w42.feature(x0,xp),m.w42.feature(x0,xn))
            p0,n0=observations["intact"]
            pm,nm=observations["minus1"]
            pp,np=observations["plus1"]
            assert len(p0)==len(pm)==len(pp)==16
            reports["cases"]+=1
            reports["minus1_positive_equal"]+=int(p0==pm)
            reports["minus1_negative_equal"]+=int(n0==nm)
            reports["plus1_positive_equal"]+=int(p0==pp)
            reports["plus1_negative_equal"]+=int(n0==np)
            reports["minus1_positive_difference_l1"]+=sum(abs(a-b) for a,b in zip(p0,pm))
            reports["plus1_positive_difference_l1"]+=sum(abs(a-b) for a,b in zip(p0,pp))
assert reports["cases"]>0
print(json.dumps({"classification":"DEVELOPMENT_ONLY_FEATURE_CAUSALITY","no_sealed_replay":True,"report":reports},sort_keys=True))
