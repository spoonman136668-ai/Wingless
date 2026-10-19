"""T23 TRAIN-learned role-projected prefix-state DAG for four shared-causal queries.

Exactness: cache keys retain the same observable input and the actual global
operator for every role appearing in each prefix. Assignments to other roles
may be forgotten only because that prefix cannot depend on them. No approximate
candidate answer, T21 bitmatrix, T22 identity shortcut or private source label.
"""
from __future__ import annotations
import importlib.util,itertools,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path

ROOT=Path(__file__).resolve().parent
ARMS=("learned-prefix-DAG","trained-frequency-prefix-DAG",
      "dynamic-exact-prefix-trie","fully-trained-T16B-batched","full-public-exact-batched")
CAP=32
TRAIN_UPDATES=2048
def load(n,f):
    spec=importlib.util.spec_from_file_location(n,ROOT/f)
    m=importlib.util.module_from_spec(spec);sys.modules[n]=m;spec.loader.exec_module(m)
    return m
def deps():
    return (load("t23_cert","20261019-wingless-t16b-persistent-bindings.py"),
            load("t23_source","20261019-wingless-t23-new-multitask-source.py"))
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
        lib=ctypes.WinDLL("kernel32");ps=ctypes.WinDLL("psapi")
        lib.GetCurrentProcess.restype=wintypes.HANDLE
        ps.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        ps.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not ps.GetProcessMemoryInfo(lib.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T23_RSS_UNAVAILABLE")
        return int(m.PeakWorkingSetSize)
    import resource
    value=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(value)*1024 if sys.platform.startswith("linux") else int(value)

def prefixes(words):
    # A prefix is a binding-preserving relational state, not an identity.
    for w in words:
        for i in range(2,len(w)+1):
            yield tuple(w[:i])

def train_cache(train,arm):
    if arm not in ARMS[:2]:return frozenset(),0
    counts=Counter();updates=0
    for case in train:
        for prefix in prefixes(case["words"]):
            if updates>=TRAIN_UPDATES:break
            counts[prefix]+=1;updates+=1
        if updates>=TRAIN_UPDATES:break
    if arm=="learned-prefix-DAG":
        # Public-only score: likely native computation saved per cache byte
        # with capacity for unresolved roles retained in the projection.
        ordered=sorted(counts,key=lambda x:(
             -((counts[x]-1)*max(1,len(x)-1) * (4-len({j for j,_ in x})))/(24+2*len(x)),
             -counts[x],-len(x),x))
    else:
        ordered=sorted(counts,key=lambda x:(-counts[x],-len(x),x))
    return frozenset(ordered[:CAP]),updates

def dynamic_cache(words):
    counts=Counter(prefixes(words))
    ranked=sorted(counts,key=lambda x:(-(counts[x]-1)*len(x),-len(x),x))
    return frozenset(ranked[:CAP])

def exact_batch(cert,visible,request,words,starts,rank,arm,global_cache):
    rows=[tuple(pair) for pair in visible]
    used=tuple(sorted({r for word in words for r,_ in cert.reduce_word(word)}))
    purchases=0;memo_hits=0;memo_entries_peak=0;assignments_visited=0
    while True:
        answer_sets=[set() for _ in words]
        if arm in ("fully-trained-T16B-batched","full-public-exact-batched"):
            all_certified=True
            for i,(word,start) in enumerate(zip(words,starts)):
                value,proof=cert.shared_role_certificate(1,rows,word,start,budget=4096)
                if not proof["certified"]:
                    all_certified=False
                else:
                    answer_sets[i].add(value)
            if all_certified:
                return tuple(next(iter(v)) for v in answer_sets),purchases,memo_hits,memo_entries_peak,assignments_visited
        else:
            # Exact global binding enumeration per public observation;
            # projection drops only roles not used in the PREFIX, not those
            # occurring later in the word, preserving all repeated bindings.
            options=[cert.consistent_ops(1,rows[r]) for r in range(3)]
            if any(not v for v in options):
                raise RuntimeError("T23_INCONSISTENT_PUBLIC")
            memo={};assignment_list=tuple(itertools.product(*options))
            cacheable=(dynamic_cache(words) if arm=="dynamic-exact-prefix-trie" else global_cache)
            for assignment in assignment_list:
                assignments_visited+=1
                for i,(word,start) in enumerate(zip(words,starts)):
                    value=start
                    for step,(role,inv) in enumerate(word):
                        prefix=tuple(word[:step+1])
                        key=None
                        if prefix in cacheable:
                            relevant=tuple(sorted({r for r,_ in prefix}))
                            key=(prefix,start,tuple(assignment[r] for r in relevant))
                            if key in memo:
                                value=memo[key];memo_hits+=1
                                continue
                        value=cert.apply_op(1,assignment[role],value,inv)
                        if key is not None:
                            memo[key]=value
                    answer_sets[i].add(value)
            memo_entries_peak=max(memo_entries_peak,len(memo))
            if all(len(x)==1 for x in answer_sets):
                return tuple(next(iter(x)) for x in answer_sets),purchases,memo_hits,memo_entries_peak,assignments_visited
        remaining=[r for r in used if len(rows[r])==1]
        if not remaining:
            raise RuntimeError("T23_FULL_PUBLIC_NOT_CERTIFIED")
        priority=max(remaining,key=lambda r:(rank[r],-r))
        rows[priority]=(rows[priority][0],request(priority));purchases+=1

def trial(arm,rep):
    assert arm in ARMS and 0<=rep<5
    startcpu=time.process_time_ns()
    cert,source=deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64
    assert not set(source.TRAIN)&set(source.DEV)
    train=[source.train_public(source.build(name)) for name in source.TRAIN]
    rank=cert.score_role_requests([(1,c["role_pairs"],word,start)
                  for c in train for word,start in zip(c["words"],c["starts"])])
    learned_cache,updates=train_cache(train,arm)
    score_rows=[];cache_hits=0;cache_entries=0;bindings=0
    for name in source.DEV:
        item=source.build(name)
        view=source.dev_public(item)
        if arm=="full-public-exact-batched":
            for role in sorted({r for w in view["words"] for r,_ in cert.reduce_word(w)}):
                view["request"](role)
            # Include initial full disclosure in charged public queries,
            # but exact inference uses all available second pairs.
            rows=item["role_pairs"]
            solutions=[]
            for word,start in zip(view["words"],view["starts"]):
                solutions.append(cert.public_exact_fallback(1,rows,word,start))
            answers=tuple(solutions);purchased=view["charged"]()
            hit=0;peak=0;enumerated=0
        else:
            answers,purchased,hit,peak,enumerated=exact_batch(
                 cert,view["first"],view["request"],view["words"],
                 view["starts"],rank,arm,learned_cache)
        if tuple(answers)!=item["private_truths"] or purchased!=view["charged"]():
            raise RuntimeError("T23_HARD_FALSE_POSITIVE_OR_BAD_QUERY:"+name)
        score_rows.append({"id":name,"correct_answers":4,"purchased_public_pairs":purchased})
        cache_hits+=hit;cache_entries=max(cache_entries,peak);bindings+=enumerated
    return {"schema":"wingless.t23.internal-multiquery-original-process.v1",
        "arm":arm,"repetition":rep,"rows":score_rows,
        "TRAIN_cases":40,"DEV_cases":64,"DEV_questions":256,
        "trained_prefixes":len(learned_cache),"TRAIN_updates":updates,
        "prefix_cache_hits":cache_hits,"peak_prefix_state_entries":cache_entries,
        "exact_assignment_iterations":bindings,
        "complete_process_cpu_ns":time.process_time_ns()-startcpu,
        "peak_process_rss_bytes":rss(),"source_origin":"INTERNAL_NOT_EXTERNAL",
        "prior_DEV_reuse":False,"optimizer_admission":False}

def run():
    processes=[]
    for arm in ARMS:
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                "--arm",arm,str(rep)],capture_output=True,text=True)
            if p.returncode:
                raise RuntimeError("T23_CHILD_FAILURE "+arm+" "+str(rep)+"\n"+p.stderr)
            processes.append(json.loads(p.stdout))
    output={}
    for arm in ARMS:
        records=[x for x in processes if x["arm"]==arm]
        assert len(records)==5 and all(r["rows"]==records[0]["rows"] for r in records)
        output[arm]={"correct_questions":sum(x["correct_answers"] for x in records[0]["rows"]),
           "correct_cases":len(records[0]["rows"]),
           "paid_queries":sum(x["purchased_public_pairs"] for x in records[0]["rows"]),
           "per_case":records[0]["rows"],
           "trained_prefixes":records[0]["trained_prefixes"],
           "TRAIN_updates":records[0]["TRAIN_updates"],
           "prefix_cache_hits":records[0]["prefix_cache_hits"],
           "peak_prefix_state_entries":records[0]["peak_prefix_state_entries"],
           "assignment_iterations":records[0]["exact_assignment_iterations"],
           "five_whole_cpu_ns":[r["complete_process_cpu_ns"] for r in records],
           "five_peak_rss_bytes":[r["peak_process_rss_bytes"] for r in records],
           "median_whole_cpu_ns":int(statistics.median(r["complete_process_cpu_ns"] for r in records)),
           "max_peak_rss_bytes":max(r["peak_process_rss_bytes"] for r in records)}
    assert all(v["correct_questions"]==256 and v["correct_cases"]==64 for v in output.values())
    cand=output[ARMS[0]]
    comparison={a:{
       "all_answers_correct":cand["correct_questions"]>=ref["correct_questions"],
       "queries_no_worse":cand["paid_queries"]<=ref["paid_queries"],
       "cpu_no_worse":cand["median_whole_cpu_ns"]<=ref["median_whole_cpu_ns"],
       "rss_no_worse":cand["max_peak_rss_bytes"]<=ref["max_peak_rss_bytes"]}
       for a,ref in output.items() if a!=ARMS[0]}
    return {"schema":"wingless.t23.source-blind-internal-multiquery-comparison.v1",
       "source_families":"I_and_J_new_internal","arms":output,"comparison":comparison,
       "strict_learned_query_gain_vs_dynamic_exact":cand["paid_queries"]<output["dynamic-exact-prefix-trie"]["paid_queries"],
       "strict_learned_CPU_gain_vs_dynamic_exact":cand["median_whole_cpu_ns"]<output["dynamic-exact-prefix-trie"]["median_whole_cpu_ns"],
       "scientific_joint_acceptance":False,"real_remote_query_cost":"UNMEASURED",
       "independently_replicated":False,"optimizer_admission":False,
       "old_DEV_reopened":False,"protected_originals_replayed":False}

if __name__=="__main__":
    if len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    elif len(sys.argv)==1:print(json.dumps(run(),sort_keys=True))
    else:raise SystemExit("T23_INVALID_ARGUMENTS")
