# Wingless UP-155B — post-REPORT rehearsal dose

Status: preregistered scientific lexical-stability intervention experiment.

Scientific parent: sealed UP-154B 315f2b88ae66d87c15b2765b4de6fba9943b886b.

## Question

UP-154B showed that moving the same fixed rehearsal budget after a damaging REPORT block can partially repair old retention, while a 7-update post-REPORT cleanup was weak and 15 updates was clearly stronger. What is the repair-dose curve under the same fixed total budget?

## Frozen starting state and epoch structure

Reconstruct the exact accepted 15-epoch common prefix.

For each terminal epoch and subject independently (Mia, Noah):
- use the same 20 non-REPORT-tail new-family examples first;
- use the fixed four-example REPORT block: relays, announces, cites, summarizes;
- use exactly the same 15 old-rehearsal examples and accepted subject alternation;
- learning rate remains 0.08;
- total new updates = 24;
- total old updates = 15.

## Post-REPORT rehearsal-dose arms

Move the following number of the 15 rehearsal updates after REPORT:
- after_0: 15 before / 0 after
- after_4: 11 before / 4 after
- after_8: 7 before / 8 after
- after_12: 3 before / 12 after
- after_15: 0 before / 15 after

The old-rehearsal sequence is identical; each arm cuts the same ordered list at a different frozen split point.

## Measurements

Per subject × dose:
- common-prefix old retention;
- 20-example new-prefix effect;
- pre-REPORT rehearsal recovery;
- REPORT-block effect;
- post-REPORT rehearsal recovery;
- net epoch change;
- final old held-out/unseen retention;
- final mean old retention;
- final primary/secondary new-family accuracy;
- final mean new-family accuracy.

## Interpretation

The exact five-point dose curve is the result. A graded repair curve would establish a controllable fixed-budget stability knob; a threshold or plateau would locate the minimum useful post-REPORT cleanup region.

No post-result threshold is introduced and no dose is selected adaptively.

## Bounds

No extra updates, no memory change, no projector recomputation, no architecture change, no adaptive scheduling, no result-informed retry, no live activation, no production authority.
