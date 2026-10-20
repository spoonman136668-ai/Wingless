"""T35: source-free TRAIN-learned exact multi-pivot GF2 elimination cache.

A public TRAIN fitted coordinate order groups correlated pivot coordinates
into four-pivot lookup blocks. A cache hit reduces several active pivot rows
with ONE exact equation XOR. Cache construction and invalidation XOR work
are measured explicitly. Exact scalar and fixed-block alternatives share
identical public TRAIN and source-free validation rows.
No actual Wingless scientific cohort, held-out world or DEV is generated.
"""
from __future__ import annotations
import itertools,json,random,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
DIMS=(32,64)
TRAIN_PER_DIM=24
VALIDATE_PER_DIM=20
OBS=64
CHECKPOINT_EVERY=8
QUERIES_PER_CHECKPOINT=12
ARMS=("TRAIN-learned-block-exact-GF2",
      "TRAIN-frequency-optimized-scalar-GF2",
      "deterministic-block-exact-GF2",
      "optimized-exact-incremental-GF2",
      "optimized-batch-rebuild-GF2")
SEED=35023019
VALIDATION_SEED=35088031
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
def trained_pivot(mask,rank):
    return max(bits(mask),key=lambda bit:rank[bit])
def add_trained(basis,mask,rhs,rank):
    work=0
    while mask:
        p=trained_pivot(mask,rank)
        row=basis.get(p)
        if row is None:
            basis[p]=(mask,rhs)
            return work
        mask^=row[0];rhs^=row[1];work+=1
    if rhs:raise RuntimeError("T34_INCONSISTENT_TRAINED_GF2")
    return work
def decide_trained(basis,mask,rank):
    rhs=0
    while mask:
        p=trained_pivot(mask,rank)
        row=basis.get(p)
        if row is None:return None
        mask^=row[0];rhs^=row[1]
    return rhs
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
def fit(dim):
    frequency=Counter();joint=Counter()
    for index in range(TRAIN_PER_DIM):
        rows,_,_=stream(dim,index,"TRAIN")
        for mask,_ in rows:
            nodes=list(bits(mask))
            frequency.update(nodes)
            joint.update((min(a,b),max(a,b)) for a,b in itertools.combinations(nodes,2))
    groups={i:{i} for i in range(dim)}
    parent=list(range(dim))
    def root(i):
        while parent[i]!=i:
            parent[i]=parent[parent[i]]
            i=parent[i]
        return i
    for (a,b),n in sorted(joint.items(),key=lambda z:(-z[1],z[0])):
        ia,ib=root(a),root(b)
        if ia==ib or len(groups[ia])+len(groups[ib])>4:continue
        if ia>ib:ia,ib=ib,ia
        groups[ia].update(groups.pop(ib));parent[ib]=ia
    learned=[]
    for group in sorted(groups.values(),key=lambda g:(-sum(frequency[i] for i in g),min(g))):
        learned.extend(sorted(group,key=lambda v:(-frequency[v],v)))
    freq=tuple(sorted(range(dim),key=lambda i:(-frequency[i],i)))
    if set(learned)!=set(range(dim)) or len(learned)!=dim:
        raise RuntimeError("T35_TRAINED_GF2_COORDINATE_NONBIJECTION")
    return tuple(learned),freq

def remap(mask,order):
    return sum(1<<i for i,v in enumerate(order) if mask&(1<<v))

class CachedBasis:
    """Incremental exact triangular basis with lazy four-pivot XOR words.

    Cache materialization XOR work is counted separately from row reduction.
    Cache entries are invalidated ONLY when a pivot enters their block.
    """
    def __init__(self,dim,order):
        self.dim=dim;self.order=order;self.basis={}
        self.pivots=[0]*((dim+3)//4)
        self.cache=[{} for _ in self.pivots]
        self.reduction_xors=0;self.cache_xors=0
        self.hits=0;self.avoided=0
    def reduce(self,word):
        for block in range(len(self.pivots)-1,-1,-1):
            available=self.pivots[block]
            if not available:continue
            key=word&available
            if not key:continue
            if key in self.cache[block]:
                combination,n=self.cache[block][key]
                word^=combination;self.reduction_xors+=1
                self.hits+=1;self.avoided+=max(0,n-1)
                continue
            combination=0;count=0
            for pivot in sorted(bits(available),reverse=True):
                if (word>>pivot)&1:
                    row=self.basis[pivot]
                    word^=row;combination^=row
                    self.reduction_xors+=1
                    self.cache_xors+=1
                    count+=1
            if count:self.cache[block][key]=(combination,count)
        return word
    def add(self,mask,rhs):
        word=remap(mask,self.order)|(rhs<<self.dim)
        word=self.reduce(word)
        remaining=word&((1<<self.dim)-1)
        if not remaining:
            if word>>self.dim:raise RuntimeError("T35_CONTRADICTORY_PUBLIC_EQUATION")
            return
        pivot=remaining.bit_length()-1
        if pivot in self.basis:raise RuntimeError("T35_NONTRIANGULAR_ROW")
        self.basis[pivot]=word
        i=pivot//4
        self.pivots[i]|=1<<pivot
        self.cache[i].clear()
    def query(self,mask):
        word=self.reduce(remap(mask,self.order))
        return (word>>self.dim)&1 if not(word&((1<<self.dim)-1)) else None

class ScalarBasis:
    def __init__(self,dim,order):
        self.dim=dim;self.order=order;self.basis={};self.xors=0
    def add(self,mask,rhs):
        word=remap(mask,self.order)|(rhs<<self.dim)
        while word&((1<<self.dim)-1):
            pivot=(word&((1<<self.dim)-1)).bit_length()-1
            if pivot not in self.basis:
                self.basis[pivot]=word;return
            word^=self.basis[pivot];self.xors+=1
        if word>>self.dim:raise RuntimeError("T35_SCALAR_EQUATION_CONTRADICTION")
    def query(self,mask):
        word=remap(mask,self.order)
        while word&((1<<self.dim)-1):
            pivot=(word&((1<<self.dim)-1)).bit_length()-1
            if pivot not in self.basis:return None
            word^=self.basis[pivot];self.xors+=1
        return (word>>self.dim)&1

def test_stream(dim,index,arm,learned,frequency):
    public,queries,truth=stream(dim,index,"VALIDATION")
    if arm==ARMS[0]:solver=CachedBasis(dim,learned)
    elif arm==ARMS[2]:solver=CachedBasis(dim,tuple(range(dim)))
    elif arm==ARMS[1]:solver=ScalarBasis(dim,frequency)
    else:solver=ScalarBasis(dim,tuple(range(dim)))
    observed=[];answers=[];j=0;batch_xors=0
    for t,(mask,rhs) in enumerate(public):
        observed.append((mask,rhs))
        if arm!=ARMS[4]:solver.add(mask,rhs)
        if (t+1)%CHECKPOINT_EVERY:continue
        for query in queries[j]:
            if arm==ARMS[4]:
                batch=ScalarBasis(dim,tuple(range(dim)))
                for m,v in observed:batch.add(m,v)
                outcome=batch.query(query);batch_xors+=batch.xors
            else:outcome=solver.query(query)
            independent=independent_exact(observed,query)
            if outcome!=independent:raise RuntimeError("T35_NOT_EXACT:"+arm+":"+str(dim))
            if outcome is not None and outcome!=((query&truth).bit_count()&1):
                raise RuntimeError("T35_HARD_FALSE_CERTIFICATION")
            answers.append(-1 if outcome is None else outcome)
        j+=1
    if len(answers)!=(OBS//CHECKPOINT_EVERY)*QUERIES_PER_CHECKPOINT:
        raise RuntimeError("T35_MISSING_CHECKPOINT_QUERIES")
    reduction=solver.reduction_xors if arm in (ARMS[0],ARMS[2]) else solver.xors
    cache=solver.cache_xors if arm in (ARMS[0],ARMS[2]) else 0
    return {"dim":dim,"index":index,"answers":answers,
       "row_reduction_xors":reduction if arm!=ARMS[4] else batch_xors,
       "cache_materialization_xors":cache,
       "all_actual_xors":reduction+cache if arm!=ARMS[4] else batch_xors,
       "cached_block_hits":solver.hits if arm in (ARMS[0],ARMS[2]) else 0,
       "saved_row_reduction_xors":solver.avoided if arm in (ARMS[0],ARMS[2]) else 0}

def ablation():
    different=learned_xors=fixed_xors=learned_row=fixed_row=saved=0
    for dim in DIMS:
        order,freq=fit(dim)
        for i in range(VALIDATE_PER_DIM):
            on=test_stream(dim,i,ARMS[0],order,freq)
            off=test_stream(dim,i,ARMS[2],order,freq)
            if on["answers"]!=off["answers"]:raise RuntimeError("T35_ABLATION_CHANGED_EXACT_ANSWER")
            if on["all_actual_xors"]!=off["all_actual_xors"]:different+=1
            learned_xors+=on["all_actual_xors"];fixed_xors+=off["all_actual_xors"]
            learned_row+=on["row_reduction_xors"];fixed_row+=off["row_reduction_xors"]
            saved+=on["saved_row_reduction_xors"]
    return {"schema":"wingless.T35.sourcefree-trained-block-XOR-ablation.v1",
       "TRAIN_public_streams":TRAIN_PER_DIM*len(DIMS),
       "sourcefree_validation_streams":VALIDATE_PER_DIM*len(DIMS),
       "checked_exact_queries":VALIDATE_PER_DIM*len(DIMS)*(OBS//CHECKPOINT_EVERY)*QUERIES_PER_CHECKPOINT,
       "scientific_DEV_generated":0,"streams_learned_parameters_change_real_XOR_work":different,
       "learned_row_reduction_XORs":learned_row,"fixed_block_row_reduction_XORs":fixed_row,
       "learned_all_XORs_INCLUDING_cache_construction":learned_xors,
       "fixed_block_all_XORs_INCLUDING_cache_construction":fixed_xors,
       "learned_minus_fixed_total_XORs":learned_xors-fixed_xors,
       "learner_cache_hit_reduction_savings":saved,"answers_identical":True}

def child(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T35_CHILD_ARGUMENT_INVALID")
    start=time.process_time_ns()
    cases=[]
    for dim in DIMS:
        learned,freq=fit(dim)
        for i in range(VALIDATE_PER_DIM):
            cases.append(test_stream(dim,i,arm,learned,freq))
    return {"schema":"wingless.T35.complete-sourcefree-process.v1","arm":arm,"repeat":rep,
        "cases":cases,"complete_native_CPU_ns":time.process_time_ns()-start,
        "peak_RSS_bytes":ram(),"scientific_DEV_generated":0}

def run():
    audit=ablation();results={}
    for arm in ARMS:
        reps=[]
        for i in range(5):
            p=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                    "--child",arm,str(i)],capture_output=True,text=True)
            if p.returncode:raise RuntimeError("T35_CHILD_PROOF_FAILED:"+arm+":"+str(i)+"\n"+p.stderr)
            reps.append(json.loads(p.stdout))
        first=reps[0]
        if any(x["cases"]!=first["cases"] for x in reps):
            raise RuntimeError("T35_NONDETERMINISTIC_SOURCEFREE_STRESS")
        cases=first["cases"]
        results[arm]={"streams":len(cases),"query_results":sum(len(c["answers"]) for c in cases),
            "correct_certified":sum(sum(x!=-1 for x in c["answers"]) for c in cases),
            "sound_abstentions":sum(sum(x==-1 for x in c["answers"]) for c in cases),
            "all_total_GF2_XORs":sum(c["all_actual_xors"] for c in cases),
            "row_reduction_XORs":sum(c["row_reduction_xors"] for c in cases),
            "cache_materialization_XORs":sum(c["cache_materialization_xors"] for c in cases),
            "cached_multi_pivot_hits":sum(c["cached_block_hits"] for c in cases),
            "five_complete_CPU_ns":[r["complete_native_CPU_ns"] for r in reps],
            "five_peak_RSS_bytes":[r["peak_RSS_bytes"] for r in reps],
            "median_complete_CPU_ns":int(statistics.median(r["complete_native_CPU_ns"] for r in reps)),
            "max_peak_RSS_bytes":max(r["peak_RSS_bytes"] for r in reps)}
    reference=results[ARMS[3]]
    if any((v["correct_certified"],v["sound_abstentions"])!=(
       reference["correct_certified"],reference["sound_abstentions"]) for v in results.values()):
        raise RuntimeError("T35_STRONG_EXACT_COVERAGE_DIFFERENCE")
    learned=results[ARMS[0]]
    exact=results[ARMS[3]]
    cpu=all(learned["median_complete_CPU_ns"]<results[a]["median_complete_CPU_ns"]
            for a in ARMS[1:])
    ramok=learned["max_peak_RSS_bytes"]<=min(v["max_peak_RSS_bytes"] for v in results.values())
    xorok=(learned["all_total_GF2_XORs"]<exact["all_total_GF2_XORs"] and
           audit["learned_minus_fixed_total_XORs"]<0)
    return {"schema":"wingless.T35.sourcefree-block-GF2-native-result.v1",
       "source_free_proof_all_queries":True,"TRAIN_only_public_streams":TRAIN_PER_DIM*len(DIMS),
       "source_free_validation_streams":VALIDATE_PER_DIM*len(DIMS),
       "scientific_DEV_generated":0,"ablation":audit,"arms":results,
       "XOR_gate":xorok,"native_CPU_gate":cpu,"native_RSS_gate":ramok,
       "strict_scientific_source_admission":xorok and cpu and ramok,
       "optimizer_admission":False,"scientific_positive":False}
if __name__=="__main__":
    if len(sys.argv)==1:print(json.dumps(run(),sort_keys=True))
    elif len(sys.argv)==2 and sys.argv[1]=="--ablation":
        print(json.dumps(ablation(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(child(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T35_SOURCE_FREE_ONLY")
