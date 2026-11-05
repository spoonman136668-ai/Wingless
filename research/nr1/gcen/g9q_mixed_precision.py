"""G9Q bounded mixed companded-INT2/INT4 expert CPU representation.

Weight-only development selector chooses high-error rows for INT4; frozen
uniform-row selection is an equal-byte negative control. Inference uses
packed blocks only, never reconstructs the full dense weight matrix.
NOT original Qwen activations, GPU inference or quality acceptance.
"""
import numpy as np
from packed_companded_int2 import (
    LEVELS, pack_companded_int2, apply_companded_int2,
    stored_bytes as companded_bytes,
)
from packed_int4 import (
    pack_int4, apply_packed_int4, stored_bytes as int4_bytes,
)

POLICIES = ("ranked", "uniform")
DENOMINATORS = (4, 8)
BLOCK = 256


def _row_error(weight, data, scales):
    rows, cols = weight.shape
    score = np.zeros(rows, dtype=np.float64)
    for k in range(cols // BLOCK):
        encoded = data[:, k * BLOCK // 4:(k + 1) * BLOCK // 4]
        codes = np.empty((rows, BLOCK), dtype=np.uint8)
        for shift in range(4):
            codes[:, shift::4] = (encoded >> (2 * shift)) & 3
        restored = LEVELS[codes] * scales[:, k, None].astype(np.float32)
        diff = weight[:, k * BLOCK:(k + 1) * BLOCK].astype(np.float64) - restored
        score += np.einsum("ij,ij->i", diff, diff)
    if not np.isfinite(score).all():
        raise ValueError("nonfinite source row reconstruction error")
    return score


def pack_mixed(weight, *, denominator=8, policy="ranked"):
    if type(denominator) is not int or denominator not in DENOMINATORS or policy not in POLICIES:
        raise ValueError("nonfrozen mixed precision policy")
    w = np.asarray(weight, dtype=np.float32)
    if (w.ndim != 2 or w.shape[0] < 8 or w.shape[0] > 65535
        or w.shape[1] <= 0 or w.shape[1] % BLOCK
        or not np.isfinite(w).all()):
        raise ValueError("invalid mixed precision source matrix")
    rows, cols = w.shape
    # Fit selection to SOURCE WEIGHTS alone, never the synthetic evaluation inputs.
    # The initial full INT2 encoding is encoder-only and is discarded afterward.
    first_data, first_scales = pack_companded_int2(w, block=BLOCK)
    n_high = max(1, rows // denominator)
    if policy == "ranked":
        score = _row_error(w, first_data, first_scales)
        selected = np.lexsort((np.arange(rows), -score))[:n_high]
    else:
        selected = np.arange(0, rows, denominator)[:n_high]
    high = np.sort(selected.astype(np.uint16))
    low = np.setdiff1d(np.arange(rows, dtype=np.uint16), high, assume_unique=True)
    if not len(high) or not len(low):
        raise ValueError("empty mixed precision arm")
    hi_data, hi_scale = pack_int4(w[high], block=64)
    lo_data, lo_scale = pack_companded_int2(w[low], block=BLOCK)
    state = (high, low, hi_data, hi_scale, lo_data, lo_scale)
    validate(state, rows, cols)
    return state


def validate(state, rows, cols):
    if (not isinstance(state, tuple) or len(state) != 6
        or type(rows) is not int or not 8 <= rows <= 65535
        or type(cols) is not int or cols <= 0 or cols % BLOCK):
        raise ValueError("invalid mixed operator dimensions/state")
    hi, lo, hd, hs, ld, ls = state
    if (not isinstance(hi, np.ndarray) or not isinstance(lo, np.ndarray)
        or hi.dtype != np.uint16 or lo.dtype != np.uint16
        or hi.ndim != 1 or lo.ndim != 1
        or not len(hi) or not len(lo)
        or not np.array_equal(np.sort(np.concatenate((hi, lo))),
                              np.arange(rows, dtype=np.uint16))):
        raise ValueError("mixed row partition is incomplete or duplicated")
    if hd.shape[0] != len(hi) or ld.shape[0] != len(lo):
        raise ValueError("mixed packed row inventory mismatch")
    int4_bytes(hd, hs, cols, block=64)
    companded_bytes(ld, ls, cols, block=BLOCK)


def apply_mixed(inputs, state):
    x = np.asarray(inputs, dtype=np.float32)
    if x.ndim != 2 or not x.shape[0] or not x.shape[1] or not np.isfinite(x).all():
        raise ValueError("invalid mixed activation")
    hi, lo, hd, hs, ld, ls = state
    rows = len(hi) + len(lo)
    validate(state, rows, x.shape[1])
    result = np.empty((len(x), rows), dtype=np.float32)
    result[:, hi] = apply_packed_int4(x, hd, hs, block=64)
    result[:, lo] = apply_companded_int2(x, ld, ls, block=BLOCK)
    if not np.isfinite(result).all():
        raise ValueError("nonfinite mixed-precision operator result")
    return result


def stored_bytes(state, rows, cols):
    validate(state, rows, cols)
    hi, lo, hd, hs, ld, ls = state
    return int(hi.nbytes + lo.nbytes +
               int4_bytes(hd, hs, cols, block=64) +
               companded_bytes(ld, ls, cols, block=BLOCK) + 24)
