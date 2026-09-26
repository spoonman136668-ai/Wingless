# Wingless UP-21: anonymous phase-code bottleneck

Status: Windows-qualified scientific pass; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-20 qualified partial positive sealed at `bd735f81c8586a462bc9992f7de37da52563f45d`.

## Question

UP-20 proved that the anonymous quadratic representation is learnable without semantic demixing, but generic one-hot ridge generalization remained too weak for recurrent mutable use.

UP-21 asks whether the readout should exploit the circular value code already learned in UP-13 instead of learning four unrelated class surfaces.

## Learned target

For each value, UP-21 uses the frozen UP-13 learned phase alphabet:

`[0, -1.7700278912430298, 1.1656136621055615, 2.5917033093902164]`.

Each class is converted to a two-dimensional unit-circle target:

`[cos(phi), sin(phi)]`.

This is not a hand-designed new code. It reuses a value representation learned earlier in the same research lineage.

## Readout

The recurrent state remains:

- one 96-dimensional complex state;
- six anonymous 16-dimensional blocks;
- dense hidden channel mixing;
- no semantic channel locations;
- no hidden mixer or inverse exposed to learning.

The observer uses the same complete 702-feature quadratic anonymous Gram map from UP-19/UP-20.

For each entity:

1. deterministic ridge regression maps 702 features to two learned phase coordinates;
2. a fresh two-feature linear softmax head calibrates those coordinates into four value probabilities.

The ridge solve is converted to primal weights. Training examples are not retained at runtime.

## Training

- all 128 balanced training tables;
- four noisy trials per table;
- memory noise 0.05;
- memory-only global-phase nuisance;
- ridge lambda 1e-6;
- fresh 2D softmax calibration head per entity.

No oracle witness or hidden mixer is used for learning.

## Evaluation

- all 128 disjoint held-out tables;
- unseen depths 32, 128, 512, 1024;
- two noisy trials per table/depth;
- matched non-unitary path uses the exact same learned model.

UP-21 also reports mean cosine alignment between the predicted two-dimensional bottleneck and the learned UP-13 phase target.

## Mutable integration

The same phase-code model is carried into:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relation head.

## Scientific gates

**Phase-code learning**

- canonical training accuracy >= 0.99.

**Unseen depth**

- unitary held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-12.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-12.

## Interpretation boundary

A positive result means anonymous-state readout becomes sufficiently robust when the learner is constrained to the task-relevant two-dimensional phase code. That would support removing the historical semantic channel decomposition entirely from readout.

A negative result means that reusing the learned value phase alphabet is still insufficient to produce recurrently stable anonymous-state decoding. The next experiment should target calibration/commit confidence or a learned low-rank factorization directly, not enlarge the polynomial feature space.

No result here establishes phase-specific or unitary superiority. A real orthogonal transport comparator remains required before stronger causal claims.

## Plain speak

UP-20 proved we can read almost everything from the scrambled state, but one part is still too shaky.

UP-21 stops asking the reader to learn four unrelated answers. It teaches the reader to first recover the two-dimensional “value compass” Wingless already learned earlier, then decide which value that direction means.

If it works, we keep the state scrambled and never have to rebuild the old labeled drawers.

## Authority boundary

UP-21 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35852739971` completed successfully at source head `6b150300b2782625e11ad97136861dece6ebac37`.

The anonymous phase-code bottleneck reached 1.0 train and 1.0 held-out accuracy across every entity at unseen depths 32, 128, 512, and 1024. Mean cosine alignment with the learned UP-13 phase target was 0.9861381590317367.

The exact same learned regressors and calibration heads on the matched non-unitary transport reached 0.59765625 held-out accuracy while its norm drift reached 2211700831922.4893.

Mutable integration used 48 scenarios × 16 writes and achieved 1.0 commit decode, 1.0 exact final-table, and 1.0 relational-query accuracy. Minimum value margin was 0.5043631193460631 and minimum relation margin was 0.9014937818220882.

Interpretation: historical semantic channel recovery is not required. A low-rank task-relevant phase bottleneck can decode the anonymously mixed recurrent state directly and remain exact under deep unitary transport and repeated irreversible commits. The remaining explicit structural assumption is the six anonymous 16-dimensional block factorization.
