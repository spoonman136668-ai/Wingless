"""T6 public-confidence learning with direct F11 candidate evaluation.

Fresh T6 negative TRAIN and OPEN source identity namespaces only. Fixed
16:1 source-blind margin protects some exact baseline decisions, but no
post-hoc threshold search or protected cases are permitted.
"""
from __future__ import annotations
import importlib.util,types,json,time,tracemalloc
from pathlib import Path
HERE=Path(__file__).resolve().parent
FROZEN=HERE/"20261018-wingless-t6-freeze.json"
T5=HERE/"20261018-wingless-t5-learning.py"
TRAIN=range(-14031,-13999)
OPEN=range(-14163,-14099)
ODDS=16
def need(ok,tag):
    if not ok:raise ValueError("WING_T6_"+tag)
def fresh_source(t5,seed,k):
    need(seed in TRAIN or seed in OPEN,"REGISTERED_SOURCE_NOT_PERMITTED")
    fn=types.FunctionType(t5.source.__code__,
             dict(t5.source.__globals__,TRAIN=TRAIN,OPEN=OPEN),
             t5.source.__name__,t5.source.__defaults__)
    return fn(seed,k)
def learned_direct(k,domain,pairs,prior):
    need(len(pairs)==1 and len(prior)==12,"ONE_PUBLIC_PAIR_AND_TWELVE_BYTES")
    (x,y),=pairs
    if domain==k.AFFINE:
        baseline_slope=1
        baseline=(baseline_slope,(y-baseline_slope*x)%11)
        slope=max(range(1,11),key=lambda a:(prior[a-1],-a))
        best=(slope,(y-slope*x)%11)
        select=(slope!=baseline_slope and
                prior[slope-1]>=ODDS*max(1,prior[baseline_slope-1]))
        return best if select else baseline
    possible=[op for op in k.PERMS if op[x]==y]
    baseline=min(possible,key=k.PERM_ID.__getitem__)
    preferred=max(possible,key=lambda op:(prior[10+t5_parity(op)],-k.PERM_ID[op]))
    # Quantized 8bit odds; existing S3 parity posterior is <16:1 in TRAIN.
    select=(preferred!=baseline and
      prior[10+t5_parity(preferred)]>=ODDS*max(1,prior[10+t5_parity(baseline)]))
    return preferred if select else baseline
def t5_parity(p):
    return sum(p[i]>p[j] for i in range(3) for j in range(i+1,3))%2
def predict(k,public,prior,mode,t5):
    domain,roles,chain,start=public
    ops=[]
    for pairs in roles:
        if mode=="direct":
            ops.append(learned_direct(k,domain,(pairs[0],),prior))
        else:
            ops.append(t5.decide(k,domain,(pairs[0],),prior,mode))
    state=k.start(domain)
    for role,inverse in chain:
        state=k.compose_state(state,ops[role],domain,inverse)
    return k.apply(state,start)
def measure(fn):
    tracemalloc.start();t=time.process_time_ns()
    result=fn();cost={"cpu_ns":time.process_time_ns()-t,
                        "python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    tracemalloc.stop();return result,cost
def main():
    f=json.loads(FROZEN.read_text())
    need(f["schema"]=="research.northstar.t6.frozen.v1"
         and f["parent_exact_head"]=="ee45c3f157f01292074fb7a14d4528bb8f980ccb"
         and not f["protected_source_first_use"],"PREIMPLEMENTATION")
    t5=__import__("importlib").util.spec_from_file_location("wing_t6_parent",T5)
    m=__import__("importlib").util.module_from_spec(t5)
    import sys
    sys.modules[t5.name]=m;t5.loader.exec_module(m)
    k=m.load(m.KERNEL,"wing_t6_typed")
    train=[fresh_source(m,i,k) for i in TRAIN]
    open_cases=[fresh_source(m,i,k) for i in OPEN]
    (prior,receipt),tcost=measure(lambda:m.fit_prior(k,train))
    need(len(prior)==12 and receipt["public_pseudolabel_role_count"]==128,"PUBLIC_TRAINING")
    test=lambda mode:[predict(k,p,prior,mode,m) for p,y in open_cases]
    candidate,ccost=measure(lambda:test("direct"))
    baseline,bcost=measure(lambda:test("exact_public"))
    t5prior,p5cost=measure(lambda:test("learned"))
    expected=[y for p,y in open_cases]
    correct=lambda p:sum(x==y for x,y in zip(p,expected))
    regress=sum(a==y and b!=y for a,b,y in zip(baseline,candidate,expected))
    output={"schema":"wingless.t6.sourceblind-quantized-confidence.v1",
       "frozen_threshold_odds":ODDS,
       "TRAIN_tasks":len(train),"OPEN_tasks":len(open_cases),
       "protected_original_source_used":0,"learned_prior_bytes":len(prior),
       "recurrent_state_bytes":3,"public_train_roles":receipt["public_pseudolabel_role_count"],
       "baseline_correct":correct(baseline),"T5_full_prior_correct":correct(t5prior),
       "T6_conditional_prior_correct":correct(candidate),
       "T6_prediction_changes":sum(a!=b for a,b in zip(candidate,baseline)),
       "baseline_correct_regressions":regress,
       "candidate_train_cost":tcost,"candidate_open_cost":ccost,
       "exact_same_information_control_cost":bcost,"T5_prior_cost":p5cost,
       "train_amortized_candidate_CPU_ns":ccost["cpu_ns"]+tcost["cpu_ns"]/len(open_cases),
       "control_CPU_ns":bcost["cpu_ns"],
       "candidate_total_peak_bytes":max(ccost["python_peak_bytes"],tcost["python_peak_bytes"]),
       "control_peak_bytes":bcost["python_peak_bytes"],
       "CPU_nonincrease":ccost["cpu_ns"]+tcost["cpu_ns"]/len(open_cases)<=bcost["cpu_ns"],
       "peak_nonincrease":max(ccost["python_peak_bytes"],tcost["python_peak_bytes"])<=bcost["python_peak_bytes"],
       "protected_first_use":False,"science_accepted":False,"production":False}
    print(json.dumps(output,sort_keys=True))
if __name__=="__main__":main()
