"""T34: source-free learned exact causal-state factorization feasibility.

This program deliberately creates NO Wingless experimental source cohorts.
Its row streams are synthetic public GF(2) algebra and cannot be described
as prospective evidence of a new learned cognitive mechanism.

Every stream is fully constraint-consistent. The learned object is a public
TRAIN-only coordinate partition. A factor node maintains an exact affine
GF2 subspace via an incremental basis; cross-partition observations merge
factors without dropping a row. Queries XOR exactly certified projections
from the independent factors, falling back to UNKNOWN when any one is not
entailed. Two optimized full-width controls share identical observations.

Ablating the learned grouping changes actual factor merges / working
component dimensions, but never the certified answer.
"""
from __future__ import annotations
import itertools,json,random,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
DIMS=(32,64)
TRAIN_PER_DIM=32
VALIDATE_PER_DIM=24
OBS=64
CHECKPOINT_EVERY=8
QUERIES_PER_CHECKPOINT=12
ARMS=("TRAIN-learned-incremental-exact-factorization",
      "untrained-adaptive-component-GF2",
      "TRAIN-informed-frequency-pivot-optimized-incremental-GF2",
      "optimized-exact-incremental-GF2",
      "optimized-batch-rebuild-GF2")
SEED=340119019
VALIDATION_SEED=340119813
def ram():
    if sys.platform=="win32":
        import ctypes
        from ctypes import wintypes
        class Mem(ctypes.Structure):
            _fields_=[("cb",wintypes.DWORD),("PageFaultCount",wintypes.DWORD),
             ("PeakWorkingSetSize",ctypes.c_size_t),("WorkingSetSize",ctypes.c_size_t),
             ("QuotaPeakPagedPoolUsage",ctypes.c_size_t),("QuotaPagedPoolUsage",ctypes.c_size_t),
             ("QuotaPeakNonPagedPoolUsage",ctypes.c_size_t),("QuotaNonPagedPoolUsage",ctypes.c_size_t),
             ("PagefileUsage",ctypes.c_size_t),("PeakPagefileUsage",ctypes.c_size_t)]
        kernel=ctypes.WinDLL("kernel32",use_last_error=True)
        psapi=ctypes.WinDLL("psapi",use_last_error=True)
        kernel.GetCurrentProcess.restype=wintypes.HANDLE
        psapi.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        psapi.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not psapi.GetProcessMemoryInfo(kernel.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T34_PROCESS_RSS_UNAVAILABLE:"+str(ctypes.get_last_error()))
        return int(m.PeakWorkingSetSize)
    import resource
    v=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(v)*1024 if sys.platform.startswith("linux") else int(v)

def bits(x):
    while x:
        low=x&-x
        yield low.bit_length()-1
        x^=low
def gram(dim):
    g=random.Random(SEED+dim*5519)
    perm=list(range(dim));g.shuffle(perm)
    return tuple(tuple(perm[i:i+8]) for i in range(0,dim,8))
def stream(dim,index,phase):
    if phase not in ("TRAIN","VALIDATION"):raise RuntimeError("T34_PHASE_UNKNOWN")
    r=random.Random((SEED if phase=="TRAIN" else VALIDATION_SEED)+dim*8791+index*9187)
    modules=gram(dim)
    truth=r.getrandbits(dim)
    public=[];questions=[]
    for timepoint in range(OBS):
        local=r.randrange(16)!=0
        if local:
            module=r.choice(modules)
            coords=r.sample(module,r.choice((2,3,4)))
        else:
            ma,mb=r.sample(range(len(modules)),2)
            coords=r.sample(modules[ma],2)+r.sample(modules[mb],r.choice((1,2)))
        mask=sum(1<<p for p in coords)
        public.append((mask,(mask&truth).bit_count()&1))
        if (timepoint+1)%CHECKPOINT_EVERY==0:
            query=[]
            for j in range(QUERIES_PER_CHECKPOINT):
                if j%4==0:
                    ma,mb=r.sample(range(len(modules)),2)
                    coords=[r.choice(modules[ma]),r.choice(modules[mb])]
                else:
                    module=r.choice(modules)
                    coords=r.sample(module,r.choice((1,2,3)))
                query.append(sum(1<<p for p in coords))
            questions.append(tuple(query))
    return tuple(public),tuple(questions),truth
def add(basis,mask,rhs,reverse=False):
    work=0
    while mask:
        pivot=(mask&-mask).bit_length()-1 if reverse else mask.bit_length()-1
        prev=basis.get(pivot)
        if prev is None:
            basis[pivot]=(mask,rhs)
            return work
        mask^=prev[0];rhs^=prev[1];work+=1
    if rhs:raise RuntimeError("T34_INCONSISTENT_PUBLIC_CAUSAL_EQUATIONS")
    return work
def decide(basis,query,reverse=False):
    rhs=0
    while query:
        pivot=(query&-query).bit_length()-1 if reverse else query.bit_length()-1
        row=basis.get(pivot)
        if row is None:return None
        query^=row[0];rhs^=row[1]
    return rhs
def independent_exact(rows,query):
    # Independent lower-pivot GF2 reduction, rebuilt from all public history.
    basis={}
    for m,r in rows:
        while m:
            p=(m&-m).bit_length()-1
            if p not in basis:
                basis[p]=(m,r);break
            x,y=basis[p];m^=x;r^=y
        else:
            if r:raise RuntimeError("T34_EXACT_ORACLE_INCONSISTENT")
    rhs=0
    while query:
        p=(query&-query).bit_length()-1
        if p not in basis:return None
        m,r=basis[p];query^=m;rhs^=r
    return rhs
def learned_partition(dim,training_rows):
    # Greedy TRAIN graph clustering of repeatedly co-observed causal
    # variables, max 8 bits/factor. Completely public; no hidden truth.
    affinity=Counter()
    for m in training_rows:
        selected=list(bits(m))
        for a,b in itertools.combinations(selected,2):
            affinity[min(a,b),max(a,b)]+=1
    member={i:{i} for i in range(dim)}
    parent=list(range(dim))
    def root(i):
        while parent[i]!=i:
            parent[i]=parent[parent[i]];i=parent[i]
        return i
    for (a,b),weight in sorted(affinity.items(),key=lambda item:(-item[1],item[0])):
        i,j=root(a),root(b)
        if i==j or len(member[i])+len(member[j])>8:continue
        if i>j:i,j=j,i
        parent[j]=i;member[i].update(member.pop(j))
    clusters=tuple(tuple(sorted(group)) for group in sorted(member.values(),key=lambda x:min(x)))
    if sorted(i for cluster in clusters for i in cluster)!=list(range(dim)):
        raise RuntimeError("T34_LEARNED_PARTITION_NOT_COMPLETE")
    return clusters,affinity
class Factors:
    def __init__(self,dim,groups):
        self.dim=dim
        self.nodes={}
        self.mapping=[-1]*dim
        self.nextid=0
        self.merges=0
        self.xor_work=0
        self.max_component_bits=0
        for group in groups:
            self._new(set(group),{},())
    def _new(self,bits_,basis,oldrows):
        ident=self.nextid;self.nextid+=1
        mask=sum(1<<x for x in bits_)
        self.nodes[ident]=[mask,dict(basis),list(oldrows)]
        for x in bits_:self.mapping[x]=ident
        self.max_component_bits=max(self.max_component_bits,mask.bit_count())
        return ident
    def add(self,mask,rhs):
        touched=sorted({self.mapping[bit] for bit in bits(mask)})
        if not touched:
            if rhs:raise RuntimeError("T34_FALSE_CONSTANT_OBSERVATION")
            return
        if len(touched)==1:
            ident=touched[0]
        else:
            # Previously disjoint local affine bases have disjoint support,
            # thus their pivotal rows can be merged EXACTLY without
            # rebuilding or changing previously observed constraints.
            mergedmask=0;mergedbasis={};mergedrows=[]
            for ident in touched:
                node=self.nodes.pop(ident)
                mergedmask|=node[0]
                if set(mergedbasis).intersection(node[1]):
                    raise RuntimeError("T34_DISJOINT_PIVOT_IDENTITY_COLLISION")
                mergedbasis.update(node[1]);mergedrows.extend(node[2])
            ident=self._new(set(bits(mergedmask)),mergedbasis,mergedrows)
            self.merges+=len(touched)-1
        node=self.nodes[ident]
        self.xor_work+=add(node[1],mask,rhs)
        node[2].append((mask,rhs))
    def decide(self,mask):
        value=0
        for node_mask,basis,_ in self.nodes.values():
            masked=mask&node_mask
            if not masked:continue
            response=decide(basis,masked)
            if response is None:return None
            value^=response
        return value
    def signature(self):
        return (len(self.nodes),self.merges,self.max_component_bits,self.xor_work)

def fit(dim):
    training=[]
    for index in range(TRAIN_PER_DIM):
        rows,_,_=stream(dim,index,"TRAIN")
        training.extend(mask for mask,_ in rows)
    groups,cooccurrence=learned_partition(dim,training)
    frequencies=Counter()
    for m in training:
        frequencies.update(bits(m))
    pivotorder=tuple(sorted(range(dim),key=lambda x:(frequencies[x],x)))
    return groups,pivotorder,len(cooccurrence)

def per_stream(dim,idx,groups,arm,pivotorder,validate_proofs=True):
    public,checkpoints,truth=stream(dim,idx,"VALIDATION")
    if arm==ARMS[0]:state=Factors(dim,groups)
    elif arm==ARMS[1]:state=Factors(dim,tuple((i,) for i in range(dim)))
    else:state={}
    rowhistory=[]
    outputs=[]
    ops=0;checkpoint=0
    for i,(mask,rhs) in enumerate(public):
        rowhistory.append((mask,rhs))
        if arm in ARMS[:2]:
            state.add(mask,rhs)
        elif arm==ARMS[2]:
            # Full-width, TRAIN-informed pivot strategy retained as exact
            # explicit basis order. Count no speculative learning gain.
            ops+=add(state,mask,rhs,reverse=True)
        elif arm==ARMS[3]:
            ops+=add(state,mask,rhs)
        elif arm==ARMS[4]:
            pass
        if (i+1)%CHECKPOINT_EVERY:
            continue
        for query in checkpoints[checkpoint]:
            if arm in ARMS[:2]:out=state.decide(query)
            elif arm==ARMS[4]:
                basis={}
                for m,r in rowhistory:add(basis,m,r)
                out=decide(basis,query)
            else:out=decide(state,query,reverse=(arm==ARMS[2]))
            if validate_proofs:
                expect=independent_exact(rowhistory,query)
                if out!=expect:
                    raise RuntimeError("T34_EXACT_GF2_INFERENCE_MISMATCH:"+arm+":"+str(dim))
                if out is not None and out!=((query&truth).bit_count()&1):
                    raise RuntimeError("T34_FALSE_CERTIFICATION")
            outputs.append(-1 if out is None else out)
        checkpoint+=1
    if len(outputs)!=CHECKPOINT_EVERY and not outputs:
        raise RuntimeError("T34_MISSING_ALGEBRA_STREAM")
    proof={"dimension":dim,"stream_id":idx,
           "certified":sum(x!=-1 for x in outputs),
           "sound_abstentions":sum(x==-1 for x in outputs),
           "all_outputs":outputs}
    if arm in ARMS[:2]:
        nodes,merged,maximum,work=state.signature()
        proof.update({"factor_count":nodes,"factor_merges":merged,
                      "max_component_width":maximum,"basis_xor_updates":work})
    else:proof.update({"factor_count":1,"factor_merges":0,
                       "max_component_width":dim,"basis_xor_updates":ops})
    return proof
def child(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T34_CHILD_ARGUMENT_INVALID")
    tick=time.process_time_ns()
    complete=[]
    differences=0
    fitted={}
    for dim in DIMS:
        groups,pivotorder,edges=fit(dim)
        fitted[dim]=groups
        for index in range(VALIDATE_PER_DIM):
            result=per_stream(dim,index,groups,arm,pivotorder)
            complete.append(result)
            if arm==ARMS[0] and result["factor_count"]!=dim:differences+=1
    return {"schema":"wingless.T34.source-free-factorization-process.v1",
      "arm":arm,"repeat":rep,"TRAIN_only_public_streams":TRAIN_PER_DIM*len(DIMS),
      "source_free_validation_streams":VALIDATE_PER_DIM*len(DIMS),
      "actual_scientific_DEV_cases":0,"source_free_query_results":complete,
      "TRAIN_fitted_groups":{str(k):v for k,v in fitted.items()} if arm==ARMS[0] else {},
      "factorization_affected_validation_streams":differences,
      "complete_native_CPU_ns":time.process_time_ns()-tick,
      "peak_process_RSS_bytes":ram(),
      "no_scientific_positive":True,"optimizer_admission":False}
def run():
    results={}
    for arm in ARMS:
        reps=[]
        for rep in range(5):
            proc=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                     "--child",arm,str(rep)],capture_output=True,text=True)
            if proc.returncode:
                raise RuntimeError("T34_SOURCEFREE_CHILD_FAILED:"+arm+":"+str(rep)+"\n"+proc.stderr)
            reps.append(json.loads(proc.stdout))
        if any(r["source_free_query_results"]!=reps[0]["source_free_query_results"] for r in reps):
            raise RuntimeError("T34_ALGEBRA_VALIDATION_NOT_DETERMINISTIC")
        first=reps[0]
        flattened=[r for r in first["source_free_query_results"]]
        result={
          "source_free_streams":len(flattened),
          "source_free_checkpoint_questions":sum(r["certified"]+r["sound_abstentions"] for r in flattened),
          "correct_certified":sum(r["certified"] for r in flattened),
          "sound_abstentions":sum(r["sound_abstentions"] for r in flattened),
          "factorization_changed_streams":first["factorization_affected_validation_streams"],
          "factor_merges":sum(r["factor_merges"] for r in flattened),
          "total_factor_update_xor_ops":sum(r["basis_xor_updates"] for r in flattened),
          "mean_max_factor_width":statistics.mean(r["max_component_width"] for r in flattened),
          "trained_factor_layouts":first["TRAIN_fitted_groups"],
          "five_complete_native_CPU_ns":[x["complete_native_CPU_ns"] for x in reps],
          "five_peak_process_RSS_bytes":[x["peak_process_RSS_bytes"] for x in reps],
          "median_complete_native_CPU_ns":int(statistics.median(x["complete_native_CPU_ns"] for x in reps)),
          "max_peak_process_RSS_bytes":max(x["peak_process_RSS_bytes"] for x in reps)
        }
        results[arm]=result
    exact=results[ARMS[3]]
    for name,r in results.items():
        if (r["correct_certified"],r["sound_abstentions"])!=(exact["correct_certified"],exact["sound_abstentions"]):
            raise RuntimeError("T34_SOURCEFREE_NOT_EXACT_EQUAL_TO_OPTIMIZED_BASELINE")
    learned=results[ARMS[0]]
    native_admission=all(
        learned["median_complete_native_CPU_ns"]<results[a]["median_complete_native_CPU_ns"]
        for a in ARMS[1:]
    ) and learned["max_peak_process_RSS_bytes"]<=min(
        item["max_peak_process_RSS_bytes"] for item in results.values())
    return {
      "schema":"wingless.T34.source-free-mathematical-and-native-feasibility.v1",
      "scientific_DEV_cohort_created":False,"all_old_DEV_untouched":True,
      "TRAIN_public_streams":TRAIN_PER_DIM*len(DIMS),
      "SOURCE_FREE_validation_streams":VALIDATE_PER_DIM*len(DIMS),
      "full_exact_per_query_independent_rebuild_verified":True,
      "dimensions":DIMS,"arms":results,
      "strict_native_admission":native_admission,
      "scientific_positive":False,"optimizer_admission":False,
      "conclusion":"PRE_SOURCE_FEASIBLE" if native_admission else "PRE_SOURCE_NATIVE_NO_GO"
    }
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(run(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(child(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T34_SOURCEFREE_ONLY")
