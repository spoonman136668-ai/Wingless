# Wingless UP-LM0K — learned-routing stream-length scaling

Status: preregistered scientific language/memory robustness experiment.

Scientific parent: sealed UP-LM0J `01af2d8f9a0aa00245a567be70413734cb7daaaa`.

## Question

UP-LM0J matched explicit three-way event routing through stream4. Does learned prefix-timed STORE/OBSERVE/REPORT routing remain stable across longer continuous streams without resetting recurrent or exact-memory state between paragraphs?

## Frozen language system

Identical to UP-LM0J:
- deterministic 64-D recurrent language state;
- 512 recurrent-state bytes;
- 20 language-model epochs;
- learning rate 0.08;
- exact recall cap 16;
- same held-out paragraph corpus;
- same three-way prefix classifier;
- no attention;
- no future oracle.

## Frozen stream ladder

Continuous stream sizes:
- 1 paragraph
- 4 paragraphs
- 8 paragraphs
- 16 paragraphs

Within each stream chunk:
- recurrent state persists across paragraph boundaries;
- exact recall persists across paragraph boundaries;
- newline bytes are processed exactly as in UP-LM0J stream4;
- state and recall reset only at the next chunk.

## Frozen arms

1. explicit_event_routing
2. learned_prefix_threeway

## Evaluation

For every arm x stream-size cell:
- top-1 byte accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- exact accuracy by update count;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- maximum recall entries.

Classifier train-subject and held-out-subject metrics are also reported unchanged from UP-LM0J.

## Interpretation

Stable learned routing across stream8/16 would support longer-horizon integration without state or memory expansion. Degradation is sealed as a stream-length boundary rather than tuned away.

## Bounds

No state expansion, no recall-cap increase, no attention, no threshold search, no retraining per stream length, no result-informed retry, no live activation, no production authority.
