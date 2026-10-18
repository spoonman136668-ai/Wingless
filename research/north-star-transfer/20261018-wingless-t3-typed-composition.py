"""T3 Wingless learned PUBLIC role-prior with typed F11 and S3 composition.

Corrects T1's F7 vs frozen W-F55 F11 mismatch and F7 inability to express
S3. Inference receives a PUBLIC algebra descriptor, never an oracle family
identity, held answer, true operator or source-generation key. This is a
nonconsuming mathematical capacity proof, not transfer evidence.
"""
from __future__ import annotations
from itertools import permutations
import math

AFFINE=0
NONCOMMUTING_PERM=1
P11=11
PERMS=tuple(permutations((0,1,2)))
PARAM_BYTES=12
STATE_BYTES=3
PERM_ID={p:i for i,p in enumerate(PERMS)}
def need(ok,why):
    if not ok:raise ValueError("WING_T3_"+why)
def signed(data):
    need(type(data) is bytes and len(data)==PARAM_BYTES,"TRAINED_12_WEIGHT_BYTES")
    return [v if v<128 else v-256 for v in data]
def flatten(cases):
    # 12 trainable signed confidence/prior units over PUBLIC consistency
    # margins. Public pseudolabels only; zero hidden operator supervision.
    score=[0.0]*12
    need(isinstance(cases,(list,tuple)) and 1<=len(cases)<=2048,
         "ONLY_PUBLIC_TRAIN_CASES")
    steps=0
    for _ in range(2):
        for domain,pairs in cases:
            choices=candidates(domain)
            agreements=[sum(operator_apply(domain,k,x)==y for x,y in pairs)
                        for k in choices]
            need(len(pairs) in (2,3),"FROZEN_PUBLIC_EXAMPLES")
            best=max(agreements)
            for j in range(12):
                # fixed public histogram; no true slope/permutation labels.
                feature=(agreements[j%len(agreements)]/len(pairs))
                score[j]+=0.03125*(feature-best/len(pairs)/2)
                steps+=1
    encoded=bytes(max(0,min(255,round(v*32)+128)) for v in score)
    need(len(encoded)==12,"EXACT_SHARED_PRIOR_BYTES")
    return encoded,{"public_only_training_feature_updates":steps,
                    "cross_family_training_executed":False,
                    "full_process_cost_qualified":False}
def candidates(domain):
    if domain==AFFINE:return tuple((a,b) for a in range(1,P11) for b in range(P11))
    need(domain==NONCOMMUTING_PERM,"ONLY_REGISTERED_PUBLIC_DOMAIN")
    return PERMS
def operator_apply(domain,k,x):
    if domain==AFFINE:
        a,b=k;need(isinstance(x,int) and 0<=x<P11,"F11_INPUT")
        return (a*x+b)%P11
    need(domain==NONCOMMUTING_PERM and 0<=x<3,"S3_INPUT")
    return k[x]
def infer_role(domain,pairs,prior):
    need(domain in (AFFINE,NONCOMMUTING_PERM) and type(pairs) in (tuple,list)
         and len(pairs) in (2,3),"PUBLIC_ROLE_CALIBRATION")
    bounds=P11 if domain==AFFINE else 3
    need(all(len(z)==2 and 0<=z[0]<bounds and 0<=z[1]<bounds
             for z in pairs),"PUBLIC_IN_DOMAIN")
    weights=signed(prior)
    choices=candidates(domain)
    # Public evidence, not hidden source type, dominates weights.
    # This necessarily ties an exact symbolic control when uncorrupted.
    def evidence(k):
        fits=sum(operator_apply(domain,k,x)==y for x,y in pairs)
        code=PERM_ID[k] if domain==NONCOMMUTING_PERM else k[0]*P11+k[1]
        return (fits*1024+weights[code%12],-code)
    winner=max(choices,key=evidence)
    return winner,{"public_pairs":len(pairs),
                   "correct_public_pairs":sum(operator_apply(domain,winner,x)==y
                                              for x,y in pairs),
                   "hidden_class_used":False}
def compose_state(state,operator,domain,inverse=False):
    need(type(state) is bytes and len(state)==STATE_BYTES,"EXACT_STATE")
    need(state[0]==domain and domain in (AFFINE,NONCOMMUTING_PERM),
         "NO_CROSS_DOMAIN_COERCION")
    if domain==AFFINE:
        a,b=state[1:3];u,v=operator
        need(1<=a<P11 and 0<=b<P11 and 1<=u<P11 and 0<=v<P11,
             "F11_AFFINE_BIJECTION")
        if inverse:u,v=pow(u,-1,P11),(-v*pow(u,-1,P11))%P11
        return bytes((domain,(u*a)%P11,(u*b+v)%P11))
    need(state[1] in range(6) and state[2]==0
         and tuple(operator) in PERM_ID,"S3_PERMUTATION_BIJECTION")
    op=tuple(operator)
    if inverse:
        inv=[0,0,0]
        for i,j in enumerate(op):inv[j]=i
        op=tuple(inv)
    old=PERMS[state[1]]
    new=tuple(op[old[i]] for i in range(3))
    return bytes((domain,PERM_ID[new],0))
def start(domain):
    need(domain in (AFFINE,NONCOMMUTING_PERM),"DOMAIN_REQUIRED")
    return bytes((domain,1,0)) if domain==AFFINE else bytes((domain,PERM_ID[(0,1,2)],0))
def apply(state,x):
    need(type(state) is bytes and len(state)==3,"PACKED_STATE")
    if state[0]==AFFINE:return (state[1]*x+state[2])%P11
    need(state[0]==NONCOMMUTING_PERM and state[1]<6,"S3_STATE")
    return PERMS[state[1]][x]
def selftest():
    prior,meta=flatten([
        (AFFINE,[(0,3),(1,5),(2,7)]),
        (NONCOMMUTING_PERM,[(0,1),(1,2),(2,0)])])
    need(len(prior)==12 and meta["public_only_training_feature_updates"]==48,
         "TRAINED_WEIGHT_SHAPE")
    for domain in (AFFINE,NONCOMMUTING_PERM):
        for op in candidates(domain):
            state=compose_state(start(domain),op,domain)
            restored=compose_state(state,op,domain,inverse=True)
            need(restored==start(domain),"INVERTIBLE_TYPED_ROLE_COMPOSITION")
    p=PERM_ID[(1,2,0)];q=PERM_ID[(1,0,2)]
    left=compose_state(compose_state(start(1),PERMS[p],1),PERMS[q],1)
    right=compose_state(compose_state(start(1),PERMS[q],1),PERMS[p],1)
    need(left!=right,"TRUE_NONCOMMUTING_COMPOSITION")
    learned,_=infer_role(0,[(0,3),(1,5),(2,7)],prior)
    need(learned==(2,3),"F11_COMPATIBILITY_RESTORED")
    perm,_=infer_role(1,[(0,1),(1,2),(2,0)],prior)
    need(perm==(1,2,0),"S3_TYPED_PUBLIC_INFERENCE")
    print("T3_WINGLESS_F11_AND_S3_REPRESENTATION=PASS")
    print("T3_WINGLESS_PRIOR_PLUS_PACKED_STATE_BYTES=15")
    print("T3_WINGLESS_S3_NONCOMMUTING=true")
    print("T3_WINGLESS_STRICT_OOD_GAIN_OVER_EXACT_PUBLIC_SOLVER=false")
    print("T3_WINGLESS_FULL_RESOURCE_PARITY=false")
    print("T3_WINGLESS_SCIENTIFIC_FIRST_USE=false")
if __name__=="__main__":selftest()
