"""T41 source-free, new generator. No previous T40 public META or DEV replay.

The only outputs given to the learner are (mask, revealed signature) pairs.
Hidden truth tables and private source regime labels stay within scorer.
Not independently originated scientific first-use.
"""
from __future__ import annotations
import random

REGIMES=("unrestricted256","uniform_composition768","nonuniform_composition768")
N_TRAIN=64
N_EVAL=96
SEED=202610204141
MIX=10657
COMPOSED_ROLES=16

def nonuniform_program_weights(core):
    # Deliberately different from the T40 private source weights.
    values=[(1+(11*g+3)%17)*(1+(7*h+5)%19)*(3,5,2)[op]
            for g,h,op in core.PROGRAMS]
    total=sum(values)
    return tuple(z/total for z in values)

def choice_weighted(rng,weights):
    threshold=rng.random()
    for i,w in enumerate(weights):
        threshold-=w
        if threshold<0:
            return i
    return len(weights)-1

def choose_rule(regime,phase,index,rng,core,program_weights):
    if regime=="unrestricted256":
        # Source-scoped stratification of the eval cohort alone guarantees
        # genuinely out-of-grammar laws at one quarter of eval positions.
        if phase=="eval":
            representable=set(core.LAW)
            outside=tuple(i for i in range(256) if i not in representable)
            inside=tuple(sorted(representable))
            assert len(outside)==64 and len(inside)==192
            if index%4==0:
                return rng.choice(outside),None
            return rng.choice(inside),None
        return rng.randrange(256),None
    # Forced coverage is a source-only generation rule, never shown to learner.
    if phase=="train" and index<16:
        op=rng.randrange(3)
        pid=core.PROGRAMS.index((index,(index+3)%16,op))
    elif phase=="train" and index<32:
        op=rng.randrange(3)
        pid=core.PROGRAMS.index(((index+1)%16,index-16,op))
    elif regime=="uniform_composition768":
        pid=rng.randrange(768)
    else:
        pid=choice_weighted(rng,program_weights)
    return core.LAW[pid],pid

def cases(regime,phase,core):
    if regime not in REGIMES or phase not in ("train","eval"):
        raise ValueError("T41_SOURCE_SCOPE_INVALID")
    weights=nonuniform_program_weights(core)
    n=N_TRAIN if phase=="train" else N_EVAL
    output=[]
    coverage_g=set();coverage_h=set();outside=0
    for i in range(n):
        realm=REGIMES.index(regime)
        phase_shift=0 if phase=="train" else 773933
        rng=random.Random(SEED+realm*3918491+phase_shift+i*MIX)
        law,pid=choose_rule(regime,phase,i,rng,core,weights)
        if pid is not None:
            g,h,_=core.PROGRAMS[pid]
            coverage_g.add(g);coverage_h.add(h)
        if law not in set(core.LAW):
            outside+=1
        if phase=="train":
            spots=rng.sample(range(8),4)
            mask=sum(1<<q for q in spots)
            output.append({"id":"t41-public-train-%s-%03d"%(regime,i),
                           "public":(mask,law&mask),
                           "source_birth":"SOURCEFREE_META_T41_ONLY"})
        else:
            evidence={}
            for k in (2,4,6):
                mask=sum(1<<q for q in rng.sample(range(8),k))
                evidence[k]=(mask,law&mask)
            output.append({"id":"t41-sourcefree-eval-%s-%03d"%(regime,i),
                           "public_observations":evidence,
                           "private_truth_rule":law,
                           "source_birth":"SOURCEFREE_INTERNAL_EVAL_NOT_SCIENCE"})
    if regime!="unrestricted256" and phase=="train":
        assert len(coverage_g)==len(coverage_h)==16
    if regime=="unrestricted256" and phase=="eval":
        assert outside==N_EVAL//4
    return tuple(output),{"regime_private_auditor_only":regime,"phase":phase,
            "sourcefree_cases":n,"coverage_G":len(coverage_g),
            "coverage_H":len(coverage_h),"outside_composition":outside,
            "no_T40_case_replay":True,"independent_scientific_source":False}
