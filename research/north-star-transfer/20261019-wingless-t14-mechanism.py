"""T14-B Wingless original candidate: cap-priced sound public certificate with
TRAIN-learned query priority; never trust confidence alone as a certificate.
Nonprotected diagnostic; no independent scientific acceptance or production.
"""
from __future__ import annotations
import json,time,tracemalloc
from pathlib import Path
D=Path(__file__).parent
def timed(fn):
    tracemalloc.start();t=time.process_time_ns()
    try:return fn(),{"cpu_ns":time.process_time_ns()-t,"python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    finally:tracemalloc.stop()
def reachable(k,t10,case,rows,limit=48):
    domain,_,word,start=case
    states={start};work=0
    for role,inverse in t10.word(word):
        choices=t10.possible(k,domain,rows[role])
        new=set()
        for v in states:
            for op in choices:
                work+=1
                if work>limit:return None,work
                if inverse:
                    if domain==0:
                        a,b=op;new.add(((v-b)*pow(a,-1,11))%11)
                    else:new.add(op.index(v))
                else:new.add(k.operator_apply(domain,op,v))
        states=new
    return states,work
def candidate(k,t10,t11,case,weights):
    dom,roles,word,start=case
    active=sorted({role for role,_ in t10.word(word)})
    rows=[tuple(p[:1]) for p in roles]
    queried=set()
    cert_calls=0;work=0
    while True:
        # The work limit bounds EACH certificate. No bounded proof => full fallback.
        if work>48:break
        answers,ops=reachable(k,t10,case,rows,48-work)
        work+=ops;cert_calls+=1
        if answers is not None and len(answers)==1:
            return next(iter(answers)),{"public_queries":len(queried),"certificate_work":work,
             "certified":True,"fallback":False,"cert_calls":cert_calls}
        if answers is None:break
        remaining=[j for j in active if j not in queried]
        if not remaining:break
        # Public-only learned rank; no second pair or hidden truth enters features.
        role=max(remaining,key=lambda j:(
            sum(a*b for a,b in zip(weights,t11.encode(t10,k,case,j))),-j))
        rows[role]=tuple(roles[role]);queried.add(role)
    full,_=t11.choose(t10,k,case,weights,"full")
    return full,{"public_queries":len(active),"certificate_work":work,
     "certified":False,"fallback":True,"cert_calls":cert_calls}
def main():
    from importlib.util import spec_from_file_location,module_from_spec
    import sys
    def load(path,name):
        sp=spec_from_file_location(name,path);mo=module_from_spec(sp)
        sys.modules[name]=mo;sp.loader.exec_module(mo);return mo
    f=json.loads((D/"20261019-wingless-t14-complete-source-and-assay-freeze.json").read_text())
    assert f["schema"]=="northstar.t14.wingless.complete-source-assay-freeze.v1"
    assert f["source_blob"]=="3b3d042254c0a123d9058800eec76439096ee47c"
    src=load(D/"20261019-wingless-t14-source.py","w_t14_src")
    t10=load(D/"20261019-wingless-t10-learned-query.py","w_t14_t10")
    t11=load(D/"20261019-wingless-t11-mechanism.py","w_t14_t11")
    k=load(D/"20261018-wingless-t3-typed-composition.py","w_t14_k")
    # Original source birth occurs only when launched under verified original CI.
    train=[src.birth(s,t10,k) for s in src.TRAIN]
    opens=[src.birth(s,t10,k) for s in src.OPEN]
    (weights,fit),train_cost=timed(lambda:t11.fit(t10,k,train))
    assert len(weights)==6 and fit["parameter_updates"]>0
    predictions,cost=timed(lambda:[candidate(k,t10,t11,case,weights) for case,y in opens])
    full,fcost=timed(lambda:[t11.choose(t10,k,c,weights,"full") for c,y in opens])
    exact=lambda rows:sum(x[0]==y for x,(_,y) in zip(rows,opens))
    lost=sum(a[0]==y and b[0]!=y for a,b,(_,y) in zip(full,predictions,opens))
    queries=lambda rows:sum(x[1].get("public_queries",x[1].get("queries",0)) for x in rows)
    report={"schema":"wingless.t14.original-sourceonly-diagnostic.v1",
      "source_definition_git_blob":f["source_blob"],"TRAIN":len(train),"OPEN":len(opens),
      "learned_updates":fit["parameter_updates"],"candidate_correct":exact(predictions),
      "full_strong_correct":exact(full),"baseline_positive_regressions":lost,
      "candidate_queries":queries(predictions),"full_queries":queries(full),
      "certified":sum(x[1]["certified"] for x in predictions),
      "fallback_count":sum(x[1]["fallback"] for x in predictions),
      "candidate_certificate_work":sum(x[1]["certificate_work"] for x in predictions),
      "candidate_CPU_ns":cost["cpu_ns"],"strong_full_CPU_ns":fcost["cpu_ns"],
      "TRAIN_CPU_ns":train_cost["cpu_ns"],
      "TRAIN_amortized_candidate_CPU_ns":cost["cpu_ns"]+train_cost["cpu_ns"]//len(opens),
      "complete_candidate_python_peak_bytes":max(cost["python_peak_bytes"],train_cost["python_peak_bytes"]),
      "full_python_peak_bytes":fcost["python_peak_bytes"],
      "remote_public_service_cost_qualified":False,
      "independent_novel_relationship_replication":False,
      "signed_scientific_acceptance":False,"optimizer_admission":False}
    assert lost==0,"T14_WINGLESS_UNSAFE_BASELINE_REGRESSION"
    assert exact(predictions)==len(opens),"T14_WINGLESS_ORIGINAL_WRONG"
    assert report["candidate_queries"]<=report["full_queries"]
    print(json.dumps(report,sort_keys=True))
if __name__=="__main__":main()
