# Wingless UP-8: confound-controlled observer ablation

Status: Windows-qualified observer-ablation result; not activated, promoted, or connected to ckb-plane.

Parent research record: UP-7 qualified as a negative result, then two post-qualification confounds were identified and recorded at `c9522ed665e4d1edfaa2d6c5001b00e15f9e3e47`.

## Why UP-8 exists

UP-7 cannot cleanly support a causal conclusion about direct transported-state readability because:

1. its 128/128 table-index-parity split accidentally withheld entity-3 values 1/3 from training and values 0/2 from held-out evaluation;
2. its single linear head concatenated a global state vector with an entity one-hot selector, which can alter only class bias terms and cannot make state-feature weights depend on the queried entity.

UP-8 corrects both before changing the transport or adding an internal clock.

## Balanced combination split

The full 256 memory tables are partitioned by:

`sum(entity values) mod 2`

This produces 128 train-pool and 128 held-out-pool tables while preserving every entity/value marginal in both pools.

For bounded runtime, UP-8 selects 64 tables from each pool through a deterministic permutation. The harness explicitly checks that every entity/value marginal remains represented.

## Legitimate entity conditioning

UP-8 trains four separate but architecture- and budget-matched heads, one for each queried entity.

This makes the question simply: can the transported representation for a table reveal each entity's value?

It does not test parameter sharing across entity queries.

## Feature ablation

### Magnitude only

Sixteen normalized coordinate probabilities.

This corresponds closely to the measurement-only information exposed in UP-7.

### Phase/coherence aware

The full normalized pure-state density representation:

`rho = psi psi^dagger`

Features contain:

- 16 diagonal probabilities;
- real and imaginary components of all 120 unique off-diagonal coherences;
- 256 real features total.

This preserves relative phase/coherence while being invariant to a global phase rotation.

## Frame/depth ablation

### Fixed frame

Observer training and held-out table evaluation both use propagation depth 128.

The table combinations are disjoint, but the transport frame is the same.

### Unseen depth

Training depths:

`8, 24, 72, 216, 432, 648`

Held-out depths:

`32, 128, 512, 1024`

No explicit depth or inverse is provided.

## Cells

UP-8 measures:

1. unitary magnitude-only, fixed frame;
2. unitary coherence-aware, fixed frame;
3. unitary magnitude-only, unseen depth;
4. unitary coherence-aware, unseen depth;
5. matched non-unitary coherence-aware, fixed frame;
6. matched non-unitary coherence-aware, unseen depth.

Every head receives 100 deterministic full-batch gradient steps at learning rate 1.0.

## Scientific interpretation

UP-8 reports separate flags rather than forcing the experiment to succeed:

- **measurement-loss supported** when fixed-frame coherence improves held-out accuracy over fixed-frame magnitude by at least 0.15;
- **fixed coherence pass** at >= 0.95 held-out accuracy;
- **cross-depth coherence pass** at >= 0.95;
- **frame ambiguity supported** when fixed-frame coherence passes but unseen-depth coherence does not.

If fixed coherence passes while unseen-depth coherence fails, the next experiment should add a co-evolving internal frame/clock rather than increasing generic observer size blindly.

If fixed coherence also fails, the observer representation itself remains insufficient.

## Harness qualification

The Windows harness fails only for build/regression/determinism/invariant defects.

The scientific flags are reported independently and do not control harness acceptance.

## Authority boundary

UP-8 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-22 against source head `d4ec92378407e46b2d38577a0ef096c2d0f526d0` passed the harness and complete Wingless regression.

The corrected balanced split remained valid with minimum entity/value marginal counts 14 in training and 13 in held-out data.

Observed held-out accuracies:

- unitary magnitude-only, fixed frame: 0.62109375;
- unitary coherence-aware, fixed frame: 1.0;
- unitary magnitude-only, unseen depth: 0.294921875;
- unitary coherence-aware, unseen depth: 0.373046875;
- matched-control coherence-aware, fixed frame: 1.0;
- matched-control coherence-aware, unseen depth: 0.3896484375.

The fixed-frame coherence advantage over magnitude-only observation was 0.37890625.

Interpretation: the corrected ablation supports two distinct effects. First, relative phase/coherence contains memory information lost by magnitude-only measurement. Second, even coherence-aware observation remains frame-dependent and does not generalize across unseen propagation depths without additional frame information.
