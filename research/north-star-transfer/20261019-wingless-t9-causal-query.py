"""T9 causally sufficient PUBLIC information acquisition in Wingless.

This is a SYMBOLIC exact acquisition policy, not evidence of learned OOD
cognition. One role calibration pair is insufficient for F11/S3. Adjacent
inverse cancellation identifies the roles that actually affect the final
answer. Request exactly one new, distinct public input for each active role.
No evaluator-only operator is passed into selection or inference.
"""
from __future__ import annotations
import random,json,sys,importlib.util,time,tracemalloc
from pathlib import Path
HERE=Path(__file__).resolve().parent
FREEZE=HERE/"20261019-wingless-t9-source-freeze.json"
KERNEL=HERE/"20261018-wingless-t3-typed-composition.py"
SEEDS=tuple(range(-20031,-19999))
def need(x,n):
    if not x:raise AssertionError("WING_T9_"+n)
def load():
    s=importlib.util.spec_from_file_location("wing_t9_typed",KERNEL)
    k=importlib.util.module_from_spec(s);sys.modules[s.name]=k;s.loader.exec_module(k)
    return k
def reduce_word(word):
    stack=[]
    for role,inv in word:
        if stack and stack[-1][0]==role and stack[-1][1]!=inv:stack.pop()
        else:stack.append((role,inv))
    return tuple(stack)
def born(seed,k):
    need(seed in SEEDS,"ONLY_NEW_SOURCE_NEUTRAL_FIXTURES")
    rng=random.Random(971911+abs(seed)*6113)
    domain=k.AFFINE if abs(seed)%2==0 else k.NONCOMMUTING_PERM
    n=11 if domain==0 else 3
    truth=[]
    public=[]
    for j in range(4):
        op=k.candidates(domain)[rng.randrange(len(k.candidates(domain)))]
        xs=rng.sample(range(n),2)
        truth.append(op)
        public.append(tuple((x,k.operator_apply(domain,op,x)) for x in xs))
    core=[(0,False),(1,False)] if abs(seed)%4<2 else [(0,False),(1,True),(0,False),(2,False)]
    word=core[:]
    while len(word)<(10 if len(core)==2 else 12):
        j=rng.randrange(4);inverse=bool(rng.randrange(2))
        at=rng.randrange(len(word)+1)
        word[at:at]=[(j,inverse),(j,not inverse)]
    need(len(word) in (10,12),"SOURCE_TEN_TWELVE")
    start=rng.randrange(n)
    state=k.start(domain)
    for role,inverse in word:state=k.compose_state(state,truth[role],domain,inverse)
    return {"domain":domain,"two_public_role_pairs":tuple(public),
            "word":tuple(word),"start":start}, k.apply(state,start)
def infer(k,domain,pairs):
    eligible=[p for p in k.candidates(domain)
              if all(k.operator_apply(domain,p,x)==y for x,y in pairs)]
    need(len(eligible)>=1,"TRUE_ROLE_CONSISTENT_WITH_PUBLIC_EVIDENCE")
    return eligible[0],len(eligible)
def predict(k,public,mode):
    domain=public["domain"]
    word=reduce_word(public["word"])
    active={r for r,_ in word}
    roles=[]
    acquired=0
    for role,pairs in enumerate(public["two_public_role_pairs"]):
        if mode=="adaptive":
            chosen=pairs if role in active else (pairs[0],)
        elif mode=="full":chosen=pairs
        else:chosen=(pairs[0],)
        if len(chosen)==2:acquired+=1
        operator,poss=infer(k,domain,chosen)
        if role in active and mode!="one":
            need(poss==1,"TWO_DISTINCT_PUBLIC_INPUTS_UNIQUELY_INFER_OPERATOR")
        roles.append(operator)
    state=k.start(domain)
    for r,inv in word:state=k.compose_state(state,roles[r],domain,inv)
    return k.apply(state,public["start"]),{
      "active_roles":len(active),
      "extra_public_queries":acquired,
      "first_public_queries":4,"total_public_queries":4+acquired,
      "fully_identified_relevant_roles":mode!="one",
      "reduced_word_length":len(word),"original_word_length":len(public["word"])}
def measured(fn):
    tracemalloc.start();cpu=time.process_time_ns()
    v=fn();d={"CPU_ns":time.process_time_ns()-cpu,
              "Python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    tracemalloc.stop();return v,d
def main():
    f=json.loads(FREEZE.read_text())
    need(f["schema"]=="northstar.t9.causal-sufficiency.preimplementation.v1"
        and f["parent_exact_head"]=="08e84aab4b5572a21027cb417732c829c1e4f991"
        and f["authority"]["protected_source_first_use"] is False,"FREEZE")
    k=load()
    cases=[born(s,k) for s in SEEDS]
    ((adapt,meta),acost)=measured(lambda:(
        [predict(k,public,"adaptive")[0] for public,truth in cases],
        [predict(k,public,"adaptive")[1] for public,truth in cases]))
    (full,fcost)=measured(lambda:[
        predict(k,public,"full")[0] for public,truth in cases])
    (once,ocost)=measured(lambda:[
        predict(k,public,"one")[0] for public,truth in cases])
    targets=[truth for _,truth in cases]
    need(adapt==full==targets,"ADAPTIVE_CAUSALLY_SUFFICIENT_EVIDENCE")
    # Public query acquisition and source creation are accounted
    # separately from model inference; all methods share same source.
    total=lambda key:sum(x[key] for x in meta)
    result={"schema":"wingless.t9.certified-causal-public-query-policy.v1",
       "fresh_source_neutral_tasks":len(cases),
       "exact_answers_adaptive":sum(a==b for a,b in zip(adapt,targets)),
       "exact_answers_full_two_pair":sum(a==b for a,b in zip(full,targets)),
       "exact_answers_one_pair":sum(a==b for a,b in zip(once,targets)),
       "baseline_correct_regressions":sum(a==b and c!=b for a,c,b
                                         in zip(once,adapt,targets)),
       "adaptive_added_public_queries":total("extra_public_queries"),
       "full_added_public_queries":len(cases)*4,
       "adaptive_fewer_public_queries":total("extra_public_queries")<len(cases)*4,
       "compressed_causal_word_operations":total("reduced_word_length"),
       "original_public_word_operations":total("original_word_length"),
       "adaptive_model_cpu_peak":acost,"full_public_model_cpu_peak":fcost,
       "one_pair_model_cpu_peak":ocost,
       "additional_public_observation_request_cost_measured":False,
       "cross_task_learned_parameter_updates":0,
       "separately_learned_transfer_demonstrated":False,
       "original_protected_WNS_used":0,
       "consumed_T6_T7_sources_replayed":False,
       "scientific_acceptance":False}
    print(json.dumps(result,sort_keys=True))
if __name__=="__main__":main()
