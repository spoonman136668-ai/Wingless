# Wingless UP-26: dynamics-discovered observable breadth

Status: Windows-qualified strong partial positive with mutable commit gate negative; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-25 qualified partial positive sealed at `6dcbec3c567e2dc39f4052948925a62c758377fb`.

## Question

UP-25 proved that dynamics-only observable discovery works:

- the discovered operators approximately commute with the unitary transport;
- their features remain depth-stable;
- the task learner fits the training set above 99%.

But a bank of 32 generic discovered observables reached only 92.48% held-out accuracy and failed recurrent mutable integration.

UP-26 asks one narrow question:

**Was UP-25 limited primarily by discovered observable-bank breadth?**

## Frozen comparison

UP-26 discovers a single deterministic bank of 64 observables using the exact UP-25 procedure.

The first 32 observables form the baseline arm.

This guarantees:

- identical seed construction;
- identical transport projection;
- identical projection rounds;
- identical ordering;
- identical training data;
- identical held-out data;
- identical phase-code target;
- identical ridge lambda;
- identical classifier training.

Only the number of available discovered observables differs.

## Discovery

Discovery remains:

- full 96-dimensional latent transport only;
- no hidden six-channel multiplicity;
- no hidden full-state mixer;
- no semantic labels;
- 20 binary conjugation-averaging rounds;
- effective orbit size 1,048,576;
- offline transport adjoint allowed only during discovery;
- no runtime adjoint or unmix.

## Arms

**Baseline 32**

- 32 complex observable expectations;
- 64 real raw features;
- complete degree-2 lift: 2,144 features.

This arm must reproduce the UP-25 held-out accuracy of 0.9248046875 to numerical tolerance.

**Primary 64**

- 64 complex observable expectations;
- 128 real raw features;
- complete degree-2 lift: 8,384 features.

All other learning settings remain frozen.

## Training

Both arms train only at depth 0 using:

- all 128 balanced training tables;
- two deterministic noisy trials per table;
- memory noise 0.05;
- memory-only global-phase nuisance;
- ridge lambda 1e-6;
- UP-13 learned 2-D phase-code supervision;
- fresh 2-D softmax calibration heads.

## Held-out evaluation

Both arms use:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two deterministic noisy trials per table/depth.

The matched non-unitary control is evaluated only for the primary 64-observable model using the exact same learned regressors and classifiers.

## Mutable integration

The primary 64-observable model runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relational-query head.

## Scientific gates

**Baseline reproduction**

- absolute held-out accuracy difference from UP-25 <= 1e-12.

**Discovery commutator**

- maximum entrywise commutator error <= 1e-5.

**Feature invariance**

- maximum degree-2 feature drift at depth 1024 <= 5e-3.

**Primary phase-code learning**

- depth-0 training accuracy >= 0.99.

**Primary unseen depth**

- held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-10.

**Primary mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-10.

## Interpretation boundary

A positive primary result supports the hypothesis that UP-25 was coverage-limited: the discovery mechanism was sound, but 32 generic commutant directions did not span enough task-relevant stable structure.

A negative result with a reproduced 32-observable baseline means simply doubling generic breadth is insufficient. The next experiment should change *selection or organization* of discovered observables rather than continue increasing the bank blindly.

No post-result change to lambda, data volume, rounds, or thresholds is authorized.

No result establishes phase-specific or unitary superiority. A real orthogonal comparator remains required before stronger causal claims.

## Plain speak

UP-25 found the right kind of stable internal measuring tools, but perhaps not enough of them.

UP-26 repeats the exact same discovery process and simply keeps twice as many.

The first half must behave exactly like UP-25. If the full bank closes the gap, we know the discovery method was right and the earlier bank was just too narrow.

## Authority boundary

UP-26 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35870174317` completed successfully at source head `663ae33f9ef84d425e9a9da687ae3281af266113`.

The 32-observable control reproduced UP-25 exactly: held-out accuracy was 0.9248046875 with zero delta.

The 64-observable primary arm reached 1.0 training accuracy and 0.998779296875 held-out accuracy, so the preregistered unseen-depth gate passed. Per-entity held-out accuracies were 0.9970703125, 1.0, 1.0, and 0.998046875.

Mutable performance improved sharply but remained just below the preregistered commit threshold: 756/768 commits were correct (0.984375), exact-final accuracy was 0.9583333333333334, and relational-query accuracy was 1.0.

Interpretation: dynamics-discovered observable breadth is a real capacity variable, but simply retaining more generic directions should not be extended blindly. UP-27 keeps runtime width fixed at 64 and tests whether training-only task relevance can choose a better 64-observable subset from a 128-observable discovered candidate bank.
