"""T15 Wingless: joint-operator sound certificate instead of T14 DP.
Only public role pairs enter inference; hidden target is evaluator-only.
A budget-exhausted / non-singleton certificate always uses exact full evidence.
This synthetic source-only probe does not authorize scientific acceptance.
"""
from __future__ import annotations
import itertools,importlib.util,json,sys,time,tracemalloc
from pathlib import Path
D=Path(__file__).parent;CAP=32
def load(file,name):
    s=importlib.util.spec_from_file_location(name,file);m=importlib.util.module_from_spec(s)
    sys.modules[name]=m;s.loader.exec_module(m);return m
def measured(fn):
    tracemalloc.start();tick=time.process_time_ns()
    try:return fn(),{"cpu_ns":time.process_time_ns()-tick,"python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    finally:tracemalloc.stop()
def certificate(t10,t11,k,case,public,limit):
    dom,roles,word,start=case
    active=sorted({j for j,inv in t10.word(word)})
    options=[t10.possible(k,dom,public[j]) for j in active]
    work=1
    for opt in options:
        work*=len(opt)
        if work>limit:return None,0
    values=set()
    for joint in itertools.product(*options):
        chosen={j:op for j,op in zip(active,joint)}
        state=k.start(dom)
        for j,inv in t10.word(word):
            state=k.compose_state(state,chosen[j],dom,inv)
        values.add(k.apply(state,start))
        if len(values)>1:break
    # Returning a singleton is sound only if ALL assignments were evaluated.
    # If early break happened, more than one distinct result was observed.
    return (next(iter(values)) if len(values)==1 else None),work
def candidate(t10,t11,k,case,weights):
    dom,roles,word,start=case
    active=sorted({j for j,inv in t10.word(word)})
    visible=[tuple(p[:1]) for p in roles]
    revealed=set();work=0;checks=0
    while work<CAP:
        value,budget=certificate(t10,t11,k,case,visible,CAP-work)
        checks+=1;work+=budget
        if value is not None:
            return value,{"queries":len(revealed),"certified":True,
               "fallback":False,"joint_public_assignments":work}
        remaining=[j for j in active if j not in revealed]
        if not remaining:break
        # Trained causal score only selects which PUBLIC second pair to buy.
        j=max(remaining,key=lambda j:(
            sum(a*b for a,b in zip(weights,t11.encode(t10,k,case,j))),-j))
        revealed.add(j);visible[j]=tuple(roles[j])
    result,_=t11.choose(t10,k,case,weights,"full")
    return result,{"queries":len(active),"certified":False,"fallback":True,
         "joint_public_assignments":work}
def main():
    f=json.loads((D/"20261019-wingless-t15-preimplementation-freeze.json").read_text())
    assert f["schema"]=="northstar.t15.wingless.freeze.v1"
    assert f["source_blob"]=="ffa1e00cf2ceec3ab412fd06f8da5d25d5794f49"
    src=load(D/"20261019-wingless-t15-source.py","w15src")
    t10=load(D/"20261019-wingless-t10-learned-query.py","w15t10")
    t11=load(D/"20261019-wingless-t11-mechanism.py","w15t11")
    t14=load(D/"20261019-wingless-t14-mechanism.py","w15t14")
    k=load(D/"20261018-wingless-t3-typed-composition.py","w15k")
    train=[src.birth(i,k) for i in src.TRAIN]
    opens=[src.birth(i,k) for i in src.OPEN]
    (weights,diag),tcost=measured(lambda:t11.fit(t10,k,train))
    assert diag["parameter_updates"]>0 and len(weights)==6
    candidate_rows,ccost=measured(lambda:[candidate(t10,t11,k,c,weights) for c,y in opens])
    full_rows,fcost=measured(lambda:[t11.choose(t10,k,c,weights,"full") for c,y in opens])
    t14_rows,t14cost=measured(lambda:[t14.candidate(k,t10,t11,c,weights) for c,y in opens])
    correct=lambda rows:sum(x[0]==truth for x,(_,truth) in zip(rows,opens))
    lost=lambda baseline,rows:sum(x[0]==truth and y[0]!=truth for x,y,(_,truth) in zip(baseline,rows,opens))
    q=lambda rows:sum(v[1]["queries"] if "queries" in v[1] else v[1]["public_queries"] for v in rows)
    out={"schema":"wingless.t15.original-public-joint-certificate.v1","source_git_blob":f["source_blob"],
      "TRAIN":len(train),"OPEN":len(opens),"TRAIN_weight_updates":diag["parameter_updates"],
      "candidate_correct":correct(candidate_rows),"T14_control_correct":correct(t14_rows),
      "strong_full_correct":correct(full_rows),"candidate_regressions_vs_full":lost(full_rows,candidate_rows),
      "candidate_public_queries":q(candidate_rows),"T14_public_queries":q(t14_rows),
      "strong_full_public_queries":q(full_rows),
      "candidate_certified":sum(v[1]["certified"] for v in candidate_rows),
      "candidate_fallback":sum(v[1]["fallback"] for v in candidate_rows),
      "joint_assignment_budget":CAP,
      "amortized_candidate_cpu_ns":ccost["cpu_ns"]+tcost["cpu_ns"]//len(opens),
      "full_cpu_ns":fcost["cpu_ns"],"T14_cpu_ns":t14cost["cpu_ns"],
      "candidate_python_peak_training_inclusive":max(ccost["python_peak_bytes"],tcost["python_peak_bytes"]),
      "strong_full_python_peak":fcost["python_peak_bytes"],"TRAIN_cpu_ns":tcost["cpu_ns"],
      "remote_query_service_cost_authenticated":False,"independently_signed_source":False,
      "protected_first_use":False,"signed_scientific_acceptance":False,
      "optimizer_admission":False}
    assert out["candidate_regressions_vs_full"]==0 and out["candidate_correct"]==len(opens)
    assert out["candidate_public_queries"]<=out["strong_full_public_queries"]
    print(json.dumps(out,sort_keys=True))
if __name__=="__main__":main()
