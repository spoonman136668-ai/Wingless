UP-LM6A PREREGISTRATION — PREPRESSURE ORIGIN OF LENGTH-16 SUPPORT

Parent UP-LM5Z qualified scientific result.
Observed parent anchor: raw sub-16 candidate count=0, eligible sub-16=0; every non-empty eligible candidate lies at length 16.

Question: is the hard length-16 support created by the frozen prepressure transformation, or is it already present in the underlying arm dynamics?

Freeze exact LM5Z condition grid and canonical hazard selector:
- same deadline profiles, resource reductions, rotations, permutations, windows, budgets, starts, throughputs;
- same selector and action policy;
- same six rounds;
- no future-schedule information;
- no adaptive policy selection;
- counterfactual only;
- no live activation.

Evaluate exactly two matched arms for every frozen condition:
PREPRESSURE_ON: exact LM5Z trajectory.
PREPRESSURE_OFF: skip only uplm3tPrepressure; change nothing else.

Immediately before each selector decision, record raw non-empty pending-order lengths across all arms.

Report per arm:
- total conditions and decision points;
- exact raw order-length histogram;
- raw sub-16 candidate count.

Classification:
PREPRESSURE_CREATES_LENGTH16_FLOOR if ON reproduces raw_sub16=0 and OFF has raw_sub16>0.
BASE_DYNAMICS_LENGTH16_FLOOR if both ON and OFF have raw_sub16=0.
ANCHOR_NOT_REPRODUCED if ON does not reproduce LM5Z.
OTHER_VALID_PATTERN otherwise.

This experiment does not authorize policy changes or prepressure removal.
Scientific negatives are valid. No post-result tuning.
