# Wingless UP-101B — held-out verb transfer

Status: preregistered scientific semantic-routing generalization experiment.

Scientific parent: sealed UP-100B `0c91e700b58bc2873e932bfdd7cc2d0260beadf8`.

## Question

UP-100B showed perfect three-way routing when all nine event verbs were represented during training and read out at the verb boundary. Does that prefix-timed representation generalize to event verbs that are completely absent from classifier training?

## Frozen lexical domain

Names:
- ada
- ben
- cy
- dee
- eli
- fay

Classes and verbs:

STORE:
- training: stores, keeps
- held out: holds

OBSERVE:
- training: observes, sees
- held out: notes

REPORT:
- training: reports, recalls
- held out: tells

No held-out verb appears in classifier training.

## Frozen encoder/classifier

- exact UP-99B/UP-100B deterministic 64-D event encoder;
- classifier input ends at `<name> <verb>`;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- deterministic lexical order;
- no threshold tuning;
- no explicit event class at learned inference.

## Evaluation

Classifier:
- train-verb overall accuracy and per-class precision/recall;
- held-out-verb overall accuracy and per-class precision/recall;
- per-verb accuracy for all nine verbs;
- held-out 3x3 confusion matrix.

Routing probe:
- compare explicit event class against learned prefix classifier;
- event streams use held-out verbs only for STORE/OBSERVE/REPORT;
- exact recall cap 16;
- target sets 4, 8, 16;
- total writes 32, 64, 128, 256;
- seeds 159M and 160M.

## Interpretation

Success would show lexical-class transfer beyond trained verb surfaces. Failure would establish that verb-boundary timing solves suffix interference but not unseen lexical semantics.

## Bounds

No pretrained semantics, no class weights, no threshold search, no encoder change, no attention, no recurrent-state expansion, no recall-cap increase, no result-informed retry, no live activation, no production authority.
