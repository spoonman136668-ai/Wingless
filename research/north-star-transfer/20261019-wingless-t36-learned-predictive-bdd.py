"""T36 source-free learned exact predictive-state ROBDD vs optimized GF2.

Pure exact Boolean-version-space constraint state; NO Gaussian basis in
the candidate or identity-ROBDD paths. A TRAIN-learned variable order
controls BDD reduction, sharing and downstream query costs. An independent
GF2 solver recomputes every source-free query for soundness. No real
unfamiliar causal source or consumed DEV is ever created or opened.
"""
from __future__ import annotations
import itertools,json,random,statistics,subprocess,sys,time
from collections import Counter
from pathlib import Path
DIMS=(16,24)
TRAIN_PER_DIM=20
VALIDATE_PER_DIM=12
OBS=40
STEP=8
QUERIES=8
SEED=360119019
VALIDATION_SEED=360119813
MAX_BDD_NODES=150000
ARMS=("TRAIN-learned-exact-predictive-ROBDD",
      "deterministic-identity-exact-ROBDD",
      "TRAIN-informed-pivot-exact-incremental-GF2",
      "optimized-incremental-exact-GF2",
      "independent-lower-pivot-batch-GF2")

def peak_rss():
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
        ps=ctypes.WinDLL("psapi",use_last_error=True)
        k.GetCurrentProcess.restype=wintypes.HANDLE
        ps.GetProcessMemoryInfo.argtypes=[wintypes.HANDLE,ctypes.POINTER(Mem),wintypes.DWORD]
        ps.GetProcessMemoryInfo.restype=wintypes.BOOL
        m=Mem();m.cb=ctypes.sizeof(m)
        if not ps.GetProcessMemoryInfo(k.GetCurrentProcess(),ctypes.byref(m),m.cb):
            raise RuntimeError("T36_WIN_PROCESS_COUNTER="+str(ctypes.get_last_error()))
        return int(m.PeakWorkingSetSize)
    import resource
    result=resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    return int(result)*1024 if sys.platform.startswith("linux") else int(result)

def indices(mask):
    while mask:
        bit=mask&-mask
        yield bit.bit_length()-1
        mask^=bit

def grammar(dim):
    r=random.Random(SEED+dim*317)
    order=list(range(dim));r.shuffle(order)
    return tuple(tuple(order[j:j+6]) for j in range(0,dim,6))

def sourcefree(dim,index,part):
    if part not in ("TRAIN","VALIDATION"):raise RuntimeError("T36_BAD_SOURCEFREE_PARTITION")
    r=random.Random((SEED if part=="TRAIN" else VALIDATION_SEED)+dim*4021+index*7499)
    blocks=grammar(dim)
    true_law=r.getrandbits(dim)
    equations=[];checkpoints=[]
    for t in range(OBS):
        if r.randrange(14):
            group=r.choice(blocks)
            chosen=r.sample(group,r.choice((2,3,4)))
        else:
            a,b=r.sample(range(len(blocks)),2)
            chosen=r.sample(blocks[a],2)+r.sample(blocks[b],1)
        mask=sum(1<<b for b in chosen)
        equations.append((mask,(mask&true_law).bit_count()&1))
        if (t+1)%STEP==0:
            out=[]
            for k in range(QUERIES):
                if k%4==0:
                    a,b=r.sample(range(len(blocks)),2)
                    select=(r.choice(blocks[a]),r.choice(blocks[b]))
                else:
                    select=r.sample(r.choice(blocks),r.choice((1,2,3)))
                out.append(sum(1<<b for b in select))
            checkpoints.append(tuple(out))
    return tuple(equations),tuple(checkpoints),true_law

def learn_order(dim):
    adjacency=Counter();freq=Counter()
    for i in range(TRAIN_PER_DIM):
        rows,_,_=sourcefree(dim,i,"TRAIN")
        for mask,_ in rows:
            items=list(indices(mask))
            freq.update(items)
            adjacency.update((min(a,b),max(a,b)) for a,b in itertools.combinations(items,2))
    # TRAIN-only learned linear arrangement maximizing local correlation.
    remaining=set(range(dim))
    anchor=max(remaining,key=lambda bit:(freq[bit],-bit))
    ordered=[anchor];remaining.remove(anchor)
    while remaining:
        latest=ordered[-1]
        best=max(remaining,key=lambda bit:(
            3*adjacency[min(latest,bit),max(latest,bit)]+
            sum(adjacency[min(bit,v),max(bit,v)] for v in ordered[-3:]),
            freq[bit],-bit))
        ordered.append(best);remaining.remove(best)
    if sorted(ordered)!=list(range(dim)):
        raise RuntimeError("T36_TRAIN_VARIABLE_ORDER_NOT_BIJECTION")
    sorted_freq=tuple(sorted(range(dim),key=lambda i:(-freq[i],i)))
    return tuple(ordered),sorted_freq

def remap(mask,order):
    # A trained order changes canonical diagram structure, not truth.
    return sum(1<<rank for rank,bit in enumerate(order) if mask&(1<<bit))

class ROBDD:
    """Canonical reduced ordered BDD for conjunction of parity equations."""
    def __init__(self,dim,order):
        self.dim=dim;self.order=order;self.root=1
        self.nodes=[None,None]
        self.unique={}
        self.apply_calls=0;self.apply_misses=0;self.query_visits=0
        self.max_node_count=2
    def node(self,var,low,high):
        if low==high:return low
        key=(var,low,high)
        old=self.unique.get(key)
        if old is not None:return old
        ident=len(self.nodes)
        if ident>=MAX_BDD_NODES:raise RuntimeError("T36_BDD_NODE_BUDGET_EXCEEDED")
        self.nodes.append(key);self.unique[key]=ident
        self.max_node_count=max(self.max_node_count,len(self.nodes))
        return ident
    def eq_bdd(self,mask,rhs):
        positions=list(indices(remap(mask,self.order)))
        memo={}
        def recurse(j,need):
            if j==len(positions):return 1 if need==0 else 0
            key=j,need
            if key not in memo:
                p=positions[j]
                memo[key]=self.node(p,recurse(j+1,need),recurse(j+1,need^1))
            return memo[key]
        return recurse(0,rhs)
    def intersection(self,a,b):
        calls={}
        def and2(x,y):
            self.apply_calls+=1
            if x==0 or y==0:return 0
            if x==1:return y
            if y==1:return x
            if x==y:return x
            key=(min(x,y),max(x,y))
            if key in calls:return calls[key]
            self.apply_misses+=1
            p=min(self.nodes[x][0],self.nodes[y][0])
            def sub(z,branch):
                return z if self.nodes[z][0]!=p else self.nodes[z][branch+1]
            low=and2(sub(x,0),sub(y,0))
            high=and2(sub(x,1),sub(y,1))
            node=self.node(p,low,high)
            calls[key]=node
            return node
        return and2(a,b)
    def add(self,mask,rhs):
        condition=self.eq_bdd(mask,rhs)
        self.root=self.intersection(self.root,condition)
        if self.root==0:raise RuntimeError("T36_TRUE_CAUSAL_STATE_ELIMINATED")
    def query(self,mask):
        requested=remap(mask,self.order)
        memory={}
        def possibilities(node,start):
            self.query_visits+=1
            if node==0:return 0
            if node==1:
                return 3 if requested & ~((1<<start)-1) else 1
            key=(node,start)
            if key in memory:return memory[key]
            var,lo,hi=self.nodes[node]
            # Skipped BDD variables are FREE and can flip parity.
            if requested & (((1<<var)-1)^((1<<start)-1)):
                result=3
            else:
                a=possibilities(lo,var+1)
                b=possibilities(hi,var+1)
                if requested&(1<<var):
                    b=((b&1)<<1)|((b&2)>>1)
                result=a|b
            memory[key]=result
            return result
        mask=possibilities(self.root,0)
        return 0 if mask==1 else 1 if mask==2 else None if mask==3 else RuntimeError("T36_EMPTY_QUERY_BDD")

class GF2:
    def __init__(self,dim,order,low_pivot=False):
        self.dim=dim;self.order=order;self.low_pivot=low_pivot;self.basis={};self.xors=0
    def pivot(self,mask):
        return (mask&-mask).bit_length()-1 if self.low_pivot else mask.bit_length()-1
    def add(self,mask,rhs):
        mask=remap(mask,self.order)
        while mask:
            p=self.pivot(mask)
            if p not in self.basis:
                self.basis[p]=(mask,rhs);return
            m,bit=self.basis[p];mask^=m;rhs^=bit;self.xors+=1
        if rhs:raise RuntimeError("T36_GF2_TRUE_HYPOTHESIS_LOST")
    def query(self,mask):
        mask=remap(mask,self.order);rhs=0
        while mask:
            p=self.pivot(mask)
            if p not in self.basis:return None
            m,b=self.basis[p];mask^=m;rhs^=b
            self.xors+=1
        return rhs

def independent_exact(rows,query,dim):
    # Orthogonal low-pivot rebuild (not a BDD or the optimized upper GF2).
    basis={}
    for m,r in rows:
        while m:
            p=(m&-m).bit_length()-1
            if p not in basis:
                basis[p]=(m,r);break
            x,y=basis[p];m^=x;r^=y
        else:
            if r:raise RuntimeError("T36_ORACLE_INCONSISTENT")
    rhs=0
    while query:
        p=(query&-query).bit_length()-1
        if p not in basis:return None
        m,v=basis[p];query^=m;rhs^=v
    return rhs

def evaluate_stream(dim,index,part,arm,learned,frequency):
    equations,queries,true_law=sourcefree(dim,index,part)
    if arm==ARMS[0]:state=ROBDD(dim,learned)
    elif arm==ARMS[1]:state=ROBDD(dim,tuple(range(dim)))
    elif arm==ARMS[2]:state=GF2(dim,frequency)
    else:state=GF2(dim,tuple(range(dim)),low_pivot=(arm==ARMS[4]))
    seen=[];out=[];checkpoint=0
    for t,(mask,rhs) in enumerate(equations):
        seen.append((mask,rhs))
        if arm!=ARMS[4]:state.add(mask,rhs)
        if (t+1)%STEP:continue
        for query in queries[checkpoint]:
            if arm==ARMS[4]:
                # The independent comparator recomputes from all history.
                exact=independent_exact(seen,query,dim)
                result=exact
            else:result=state.query(query)
            expected=independent_exact(seen,query,dim)
            if result!=expected:
                raise RuntimeError("T36_INDEPENDENT_EXACTNESS_FAIL:"+arm+":"+str(dim)+":"+str(index))
            if result is not None and result!=((query&true_law).bit_count()&1):
                raise RuntimeError("T36_FALSE_CERTIFIED_QUERY")
            out.append(-1 if result is None else result)
        checkpoint+=1
    if len(out)!=(OBS//STEP)*QUERIES:raise RuntimeError("T36_MISSING_QUERIES")
    if isinstance(state,ROBDD):
        return {"dim":dim,"id":index,"outputs":out,"bdd_nodes":len(state.nodes),
           "bdd_apply_misses":state.apply_misses,
           "bdd_apply_calls":state.apply_calls,"bdd_query_visits":state.query_visits,
           "basis_xors":0}
    return {"dim":dim,"id":index,"outputs":out,"bdd_nodes":0,
           "bdd_apply_misses":0,"bdd_apply_calls":0,"bdd_query_visits":0,
           "basis_xors":state.xors}

def audit():
    changed=more=less=total_on=total_off=proofs=0
    for dim in DIMS:
        learned,frequency=learn_order(dim)
        for idx in range(VALIDATE_PER_DIM):
            on=evaluate_stream(dim,idx,"VALIDATION",ARMS[0],learned,frequency)
            off=evaluate_stream(dim,idx,"VALIDATION",ARMS[1],learned,frequency)
            if on["outputs"]!=off["outputs"]:raise RuntimeError("T36_ABLATION_CHANGED_EXACT_CERTIFICATION")
            if (on["bdd_nodes"],on["bdd_apply_misses"])!=(
               off["bdd_nodes"],off["bdd_apply_misses"]):changed+=1
            if on["bdd_nodes"]<off["bdd_nodes"]:less+=1
            if on["bdd_nodes"]>off["bdd_nodes"]:more+=1
            total_on+=on["bdd_nodes"];total_off+=off["bdd_nodes"]
            proofs+=len(on["outputs"])
    return {"schema":"wingless.T36.real-learned-ROBDD-state-structure-ablation.v1",
      "TRAIN_streams":TRAIN_PER_DIM*len(DIMS),"sourcefree_streams":VALIDATE_PER_DIM*len(DIMS),
      "independent_exact_checkpoint_queries":proofs,
      "learned_order_changes_actual_BDD_structure_streams":changed,
      "learned_uses_fewer_allocated_nodes_streams":less,
      "learned_uses_more_allocated_nodes_streams":more,
      "learned_total_allocated_nodes":total_on,"identity_total_allocated_nodes":total_off,
      "exact_certifications_preserved":True,"scientific_DEV_created":0}

def process(arm,rep):
    if arm not in ARMS or rep not in range(5):raise RuntimeError("T36_BAD_TRIAL")
    begin=time.process_time_ns()
    cases=[]
    for dim in DIMS:
        learned,frequency=learn_order(dim)
        for idx in range(VALIDATE_PER_DIM):
            cases.append(evaluate_stream(dim,idx,"VALIDATION",arm,learned,frequency))
    return {"schema":"wingless.T36.sourcefree-process.v1",
       "arm":arm,"repeat":rep,"cases":cases,
       "complete_process_CPU_ns":time.process_time_ns()-begin,
       "max_peak_process_RSS_bytes":peak_rss(),"scientific_DEV_created":0}

def run():
    ab=audit();scores={}
    for arm in ARMS:
        trials=[]
        for rep in range(5):
            r=subprocess.run([sys.executable,str(Path(__file__).resolve()),
                 "--child",arm,str(rep)],capture_output=True,text=True)
            if r.returncode:
                raise RuntimeError("T36_SUBPROCESS_FAILED:"+arm+":"+str(rep)+"\n"+r.stderr[-4000:])
            trials.append(json.loads(r.stdout))
        c=trials[0]["cases"]
        if any(x["cases"]!=c for x in trials):raise RuntimeError("T36_NONDETERMINISTIC_FROZEN_STRESS")
        scores[arm]={"sourcefree_cases":len(c),
            "query_count":sum(len(x["outputs"]) for x in c),
            "certified":sum(sum(v!=-1 for v in x["outputs"]) for x in c),
            "sound_abstentions":sum(sum(v==-1 for v in x["outputs"]) for x in c),
            "BDD_nodes_allocated":sum(x["bdd_nodes"] for x in c),
            "BDD_apply_misses":sum(x["bdd_apply_misses"] for x in c),
            "BDD_query_visits":sum(x["bdd_query_visits"] for x in c),
            "basis_XORs":sum(x["basis_xors"] for x in c),
            "five_full_CPU_ns":[x["complete_process_CPU_ns"] for x in trials],
            "five_peak_RSS_bytes":[x["max_peak_process_RSS_bytes"] for x in trials],
            "median_full_CPU_ns":int(statistics.median(x["complete_process_CPU_ns"] for x in trials)),
            "max_peak_RSS_bytes":max(x["max_peak_process_RSS_bytes"] for x in trials)}
    baseline=scores[ARMS[3]]
    if any((s["certified"],s["sound_abstentions"])!=(
        baseline["certified"],baseline["sound_abstentions"]) for s in scores.values()):
        raise RuntimeError("T36_EXACT_CERTIFIED_COUNTS_MISMATCH")
    learned=scores[ARMS[0]]
    admission=(ab["learned_order_changes_actual_BDD_structure_streams"]>0 and
       all(learned["median_full_CPU_ns"]<scores[a]["median_full_CPU_ns"] for a in ARMS[1:]) and
       learned["max_peak_RSS_bytes"]<=min(x["max_peak_RSS_bytes"] for x in scores.values()))
    return {"schema":"wingless.T36.sourcefree-learned-BDD-inference-native-feasibility.v1",
       "TRAIN_public_streams":TRAIN_PER_DIM*len(DIMS),
       "sourcefree_validation_streams":VALIDATE_PER_DIM*len(DIMS),
       "scientific_DEV_created":0,"sourcefree_all_queries_exact":True,
       "ablation":ab,"arms":scores,
       "strict_scientific_source_admission":admission,
       "scientific_positive":False,"optimizer_admission":False,
       "conclusion":"PRE_SOURCE_FEASIBLE" if admission else "PRE_SOURCE_NATIVE_NO_GO"}
if __name__=="__main__":
    if len(sys.argv)==2 and sys.argv[1]=="--ablation":
        print(json.dumps(audit(),sort_keys=True))
    elif len(sys.argv)==1:
        print(json.dumps(run(),sort_keys=True))
    elif len(sys.argv)==4 and sys.argv[1]=="--child":
        print(json.dumps(process(sys.argv[2],int(sys.argv[3])),sort_keys=True))
    else:raise SystemExit("T36_SOURCEFREE_MODES_ONLY")
