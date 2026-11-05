"""CNC-G9P: capacity reconciliation to independently witnessed FORGE G59 GPU bytes.

No current hardware probing, controller/runner interaction, CUDA model use, or
claim of present free memory. Original code/weight planning remains G9B/D/H/K.
"""
import json
from g9b_tensor_budget import VRAM_CEILING, FROZEN_CONFIG, expected_tensors, make_report
from g9h_sparse_budget import sparse_budget
from g9k_lowrank_budget import report as lowrank_report
from g9d_mixed_precision_budget import mixed_report

MIB = 1 << 20
HARDWARE_EVIDENCE = {
    "original_repo": "spoonman136668-ai/FORGE",
    "original_run_id": 41030389354,
    "original_artifact_id": 12966370600,
    "original_zip_sha256": "c6d368b70e6697e4ec699c7943100401214eed5fcea17688e3ed6351b564bac2",
    "gpu_name": "NVIDIA GeForce RTX 5070 Ti Laptop GPU",
    "nvidia_driver": "596.13",
    "physical_total_mib": 12227,
    "historical_free_mib": 11812,
}


def _archive_total(case, physical_total, historical_free):
    cost = case["planned_total_bytes"]
    if type(cost) is not int or cost <= 0:
        raise ValueError("unverified planning cost")
    return {
        "name": case["name"],
        "total_planned_bytes_including_hypothetical_workspace": cost,
        "remaining_against_historical_physical_total_bytes": physical_total-cost,
        "would_fit_physical_total_analytically": cost<=physical_total,
        "remaining_against_historical_free_snapshot_bytes": historical_free-cost,
        "would_fit_historical_free_snapshot_analytically": cost<=historical_free,
        "current_available_vram_bytes": None,
        "cuda_context_peak_measured": False,
        "full_original_model_forward_executed": False,
    }


def report():
    # Fixture uses G9B's exactly frozen architecture, not fresh original shard data.
    shapes=expected_tensors(dict(FROZEN_CONFIG))
    index={"weight_map": {k: f"model-{i%16+1:05d}-of-00016.safetensors"
                           for i,k in enumerate(sorted(shapes))},
           "metadata": {"total_size":sum(2*__import__("math").prod(shape)
                                    for shape,_ in shapes.values())}}
    baseline=make_report(dict(FROZEN_CONFIG),index,source_sha="9"*40)
    int2=next(c for c in baseline["scenarios"] if c["name"]=="int2_block256_fp16")
    sparse=sparse_budget()
    lowrank=lowrank_report()
    mixed=mixed_report(dict(FROZEN_CONFIG),index,source_sha="9"*40)
    total=HARDWARE_EVIDENCE["physical_total_mib"]*MIB
    free=HARDWARE_EVIDENCE["historical_free_mib"]*MIB
    if VRAM_CEILING!=12*(1<<30) or total>=VRAM_CEILING or not 0<free<=total:
        raise ValueError("frozen physical VRAM evidence or nominal budget drift")
    if any(r["baseline_remaining_bytes"]!=int2["remaining_bytes"]
           for r in (sparse, lowrank)):
        raise ValueError("G9B/G9H/G9K planning source disagreement")
    candidates=[{"name":"int2_block256_fp16","planned_total_bytes":int2["planned_total_bytes"]}]
    candidates += [{"name":"sparse_k"+str(x["corrections_per_row"]),
                    "planned_total_bytes":x["model_total_planned_bytes"]}
                   for x in sparse["cases"]]
    candidates += [{"name":"lowrank_rank"+str(x["rank"]),
                    "planned_total_bytes":VRAM_CEILING-x["remaining_after_hypothetical_workspace_bytes"]}
                   for x in lowrank["cases"]]
    candidates += [{"name":"mixed_"+x["name"],"planned_total_bytes":x["total_planned_bytes"]}
                   for x in mixed["scenarios"]]
    if len({x["name"] for x in candidates})!=len(candidates):
        raise ValueError("duplicate budget method")
    scenarios=[_archive_total(x,total,free) for x in candidates]
    return {
        "schema":"wingless.cnc.g9p.historical-hardware-rebased-budget.v1",
        "result":"HISTORICAL_HARDWARE_STORAGE_REBASE_NOT_RUNTIME_QUALIFICATION",
        "hardware_receipt":dict(HARDWARE_EVIDENCE),
        "gpu_total_bytes_from_archived_driver":total,
        "gpu_free_bytes_at_archived_probe":free,
        "historical_total_below_nominal_12gib_bytes":VRAM_CEILING-total,
        "historical_free_below_nominal_12gib_bytes":VRAM_CEILING-free,
        "context_tokens":8192,
        "hypothetical_workspace_bytes":1<<30,
        "scenario_count":len(scenarios),
        "scenarios":scenarios,
        "source_architecture":"G9B_SOURCE_RECONSTRUCTED_FROM_PRIOR_PINNED_INDEX_CUSTODY",
        "hardware_measurement_is_current":False,
        "runner_online_now":"UNVERIFIED",
        "runner_usage_authorized":"UNVERIFIED",
        "cuda_runtime_installed":"UNVERIFIED",
        "original_model_shards_on_runner":"UNVERIFIED",
        "real_gpu_residency":"UNVERIFIED",
        "original_model_activations":"UNVERIFIED",
        "model_quality":"UNVERIFIED",
        "scientific_acceptance":"NOT_GRANTED",
    }


if __name__=="__main__":
    print("CNC_G9P_HISTORICAL_GPU_REBASE "+json.dumps(report(),sort_keys=True))
