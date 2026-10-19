"""T21-C one-time internal prospective learned relational motif experiment.

Fresh source generator is a separate module and cannot import this learner.
All arms get the identical source, public TRAIN cases and source generation
cost. Private DEV oracle enters only the independent scorer after prediction.
"""
from __future__ import annotations
import importlib.util
import itertools
import json
import statistics
import subprocess
import sys
import time
from collections import Counter
from pathlib import Path

HERE=Path(__file__).resolve().parent
ARMS=("learned-certified-motif-basis","trained-frequency-motif-control",
      "untrained-same-basis-fixed-motifs","equally-trained-T16B","full-public-exact")
MOTIF_ENTRIES=16
WORD_LIMIT=8
TRAIN_UPDATE_BUDGET=2048
BYTE_BUDGET=8192

def load(name,filename):
    spec=importlib.util.spec_from_file_location(name,HERE/filename)
    obj=importlib.util.module_from_spec(spec)
    sys.modules[name]=obj
    spec.loader.exec_module(obj)
    return obj

def deps():
    basis=load("t21c_exact_basis","20261019-wingless-t21-relational-bitmatrix.py")
    cert=load("t21c_frozen_cert","20261019-wingless-t16b-persistent-bindings.py")
    source=load("t21c_sources","20261019-wingless-t21c-internal-sources.py")
    return basis,cert,source

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
        kernel=ctypes.WinDLL("kernel32");psapi=ctypes.WinDLL("psapi")
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        psapi.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        psapi.GetProcessMemoryInfo.restype=wintypes.BOOL
        v=Mem();v.cb=ctypes.sizeof(v)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(v),v.cb):
            raise RuntimeError("T21C_PROCESS_RSS_UNAVAILABLE")
        return int(v.PeakWorkingSetSize)
    import resource
    value=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(value)*1024 if sys.platform.startswith("linux") else int(value)

def train_motif_candidates(train):
    """Fit only TRAIN public words; no true operators, source tags or DEV."""
    counts=Counter();updates=0
    for public in train:
        word=tuple(public["word"])
        for n in range(2,WORD_LIMIT+1):
            for start in range(len(word)-n+1):
                if updates>=TRAIN_UPDATE_BUDGET:
                    return counts,updates
                counts[word[start:start+n]]+=1
                updates+=1
    return counts,updates

def frozen_dictionary(basis,cert,train,arm):
    """Genuinely learned MDL motif selection versus frequency and untrained."""
    assert arm in ARMS[:3]
    counts,updates=train_motif_candidates(train)
    if arm=="learned-certified-motif-basis":
        # Select reusable exact representations by TRAIN-observed work saved.
        ranked=sorted(counts,key=lambda w:(
            -(counts[w]-1)*(len(w)-1),-len(w),w))
    elif arm=="trained-frequency-motif-control":
        ranked=sorted(counts,key=lambda w:(-counts[w],-len(w),w))
    else:
        # Same 16 certified entries + payload budget, but not fitted to TRAIN.
        alphabet=tuple(itertools.product(range(3),(False,True)))
        ranked=list(itertools.product(alphabet,repeat=2))
    selected=tuple(ranked[:MOTIF_ENTRIES])
    if len(selected)!=MOTIF_ENTRIES or any(len(w)>WORD_LIMIT for w in selected):
        raise RuntimeError("T21C_MOTIF_CAPACITY_MISMATCH")
    # A bitset relation stores 9 times 216 bits = 243 bytes. Count all
    # canonical metadata and motif word bytes too; size is explicitly capped.
    encoded_cost=sum(243+2*len(w)+64 for w in selected)
    if encoded_cost>BYTE_BUDGET:
        raise RuntimeError("T21C_DICTIONARY_BYTE_BUDGET_EXCEEDED")
    dictionary={}
    for word in selected:
        m=basis.construct(word)
        if not basis.independently_verify(m,cert)["verified"]:
            raise RuntimeError("T21C_PROOF_FAILURE")
        dictionary[word]=m
    return dictionary,updates,encoded_cost

def compile_word(basis,cert,word,dictionary,one):
    """Greedy frozen-motif tiling plus audited exact single-role fallback."""
    segments=[];i=0;hits=0
    while i<len(word):
        matching=[w for w in dictionary if len(w)<=len(word)-i
                  and word[i:i+len(w)]==w]
        if matching:
            best=max(matching,key=lambda w:(len(w),w))
            part=dictionary[best]
            i+=len(best);hits+=1
        else:
            token=(word[i],)
            part=one[token]
            i+=1
        segments.append(part)
    if not segments:
        part=basis.construct(())
        return part,0,0
    whole=segments[0]
    for item in segments[1:]:
        whole=basis.compose(whole,item)
    if whole.word!=tuple(word):
        raise RuntimeError("T21C_WRONG_GLOBAL_WORD")
    return whole,hits,len(segments)

def solve_certified(basis,cert,public,request,word,start,rank,compiled):
    rows=[tuple(z) for z in public]
    active={role for role,_ in cert.reduce_word(word)}
    queries=0
    while True:
        value,proof=basis.infer(compiled,rows,start)
        if proof["certified"]:
            return value,queries,True
        remaining=[j for j in active if len(rows[j])==1]
        if not remaining:
            break
        selected=max(remaining,key=lambda j:(rank[j],-j))
        rows[selected]=(rows[selected][0],request(selected))
        queries+=1
    # Exact proof only; ambiguities cannot become guesses.
    return cert.public_exact_fallback(1,rows,word,start),queries,True

def trial(arm,rep):
    assert arm in ARMS and 0<=rep<5
    begin=time.process_time_ns()
    basis,cert,source=deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert set(source.TRAIN).isdisjoint(source.DEV)
    train=[source.public_train_case(source.build(label)) for label in source.TRAIN]
    rank=cert.score_role_requests([(1,c["roles"],c["word"],c["start"]) for c in train])
    dictionary={};updates=0;dictionary_bytes=0
    singletons={}
    if arm in ARMS[:3]:
        dictionary,updates,dictionary_bytes=frozen_dictionary(basis,cert,train,arm)
        for role in range(3):
            for inv in (False,True):
                word=((role,inv),)
                singleton=basis.construct(word)
                if not basis.independently_verify(singleton,cert)["verified"]:
                    raise RuntimeError("T21C_ATOMIC_CERTIFICATE_INVALID")
                singletons[word]=singleton
    results=[];motif_hits=0;compiled_segments=0
    for label in source.DEV:
        case=source.build(label)
        if case["source_id"]!=label:
            raise RuntimeError("T21C_DEV_ID_DRIFT")
        exposed=source.observed_interface(case)
        word=exposed["word"];start=exposed["start"]
        if arm in ARMS[:3]:
            compiled,hits,segments=compile_word(basis,cert,word,dictionary,singletons)
            value,queries,sound=solve_certified(
                basis,cert,exposed["first_public_pairs"],
                exposed["query_second_pair"],word,start,rank,compiled)
            if queries!=exposed["charged_queries"]():
                raise RuntimeError("T21C_PUBLIC_REQUEST_ACCOUNTING_DRIFT")
            motif_hits+=hits;compiled_segments+=segments
        elif arm=="equally-trained-T16B":
            value,meta=cert.bound_query_solve(1,case["roles"],word,start,rank,cap=128)
            queries=meta["queries"];sound=meta["sound"]
        else:
            active={role for role,_ in cert.reduce_word(word)}
            for j in sorted(active):
                exposed["query_second_pair"](j)
            value=cert.public_exact_fallback(1,case["roles"],word,start)
            queries=exposed["charged_queries"]();sound=True
        if not sound or value!=case["private_truth"]:
            raise RuntimeError("T21C_FALSE_POSITIVE_OR_BASELINE_LOSS:"+label)
        results.append({"id":label,"correct":True,
                        "public_queries":queries})
    assert len(results)==64
    return {"schema":"wingless.t21c.internal-prospective-original-per-process.v1",
      "arm":arm,"repetition":rep,"results":results,
      "source_families":2,"train_cases":40,"dev_cases":64,
      "motif_hits":motif_hits,"segments":compiled_segments,
      "train_updates":updates,"dictionary_entries":len(dictionary),
      "dictionary_estimated_bytes":dictionary_bytes,
      "whole_cpu_ns":time.process_time_ns()-begin,
      "peak_rss_bytes":rss(),"real_service_cost":"UNMEASURED",
      "source_origin":"INTERNAL_SEPARATE_SOURCE_FILE_NOT_EXTERNAL",
      "independent_external_replication":False,
      "T16_T20_originals_used":0,"optimizer_admission":False,
      "scientific_acceptance":False}

def evaluate():
    processes=[]
    for arm in ARMS:
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                "--arm",arm,str(rep)],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T21C_CHILD_FAILED "+arm+" rep="+str(rep)+"\n"+p.stderr)
            processes.append(json.loads(p.stdout))
    result={}
    for arm in ARMS:
        cases=[p for p in processes if p["arm"]==arm]
        if len(cases)!=5 or any(c["results"]!=cases[0]["results"] for c in cases):
            raise RuntimeError("T21C_NONDETERMINISTIC_CASE_OUTCOMES")
        result[arm]={
          "correct":sum(x["correct"] for x in cases[0]["results"]),
          "public_queries":sum(x["public_queries"] for x in cases[0]["results"]),
          "source_rows":cases[0]["results"],
          "whole_cpu_five_process_ns":[x["whole_cpu_ns"] for x in cases],
          "peak_rss_five_process_bytes":[x["peak_rss_bytes"] for x in cases],
          "median_cpu_ns":int(statistics.median(x["whole_cpu_ns"] for x in cases)),
          "max_rss_bytes":max(x["peak_rss_bytes"] for x in cases),
          "motif_hits":cases[0]["motif_hits"],
          "compiled_segments":cases[0]["segments"],
          "train_updates":cases[0]["train_updates"],
          "dictionary_entries":cases[0]["dictionary_entries"],
          "dictionary_estimated_bytes":cases[0]["dictionary_estimated_bytes"]}
    learned=result[ARMS[0]]
    comparisons={arm:{
      "correctness_nonregression":learned["correct"]>=ref["correct"],
      "query_nonregression":learned["public_queries"]<=ref["public_queries"],
      "cpu_nonregression":learned["median_cpu_ns"]<=ref["median_cpu_ns"],
      "rss_nonregression":learned["max_rss_bytes"]<=ref["max_rss_bytes"]}
      for arm,ref in result.items() if arm!=ARMS[0]}
    assert all(v["correct"]==64 for v in result.values())
    return {"schema":"wingless.t21c.internal-blind-prospective.v1",
      "origin":"INTERNAL_TWO_FAMILIES_E_INTERVENTION_F_BRAID",
      "results":result,"comparisons":comparisons,
      "learned_strict_query_savings_vs_trained_T16B":
        learned["public_queries"]<result["equally-trained-T16B"]["public_queries"],
      "learned_strict_cpu_savings_vs_trained_T16B":
        learned["median_cpu_ns"]<result["equally-trained-T16B"]["median_cpu_ns"],
      "joint_acceptance":False,"joint_acceptance_reason":"REAL_REMOTE_SERVICE_COST_UNMEASURED",
      "source_origin_independent_external":False,
      "prospective_internal_first_use":True,
      "independent_replication":False,
      "scientific_acceptance":False,"optimizer_admission":False,
      "no_paid_api":True,"no_controllers":True,
      "no_old_original_source_replay":True,"no_merge":True}

def main():
    if len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(evaluate(),sort_keys=True))
    else:
        raise SystemExit("T21C_INVALID_ARGS")
if __name__=="__main__":
    main()
