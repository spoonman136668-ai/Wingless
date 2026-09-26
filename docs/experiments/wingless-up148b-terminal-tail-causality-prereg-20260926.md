# Wingless UP-148B — terminal-tail causal contrast

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-147B 8a6c31cc4f3c3c439d8d297d953b847cfe8b53ae.

## Question

UP-147B showed an identity × curriculum-position interaction: shifts 6/18 were safer, shifts 0/12 were damaging, and schedules ending in safer tails finished with better old/new accuracy. Holding the entire first 15 epochs fixed, does swapping only the final five-epoch tail from a safe identity to a damaging identity causally produce a final retention gap?

## Frozen common prefix

Both arms start from the exact same UP-147B pre-adaptation gate and run:

- epochs 1–5: shift 0;
- epochs 6–10: shift 6;
- epochs 11–15: shift 12.

Every epoch:
- learning rate 0.08;
- 24 new-family updates;
- 15 old-family rehearsal updates;
- final four new-family updates after the old anchor.

## Terminal arms

Only epochs 16–20 differ:

- safe_terminal: shift 18 for five epochs.
- damaging_terminal: shift 12 for five epochs.

The second arm intentionally repeats shift 12; this is the causal manipulation. Both arms have identical total updates and identical history through epoch 15.

## Measurements

- old retention at the end of the common 15-epoch prefix;
- terminal-block mean prefix damage;
- terminal-block mean anchor recovery;
- terminal-block mean tail damage;
- terminal-block mean net epoch change;
- final old held-out/unseen retention;
- final primary/secondary new-family accuracy;
- final mean old retention;
- final mean new-family accuracy.

## Interpretation

If the damaging-terminal arm finishes materially below the safe-terminal arm from the same 15-epoch state, terminal tail identity causally contributes to the retention gap. If final outcomes converge, the UP-147B schedule differences depended on earlier identity × position history rather than the terminal block itself.

No post-result threshold is introduced; the exact contrast is the result.

## Bounds

No adaptive tail selection, no extra updates, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
