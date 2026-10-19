"""T16-E external-authored S3 multiplication semantic reference, NONCONSUMING.

Independent upstream publication dashstander/sn-grok, commit
cb0673dba621a92b1ae4566cb2baf449b7fe6e06, sngrok/permutations.py
Git blob aa2489d8e528bad04ba60af36cf12fff11a52dd7, MIT licensed.
Upstream Permutation.__mul__ uses right.sigma[left.sigma[i]].

This source-neutral public algebra check is NOT independently originated
held-out variable-binding tasks, nor independent first-use or science.
"""
from itertools import permutations

UPSTREAM_REPO="dashstander/sn-grok"
UPSTREAM_COMMIT="cb0673dba621a92b1ae4566cb2baf449b7fe6e06"
UPSTREAM_PERMUTATION_BLOB="aa2489d8e528bad04ba60af36cf12fff11a52dd7"
UPSTREAM_LICENSE_BLOB="fb58af050cc5bf48009f3478ca02f7200c0b9bf1"
DOMAIN=tuple(permutations(range(3)))

def upstream_s3_compose(left,right):
    """Structural transcription of upstream Permutation.__mul__,
    not a representation of any unopened original science source."""
    assert type(left) is tuple and type(right) is tuple
    assert left in DOMAIN and right in DOMAIN
    return tuple(right[left[i]] for i in range(3))

def inverse(p):
    assert p in DOMAIN
    return tuple(p.index(i) for i in range(3))

def check_wingless_public_algebra(k):
    assert len(DOMAIN)==6
    proof=0
    for left in DOMAIN:
        for right in DOMAIN:
            state=k.compose_state(k.start(k.NONCOMMUTING_PERM),left,k.NONCOMMUTING_PERM)
            state=k.compose_state(state,right,k.NONCOMMUTING_PERM)
            expected=upstream_s3_compose(left,right)
            assert tuple(k.apply(state,i) for i in range(3))==expected
            for i in range(3):
                assert k.operator_apply(k.NONCOMMUTING_PERM,expected,i)==right[left[i]]
            proof+=1
    return {"source":"EXTERNAL_AUTHOR_S3_PUBLIC_ALGEBRA_ONLY",
       "upstream_repo":UPSTREAM_REPO,"upstream_commit":UPSTREAM_COMMIT,
       "upstream_blob":UPSTREAM_PERMUTATION_BLOB,"exhaustive_products":proof,
       "exact_semantics_preserved":True,"same_as_previously_known_S3":True,
       "novel_family_discovered":False,"original_firstuse_performed":False,
       "independent_witness_verified":False,"learned_transfer_improvement":False,
       "full_resource_parity_measured":False,"optimizer_admission":False}
