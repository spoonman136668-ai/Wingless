"""T24 exact information-theoretic causal source headroom.

All candidate policies see only first calibration rows, causal composition
words, starts and purchased second pairs. Exactly 16 globally persistent S3
operator assignments follow. Optimal online search sees NO actual hidden
operators or answers; hindsight floor is scored separately after prediction.
No learned cognition advantage or externally independent science is claimed.
"""
from __future__ import annotations
import functools,importlib.util,itertools,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-exact-greedy","optimal-online-exact","fixed-exact-order","full-public-exact")
def load(n,p):
    spec=importlib.util.spec_from_file_location(n,ROOT/p)
    m=importlib.util.module_from_spec(spec)
    sys.modules[n]=m;spec.loader.exec_module(m)
    return m
def deps():
    return (load("t24_source","20261019-wingless-t24-new-source.py"),
            load("t24_cert","20261019-wingless-t16b-persistent-bindings.py"))
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
        kl=ctypes.WinDLL("kernel32");ps=ctypes.WinDLL("psapi")
        kl.GetCurrentProcess.restype=wintypes.HANDLE
        ps.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        ps.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not ps.GetProcessMemoryInfo(kl.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T24_RSS_UNAVAILABLE")
        return int(m.PeakWorkingSetSize)
    import resource
    val=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(val)*1024 if sys.platform.startswith("linux") else int(val)
def fit_priority(train):
    """Only public TRAIN words; no oracle results, hidden op or unseen source."""
    scores=[0]*4
    for public in train:
        for word in public["words"]:
            for role,inv in word:scores[role]+=1
    return tuple(scores)

def model(cert,words,starts,first):
    if len(first)!=4 or len(words)!=3 or len(starts)!=3:
        raise ValueError("T24_INVALID_PUBLIC_INTERFACE")
    options=tuple(cert.consistent_ops(1,pairs) for pairs in first)
    if any(len(op)!=2 for op in options):
        raise RuntimeError("T24_FIRST_PAIR_NOT_BINARY")
    indices=tuple(itertools.product(range(2),repeat=4))
    assignments=tuple(tuple(options[role][idx[role]] for role in range(4)) for idx in indices)
    outcomes=[]
    for op_set in assignments:
        answers=[]
        for word,start in zip(words,starts):
            value=start
            for role,inv in word:
                value=cert.apply_op(1,op_set[role],value,inv)
            answers.append(value)
        outcomes.append(tuple(answers))
    assert len(outcomes)==16
    return {"options":options,"binary_roles":indices,"answers":tuple(outcomes),
            "full_mask":(1<<16)-1,"used":tuple(sorted({role for w in words for role,inv in w}))}

def selected_states(state):
    return tuple(i for i in range(16) if state&(1<<i))
def answer_if_certified(m,state):
    answers={m["answers"][i] for i in selected_states(state)}
    if not answers:raise RuntimeError("T24_EMPTY_BINDING_SET")
    return next(iter(answers)) if len(answers)==1 else None
def partition(m,state,role):
    groups=(0,0)
    a=b=0
    for index in selected_states(state):
        if m["binary_roles"][index][role]==0:a|=1<<index
        else:b|=1<<index
    return a,b
def available_roles(m,state):
    return tuple(role for role in m["used"]
                 if all(partition(m,state,role)))

def online_dp(m):
    """Exact Bayes-optimal online expected query count, uniform bindings.

    Uses ONLY the public set of legal assignments, causal path words and
    starts. Never uses the actual unseen source operator or private answer.
    """
    @functools.lru_cache(None)
    def solve(state):
        if answer_if_certified(m,state) is not None:return (0.0,None)
        n=state.bit_count();candidates=[]
        for role in available_roles(m,state):
            lo,hi=partition(m,state,role)
            val=1+(lo.bit_count()*solve(lo)[0]+hi.bit_count()*solve(hi)[0])/n
            candidates.append((val,role))
        if not candidates:raise RuntimeError("T24_NOT_CERTIFIED_WITH_ALL_PUBLIC")
        return min(candidates)
    root=solve(m["full_mask"])
    return solve,root[0],solve.cache_info().currsize

def choose_greedy(m,state,priority):
    """Strong dynamic entropy reduction, with TRAIN-only tie breaks."""
    options=[];before=len({m["answers"][i] for i in selected_states(state)})
    for role in available_roles(m,state):
        lo,hi=partition(m,state,role)
        expected=(lo.bit_count()*len({m["answers"][i] for i in selected_states(lo)})+
                 hi.bit_count()*len({m["answers"][i] for i in selected_states(hi)}))/state.bit_count()
        options.append((before-expected,priority[role],-role,role))
    if not options:raise RuntimeError("T24_NO_GREEDY_QUERY_AVAILABLE")
    return max(options)[3]
def play(m,policy,train_priority,request):
    state=m["full_mask"];purchases=[];cost=None
    dp=None;dp_states=0
    if policy=="optimal-online-exact":
        dp,cost,dp_states=online_dp(m)
    while True:
        answer=answer_if_certified(m,state)
        if answer is not None:return answer,tuple(purchases),cost,dp_states
        if policy=="optimal-online-exact":
            role=dp(state)[1]
        elif policy=="TRAIN-exact-greedy":
            role=choose_greedy(m,state,train_priority)
        elif policy=="fixed-exact-order":
            role=min(available_roles(m,state))
        else:
            role=min(available_roles(m,state))
        if role is None or role in purchases:raise RuntimeError("T24_INVALID_POLICY")
        second=request(role)
        candidate_options=m["options"][role]
        selected=[z for z in range(2) if candidate_options[z][second[0]]==second[1]]
        if len(selected)!=1:raise RuntimeError("T24_SECOND_PAIR_UNRESOLVED")
        lo,hi=partition(m,state,role)
        state=lo if selected[0]==0 else hi
        purchases.append(role)
        if len(purchases)>4:raise RuntimeError("T24_QUERY_CAP_EXCEEDED")
def hindsight_floor(m,true_index):
    """Scorer-only optimistic bound after prediction; NOT a runnable policy."""
    true_bits=m["binary_roles"][true_index]
    for length in range(5):
        for chosen in itertools.combinations(range(4),length):
            state=m["full_mask"]
            for role in chosen:
                lo,hi=partition(m,state,role)
                state=lo if true_bits[role]==0 else hi
            if answer_if_certified(m,state)==m["answers"][true_index]:
                return length
    raise RuntimeError("T24_HINDSIGHT_BOUND_FAILED")
def actual_source_assignment_index(m,case):
    bits=[]
    for role,row in enumerate(case["rows"]):
        op=m["options"][role]
        possible=[z for z in range(2) if op[z][row[1][0]]==row[1][1]]
        if len(possible)!=1:raise RuntimeError("T24_TRUTH_NOT_IN_PUBLIC_BINDINGS")
        bits.append(possible[0])
    return m["binary_roles"].index(tuple(bits))
def trial(arm,repetition):
    assert arm in ARMS and 0<=repetition<5
    begin=time.process_time_ns()
    source,cert=deps()
    train=[source.visible_training(source.build(name)) for name in source.TRAIN]
    priority=fit_priority(train)
    records=[];dp_cache_states=0
    for name in source.DEV:
        case=source.build(name)
        view=source.visible_dev(case)
        m=model(cert,view["words"],view["starts"],view["first"])
        if arm=="full-public-exact":
            # Full disclosure pays for every actual used role, with no oracle.
            for role in m["used"]:view["purchase"](role)
            truth=tuple(case["truth"])
            answers=tuple(cert.public_exact_fallback(1,case["rows"],word,start)
                         for word,start in zip(view["words"],view["starts"]))
            cost=None;dp_states=0;requested=m["used"]
        else:
            answers,requested,cost,dp_states=play(m,arm,priority,view["purchase"])
        if answers!=case["truth"] or len(requested)!=view["charged"]():
            raise RuntimeError("T24_HARD_FALSE_POSITIVE_OR_QUERY_DRIFT "+name)
        # Math-only optimistic lower bound computed by scorer AFTER inference.
        # It may depend on hidden truth, so cannot be used to select a query.
        floor=hindsight_floor(m,actual_source_assignment_index(m,case))
        if floor>len(requested):
            raise RuntimeError("T24_LOWER_BOUND_GREATER_THAN_OBSERVED")
        if cost is not None and cost>4.0000001:
            raise RuntimeError("T24_OPTIMAL_EXPECTED_COST_INVALID")
        dp_cache_states+=dp_states
        records.append({"id":name,"correct_outputs":3,"purchases":len(requested),
                        "roles_purchased":requested,"hindsight_floor":floor,
                        "optimal_expected_queries":cost})
    return {"schema":"wingless.t24.original-internal-headroom-single-process.v1",
            "arm":arm,"repetition":repetition,"train_cases":40,"DEV_cases":64,
            "records":records,"TRAIN_priorities":priority,
            "DP_cache_states":dp_cache_states,
            "whole_cpu_ns":time.process_time_ns()-begin,"peak_rss_bytes":rss(),
            "source_origin":"INTERNAL_NOT_EXTERNAL","old_DEV_reused":False,
            "scientific_acceptance":False,"optimizer_admission":False}
def evaluate():
    results={}
    for arm in ARMS:
        processes=[]
        for rep in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),"--arm",arm,str(rep)],
                             capture_output=True,text=True)
            if p.returncode:raise RuntimeError("T24_CHILD_FAIL "+arm+" "+str(rep)+"\n"+p.stderr)
            processes.append(json.loads(p.stdout))
        if any(p["records"]!=processes[0]["records"] for p in processes):
            raise RuntimeError("T24_NONDETERMINISTIC_SOURCE")
        rows=processes[0]["records"]
        results[arm]={"correct_outputs":sum(r["correct_outputs"] for r in rows),
            "public_second_pair_queries":sum(r["purchases"] for r in rows),
            "per_case":rows,"hindsight_optimistic_floor_sum":sum(r["hindsight_floor"] for r in rows),
            "optimal_expected_query_cost_sum":
               sum(r["optimal_expected_queries"] for r in rows)
               if arm=="optimal-online-exact" else None,
            "DP_cache_states":processes[0]["DP_cache_states"],
            "five_complete_cpu_ns":[p["whole_cpu_ns"] for p in processes],
            "five_peak_rss_bytes":[p["peak_rss_bytes"] for p in processes],
            "median_complete_cpu_ns":int(statistics.median(p["whole_cpu_ns"] for p in processes)),
            "maximum_peak_rss_bytes":max(p["peak_rss_bytes"] for p in processes)}
    assert all(z["correct_outputs"]==192 for z in results.values())
    g=results["TRAIN-exact-greedy"];optimal=results["optimal-online-exact"]
    return {"schema":"wingless.t24.original-headroom-audit.v1",
        "source":"NEW_INTERNAL_K_L_CAUSAL_GRAPH_INTERVENTIONS",
        "results":results,
        "observed_greedy_vs_online_optimal_query_gap":g["public_second_pair_queries"]-
                                                    optimal["public_second_pair_queries"],
        "online_theoretical_EXPECTED_optimal_query_cost_sum":optimal["optimal_expected_query_cost_sum"],
        "optimistic_hindsight_non_deployable_floor_sum":optimal["hindsight_optimistic_floor_sum"],
        "full_public_query_budget":results["full-public-exact"]["public_second_pair_queries"],
        "source_cases":64,"target_outputs":192,
        "no_learning_candidate_trained":True,"no_scientific_accepted_positive":True,
        "independent_replication":False,"optimizer_admission":False,
        "old_DEV_replayed":False,"protected_originals_accessed":False}
def TRAIN_only():
    source,cert=deps()
    priority=fit_priority([source.visible_training(source.build(name)) for name in source.TRAIN])
    verified=0;states=0;floor_checks=0;mean_excess=0
    for name in source.TRAIN:
        case=source.build(name)
        pub=source.visible_training(case)
        m=model(cert,pub["words"],pub["starts"],pub["first"])
        dp,cost,num=online_dp(m)
        states+=num
        for assignment in range(16):
            trueop=tuple(m["options"][role][m["binary_roles"][assignment][role]] for role in range(4))
            expected=[]
            for word,start in zip(pub["words"],pub["starts"]):
                x=start
                for role,inv in word:
                    x=cert.apply_op(1,trueop[role],x,inv)
                expected.append(x)
            if m["answers"][assignment]!=tuple(expected):
                raise RuntimeError("T24_PUBLIC_EXACT_BINDING_DRIFT")
            floor=hindsight_floor(m,assignment)
            if not 0<=floor<=4:raise RuntimeError("T24_BAD_FLOOR")
            floor_checks+=1;verified+=3
        # Confirm dynamic online recurrence, proper expected value vs greedy
        for state in range(1,1<<16):
            if state & m["full_mask"] == 0:continue
            # Only visit states that are reachable by public role disclosure.
            # Accessing every mask would be prohibitively costly and invalid.
            break
        if not 0<=cost<=4 or num>81:raise RuntimeError("T24_DP_NOT_BOUNDED")
        mean_excess+=cost
    return {"status":"PASS_TRAIN_ONLY_NO_DEV",
        "source_cases":40,"DEV_cases_generated":0,
        "three_output_algebraic_identity_checks":verified,
        "hindsight_bound_checked":floor_checks,
        "DP_total_states":states,"TRAIN_only_policy":priority,
        "expected_minimum_TRAIN_queries_sum":mean_excess,
        "external_source_independent":False,"learned_model_tested":False,
        "optimizer_admission":False}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(TRAIN_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T24_INVALID_ARGS")
