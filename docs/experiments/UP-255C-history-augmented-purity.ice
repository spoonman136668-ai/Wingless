UP-255C PREREGISTRATION — HISTORY-AUGMENTED PURITY

Parent UP-254C qualified FULL_NATIVE_STATE_NONSEPARABLE: all nine frozen initial native coordinates still leave 4 mixed groups across the 64 paired initial states.

Question: is the residual ambiguity static-state aliasing that becomes separable when the already-existing bounded trajectory diagnostics from the paired arms are included?

Freeze exact UP-254C paired substrate and outcomes:
A = schedule2|advance22
B = schedule3|advance0
same 64 initial states, same nine initial native coordinates, same intervention logic, no new policy, no fitting, no classifier, no adaptive feature selection, no live activation.

Use only existing UP-237C diagnostic outputs from each paired arm:
event_present
current_signature
trajectory_signature
These are observational diagnostics and may not affect execution.

Evaluate three exact partitions:
STATIC9 = all nine initial native coordinates.
STATIC9_CURRENT = STATIC9 plus A/B event_present and current_signature.
STATIC9_TRAJECTORY = STATIC9_CURRENT plus A/B trajectory_signature.

Report groups, mixed groups, max classes/group for all three.

Classification:
TRAJECTORY_RESOLVES_ALIASING if STATIC9 mixed>0, CURRENT mixed>0, TRAJECTORY mixed=0.
CURRENT_STATE_RESOLVES_ALIASING if STATIC9 mixed>0 and CURRENT mixed=0.
HISTORY_REDUCES_NOT_RESOLVES if augmented mixed groups are fewer than STATIC9 but remain >0.
HISTORY_DOES_NOT_REDUCE_ALIASING if neither augmented partition reduces mixed groups.
ANCHOR_NOT_REPRODUCED if STATIC9 does not reproduce UP-254C.
Scientific negatives valid. No post-result tuning.
