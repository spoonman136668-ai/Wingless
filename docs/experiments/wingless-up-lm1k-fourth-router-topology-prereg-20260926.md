# Wingless UP-LM1K — fourth-family router confusion topology

Status: preregistered scientific semantic-grounding diagnostic.

Scientific parent: sealed UP-LM1J 077930c37f583288be2c56ea157d797ad41207a7.

## Question

LM1J preserved perfect routing on the first three lexical families but achieved only 66.7% held-out accuracy on the fourth family, with REPORT routing at 33.3%. Which fourth-family surface and competitor class create the failure?

## Frozen router

Recreate the exact deterministic LM1J router:
- base/paraphrase/third grounded router unchanged;
- fourth family retains / inspects / states;
- 20 fourth-grounding epochs at lr 0.08;
- training subjects ada, ben, cy, dee;
- one-subject rehearsal of each prior family per epoch;
- no further training in this diagnostic.

## Evaluation

Fourth family:
- train subjects: ada, ben, cy, dee
- heldout subjects: eli, fay
- unseen subjects: gia, hal, ivy, jon, kia, leo

Prior-family retention:
- base, paraphrase, and third surfaces on heldout subjects.

For every surface x split:
- target class;
- predicted class;
- class accuracy;
- mean STORE probability;
- mean OBSERVE probability;
- mean REPORT probability;
- minimum target-class margin over strongest competitor.

Also report:
- confusion counts;
- per-class recall.

## Interpretation

A single failing fourth surface/class supports a class-local representation or grounding intervention. Broad fourth-family confusion would instead indicate router-capacity or rehearsal interference.

## Bounds

Diagnostic only. No router retraining beyond exact LM1J reconstruction, no byte-model training, no threshold tuning, no recall changes, no adaptive rehearsal, no result-informed retry, no live activation, no production authority.
