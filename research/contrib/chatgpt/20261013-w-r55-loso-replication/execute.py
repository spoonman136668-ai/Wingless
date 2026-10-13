"""W-R55 frozen prospective new-cohort reuse of W-R54 source-holdout experiment."""
import importlib.util,json
from pathlib import Path
p=Path(__file__).resolve().parent.parent/"20261013-w-r54-loso"/"execute.py"
spec=importlib.util.spec_from_file_location("w54_custody",p)
assert spec and spec.loader
w54=importlib.util.module_from_spec(spec);spec.loader.exec_module(w54)
DEV=(tuple(range(5408,5424)),tuple(range(5440,5456)),tuple(range(5472,5488)))
HELDOUT=tuple(range(5520,5568))
def run():
    previous=(w54.DEV,w54.HELDOUT)
    try:
        w54.DEV=DEV;w54.HELDOUT=HELDOUT
        result=w54.run()
    finally:
        w54.DEV,w54.HELDOUT=previous
    assert result["development_blocks"]==DEV
    assert result["sealed_heldout_seeds"]==HELDOUT
    family="xorshift_independent"
    vals=[result["all_models_by_encoding"]["identity"][f"{i}:withhold_{family}"]["withheld_accuracy"] for i in (1,2,3)]
    broad=all(x["sources_above_55pct"]>=2 and x["worst_withheld_source_accuracy"]>=.48
              for x in result["cohort_comparisons"])
    if min(vals)>=.65 and broad:label="STABLE_XORSHIFT_UNSEEN_TRANSFER"
    elif min(vals)>=.55:label="PARTIAL_UNSEEN_REPLICATION"
    else:label="NOT_REPLICATED"
    return {"schema":"wingless.w-r55-source-heldout-replication.v1",
            "parent_head":"2e1fc9d5ce47c7c469c3f5f7fade8313aab0d5a9",
            "classification":label,"source_specific_xorshift_accuracy":vals,
            "broad_source_transfer_gate":broad,
            "development_blocks":DEV,"sealed_heldout_seeds":HELDOUT,
            "w54_original_rubric_on_new_cohort":result["classification"],
            "original_frozen_full_result":result,"model_calls":0,
            "claim_limit":"Synthetic two-candidate transfer; new cohort independent but unchanged underlying source generator families."}
if __name__=="__main__":print(json.dumps(run(),sort_keys=True,indent=2))
