"""W-R27: known ADD4 conjugacy hypotheses confronted with latent XOR4."""
import hashlib
import importlib.util
import json
from pathlib import Path

PREVIOUS=Path(__file__).resolve().parent.parent/"20261013-w-r26-latent-alphabet"/"execute.py"
spec=importlib.util.spec_from_file_location("w26_frozen",PREVIOUS)
assert spec and spec.loader
w26=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w26)
PERMUTATIONS=w26.PERMUTATIONS
BUDGETS=w26.BUDGETS


def operation(seed,a,b):
    return (a+b)%4 if seed%2==0 else a^b


def sequence(seed,episode):
    assert 120<=seed<=127 and 0<=episode<16
    values=[(episode+seed)%4,((episode//4)+3*seed)%4]
    while len(values)<32:
        values.append(operation(seed,values[-2],values[-1]))
    digest=hashlib.sha256(("W27/"+str(seed)).encode()).digest()
    pi=PERMUTATIONS[digest[0]%24]
    return tuple(pi[x] for x in values)


def triples(seed,k):
    cases=[]
    for episode in range(k):
        obs=sequence(seed,episode)
        cases.append((obs[0],obs[1],obs[2]))
    assert len(cases)==k
    assert len({(a,b) for a,b,_ in cases})==k
    return cases


def evaluate(seed,table):
    cm=[[0]*4 for _ in range(4)]
    digest=hashlib.sha256()
    for episode in range(16):
        rows=sequence(seed,episode)
        for tick in range(3,32):
            a,b=rows[tick-2],rows[tick-1]
            guess=(a+b)%4 if table is None else table[4*a+b]
            actual=rows[tick]
            cm[actual][guess]+=1
            digest.update(bytes((episode,tick,actual,guess)))
    assert sum(map(sum,cm))==464
    correct=sum(cm[i][i] for i in range(4))
    return {"correct":correct,"total":464,"accuracy":correct/464,
            "confusion":cm,"trace_sha256":digest.hexdigest()}


def run():
    assert w26.fixture()=="PASS"  # artificial permutation fixture only
    rows=[]
    for seed in range(120,128):
        category="ADD4" if seed%2==0 else "XOR4"
        rows.append({
            "seed":seed,"family":category,"budget":0,
            "candidate_matches":0,"distinct_tables":1,"accepted":True,
            "test":evaluate(seed,None),
        })
        for budget in BUDGETS:
            evidence=triples(seed,budget)
            candidates=[p for p in PERMUTATIONS if w26.fits(p,evidence)]
            mapping={w26.table_for(p) for p in candidates}
            accepted=len(mapping)==1
            table=next(iter(mapping)) if accepted else None
            rows.append({
                "seed":seed,"family":category,"budget":budget,
                "candidate_matches":len(candidates),
                "distinct_tables":len(mapping),
                "accepted":accepted,
                "table_sha256":hashlib.sha256(bytes(table)).hexdigest()
                               if accepted else None,
                "test":evaluate(seed,table) if accepted else None,
            })
    cohorts={}
    for budget in BUDGETS:
        for family in ("ADD4","XOR4"):
            records=[r for r in rows if r["budget"]==budget and r["family"]==family]
            cohort=family+"_k"+str(budget)
            count=sum(r["accepted"] for r in records)
            hits=sum(r["test"]["correct"] for r in records if r["accepted"])
            cohorts[cohort]={"accepted_seeds":count,"conditional_accuracy":hits/(count*464)
                             if count else None}
    in_good=cohorts["ADD4_k16"]["accepted_seeds"]==4 and cohorts["ADD4_k16"]["conditional_accuracy"]==1
    out_reject=cohorts["XOR4_k16"]["accepted_seeds"]==0
    label=("OPERATION_FAMILY_BOUNDARY" if in_good and out_reject else
           "FALSE_ACCEPT" if not out_reject else "MIXED_OR_NEGATIVE")
    return {
        "schema":"wingless.w-r27-xor-out-of-family.v1",
        "predecessor_head":"1dfad834724594354aa622d23fd0405f76e79a07",
        "classification":label,"candidate_permutations":24,
        "cohort_scores":cohorts,
        "rows":rows,"model_calls":0,
        "claim_limit":"Synthetic ADD4 versus nonisomorphic XOR4 operations, fixed known hypothesis class."
    }


if __name__=="__main__":
    outcome=run()
    assert outcome==run()
    print(json.dumps(outcome,sort_keys=True,indent=2))
