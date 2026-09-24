# Wingless UP-58C — six-bank phase-tag geometry ablation

Status: preregistered scientific scale experiment.

Scientific parent: UP-57C seal `010642f03fd990bffcc276ca116fb971eed8b667`.

UP-57C showed five banks are separable at low noise with the confirmed golden construction, while six golden banks fail even at zero noise. UP-58C distinguishes a golden-geometry limitation from a broader fixed-state capacity boundary.

At six banks and fixed dimension 16, four tag families are frozen before execution: golden rotation, uniform hexagon spacing, the original quadratic formula extended to six banks, and a fixed irregular set. Each is tested at zero noise and noise 0.0025 with unchanged scenarios, reversible transport, coherence decoder, magnitude control, and gate.

No family is selected or promoted by this experiment. A positive family requires untouched confirmation.


## Authoritative Windows result

Workflow run: `36048603948`

Runner: `WINGLESS-UP-C`

Source head: `1cb083dde907d000fa430396507870348dd6d593`

Artifact: `10829845310`

Artifact digest: `sha256:99fb1da66a6d3464b7e639b39d9d367d7e7bb3a723c510993f41d90253e9af64`

Results:

- golden rotation: PASS at noise 0 and 0.0025; value/exact accuracy `1.0 / 1.0` at both;
- quadratic extension: PASS at noise 0 and 0.0025; value/exact accuracy `1.0 / 1.0` at both;
- fixed irregular: PASS at noise 0 and 0.0025; value/exact accuracy `1.0 / 1.0` at both;
- uniform hexagon: unencodable at both conditions because some codewords produce exact phase-code cancellation; gate FAIL by construction.

Minimum margins at noise 0.0025 were approximately:

- golden: `0.0026099105724114446`;
- fixed irregular: `0.0012825637988460592`;
- quadratic: `0.00014709039515914402`.

## Scientific classification

Six-bank failure is not a fundamental zero-noise capacity limit of the fixed 16-dimensional state. Multiple frozen phase-tag geometries remain exactly separable at six banks under low noise.

The meaningful next question is the six-bank noise boundary across the three encodable families on untouched schedules. Uniform hexagon is retired from that ladder because it is mathematically non-encodable for some frozen codewords; this is a recorded negative, not a tuning decision.
