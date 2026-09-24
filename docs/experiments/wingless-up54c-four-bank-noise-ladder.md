# Wingless UP-54C — four-bank phase-code noise ladder

Status: preregistered scientific scale diagnostic.

Scientific parent: UP-53C seal `31b5136faf6a3575728a91d18886fff8f6cce122`.

UP-53C showed perfect fixed-state coherence decoding through three simultaneous memory banks at noise 0.05, with a sharp failure at four banks. UP-54C freezes the phase tags and compares bank counts 3 and 4 across noise amplitudes `0, 0.005, 0.01, 0.02, 0.03, 0.04, 0.05`.

All other mechanics remain unchanged: dimension 16, 256 scenarios, global-phase nuisance, 64 reversible transport steps, complete per-entity prototype decoding, and the same frozen gate.

If four banks fail even at zero noise, the boundary is representational/tag-geometry limited. If four banks pass at low noise and fail later, the boundary is finite-separation robustness. No phase-tag optimization is permitted in this experiment.
