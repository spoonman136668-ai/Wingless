# Wingless UP-LM0N — grounded paraphrase plus order composition

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0M e4cee8ac2b372c89fbbc68b5e06c8747a6bfa3eb.

## Question

LM0M showed exact learned routing across four clause-order families, while B103 showed that a few grounded examples can teach new verb surfaces. Do those mechanisms compose when the language stream changes both lexical verb surface and clause order at once?

## Frozen language model

Use the exact LM0M language model:
- trained only on the original block-ordered corpus;
- no retraining on paraphrase surfaces or new orderings;
- recurrent state dimension 64;
- 512 recurrent-state bytes;
- exact recall cap 16;
- no attention;
- no future oracle.

## Grounded paraphrase classifier

Base event surfaces:
- stores -> STORE
- observes -> OBSERVE
- reports -> REPORT

New paraphrase surfaces:
- saves -> STORE
- sees -> OBSERVE
- recalls -> REPORT

Classifier:
- same deterministic 64-D prefix-through-verb encoder;
- three-class softmax head;
- zero initialization;
- base training on stores/observes/reports for subjects ada, ben, cy, dee;
- then 20 grounding epochs on the three paraphrases using subjects ada, ben, cy, dee;
- learning rate 0.08;
- after each grounding epoch, rehearse the three base surfaces on fixed subject ada;
- no value bytes at classifier inference.

Evaluation subjects remain eli and fay.

## Frozen order families

Use paraphrase verbs in the same four dependency-valid families from LM0M:
1. per_name
2. paired_names
3. stores_then_local_reports
4. reverse_report_tail

REPORT clauses use recalls; STORE uses saves; OBSERVE uses sees.

Evaluate stream1 and stream4.

## Frozen arms

1. explicit_event_routing
2. learned_grounded_paraphrase_routing

## Metrics

Classifier:
- held-out-subject accuracy and per-class precision/recall on saves/sees/recalls;
- retained accuracy on base stores/observes/reports.

Language:
- top-1 byte accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- maximum recall entries.

## Interpretation

Matching explicit routing across paraphrase-plus-order conditions would establish composition of bounded lexical grounding and structural routing robustness. A learned-only failure identifies integration interference; a shared failure identifies generic language-distribution shift.

## Bounds

No language-model retraining on paraphrases, no state expansion, no recall-cap increase, no attention, no threshold tuning, no adaptive replay, no result-informed retry, no live activation, no production authority.
