"""T43: source-free learned causal-intervention selection, no scientific DEV.

One-step objective is the exact expected reduction in Bayes Brier risk on
two disjoint hidden TARGET rule bits, excluding the observed intervention
bit itself. Compare TRAIN-learned conditional covariance and mutual-
information policies with both optimized and memoized exact Bayes oracle.
All policies see the same public 96 fully disclosed meta rule observations
and the same 48 independent source-free partial-rule validation histories.
"""
from __future__ import annotations
import functools,json,math,random,statistics,subprocess,sys,time
from pathlib import Path
ARMS=("TRAIN-covariance-approximate-one-step-policy",
      "optimized-exact-covariance-one-step-optimal",
      "cached-exact-one-step-optimal",
      "TRAIN-mutual-information-pairwise-control",
      "fixed-first-action-control")
TRAIN_COUNT=96
VALID_COUNT=48
def compose(g,h,op):
    rule=0
    for p in range(8):
        a,b,c=p&1,(p>>1)&1,(p>>2)&1
        x=(g>>(a+2*b))&1
        y=(h>>(b+2*c))&1
        bit=x^y if op==0 else (x&y if op==1 else x|y)
        rule|=bit<<p
    return rule
PROGRAMS=tuple(compose(g,h,op) for g in range(16) for h in range(16) for op in range(3))
multiplicity=[0]*256
for rule in PROGRAMS:multiplicity[rule]+=1
PRIOR=tuple(.5*count/768+.5/256 for count in multiplicity)
assert len(PROGRAMS)==768 and abs(sum(PRIOR)-1)<1e-12
def fixture(part,index):
    if part not in ("TRAIN","VALID") or index not in range(TRAIN_COUNT if part=="TRAIN" else VALID_COUNT):
        raise RuntimeError("T43_UNFROZEN_CASE_ID")
    r=random.Random((43002026 if part=="TRAIN" else 43012026)+index*17011)
    hidden_regime=r.randrange(2)
    law=r.choice(PROGRAMS) if hidden_regime==0 else r.randrange(256)
    if part=="TRAIN":
        return {"id":f"t43-public-train-{index:03d}","public_rule":law,
                "_hidden_source_regime":hidden_regime}
    order=tuple(r.sample(range(8),8))
    return {"id":f"t43-sourcefree-valid-{index:03d}",
            "public":tuple((p,(law>>p)&1) for p in order[:2]),
            "targets":order[2:4],"actions":order[4:],
            "_private_rule":law,"_private_regime":hidden_regime}
def train():
    # SOURCE-BLIND: read PUBLIC TRAIN laws only, never source regime labels.
    laws=[fixture("TRAIN",i)["public_rule"] for i in range(TRAIN_COUNT)]
    mean=[sum((r>>p)&1 for r in laws)/TRAIN_COUNT for p in range(8)]
    covariance=[[sum((((r>>i)&1)-mean[i])*(((r>>j)&1)-mean[j]) for r in laws)/TRAIN_COUNT
                 for j in range(8)] for i in range(8)]
    mi=[[0.]*8 for _ in range(8)]
    for i in range(8):
        for j in range(8):
            count=[[.5,.5],[.5,.5]]
            for law in laws:count[(law>>i)&1][(law>>j)&1]+=1
            total=sum(sum(row) for row in count)
            p=[sum(row)/total for row in count]
            q=[sum(count[a][b] for a in range(2))/total for b in range(2)]
            mi[i][j]=sum((count[a][b]/total)*math.log2((count[a][b]/total)/(p[a]*q[b]))
                         for a in range(2) for b in range(2))
    return {"cov":covariance,"mi":mi,"train_cases":len(laws)}
@functools.lru_cache(None)
def posterior(public):
    compatible=[(law,w) for law,w in enumerate(PRIOR)
                if all(((law>>bit)&1)==value for bit,value in public)]
    mass=sum(w for _,w in compatible)
    if mass<=0:raise RuntimeError("T43_TRUE_CAUSAL_RULE_EXCLUDED")
    return tuple((law,w/mass) for law,w in compatible)
def gain(items,action,targets):
    pa=sum(w*((law>>action)&1) for law,w in items)
    variance=pa*(1-pa)
    if variance<=1e-15:return 0.
    total=0.
    for target in targets:
        pt=sum(w*((law>>target)&1) for law,w in items)
        pair=sum(w*(((law>>action)&1)*((law>>target)&1)) for law,w in items)
        total+=(pair-pa*pt)**2/(2*variance)
    return total
def direct_exhaustive_gain(items,action,targets):
    before=sum((lambda p:p*(1-p))(sum(w*((law>>t)&1) for law,w in items))
               for t in targets)/2
    after=0.
    for outcome in (0,1):
        subset=[(law,w) for law,w in items if ((law>>action)&1)==outcome]
        mass=sum(w for _,w in subset)
        if mass<=1e-15:continue
        for t in targets:
            prob=sum(w*((law>>t)&1) for law,w in subset)/mass
            after+=mass*prob*(1-prob)/2
    return before-after
def conditional_cov(model,public):
    cov=model["cov"]
    i,j=[bit for bit,output in public]
    a,b,c=cov[i][i]+.025,cov[i][j],cov[j][j]+.025
    determinant=a*c-b*b
    if determinant<=1e-12:raise RuntimeError("T43_TRAINED_COVARIANCE_NOT_POSITIVE")
    inverse=((c/determinant,-b/determinant),
             (-b/determinant,a/determinant))
    def conditional(x,y):
        return cov[x][y]-sum(cov[x][(i,j)[k]]*inverse[k][m]*cov[(i,j)[m]][y]
                             for k in range(2) for m in range(2))
    return conditional
def choose(arm,case,model,memo):
    actions=case["actions"]
    targets=case["targets"]
    if arm==ARMS[4]:return actions[0]
    if arm==ARMS[3]:score=lambda action:sum(model["mi"][action][t] for t in targets)
    elif arm==ARMS[0]:
        cond=conditional_cov(model,case["public"])
        def score(action):
            var=cond(action,action)
            return sum(cond(action,t)**2/(2*var) for t in targets) if var>1e-12 else 0.
    else:
        key=(case["public"],targets,actions)
        if arm==ARMS[2] and key in memo:return memo[key]
        items=posterior(case["public"])
        selected=max(actions,key=lambda a:(gain(items,a,targets),-actions.index(a)))
        if arm==ARMS[2]:memo[key]=selected
        return selected
    return max(actions,key=lambda a:(score(a),-actions.index(a)))
def evaluate(arm,model):
    total=optimal=baseline=loss=0.
    certified=abstentions=incorrect=0
    rows=[];memo={}
    for index in range(VALID_COUNT):
        case=fixture("VALID",index)
        chosen=choose(arm,case,model,memo)
        if chosen not in case["actions"]:raise RuntimeError("T43_NOT_AN_ALLOWED_CAUSAL_INTERVENTION")
        items=posterior(case["public"])
        utility={a:gain(items,a,case["targets"]) for a in case["actions"]}
        # Separately verify exact covariance identity using all outcome paths.
        for action in case["actions"]:
            direct=direct_exhaustive_gain(items,action,case["targets"])
            if abs(utility[action]-direct)>1e-10:
                raise RuntimeError("T43_BAYES_EXPECTED_BRIER_COVARIANCE_PROOF_FAILED")
        best=max(utility.values())
        chosen_gain=utility[chosen]
        if chosen_gain>best+1e-10:raise RuntimeError("T43_LEARNED_POLICY_OUTPERFORMS_EXACT_UPPER_BOUND")
        total+=chosen_gain;optimal+=best;baseline+=utility[case["actions"][0]]
        # Observe the selected causal-pattern bit, THEN predict disjoint targets.
        rule=case["_private_rule"]
        public_after=case["public"]+((chosen,(rule>>chosen)&1),)
        updated=posterior(public_after)
        if rule not in [law for law,w in updated]:
            raise RuntimeError("T43_SOURCE_TRUE_RULE_LOST_AFTER_INTERVENTION")
        for t in case["targets"]:
            p=sum(w*((law>>t)&1) for law,w in updated)
            y=(rule>>t)&1
            loss+=(p-y)**2
            possible={((law>>t)&1) for law,w in updated}
            if len(possible)==1:
                certified+=1
                if next(iter(possible))!=y:incorrect+=1
            else:abstentions+=1
        rows.append({"id":case["id"],"public":case["public"],
                     "action":chosen,"regret":round(best-chosen_gain,12)})
    return {"total_selected_expected_gain":total,"total_exact_optimal_gain":optimal,
            "total_first_action_expected_gain":baseline,
            "fraction_exact_optimal_utility":total/optimal if optimal>0 else 1.,
            "realized_disjoint_target_Brier":loss/(VALID_COUNT*2),
            "correct_certified":certified,"sound_abstentions":abstentions,
            "hard_false_certifications":incorrect,"trace":rows}
def peak_rss():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class M(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
               ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
               ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
               ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
               ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        kernel=ctypes.WinDLL("kernel32",use_last_error=True)
        psapi=ctypes.WinDLL("psapi",use_last_error=True)
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        psapi.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(M),wintypes.DWORD]
        psapi.GetProcessMemoryInfo.restype=wintypes.BOOL
        mem=M();mem.cb=ctypes.sizeof(mem)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(mem),mem.cb):
            raise RuntimeError("T43_NATIVE_WINDOWS_RSS_UNAVAILABLE")
        return int(mem.PeakWorkingSetSize)
    import resource
    m=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(m)*1024 if sys.platform.startswith("linux") else int(m)
def ablation():
    model=train()
    results={arm:evaluate(arm,model) for arm in (ARMS[0],ARMS[1],ARMS[3],ARMS[4])}
    candidate=results[ARMS[0]]
    changed=sum(a["action"]!=b["action"] for a,b in zip(candidate["trace"],results[ARMS[4]]["trace"]))
    if candidate["hard_false_certifications"] or results[ARMS[1]]["hard_false_certifications"]:
        raise RuntimeError("T43_FALSE_HARD_CERTIFICATION")
    if results[ARMS[1]]["fraction_exact_optimal_utility"]<1-1e-10:
        raise RuntimeError("T43_OPTIMIZED_EXACT_CONTROL_NOT_OPTIMAL")
    return {"schema":"wingless.T43.frozen-learned-active-intervention-ablation.v1",
       "TRAIN_public_cases":TRAIN_COUNT,"sourcefree_validation_cases":VALID_COUNT,
       "scientific_DEV_generated":0,
       "learned_vs_disabled_first_action_changed_cases":changed,
       "policy":{arm:{k:v for k,v in res.items() if k!="trace"} for arm,res in results.items()},
       "original_learned_case_trace":candidate["trace"],
       "exact_covariance_equal_exhaustive_outcome_risk":True,
       "no_previous_scientific_cases_replayed":True}
def process(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T43_NATIVE_ARM_UNFROZEN")
    begin=time.process_time_ns()
    model=train()
    out=evaluate(arm,model)
    return {"arm":arm,"repeat":rep,"TRAIN_public_cases":TRAIN_COUNT,
      "sourcefree_validation_cases":VALID_COUNT,"scientific_DEV_generated":0,
      "metrics":{k:v for k,v in out.items() if k!="trace"},
      "case_fingerprints":[(row["id"],row["public"]) for row in out["trace"]],
      "full_process_CPU_ns":time.process_time_ns()-begin,
      "peak_RSS_bytes":peak_rss()}
def native():
    results={}
    for arm in ARMS:
        samples=[]
        for rep in range(5):
            child=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                     "--arm",arm,str(rep)],capture_output=True,text=True)
            if child.returncode:raise RuntimeError("T43_NATIVE_CHILD_FAILURE:"+arm+"\n"+child.stderr[-4000:])
            samples.append(json.loads(child.stdout))
        first=samples[0]
        if any(x["metrics"]!=first["metrics"] or x["case_fingerprints"]!=first["case_fingerprints"] for x in samples):
            raise RuntimeError("T43_NONDETERMINISTIC_SOURCEFREE_EVIDENCE")
        x=first["metrics"]
        results[arm]={"fraction_exact_optimal_utility":x["fraction_exact_optimal_utility"],
          "selected_expected_gain":x["total_selected_expected_gain"],
          "optimal_expected_gain":x["total_exact_optimal_gain"],
          "correct_certified":x["correct_certified"],
          "sound_abstentions":x["sound_abstentions"],
          "hard_false_certifications":x["hard_false_certifications"],
          "five_full_CPU_ns":[t["full_process_CPU_ns"] for t in samples],
          "five_peak_RSS_bytes":[t["peak_RSS_bytes"] for t in samples],
          "median_full_CPU_ns":int(statistics.median(t["full_process_CPU_ns"] for t in samples)),
          "max_peak_RSS_bytes":max(t["peak_RSS_bytes"] for t in samples)}
    learned=results[ARMS[0]]
    exact=(ARMS[1],ARMS[2])
    CPU=all(learned["median_full_CPU_ns"]<results[a]["median_full_CPU_ns"] for a in exact)
    RSS=learned["max_peak_RSS_bytes"]<=min(t["max_peak_RSS_bytes"] for a,t in results.items() if a!=ARMS[0])
    correct=all(results[a]["fraction_exact_optimal_utility"]>=1-1e-10 for a in exact)
    return {"schema":"wingless.T43.frozen-sourcefree-full-native-admission.v1",
       "TRAIN_public_cases":TRAIN_COUNT,"sourcefree_validation_cases":VALID_COUNT,
       "scientific_DEV_generated":0,
       "all_independent_exact_controls_truly_optimal":correct,
       "learned_99pct_optimal_utility":learned["fraction_exact_optimal_utility"]>=.99,
       "strict_CPU_below_both_exact_controls":CPU,
       "peak_RSS_no_more_than_best_control":RSS,
       "prospective_scientific_source_admitted":correct and CPU and RSS and learned["fraction_exact_optimal_utility"]>=.99,
       "arms":results,"scientific_positive":False,"optimizer_promotion":False}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--ablation":print(json.dumps(ablation(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--native":print(json.dumps(native(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":print(json.dumps(process(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T43_SOURCEFREE_ONLY_NO_SCIENTIFIC_DEV")
