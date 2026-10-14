"""Archive-only independent W-R64 recomputation; does NOT replay first-use."""
import hashlib
import json
import math
from pathlib import Path
p=Path("evidence/W-R64-original.json")
raw=p.read_bytes()
z=json.loads(raw)
assert z["schema"]=="wingless.w-r64-source-channel-specificity.v1"
assert z["parent_head"]=="635d8409a175c483dde48b624149ba1f6505d1a7"
assert z["sealed_heldout_seeds"]==list(range(9664,9688))
assert z["model_calls"]==0 and z["accepted_ref_mutation"] is False
rows=z["all_sealed_records"]
assert len(rows)==180
seen=set()
scores={}
seedrows=0
for r in rows:
    key=(r["model"],r["amplitude"],r["permutation"],r["condition"])
    assert key not in seen
    seen.add(key)
    per=r["by_seed"]
    assert [s["seed"] for s in per]==list(range(9664,9688))
    assert len(per)==24
    seedrows+=len(per)
    credit=sum(s["correct_credit"] for s in per)
    denom=sum(s["presentations"] for s in per)
    assert denom>0
    assert math.isclose(credit,r["credit"],abs_tol=1e-9)
    assert denom==r["presentations"]
    assert math.isclose(credit/denom,r["accuracy"],abs_tol=1e-12)
    scores[key]=credit/denom
assert len(scores)==3*3*4*5 and seedrows==4320
independent=[]
for model in (f"{i}:withhold_xorshift_independent" for i in (1,2,3)):
    for amp in (1,2,3):
        for perm in ("identity","rotate","swap01","reverse"):
            key=(model,amp,perm)
            intact=scores[key+("intact",)]
            delta={str(channel):scores[key+(f"erase_source{channel}",)]-intact for channel in range(4)}
            independent.append({"model":model,"amplitude":amp,"permutation":perm,
                                "intact_accuracy":intact,"loss_vs_intact":delta})
assert len(independent)==36
for recorded,derived in zip(z["comparisons"],independent):
    assert recorded["model"]==derived["model"]
    assert recorded["amplitude"]==derived["amplitude"]
    assert recorded["permutation"]==derived["permutation"]
    assert math.isclose(recorded["intact_accuracy"],derived["intact_accuracy"],abs_tol=1e-12)
    for c in ("0","1","2","3"):
        assert math.isclose(recorded["loss_vs_intact"][c],derived["loss_vs_intact"][c],abs_tol=1e-12)
mean={str(c):sum(x["loss_vs_intact"][str(c)] for x in independent)/36 for c in range(4)}
for key,derived in mean.items():
    assert math.isclose(z["source_erase_mean_loss"][key],derived,abs_tol=1e-12)
if mean["0"]<=-.07 and mean["0"]<=min(mean[str(c)] for c in (1,2,3))-.035:
    label="SOURCE0_DOMINANT_REPLICATED"
elif sum(v<=-.05 for v in mean.values())>=2:
    label="MULTICHANNEL_SENSITIVITY"
else:
    label="NO_CLEAR_SOURCE_RANK"
assert label==z["classification"]
witness={"schema":"research.archived-evidence-witness.v1","experiment":"W-R64",
         "classification":label,"source_run":38530053003,"source_artifact_id":11878040664,
         "original_json_sha256":hashlib.sha256(raw).hexdigest(),
         "independent_case_rows":len(rows),"independent_seed_rows":seedrows,
         "source_channel_mean_loss":mean,"sealed_execution_replayed":False,
         "scope":"Read-only independent archive recomputation; not external FSA authority"}
Path("W-R64-witness.json").write_text(json.dumps(witness,sort_keys=True,indent=2)+"\n")
print("WITNESS_PASS",json.dumps(witness,sort_keys=True))
