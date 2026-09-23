# Wingless UP-36: task-weighted continuous symmetry fusion

Status: Windows-qualified scientific negative; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-35 qualified scientific positive sealed at `09284aa70203adf8881be42f67084c49b338dabe`.

## Question

UP-35 showed that useful repeated-spectrum symmetry can be constructed incrementally from the minimum-capacity singleton state through training-only pairwise hard merges.

UP-36 removes the hard merge search.

It asks:

**Can six independently parameterized spectral directions move continuously under task-derived attraction and spontaneously form exact repeated-spectrum groups without evaluating discrete merge candidates or receiving a finished partition menu?**

## Starting point

The transport begins from the same six distinct singleton offsets used in the multiplicity-ablation line:

- maximum multiplicity 1;
- commutant capacity 6;
- zero mean;
- RMS offset magnitude 0.05 radians per step.

No partition labels or finished symmetry organization are supplied.

## Training-only pair affinities

For each of the 15 offset pairs, UP-36 performs a continuous probe:

- move both offsets 50% toward their common midpoint;
- do not make them equal;
- renormalize the six offsets to zero mean and RMS 0.05;
- evaluate the result only on the balanced training-side fit/validation split.

The probe remains multiplicity 1. It is not a merge.

Its attraction weight is the positive validation-score gain over the singleton baseline. Negative or zero gains receive zero attraction.

The validation score remains the UP-34/UP-35 resource-aware objective:

`score = min(validation held, validation commit) - 0.02 * capacity / 36`.

Since all pair probes remain capacity 6, their weight is determined only by task benefit.

## Continuous fusion flow

All positive pair attractions are normalized into one symmetric weight matrix.

The six offsets then evolve simultaneously:

`delta_i <- delta_i + 0.45 * sum_j w_ij (delta_j - delta_i)`

After every flow step:

- offsets are re-centered to zero mean;
- RMS is restored to 0.05;
- already fused coordinates remain fused;
- any pair whose gap falls below 0.0025 radians is fused exactly and remains sticky.

The flow runs for at most 10 fixed steps.

No pair is selected as a discrete architecture mutation. Fusion is triggered only by the continuous trajectory crossing the preregistered tolerance.

## Model selection

Every flow checkpoint is evaluated only on the training-side validation split.

Checkpoint selection uses the same task-minus-capacity score. Exact fusions increase commutant capacity and therefore incur the same frozen capacity price as earlier experiments.

The best checkpoint is frozen before any true held-out evaluation.

## Final evaluation

After the continuous trajectory and checkpoint choice are complete:

1. the selected continuous-fusion structure is retrained on the full original training pool;
2. it is evaluated once on the untouched 128-table true held-out pool;
3. the capacity-36 repeated-spectrum control is evaluated under the same final protocol.

## Scientific gates

Preflight:

- balanced fit/validation split;
- no true held-out selection;
- at least one task-positive continuous pair affinity.

Continuous formation:

- selected flow step >= 1;
- selected commutant capacity > 6;
- selected commutant capacity < 36.

Final capability:

- held-out accuracy >= 0.985;
- mutable commit accuracy >= 0.95;
- exact final table accuracy >= 0.90;
- relational query accuracy >= 0.95.

Retention versus full capacity:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

## Interpretation boundary

A positive result would be stronger than UP-35.

UP-35 supplied a discrete merge operation and explicitly evaluated merge candidates.

UP-36 supplies only continuously parameterized spectral directions, task-derived attraction weights, and a global fusion flow. Exact degeneracies arise only when the continuous dynamics drive offsets within the fixed fusion tolerance.

The experiment still contains important scaffolding:

- the six spectral directions are predefined;
- the fusion flow equation and tolerance are hand-designed;
- pairwise attraction is estimated by finite continuous probes.

A later positive-result successor should remove the hand-designed flow and learn the continuous transport parameters directly under a differentiable task/resource objective.

A negative result would show that the discrete developmental operator in UP-35 was doing essential search work that the continuous flow cannot yet replace.

UP-29 already established exact real-orthogonal equivalence; no result here supports uniquely complex or quantum-computation claims.

## Plain speak

UP-35 let the system decide which walls to knock down.

UP-36 removes the hammer.

Every wall can now slide continuously. The task measures whether moving any two walls closer helps, and those local signals create a simultaneous attraction field.

If some walls naturally collapse onto the same location and the resulting structure still works perfectly on unseen data, that is a stronger form of self-organization than choosing explicit merges.

## Authority boundary

UP-36 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35913964677` completed successfully at source head `c5dcbd283d27c1757102f2f5359658b23bd3c73a`.

The continuous process generated 14 positive pair affinities and created an exact fusion without evaluating hard merge candidates. The best validation checkpoint was flow step 1, where channels 1 and 2 fused and commutant capacity increased from 6 to 8.

That checkpoint did not meet the scientific performance gates. After retraining on the complete original training pool, true held-out accuracy was 0.974609375, mutable-commit accuracy was 0.8502604166666666, exact-final-table accuracy was 0.7916666666666666, and relational-query accuracy was 0.8333333333333334. The full-capacity control remained perfect.

Interpretation: continuous symmetry formation itself is viable, but the fixed affinity field estimated only at the singleton starting geometry does not supply enough adaptive guidance. The next experiment should recompute task-derived continuous affinities after each flow step rather than freezing them at initialization.
