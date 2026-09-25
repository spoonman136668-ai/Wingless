# Wingless UP-SQ0 — serialized pre-language sequence qualification gate

Status: preregistered design only; not dispatched while A/B/C lane qualifications are active.

Purpose: first serialized integration gate before promoting A/B/C mechanisms toward UP-LM-0. This design is frozen independently of the unobserved outcomes of UP-75A, UP-61B margin calibration, and UP-69C.

Evidence basis: sealed Wingless A/B/C results through UP-74A, UP-61B write-scale, and UP-68C, plus the September 24, 2026 Deep Research review supplied in the research conversation. No currently running experiment is modified by this document.

## Scientific question

Can a bounded Wingless sequence machine preserve exact information, update evolving state, maintain role/filler structure, and expose useful endogenous confidence under workloads that are known to separate simple recurrent persistence from language-relevant sequence memory?

## Frozen task families

All families use deterministic generators and structural train/test splits. Random example-level train/test shuffling is prohibited.

1. **Delayed copy**
   - alphabet size: 16
   - training sequence lengths: 16, 32
   - held-out lengths: 64, 128
   - delays: 4, 8, 16, 32
   - metric: exact token accuracy and whole-sequence exact accuracy.

2. **Selective copy**
   - alphabet size: 16
   - 25% of presented symbols marked for later recall
   - training lengths: 32, 64
   - held-out lengths: 128, 256
   - metric: selected-symbol exact accuracy and ordered-sequence exact accuracy.

3. **Multi-query associative recall**
   - key vocabulary: 32
   - value vocabulary: 32
   - loads: 4, 8, 16 associations
   - queries per episode: 4
   - collision arm: repeated keys with latest-value-wins semantics
   - metric: exact query accuracy, stratified by load and collision status.

4. **Overwrite / latest-value-wins**
   - 8 logical slots
   - 16-value alphabet
   - writes per episode: 8, 16, 32
   - exactly half of episodes contain at least one overwrite
   - metric: exact final-slot accuracy and overwrite-only accuracy.

5. **State-machine composition**
   - 8 latent states
   - 6 input operators, each a fixed permutation
   - training composition depths: 1–4
   - held-out depths: 8, 16
   - held-out operator compositions are absent from training
   - metric: exact final-state accuracy.

6. **Role/filler recombination**
   - 5 roles
   - 3 fillers per role
   - training tuples: the same 54-state structural partition used by the compact factorized A family
   - held-out tuples: the same 162-state complement
   - inputs are presented sequentially as role/filler assignments
   - metric: per-role and whole-tuple exact accuracy.

7. **Multi-bank interference**
   - bank counts: 2, 4, 6
   - 4 entities per bank
   - load levels: 25%, 50%, 100% of available bank/entity slots
   - metric: value accuracy, exact-episode accuracy, and cross-bank corruption rate.

## Frozen sequence lengths and extrapolation

Every applicable family reports in-distribution performance and 2×/4× held-out length performance. No context-length claim is permitted from propagation depth alone.

## Frozen integration arms

The gate must compare these preregistered arms under the same generated episodes:

- **transport-only**: persistent reversible/unitary state with no confidence-gated correction and no exact-recall side channel;
- **transport + gated correction**: same carrier with an endogenous confidence/margin-controlled correction; no oracle state;
- **transport + gated correction + bounded exact recall**: same recurrent state plus a separately bounded exact-recall channel.

The exact-recall channel is capped at 16 key/value entries in this gate. Full global self-attention is prohibited.

## Controls

- no oracle state;
- no adaptive threshold selection;
- no result-informed retries;
- no changing train/test generators after execution;
- no hidden increase in recurrent-state dimension between arms;
- no production authority or live activation;
- no CKB, ckb-plane, KTRADE, Nemotron, broker, accepted-ref, or competing-scheduler dependency.

## Seeds

Generator seed bases are frozen at 120M, 121M, and 122M for every arm. Any infrastructure repair must rerun the same frozen experiment.

## Measurements

For every arm and task family record:

- task accuracy metrics above;
- peak recurrent-state bytes;
- exact-recall bytes;
- peak process RAM;
- peak VRAM when a GPU path exists;
- batch-one examples/tokens per second;
- wall-clock qualification time.

## Interpretation

This is a diagnostic gate, not a promotion guarantee. Scientific negatives are sealed exactly as observed. No arm is selected by aggregate score after execution. Mechanism promotion requires a later preregistered rule using this sealed evidence.

UP-SQ0 must execute serially, not concurrently with an A/B/C qualification, unless a future preregistration explicitly changes that scheduling rule before execution.
