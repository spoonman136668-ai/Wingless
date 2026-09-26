# Wingless UP-127B — class-balanced STORE gate training

Status: preregistered scientific representation-routing experiment.

Scientific parent: sealed UP-126B ad44f3b77c0914f60198de80f22cd8094f967366.

## Question

UP-126B showed that the learned STORE gate has perfect precision but broad false negatives across stores, holds, saves, and archives, while every non-STORE surface remains below probability 0.421. Does correcting the 1:2 positive/negative training exposure recover STORE recall at the same frozen threshold without introducing non-STORE false positives?

## Frozen encoder and gate form

- raw deterministic 64-D prefix-through-verb encoder;
- linear logistic gate;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- threshold 0.5;
- same training surfaces and subjects as UP-125B/126B;
- no projector or downstream classifier in this probe.

## Arms

1. original_unbalanced
   - exact UP-125B training procedure.

2. balanced_positive_weight2
   - exact same example order and subjects;
   - STORE-positive logistic gradient contribution multiplied by 2.0;
   - non-STORE contribution unchanged;
   - no threshold adjustment.

This equalizes total positive and negative loss weight while leaving sample order and threshold fixed.

## Evaluation

Subject splits:
- train: ada, ben, cy, dee
- heldout: eli, fay
- unseen: gia, hal, ivy, jon, kia, leo

For every arm and split:
- accuracy;
- precision;
- recall.

For every arm x split x surface:
- mean probability;
- minimum probability;
- maximum probability;
- positive rate.

## Interpretation

If balanced weighting raises STORE recall while preserving high precision, the B125 failure is primarily training-objective prior bias. If specific surfaces remain below threshold, the next stage returns to surface geometry rather than further class weighting.

## Bounds

No threshold tuning, no learning-rate search, no encoder change, no projector use, no downstream classifier training, no adaptive weighting, no result-informed retry, no live activation, no production authority.
