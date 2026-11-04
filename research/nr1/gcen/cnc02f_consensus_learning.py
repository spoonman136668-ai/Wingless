"""CNC-02F: frozen cross-FIT-consensus direct packed-code adaptation.

Fresh development cohort (not CNC-02C/E FIT or TRANSFER). For each of two
steps, propose one INT4 coordinate edit on FIT-A, and accept it only if a
disjoint FIT-B batch also improves. Transfer is never used for acceptance.
All persistent learning state is packed codes and frozen block scales.

The FP32 SGD arm starts from the identical INT4 reconstruction and uses the
same FIT-A inputs for exactly two steps; step count does not equalize FLOPs,
trainable parameters or transient workspace. This does NOT run Qwen inference.
"""
import numpy as np
from compressed_learning import PackedCodeLearner
from cnc02b_fit_transfer import ROWS, COLS, STEPS, _error
from cnc02d_matched_controls import reconstruct, _sgd_step
from packed_int4 import pack_int4, apply_packed_int4

FIT_A_SEED = 26061
FIT_B_SEED = 26087
TRANSFER_SEED = 26111
COHORT_SIZE = 32


def _cohort(weight, seed):
    x = np.random.default_rng(seed).standard_normal(
        (COHORT_SIZE, COLS)).astype(np.float32)
    return x, x @ weight.T


def compare_consensus(weight, *, steps=STEPS):
    w = np.asarray(weight, dtype=np.float32)
    if (w.shape != (ROWS, COLS) or not np.isfinite(w).all() or
            not np.any(w) or type(steps) is not int or steps != STEPS):
        raise ValueError("unfrozen CNC-02F source or step count")
    x_a, y_a = _cohort(w, FIT_A_SEED)
    x_b, y_b = _cohort(w, FIT_B_SEED)
    x_t, y_t = _cohort(w, TRANSFER_SEED)
    packed, scales = pack_int4(w)
    direct = PackedCodeLearner(packed, scales)
    fp32 = reconstruct(packed, scales)
    frozen = apply_packed_int4(x_t, packed, scales)
    baseline_transfer = _error(frozen, y_t)
    accepted = rejected = 0
    for _ in range(steps):
        old_codes = direct.packed.copy()
        # Proposal generated ONLY from FIT-A, never TRANSFER.
        proposal = direct.train_step(x_a, y_a)
        if proposal["changed"]:
            before_b = _error(apply_packed_int4(x_b, old_codes, scales), y_b)
            after_b = _error(direct.predict(x_b), y_b)
            if after_b < before_b:
                accepted += 1
            else:
                direct.packed[:] = old_codes
                rejected += 1
        _sgd_step(fp32, x_a, y_a)
    if not np.array_equal(direct.scales, scales) or direct.steps != steps:
        raise ValueError("direct-state invariant drift")
    if accepted + rejected > steps:
        raise ValueError("invalid compressed update budget")
    return {
        "status": "FROZEN_FRESH_COHORT_SYNTHETIC_DEVELOPMENT_ONLY",
        "fit_a_seed": FIT_A_SEED, "fit_b_seed": FIT_B_SEED,
        "transfer_seed": TRANSFER_SEED, "examples_per_cohort": COHORT_SIZE,
        "step_count_each": steps,
        "compressed_candidate_updates_accepted": accepted,
        "compressed_candidate_updates_rejected": rejected,
        "frozen_int4_transfer_error": baseline_transfer,
        "consensus_int4_transfer_error": _error(direct.predict(x_t), y_t),
        "fp32_sgd_transfer_error": _error(x_t @ fp32.T, y_t),
        "consensus_fit_a_error": _error(direct.predict(x_a), y_a),
        "consensus_fit_b_error": _error(direct.predict(x_b), y_b),
        "persistent_weight_bytes": {
            "direct_compressed": direct.resident_weight_bytes,
            "fp32_sgd": int(fp32.nbytes),
        },
        "equal_compute_or_parameter_updates": False,
        "transient_peak_memory": "UNMEASURED",
        "original_model_activations": "NOT_CAPTURED",
        "full_model_training": False, "scientific_acceptance": "NOT_GRANTED",
    }
