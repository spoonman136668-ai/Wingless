# Wingless recurrent JEPA-style memory pilot — isolated DEV

**Parent source:** `9a92fcae40c811c724b4214da03fff34126f6b4f`.
**Motivation:** Balestriero & LeCun (2026), https://arxiv.org/abs/2610.39302. Good-view latent training motivates this source-only temporal adaptation; their theorem does **not** prove recurrent memory or native causal reasoning.

## Implementation

`stream.go` is an independent, two-weight Go recurrent student: `state_t = alpha*state_(t-1) + observation_t`; a learned predictor maps it to a fixed, noncollapsed independent target observation from the same underlying sequence. Backpropagation through the bounded recurrence learns `alpha` and the predictor weight. A full-gradient-steps memoryless control (`alpha=0`), shuffled/incorrectly paired target control, strict bounded input/budget checks, and SHA-256 versioned cold checkpoint are included. No external deps, token model, LLM calls, owner routes, queues, production integration or scientific first-use.

`stream_test.go` demonstrates held-out development learning vs a competent memoryless learner; it also verifies bad views, exact cold restoration and tamper resistance. Tests are NOT scientific first-use, have no sealed world-family splits and cannot establish native intelligence or resource efficiency versus strong recurrent controls.

## Future scientific acceptance (NOT executed here)

Freeze independent unfamiliar temporal world families and paired observations, train/validation/sealed-heldout, original baseline and fully trained RNN/sequence-model controls, equal actions/information/optimizer/compute, cold memory/retention, action-consequence transfer, calibration, interference, whole-process wall/CPU/RSS, multiple independent seeds and explicit pre-result thresholds. Require native Windows/Linux full regression plus research owner acceptance before a scientific claim. The ongoing Wingless UP research and model inference infrastructure are unchanged.

Reproduction: `go test ./research/jepa -count=1`.
