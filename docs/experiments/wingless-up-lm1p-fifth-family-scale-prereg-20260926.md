# Wingless UP-LM1P — fifth-family balanced integration scale test

Status: preregistered scientific language/router/memory scaling experiment.

Scientific parent: sealed UP-LM1O fdfc618ce6f6a323683714870ae1ff6b6eeba7c8.

## Question

LM1O showed that the four-family mirrored and rotating-palindromic byte schedules remain fully compatible with corrected semantic routing and exact memory. Does the same balanced scheduling principle scale when a fifth lexical family is added without changing recurrent state or recall capacity?

## Frozen starting point

Use the exact LM1O rotating-palindromic integrated arm as the starting byte model:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family byte training history;
- exact corrected fourth-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Fifth lexical family

New event surfaces:
- STORE: banks
- OBSERVE: surveys
- REPORT: declares

Byte corpus uses the same names, values, update-count distribution, and structural order families as prior lexical families.

The fifth-family surfaces are absent from the four-family starting model's training data.

## Fifth-family router grounding

Starting from the exact corrected four-family router:
- 20 grounding epochs;
- learning rate 0.08;
- fifth-family training subjects: ada, ben, cy, dee;
- fifth-family evaluation subjects: eli, fay plus the existing unseen subject set;
- each epoch rehearses one fixed subject from every prior lexical family;
- no representation correction is introduced for fifth-family surfaces in this experiment.

## Byte adaptation arms

All adaptation arms run 4 epochs and use matched base / paraphrase / third / fourth / fifth examples.

1. no_fifth_adaptation
   - no additional byte-model updates.

2. five_family_cyclic
   - one full lr 0.08 update per family per matched example in order:
     base, paraphrase, third, fourth, fifth.

3. five_family_rotating_palindromic
   - center family rotates by example index modulo 5;
   - the four non-center families each receive one 0.04 half-step before center in cyclic order;
   - center receives one 0.08 full step;
   - the four non-center half-steps repeat in exact reverse order;
   - nominal learning-rate mass remains 0.08 per family per matched example.

No arm changes recurrent parameters.

## Evaluation

Byte metrics:
- base / paraphrase / third / fourth / fifth held-out top-1 accuracy and perplexity;
- minimum, mean, and spread across all five families.

Integrated fifth-family evaluation:
- block;
- per_name;
- paired_names;
- stores_then_local_reports;
- reverse_report_tail;
- stream1 and stream4.

Metrics:
- fifth-family byte top-1/perplexity;
- dependent first-byte accuracy;
- query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries;
- fifth-router train/heldout/unseen accuracy and per-surface recall.

## Interpretation

If rotating-palindromic scheduling preserves a stronger five-family minimum/mean balance than cyclic scheduling while router/memory remain exact, the operator-balance mechanism has scaled one lexical family further. Router failure would instead identify a semantic representation boundary before byte capacity can be judged cleanly.

## Bounds

No recurrent training, no recall-cap increase, no attention, no adaptive schedule, no fifth-family representation correction, no threshold tuning, no result-informed retry, no live activation, no production authority.
