# Wingless UP-LM0Q — equal-update adaptation order

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0P 7321157517e23061a4bc15d498d535d717d8103a.

## Question

LM0P showed a strong retention/adaptation tradeoff as base replay volume changed. With replay volume held exactly constant, how much of that tradeoff is caused by update order and recency?

## Frozen starting point

Identical to LM0P:
- deterministic 64-D recurrent byte language model;
- 20 base training epochs before adaptation;
- paraphrase adaptation depth fixed at 4 epochs;
- learning rate 0.08;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen update counts

Every arm performs, per adaptation epoch:
- exactly one complete paraphrase training corpus worth of examples;
- exactly one complete original base training corpus worth of examples.

The base and paraphrase corpora use the same structural split and have equal example counts. No arm receives more examples than another.

## Frozen schedule arms

1. para_then_base
   - complete paraphrase pass, then complete base pass.
   - exact LM0O replay-1 schedule.

2. base_then_para
   - complete base pass, then complete paraphrase pass.

3. example_interleaved
   - for index i in deterministic corpus order:
     - train paraphrase example i;
     - then train base example i.
   - repeat for the full corpus.

All arms start from identical copies of the same 20-epoch base model.

## Evaluation

For every arm:
1. base held-out:
   - top-1 byte accuracy
   - perplexity
2. paraphrase block held-out:
   - top-1 byte accuracy
   - perplexity
3. paraphrase block plus four LM0N order families, stream1:
   - dependent first-byte accuracy
   - whole-query-set exact accuracy
   - event-routing accuracy
   - report-routing accuracy
   - maximum recall entries
   - top-1 byte accuracy and perplexity

## Interpretation

A schedule that improves both distributions at identical update counts would identify recency/order as a controllable mechanism rather than a pure replay-volume tradeoff.

## Bounds

No router retraining, no adaptive ordering, no extra examples, no state expansion, no recall-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
