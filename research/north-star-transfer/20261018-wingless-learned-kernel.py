"""T1 source-blind learned recurrent affine composition: implementation only.

The 12 signed int8 prior parameters learn from PUBLIC calibration pairs
without reading true operator classes or a source family. No new real
TRAIN, VAL, SEALED or replication source is imported here. All resource
claims beyond the packed 15-byte state are explicitly UNQUALIFIED.
"""
from __future__ import annotations
from dataclasses import dataclass
import math

PRIME=7
SLOPES=tuple(range(1,7))
PARAMETER_BYTES=12
RECURRENT_BYTES=3
PACKED_TOTAL_BYTES=15
MAX_PACKED_BYTES=23

def need(condition,tag):
    if not condition:raise ValueError("WING_TRANSFER_"+tag)

def int8(x):
    return max(-128,min(127,int(round(x))))

def public_features(pairs):
    need(type(pairs) in (list,tuple) and 2<=len(pairs)<=8,"PUBLIC_PAIR_COUNT")
    need(all(len(p)==2 and all(type(x) is int and 0<=x<PRIME for x in p) for p in pairs),
         "MOD7_PUBLIC_EXAMPLES")
    a=[]
    for slope in SLOPES:
        hist={}
        for x,y in pairs:
            shift=(y-slope*x)%PRIME
            hist[shift]=hist.get(shift,0)+1
        a.append(max(hist.values())/len(pairs))
    return tuple(a)

def learn_shared_prior(public_tasks,passes=2):
    """Train with public (x,y) observations only; no hidden class oracle."""
    need(type(public_tasks) in (list,tuple) and 1<=len(public_tasks)<=4096,"TRAIN_ONLY_INPUT_BOUND")
    need(passes==2,"FROZEN_TWO_PASSES")
    weights=[0.0]*6;biases=[0.0]*6
    work=0
    for _ in range(passes):
        for samples in public_tasks:
            features=public_features(samples)
            # Soft pseudo-posterior learned by minimizing public reconstruction
            # inconsistency; no secretly supplied true slope labels.
            public_err=[len(samples)*(1.0-f) for f in features]
            teacher=[math.exp(-2*v) for v in public_err]
            z=sum(teacher)
            teacher=[v/z for v in teacher]
            raw=[biases[i]+weights[i]*features[i] for i in range(6)]
            m=max(raw)
            pred=[math.exp(v-m) for v in raw]
            q=sum(pred);pred=[v/q for v in pred]
            for i in range(6):
                diff=pred[i]-teacher[i]
                weights[i]-=0.08*diff*features[i]
                biases[i]-=0.08*diff
                work+=1
    packed=bytes((int8(16*x)&255) for x in weights+biases)
    need(len(packed)==12,"MODEL_PARAMETERS_12_BYTES")
    return packed,{"training_parameter_updates":work,
                   "training_cases":len(public_tasks),"passes":passes,
                   "extra_training_transient_bytes_not_qualified":True,
                   "scientific_transfer_qualified":False}

def unpack(packed):
    need(type(packed) is bytes and len(packed)==12,"EXACT_WEIGHT_BYTES")
    vals=[x if x<128 else x-256 for x in packed]
    return vals[:6],vals[6:]

def public_role_model(samples,packed):
    features=public_features(samples)
    w,b=unpack(packed)
    # Match public consistency before allowing learned prior to break ties.
    counts=[len(samples)*f for f in features]
    costs=[(len(samples)-counts[i])*64-(b[i]+w[i]*features[i])/16
           for i in range(6)]
    k=min(range(6),key=lambda i:(costs[i],i))
    a=SLOPES[k]
    counts_shift=[0]*PRIME
    for x,y in samples:counts_shift[(y-a*x)%PRIME]+=1
    shift=min(range(PRIME),key=lambda s:(-counts_shift[s],s))
    return a,shift,{"public_errors":len(samples)-counts_shift[shift],
                    "hidden_family_used":False,"role_table_persisted":False}

def compose(state,role,packed,inverse=False):
    """Carry output=a*input+b in a 3-byte source-independent recurrent state."""
    need(type(state) is bytes and len(state)==3
         and state[0] in SLOPES and state[1]<PRIME and state[2] in (0,1),
         "RECURRENT_EXACT_PACK")
    a,b,_=state
    u,v,_=public_role_model(role,packed)
    if inverse:
        u=pow(u,-1,PRIME)
        v=(-u*v)%PRIME
    return bytes(((u*a)%PRIME,(u*b+v)%PRIME,1)),{
        "persistent_bytes":PACKED_TOTAL_BYTES,
        "new_public_role_examples":len(role),
        "source_true_label_access":False}

def start():
    return bytes((1,0,0))

def apply(state,x):
    need(type(x) is int and 0<=x<PRIME,"PUBLIC_QUERY")
    return (state[0]*x+state[1])%PRIME

def selftest():
    # Synthetic closure tests only: no cohort ids or generator.
    public=[[(0,2),(1,4),(2,6)],
            [(0,5),(1,1),(2,4)],
            [(0,1),(1,5),(2,2)]]
    packed,meta=learn_shared_prior(public)
    need(len(packed)+len(start())==15 and meta["training_parameter_updates"]==36,
         "MODEL_BINARY_AND_TRAINING_BUDGET")
    for c in public:
        state,_=compose(start(),c,packed)
        inverse,_=compose(state,c,packed,inverse=True)
        need(all(apply(inverse,x)==x for x in range(PRIME)),"INVERSE_COMPOSITION")
    a,_=compose(start(),public[0],packed)
    b,_=compose(a,public[1],packed)
    u,v,_=public_role_model(public[0],packed)
    p,q,_=public_role_model(public[1],packed)
    need(all(apply(b,x)==(p*(u*x+v)+q)%PRIME
             for x in range(PRIME)),"PUBLIC_TWO_ROLE_COMPOSITION")
    print("WING_T1_PUBLIC_ONLY_COMPOSITION_SELFTEST=PASS")
    print("WING_T1_PERSISTENT_MODEL_AND_STATE_BYTES=15")
    print("WING_T1_FULL_COMPUTE_MEMORY_QUALIFIED=false")
    print("WING_T1_REAL_SOURCE_OR_HOLDOUT_CONSUMED=false")
if __name__=="__main__":selftest()
