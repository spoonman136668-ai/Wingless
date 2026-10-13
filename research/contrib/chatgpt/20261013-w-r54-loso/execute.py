"""W-R54: preregistered leave-one-generator-out synthetic 16D ridge transfer."""
import hashlib
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r53-ridge"/"execute.py"
spec=importlib.util.spec_from_file_location("w53_pinned",SOURCE)
assert spec and spec.loader
w53=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w53)
w37=w53.w37
w42=w53.w42
w45=w53.w45

DEV=(tuple(range(5200,5216)),tuple(range(5232,5248)),tuple(range(5264,5280)))
HELDOUT=tuple(range(5312,5360))
GENERATORS=w53.GENERATORS
TRAIN_SITES=w53.w52.TRAIN_SITES
TEST_SITES=w53.TEST_SITES
ENCODINGS=("identity","xor55")

def build_examples(seeds,excluded):
    training=[]
    counts={g:0 for g in GENERATORS}
    for seed in seeds:
        for index,generator in enumerate(GENERATORS):
            if generator==excluded: continue
            for entity in range(9):
                # Exactly 4 or 5 entities, alternating across consecutive
                # seed parity; aggregate exactly 864 signed items per family.
                if entity%2!=(seed+index)%2:continue
                for site in TRAIN_SITES:
                    for amplitude in (1,2):
                        base,pos,neg=w37.recoded_example(
                            seed,entity,generator,"identity",site,amplitude,True)
                        p=w42.feature(base,pos)
                        n=w42.feature(base,neg)
                        delta=tuple(a-b for a,b in zip(p,n))
                        training.extend(((delta,1),(tuple(-x for x in delta),-1)))
                        counts[generator]+=2
    assert len(training)==1728
    assert counts[excluded]==0
    assert all(counts[g]==864 for g in GENERATORS if g!=excluded)
    corpus=hashlib.sha256(json.dumps([[list(v),y] for v,y in training],
                                  separators=(",",":")).encode()).hexdigest()
    return training,counts,corpus

def run():
    fitted={}
    for idx,block in enumerate(DEV,1):
        for excluded in GENERATORS:
            examples,counts,corpus=build_examples(block,excluded)
            weights=w53.ridge_fit(examples)
            assert len(weights)==16
            key=f"{idx}:withhold_{excluded}"
            fitted[key]={"weights":weights,"excluded":excluded,
                         "development_cohort":idx,
                         "source_counts":counts,"corpus_sha256":corpus,
                         "weight_sha256":hashlib.sha256(json.dumps(
                             weights,separators=(",",":")).encode()).hexdigest()}
    assert len(fitted)==9
    results={}
    for encoding in ENCODINGS:
        observed={}
        for key,model in fitted.items():
            by_source={g:{"credit":0.0,"n":0} for g in GENERATORS}
            by_seed=[]
            for seed in HELDOUT:
                record={"seed":seed,"correct_credit":0.0,
                        "presentations":0,"ties":0,"by_source_credit":{}}
                for generator in GENERATORS:
                    subtotal=0.0
                    for site in TEST_SITES:
                        for entity in range(9):
                            base,pos,neg=w37.recoded_example(
                                seed,entity,generator,encoding,site,3,False)
                            b0,p0,n0=w37.recoded_example(
                                seed,entity,generator,"identity",site,3,False)
                            a=w42.feature(base,pos)
                            b=w42.feature(base,neg)
                            assert a==w42.feature(b0,p0)
                            assert b==w42.feature(b0,n0)
                            margin=sum(w*(x-y) for w,x,y in zip(model["weights"],a,b))
                            credit=w45.credit(margin,True)+w45.credit(-margin,False)
                            subtotal+=credit
                            record["correct_credit"]+=credit
                            record["presentations"]+=2
                            record["ties"]+=2*int(abs(margin)<=1e-10)
                            by_source[generator]["credit"]+=credit
                            by_source[generator]["n"]+=2
                    record["by_source_credit"][generator]=subtotal
                assert record["presentations"]==108
                by_seed.append(record)
            accuracy={g:x["credit"]/x["n"] for g,x in by_source.items()}
            observed[key]={"withheld_source":model["excluded"],
                           "withheld_accuracy":accuracy[model["excluded"]],
                           "full_source_accuracy":accuracy,
                           "overall_accuracy":sum(x["correct_credit"] for x in by_seed)/(48*108),
                           "by_seed":by_seed,
                           "source_training_counts":model["source_counts"],
                           "training_corpus_sha256":model["corpus_sha256"],
                           "frozen_weight_sha256":model["weight_sha256"]}
        results[encoding]=observed
    for key in results["identity"]:
        assert results["identity"][key]["withheld_accuracy"]==results["xor55"][key]["withheld_accuracy"]
        assert results["identity"][key]["full_source_accuracy"]==results["xor55"][key]["full_source_accuracy"]
    comparisons=[]
    for idx in range(1,4):
        row={g:results["identity"][f"{idx}:withhold_{g}"]["withheld_accuracy"]
             for g in GENERATORS}
        comparisons.append({"cohort":idx,"withheld_source_accuracies":row,
                            "sources_above_55pct":sum(v>=.55 for v in row.values()),
                            "worst_withheld_source_accuracy":min(row.values())})
    strict=all(c["sources_above_55pct"]>=2 and c["worst_withheld_source_accuracy"]>=.48
               for c in comparisons)
    partial=any(all(results["identity"][f"{idx}:withhold_{g}"]["withheld_accuracy"]>=.53
                    for idx in (1,2,3)) for g in GENERATORS)
    label=("UNSEEN_SOURCE_GENERALIZATION" if strict
           else "PARTIAL_UNSEEN_SIGNAL" if partial else "NO_GAIN")
    return {"schema":"wingless.w-r54-leave-one-generator-out.v1",
            "parent_head":"d9da7cc2cf59d61d186be14a26ba20d532d01fe9",
            "classification":label,
            "development_blocks":DEV,"sealed_heldout_seeds":HELDOUT,
            "training_generators_per_model":2,
            "training_examples_per_model":1728,
            "ridge_lambda_normalized":1.0,
            "cohort_comparisons":comparisons,
            "all_models_by_encoding":results,"new_training_from_heldout":False,
            "model_calls":0,
            "claim_limit":"Two-generator-trained ridge models scored on unseen synthetic third generator, fixed two-candidate interface, not real-world source-born general cognition."}

if __name__=="__main__":
    print(json.dumps(run(),indent=2,sort_keys=True))
