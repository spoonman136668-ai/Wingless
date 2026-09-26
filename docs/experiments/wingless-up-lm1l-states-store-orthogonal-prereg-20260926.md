# Wingless UP-LM1L — fourth-family states-vs-STORE orthogonal correction

Status: preregistered scientific semantic-routing representation experiment.

Scientific parent: sealed UP-LM1K e5944e3a6a52a1e4bd98c5c880c5eca964075fd5.

## Question

LM1K localized the fourth-family router failure entirely to the REPORT surface "states", which is misrouted toward STORE. Can one frozen representation correction that removes the STORE competitor direction from "states" restore REPORT routing without changing any other fourth-family surface or grounding schedule?

## Frozen router/training

Use the exact LM1K reconstruction:
- same 64-D router representation;
- same base/paraphrase/third/fourth grounding schedule;
- 20 fourth-family grounding epochs;
- learning rate 0.08;
- same prior-family rehearsal;
- no threshold changes.

## Frozen correction

Build one STORE competitor direction before router training:
- training subjects only: ada, ben, cy, dee;
- STORE surfaces: stores, saves, archives, retains;
- direction = mean raw 64-D representation across those surfaces and subjects;
- normalize direction.

For the correction arm only:
- when the surface is exactly "states", subtract its projection onto the frozen STORE direction;
- rescale the residual to the original vector norm;
- use the corrected vector for both fourth-family grounding and evaluation of "states".

All other surfaces remain raw and unchanged.

## Arms

1. raw_control — exact LM1K router.
2. states_store_orthogonal — frozen correction above.

## Evaluation

For each arm:
- fourth train / heldout / unseen family accuracy;
- STORE / OBSERVE / REPORT recall;
- per-surface predicted class rates and probabilities;
- states REPORT accuracy and STORE-confusion rate;
- prior base/paraphrase/third heldout accuracy.

## Interpretation

If the correction repairs "states" while prior families and the other fourth surfaces remain intact, the LM1K failure is a local representation-conflict mechanism rather than a grounding-capacity limit.

## Bounds

No byte-model training, no adaptive correction, no correction recomputation after results, no threshold search, no extra grounding epochs, no recurrent training, no result-informed retry, no live activation, no production authority.
