"""G22 read-only route audit. Declarations never grant source or execution authority.

Validates whether an external inference/capture route is even *compatible* with
the frozen original-model activation protocol. No model execution or downloads.
"""
import re

from activation_replay import MODEL, REVISION

SHA256 = re.compile(r"[a-f0-9]{64}\\Z")
GIT_SHA = re.compile(r"[a-f0-9]{40}\\Z")
EXPERT = re.compile(r"model\\.layers\\.\\d+\\.mlp\\.experts\\.\\d+\\.(?:down_proj|gate_proj|up_proj)\\.weight\\Z")


def audit_route(proposal, pinned_manifest):
    """Return bounded incompatibility reasons, never a PASS/authorization.

    The proposal is self-declared. A structurally complete proposal still needs
    independently inspected original execution, actual byte custody, model
    forward-hook evidence, and approved execution authority.
    """
    blockers = []
    if not isinstance(proposal, dict):
        return {"status": "BLOCKED", "blockers": ["ROUTE_DESCRIPTOR_MISSING"]}
    if not isinstance(pinned_manifest, dict) or pinned_manifest.get("model") != MODEL or pinned_manifest.get("revision") != REVISION:
        return {"status": "BLOCKED", "blockers": ["PINNED_MANIFEST_IDENTITY_MISMATCH"]}
    reference_files = pinned_manifest.get("files")
    if not isinstance(reference_files, dict) or len(reference_files) != 18 or any(not isinstance(v, str) or not SHA256.fullmatch(v) for v in reference_files.values()):
        blockers.append("INCOMPLETE_PINNED_CHECKPOINT_MANIFEST")
    if proposal.get("model") != MODEL or proposal.get("revision") != REVISION:
        blockers.append("WRONG_ORIGINAL_MODEL")
    if proposal.get("checkpoint_format") != "original_safetensors":
        blockers.append("INCOMPATIBLE_CHECKPOINT_FORMAT")
    if proposal.get("runtime") != "native_full_model_forward":
        blockers.append("NO_ORIGINAL_MODEL_FORWARD_ROUTE")
    if proposal.get("capture_kind") != "native_expert_input_forward_hook":
        blockers.append("ROUTER_TELEMETRY_IS_NOT_EXPERT_ACTIVATIONS")
    tensor = proposal.get("tensor_name")
    if not isinstance(tensor, str) or not EXPERT.fullmatch(tensor):
        blockers.append("INVALID_EXPERT_TENSOR")
    if not isinstance(proposal.get("source_commit"), str) or not GIT_SHA.fullmatch(proposal["source_commit"]):
        blockers.append("UNPINNED_RUNTIME_SOURCE")
    manifest_files = proposal.get("checkpoint_files")
    if not isinstance(manifest_files, dict) or manifest_files != reference_files:
        blockers.append("ORIGINAL_CHECKPOINT_FILE_IDENTITY_MISSING")
    if proposal.get("cohort") != "development":
        blockers.append("PROTECTED_OR_UNFROZEN_COHORT_FORBIDDEN")
    calibration = proposal.get("calibration_record_ids")
    evaluation = proposal.get("evaluation_record_ids")
    if (not isinstance(calibration, list) or not calibration or
        not isinstance(evaluation, list) or not evaluation or
        any(not isinstance(v, str) or not v for v in calibration + evaluation) or
        len(set(calibration)) != len(calibration) or len(set(evaluation)) != len(evaluation) or
        not set(calibration).isdisjoint(evaluation)):
        blockers.append("CALIBRATION_EVALUATION_SEPARATION_UNVERIFIED")
    if not isinstance(proposal.get("authority_reference"), str) or not proposal["authority_reference"].strip():
        blockers.append("EXECUTION_AUTHORITY_NOT_IDENTIFIED")
    if not isinstance(proposal.get("qualified_runner_reference"), str) or not proposal["qualified_runner_reference"].strip():
        blockers.append("QUALIFIED_RUNTIME_NOT_IDENTIFIED")

    return {
        "status": "BLOCKED" if blockers else "DECLARATION_COMPATIBLE_SOURCE_UNVERIFIED",
        "blockers": blockers,
        "source_authenticity": "UNVERIFIED",
        "execution_authorized_by_this_audit": False,
        "scientific_acceptance": "NOT_GRANTED",
        "requires": [
            "independent original runtime source and runner qualification",
            "actual immutable 18-file checkpoint byte verification in the execution session",
            "original model forward-hook and activation evidence with immutable source artifact",
            "separately verified operation authorization and native memory feasibility",
        ],
    }
