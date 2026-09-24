# Wingless UP-54C — four-bank phase-code noise ladder

Status: preregistered scientific scale diagnostic.

Scientific parent: UP-53C seal `31b5136faf6a3575728a91d18886fff8f6cce122`.

UP-53C showed perfect fixed-state coherence decoding through three simultaneous memory banks at noise 0.05, with a sharp failure at four banks. UP-54C freezes the phase tags and compares bank counts 3 and 4 across noise amplitudes `0, 0.005, 0.01, 0.02, 0.03, 0.04, 0.05`.

All other mechanics remain unchanged: dimension 16, 256 scenarios, global-phase nuisance, 64 reversible transport steps, complete per-entity prototype decoding, and the same frozen gate.

If four banks fail even at zero noise, the boundary is representational/tag-geometry limited. If four banks pass at low noise and fail later, the boundary is finite-separation robustness. No phase-tag optimization is permitted in this experiment.


## Authoritative Windows result

Workflow run: `36032865894`

Runner: `WINGLESS-UP-C`

Source head: `3c2f25c8264525e193f0b211d3826e45ae24c00e`

Artifact: `10823966778`

Artifact digest: `sha256:d75130a1bdd2c4a949c3920fdc68916b50259fca6fe62468ee718e18fa8604da`

Three-bank coherence decoding remains perfect through noise `0.05`.

Four-bank coherence decoding:

- noise 0 through 0.03: value/exact scenario accuracy `1.0`, PASS;
- noise 0.04: value `0.9853515625`, exact scenario `0.8828125`, FAIL;
- noise 0.05: value `0.962158203125`, exact scenario `0.6953125`, FAIL.

The four-bank coherence margin falls from `0.018907745714141222` at zero noise to `0.00203256571848609` at noise 0.03 before collapsing near zero at the failure boundary.

## Scientific classification

Four-bank phase multiplexing is not blocked by an exact representational collision. It is perfectly decodable at zero and low noise; the current phase-tag geometry simply has a much smaller separation margin than the three-bank code.

The next C-lane experiment compares several preregistered phase-tag packing families at four banks under the frozen boundary noise levels. This is architecture exploration, not post-result coefficient tuning.
