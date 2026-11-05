"""G9P: frozen full-original-expert CPU operator comparison.

One authenticated Qwen BF16 down projection, NOT original activations or
model inference. Matched synthetic probes for BF16, INT4, uniform INT2,
and companded INT2. No heldout, GPU, model, or promotion.
"""
import statistics
import time

BATCH = 8
SEEDS = (751, 811)
REPEATS = 3
RESULT = "G9P_ORIGINAL_FULL_EXPERT_OPERATOR_DEVELOPMENT_ONLY"


def compare_original_expert(weight, *, seeds=SEEDS):
    """Freeze selection/budgets; return each encoding, error and warm CPU timings.

    Weight must be original BF16 converted to FP32 after authenticated admission.
    Original inputs are not available: these synthetic Gaussian probes are
    development-only and cannot establish model accuracy or GPU speed.
    """
    import numpy as np
    from packed_int4 import pack_int4, apply_packed_int4, stored_bytes as bytes4
    from packed_int2 import pack_int2, apply_packed_int2, stored_bytes as bytes2
    from packed_companded_int2 import (
        pack_companded_int2, apply_companded_int2, stored_bytes as bytesc
    )
    from g9q_mixed_precision import pack_mixed, apply_mixed, stored_bytes as mixed_bytes
    w=np.asarray(weight, dtype=np.float32)
    if (w.ndim!=2 or w.shape[0]<4 or w.shape[1]<256
        or w.shape[1]%256 or w.size>4_194_304
        or not np.isfinite(w).all() or not np.any(w)):
        raise ValueError("invalid full original BF16 expert shape/data")
    if tuple(seeds)!=SEEDS:
        raise ValueError("development probe seeds must remain frozen")
    original_bytes=w.size*2
    arms=(
        ("original_bf16_float32_compute",None,None,None),
        ("packed_int4_block64",pack_int4,apply_packed_int4,bytes4),
        ("uniform_int2_block256",pack_int2,apply_packed_int2,bytes2),
        ("companded_int2_block256",pack_companded_int2,apply_companded_int2,bytesc),
    )
    state={}
    for name,packer,operator,sizer in arms:
        if packer is None:
            state[name]={"storage_bytes":int(original_bytes),"encode_seconds":0.0,
                         "operator":lambda x:x@w.T}
            continue
        kwargs={"block":64 if "int4" in name else 256}
        start=time.perf_counter()
        packed,scale=packer(w,**kwargs)
        encode=time.perf_counter()-start
        nbytes=sizer(packed,scale,w.shape[1],**kwargs)
        if not 0<nbytes<original_bytes:
            raise ValueError("compressed weight exceeds original storage")
        state[name]={"storage_bytes":int(nbytes),"encode_seconds":encode,
                     "operator":lambda x,p=packed,s=scale,k=kwargs,op=operator:
                         op(x,p,s,**k)}
    # The uniform-eighth arm is an equal-byte, weight-independent control.
    # Fractions/policies are frozen; FIT uses weights only, never probe inputs.
    mixed_arms=(
        ("mixed_uniform_eighth_INT4",8,"uniform"),
        ("mixed_ranked_eighth_INT4",8,"ranked"),
        ("mixed_ranked_quarter_INT4",4,"ranked"),
    )
    for name, denominator, policy in mixed_arms:
        start=time.perf_counter()
        representation=pack_mixed(w,denominator=denominator,policy=policy)
        encode=time.perf_counter()-start
        nbytes=mixed_bytes(representation,*w.shape)
        if not 0<nbytes<original_bytes:
            raise ValueError("mixed representation exceeds BF16 storage")
        state[name]={"storage_bytes":int(nbytes),"encode_seconds":encode,
                     "operator":lambda x,p=representation:apply_mixed(x,p)}
    results=[]
    for seed in SEEDS:
        rng=np.random.default_rng(seed)
        x=rng.standard_normal((BATCH,w.shape[1])).astype(np.float32)
        reference=x@w.T
        if not np.isfinite(reference).all() or not np.any(reference):
            raise ValueError("invalid original BF16 reference output")
        denominator=float(np.sum(reference.astype(np.float64)**2))
        if denominator<=0 or not np.isfinite(denominator):
            raise ValueError("invalid baseline energy")
        for name in tuple(x[0] for x in arms)+tuple(x[0] for x in mixed_arms):
            op=state[name]["operator"]
            warm=op(x)
            if warm.shape!=reference.shape or not np.isfinite(warm).all():
                raise ValueError("bad compressed operator output")
            timings=[]
            for _ in range(REPEATS):
                t0=time.perf_counter()
                out=op(x)
                timings.append(time.perf_counter()-t0)
            if not np.allclose(out,warm,rtol=1e-5,atol=1e-5):
                raise ValueError("non-deterministic compressed operator")
            error=float(np.sum((warm.astype(np.float64)-reference.astype(np.float64))**2)/denominator)
            if not np.isfinite(error) or error<0:
                raise ValueError("invalid normalized squared error")
            results.append({"seed":seed,"arm":name,
                "original_bf16_weight_bytes":int(original_bytes),
                "encoded_bytes_including_accounted_metadata":state[name]["storage_bytes"],
                "compression_ratio_vs_original_bf16":original_bytes/state[name]["storage_bytes"],
                "relative_output_squared_error":error,
                "encode_seconds":state[name]["encode_seconds"],
                "warm_cpu_operator_median_seconds":statistics.median(timings)})
    assert len(results)==len(SEEDS)*(len(arms)+len(mixed_arms))
    return {"schema":"wingless.gcen.g9p.one-original-full-expert-comparison.v1",
        "result":RESULT,"frozen_seeds":list(SEEDS),"batch":BATCH,
        "mixed_precision_selection":"WEIGHT_ONLY_ROW_MSE_VS_UNIFORM_ROW_CONTROL",
        "timing_repetitions_per_arm":REPEATS,"shape":list(w.shape),
        "comparison_count":len(results),"comparisons":results,
        "original_weights_authenticated_by_caller":True,
        "probe_inputs":"SYNTHETIC_GAUSSIAN_NOT_QWEN_ACTIVATIONS",
        "cpu_computation":"NUMPY_FLOAT32_NOT_NATIVE_BF16_EXECUTION",
        "timing_scope":"DEV_RUNNER_PYTHON_CPU_ONLY_NOT_COMPARATIVE_GPU",
        "full_model_forward":False,"gpu_vram_measured":False,
        "protected_heldout_used":False,"scientific_acceptance":"NOT_GRANTED"}
