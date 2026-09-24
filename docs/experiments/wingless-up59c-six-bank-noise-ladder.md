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
