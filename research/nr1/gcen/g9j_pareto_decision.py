"""CNC-G9J: immutable evidence-based compression Pareto decision.

Source data are explicitly frozen G9I authentic original-weight development
aggregate receipts. Only computes relative storage/fidelity diagnostics.
No new original shard consumption, inference, GPU, heldout or acceptance.
"""
import json
from fractions import Fraction

SOURCE_SHA = "f4b60b8de510d654211f5ecac483084e79997b62"
RUN_ID = 41802853750
ARTIFACT_ID = 13308383049
ARTIFACT_ZIP_SHA256 = "01688b9f0a999b73fdc4a8dd140f4eced590cfed41856940a5dfbe9b336643c8"
REPORT_SHA256 = "b5c57f908c8f1d2050be5f77db70ef7990a964806b313f75b4d958ab80f221ad"
G9I_SEEDS = (107, 163, 229, 271, 347)
G9I_BASELINE = (
    # Method, average relative output squared error (16 expert), actual
    # packed bytes for 16 BF16 128x256 original expert weight windows.
    ("bf16_reference", 0.0, 1048576),
    ("int4_b64_fp32", 0.012366740722513321, 295296),
    ("companded_int2_b256", 0.12540784166370125, 135552),
    ("sparse_k1", 0.11289221801679578, 152000),
    ("sparse_k2", 0.10552646324604775, 160192),
    ("sparse_k3", 0.10006517961641721, 168384),
)
# G9H exact architectural static storage model with hypothetical workspace.
G9H_HEADROOM = {
    "companded_int2_b256": 448974848,
    "sparse_k1": 272740352,
    "sparse_k2": 184659968,
    "sparse_k3": 96579584,
}


def decision_report():
    names = [a[0] for a in G9I_BASELINE]
    if len(set(names)) != len(names):
        raise ValueError("duplicate frozen arm")
    cases = []
    by_name = {name: (error, size) for name, error, size in G9I_BASELINE}
    if by_name["bf16_reference"] != (0.0, 1048576):
        raise ValueError("frozen BF16 extent drift")
    prior_error, prior_bytes = by_name["companded_int2_b256"]
    for name, error, size in G9I_BASELINE:
        if not (type(error) is float and error >= 0 and
                type(size) is int and 0 < size <= 1048576):
            raise ValueError("invalid frozen numerical/storage evidence")
        row = {
            "method": name,
            "relative_output_squared_error": error,
            "actual_window_stored_bytes": size,
            "compression_ratio_vs_original_bf16": 1048576 / size,
            "original_weight_window_evidence": "AUTHENTIC_PINNED_BF16",
        }
        if name in G9H_HEADROOM:
            row["whole_model_analytical_headroom_bytes"] = G9H_HEADROOM[name]
            row["real_gpu_measured"] = False
        if name.startswith("sparse_k"):
            row["relative_output_squared_error_reduction_vs_int2"] = (
                (prior_error - error) / prior_error
            )
            row["extra_physical_window_bytes_vs_int2"] = size-prior_bytes
            row["additional_bytes_per_unit_output_squared_error_reduced"] = (
                (size-prior_bytes)/(prior_error-error)
            )
        cases.append(row)
    for a, b in zip(("sparse_k1","sparse_k2"),("sparse_k2","sparse_k3")):
        if not (by_name[b][0] < by_name[a][0] and by_name[b][1] > by_name[a][1]):
            raise ValueError("sparse Pareto shape changed")
    if not all(by_name[name][0] > by_name["int4_b64_fp32"][0] * 8
               for name in G9H_HEADROOM):
        raise ValueError("frozen INT4 gap assumption changed")
    return {
        "schema": "wingless.cnc.g9j.authentic-weight-development-pareto.v1",
        "source_sha": SOURCE_SHA,
        "source_run": RUN_ID,
        "authentic_original_artifact_id": ARTIFACT_ID,
        "artifact_zip_sha256": ARTIFACT_ZIP_SHA256,
        "inside_report_sha256": REPORT_SHA256,
        "frozen_seeds": list(G9I_SEEDS),
        "comparison_scope": "16_REAL_QWEN_BF16_EXPERT_WINDOWS_SYNTHETIC_ACTIVATION_INPUTS",
        "cases": cases,
        "decision": "NO_GO_DIRECT_MODEL_TRIAL_FOR_UNIFORM_K1_K2_K3_SPARSE_CORRECTIONS",
        "reason": "Each G9I variant retains >8x INT4 output squared distortion and consumes extra bytes; no model quality or GPU admission established.",
        "next_gate": "NEW_DISTINCT_REPRESENTATION_WITH_SEPARATE_FROZEN_CUSTODY_AND_ACTUAL_MODEL_FORWARD_GUARDS",
        "source_original_weights_authenticated": True,
        "original_qwen_activations_captured": False,
        "full_model_forward": False,
        "gpu_residency_measured": False,
        "protected_heldout_used": False,
        "scientific_acceptance": "NOT_GRANTED",
    }


if __name__ == "__main__":
    print("CNC_G9J_PARETO_DEVELOPMENT_ONLY "+json.dumps(decision_report(),sort_keys=True))
