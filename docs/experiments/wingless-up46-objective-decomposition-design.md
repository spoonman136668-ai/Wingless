# Wingless UP-46 design — UP-45 objective-decomposition audit

Status: staged only; do not launch until the seed-matched UP-44 FIXA control is sealed and the Branch-C continuation remains scientifically supported.

Parent: UP-45 seal `48fb886f517db0cb21fdbc757cd7bc9500ed9f59`.

Frozen decision-tree branch: C.

## Question

UP-45's reversible proximal operator created exact fusion at step 29, but the unchanged task objective selected non-fused step 24.

Which component of the frozen objective makes step 29 less preferred, and does exact equality change hard recurrent capability relative to the task-selected state?

## Frozen state pair

No optimization occurs in UP-46.

The two states come directly from the authoritative UP-45 trajectory:

- selected step 24:
  `[-0.05841298929708432, -0.03887820934987345, 0.0051658188506906524, -0.03702321187042985, 0.07759973733795059, 0.05154885432874637]`;
- exact-fused step 29:
  `[-0.045664698831175514, -0.04501676812740249, 0.007218659931181362, -0.045664698831175514, 0.07894943527059196, 0.05017807058798021]`.

Step 29 is the only exact-fused state in the UP-45 trajectory, so no post-result choice among multiple fused states is required.

## Measurements

Using the unchanged UP-45 evaluator, measure for both states:

- normalized phase score;
- mean correct value probability;
- mean correct relation probability;
- harmonic task score;
- soft capacity;
- resource penalty;
- final smooth objective.

Also evaluate the unchanged hard recurrence and report:

- held-out accuracy;
- commit accuracy;
- exact final-table accuracy;
- relation accuracy.

Report fused-minus-selected deltas.

## Boundaries

- no optimizer;
- no lambda change;
- no stronger fusion pressure;
- no checkpoint search;
- no held-out selection;
- no target geometry;
- no activation or external integration.

The hard held-out comparison is diagnostic only and cannot choose or tune a subsequent geometry.

## Frozen interpretation

- If one smooth task component clearly accounts for the objective loss at the exact-fused state while hard capability is unchanged, isolate that surrogate component next.
- If the resource penalty accounts for the rejection while task-score components are preserved or improved, audit the soft-capacity regularizer next; do not tune lambda.
- If exact fusion materially improves hard capability despite a worse smooth objective, the surrogate objective is misaligned with the hard task and should be repaired before further architecture pressure.
- If exact fusion worsens hard capability as well as the smooth objective, the task-supported near-symmetry state is preferable under the current mechanism; do not force equality.
- If the seed-matched UP-44 FIXA control fails to replicate the representation-only geometry alignment, do not launch this audit automatically; first revise the causal lineage.

## Plain speak

UP-46 does not try to make anything better.

It takes the state the score liked and the exact-fusion state the optimizer briefly reached, then puts them side by side to see exactly what changed.

That tells us whether the score is rejecting exact symmetry for a good reason or because one part of the surrogate objective is misleading it.
