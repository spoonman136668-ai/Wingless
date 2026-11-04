# CNC-G9F original-weight fidelity: preregistered development-only run

Qwen/Qwen3-30B-A3B original revision ad44e777bcd18fa416d9da3bd8f70d33ebb85d39.
Reuse G8/G9E SHA-256 single-original-shard custody; same first-complete
16-expert layer, five frozen seeds 11/23/37/53/71, 128x256 authentic BF16
weight windows and identical 32-row Gaussian activation probes. Frozen
4/8/16 expert prefixes. INT4 block64 FP32-scale control, prior G9E
peak-scaled INT2 block128/256, and new G9F fixed-normal-codebook
RMS-companded INT2 block128/256. The two INT2 formats must consume
**exactly equal physical packed+scale+metadata bytes** per tensor and group.
Report relative output/weight squared errors without altering codebook,
threshold, seed, window, expert or activation selection after observing
original weights. Preserve G9E NO_GO as an immutable negative.

No authentic model activations, real task/reasoning claims, GPU residency,
model forward, held-out science, scientific acceptance, promotion, or
production activation. A separate original-weight one-shot only after
exact-head Linux/Windows mandatory qualification. Fail closed on SHA
identity or original staging provenance mismatch; clean scratch and
preserve checksum-bound receipt.
