"""T32: TRAIN-learned parity-resolution rule construction, not sorting.

Public parity observations imply linear GF(2) equations over hidden causal
truth-table bits. A TRAIN-induced collection of row-pair inference macros
adds *new exact singleton rule constraints* to the posterior before
optimized Gaussian elimination. These rules are not conjectures:
each is the XOR of two witnessed public rows. The macro ablation changes
the CONSTRUCTED constraint state (and candidate rule filtering), while
both candidates must match independent brute simulation of ALL laws.
"""
from __future__ import annotations
import importlib.util,json,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
ROOT=Path(__file__).resolve().parent
ARMS=("TRAIN-learned-parity-constraint-propagation",
      "TRAIN-frequency-compiled-parity-closure",
      "optimized-exact-GF2-solver",
      "source-blind-exhaustive-rule-enumerator")
def source():
    spec=importlib.util.spec_from_file_location("t32_source",ROOT/"20261019-wingless-t32-parity-causal-sources.py")
    obj=importlib.util.module_from_spec(spec);sys.modules[spec.name]=obj
    spec.loader.exec_module(obj);return obj
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
        mem=Mem();mem.cb=ctypes.sizeof(mem)
        if not p.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(mem),mem.cb):
            raise RuntimeError("T32_NATIVE_MEMORY_ERROR="+str(ctypes.get_last_error()))
        return int(mem.PeakWorkingSetSize)
    import resource
    val=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(val)*1024 if sys.platform.startswith("linux") else int(val)
def build_rows(public,src,graph_id,mode):
    f=public["family"];n,k=src.FAMILIES[f];graph=src.TEMPLATES[f][graph_id]
    rows=[]
    for probe in public["probes"]:
        pre=probe["pre"];action=probe["do"]
        for nodes,value in probe["observed_parities"]:
            if len(nodes) not in (2,3) or len(set(nodes))!=len(nodes):
                raise RuntimeError("T32_MALFORMED_PUBLIC_PARITY")
            mask=0;known=0
            for j in nodes:
                if action is not None and action[0]==j:
                    known^=action[1]
                else:
                    pat=src.pattern(pre,graph[j])
                    mask^=1<<pat
                    known^=(pre>>j)&1
                    if (j%2):known^=mode*((pre>>((j+1)%n))&1)
            rows.append((mask,value^known))
    return tuple(rows)
def xor_masks(rows):
    """Find exact derived degree-one rows and their witnessed parents."""
    found=[]
    seen=set()
    for i,a in enumerate(rows):
        for b in rows[i+1:]:
            if a[0].bit_count()<=1 or b[0].bit_count()<=1:continue
            outcome=(a[0]^b[0],a[1]^b[1])
            if outcome[0].bit_count()==1:
                parents=tuple(sorted((a[0],b[0])))
                if parents not in seen:found.append(parents);seen.add(parents)
    return found
def fit(src,publics):
    counts=Counter();fcounts={f:Counter() for f in src.FAMILIES}
    for public in publics:
        family=public["family"]
        for graph in range(8):
            for mode in (0,1):
                matches=xor_masks(build_rows(public,src,graph,mode))
                counts.update((family,pair) for pair in matches)
                fcounts[family].update(matches)
    learned={}
    frequency={}
    for family in src.FAMILIES:
        ranked=sorted(fcounts[family],key=lambda pair:(-fcounts[family][pair],pair))
        learned[family]=tuple(ranked[:32])
        # Trained count-only baseline works with the same public TRAIN,
        # uses its own most frequent raw row signatures and closure.
        frequency[family]=tuple(sorted(ranked,key=lambda pair:(-sum(fcounts[family][p] for p in ranked if pair[0] in p),pair))[:32])
    return learned,frequency,{"TRAIN_public_macro_candidates":len(counts),
                              "TRAIN_learned_macro_templates":sum(map(len,learned.values()))}
def apply_macros(rows,macros):
    """Derive new UNIT constraints from TRAIN-selected exact parity rules."""
    grouped={}
    for rowmask,bit in rows:
        if rowmask in grouped and grouped[rowmask]!=bit:
            return tuple(rows),0,0
        grouped[rowmask]=bit
    added=[];proofs=0
    original=set(rows)
    for a,b in macros:
        if a not in grouped or b not in grouped:continue
        mask=a^b;val=grouped[a]^grouped[b]
        if mask.bit_count()!=1:raise RuntimeError("T32_TRAIN_MACRO_NOT_UNIT")
        proofs+=1
        derived=(mask,val)
        if derived not in original and derived not in added:
            added.append(derived)
    return tuple(rows)+tuple(added),len(added),proofs
def gf2_laws(rows,width):
    """Exact rule-space from canonical full-rank Gaussian GF(2) constraints."""
    basis={}
    for original_mask,original_rhs in rows:
        mask=original_mask;rhs=original_rhs
        while mask:
            pivot=mask.bit_length()-1
            if pivot not in basis:
                basis[pivot]=(mask,rhs);break
            old,oldrhs=basis[pivot]
            mask^=old;rhs^=oldrhs
        else:
            if rhs:return ()
    pivots=set(basis)
    free=[i for i in range(width) if i not in pivots]
    laws=[]
    for assignments in range(1<<len(free)):
        law=sum(((assignments>>j)&1)<<bit for j,bit in enumerate(free))
        for pivot in sorted(pivots):
            mask,bit=basis[pivot]
            implied=bit^((law& (mask^(1<<pivot))).bit_count()&1)
            law|=implied<<pivot
        laws.append(law)
    return tuple(laws)
def all_worlds(public,src,macros=None):
    f=public["family"];width=1<<src.FAMILIES[f][1]
    result=[];derived=proofs=0
    for gid in range(8):
        for mode in (0,1):
            original=build_rows(public,src,gid,mode)
            new,extra,found=apply_macros(original,macros[f]) if macros is not None else (original,0,0)
            # Independent proof: every generated clause must hold in every
            # assignment satisfying original source equations.
            laws=gf2_laws(new,width)
            if extra:
                baseline=gf2_laws(original,width)
                if set(laws)!=set(baseline):
                    raise RuntimeError("T32_LEARNED_CONSTRAINT_DROPPED_POSSIBLE_LAW")
            derived+=extra;proofs+=found
            result.extend((gid,mode,law) for law in laws)
    if not result:raise RuntimeError("T32_NO_HIDDEN_CAUSAL_WORLDS")
    return tuple(result),derived,proofs
def brute_worlds(public,src):
    f=public["family"];width=1<<src.FAMILIES[f][1];worlds=[]
    for gid,graph in enumerate(src.TEMPLATES[f]):
        for mode in (0,1):
            for law in range(1<<width):
                correct=True
                for probe in public["probes"]:
                    after=src.update(f,probe["pre"],graph,mode,law,probe["do"])
                    for nodes,expected in probe["observed_parities"]:
                        observed=sum((after>>j)&1 for j in nodes)%2
                        if observed!=expected:correct=False;break
                    if not correct:break
                if correct:worlds.append((gid,mode,law))
    if not worlds:raise RuntimeError("T32_BRUTE_NO_CAUSAL_WORLDS")
    return tuple(worlds)
def future(pub,src,world):
    f=pub["family"];gid,mode,law=world;start=pub["query"]["start"]
    for step in range(3):
        start=src.update(f,start,src.TEMPLATES[f][gid],mode,law,
                         pub["query"]["do"] if step==0 else None)
    return tuple((start>>node)&1 for node in pub["query"]["targets"])
def certify(public,src,worlds):
    possible=[set(),set(),set()]
    for world in worlds:
        for idx,answer in enumerate(future(public,src,world)):
            possible[idx].add(answer)
    result=tuple(next(iter(s)) if len(s)==1 else None for s in possible)
    return result,tuple(tuple(sorted(s)) for s in possible)
def train_only():
    src=source();TRAIN=[src.train_public(src.sample(name)) for name in src.TRAIN]
    learned,freq,stats=fit(src,TRAIN)
    retained=equivalent=certified=unknown=affected=unit_rows=proofs=0
    for name,public in zip(src.TRAIN,TRAIN):
        case=src.sample(name)
        raw,_,_=all_worlds(public,src,None)
        enhanced,added,witnessed=all_worlds(public,src,learned)
        control,_,_=all_worlds(public,src,freq)
        independent=brute_worlds(public,src)
        if set(raw)!=set(enhanced) or set(raw)!=set(independent) or set(raw)!=set(control):
            raise RuntimeError("T32_TRAIN_POSTERIOR_NOT_FULLY_SOUND")
        truth=(case["_true_graph"],case["_true_mode"],case["_true_law"])
        if truth not in raw:raise RuntimeError("T32_TRUE_CAUSAL_RULE_EXCLUDED")
        retained+=1
        r,poss=certify(public,src,enhanced)
        exact,_=certify(public,src,raw)
        if r!=exact:raise RuntimeError("T32_LEARNED_CERTIFICATE_MISMATCH")
        equivalent+=3
        if added>0:affected+=1
        unit_rows+=added;proofs+=witnessed
        for val,options,answer in zip(r,poss,case["_answers"]):
            if answer not in options:raise RuntimeError("T32_TRAIN_TRUE_ANSWER_NOT_IN_POSTERIOR")
            if val is None:unknown+=1
            else:
                if val!=answer:raise RuntimeError("T32_TRAIN_HARD_FALSE_POSITIVE")
                certified+=1
    return {"schema":"wingless.T32.TRAIN-only-real-clause-construction.v1",
      "status":"PASS_TRAIN_NO_DEV", "TRAIN":40, "DEV_opened":0,
      "TRUE_hidden_rules_retained":retained,"exact_counterfactual_equivalences":equivalent,
      "certified":certified,"safe_abstentions":unknown,
      "uncertainty_fraction":unknown/equivalent,
      "learned_parameters_change_actual_constructed_constraint_cases":affected,
      "newly_constructed_unit_constraints_from_TRAIN_macros":unit_rows,
      "witnessed_public_parity_inference_proofs":proofs,
      "TRAIN_public_macros_learned":stats["TRAIN_learned_macro_templates"],
      "full_brute_exact_posterior_equivalences":40,
      "no_previous_consumed_scientific_cases":True,
      "optimizer_admission":False,"scientific_positive":False}
def one(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T32_BAD_ARM")
    begin=time.process_time_ns()
    src=source()
    training=[src.train_public(src.sample(name)) for name in src.TRAIN]
    trained,frequency,_=fit(src,training) if arm in ARMS[:2] else ({},{},{})
    active=trained if arm==ARMS[0] else frequency if arm==ARMS[1] else None
    rows=[];correct=unknown=derived=0
    for name in src.DEV:
        case=src.sample(name);public=src.dev_public(case)
        if arm==ARMS[3]:worlds=brute_worlds(public,src);extra=0
        else:worlds,extra,_=all_worlds(public,src,active)
        output,possible=certify(public,src,worlds)
        numcorrect=numunknown=0
        for pred,opts,truth in zip(output,possible,case["_answers"]):
            if truth not in opts:raise RuntimeError("T32_DEV_TRUE_ANSWER_EXCLUDED")
            if pred is None:unknown+=1;numunknown+=1
            else:
                if pred!=truth:raise RuntimeError("T32_ORIGINAL_HARD_FALSE_POSITIVE")
                correct+=1;numcorrect+=1
        derived+=extra
        rows.append({"source_id":name,"correct_certified":numcorrect,
          "safe_abstentions":numunknown,"hypotheses":len(worlds),
          "added_learned_unit_clauses":extra})
    if len(rows)!=64 or correct+unknown!=192:raise RuntimeError("T32_ORIGINAL_INCOMPLETE")
    return {"schema":"wingless.T32.original-constraint-learning-process.v1",
       "arm":arm,"repeat":rep,"per_case":rows,"correct_certified":correct,
       "safe_abstentions":unknown,"hard_false_positives":0,
       "learned_extra_unit_clauses":derived,
       "full_process_CPU_ns":time.process_time_ns()-begin,
       "peak_RSS_bytes":rss(),"DEV_cases":64,
       "no_older_DEV_reuse":True,"optimizer_admission":False}
def evaluate():
    outcome={}
    for arm in ARMS:
        batch=[]
        for i in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                   "--arm",arm,str(i)],capture_output=True,text=True)
            if p.returncode:raise RuntimeError("T32_ORIGINAL_CHILD_FAIL:"+arm+":"+str(i)+"\n"+p.stderr)
            batch.append(json.loads(p.stdout))
        base=batch[0]
        if any(item["per_case"]!=base["per_case"] for item in batch):
            raise RuntimeError("T32_NONDETERMINISTIC_CASES")
        outcome[arm]={"correct_certified":base["correct_certified"],
          "safe_abstentions":base["safe_abstentions"],"hard_false_positives":0,
          "learned_extra_unit_clauses":base["learned_extra_unit_clauses"],
          "per_case":base["per_case"],
          "five_complete_CPU_ns":[a["full_process_CPU_ns"] for a in batch],
          "five_peak_RSS_bytes":[a["peak_RSS_bytes"] for a in batch],
          "median_CPU_ns":int(statistics.median(a["full_process_CPU_ns"] for a in batch)),
          "max_peak_RSS_bytes":max(a["peak_RSS_bytes"] for a in batch)}
    if any(x["correct_certified"]+x["safe_abstentions"]!=192 or x["hard_false_positives"] for x in outcome.values()):
        raise RuntimeError("T32_FIRSTUSE_SAFETY_FAILED")
    learned=outcome[ARMS[0]]
    return {"schema":"wingless.T32.original-unknown-causal-parity-comparison.v1",
      "TRAIN_cases":40,"prospective_DEV_cases":64,"counterfactuals":192,
      "arms":outcome,
      "comparison":{name:{"correct_baseline_preserved":learned["correct_certified"]>=x["correct_certified"],
            "complete_CPU_no_regression":learned["median_CPU_ns"]<=x["median_CPU_ns"],
            "peak_RSS_no_regression":learned["max_peak_RSS_bytes"]<=x["max_peak_RSS_bytes"]}
          for name,x in outcome.items() if name!=ARMS[0]},
      "external_replication":False,"scientific_positive":False,
      "optimizer_admission":False,"old_DEV_replayed":False}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(evaluate(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--train-only":
        print(json.dumps(train_only(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--arm":
        print(json.dumps(one(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T32_INVALID_ARGS")
