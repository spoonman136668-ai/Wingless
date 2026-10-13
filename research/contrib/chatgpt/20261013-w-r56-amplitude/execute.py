"""W-R56 frozen cross-amplitude source-heldout synthetic ridge challenge."""
import importlib.util,json
from pathlib import Path
p=Path(__file__).resolve().parent.parent/"20261013-w-r54-loso"/"execute.py"
spec=importlib.util.spec_from_file_location("w54_pinned",p)
assert spec and spec.loader
m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
DEV=(tuple(range(5700,5716)),tuple(range(5732,5748)),tuple(range(5764,5780)))
HELDOUT=tuple(range(5824,5872))
AMPS=(1,2,4)
def run():
    models={}
    for i,dev in enumerate(DEV,1):
        for source in m.GENERATORS:
            x,counts,fingerprint=m.build_examples(dev,source)
            weights=m.w53.ridge_fit(x)
            models[f"{i}:{source}"]={"weights":weights,"training_counts":counts,
                                     "dev_corpus_sha256":fingerprint}
    results={}
    for amp in AMPS:
        by_model={}
        for name,model in models.items():
            families={g:{"credit":0.0,"presentations":0} for g in m.GENERATORS}
            seed_rows=[]
            for seed in HELDOUT:
                row={"seed":seed,"by_source":{}}
                for generator in m.GENERATORS:
                    credit=0.0
                    for site in m.TEST_SITES:
                        for entity in range(9):
                            base,pos,neg=m.w37.recoded_example(seed,entity,generator,"identity",site,amp,False)
                            b2,p2,n2=m.w37.recoded_example(seed,entity,generator,"xor55",site,amp,False)
                            a=m.w42.feature(base,pos);b=m.w42.feature(base,neg)
                            assert a==m.w42.feature(b2,p2)
                            assert b==m.w42.feature(b2,n2)
                            margin=sum(w*(v-u) for w,v,u in zip(model["weights"],a,b))
                            credit+=m.w45.credit(margin,True)+m.w45.credit(-margin,False)
                            families[generator]["presentations"]+=2
                    row["by_source"][generator]=credit
                    families[generator]["credit"]+=credit
                seed_rows.append(row)
            by_model[name]={"family_accuracy":{g:v["credit"]/v["presentations"]
                                              for g,v in families.items()},
                            "by_seed":seed_rows,
                            "development_source_counts":model["training_counts"],
                            "training_sha256":model["dev_corpus_sha256"]}
        results[str(amp)]=by_model
    compared={}
    for amp in AMPS:
        compared[str(amp)]={g:[results[str(amp)][f"{i}:{g}"]["family_accuracy"][g]
                                for i in (1,2,3)] for g in m.GENERATORS}
    xor=all(a>=.55 for amp in compared.values()
            for a in amp["xorshift_independent"])
    strict=(all(a>=.65 for amp in compared.values()
                for a in amp["xorshift_independent"]) and
            all(sum(amp[g][i]>=.55 for g in m.GENERATORS)>=2
                for amp in compared.values() for i in range(3)))
    classification="AMPLITUDE_ROBUST_XORSHIFT" if strict else (
        "PARTIAL_AMPLITUDE_TRANSFER" if xor else "OTHERWISE_NEGATIVE")
    return {"schema":"wingless.w-r56-swap-amplitude-transfer.v1",
            "parent_head":"7656df2467df5013cdc48148804d47ab6f9cbaa3",
            "classification":classification,
            "development_blocks":DEV,"sealed_heldout_seeds":HELDOUT,
            "test_amplitudes":AMPS,"comparison":compared,
            "all_original_sealed_records":results,"model_calls":0,
            "claim_limit":"Synthetic source-heldout two-candidate ridge transfer under new swap amplitudes."}
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
