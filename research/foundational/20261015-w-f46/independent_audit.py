"""Independent read-only audit of the existing F46 Wingless TRAIN artifact.

Never invokes source_pilot.run, never buys actual study mappings, does not
generate a fresh diagnostic or original. Source constructor is used solely
to verify deterministic evidence provenance and receipt integrity.
"""
import importlib.util
import json
import math
import pathlib
import sys

ROOT=pathlib.Path(__file__).resolve().parent
SOURCE=ROOT/"source_pilot.py"
s=importlib.util.spec_from_file_location("w_f46_existing_source_audit",SOURCE)
model=importlib.util.module_from_spec(s)
sys.modules[s.name]=model
s.loader.exec_module(model)

ARMS=("query_edge_voc_per_compute","one_step_value","exact_entropy",
      "query_frontier","budget_matched_rank","no_purchase","evaluator_oracle_unfair")
FAIR=("one_step_value","exact_entropy","query_frontier","budget_matched_rank")

def audit(payload):
    data=json.loads(pathlib.Path(payload).read_text(encoding="utf-8"))
    assert data["schema"]=="wingless.f46.query-edge-voc.trainonly.v1"
    assert data["stage"]=="TRAIN_ONLY_NOT_ORIGINAL"
    assert data["questions"]==96 and len(data["rows"])==96
    assert data["seeds"]==[391100,391112] and data["families"]==[500,501,502,503]
    assert data["fit_seeds_not_consumed"]==[391000,391012]
    assert data["reserved_originals_not_consumed"]==[392000,392016]
    assert data["old_originals_not_consumed"]==[372000,372016]
    assert data["original_executed"] is False
    assert data["independent_F3_custody"] is False
    assert data["native_576_byte_resident_qualified"] is False
    assert data["accepted_ref_mutation"] is False and data["paid_api_calls"]==0
    assert len(data["scores"])==len(ARMS) and set(data["scores"])==set(ARMS)
    assert set(data["posterior_world_work"])==set(ARMS)-{"evaluator_oracle_unfair"}
    assert len(data["by_family"])==4
    seen=set()
    cost_totals={a:0 for a in ARMS[:-1]}
    counts={a:0 for a in ARMS}
    for row in data["rows"]:
        seed,fam,length=row["seed"],row["family"],row["length"]
        assert seed in model.DIAG and fam in model.FAMILIES and length in (3,5)
        assert (seed,fam,length) not in seen
        seen.add((seed,fam,length))
        private,public,start,chain=model.source(seed,fam,length)
        known={(r,x) for r,x,y in public}
        truth=model.f45.prior.prior.follow(private,start,chain)
        assert row["truth"]==truth
        assert set(row["answers"])==set(ARMS)
        assert set(row["receipts"])==set(ARMS)-{"evaluator_oracle_unfair"}
        assert set(row["work"])==set(row["receipts"])
        assert row["answers"]["evaluator_oracle_unfair"]==truth
        for arm in ARMS[:-1]:
            receipts=row["receipts"][arm]
            assert len(receipts)<=2
            assert arm!="no_purchase" or not receipts
            observed=list(public)
            purchased=set()
            for pair in receipts:
                assert len(pair)==3
                r,x,y=pair
                assert isinstance(r,int) and isinstance(x,int) and isinstance(y,int)
                assert 0<=r<4 and 0<=x<6 and 0<=y<6
                assert (r,x) not in known and (r,x) not in purchased
                assert private[r][x]==y
                purchased.add((r,x))
                observed.append((r,x,y))
            # Reconstruct terminal answer from paid receipts and public mappings.
            pred=model.f45.prior.prior.predict(observed,start,chain)
            assert row["answers"][arm]==pred
            work=row["work"][arm]
            assert isinstance(work,int) and work>=0
            if arm=="no_purchase":assert work==0
            cost_totals[arm]+=work
        for arm in ARMS:
            counts[arm]+=int(row["answers"][arm]==truth)
    assert len(seen)==96
    assert counts["evaluator_oracle_unfair"]==96
    for arm in ARMS:
        score=data["scores"][arm]
        assert score["correct"]==counts[arm] and score["total"]==96
        assert math.isclose(score["accuracy"],counts[arm]/96,abs_tol=1e-12)
    for arm in ARMS[:-1]:
        assert data["posterior_world_work"][arm]==cost_totals[arm]
    for fam in model.FAMILIES:
        subset=[row for row in data["rows"] if row["family"]==fam]
        assert len(subset)==24
        for arm in ARMS:
            assert data["by_family"][str(fam)][arm]==sum(
                row["answers"][arm]==row["truth"] for row in subset)
    strongest=max(FAIR,key=lambda a:(counts[a]/96,-cost_totals[a]))
    cand="query_edge_voc_per_compute"
    gain=(counts[cand]-counts[strongest])/96
    ratio=cost_totals[cand]/max(1,cost_totals[strongest])
    nondegrade=all(data["by_family"][str(f)][cand]>=max(
        data["by_family"][str(f)][a] for a in FAIR) for f in model.FAMILIES)
    qualifies=gain>=.03 and ratio<=1.2 and nondegrade
    assert data["strongest_budget_fair_control"]==strongest
    assert math.isclose(data["candidate_gain"],gain,abs_tol=1e-12)
    assert math.isclose(data["posterior_work_ratio"],ratio,abs_tol=1e-12)
    assert data["family_nondegradation"]==nondegrade
    assert data["frozen_criterion_met"]==qualifies
    assert data["verdict"]==("TRAIN_ONLY_SUPPORTED" if qualifies else "TRAIN_ONLY_NEGATIVE")
    summary={"schema":"wingless.f46.existing-train-audit.v1",
             "mode":"READ_ONLY_EXISTING_TRAIN_ARTIFACT_NOT_PRIMARY_ORIGINAL",
             "source_run":38803701796,"cases":96,"scores":counts,
             "strongest_fair":strongest,"gain":gain,
             "measured_posterior_work":cost_totals,"work_ratio":ratio,
             "per_family_nondegradation":nondegrade,
             "train_gate_pass":qualifies,
             "independent_f3_custody":False,"original_executed":False,
             "accepted_ref_mutated":False}
    print(json.dumps(summary,indent=2,sort_keys=True))
    return summary

if __name__=="__main__":
    if len(sys.argv)!=2:raise SystemExit("usage: independent_audit.py RESULT_JSON")
    audit(sys.argv[1])
