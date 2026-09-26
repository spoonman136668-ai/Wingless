# Wingless UP-106B — sequential vocabulary growth

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-105B 20ff005b26a6b5cd955dadfc947e66f43a7f606a.

## Question

UP-105B showed that replaying two fixed subjects per previously acquired verb preserves three sequentially learned verbs without damaging the base lexicon. Does that fixed replay rule remain stable as the acquired vocabulary doubles to six new verbs?

## Frozen representation and optimization

- deterministic 64-D prefix-through-verb encoder;
- three-class linear softmax head;
- zero-initialized base classifier trained on the frozen base lexicon;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- no threshold tuning;
- no state expansion.

Base lexicon:
- STORE: stores, keeps
- OBSERVE: observes, sees
- REPORT: reports, recalls

## Frozen acquisition sequence

1. holds -> STORE
2. notes -> OBSERVE
3. tells -> REPORT
4. saves -> STORE
5. watches -> OBSERVE
6. remembers -> REPORT

Each new verb:
- grounds on ada, ben, cy, dee;
- evaluates on eli, fay.

After every epoch:
- rehearse each base verb on fixed subject ada;
- rehearse every previously acquired verb on fixed subjects ada and ben;
- current acquisition gets only its four normal grounding examples.

## Evaluation

At stage 0 and after each of six acquisitions:
- base-lexicon accuracy over all six subjects;
- accuracy of each acquired verb on eli/fay;
- aggregate accuracy over all verbs acquired so far.

## Interpretation

Stable accuracy through six batches supports bounded vocabulary growth under fixed replay. Failure identifies the acquisition depth at which interference returns.

## Bounds

No adaptive replay, no encoder change, no extra hidden state, no class weighting, no pretrained semantics, no threshold search, no result-informed retry, no live activation, no production authority.
