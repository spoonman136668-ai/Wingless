"""T17 trained query-value against fully trained, untrained, and exact controls.

INTERNAL synthetic development only. No T16 protected cohort, no external
signed source, no scientific acceptance, no optimizer admission.
"""
from __future__ import annotations
import importlib.util
import json
import math
import os
import subprocess
import sys
import time
from pathlib import Path

HERE=Path(__file__).parent
ARMS=("learned-value","trained-T16B","untrained-same-cert","exact-full")
def load(name,path):
    spec=importlib.util.spec_from_file_location(name,path)
    obj=importlib.util.module_from_spec(spec)
    sys.modules[name]=obj
    spec.loader.exec_module(obj)
    return obj

def modules():
    k=load("t17_frozen_kernel",HERE/"20261018-wingless-t3-typed-composition.py")
    cert=load("t17_frozen_cert",HERE/"20261019-wingless-t16b-persistent-bindings.py")
    src=load("t17_internal_source",HERE/"20261019-wingless-t17-internal-synthetic-source.py")
    return k,cert,src

def public_case(fixture):
    return {"domain":fixture["domain"],"roles":fixture["roles"],
            "word":fixture["word"],"start":fixture["start"]}

def feature(case,role,cert):
    word=cert.reduce_word(case["word"])
    occurs=sum(r==role for r,_ in word)
    inverses=sum(r==role and inv for r,inv in word)
    return (case["domain"],occurs,inverses)

def outputs_for_public(k,cert,case,rows):
    word=cert.reduce_word(case["word"])
    active=sorted(set(j for j,_ in word))
    if not active:return {case["start"]}
    opts=[cert.consistent_ops(case["domain"],rows[j]) for j in active]
    chosen={}
    results=set()
    def recurse(i):
        if i==len(active):
            state=k.start(case["domain"])
            for j,inv in word:
                state=k.compose_state(state,chosen[j],case["domain"],inverse=inv)
            results.add(k.apply(state,case["start"]))
            return
        role=active[i]
        for op in opts[i]:
            chosen[role]=op
            recurse(i+1)
    recurse(0)
    return results

def train_public_value(k,cert,train):
    # Supervision is TRAIN PUBLIC second observations only, not true
    # operator IDs or a hidden task answer. A public exact certifier
    # supplies the number of reachable answers before and after reveal.
    observations={}
    for case in train:
        word=cert.reduce_word(case["word"])
        active=sorted({j for j,_ in word})
        if not active:continue
        base=[rows[:1] for rows in case["roles"]]
        before=outputs_for_public(k,cert,case,base)
        for role in active:
            changed=list(base);changed[role]=case["roles"][role]
            after=outputs_for_public(k,cert,case,changed)
            key=feature(case,role,cert)
            # Fewer possible public answers is a TRAIN-only utility label.
            reward=(len(before)-len(after))*1024//max(1,len(before))
            observations.setdefault(key,[]).append(reward)
    fitted={key:sum(values)//len(values) for key,values in sorted(observations.items())}
    return fitted, sum(len(z) for z in observations.values())

def value_decide(k,cert,case,weights):
    word=cert.reduce_word(case["word"])
    active=sorted({j for j,_ in word})
    rows=[r[:1] for r in case["roles"]]
    extra=0
    while True:
        outcome,proof=cert.shared_role_certificate(
            case["domain"],rows,case["word"],case["start"],budget=128)
        if proof["certified"]:
            return outcome,{"queries":extra,"sound":True,"certificate":proof["reason"]}
        pending=[r for r in active if len(rows[r])==1]
        if not pending:break
        # Unseen feature values revert to a source-blind fixed score.
        ranked=lambda r:(weights.get(feature(case,r,cert),0),
                          sum(t==r for t,_ in word),-r)
        nxt=max(pending,key=ranked)
        rows[nxt]=case["roles"][nxt]
        extra+=1
    full=cert.public_exact_fallback(case["domain"],case["roles"],case["word"],case["start"])
    return full,{"queries":len(active),"sound":True,"certificate":"FULL_EXACT_FALLBACK"}

def whole_process_peak_rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Counters(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
                ("PeakWorkingSetSize",ctypes.c_size_t),
                ("WorkingSetSize",ctypes.c_size_t),
                ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),
                ("QuotaPagedPoolUsage",ctypes.c_size_t),
                ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),
                ("QuotaNonPagedPoolUsage",ctypes.c_size_t),
                ("PagefileUsage",ctypes.c_size_t),
                ("PeakPagefileUsage",ctypes.c_size_t)]
        psapi=ctypes.WinDLL("psapi")
        kernel=ctypes.WinDLL("kernel32")
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        c=Counters();c.cb=ctypes.sizeof(c)
        if psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),
                                      ctypes.byref(c),ctypes.sizeof(c))==0:
            raise OSError("WINDOWS_FULL_PROCESS_PEAK_RSS_UNAVAILABLE")
        return int(c.PeakWorkingSetSize)
    import resource
    peak=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(peak)*1024 if sys.platform.startswith("linux") else int(peak)

def run_arm(arm,trial):
    assert arm in ARMS and trial in range(5)
    tick=time.process_time_ns()
    k,cert,src=modules()
    train=[public_case(src.birth_internal_fixture(i,k)) for i in src.TRAIN]
    if arm in ("learned-value","trained-T16B"):
        trained,updates=train_public_value(k,cert,train) if arm=="learned-value" else ({},0)
        frozen_rank=cert.score_role_requests([
            (c["domain"],c["roles"],c["word"],c["start"]) for c in train])
    else:
        trained={};updates=0;frozen_rank=(0,0,0,0)
    # Identical generated and scored internal DEV fixtures across controls.
    open_cases=[src.birth_internal_fixture(i,k) for i in src.DEV]
    total_queries=0
    correct=0
    baseline_correct_losses=0
    per_pattern={}
    for f in open_cases:
        case=public_case(f)
        if arm=="learned-value":
            result,meta=value_decide(k,cert,case,trained)
        elif arm=="trained-T16B":
            result,meta=cert.bound_query_solve(case["domain"],case["roles"],
                 case["word"],case["start"],frozen_rank,cap=128)
        elif arm=="untrained-same-cert":
            result,meta=cert.bound_query_solve(case["domain"],case["roles"],
                 case["word"],case["start"],(0,0,0,0),cap=128)
        else:
            result=cert.public_exact_fallback(case["domain"],case["roles"],
                                              case["word"],case["start"])
            meta={"queries":len({r for r,inv in cert.reduce_word(case["word"])}),"sound":True}
        if not meta["sound"]:raise ValueError("T17_UNSOUND_CERTIFICATE")
        total_queries+=meta["queries"]
        correct+=result==f["truth"]
        baseline_correct_losses+=result!=f["truth"]
        p=str(src.PATTERNS.index(case["word"]))
        group=per_pattern.setdefault(p,{"correct":0,"cases":0,"queries":0})
        group["correct"]+=int(result==f["truth"])
        group["cases"]+=1
        group["queries"]+=meta["queries"]
    cpu_ns=time.process_time_ns()-tick
    report={
      "schema":"wingless.t17.internal-synthetic-performance.v1",
      "arm":arm,"trial":trial,"train_cases":len(train),
      "internal_dev_cases":len(open_cases),"correct":correct,
      "full_exact_correct_answer_losses":baseline_correct_losses,
      "public_extra_queries":total_queries,
      "public_value_train_updates":updates,
      "cpu_ns_training_generation_scoring_inclusive":cpu_ns,
      "peak_whole_process_rss_bytes":whole_process_peak_rss(),
      "per_pattern":per_pattern,"same_internal_fixture_source":True,
      "independent_firstuse":False,"signed_science_accepted":False,
      "optimizer_admitted":False
    }
    assert correct==64,"T17_HARD_FALSE_POSITIVE_OR_BASELINE_LOSS"
    return report

def main():
    if len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(run_arm(sys.argv[2],int(sys.argv[3])),sort_keys=True))
        return
    if len(sys.argv)!=1:raise SystemExit("T17_ARGS")
    batch=[]
    for arm in ARMS:
        for trial in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                             "--arm",arm,str(trial)],capture_output=True,text=True,check=True)
            batch.append(json.loads(p.stdout))
    from statistics import median
    metrics={}
    for arm in ARMS:
        rows=[r for r in batch if r["arm"]==arm]
        for k in ("correct","public_extra_queries","full_exact_correct_answer_losses"):
            assert len({r[k] for r in rows})==1
        metrics[arm]={
          "correct":rows[0]["correct"],
          "queries":rows[0]["public_extra_queries"],
          "train_updates":rows[0]["public_value_train_updates"],
          "median_complete_cpu_ns":int(median(r["cpu_ns_training_generation_scoring_inclusive"] for r in rows)),
          "max_complete_process_peak_rss_bytes":max(r["peak_whole_process_rss_bytes"] for r in rows),
          "median_complete_process_peak_rss_bytes":int(median(r["peak_whole_process_rss_bytes"] for r in rows)),
          "cpu_five_process_measurements_ns":[r["cpu_ns_training_generation_scoring_inclusive"] for r in rows],
          "rss_five_process_measurements_bytes":[r["peak_whole_process_rss_bytes"] for r in rows],
          "per_pattern":rows[0]["per_pattern"]
        }
    candidate=metrics["learned-value"]
    checks={}
    for control in ARMS[1:]:
        c=metrics[control]
        checks[control]={
         "strictly_fewer_extra_queries":candidate["queries"]<c["queries"],
         "complete_cpu_nonregression":candidate["median_complete_cpu_ns"]<=c["median_complete_cpu_ns"],
         "peak_process_rss_nonregression":
             candidate["max_complete_process_peak_rss_bytes"]<=c["max_complete_process_peak_rss_bytes"],
         "correctness_nonregression":candidate["correct"]>=c["correct"]}
    eligible=all(all(v.values()) for v in checks.values())
    result={"schema":"wingless.t17.internal-training-versus-strong-controls.v1",
      "label":"INTERNAL_SYNTHETIC_DIAGNOSTIC_NOT_INDEPENDENT_SCIENCE",
      "source":"T17 disjoint self-authored synthetic 24 TRAIN 64 DEV",
      "source_born_externally":False,"protected_originals_consumed":0,
      "results":metrics,"control_checks":checks,
      "candidate_strict_joint_resource_eligible":eligible,
      "scientific_acceptance":False,
      "optimizer_admission":False,"new_controller":False,
      "genuine_public_query_service_metered":False,
      "no_original_firstuse":True}
    print(json.dumps(result,sort_keys=True))
if __name__=="__main__":main()
