"""T5 Wingless: learn actual shared public operator priors under ambiguity.

The public TRAIN roles have 3 samples; OPEN roles reveal only one public
calibration pair. Role transformations and answer are source evaluator-only.
No original WNS TRAIN/SCREEN/UNSEEN/replica identity can be instantiated.
All public algebra descriptors and observations match between controls.
"""
from __future__ import annotations
import importlib.util,sys,json,random,time,tracemalloc
from pathlib import Path
HERE=Path(__file__).resolve().parent
FREEZE=HERE/"20261018-wingless-t5-frozen.json"
KERNEL=HERE/"20261018-wingless-t3-typed-composition.py"
TRAIN=range(-13031,-12999)
OPEN=range(-13163,-13099)
def need(v,n):
    if not v:raise ValueError("WING_T5_"+n)
def load(path,name):
    sp=importlib.util.spec_from_file_location(name,path)
    m=importlib.util.module_from_spec(sp)
    sys.modules[sp.name]=m;sp.loader.exec_module(m)
    return m
def source(seed,k):
    need(seed in TRAIN or seed in OPEN,"NOT_ORIGINAL_PROTECTED_SOURCE")
    rng=random.Random(674912 + abs(seed)*104729)
    domain=k.AFFINE if abs(seed)%2==0 else k.NONCOMMUTING_PERM
    roles=[];truth=[]
    for _ in range(4):
        if domain==k.AFFINE:
            slope=2 if rng.random()<0.75 else rng.randrange(1,11)
            op=(slope,rng.randrange(11))
            xs=rng.sample(range(11),3)
        else:
            op=rng.choice([(0,1,2),(1,2,0),(2,0,1)]) if rng.random()<0.75 else rng.choice(k.PERMS)
            xs=[0,1,2]
        truth.append(op)
        public=tuple((x,k.operator_apply(domain,op,x)) for x in xs)
        roles.append(public)
    chain=tuple((rng.randrange(4),bool(rng.randrange(2)))
          for _ in range(10+(abs(seed)%2)*2))
    initial=rng.randrange(11 if domain==0 else 3)
    state=k.start(domain)
    for r,inv in chain:state=k.compose_state(state,truth[r],domain,inv)
    expected=k.apply(state,initial)
    return (domain,tuple(roles),chain,initial),expected
def code(k,domain,operator):
    if domain==0:return operator[0]*11+operator[1]
    return k.PERM_ID[operator]
def parity(p):
    return sum(p[i]>p[j] for i in range(3) for j in range(i+1,3))%2
def fit_prior(k,training):
    # Pseudolabels come from THREE disclosed PUBLIC observations only.
    # Learner never receives evaluator true operator or source distribution.
    counts=[0]*12
    role_count=0
    for public,expected in training:
        domain,roles,_,_=public
        for pairs in roles:
            possible=[op for op in k.candidates(domain)
                if all(k.operator_apply(domain,op,x)==y for x,y in pairs)]
            need(len(possible)==1,"TRAIN_PUBLIC_CONSISTENCY_NOT_IDENTIFIABLE")
            op=possible[0]
            counts[(op[0]-1) if domain==0 else 10+parity(op)]+=1
            role_count+=1
    need(role_count==128,"EXACT_TRAIN_SOURCE_BUDGET")
    return bytes(min(127,x) for x in counts),{
      "public_pseudolabel_role_count":role_count,
      "learned_12_parameters":tuple(counts),
      "private_true_operator_used":False}
def decide(k,domain,pairs,prior,mode):
    need(len(pairs)==1 and len(prior)==12,"STRICT_EQUAL_ONE_PUBLIC_PAIR")
    options=[]
    for op in k.candidates(domain):
        matches=sum(k.operator_apply(domain,op,x)==y for x,y in pairs)
        weight=prior[(op[0]-1) if domain==0 else 10+parity(op)] if mode=="learned" else 0
        options.append((matches*256+weight,-code(k,domain,op),op))
    return max(options)[2]
def predict(k,case,prior,mode):
    domain,roles,chain,start=case
    inferred=[decide(k,domain,(pairs[0],),prior,mode) for pairs in roles]
    state=k.start(domain)
    for r,inv in chain:state=k.compose_state(state,inferred[r],domain,inv)
    return k.apply(state,start)
def measure(fn):
    tracemalloc.start();t=time.process_time_ns()
    v=fn();receipt={"cpu_ns":time.process_time_ns()-t,
                      "python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    tracemalloc.stop()
    return v,receipt
def main():
    f=json.loads(FREEZE.read_text())
    need(f["schema"]=="research.northstar.t5.public-ambiguity-prior.v1"
         and f["parent_t4_head"]=="1b6167a3c8792b6000c80505fd3a7b465cdb28d0"
         and f["state"]=="IMMUTABLE_PREIMPLEMENTATION"
         and not f["authority"]["protected_original_first_use"],"FIRST_USE_FREEZE")
    k=load(KERNEL,"wing_t5_typed")
    train=[source(i,k) for i in TRAIN]
    tests=[source(i,k) for i in OPEN]
    (prior,fit),fitcost=measure(lambda:fit_prior(k,train))
    ((learned,control),predmeta)=measure(lambda:(
       [predict(k,c,prior,"learned") for c,_ in tests],
       [predict(k,c,bytes(12),"exact_public") for c,_ in tests]))
    (onlylearn,costl)=measure(lambda:[predict(k,c,prior,"learned") for c,_ in tests])
    (onlyctrl,costb)=measure(lambda:[predict(k,c,bytes(12),"exact_public") for c,_ in tests])
    need(learned==onlylearn and control==onlyctrl,"DETERMINISTIC_COMPARABLE_PREDICTIONS")
    answers=[answer for _,answer in tests]
    n=len(answers)
    gain=sum(a==b for a,b in zip(onlylearn,answers))-sum(a==b for a,b in zip(onlyctrl,answers))
    regressions=sum(onlyctrl[i]==answers[i] and onlylearn[i]!=answers[i] for i in range(n))
    out={"schema":"wingless.t5.actual-shared-public-prior-learning.v1",
       "independent_train_tasks":len(train),"open_diagnostic_tasks":n,
       "TRAIN_three_role_pairs_per_role":True,"OPEN_one_public_pair_per_role":True,
       "registered_original_protected_cases_used":0,
       "trained_prior_bytes":len(prior),"packed_recurrent_state_bytes":3,
       "public_train_role_count":fit["public_pseudolabel_role_count"],
       "public_train_count_weights":fit["learned_12_parameters"],
       "learned_changes_from_no_prior":sum(a!=b for a,b in zip(onlylearn,onlyctrl)),
       "learned_correct":sum(a==b for a,b in zip(onlylearn,answers)),
       "exact_information_matched_control_correct":sum(a==b for a,b in zip(onlyctrl,answers)),
       "learned_minus_control":gain,
       "baseline_correct_cases_regressed":regressions,
       "train_cpu_peak":fitcost,"learned_inference_cpu_peak":costl,
       "matched_exact_control_cpu_peak":costb,
       "candidate_train_amortized_CPU_ns_per_task":(fitcost["cpu_ns"]+costl["cpu_ns"])/n,
       "control_CPU_ns_per_task":costb["cpu_ns"]/n,
       "full_resource_parity":(fitcost["cpu_ns"]+costl["cpu_ns"]<=costb["cpu_ns"] and max(
                  fitcost["python_peak_bytes"],costl["python_peak_bytes"])<=costb["python_peak_bytes"]),
       "strict_advancement":False,"novel_S3_family_transfer_established":False,
       "protected_scientific_first_use":False}
    need(out["learned_changes_from_no_prior"]>=0 and len(prior)==12,"DIAGNOSTIC_VALID")
    print(json.dumps(out,sort_keys=True))
if __name__=="__main__":main()
