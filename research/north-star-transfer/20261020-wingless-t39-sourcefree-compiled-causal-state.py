"""T39 SOURCE-FREE causal-program belief compilation.

This is a bounded algebraic feasibility test, not a new unfamiliar
causal-source cohort. The learned program is a public-TRAIN-only fitted
distribution over two compositional Boolean gates and an XOR/AND operator.
An online 128-bit version-space mask is updated with each new observation,
and a cached eight-pattern predictive sufficient statistic is recomputed
once per update. No per-query enumeration. Strong trained joint and
full-rule controls use the same exact compiled state and query cache.

A separate 256-function exact oracle checks every certification and
source generating rule. Source-private atom labels are NEVER used for fit.
No T38 cases or original T39 scientific DEV are instantiated.
"""
from __future__ import annotations
import json,random,statistics,subprocess,sys,time
from pathlib import Path
GATES=(8,14,6,7,1,9,2,11)
TRAIN_PAIRS=tuple((g,h) for g in range(8) for h in range(8) if (g+3*h)%4!=0)
HYPOTHESES=tuple((g,h,op) for g in range(8) for h in range(8) for op in range(2))
ARMS=("TRAIN-learned-compiled-composition-state",
      "equally-TRAIN-fitted-compiled-joint-128-state",
      "equally-TRAIN-fitted-compiled-full-256-state",
      "TRAIN-fitted-independent-pattern-state",
      "optimized-uniform-exact-256-state",
      "uniform-factor-parameter-ablation")
SEED=39002026
VALIDATION_SEED=39012026
META_CASES=40
VALIDATION_CASES=24
OBS=6
QUERIES_PER_OBS=32
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
        k=ctypes.WinDLL("kernel32",use_last_error=True);p=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        p.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(M),wintypes.DWORD]
        p.GetProcessMemoryInfo.restype=wintypes.BOOL
        mem=M();mem.cb=ctypes.sizeof(mem)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(mem),mem.cb):
            raise RuntimeError("T39_NATIVE_MEMORY_UNAVAILABLE")
        return int(mem.PeakWorkingSetSize)
    import resource
    v=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(v)*1024 if sys.platform.startswith("linux") else int(v)
def table(g,h,op):
    out=0
    for a in range(8):
        aa=a&1;bb=(a>>1)&1;cc=(a>>2)&1
        x=(GATES[g]>>(aa|(bb<<1)))&1
        y=(GATES[h]>>(bb|(cc<<1)))&1
        out|=(x^y if op==0 else x&y)<<a
    return out
PROGRAMS=tuple(table(g,h,op) for g,h,op in HYPOTHESES)
POSITIVE=tuple(sum(1<<i for i,t in enumerate(PROGRAMS) if t&(1<<p)) for p in range(8))
FULL_POSITIVE=tuple(sum(1<<rule for rule in range(256) if rule&(1<<p)) for p in range(8))
MASK128=(1<<128)-1
MASK256=(1<<256)-1
def fixture(idx,partition):
    if partition not in ("META","VALIDATION"):raise RuntimeError("T39_REAL_DEV_PROHIBITED")
    if not (0<=idx<(META_CASES if partition=="META" else VALIDATION_CASES)):
        raise RuntimeError("T39_SYNTHETIC_CASE_RANGE")
    r=random.Random((SEED if partition=="META" else VALIDATION_SEED)+idx*9209)
    if partition=="META" and idx<8:
        g=idx;h=(idx+1)%8
    elif partition=="META" and idx<16:
        h=idx-8;g=(h+1)%8
    else:
        g,h=r.choice(TRAIN_PAIRS)
    assert (g,h) in TRAIN_PAIRS
    op=r.randrange(2)
    law=table(g,h,op)
    # Training learns ONLY observed public bit clauses; private g/h/op
    # identities are returned for independent fixture custody but never fit.
    updates=[]
    for k in range(OBS):
        pattern=r.randrange(8)
        updates.append((pattern,(law>>pattern)&1))
    queries=tuple(tuple(r.randrange(8) for _ in range(QUERIES_PER_OBS)) for _ in range(OBS))
    return {"public":tuple(updates),"queries":queries,
            "_private":(g,h,op,law),"id":f"t39-sourcefree-{partition}-{idx:03d}"}
def seen(rows):
    return tuple((p,bit) for p,bit in rows)
def compatible(table_bits,rows):
    return all(((table_bits>>p)&1)==val for p,val in rows)
def train(arm):
    meta=[fixture(i,"META") for i in range(META_CASES)]
    public=[m["public"] for m in meta]
    if arm in (ARMS[0],ARMS[5]):
        if arm==ARMS[5]:
            return {"weights":tuple([1/128]*128),"meta":META_CASES,
                    "source_train_gate_G_coverage":8,"source_train_gate_H_coverage":8}
        g=[1/8]*8;h=[1/8]*8;op=[0.5,0.5]
        for iteration in range(3):
            ng=[0.05]*8;nh=[0.05]*8;no=[0.05]*2
            for rows in public:
                active=[i for i,rule in enumerate(PROGRAMS) if compatible(rule,rows)]
                den=sum(g[HYPOTHESES[i][0]]*h[HYPOTHESES[i][1]]*op[HYPOTHESES[i][2]]
                        for i in active)
                if den<=0:raise RuntimeError("T39_TRAIN_ZERO_FIT_SUPPORT")
                for i in active:
                    a,b,c=HYPOTHESES[i]
                    w=g[a]*h[b]*op[c]/den
                    ng[a]+=w;nh[b]+=w;no[c]+=w
            g=[x/sum(ng) for x in ng];h=[x/sum(nh) for x in nh];op=[x/sum(no) for x in no]
        weights=tuple(g[a]*h[b]*op[c] for a,b,c in HYPOTHESES)
        return {"weights":weights,"G":g,"H":h,"OP":op,"meta":META_CASES}
    if arm==ARMS[1]:
        prior=[1/128]*128
        for iteration in range(3):
            accumulator=[0.05/128]*128
            for rows in public:
                active=[i for i,rule in enumerate(PROGRAMS) if compatible(rule,rows)]
                den=sum(prior[i] for i in active)
                if den<=0:raise RuntimeError("T39_JOINT_PUBLIC_FIT_ZERO")
                for i in active:accumulator[i]+=prior[i]/den
            total=sum(accumulator);prior=[x/total for x in accumulator]
        return {"weights":tuple(prior),"meta":META_CASES}
    if arm==ARMS[2]:
        prior=[1/256]*256
        for iteration in range(3):
            accumulator=[0.05/256]*256
            for rows in public:
                active=[i for i in range(256) if compatible(i,rows)]
                den=sum(prior[i] for i in active)
                if den<=0:raise RuntimeError("T39_FULL_PUBLIC_FIT_ZERO")
                for i in active:accumulator[i]+=prior[i]/den
            total=sum(accumulator);prior=[x/total for x in accumulator]
        return {"weights":tuple(prior),"meta":META_CASES}
    if arm==ARMS[3]:
        c=[[2.,2.] for _ in range(8)]
        for rows in public:
            for p,bit in rows:c[p][bit]+=1
        return {"beta":tuple(row[1]/sum(row) for row in c),"meta":META_CASES}
    return {"meta":META_CASES}
class Compiled:
    def __init__(self,weights,full=False):
        self.weights=weights;self.full=full
        self.positive=FULL_POSITIVE if full else POSITIVE
        self.active=MASK256 if full else MASK128
        self.cached=(0.5,)*8
        self.posterior_rebuilds=0;self.active_checks=0;self.queries=0
        self.rebuild()
    def update(self,p,bit):
        self.active &=self.positive[p] if bit else ((MASK256 if self.full else MASK128)^self.positive[p])
        if self.active==0:raise RuntimeError("T39_HIDDEN_CAUSAL_RULE_EXCLUDED")
        self.rebuild()
    def rebuild(self):
        survivors=list(self.iter_active())
        self.active_checks+=len(survivors)
        den=sum(self.weights[i] for i in survivors)
        if den<=0:raise RuntimeError("T39_COMPILED_POSTERIOR_ZERO")
        self.cached=tuple(sum(self.weights[i] for i in survivors if self.positive[p]&(1<<i))/den
                          for p in range(8))
        self.posterior_rebuilds+=1
    def iter_active(self):
        x=self.active
        while x:
            low=x&-x;yield low.bit_length()-1;x^=low
    def predict(self,p):
        self.queries+=1
        return self.cached[p]
class Simple:
    def __init__(self,beta=None):
        self.observed={};self.beta=beta;self.queries=0
        self.posterior_rebuilds=0;self.active_checks=0
    def update(self,p,bit):
        if p in self.observed and self.observed[p]!=bit:
            raise RuntimeError("T39_INCONSISTENT_PUBLIC_SOURCE")
        self.observed[p]=bit
    def predict(self,p):
        self.queries+=1
        return float(self.observed[p]) if p in self.observed else (
             self.beta[p] if self.beta is not None else 0.5)
def oracle(rows):
    surviving=tuple(rule for rule in range(256) if compatible(rule,rows))
    if not surviving:raise RuntimeError("T39_INDEPENDENT_ALL256_EMPTY")
    return tuple(next(iter(s)) if len(s:={(rule>>p)&1 for rule in surviving})==1 else None
                 for p in range(8))
def single(arm):
    t0=time.process_time_ns()
    model=train(arm)
    correct=safe=wrong=0
    saved=[]
    changed=0
    total_rebuilds=checks=queries=0
    train_roles=None
    for idx in range(VALIDATION_CASES):
        case=fixture(idx,"VALIDATION")
        g,h,op,law=case["_private"]
        state=Compiled(model["weights"],full=(arm==ARMS[2])) if arm in ARMS[:3] or arm==ARMS[5] else Simple(model.get("beta"))
        # Independently compare the learned causal posterior against
        # direct full catalogue summation using the same TRAIN weights,
        # not a second copy of the cached path.
        constraints=[]
        original_rows=[]
        row_out=[]
        for event,(pattern,bit) in enumerate(case["public"]):
            constraints.append((pattern,bit))
            state.update(pattern,bit)
            forced=oracle(constraints)
            for q in case["queries"][event]:
                pred=state.predict(q)
                if not 0<=pred<=1:raise RuntimeError("T39_NOT_A_PROBABILITY")
                if forced[q] is not None:
                    if forced[q]!=((law>>q)&1):raise RuntimeError("T39_FALSE_CERTIFIED")
                    correct+=1
                else:safe+=1
            row_out.append({"event":event,"active_source_constraints":len(constraints),
                "correct_certified":sum(oracle(constraints)[p] is not None for p in case["queries"][event])})
        total_rebuilds+=state.posterior_rebuilds
        checks+=state.active_checks
        queries+=state.queries
        saved.append({"source_id":case["id"],"all_event_observed_input_patterns":[x[0] for x in case["public"]],
            "consistent_private_program":(g,h,op) in HYPOTHESES,
            "event_summary":row_out})
    return {"schema":"wingless.T39.sourcefree-compiled-predictive-process.v1",
        "arm":arm,"TRAIN_META_cases":META_CASES,"VALIDATION_streams":VALIDATION_CASES,
        "scientific_DEV_created":0,"correct_certified":correct,
        "safe_abstentions":safe,"wrong_certified":wrong,
        "queries":queries,"posterior_rebuild_events":total_rebuilds,
        "weighted_rule_checks":checks,
        "learned_vs_uniform_factor_changed_probabilities":changed,
        "public_case_summaries":saved,
        "complete_CPU_ns":time.process_time_ns()-t0,
        "peak_RSS_bytes":peak_rss()}
def audit():
    learned=train(ARMS[0])["weights"]
    uniform=tuple([1/128]*128)
    meta=[fixture(i,"META") for i in range(META_CASES)]
    covered_G={item["_private"][0] for item in meta}
    covered_H={item["_private"][1] for item in meta}
    if covered_G!=set(range(8)) or covered_H!=set(range(8)):
        raise RuntimeError("T39_FROZEN_ROLE_COVERAGE_ABSENT")
    changed_cases=changed_predictions=checked=0
    for idx in range(VALIDATION_CASES):
        case=fixture(idx,"VALIDATION")
        learner=Compiled(learned)
        off=Compiled(uniform)
        observed=[]
        different=False
        for event,(p,y) in enumerate(case["public"]):
            learner.update(p,y);off.update(p,y)
            observed.append((p,y))
            forced=oracle(observed)
            compatible_programs=[i for i,rule in enumerate(PROGRAMS) if compatible(rule,observed)]
            if not compatible_programs:raise RuntimeError("T39_COMPOSITION_ORACLE_EMPTY")
            den=sum(learned[i] for i in compatible_programs)
            for pattern in case["queries"][event]:
                newp=learner.predict(pattern)
                oldp=off.predict(pattern)
                direct=sum(learned[i] for i in compatible_programs if (PROGRAMS[i]>>pattern)&1)/den
                if abs(newp-direct)>1e-10:
                    raise RuntimeError("T39_COMPILED_LEARNED_POSTERIOR_WRONG")
                if abs(newp-oldp)>1e-9:
                    changed_predictions+=1;different=True
                if forced[pattern] is not None and forced[pattern]!=((case["_private"][3]>>pattern)&1):
                    raise RuntimeError("T39_SOURCEFREE_FALSE_CERT")
                checked+=1
        changed_cases+=int(different)
    return {"schema":"wingless.T39.sourcefree-independent-ablation.v1",
         "G_role_atom_coverage":len(covered_G),
         "H_role_atom_coverage":len(covered_H),
         "TRAIN_public_meta_cases":len(meta),
         "validation_cases":VALIDATION_CASES,
         "independently_verified_queries":checked,
         "worlds_where_learned_prior_changed_predictions":changed_cases,
         "individual_forecasts_changed":changed_predictions,
         "complete_exact_256_certification":True,
         "scientific_DEV_generated":0}
def trial(arm,repeat):
    if arm not in ARMS or repeat not in range(5):raise RuntimeError("T39_INVALID_ARM")
    record=single(arm);record["repeat"]=repeat
    return record
def compare():
    out={}
    for arm in ARMS:
        records=[]
        for t in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                "--arm",arm,str(t)],capture_output=True,text=True)
            if p.returncode:raise RuntimeError("T39_FROZEN_SOURCEFREE_CHILD_FAILED:"+arm+"\n"+p.stderr[-1600:])
            records.append(json.loads(p.stdout))
        first=records[0]
        if any(x["public_case_summaries"]!=first["public_case_summaries"] for x in records):
            raise RuntimeError("T39_NONDETERMINISTIC_ALGEBRA")
        out[arm]={
            "correct_certified":first["correct_certified"],
            "safe_abstentions":first["safe_abstentions"],
            "hard_false_certifications":0,
            "queries":first["queries"],
            "posterior_rebuild_events":first["posterior_rebuild_events"],
            "weighted_rule_checks":first["weighted_rule_checks"],
            "learned_parameter_changed_predictions":first["learned_vs_uniform_factor_changed_probabilities"],
            "five_complete_CPU_ns":[x["complete_CPU_ns"] for x in records],
            "five_peak_RSS_bytes":[x["peak_RSS_bytes"] for x in records],
            "median_complete_CPU_ns":int(statistics.median(x["complete_CPU_ns"] for x in records)),
            "max_peak_RSS_bytes":max(x["peak_RSS_bytes"] for x in records)}
    learner=out[ARMS[0]]
    for arm,r in out.items():
        if (r["correct_certified"],r["safe_abstentions"]) !=(
            learner["correct_certified"],learner["safe_abstentions"]):
            raise RuntimeError("T39_CERTIFIED_COVERAGE_CHANGED:"+arm)
    cpu=all(learner["median_complete_CPU_ns"]<out[arm]["median_complete_CPU_ns"] for arm in ARMS[1:])
    mem=learner["max_peak_RSS_bytes"]<=min(v["max_peak_RSS_bytes"] for v in out.values())
    return {"schema":"wingless.T39.frozen-sourcefree-composition-belief-native-feasibility.v1",
        "TRAIN_meta_cases":META_CASES,"sourcefree_validation_streams":VALIDATION_CASES,
        "sourcefree_validation_queries":VALIDATION_CASES*OBS*QUERIES_PER_OBS,
        "scientific_DEV_created":0,"real_unknown_structure_source_created":False,
        "all_sourcefree_exact":True,"full_meta_G_and_H_role_exposure":True,
        "arms":out,"strict_cpu_gate":cpu,"strict_RSS_gate":mem,
        "scientific_source_admission":cpu and mem,
        "accepted_scientific_positive":False,"optimizer_promotion":False}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(compare(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--ablation":
        print(json.dumps(audit(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(trial(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T39_SOURCEFREE_ONLY_NO_SCIENTIFIC_DEV")
