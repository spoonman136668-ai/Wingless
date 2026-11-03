"""Development-only activation-output replay; NOT model-quality acceptance.

This module does not capture activations and cannot authenticate their provenance.
Real model-forward-hook traces require independent custody before interpretation.
"""
import hashlib
import json
import re

import numpy as np

MODEL = "Qwen/Qwen3-30B-A3B"
REVISION = "ad44e777bcd18fa416d9da3bd8f70d33ebb85d39"
_SHA256 = re.compile(r"[0-9a-f]{64}\Z")
_SHA40 = re.compile(r"[0-9a-f]{40}\Z")
_TENSOR = re.compile(r"model\.layers\.\d+\.mlp\.experts\.\d+\.(?:down_proj|gate_proj|up_proj)\.weight\Z")


def activation_digest(activations):
    """SHA256 of shape, dtype and little-endian row-major float32 activation bytes."""
    x = np.asarray(activations)
    if x.ndim != 2 or not np.issubdtype(x.dtype, np.floating) or not np.isfinite(x).all():
        raise ValueError("activations must be a finite 2D floating array")
    data = np.ascontiguousarray(x, dtype="<f4")
    h = hashlib.sha256()
    h.update(json.dumps({"dtype": "float32-le", "shape": list(data.shape)}, sort_keys=True, separators=(",", ":")).encode("utf-8"))
    h.update(b"\n")
    h.update(data.tobytes(order="C"))
    return h.hexdigest()


def validate_provenance(activations, provenance):
    """Validate declared identity, separation and checksum; NOT provenance authenticity."""
    if not isinstance(provenance, dict):
        raise ValueError("missing replay provenance")
    if provenance.get("model") != MODEL or provenance.get("revision") != REVISION:
        raise ValueError("checkpoint identity mismatch")
    if provenance.get("cohort") != "development":
        raise ValueError("protected or unfrozen cohort not permitted")
    name = provenance.get("tensor_name")
    if not isinstance(name, str) or _TENSOR.fullmatch(name) is None:
        raise ValueError("unsupported expert tensor")
    x = np.asarray(activations)
    if x.ndim != 2 or not 1 <= x.shape[0] <= 8192 or not 1 <= x.shape[1] <= 8192:
        raise ValueError("invalid trace dimensions")
    ids = provenance.get("record_ids")
    if not isinstance(ids, list) or len(ids) != len(x) or any(not isinstance(i, str) or not i for i in ids) or len(set(ids)) != len(ids):
        raise ValueError("trace identifiers missing or duplicated")
    digest = provenance.get("activation_sha256")
    if not isinstance(digest, str) or _SHA256.fullmatch(digest) is None or digest != activation_digest(x):
        raise ValueError("activation checksum mismatch")
    kind = provenance.get("capture_kind")
    if kind not in ("synthetic_test_fixture", "native_model_forward_hook"):
        raise ValueError("capture kind must be explicit")
    if kind == "native_model_forward_hook":
        source = provenance.get("source_commit")
        receipt = provenance.get("source_artifact_sha256")
        run_id = provenance.get("source_workflow_run_id")
        calibration = provenance.get("calibration_record_ids")
        if not isinstance(source, str) or _SHA40.fullmatch(source) is None:
            raise ValueError("original forward-hook source commit missing")
        if not isinstance(receipt, str) or _SHA256.fullmatch(receipt) is None:
            raise ValueError("original activation trace custody digest missing")
        if not isinstance(run_id, int) or isinstance(run_id, bool) or run_id <= 0:
            raise ValueError("source workflow run missing")
        if not isinstance(calibration, list) or not calibration or any(not isinstance(v, str) or not v for v in calibration) or len(calibration) != len(set(calibration)) or set(calibration).intersection(ids):
            raise ValueError("calibration and evaluation records not independently separated")
    return kind


def evaluate_activation_output(reference_weight, candidate_weights, activations, provenance):
    """Compute output distortion for frozen expert-input activations and candidates.

    Cannot claim authentic trace custody, model quality, peak VRAM or acceptance.
    """
    kind = validate_provenance(activations, provenance)
    w = np.asarray(reference_weight, dtype=np.float32)
    x = np.asarray(activations, dtype=np.float32)
    if w.ndim != 2 or w.shape[1] != x.shape[1] or not np.isfinite(w).all():
        raise ValueError("invalid reference matrix or input dimension")
    if not isinstance(candidate_weights, dict) or not candidate_weights:
        raise ValueError("at least one named candidate required")
    y = x @ w.T
    energy = float(np.sum(y.astype(np.float64) ** 2))
    if not np.isfinite(energy) or energy <= 0:
        raise ValueError("zero-energy reference expert output")
    w_energy = float(np.sum(w.astype(np.float64) ** 2))
    if w_energy <= 0:
        raise ValueError("zero-energy reference weights")
    result = {}
    for name, value in candidate_weights.items():
        if not isinstance(name, str) or not name or name == "original":
            raise ValueError("invalid candidate name")
        estimate = np.asarray(value, dtype=np.float32)
        if estimate.shape != w.shape or not np.isfinite(estimate).all():
            raise ValueError("invalid candidate matrix")
        output_delta = x @ (w - estimate).T
        result[name] = {
            "relative_output_squared_error": float(np.sum(output_delta.astype(np.float64) ** 2) / energy),
            "relative_weight_squared_error": float(np.sum((w.astype(np.float64) - estimate.astype(np.float64)) ** 2) / w_energy),
        }
    return {
        "source_kind": kind,
        "trace_activation_sha256": provenance["activation_sha256"],
        "tensor_name": provenance["tensor_name"],
        "trace_records": len(x),
        "candidate_metrics": result,
        "source_authenticity": "UNVERIFIED_REQUIRES_INDEPENDENT_ORIGINAL_ARTIFACT_INSPECTION",
        "acceptance": "DEVELOPMENT_DIAGNOSTIC_ONLY_NO_TASK_QUALITY_OR_INFERENCE_CLAIM",
    }
