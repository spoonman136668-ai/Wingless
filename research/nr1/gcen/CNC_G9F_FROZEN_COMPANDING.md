# CNC-G9F — fixed-companded nonuniform INT2 screen (frozen)

Motivation: G9E original-weight Qwen window result at
`c7eab82ca6a609e6ffdf4c61b90d6ceccf09a4d6` shows
~35.65%/42.65% relative output squared error for the peak-scaled
midrise INT2 block128/256, versus ~1.24% for G8 INT4 block64.
Preserve G9E negative and authentication receipts unchanged.

The G9F hypothesis is that using a priori standard-normal Lloyd-Max
reconstruction points `[-1.5104,-0.4528,0.4528,1.5104]` and
frozen per-row/block RMS with FP16 scaling will improve Gaussian
weight fidelity at the **same packed bytes** as G9C 2-bit codes.
No Qwen-specific calibration, empirical codebook fitting, or training.
Use blocks 128,256, threshold `±0.9816` and 0, fixed before
any authentic-source trial. Symmetric extremes stay inside 2 bits.

G9F first-stage code qualification: deterministic synthetic seed 23,
weight window 128×256, input batch32; compare BF16 direct output,
G9C uniform INT2 same block and bytes, and G8 INT4 at higher bytes.
Test malformed shapes, NaN/Inf, FP16 scale limits, zero blocks,
exact packed codes and blockwise direct-operator behavior.
This is a development mechanistic screen, **not** scientific acceptance.

G9F original Qwen expert-weight fidelity is UNVERIFIED until a separate
original SHA256-custodied one-shot run, with frozen G8/G9E seeds and
selection and no result-dependent retuning. Authentic Qwen activations,
task quality, GPU speed, actual 12 GiB residency, and model inference
also UNVERIFIED. No production activation or protected-main mutation.
