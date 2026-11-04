"""CNC-02B: independent FIT/TRANSFER adaptation on authentic-weight-compatible windows.

This pure module accepts a weight window externally; it does not download source
weights, read protected heldout labels, or run Qwen model inference. Development
targets are the original weight's output under independent synthetic Gaussian
inputs. FIT and TRANSFER have disjoint RNG seeds; no TRANSFER update or labels
enter candidate selection. Direct packed edits never retain an FP32 master.
"""
import numpy as np
from compressed_learning import PackedCodeLearner
from packed_int4 import pack_int4, apply_packed_int4

FIT_SEED = 6011
TRANSFER_SEED = 17201
ROWS = 128
COLS = 256
FIT_SIZE = 32
TRANSFER_SIZE = 32
STEPS = 2


def _error(estimate, truth):
    a = np.asarray(estimate, dtype=np.float64)
    b = np.asarray(truth, dtype=np.float64)
    energy = float(np.sum(b*b))
    if energy <= 0 or a.shape != b.shape:
        raise ValueError("invalid normalized error")
    return float(np.sum((a-b)**2) / energy)


def compare_window(weight, *, steps=STEPS):
    w = np.asarray(weight, dtype=np.float32)
    if w.shape != (ROWS, COLS) or not np.isfinite(w).all() or not np.any(w):
        raise ValueError("invalid original-weight-compatible window")
    if type(steps) is not int or not 1 <= steps <= STEPS:
        raise ValueError("unfrozen update budget")
    x_fit = np.random.default_rng(FIT_SEED).standard_normal((FIT_SIZE, COLS)).astype(np.float32)
    x_transfer = np.random.default_rng(TRANSFER_SEED).standard_normal(
        (TRANSFER_SIZE, COLS)).astype(np.float32)
    target_fit = x_fit @ w.T
    target_transfer = x_transfer @ w.T
    learner = PackedCodeLearner.from_weight(w)
    initial_codes = learner.packed.copy()
    initial_scales = learner.scales.copy()
    baseline_fit = _error(learner.predict(x_fit), target_fit)
    baseline_transfer = _error(learner.predict(x_transfer), target_transfer)
    decisions = []
    for _ in range(steps):
        decisions.append(learner.train_step(x_fit, target_fit))
    final_fit = _error(learner.predict(x_fit), target_fit)
    final_transfer = _error(learner.predict(x_transfer), target_transfer)
    # Exact deterministic frozen/re-encoded baseline: same codes, no updates.
    frozen_fit = _error(apply_packed_int4(x_fit, initial_codes, initial_scales), target_fit)
    frozen_transfer = _error(
        apply_packed_int4(x_transfer, initial_codes, initial_scales), target_transfer)
    if baseline_fit != frozen_fit or baseline_transfer != frozen_transfer:
        raise ValueError("frozen INT4 control drift")
    if not np.array_equal(initial_scales, learner.scales):
        raise ValueError("compressed scales changed")
    if learner.steps != steps:
        raise ValueError("invalid training step count")
    return {
        "scope": "AUTHENTIC_WEIGHT_COMPATIBLE_SYNTHETIC_FIT_TRANSFER_ONLY",
        "fit_seed": FIT_SEED,
        "transfer_seed": TRANSFER_SEED,
        "fit_size": FIT_SIZE,
        "transfer_size": TRANSFER_SIZE,
        "updates_allowed": steps,
        "frozen_int4_fit_error": frozen_fit,
        "frozen_int4_transfer_error": frozen_transfer,
        "compressed_direct_fit_error": final_fit,
        "compressed_direct_transfer_error": final_transfer,
        "code_edits": sum(int(d["changed"]) for d in decisions),
        "persistent_compressed_bytes": learner.resident_weight_bytes,
        "fp32_master_copy_retained": False,
        "native_qwen_activations": "NOT_CAPTURED",
        "full_model_training": False,
        "gpu_execution": False,
        "scientific_acceptance": "NOT_GRANTED",
    }
