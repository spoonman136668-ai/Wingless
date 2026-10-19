"""T21-B source-free finite S3 algebra qualification.

No TRAIN, DEV, independently source-born tasks, or consumed original cohorts
are imported. Everything below is synthetic math property verification.
"""
import importlib.util
import itertools
import json
import random
import sys
from pathlib import Path

ROOT=Path(__file__).resolve().parent

def load(name, file):
    spec=importlib.util.spec_from_file_location(name,ROOT/file)
    module=importlib.util.module_from_spec(spec)
    sys.modules[name]=module
    spec.loader.exec_module(module)
    return module

def main():
    basis=load("t21_basis","20261019-wingless-t21-relational-bitmatrix.py")
    cert=load("t21_exact_existing","20261019-wingless-t16b-persistent-bindings.py")
    rng=random.Random(210619)
    words=[(),((0,False),),((0,True),(0,False)),
        ((0,False),(1,False),(0,True),(1,True)),
        ((0,False),(1,True),(2,False),(0,True),(2,True)),
        ((0,True),(1,False),(0,False),(2,True),(1,True),(2,False))]
    alphabet=tuple(itertools.product(range(3),(False,True)))
    for k in range(144):
        size=3+(k%12)
        words.append(tuple(rng.choice(alphabet) for _ in range(size)))
    checks={"universal_certifications":0,"composed_relations":0,
            "oracle_public_agreement":0,"alpha_renaming":0,
            "tamper_rejection":0,"source_cases_opened":0}
    for word in words:
        full=basis.construct(word)
        proof=basis.independently_verify(full,cert)
        assert proof["verified"] and proof["full_assignment_comparisons"]==648
        checks["universal_certifications"]+=1
        if len(word)>=2:
            split=1+(len(word)-1)//2
            composed=basis.compose(basis.construct(word[:split]),
                                   basis.construct(word[split:]))
            assert composed.relation==full.relation
            assert composed.identity==full.identity
            basis.independently_verify(composed,cert)
            checks["composed_relations"]+=1
        for sample in (0,4,23,102,215):
            inds=basis.ASSIGNMENTS[sample]
            rows_full=tuple(((0,basis.OPS[inds[role]][0]),
                             (1,basis.OPS[inds[role]][1])) for role in range(3))
            rows_first=tuple((row[0],) for row in rows_full)
            for rows in (rows_first,rows_full):
                for start in range(3):
                    got,detail=basis.infer(full,rows,start)
                    ref,reference=cert.shared_role_certificate(1,rows,word,start,
                                                               budget=4096)
                    assert bool(detail["certified"])==bool(reference["certified"])
                    assert got==ref
                    checks["oracle_public_agreement"]+=1
        if word:
            rename=(2,0,1)
            new=tuple((rename[r],inv) for r,inv in word)
            variant=basis.construct(new)
            for case in (0,17,58,215):
                original=basis.ASSIGNMENTS[case]
                renamed=[None]*3
                for role in range(3):
                    renamed[rename[role]]=original[role]
                for start in range(3):
                    assert basis.direct(word,start,original)==basis.direct(
                        new,start,renamed)
                    checks["alpha_renaming"]+=1
        # Recompute a plausible but false digest: independent algebra must reject.
        relation=[list(row) for row in full.relation]
        relation[0][0],relation[0][1]=relation[0][1],relation[0][0]
        forged=tuple(tuple(row) for row in relation)
        fake=basis.Summary(full.word,forged,basis.digest(full.word,forged))
        try:
            basis.independently_verify(fake,cert)
        except ValueError as e:
            assert str(e)=="T21_UNIVERSAL_RELATION_NOT_EQUIVALENT"
            checks["tamper_rejection"]+=1
        else:
            raise AssertionError("T21_FALSE_RECOMPUTED_DIGEST_ACCEPTED")
    assert checks["universal_certifications"]==len(words)
    assert checks["source_cases_opened"]==0
    return {"schema":"wingless.t21.algebraic-foundation-only.v1",
        "result":"PASS_ALGEBRA_NOT_SCIENCE", "checks":checks,
        "independent_source_family_admitted":False,
        "learned_model_trained":False,"scientific_acceptance":False,
        "optimizer_admission":False,"T16_originals_consumed":0,
        "T19_T20_DEV_reused":False}

if __name__=="__main__":
    print(json.dumps(main(),sort_keys=True))
