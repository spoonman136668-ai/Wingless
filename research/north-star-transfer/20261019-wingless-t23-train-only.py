"""T23 TRAIN-only equality and learned-prefix capacity proof; DEV never built."""
import importlib.util,json,sys
from pathlib import Path
ROOT=Path(__file__).parent
def load(n,p):
    spec=importlib.util.spec_from_file_location(n,ROOT/p)
    m=importlib.util.module_from_spec(spec);sys.modules[n]=m;spec.loader.exec_module(m)
    return m
def main():
    model=load("t23_training_model","20261019-wingless-t23-prefix-dag-learner.py")
    cert,source=model.deps()
    assert len(source.TRAIN)==40 and len(source.DEV)==64 and not set(source.TRAIN)&set(source.DEV)
    train=[source.train_public(source.build(x)) for x in source.TRAIN]
    rank=cert.score_role_requests([(1,c["role_pairs"],word,start)
              for c in train for word,start in zip(c["words"],c["starts"])])
    pools={a:model.train_cache(train,a) for a in model.ARMS[:2]}
    assert all(len(nodes)<=32 and updates<=2048 for nodes,updates in pools.values())
    checks=0
    for case in train:
        for variation in (0,1):
            def requested():
                pairs=case["role_pairs"];first=tuple((r[0],) for r in pairs)
                charged=set()
                def request(j):
                    if j in charged:raise ValueError("T23_TRAIN_DOUBLE_QUERY")
                    charged.add(j)
                    return pairs[j][1]
                return first,request
            reference_first,request=requested()
            ref,ref_n,*_=model.exact_batch(cert,reference_first,request,
                 case["words"],case["starts"],rank,
                 "fully-trained-T16B-batched",frozenset())
            if variation==1:
                # Force every public role to be fully known before verifying.
                role_rows=case["role_pairs"]
                ref=tuple(cert.public_exact_fallback(1,role_rows,w,s)
                         for w,s in zip(case["words"],case["starts"]))
            for arm in model.ARMS[:3]:
                first,request=requested()
                val,num,*_=model.exact_batch(cert,first,request,case["words"],
                    case["starts"],rank,arm,pools.get(arm,(frozenset(),0))[0])
                if val!=ref or (variation==0 and num!=ref_n):
                    raise RuntimeError("T23_TRAIN_VARIABLE_IDENTITY_UNSOUND")
                checks+=1
    return {"schema":"wingless.t23.TRAIN-only-admission.v1","status":"PASS",
        "TRAIN_cases":40,"DEV_generated":0,"private_oracle_not_used":True,
        "equality_checks":checks,"learned_prefix_nodes":len(pools[model.ARMS[0]][0]),
        "frequency_prefix_nodes":len(pools[model.ARMS[1]][0]),
        "TRAIN_updates":pools[model.ARMS[0]][1],
        "external_independence_claim":False,"scientific_admission":False,
        "optimizer_admission":False}
if __name__=="__main__":print(json.dumps(main(),sort_keys=True))
