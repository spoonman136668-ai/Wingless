UP-256C PREREGISTRATION — CURRENT-STATE COMPONENT SUFFICIENCY

Parent UP-255C qualified CURRENT_STATE_RESOLVES_ALIASING: STATIC9 leaves the frozen residual mixed groups, while STATIC9 plus current event/signature diagnostics resolves all mixed groups across the same 64 paired initial states.

Question: which current-state component is sufficient for that separation?

Freeze exact UP-255C paired substrate and outcomes:
A = schedule2|advance22
B = schedule3|advance0
same 64 initial states;
same nine initial native coordinates;
same intervention logic;
no new policy;
no fitting;
no classifier;
no adaptive feature selection;
no live activation.

Use only existing UP-237C observational diagnostics:
- event_present for A and B;
- current_signature for A and B.
Trajectory signatures are excluded from this experiment.

Evaluate exactly four partitions:
STATIC9
STATIC9_EVENTS = STATIC9 plus A/B event_present
STATIC9_SIGNATURES = STATIC9 plus A/B current_signature
STATIC9_EVENTS_SIGNATURES = STATIC9 plus both components

The joint partition must reproduce the UP-255C current-state partition exactly.

Classification:
BOTH_INDIVIDUALLY_SUFFICIENT if EVENTS and SIGNATURES each have zero mixed groups.
EVENTS_SUFFICIENT if EVENTS has zero mixed groups and SIGNATURES does not.
CURRENT_SIGNATURE_SUFFICIENT if SIGNATURES has zero mixed groups and EVENTS does not.
JOINT_ONLY if each component alone remains mixed but the joint partition has zero mixed groups.
CURRENT_COMPONENTS_REDUCE_NOT_RESOLVE if the joint partition reduces mixed groups but remains nonzero.
CURRENT_COMPONENTS_DO_NOT_REDUCE if the joint partition does not reduce mixed groups.
ANCHOR_NOT_REPRODUCED if UP-255C no longer reproduces CURRENT_STATE_RESOLVES_ALIASING or the joint partition differs from its current-state partition.

Scientific negatives are valid. No post-result tuning.
