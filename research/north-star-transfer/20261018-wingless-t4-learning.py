"""T4 Wingless: actual public TRAIN learning vs untrained and exact controls.

The only instantiated IDs are negative-numbered TRAIN and OPEN diagnostic,
not earlier sealed original/replication sources. FIT never sees a true role
class. The source evaluator keeps hidden role maps solely to score answers.
"""
from __future__ import annotations
import importlib.util,sys,json,random,time,tracemalloc
from pathlib import Path
HERE=Path(__file__).resolve().parent
T3=HERE/"20261018-wingless-t3-typed-composition.py"
FROZEN=HERE/"20261018-wingless-t4-frozen-train-only.json"
TRAIN=tuple(range(-12031,-11999))
OPEN=tuple(range(-12115,-12099))
def need(x,msg):
    if not x:raise ValueError("WNS_T4_"+msg)
def load():
    spec=importlib.util.spec_from_file_location("wing_t4_native",T3)
    module=importlib.util.module_from_spec(spec)
    sys.modules[spec.name]=module
    spec.loader.exec_module(module)
    return module
def source(seed,m):
    need(seed in TRAIN or seed in OPEN,"REGISTERED_PROTECTED_COHORT")
    domain=m.AFFINE if (-seed)%2==0 else m.NONCOMMUTING_PERM
    rng=random.Random(317123+(-seed)*991)
    choices=m.candidates(domain)
    n=11 if domain==m.AFFINE else 3
    role_truth=[choices[rng.randrange(len(choices))] for _ in range(4)]
    public=[]
    for truth in role_truth:
        xs=rng.sample(range(n),2)
        public.append(tuple((x,m.operator_apply(domain,truth,x)) for x in xs))
    chain=tuple((rng.randrange(4),bool(rng.randrange(2)))
                for _ in range(10+((-seed)%2)*2))
    start=rng.randrange(n)
    state=m.start(domain)
    for role,invert in chain:state=m.compose_state(state,role_truth[role],domain,invert)
    evaluator=m.apply(state,start)
    return domain,tuple(public),chain,start,evaluator
def predict(m,case,weights,exact=False):
    domain,public,chain,start,_=case
    inferred=[]
    for role in public:
        if exact:
            choices=m.candidates(domain)
            def evidence(k):
                c=sum(m.operator_apply(domain,k,x)==y for x,y in role)
                code=m.PERM_ID[k] if domain==1 else k[0]*11+k[1]
                return (c,-code)
            inferred.append(max(choices,key=evidence))
        else:inferred.append(m.infer_role(domain,role,weights)[0])
    state=m.start(domain)
    for role,inv in chain:
        state=m.compose_state(state,inferred[role],domain,inv)
    return m.apply(state,start)
def measure(f):
    tracemalloc.start();start=time.process_time_ns()
    out=f()
    cost={"cpu_ns":time.process_time_ns()-start,
          "python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    tracemalloc.stop()
    return out,cost
def main():
    d=json.loads(FROZEN.read_text())
    need(d["schema"]=="research.northstar.t4.train-only-learning-resource-parity.v1"
         and d["frozen_parent_head"]=="f49246ab6f4858c1f994d806cbbe31264f4ca4a1"
         and not d["original_protected_holdouts_allowed"],"FROZEN_CONTRACT")
    m=load()
    training=[source(seed,m) for seed in TRAIN]
    cases=[source(seed,m) for seed in OPEN]
    samples=[(domain,pair) for domain,pairs,*_ in training for pair in pairs]
    (learned,trainmeta),traincost=measure(lambda:m.flatten(samples))
    zero=bytes(12)
    def batch(weights,exact=False):
        return [predict(m,case,weights,exact) for case in cases]
    l,lc=measure(lambda:batch(learned))
    z,zc=measure(lambda:batch(zero))
    b,bc=measure(lambda:batch(zero,True))
    expected=[case[-1] for case in cases]
    correct=lambda results:sum(x==y for x,y in zip(results,expected))
    changed=sum(x!=y for x,y in zip(l,z))
    n=len(OPEN)
    out={"schema":"wingless.t4.public-learning-parity.v1",
         "public_only_training_samples":len(samples),"public_training_tasks":len(TRAIN),
         "open_diagnostic_tasks":n,"registered_original_holdouts_used":0,
         "model_weights_changed_from_zero":learned!=zero,
         "model_parameters_bytes":len(learned),"recurrent_state_bytes":3,
         "learned_correct":correct(l),"zero_prior_correct":correct(z),
         "exact_public_control_correct":correct(b),
         "learned_prediction_changes_from_untrained":changed,
         "training_cost":traincost,"learned_inference_cost":lc,
         "zero_prior_cost":zc,"exact_public_control_cost":bc,
         "train_parameter_updates":trainmeta["public_only_training_feature_updates"],
         "learned_train_amortized_CPU_per_task_ns":
             lc["cpu_ns"]/n+traincost["cpu_ns"]/n,
         "learned_complete_peak_bytes_lower_bound":
             max(lc["python_peak_bytes"],traincost["python_peak_bytes"]),
         "control_peak_bytes":bc["python_peak_bytes"],
         "strongest_control_accuracy_undercut":correct(l)<correct(b),
         "strict_accuracy_gain_over_strongest_control":correct(l)>correct(b),
         "full_cost_nonincrease":lc["cpu_ns"]+traincost["cpu_ns"]<=bc["cpu_ns"]
               and max(lc["python_peak_bytes"],traincost["python_peak_bytes"])
                   <=bc["python_peak_bytes"],
         "eligible_protected_first_use":False,
         "scientific_transfer_accepted":False}
    need(out["model_parameters_bytes"]==12 and correct(b)==n
         and changed==0,"LEARNED_PRIOR_CAN_ONLY_TIE_EXACT_PUBLIC_SOLVER")
    print(json.dumps(out,sort_keys=True))
if __name__=="__main__":main()
