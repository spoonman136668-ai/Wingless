"""T7 frozen original OOD-distribution source: NO prospective birth in preflight.

New F11/S3 operator *distributions*, same typed group representations.
Truth and source distribution are never passed to T6 sourceblind learner.
A/B cohorts use distinct random-generation implementations.
"""
from __future__ import annotations
import importlib.util,json,random,sys,time,tracemalloc
from pathlib import Path
HERE=Path(__file__).resolve().parent
FROZEN=HERE/"20261018-wingless-t7-preregister.json"
T6=HERE/"20261018-wingless-t6-learning.py"
A=range(-15163,-15099);B=range(-15263,-15199)
TRAIN=range(-14031,-13999)
def need(ok,tag):
    if not ok:raise ValueError("WING_T7_"+tag)
def load(path,name):
    sp=importlib.util.spec_from_file_location(name,path)
    module=importlib.util.module_from_spec(sp);sys.modules[name]=module
    sp.loader.exec_module(module);return module
def born(seed,k,mode="prospective"):
    if mode=="fixture":need(seed==-15000,"ONLY_FIXTURE")
    else:need(seed in A or seed in B,"FIRST_USE_IDENTITY")
    distinct=seed in B
    rng=random.Random((314159265 ^ (abs(seed)*1103515245)) if distinct
                      else 981212+abs(seed)*90217)
    domain=k.AFFINE if abs(seed)%2==0 else k.NONCOMMUTING_PERM
    roles=[];truth=[]
    for j in range(4):
        if domain==k.AFFINE:
            if rng.random()<0.75:
                a=(7,9)[rng.randrange(2)]
            else:a=rng.randrange(1,11)
            op=(a,rng.randrange(11))
            x=rng.randrange(11)
        else:
            odd=((1,0,2),(2,1,0),(0,2,1))
            op=odd[rng.randrange(3)] if rng.random()<0.75 else k.PERMS[rng.randrange(6)]
            x=rng.randrange(3)
        roles.append(((x,k.operator_apply(domain,op,x)),))
        truth.append(op)
    chain=tuple((rng.randrange(4),bool(rng.randrange(2)))
                for _ in range(10+(abs(seed)%2)*2))
    start=rng.randrange(11 if domain==0 else 3)
    state=k.start(domain)
    for role,inverse in chain:
        state=k.compose_state(state,truth[role],domain,inverse)
    return (domain,tuple(roles),chain,start),k.apply(state,start)
def measured(fn):
    tracemalloc.start();t=time.process_time_ns()
    result=fn()
    d={"cpu_ns":time.process_time_ns()-t,
       "python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    tracemalloc.stop();return result,d
def prereg():
    f=json.loads(FROZEN.read_text())
    need(f["schema"]=="northstar.t7.wingless.preimplementation-prospective.v1"
         and f["source_parent_head"]=="f3da303e06e7c4d56e1c6c27db0f32230c110030",
         "FREEZE_EXACT_PARENT")
    t6=load(T6,"wing_t7_t6")
    t5=t6.__import__("importlib").util.spec_from_file_location("unused",T6) if False else load(t6.T5,"wing_t7_t5")
    k=t5.load(t5.KERNEL,"wing_t7_typed")
    case,answer=born(-15000,k,"fixture")
    need(len(case[1])==4 and all(len(v)==1 for v in case[1])
         and len(case[2])==10 and 0<=answer<(11 if case[0]==0 else 3),
         "NEW_DISTRIBUTION_SOURCE_SHAPE")
    for forbidden in (-15100,-15200,-13000,882000):
        try:born(forbidden,k,"fixture")
        except ValueError:pass
        else:raise AssertionError("prospective or legacy source opened")
    return {"schema":"wingless.t7.nonconsuming.source-only.v1",
       "T7_prospective_births":0,"original_WNS_used":0,
       "new_distribution_operators":"F11 slopes 7/9 and S3 odd transpositions",
       "source_truth_not_model_input":True,"preflight_pass":True}
def cohort(t6,t5,k,prior,seeds):
    cases=[born(i,k) for i in seeds]
    names=["direct","exact_public"]
    predictions={mode:[t6.predict(k,case,prior,mode,t5)
                       for case,_ in cases] for mode in names}
    expected=[truth for _,truth in cases]
    correct=lambda p:sum(x==y for x,y in zip(p,expected))
    reg=sum(a==y and b!=y for a,b,y in zip(
            predictions["exact_public"],predictions["direct"],expected))
    return {"tasks":len(cases),"learned_correct":correct(predictions["direct"]),
       "same_information_control_correct":correct(predictions["exact_public"]),
       "baseline_correct_regressions":reg,
       "learned_changes":sum(a!=b for a,b in zip(predictions["direct"],
                                               predictions["exact_public"]))}
def original():
    prereg()
    t6=load(T6,"wing_t7_t6orig")
    t5=load(t6.T5,"wing_t7_t5orig")
    k=t5.load(t5.KERNEL,"wing_t7_kernelorig")
    # Previously public TRAIN only, no fresh training schedule or threshold
    # selected after reviewing T6's consumed OPEN.
    original_train=[t6.fresh_source(t5,i,k) for i in t6.TRAIN]
    (prior,_),traincost=measured(lambda:t5.fit_prior(k,original_train))
    first,diagnostic=measured(lambda:cohort(t6,t5,k,prior,A))
    # Independent replica is first-used only on admission of initial A
    # nonregression and accuracy. Cost must also pass, checked separately.
    def compared(seeds):
        cases=[born(i,k) for i in seeds]
        candidate,cc=measured(lambda:[t6.predict(k,c,prior,"direct",t5) for c,_ in cases])
        control,bc=measured(lambda:[t6.predict(k,c,prior,"exact_public",t5) for c,_ in cases])
        return cc,bc
    candidate,control=compared(A)
    cost_ok=(candidate["cpu_ns"]+traincost["cpu_ns"]/64<=control["cpu_ns"]
             and max(candidate["python_peak_bytes"],traincost["python_peak_bytes"])
                     <=control["python_peak_bytes"])
    eligible=(first["learned_correct"]>first["same_information_control_correct"]
              and first["baseline_correct_regressions"]==0 and cost_ok)
    replica=None
    if eligible:
        replica=cohort(t6,t5,k,prior,B)
        # Prospective B must preserve original nonregression and gain.
    return {"schema":"wingless.t7.first-original-new-distribution.v1",
        "TRAIN_prior_reused_unchanged":True,"TRAIN_count":len(original_train),
        "source_parent_T6_unchanged":True,
        "A":first,"A_candidate_runtime":candidate,"A_control_runtime":control,
        "original_amortized_train_runtime":traincost,
        "A_cost_parity":cost_ok,"A_eligible":eligible,
        "B_replica":replica,
        "B_constructed":replica is not None,
        "A_prospective_sources_constructed":len(A),
        "B_prospective_sources_constructed":len(B) if replica else 0,
        "WNS_protected_source_constructed":0,
        "scientific_OOD_acceptance":False,
        "original_science_promotion":False}
if __name__=="__main__":
    import argparse
    p=argparse.ArgumentParser();p.add_argument("--mode",choices=("preflight","original"),default="preflight")
    args=p.parse_args()
    print(json.dumps(prereg() if args.mode=="preflight" else original(),sort_keys=True))
