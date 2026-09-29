UP-LM6B PREREGISTRATION — PREPRESSURE ARM EXTENSION MAP

Parent UP-LM6A qualified PREPRESSURE_CREATES_LENGTH16_FLOOR.

Question: how does the frozen prepressure transform produce the hard length-16 support?

Freeze the exact arm construction and prepressure definitions from LM6A/LM3T.
Evaluate every profile × rotation × permutation × arm exactly once:
profiles={deferred_only,layout_only,hybrid_min}
rotations={5,13}
permutations={identity,reverse,rotate2}
12 arms each = 216 rows.

For each row record:
- prepressure target;
- pending-order length before prepressure;
- first-pending countdown before prepressure;
- exact number of prepressure writes;
- pending-order length after prepressure;
- first-pending countdown after prepressure.

Apply the exact uplm3tPrepressure stopping rule; do not run selector actions or resource schedules.

Classification:
UNIFORM_EXTENSION if every row receives the same positive write count.
TARGET_DEPENDENT_EXTENSION if write counts vary across rows and at least one write is applied.
NO_PREPRESSURE_EXTENSION if all write counts are zero.
ANCHOR_NOT_REPRODUCED if LM6A does not reproduce PREPRESSURE_CREATES_LENGTH16_FLOOR.
OTHER_VALID_PATTERN otherwise.

Observational mechanism map only. No policy change, no live activation, no post-result tuning.
