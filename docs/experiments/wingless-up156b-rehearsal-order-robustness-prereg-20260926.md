# Wingless UP-156B — rehearsal-order robustness

Status: preregistered scientific lexical-stability intervention experiment.

Scientific parent: sealed UP-155B 78f13eaa8774efe6ee17cc8343d4e19e8bff025e.

## Question

UP-155B showed stronger local repair when more of the fixed 15 old-rehearsal updates occur after a damaging REPORT block, but dose is confounded with which specific old examples occupy the post-REPORT tail. Does the repair-dose pattern survive deterministic changes to old-rehearsal order?

## Frozen starting state and new-family ordering

Reconstruct the exact accepted 15-epoch common prefix.

For each terminal epoch and subject independently:
- use the same 20 non-REPORT-tail new-family examples first;
- use the fixed four-example REPORT block: relays, announces, cites, summarizes;
- use exactly the same 15 old rehearsal example identities and the accepted subject assignment for that epoch;
- learning rate remains 0.08;
- total new updates = 24;
- total old updates = 15.

## Frozen post-REPORT doses

- after_0
- after_8
- after_12
- after_15

## Frozen rehearsal orders

1. canonical
   - accepted 0→14 old-example order.

2. reverse
   - accepted identities in 14→0 order.

3. rotating
   - deterministic cyclic start = ((terminal epoch - 15) * 3) mod 15, wrapping through all 15 identities.

The subject paired with each old-example identity is computed from its original accepted index and epoch before reordering, so order changes without changing the identity/subject multiset.

## Measurements

Per subject × dose × order:
- common-prefix old retention;
- new-prefix effect;
- pre-REPORT rehearsal recovery;
- REPORT effect;
- post-REPORT rehearsal recovery;
- net epoch change;
- final old held-out/unseen retention;
- final mean old retention;
- final primary/secondary new-family accuracy;
- final mean new-family accuracy.

## Interpretation

If larger post-REPORT doses retain stronger repair across rehearsal orders, dose is a robust mechanism. If order materially changes which doses help, rehearsal identity/order remains part of the causal mechanism.

No post-result threshold or arm selection is introduced.

## Bounds

No extra updates, no memory change, no projector recomputation, no architecture change, no adaptive scheduling, no result-informed retry, no live activation, no production authority.
