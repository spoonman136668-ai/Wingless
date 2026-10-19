"""T16 Wingless constant-domain *sound* public-image propagation.
TRAIN-only role ordering; exact algebraic two-pair full fallback.
One-pair role reuse is treated independently (safe superset).
All true operators and target answer are scoring-only source data.
"""
from __future__ import annotations
import sys,importlib.util,json,time,tracemalloc
from pathlib import Path
D=Path(__file__).parent
def load(p,name):
    s=importlib.util.spec_from_file_location(name,p);m=importlib.util.module_from_spec(s)
    sys.modules[name]=m;s.loader.exec_module(m);return m
def metered(fn):
    tracemalloc.start();start=time.process_time_ns()
    try:return fn(),{"cpu_ns":time.process_time_ns()-start,"python_peak_bytes":tracemalloc.get_traced_memory()[1]}
    finally:tracemalloc.stop()
def exact_operator(dom,rows):
    assert len(rows)==2 and rows[0][0]!=rows[1][0]
    (x,y),(u,v)=rows
    if dom==0:
        a=((v-y)*pow((u-x)%11,-1,11))%11
        b=(y-a*x)%11
        assert a and (a*u+b)%11==v
        return (a,b)
    mapping=[None]*3
    mapping[x]=y;mapping[u]=v
    missing_input=next(i for i,z in enumerate(mapping) if z is None)
    missing_output=next(i for i in range(3) if i not in (y,v))
    mapping[missing_input]=missing_output
    return tuple(mapping)
def propagate(t10,k,case,rows):
    dom,_,word,start=case
    n=11 if dom==0 else 3;allbits=(1<<n)-1
    states=1<<start;ops=0
    for role,inverse in t10.word(word):
        samples=rows[role]
        new=0
        if len(samples)==1:
            x,y=samples[0]
            special_out=(x if inverse else y)
            special_in=(y if inverse else x)
            # A single known bijection pair guarantees the same output
            # only on its sampled input. All other outputs are possible
            # except that mapped output; safe superset for F11 and S3.
            if states & (1<<special_in):new|=1<<special_out
            if states & (allbits^(1<<special_in)):
                new|=allbits^(1<<special_out)
            ops+=2
        else:
            op=exact_operator(dom,samples)
            if inverse:
                if dom==0:
                    a,b=op;lookup=lambda j:((j-b)*pow(a,-1,11))%11
                else:lookup=lambda j:op.index(j)
            else:lookup=lambda j:k.operator_apply(dom,op,j)
            for j in range(n):
                if states & (1<<j):new|=1<<lookup(j);ops+=1
        states=new
    return states,ops
def complete_exact(t10,k,case):
    dom,roles,word,start=case
    op={j:exact_operator(dom,roles[j]) for j in {r for r,v in t10.word(word)}}
    st=k.start(dom)
    for role,inv in t10.word(word):st=k.compose_state(st,op[role],dom,inv)
    return k.apply(st,start)
def learn(t10,k,source):
    weight=[0]*4;updates=0
    for seed in source.TRAIN:
        case,_=source.birth(seed,k)
        dom,roles,word,start=case
        rows=[tuple(z[:1]) for z in roles]
        before,_=propagate(t10,k,case,rows)
        for j in sorted({x for x,_ in t10.word(word)}):
            row=list(rows);row[j]=tuple(roles[j])
            after,_=propagate(t10,k,case,row)
            # TRAIN-only public ambiguity reduction, not hidden answer.
            signal=before.bit_count()-after.bit_count()
            weight[j]+=max(-8,min(8,signal))
            updates+=1
    return tuple(weight),updates
def decide(t10,k,case,weights):
    dom,roles,word,start=case
    active=sorted({j for j,inv in t10.word(word)})
    used=set();rows=[tuple(role[:1]) for role in roles];work=0
    while True:
        states,cost=propagate(t10,k,case,rows);work+=cost
        if states.bit_count()==1:
            return states.bit_length()-1,{"queries":len(used),"sound":True,
                       "fallback":False,"state_work":work}
        remaining=[j for j in active if j not in used]
        if not remaining:break
        count={j:sum(role==j for role,inv in t10.word(word)) for j in active}
        j=max(remaining,key=lambda i:(weights[i]+count[i],-i))
        used.add(j);rows[j]=tuple(roles[j])
    return complete_exact(t10,k,case),{"queries":len(active),
                                      "sound":False,"fallback":True,"state_work":work}
def main():
    freeze=json.loads((D/"20261019-wingless-t16-freeze.json").read_text())
    assert freeze["schema"]=="northstar.t16.wingless.preimplementation.freeze.v1"
    assert freeze["source_blob"]=="5bcec7c2440412bdb8c29f8f162edb485b655ab8"
    src=load(D/"20261019-wingless-t16-source.py","w16_src")
    t10=load(D/"20261019-wingless-t10-learned-query.py","w16_t10")
    t11=load(D/"20261019-wingless-t11-mechanism.py","w16_t11")
    t15=load(D/"20261019-wingless-t15-mechanism.py","w16_t15")
    t14=load(D/"20261019-wingless-t14-mechanism.py","w16_t14")
    k=load(D/"20261018-wingless-t3-typed-composition.py","w16_k")
    (weight,changes),traincost=metered(lambda:learn(t10,k,src))
    opens=[src.birth(seed,k) for seed in src.OPEN]
    cand,ccost=metered(lambda:[decide(t10,k,case,weight) for case,y in opens])
    structural,scost=metered(lambda:[decide(t10,k,case,(0,0,0,0)) for case,y in opens])
    full,fcost=metered(lambda:[t11.choose(t10,k,case,(0,)*6,"full") for case,y in opens])
    analytic,acost=metered(lambda:[complete_exact(t10,k,case) for case,y in opens])
    prior,pcost=metered(lambda:[t15.candidate(t10,t11,k,case,(0,)*6) for case,y in opens])
    score=lambda rows:sum((r if isinstance(r,int) else r[0])==y for r,(case,y) in zip(rows,opens))
    loss=lambda base:sum((r if isinstance(r,int) else r[0])==y and cand[i][0]!=y
        for i,(r,(case,y)) in enumerate(zip(base,opens)))
    report={"schema":"wingless.t16.symbolic-domain-safe-learned-query.v1",
        "TRAIN":len(src.TRAIN),"OPEN":len(opens),"TRAIN_public_updates":changes,
        "learned_weights":weight,"candidate_correct":score(cand),
        "untrained_same_representation_correct":score(structural),
        "frozen_full_exact_correct":score(full),
        "analytic_exact_correct":score(analytic),"frozen_T15_correct":score(prior),
        "candidate_regressed_frozen_full":loss(full),
        "candidate_regressed_untrained_structural":loss(structural),
        "candidate_queries":sum(x[1]["queries"] for x in cand),
        "untrained_queries":sum(x[1]["queries"] for x in structural),
        "strong_full_queries":sum(x[1]["queries"] for x in full),
        "certified":sum(x[1]["sound"] for x in cand),
        "public_only_sound_image":True,"candidate_cpu_ns":ccost["cpu_ns"],
        "train_cpu_ns":traincost["cpu_ns"],
        "candidate_amortized_cpu_ns":ccost["cpu_ns"]+traincost["cpu_ns"]//len(opens),
        "full_cpu_ns":fcost["cpu_ns"],"analytic_full_cpu_ns":acost["cpu_ns"],
        "frozen_T15_cpu_ns":pcost["cpu_ns"],
        "candidate_training_inclusive_python_peak_bytes":max(ccost["python_peak_bytes"],traincost["python_peak_bytes"]),
        "full_python_peak_bytes":fcost["python_peak_bytes"],
        "untrained_python_peak_bytes":scost["python_peak_bytes"],
        "remote_query_cost_qualified":False,"independent_novel_source_custody":False,
        "signed_scientific_acceptance":False,"optimizer_admission":False}
    assert score(cand)==len(opens) and report["candidate_regressed_frozen_full"]==0
    assert score(analytic)==score(full)
    print(json.dumps(report,sort_keys=True))
if __name__=="__main__":main()
