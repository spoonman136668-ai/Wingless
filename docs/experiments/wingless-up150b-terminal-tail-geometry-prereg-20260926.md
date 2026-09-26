# Wingless UP-150B — terminal-tail frozen update geometry

Status: preregistered scientific lexical-geometry diagnostic.

Scientific parent: sealed UP-149B 6d67cc6d4b3b582879e7771c44434fea3fca853f.

## Question

UP-149B established a causal split: terminal shifts 6/18 protect old and new accuracy, while shifts 0/12 damage them from the identical 15-epoch state. What frozen update geometry separates those tail identities?

## Frozen state

Reconstruct exactly the accepted 15-epoch common prefix:
- epochs 1–5 shift 0;
- epochs 6–10 shift 6;
- epochs 11–15 shift 12;
- identical learning rate, new updates, old anchor, and final four-update structure.

After epoch 15 the gate is frozen for diagnosis.

## Frozen diagnostic

For terminal shifts {0,6,12,18}:
- identify the exact four examples that would occupy the post-anchor tail;
- compute each example's one-step gate update vector from a clone of the same frozen gate;
- sum those four frozen one-step vectors without chaining them.

Compute a frozen old-anchor reference update vector using the exact epoch-16 old-anchor examples, again with every one-step vector measured from the same frozen gate.

Measurements:
- aggregate tail-update norm;
- old-anchor reference norm;
- cosine(tail update, old-anchor update);
- mean correct-class probability of the four tail examples at the frozen gate;
- pairwise cosine among all four aggregate tail-update vectors.

No diagnostic update is retained and no terminal training is executed.

## Interpretation

If the safe 6/18 tails share update geometry that differs from damaging 0/12—especially alignment with the old-anchor correction direction or a distinct pairwise cluster—the result identifies a candidate mechanism for the causal retention split. If geometry does not separate the groups, the cause lies beyond this first-order frozen update view.

No post-result threshold is introduced.

## Bounds

No adaptive tail selection, no extra training, no new memory, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
