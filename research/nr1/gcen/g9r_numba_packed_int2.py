"""G9R: opt-in Numba direct packed-companded INT2 CPU operator.
No model forward, native activations, GPU, or acceptance. JIT is not part of routine CI.
"""
import statistics
import time
import numpy as np
from packed_companded_int2 import pack_companded_int2,apply_companded_int2,_validate

SEEDS=(751,811)
REPEATS=7
BATCH=8
THREADS=2


def relative_squared_error(observed, reference):
    """Actual operator output error, never weight-only surrogate quality."""
    y=np.asarray(observed,dtype=np.float64)
    ref=np.asarray(reference,dtype=np.float64)
    if y.ndim!=2 or y.shape!=ref.shape or not np.isfinite(y).all() or not np.isfinite(ref).all():
        raise ValueError("invalid matched operator outputs")
    energy=float(np.sum(ref*ref))
    if not np.isfinite(energy) or energy<=0:
        raise ValueError("invalid dense reference energy")
    error=float(np.sum((y-ref)**2)/energy)
    if not np.isfinite(error) or error<0:
        raise ValueError("invalid relative operator output error")
    return error

def validate(x,d,s):
    x=np.asarray(x,dtype=np.float32)
    if x.ndim!=2 or x.shape[0]!=BATCH or x.shape[1]%256 or not np.isfinite(x).all():
        raise ValueError("invalid frozen G9R input")
    _validate(d,s,x.shape[1],256)
    return x

def build_operator():
    from numba import njit,prange
    @njit(parallel=True)
    def direct(x,d,s):
        batch,cols=x.shape
        rows=d.shape[0]
        out=np.zeros((batch,rows),dtype=np.float32)
        for row in prange(rows):
            for n in range(batch):
                total=np.float32(0)
                for block in range(cols//256):
                    part=np.float32(0)
                    for j in range(64):
                        byte=d[row,block*64+j]
                        for q in range(4):
                            code=(byte>>(2*q))&3
                            level=np.float32(-1.5104 if code==0 else
                                (-0.4528 if code==1 else
                                 (0.4528 if code==2 else 1.5104)))
                            part+=x[n,block*256+4*j+q]*level
                    total+=part*np.float32(s[row,block])
                out[n,row]=total
        return out
    return direct

def benchmark_original(weight):
    from numba import set_num_threads,get_num_threads
    w=np.asarray(weight,dtype=np.float32)
    if w.shape!=(2048,768) or not np.isfinite(w).all():
        raise ValueError("G9R requires frozen original expert shape and finite weights")
    data,scales=pack_companded_int2(w,block=256)
    # Numba nopython cannot accept float16 array arguments. Preserve encoded
    # FP16 storage; promote only the small runtime scale table to FP32.
    jit_scales=scales.astype(np.float32)
    set_num_threads(THREADS)
    if get_num_threads()!=THREADS:
        raise RuntimeError("incorrect benchmark thread count")
    direct=build_operator()
    rows=[]
    compile_seconds=0.
    for seed in SEEDS:
        x=validate(np.random.default_rng(seed).standard_normal((BATCH,768)).astype(np.float32),data,scales)
        t0=time.perf_counter()
        actual=direct(x,data,jit_scales)
        if seed==SEEDS[0]:
            compile_seconds=time.perf_counter()-t0
        expected=apply_companded_int2(x,data,scales,block=256)
        # Frozen matched-input dense FP32 reference. This is the same
        # BF16-derived original weight source, but CPU FP32 computation.
        dense_reference=x@w.T
        packed_vs_dense_error=relative_squared_error(actual,dense_reference)
        numpy_vs_dense_error=relative_squared_error(expected,dense_reference)
        if abs(packed_vs_dense_error-numpy_vs_dense_error)>2e-5:
            raise RuntimeError("G9R dense-reference fidelity drift")
        if not np.allclose(actual,expected,rtol=2e-5,atol=7e-5):
            raise RuntimeError("G9R direct packed INT2 parity failure")
        timings={}
        for name,fn in (("dense_fp32",lambda:x@w.T),
                        ("numpy",lambda:apply_companded_int2(x,data,scales,block=256)),
                        ("numba",lambda:direct(x,data,jit_scales))):
            measures=[]
            for _ in range(REPEATS):
                start=time.perf_counter()
                val=fn()
                measures.append(time.perf_counter()-start)
                compare=dense_reference if name=="dense_fp32" else expected
                if not np.allclose(val,compare,rtol=2e-5,atol=7e-5):
                    raise RuntimeError("G9R repeated-output parity failure")
            timings[name]=statistics.median(measures)
        rows.append({"seed":seed,
                     "dense_fp32_median_seconds":timings["dense_fp32"],
                     "numpy_median_seconds":timings["numpy"],
                     "numba_median_seconds":timings["numba"],
                     "dense_fp32_speedup_vs_numba":timings["numba"]/timings["dense_fp32"],
                     "numba_speedup_vs_dense_fp32":timings["dense_fp32"]/timings["numba"],
                     "relative_output_squared_error_vs_dense_fp32":packed_vs_dense_error,
                     "speedup":timings["numpy"]/timings["numba"],
                     "max_abs_delta":float(np.max(np.abs(actual-expected)))})
    return {"schema":"wingless.gcen.g9r.original-packed-int2-numba-cpu.v1",
            "seeds":list(SEEDS),"batch":BATCH,"threads":THREADS,"repeats":REPEATS,
            "jit_compile_and_first_call_seconds":compile_seconds,
            "runtime_fp32_scale_bytes":int(jit_scales.nbytes),
            "original_bf16_weight_bytes":int(w.size*2),
            "packed_int2_weight_bytes":int(data.nbytes+scales.nbytes+24),
            "packed_runtime_scale_representation":"FP32_EPHEMERAL_TABLE",
            "dense_reference":"ORIGINAL_BF16_WEIGHT_CONVERTED_TO_FP32_CPU_OPERATOR",
            "benchmark_blas_threads":"OPENBLAS_NUM_THREADS_1_CI_ONLY",
            "comparisons":rows,"source_weights_authenticated_by_caller":True,
            "uses_original_activations":False,"uses_gpu":False,
            "scientific_acceptance":"NOT_GRANTED"}
