# G9U — GPU operator experiment boundary

This optional module tests the existing companded INT2 expert representation
against a matched dense FP32 CUDA reference, **only when an authorized isolated
NVIDIA GPU runner already exists**. No new runner, workflow dispatch, controller,
checkpoint download, model activation or deployment is created.

Input custody: the caller must authenticate and load a full original pinned
BF16 expert tensor before invoking `run_authenticated_expert`. The module
does not verify source-shard SHA256 itself; claiming authentic weights from
unverified caller data is prohibited. The exact repository SHA and current
CUDA device/VRAM are checked by G9T in the same process.

The packed operator uses **GPU tiled unpacking then dense matrix multiply**,
not a fused CUDA packed kernel. This is a feasibility baseline, not a claim
that packed INT2 is faster or memory-resident end to end.

Protocol: layer-0 expert-0 down projection, 2048×768; batch 8;
frozen seeds 751/811; 3 warmups + 7 CUDA-event timings per arm per seed;
NumPy packed codec with block 256; compare operator squared output error.
Peak CUDA allocator observation includes concurrent dense and packed
representations, so it **cannot** prove full-model fit or isolated packed
footprint. The measured peak is not device-wide total VRAM.

Admission requirements remain: current independently verified runner
identity, isolation, GPU use authority, original-weight custody, and
source-pinned dependencies. Existing GitHub native CI is CPU-only and
must never be treated as execution evidence.

No original activations, language quality, full-model inference or
scientific acceptance is established.
