"""CNC-02D: fixed-budget compressed-learning vs FP32 and QAT controls.

Development-only numerical comparator. All methods start from the SAME
blockwise INT4 reconstruction. FP32 and QAT retain a float32 master matrix;
the direct-code arm does not. All receive two FIT gradient/update steps but
their updated parameter counts and compute are NOT equivalent. This is a
matched step-count diagnostic, not equal-FLOP, equal-byte or task evaluation.
TRANSFER targets never participate in updates. No Qwen forward/backprop.
"""
import numpy as np

from compressed_learning import PackedCodeLearner, _code_at
from cnc02b_fit_transfer import FIT_SEED, TRANSFER_SEED, ROWS, COLS, FIT_SIZE, TRANSFER_SIZE, STEPS, _error
from packed_int4 import pack_int4, apply_packed_int4

LEARNING_RATE = 0.02
BLOCK = 64


def reconstruct(codes, scales):
    """Offline control-only matrix dequantization; never used by direct learner."""
    identity = np.eye(COLS, dtype=np.float32)
    return apply_packed_int4(identity, codes, scales, block=BLOCK).T


def _sgd_step(master, x, y):
    residual = x @ master.T - y
    gradient = residual.T @ x / x.shape[0]
    denominator = max(1., float(np.sum(x.astype(np.float64)**2) / x.shape[0]))
    master -= np.float32(LEARNING_RATE/denominator) * gradient.astype(np.float32)


def _fixed_scale_requantized(master, scales):
    """Round candidate into the same immutable-scale codebook as direct learner."""
    groups = np.repeat(scales.astype(np.float32), BLOCK, axis=1)
    codes = np.clip(np.rint(master/groups), -7, 7).astype(np.int8)
    lo = codes[:, 0::2].astype(np.uint8) & 15
    hi = codes[:, 1::2].astype(np.uint8) & 15
    return (lo | (hi << 4)).astype(np.uint8)


def compare_controls(weight, *, steps=STEPS):
    w = np.asarray(weight, dtype=np.float32)
    if (w.shape != (ROWS, COLS) or not np.isfinite(w).all() or not np.any(w)
            or type(steps) is not int or steps != STEPS):
        raise ValueError("unfrozen CNC-02D source or step count")
    xfit = np.random.default_rng(FIT_SEED).standard_normal(
        (FIT_SIZE, COLS)).astype(np.float32)
    xtransfer = np.random.default_rng(TRANSFER_SEED).standard_normal(
        (TRANSFER_SIZE, COLS)).astype(np.float32)
    yfit, ytransfer = xfit @ w.T, xtransfer @ w.T

    packed, scales = pack_int4(w)
    frozen_fit = _error(apply_packed_int4(xfit, packed, scales), yfit)
    frozen_transfer = _error(apply_packed_int4(xtransfer, packed, scales), ytransfer)

    direct = PackedCodeLearner(packed, scales)
    initial_codes = direct.packed.copy()
    direct_updates = [direct.train_step(xfit, yfit) for _ in range(steps)]
    if not np.array_equal(scales, direct.scales):
        raise ValueError("direct scale drift")

    # These are conventional control arms with persistent FP32 masters.
    initial_master = reconstruct(packed, scales)
    fp32 = initial_master.copy()
    qat_master = initial_master.copy()
    for _ in range(steps):
        _sgd_step(fp32, xfit, yfit)
        qat_codes = _fixed_scale_requantized(qat_master, scales)
        qat_shadow = reconstruct(qat_codes, scales)
        # Straight-through estimator proxy: gradient evaluated at quantized
        # forward weights but applied to persistent FP32 master.
        residual = xfit @ qat_shadow.T - yfit
        gradient = residual.T @ xfit / xfit.shape[0]
        denominator = max(1., float(np.sum(xfit.astype(np.float64)**2) / xfit.shape[0]))
        qat_master -= np.float32(LEARNING_RATE/denominator)*gradient.astype(np.float32)
    final_qat_codes = _fixed_scale_requantized(qat_master, scales)

    def errors(estimate):
        return {"fit_relative_squared_error": _error(estimate(xfit), yfit),
                "transfer_relative_squared_error": _error(estimate(xtransfer), ytransfer)}
    baseline = {
        "fit_relative_squared_error": frozen_fit,
        "transfer_relative_squared_error": frozen_transfer}
    direct_errors = errors(direct.predict)
    fp32_errors = errors(lambda x: x @ fp32.T)
    qat_errors = errors(lambda x: apply_packed_int4(x, final_qat_codes, scales))
    if not np.array_equal(initial_codes, packed):
        raise ValueError("frozen representation drift")
    if direct_errors["fit_relative_squared_error"] > frozen_fit + 1e-7:
        raise ValueError("direct training unexpectedly worsened FIT")
    return {
        "status": "SYNTHETIC_ACTIVATION_DEVELOPMENT_COMPARISON_ONLY",
        "fit_seed": FIT_SEED, "transfer_seed": TRANSFER_SEED,
        "steps_each": steps, "learning_rate_fp32_qat": LEARNING_RATE,
        "full_precision_start": "SAME_INT4_RECONSTRUCTION",
        "gradient_budget_matching": "UPDATE_COUNT_ONLY_NOT_MATCHED_FLOPS_OR_PARAMETERS",
        "direct_code_edits": sum(int(u["changed"]) for u in direct_updates),
        "frozen_int4": baseline, "direct_int4": direct_errors,
        "fp32_sgd": fp32_errors, "qat_ste": qat_errors,
        "persistent_weight_bytes": {
            "direct_int4": direct.resident_weight_bytes,
            "fp32_sgd": int(fp32.nbytes),
            "qat_ste_master_plus_codes_and_scales": int(qat_master.nbytes + final_qat_codes.nbytes + scales.nbytes),
        },
        "transient_peak_rss_bytes": "UNMEASURED",
        "gpu_peak_vram_bytes": "UNMEASURED",
        "original_qwen_activations": "NOT_CAPTURED",
        "full_model_training": False, "scientific_acceptance": "NOT_GRANTED"
    }
