"""T2 independent compositional source definition and representation audit.

Only negative-numbered synthetic fixture IDs are callable. Registered source
TRAIN/VAL/UNSEEN/REPLICATION namespaces are deliberately inaccessible.
S3 heldout is specified by a fixed mathematical group presentation only;
NO unseen task, labels, or seed are constructed on this path.
"""
from __future__ import annotations
import random
from pathlib import Path
import json,hashlib

HERE=Path(__file__).resolve().parent
FREEZE=HERE/"20261018-wingless-t2-source-freeze.json"
FIXTURE_IDS=(-900,-901)
UNSEEN_IDS=range(882000,882128)
REPLICA_IDS=range(883000,883128)

def need(ok,tag):
    if not ok:raise ValueError("WNS_T2_"+tag)

def build_fixture(seed,algebra="C7",length=10):
    need(seed in FIXTURE_IDS and algebra in ("C7","C5") and length in (10,12),
         "REGISTERED_COHORTS_SEALED")
    prime=7 if algebra=="C7" else 5
    rng=random.Random(97531+seed*31+prime*197+length)
    roles={}
    for role in range(4):
        a=rng.randrange(1,prime);b=rng.randrange(prime)
        xs=tuple(rng.sample(range(prime),3))
        obs=tuple((x,(a*x+b)%prime) for x in xs)
        roles[role]=obs
    chain=tuple((rng.randrange(4),rng.choice((1,-1))) for _ in range(length))
    start=rng.randrange(prime)
    # Evaluation-only calculator, never exposed to a learner.
    value=start
    for role,direction in chain:
        x,y=roles[role][:2]
        z,t=roles[role][1:]
        a=((t-y)*pow((z-x)%prime,-1,prime))%prime
        b=(y-a*x)%prime
        if direction==1:value=(a*value+b)%prime
        else:value=(value-b)*pow(a,-1,prime)%prime
    public={"prime":prime,"public_roles":roles,"public_chain":chain,"start":start}
    evaluator={"expected_answer":value,"source_birth":"FIXTURE_ONLY","sealed":False}
    need(len(chain)==length and all(len(v)==3 for v in roles.values()),
         "PUBLIC_CALIBRATION_SHAPE")
    return public,evaluator

def affine_compose(x,y,p=7):
    a,b=x;c,d=y
    return ((a*c)%p,(a*d+b)%p)
def affine_inverse(x,p=7):
    a,b=x;ai=pow(a,-1,p)
    return (ai,(-ai*b)%p)
def affine_order(x,p=7):
    identity=(1,0);v=identity
    for n in range(1,44):
        v=affine_compose(v,x,p)
        if v==identity:return n
    raise AssertionError("AFFINE_ORDER")

def s3_nonrepresentability():
    # AGL(1,7) has no nontrivial S3 embedding: prove via all 42
    # group elements, requiring s^3=e, t^2=e, t*s*t^-1=s^-1.
    group=[(a,b) for a in range(1,7) for b in range(7)]
    identity=(1,0)
    candidates=0
    for s in group:
        if affine_order(s)!=3:continue
        for t in group:
            if affine_order(t)!=2:continue
            candidates+=1
            lhs=affine_compose(affine_compose(t,s),affine_inverse(t))
            if lhs==affine_inverse(s):
                return {"representable":True,"candidate_pairs_checked":candidates}
    return {"representable":False,"candidate_pairs_checked":candidates,
            "heldout_tasks_generated":0,
            "affine_mod7_cannot_hold_noncommuting_S3":True}

def attest():
    d=json.loads(FREEZE.read_text())
    need(d["schema"]=="northstar.t2.wingless.unseen-operator-source.v1"
         and d["sealed_namespaces_never_instantiated_in_T2"] is True
         and d["source_anchor_head"]=="f6de5cb4fd0d25b18974a14db4e81cdfe233e88a",
         "FREEZE_SOURCE_IDENTITY")
    need(not set(FIXTURE_IDS).intersection(UNSEEN_IDS)
         and not set(FIXTURE_IDS).intersection(REPLICA_IDS),"COHORT_DISJOINT")
    for id in (880000,881000,882000,883000):
        try:build_fixture(id)
        except ValueError:pass
        else:raise AssertionError("registered source instantiated")
    a=build_fixture(-900,"C7",10)
    b=build_fixture(-901,"C5",12)
    need(len(a[0]["public_chain"])==10 and len(b[0]["public_chain"])==12,
         "INDEPENDENT_SYNTHTEST_STRUCTURE")
    proof=s3_nonrepresentability()
    need(proof["representable"] is False,"OLD_LEARNER_CANNOT_TRANSFER")
    return {"schema":"wingless.t2.nonconsuming-readiness.v1",
            "TRAIN_instantitated":0,"UNSEEN_instantiated":0,
            "REPICA_instantiated":0,"algebra_source_incompatible":True,
            "S3_necessary_adaptation":True,
            "affine_group_proof":proof,
            "whole_system_cost_qualified":False,
            "native_transfer_ready":False,
            "scientific_first_use":False}
if __name__=="__main__":
    print(json.dumps(attest(),sort_keys=True))
