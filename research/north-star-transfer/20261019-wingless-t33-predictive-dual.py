"""T33: train a predictive dual-syndrome sufficient statistic.

A pair of public observation equations can directly certify one future
causal parity iff the queried hidden law-mask equals their XOR.
The TRAIN-compiled certificate dictionary stores exact source-witness
row-mask pairs, never hidden labels or private future outcomes.
Any absent/invalid certificate falls back to full optimized GF(2).
No guess, pruning of latent graphs, fake online learning, or replay.
"""
from __future__ import annotations
import importlib.util,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-compact-dual-syndrome",
      "TRAIN-frequency-dual-certificate-control",
      "deterministic-all-pairs-dual-certificate-control",
      "optimized-exact-GF2-rowspace",
      "independent-second-pivot-exact-GF2-control")
def source():
    spec=importlib.util.spec_from_file_location("t33_source",ROOT/"20261019-wingless-t33-new-syndrome-worlds.py")
    m=importlib.util.module_from_spec(spec);sys.modules[spec.name]=m;spec.loader.exec_module(m);return m
def rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Mem(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
              ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
              ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
              ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
              ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        k=ctypes.WinDLL("kernel32",use_last_error=True)
        p=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T33_WIN_MEMORY_FAILED:"+str(ctypes.get_last_error()))
        return int(m.PeakWorkingSetSize)
    import resource
    x=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(x)*1024 if sys.platform.startswith("linux") else int(x)
def rowmask(src,f,pre,graph,nodes):
    table=src.TEMPLATES[f][graph]
    m=0
    rhs=0
    for j in nodes:
        m^=1<<src.input_pattern(pre,table[j])
        rhs^=(pre>>j)&1
    return m,rhs
def read_rows(src,public,graph):
    f=public["family"]
    rows=[]
    for p in public["probes"]:
        mask,known=rowmask(src,f,p["pre"],graph,p["nodes"])
        rows.append((mask,p["parity"]^known))
    return tuple(rows)
def goalrow(src,public,graph,query):
    # Transport new j <- previous j-2; map public output nodes back
    # exactly through known deterministic transfer.
    n,_,_=src.FAMILIES[public["family"]]
    prior=tuple((j-2)%n for j in query["targets"])
    return rowmask(src,public["family"],query["pre"],graph,prior)
def gf2_basis(rows,reverse=False):
    b={}
    for orig,r in rows:
        m=orig
        while m:
            pivot=(m&-m).bit_length()-1 if reverse else m.bit_length()-1
            if pivot not in b:
                b[pivot]=(m,r)
                break
            old,oldr=b[pivot];m^=old;r^=oldr
        else:
            if r:return None
    return b
def gf2_decide(b,mask,known,reverse=False):
    while mask:
        pivot=(mask&-mask).bit_length()-1 if reverse else mask.bit_length()-1
        if pivot not in b:return None
        row,parity=b[pivot];mask^=row;known^=parity
    return known
def learn_dictionary(src,TRAIN):
    learned=Counter()
    row_frequency={f:Counter() for f in src.FAMILIES}
    for public in TRAIN:
        f=public["family"]
        for graph in range(3):
            rows=read_rows(src,public,graph)
            observed={m for m,r in rows}
            row_frequency[f].update(observed)
            for q in public["queries"]:
                mask,_=goalrow(src,public,graph,q)
                for ai,a in enumerate(sorted(observed)):
                    for b in sorted(observed)[ai+1:]:
                        if a^b==mask and mask:
                            learned[f,mask,a,b]+=1
    answer={}
    for f in src.FAMILIES:
        common=sorted(((key,count) for key,count in learned.items() if key[0]==f),
                      key=lambda z:(-z[1],z[0]))
        answer[f]=tuple((key[1],key[2],key[3]) for key,count in common[:32])
    # Equally public TRAIN-informed baseline: rows ranked by their
    # frequency alone, with all valid combinations, not query cooccurrence.
    freq={}
    for f in src.FAMILIES:
        ordered=sorted(row_frequency[f],key=lambda m:(-row_frequency[f][m],m))[:32]
        freq[f]=tuple((a^b,a,b) for i,a in enumerate(ordered) for b in ordered[i+1:])
    return answer,freq,{"learned_pairs":sum(len(x) for x in answer.values()),
                        "TRAIN_witness_pattern_types":len(learned)}
def pair_certificate(rows,mask,dictionary,arm):
    lookup=dict(rows)
    # Conflicting duplicate public masks => invalid graph, handled by
    # independent complete basis consistency check before this point.
    if mask in lookup:return lookup[mask],0,True
    if arm==ARMS[2]:
        # Complete deterministic two-row exact certificate search.
        vals=tuple(sorted(lookup))
        for a in vals:
            b=mask^a
            if b!=a and b in lookup:
                return lookup[a]^lookup[b],1,True
        return None,0,False
    for q,a,b in dictionary:
        if q!=mask:continue
        if a in lookup and b in lookup and (a^b)==q:
            return lookup[a]^lookup[b],1,True
    return None,0,False
def infer(src,public,arm,trained,freq):
    f=public["family"]
    all_answers=[set(),set(),set()]
    accepted_graphs=0
    pairhits=directhits=full_reductions=0
    for g in range(3):
        rows=read_rows(src,public,g)
        primary=gf2_basis(rows,reverse=(arm==ARMS[4]))
        if primary is None:continue
        accepted_graphs+=1
        values=[]
        for query in public["queries"]:
            qmask,known=goalrow(src,public,g,query)
            if arm in (ARMS[0],ARMS[1],ARMS[2]):
                dictionary=trained[f] if arm==ARMS[0] else freq[f] if arm==ARMS[1] else ()
                value,count,hit=pair_certificate(rows,qmask,dictionary,arm)
                if hit:
                    if count:pairhits+=1
                    else:directhits+=1
                else:
                    value=gf2_decide(primary,qmask,0)
                    full_reductions+=1
            else:
                value=gf2_decide(primary,qmask,0,reverse=(arm==ARMS[4]))
                full_reductions+=1
            values.append(None if value is None else value^known)
        for i,x in enumerate(values):
            if x is None:
                all_answers[i].update((0,1))
            else:all_answers[i].add(x)
    if not accepted_graphs:raise RuntimeError("T33_CAUSAL_GRAPH_POSTERIOR_EMPTY")
    out=tuple(next(iter(x)) if len(x)==1 else None for x in all_answers)
    return out,tuple(tuple(sorted(x)) for x in all_answers),{"graphs":accepted_graphs,
      "trained_pair_certificate_hits":pairhits,"direct_certificate_hits":directhits,
      "full_GF2_query_reductions":full_reductions}
def train_only():
    src=source()
    if len(src.TRAIN)!=40 or len(src.DEV)!=64 or set(src.TRAIN)&set(src.DEV):
        raise RuntimeError("T33_SOURCE_COHORT_DRIFT")
    tr=[src.public_train(src.sample(i)) for i in src.TRAIN]
    learned,freq,stats=learn_dictionary(src,tr)
    retained=correct=uncertain=equivalent=pairhit=saved=changed=0
    for name,public in zip(src.TRAIN,tr):
        private=src.sample(name)
        reference,possible,rstats=infer(src,public,ARMS[3],learned,freq)
        result,rs,s=infer(src,public,ARMS[0],learned,freq)
        control,*_=infer(src,public,ARMS[1],learned,freq)
        complete,*_=infer(src,public,ARMS[2],learned,freq)
        backwards,*_=infer(src,public,ARMS[4],learned,freq)
        if not (result==reference==control==complete==backwards and rs==possible):
            raise RuntimeError("T33_TRAIN_CERTIFICATE_UNSOUND")
        true_graph=private["_hidden_graph"]
        if gf2_basis(read_rows(src,public,true_graph)) is None:
            raise RuntimeError("T33_TRUE_CAUSAL_GRAPH_REMOVED")
        retained+=1
        pairhit+=s["trained_pair_certificate_hits"]
        saved+=rstats["full_GF2_query_reductions"]-s["full_GF2_query_reductions"]
        if s["full_GF2_query_reductions"]<rstats["full_GF2_query_reductions"]:changed+=1
        for v,options,y in zip(result,possible,private["_future_truth"]):
            if y not in options:raise RuntimeError("T33_TRUE_FUTURE_EXCLUDED")
            if v is not None:
                if v!=y:raise RuntimeError("T33_TRAIN_HARD_FALSE_POSITIVE")
                correct+=1
            else:uncertain+=1
            equivalent+=1
    return {"schema":"wingless.T33.TRAIN-only-learned-exact-dual-syndrome.v1",
      "status":"PASS_PROOF_TRAIN_ONLY","TRAIN_cases":40,"DEV_generated":0,
      "true_hidden_graphs_retained":retained,
      "TRAIN_target_questions":equivalent,"TRAIN_correct_certified":correct,
      "TRAIN_safe_abstentions":uncertain,"TRAIN_uncertainty_fraction":uncertain/equivalent,
      "TRAIN_exact_equivalence_all_five_arms":equivalent,
      "TRAIN_learned_pair_certificate_hits":pairhit,
      "TRAIN_elimination_query_reductions_avoided":saved,
      "TRAIN_caused_actual_query_solver_work_change_cases":changed,
      "all_source_laws_unconstrained_by_fit":True,
      "no_old_DEV_replayed":True,"optimizer_admission":False,**stats}
def one(arm,rep,partition="TRAIN"):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T33_ARM_REP_INVALID")
    start=time.process_time_ns()
    src=source()
    tr=[src.public_train(src.sample(i)) for i in src.TRAIN]
    learned,freq,stats=learn_dictionary(src,tr) if arm in ARMS[:2] else ({f:() for f in src.FAMILIES},{f:() for f in src.FAMILIES},{})
    names=src.TRAIN if partition=="TRAIN" else src.DEV
    result=[];certified=uncertain=pairhits=red=0
    for name in names:
        private=src.sample(name)
        public=src.public_train(private) if partition=="TRAIN" else src.public_dev(private)
        ans,poss,trace=infer(src,public,arm,learned,freq)
        good=unknown=0
        for val,opts,truth in zip(ans,poss,private["_future_truth"]):
            if truth not in opts:raise RuntimeError("T33_FROZEN_TRUE_CAUSAL_ANSWER_EXCLUDED")
            if val is None:unknown+=1;uncertain+=1
            else:
                if val!=truth:raise RuntimeError("T33_FROZEN_HARD_FALSE_POSITIVE")
                good+=1;certified+=1
        pairhits+=trace["trained_pair_certificate_hits"]
        red+=trace["full_GF2_query_reductions"]
        result.append({"id":name,"correct":good,"unknown":unknown,
              "graphs":trace["graphs"],"pair_hits":trace["trained_pair_certificate_hits"],
              "full_GF2_reductions":trace["full_GF2_query_reductions"]})
    return {"schema":"wingless.T33.full-native-process.v1","arm":arm,"repeat":rep,
        "partition":partition,"cases":len(names),"correct":certified,
        "unknown":uncertain,"hard_false_positives":0,"per_case":result,
        "TRAIN_learned_witness_fit":stats,
        "pair_witness_hits":pairhits,"query_GF2_reductions":red,
        "full_process_CPU_ns":time.process_time_ns()-start,
        "peak_RSS_bytes":rss(),"no_old_DEV_replay":True,"optimizer_admission":False}
def compare(partition):
    collection={}
    for arm in ARMS:
        trials=[]
        for i in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                     "--arm",arm,str(i),partition],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T33_CHILD_FAILED:"+arm+":"+str(i)+"\n"+p.stderr)
            trials.append(json.loads(p.stdout))
        z=trials[0]
        if any(a["per_case"]!=z["per_case"] for a in trials):
            raise RuntimeError("T33_UNSTABLE_SOURCE")
        collection[arm]={
          "correct_certified":z["correct"],"sound_abstentions":z["unknown"],
          "hard_false_positives":z["hard_false_positives"],
          "witness_hits":z["pair_witness_hits"],
          "full_GF2_reductions":z["query_GF2_reductions"],
          "per_case":z["per_case"],
          "five_whole_CPU_ns":[t["full_process_CPU_ns"] for t in trials],
          "five_peak_RSS_bytes":[t["peak_RSS_bytes"] for t in trials],
          "median_whole_CPU_ns":int(statistics.median(t["full_process_CPU_ns"] for t in trials)),
          "max_peak_RSS_bytes":max(t["peak_RSS_bytes"] for t in trials)}
    reference=collection[ARMS[3]]
    for arm,v in collection.items():
        if (v["correct_certified"],v["sound_abstentions"])!=(reference["correct_certified"],reference["sound_abstentions"]):
            raise RuntimeError("T33_DIFFERENT_SOUND_CERTIFIED_COVERAGE")
        if any((r["correct"],r["unknown"],r["graphs"])!=(base["correct"],base["unknown"],base["graphs"])
               for r,base in zip(v["per_case"],reference["per_case"])):
            raise RuntimeError("T33_PER_CASE_NOT_EXACTLY_EQUIVALENT")
    candidate=collection[ARMS[0]]
    compared={k:{"no_answer_loss":candidate["correct_certified"]>=v["correct_certified"],
          "whole_CPU_nonregression":candidate["median_whole_CPU_ns"]<=v["median_whole_CPU_ns"],
          "peak_RSS_nonregression":candidate["max_peak_RSS_bytes"]<=v["max_peak_RSS_bytes"]}
      for k,v in collection.items() if k!=ARMS[0]}
    return {"schema":"wingless.T33.original-new-dual-syndrome-comparison.v1",
        "partition":partition,"TRAIN_cases":40,"unseen_DEV_cases":64 if partition=="DEV" else 0,
        "target_questions":(192 if partition=="DEV" else 120),
        "arms":collection,"comparison":compared,
        "new_source_families":"U_and_V_internal",
        "scientific_positive":False,"external_replication":False,
        "optimizer_admission":False,"old_DEV_replayed":False}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(train_only(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-cost":
        print(json.dumps(compare("TRAIN"),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--original-firstuse":
        print(json.dumps(compare("DEV"),sort_keys=True))
    elif len(sys.argv)==5 and sys.argv[1]=="--arm":
        print(json.dumps(one(sys.argv[2],int(sys.argv[3]),sys.argv[4]),sort_keys=True))
    else:raise SystemExit("T33_UNRECOGNIZED_MODE")
