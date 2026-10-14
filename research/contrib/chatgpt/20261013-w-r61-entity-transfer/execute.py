"""W-R61: frozen held-out synthetic entity-identity transfer of source-withheld ridge."""
import importlib.util
import json
from pathlib import Path

PARENT=Path(__file__).resolve().parent.parent/"20261013-w-r59-replication"/"execute.py"
spec=importlib.util.spec_from_file_location("w59_pinned",PARENT)
assert spec and spec.loader
m=importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)

DEV=(tuple(range(9008,9024)),tuple(range(9040,9056)),tuple(range(9072,9088)))
HELDOUT=tuple(range(9120,9168))
AMPS=(1,2,3)
GENERATORS=m.m.GENERATORS
GROUPS={"seen_identity_control":tuple(range(9)),"unseen_identity_transfer":tuple(range(9,18))}

def fit():
    result={}
    for i,seeds in enumerate(DEV,1):
        for excluded in GENERATORS:
            pairs,counts,fingerprint=m.m.build_examples(seeds,excluded)
            assert counts[excluded]==0 and len(pairs)==1728
            weights=m.m.w53.ridge_fit(pairs)
            assert len(weights)==16
            result[f"{i}:withhold_{excluded}"]={
                "weights":weights,"training_counts":counts,"development_sha256":fingerprint}
    assert len(result)==9
    return result

def run():
    trained=fit()
    records=[]
    scores={}
    for name,model in trained.items():
        scores[name]={}
        for amplitude in AMPS:
            scored={}
            for generator in GENERATORS:
                measures={}
                for group,entities in GROUPS.items():
                    by_seed=[]
                    credit=0.0
                    count=0
                    for seed in HELDOUT:
                        subtotal=0.0
                        n=0
                        for site in m.m.TEST_SITES:
                            for entity in entities:
                                base,pos,neg=m.m.w37.recoded_example(
                                    seed,entity,generator,"identity",site,amplitude,False)
                                xb,xp,xn=m.m.w37.recoded_example(
                                    seed,entity,generator,"xor55",site,amplitude,False)
                                positive=m.m.w42.feature(base,pos)
                                negative=m.m.w42.feature(base,neg)
                                assert positive==m.m.w42.feature(xb,xp)
                                assert negative==m.m.w42.feature(xb,xn)
                                margin=sum(w*(a-b) for w,a,b in zip(
                                    model["weights"],positive,negative))
                                subtotal+=m.m.w45.credit(margin,True)+m.m.w45.credit(-margin,False)
                                n+=2
                        by_seed.append({"seed":seed,"credit":subtotal,"presentations":n})
                        credit+=subtotal
                        count+=n
                    measures[group]={"accuracy":credit/count,"correct_credit":credit,
                                     "presentations":count,"by_seed":by_seed}
                scored[generator]=measures
                records.append({"model":name,"amplitude":amplitude,"generator":generator,
                                "group_scores":{k:v["accuracy"] for k,v in measures.items()}})
            scores[name][str(amplitude)]=scored
    comparisons=[]
    for i in range(1,4):
        key=f"{i}:withhold_xorshift_independent"
        for amp in AMPS:
            entry=scores[key][str(amp)]["xorshift_independent"]
            control=entry["seen_identity_control"]["accuracy"]
            transfer=entry["unseen_identity_transfer"]["accuracy"]
            comparisons.append({"model":key,"amplitude":amp,
                                "seen_accuracy":control,"unseen_accuracy":transfer,
                                "transfer_minus_seen":transfer-control})
    strict=all(r["unseen_accuracy"]>=.55 and r["transfer_minus_seen"]>=-.05 for r in comparisons)
    partial=(sum(r["unseen_accuracy"]>=.55 for r in comparisons)>=5 and
             sum(r["transfer_minus_seen"] for r in comparisons)/len(comparisons)>=-.10)
    label="ENTITY_GENERALIZATION" if strict else (
        "PARTIAL_ENTITY_TRANSFER" if partial else "NO_ENTITY_TRANSFER")
    return {"schema":"wingless.w-r61-entity-heldout.v1",
            "parent_head":"867ca5dea0d316b4d89d8787585369a0bad508cb",
            "classification":label,"development_blocks":DEV,
            "sealed_heldout_seeds":HELDOUT,
            "training_entities":list(GROUPS["seen_identity_control"]),
            "unseen_entities":list(GROUPS["unseen_identity_transfer"]),
            "amplitudes":AMPS,"xorshift_comparisons":comparisons,
            "full_cross_family":scores,"compact_audit_rows":records,
            "training_provenance":{k:{"source_counts":v["training_counts"],
                                     "sha256":v["development_sha256"]} for k,v in trained.items()},
            "model_calls":0,"accepted_ref_mutation":False,
            "claim_limit":"Synthetic two-candidate feature-level heldout identity transfer only."}

if __name__=="__main__":
    print(json.dumps(run(),indent=2,sort_keys=True))
