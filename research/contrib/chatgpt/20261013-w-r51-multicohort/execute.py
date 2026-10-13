"""W-R51: prospective three-development-cohort reproducibility of W-R50 batch learning."""
import importlib.util
import json
from pathlib import Path

SOURCE=Path(__file__).resolve().parent.parent/"20261013-w-r50-batch"/"execute.py"
spec=importlib.util.spec_from_file_location("w50_pinned",SOURCE)
assert spec and spec.loader
w50=importlib.util.module_from_spec(spec)
spec.loader.exec_module(w50)

BLOCKS=(tuple(range(3216,3232)),tuple(range(3248,3264)),tuple(range(3280,3296)))
HELDOUT=tuple(range(3312,3360))
ARMS=w50.ARMS
GENERATORS=w50.GEN

def run():
    original_dev,original_heldout=w50.DEV,w50.HELDOUT
    studies=[]
    try:
        for i,dev in enumerate(BLOCKS):
            w50.DEV=dev
            w50.HELDOUT=HELDOUT
            result=w50.run()
            assert result["development_seeds"]==dev
            assert result["sealed_heldout_seeds"]==HELDOUT
            assert len(result["results"]["identity"])==3
            assert all(len(result["results"]["identity"][arm]["by_seed"])==48 for arm in ARMS)
            assert result["training_examples_per_arm"]==1728
            for arm in ARMS:
                assert result["results"]["identity"][arm]["accuracy"]==result["results"]["xor55"][arm]["accuracy"]
            ident=result["results"]["identity"]
            gain=ident["order_invariant_full_batch"]["accuracy"]-(ident["online_original"]["accuracy"]+ident["online_reverse"]["accuracy"])/2
            batch=ident["order_invariant_full_batch"]
            useful=(batch["accuracy"]>=.55 and gain>=.02)
            families=sum(batch["family_accuracy"][g]>.52 for g in GENERATORS)
            studies.append({
                "development_block":i+1,
                "development_seeds":dev,
                "training_examples_per_arm":1728,
                "training_corpus_sha256":result["original_training_corpus_sha256"],
                "batch_gain_vs_online_mean":gain,
                "batch_accuracy":batch["accuracy"],
                "batch_families_above_52pct":families,
                "meets_weak_per_block_gate":useful,
                "predecessor_rubric_label_this_cohort":result["classification"],
                "arms_by_encoding":result["results"]})
    finally:
        w50.DEV,w50.HELDOUT=original_dev,original_heldout
    assert len(studies)==3
    best=sum(s["meets_weak_per_block_gate"] for s in studies)
    strict=all(s["meets_weak_per_block_gate"] and s["batch_families_above_52pct"]>=2 for s in studies)
    label="BATCH_CROSS_TRAIN_ROBUST" if strict else (
        "PARTIAL_REPLICATION" if best>=2 else "OTHERWISE_NEGATIVE")
    return {"schema":"wingless.w-r51-three-development-cohort-batch.v1",
            "parent_head":"c97a5fcf97562f4320a43c27f409bd810c06fd6f",
            "classification":label,
            "prospective_development_blocks":BLOCKS,
            "sealed_heldout_seeds":HELDOUT,
            "models_and_training_corpora_chosen_before_first_use":True,
            "number_of_useful_batch_models":best,
            "minimum_batch_accuracy":min(s["batch_accuracy"] for s in studies),
            "studies":studies,"model_calls":0,
            "claim_limit":"Synthetic same-cohort evaluation of nine preregistered models trained on three disjoint development groups; heldout evidence is shared, not three independent populations."}

if __name__=="__main__":
    print(json.dumps(run(),sort_keys=True,indent=2))
