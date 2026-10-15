"""W-F36: learned soft typed relationship composition (TRAIN-only).

Role/subject/object triples are public; private 8-node bijections are source
truth only. A gradient-trained 8x8x8 row-stochastic operator is quantized to
512 bytes, then used to carry a distribution of possible variable bindings
across ordered forward/inverse role queries. Doubly-stochastic permutation MAP
and affine/XOR MAP see identical public triple evidence. Programmed synthetic
development only; not an independent F3 original or a claim about Python heap.
"""
from __future__ import annotations
import hashlib, json, math, random
from functools import lru_cache

TRAIN_SEEDS=range(241000,241012)
TRAIN_FAMILIES=(210,211,212)
RESERVED_SEEDS=range(242000,242016)
RESERVED_FAMILIES=(214,215,216,217)
PAST_RESERVED=range(202000,202016)
DEPTHS=(4,6,8)
ROLES=8; NODES=8; DEMOS=5; EPOCHS=90
ARMS=("learned_soft","permutation_MAP","affine_MAP","frequency_walk","role_scrambled","role_agnostic","privileged_oracle")
BYTES_Q8=ROLES*NODES*NODES
BYTES_WORK=2*NODES
NATIVE_BYTES=576

def rng(seed,family,salt):
    if seed not in TRAIN_SEEDS or family not in TRAIN_FAMILIES:
        raise ValueError("unregistered source access (including reserved original)")
    return random.Random(seed*74303+family*19433+salt*991+47369)

def family_graph(seed,family,depth):
    r=rng(seed,family,depth)
    grammar=[]
    for role in range(ROLES):
        if family==210:
            a=r.choice((1,3,5,7));b=r.randrange(8)
            values=[(a*x+b)%NODES for x in range(NODES)]
        elif family==211:
            values=list(range(NODES));r.shuffle(values)
        else:
            a=r.choice((1,3,5,7));b=r.randrange(8);mask=r.randrange(1,8)
            values=[((a*x+b)%NODES)^mask for x in range(NODES)]
        grammar.append(tuple(values))
    public=[]
    for role in range(ROLES):
        for d in range(DEMOS):
            src=r.randrange(NODES);dst=grammar[role][src]
            if r.random()<.09:dst=(dst+r.randrange(1,NODES))%NODES
            public.append((role,src,dst))
    r.shuffle(public)
    start=r.randrange(NODES)
    # Query specification is public, including direction, but NOT answer.
    query=tuple((r.randrange(ROLES),1 if r.random()<.7 else -1) for _ in range(depth))
    return tuple(grammar),tuple(public),start,query

def target_oracle(grammar,start,query):
    current=start
    for role,direction in query:
        mapping=grammar[role]
        if direction>0:current=mapping[current]
        else:current=mapping.index(current)
    return current

def independent_oracle(grammar,start,query):
    possibilities=[1 if i==start else 0 for i in range(NODES)]
    for role,direction in query:
        mapping=grammar[role];following=[0]*NODES
        for x,weight in enumerate(possibilities):
            dest=mapping[x] if direction>0 else mapping.index(x)
            following[dest]+=weight
        possibilities=following
    return possibilities.index(1)

def softmax(values):
    m=max(values);e=[math.exp(x-m) for x in values];den=sum(e)
    return [x/den for x in e]

def learned_quantized(public):
    # The 512 transition logits are updated by gradient on public triples ONLY.
    logits=[[[0.0]*NODES for _ in range(NODES)] for _ in range(ROLES)]
    for epoch in range(EPOCHS):
        rate=.16/(1+epoch/90)
        for role,src,dst in public:
            row=logits[role][src];p=softmax(row)
            for j in range(NODES):
                row[j]-=rate*((p[j]-(j==dst))+.0008*row[j])
        if epoch%5==0:
            # Differentiable column-collision penalty, independent of secrets.
            for role in range(ROLES):
                probs=[softmax(x) for x in logits[role]]
                columns=[sum(probs[x][j] for x in range(NODES)) for j in range(NODES)]
                for x in range(NODES):
                    for j in range(NODES):
                        logits[role][x][j]-=.015*(columns[j]-1)*probs[x][j]
    # Each transition probability occupies 1 unsigned byte.
    q=[]
    for role in range(ROLES):
        rows=[]
        for src in range(NODES):
            p=softmax(logits[role][src])
            vals=[round(255*x) for x in p]
            # Normalize the last largest component to sum exactly 255.
            index=max(range(NODES),key=lambda j:vals[j])
            vals[index]+=255-sum(vals)
            if not all(0<=v<=255 for v in vals) or sum(vals)!=255:
                raise ValueError("quantized source uncertainty out of range")
            rows.append(tuple(vals))
        q.append(tuple(rows))
    return tuple(q)

def uncertain_compose(q,start,query):
    state=[0.0]*NODES;state[start]=1.
    for role,direction in query:
        next_state=[0.0]*NODES
        if direction>0:
            for src in range(NODES):
                for dst in range(NODES):
                    next_state[dst]+=state[src]*(q[role][src][dst]/255.)
        else:
            # Column-normalized inverse evidence: sum all possible current
            # bindings, never assume the role transition is a hard bijection.
            denominators=[sum(q[role][x][dst] for x in range(NODES)) for dst in range(NODES)]
            for dst in range(NODES):
                for original in range(NODES):
                    next_state[original]+=state[dst]*q[role][original][dst]/max(1,denominators[dst])
        total=sum(next_state)
        if total<=0:raise ValueError("collapsed probability mass")
        state=[x/total for x in next_state]
    return max(range(NODES),key=lambda x:(state[x],-x))

def assignment(counts):
    @lru_cache(None)
    def solve(src,used):
        if src==NODES:return (0.,())
        choices=[]
        for dst in range(NODES):
            if (used>>dst)&1:continue
            value,path=solve(src+1,used|(1<<dst))
            choices.append((value+counts[src][dst],(dst,)+path))
        return max(choices,key=lambda t:(t[0],tuple(-x for x in t[1])))
    return solve(0,0)[1]

def candidates():
    s=set()
    for a in (1,3,5,7):
        for b in range(NODES):
            s.add(tuple((a*x+b)%NODES for x in range(NODES)))
            for mask in range(1,NODES):s.add(tuple(((a*x+b)%NODES)^mask for x in range(NODES)))
    return tuple(sorted(s))
CANDIDATES=candidates()

def hard_controls(public):
    counts=[[[0.]*NODES for _ in range(NODES)] for _ in range(ROLES)]
    for role,src,dst in public:counts[role][src][dst]+=1
    hard=tuple(assignment(counts[role]) for role in range(ROLES))
    affine=[]
    for role in range(ROLES):
        obs=[(s,d) for r,s,d in public if r==role]
        fitted=max(CANDIDATES,key=lambda op:(sum(op[s]==d for s,d in obs),tuple(-x for x in op)))
        affine.append(fitted)
    return hard,tuple(affine),counts

def source_walk(counts,start,query):
    # Source-directed, unstructured frequency baseline.
    state=start
    for role,direction in query:
        if direction>0:
            state=max(range(NODES),key=lambda dst:(counts[role][state][dst],-dst))
        else:
            state=max(range(NODES),key=lambda src:(counts[role][src][state],-src))
    return state

def score(rows,arm):
    correct=sum(r["predictions"][arm]==r["truth"] for r in rows)
    return dict(correct=correct,total=len(rows),accuracy=correct/len(rows))

def run():
    if (set(TRAIN_SEEDS)&set(RESERVED_SEEDS) or set(TRAIN_FAMILIES)&set(RESERVED_FAMILIES) or
        set(TRAIN_SEEDS)&set(PAST_RESERVED)):
        raise ValueError("identity collision")
    rows=[]
    for seed in TRAIN_SEEDS:
      for family in TRAIN_FAMILIES:
       for depth in DEPTHS:
        grammar,public,start,query=family_graph(seed,family,depth)
        truth=target_oracle(grammar,start,query)
        if independent_oracle(grammar,start,query)!=truth:raise ValueError("dual oracle mismatch")
        q=learned_quantized(public)
        hard,affine,counts=hard_controls(public)
        public_hash=hashlib.sha256(json.dumps([public,start,query]).encode()).hexdigest()
        scrambled=tuple(reversed(query))
        predictions=dict(
            learned_soft=uncertain_compose(q,start,query),
            permutation_MAP=target_oracle(hard,start,query),
            affine_MAP=target_oracle(affine,start,query),
            frequency_walk=source_walk(counts,start,query),
            role_scrambled=uncertain_compose(q,start,scrambled),
            role_agnostic=start,
            privileged_oracle=truth)
        rows.append(dict(seed=seed,family=family,depth=depth,query=query,start=start,truth=truth,
            predictions=predictions,private_order_counterfactual_changes_answer=target_oracle(grammar,start,scrambled)!=truth,
            public_evidence_sha256=public_hash,public_triples_count=len(public),
            quantized_model_bytes=BYTES_Q8))
    scores={a:score(rows,a) for a in ARMS}
    per_family={str(f):{a:score([r for r in rows if r["family"]==f],a) for a in ARMS}
                for f in TRAIN_FAMILIES}
    fair=("permutation_MAP","affine_MAP","frequency_walk","role_scrambled","role_agnostic")
    best=max(scores[a]["accuracy"] for a in fair)
    improvement=scores["learned_soft"]["accuracy"]-best
    robust=all(per_family[str(f)]["learned_soft"]["accuracy"]>=max(per_family[str(f)][a]["accuracy"] for a in fair) for f in TRAIN_FAMILIES)
    eligible=improvement>=.03 and robust
    return dict(schema="wingless.f36.soft-role-binding.trainonly.v1",
        stage="NEW_TRAIN_ONLY_NO_PRIMARY_ORIGINAL",
        train_seeds=[TRAIN_SEEDS.start,TRAIN_SEEDS.stop],train_families=TRAIN_FAMILIES,
        reserved_original_seeds=[RESERVED_SEEDS.start,RESERVED_SEEDS.stop],
        reserved_original_families=RESERVED_FAMILIES,
        questions=len(rows),depths=DEPTHS,model_q8_bytes=BYTES_Q8,
        inference_distribution_scratch_bytes=BYTES_WORK,
        native_model_and_work_bytes=BYTES_Q8+BYTES_WORK,native_limit_bytes=NATIVE_BYTES,
        training_float_optimizer_external=True,external_triple_storage_bytes_per_graph=ROLES*DEMOS*3*8,
        ordered_path_counterfactual_changes=sum(r["private_order_counterfactual_changes_answer"] for r in rows),
        independent_two_oracle_agreement=True,source_independently_born_F3=False,
        F3_accepted=False,original_executed=False,accepted_ref_mutation=False,
        scores=scores,per_family=per_family,rows=rows,
        verdict="TRAIN_ONLY_ROLE_COMPOSITION_LEARNED_QUALIFIED" if eligible else "TRAIN_ONLY_ROLE_COMPOSITION_NO_QUALIFIED_GAIN")
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
