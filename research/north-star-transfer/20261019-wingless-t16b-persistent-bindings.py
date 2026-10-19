"""Non-consuming persistent variable-binding certificate for Wingless.

Mathematically exact shared operator identities across all uses of one role.
The optional branch budget returns UNKNOWN, never a guessed answer. The
public two-pair full solver is an independent exact fallback.
No experimental TRAIN/OPEN source is imported or originated.
"""
from itertools import permutations

PERMS=tuple(permutations(range(3)))
UNKNOWN=None

def consistent_ops(domain, pairs):
    assert domain in (0,1)
    assert 1<=len(pairs)<=2
    n=11 if domain==0 else 3
    assert all(0<=x<n and 0<=y<n for x,y in pairs)
    if domain==0:
        x,y=pairs[0]
        ops=((a,(y-a*x)%11) for a in range(1,11))
        return tuple((a,b) for a,b in ops
                     if all((a*u+b)%11==v for u,v in pairs))
    return tuple(p for p in PERMS
                 if all(p[u]==v for u,v in pairs))

def apply_op(domain, op, value, inverse):
    if domain==0:
        a,b=op
        return (value-b)*pow(a,-1,11)%11 if inverse else (a*value+b)%11
    return op.index(value) if inverse else op[value]

def reduce_word(word):
    result=[]
    for role,inverse in word:
        if result and result[-1]==(role,not inverse):
            result.pop()
        else:result.append((role,inverse))
    return tuple(result)

def shared_role_certificate(domain, rows, word, start, budget=64):
    """Returns (answer|None, proof), allowing certification only when every
    globally consistent *shared binding* has been evaluated.

    No arithmetic confidence estimate is ever treated as a proof.
    When the integer assignment cap is exceeded, answer is UNKNOWN.
    """
    assert budget>=1 and 0<=start<(11 if domain==0 else 3)
    roles=reduce_word(word)
    if not roles:
        return start,dict(certified=True,assignments=0,
                          exact_shared_role_bindings=True,reason="IDENTITY")
    used=tuple(sorted({j for j,inv in roles}))
    options={j:consistent_ops(domain,rows[j]) for j in used}
    if any(not options[j] for j in used):
        return None,dict(certified=False,assignments=0,
                         exact_shared_role_bindings=True,reason="INCONSISTENT_PUBLIC_EVIDENCE")
    assignment_count=1
    for j in used:
        assignment_count*=len(options[j])
        if assignment_count>budget:
            return None,dict(certified=False,assignments=assignment_count,
                             exact_shared_role_bindings=True,reason="BOUNDED_UNKNOWN")
    # Lazy depth-first binding. A role identity is assigned once and reused.
    chosen={}
    answers=set()
    visited=0
    def search(i, value):
        nonlocal visited
        if i==len(roles):
            visited+=1
            answers.add(value)
            return
        role,inv=roles[i]
        if role in chosen:
            search(i+1,apply_op(domain,chosen[role],value,inv))
            return
        for op in options[role]:
            chosen[role]=op
            search(i+1,apply_op(domain,op,value,inv))
        del chosen[role]
    search(0,start)
    assert visited==assignment_count
    if len(answers)==1:
        return next(iter(answers)),dict(certified=True,assignments=visited,
                   exact_shared_role_bindings=True,reason="ALL_BINDINGS_AGREE")
    return None,dict(certified=False,assignments=visited,
                     exact_shared_role_bindings=True,reason="MULTIPLE_VALID_OUTCOMES")

def public_exact_fallback(domain, rows, word, start):
    assert all(len(rows[j])==2 for j,_ in reduce_word(word))
    answer,proof=shared_role_certificate(domain,rows,word,start,budget=1)
    if proof["reason"]=="BOUNDED_UNKNOWN":
        # Full public evidence yields exactly one operator per active role.
        answer,proof=shared_role_certificate(domain,rows,word,start,budget=1<<12)
    if not proof["certified"]:
        raise ValueError("PUBLIC_EVIDENCE_NOT_EXACT")
    return answer

def score_role_requests(training_cases):
    """TRAIN-only public ambiguity reduction, with no hidden task answers.
    A score is a query *priority*, not a soundness certificate.
    """
    scores=[0,0,0,0]
    for domain,rows,word,start in training_cases:
        active=tuple(sorted({j for j,inv in reduce_word(word)}))
        for role in active:
            before=len(consistent_ops(domain,rows[role][:1]))
            after=len(consistent_ops(domain,rows[role]))
            scores[role]+=(before-after)*sum(j==role for j,_ in reduce_word(word))
    return tuple(scores)

def bound_query_solve(domain,full_rows,word,start,weights,cap=64):
    used=sorted({j for j,_ in reduce_word(word)})
    rows=[r[:1] for r in full_rows]
    queries=0
    while True:
        val,proof=shared_role_certificate(domain,rows,word,start,budget=cap)
        if proof["certified"]:
            return val,dict(queries=queries,certificate=proof,
                            fallback=False,sound=True)
        remaining=[j for j in used if len(rows[j])==1]
        if not remaining:break
        priority=max(remaining,key=lambda j:(weights[j],-j))
        rows[priority]=full_rows[priority]
        queries+=1
    val=public_exact_fallback(domain,full_rows,word,start)
    return val,dict(queries=len(used),certificate=proof,
                    fallback=True,sound=True)
