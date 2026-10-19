"""T17 internal synthetic development source; NEVER a protected T16 original.

Distinct research scratch source authored by same model. TRAIN and DEV are
irreversibly separate deterministic toy fixtures, not independently born
scientific sources; zero sealed originals accessed. Hidden operator labels
live inside fixture construction and the scorer only.
"""
import random

TRAIN=tuple(f"t17-synth-train-{i:03d}" for i in range(24))
DEV=tuple(f"t17-synth-dev-{i:03d}" for i in range(64))
PATTERNS=(
 ((0,False),(1,True),(0,True),(1,False)),   # commutator
 ((0,False),(1,False),(0,True)),           # conjugacy
 ((0,False),(1,True),(0,False),(1,False)), # shared repeated role
 ((0,False),(1,False),(1,True),(0,True)),  # identity through adjacent cancel
 ((0,False),(1,False),(0,False),(1,True),(0,True)),
 ((0,True),(1,True),(0,False),(1,False)),
 ((0,False),(1,False),(0,True),(1,True),(0,False)),
 ((0,False),(0,False),(1,True),(1,False))
)
def birth_internal_fixture(label,k):
    assert label in TRAIN or label in DEV, "T17_NOT_INTERNAL_DIAGNOSTIC"
    idx=int(label.rsplit("-",1)[-1])
    is_train=label in TRAIN
    rng=random.Random(170703131+idx*30491+(0 if is_train else 10007899))
    domain=(idx+(0 if is_train else 1))%2
    n=11 if domain==0 else 3
    roles=[]
    hidden=[]
    for role in range(2):
        op=rng.choice(k.candidates(domain))
        left,right=rng.sample(range(n),2)
        roles.append(((left,k.operator_apply(domain,op,left)),
                      (right,k.operator_apply(domain,op,right))))
        hidden.append(op)
    # Hold grammar fixed, but swap pattern proportions/train ordering.
    pattern_idx=(idx*3+(0 if is_train else 5))%len(PATTERNS)
    word=PATTERNS[pattern_idx]
    start=rng.randrange(n)
    state=k.start(domain)
    for role,inverse in word:
        state=k.compose_state(state,hidden[role],domain,inverse=inverse)
    true_answer=k.apply(state,start)
    return {"label":label,"domain":domain,"roles":tuple(roles),
            "word":word,"start":start,"truth":true_answer,
            "private_hidden_operator_never_passed_to_candidate":True}
