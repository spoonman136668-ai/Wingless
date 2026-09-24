# Wingless UP-59C — six-bank low-noise boundary across encodable phase tags

Status: preregistered scientific scale experiment.

Scientific parent: UP-58C seal `98999e0eda92f6a17e5ad9ac9db6d4d5da533c58`.

## Question

UP-58C showed that six banks are exactly separable in the fixed 16-dimensional state under multiple phase-tag geometries at zero and noise 0.0025. Uniform hexagon is mathematically unencodable for some codewords and is therefore a sealed negative rather than a candidate for further noise testing.

UP-59C maps the six-bank noise boundary for all three encodable frozen geometries without selecting one.

## Frozen design

Families:

- golden rotation;
- quadratic extension;
- fixed irregular.

Untouched deterministic schedule bases: 73M and 74M.

Noise ladder: 0.0025, 0.0035, 0.0045, 0.0055, 0.0075.

Everything else remains unchanged from UP-58C: six banks, dimension 16, 256 scenarios, depth 64, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

No family selection, decoder change, dimension increase, threshold change, or result-informed retry is permitted.

## Interpretation

The experiment identifies geometry-specific robustness boundaries. It does not promote a geometry. Any promotion or architectural change requires a separate untouched confirmation or an integration gate.


## Authoritative Windows result

Workflow run: `36049267013`

Runner: `WINGLESS-UP-C`

Source head: `27460ef4f1276fdf923632c15660f859b96e963c`

Artifact: `10830135725`

Artifact digest: `sha256:ed2346a39577291a75ea3f18a03c4f8b3e8700794d13c4715d2778a5d206404f`

Results on both untouched schedules:

- golden rotation: PASS at every tested noise through `0.0075`, with value/exact accuracy `1.0 / 1.0`;
- fixed irregular: PASS at every tested noise through `0.0075`, with value/exact accuracy `1.0 / 1.0`;
- quadratic extension: PASS at `0.0025` and `0.0035`, FAIL at `0.0045`, `0.0055`, and `0.0075`.

At the first quadratic failure (0.0045), exact-scenario accuracy fell to approximately `0.906–0.918`.

## Scientific classification

Six-bank capacity at fixed dimension 16 is strongly geometry-dependent under noise. Golden and fixed-irregular codes retain a materially larger noise margin than the quadratic code under the same decoder, scenarios, and frozen gate.

The next C-lane experiment extends only the two still-passing geometries into a higher-noise ladder on new schedules. This is boundary continuation after a sealed negative, not threshold or geometry retuning.
