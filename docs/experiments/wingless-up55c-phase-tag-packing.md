# Wingless UP-55C — four-bank phase-tag packing ablation

Status: preregistered architecture experiment.

Scientific parent: UP-54C seal `84ea3654f26c21dfbbda292bc7a9f67ad3d16f69`.

UP-54C showed that four-bank phase multiplexing is exact at low noise but crosses the frozen gate between noise 0.03 and 0.04. UP-55C asks whether that robustness boundary is specific to the current phase-tag geometry.

Four tag families are frozen before execution: the current quadratic tags, a golden-angle rotation family, and two fixed irrational-spread families. Each is evaluated at noise 0.04 and 0.05 with the same dimension-16 state, scenarios, transport, coherence decoder, magnitude control, and gates.

No family is selected or promoted by this experiment; it is an ablation of packing geometry. Any promising family requires later untouched confirmation.


## Authoritative Windows result

Workflow run: `36033576507`

Runner: `WINGLESS-UP-C`

Source head: `a87276870f352bd6fb76f78753847545a2dd9d01`

Artifact: `10822913545`

Artifact digest: `sha256:d659fe57e7b2fae8f0fae68de780da0c9ed5ae87330867193beeec07adf20c1c`

At four banks / noise 0.04:

- current quadratic: value `0.98681640625`, exact `0.89453125`, FAIL;
- golden rotation: value `0.998046875`, exact `0.96875`, PASS;
- irrational spread A: value `0.984619140625`, exact `0.75390625`, FAIL;
- irrational spread B: value `0.985595703125`, exact `0.859375`, FAIL.

At noise 0.05 all four families fail; golden rotation gives value `0.964599609375`, exact `0.64453125`.

## Scientific classification

Phase-tag geometry materially changes fixed-state multiplex robustness. The preregistered golden-angle family uniquely crosses the existing four-bank gate at noise 0.04, but this is an exploratory family comparison and is not promoted.

The next C-lane experiment performs untouched confirmation of the golden-angle family on independent deterministic noise schedules and maps its 0.04–0.05 boundary without selecting among new tag families.
