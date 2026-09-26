# Wingless UP-LM2P — unreported-eviction countdown

Status: preregistered scientific fixed-capacity shadow-forecast experiment.

Scientific parent: sealed UP-LM2O a10aeb370020820e12636dd5182e91d48287288e.

## Question
Can current recall FIFO order, free capacity, and past REPORT history derive the exact number of additional unique STORE operations until the earliest never-reported entry is evicted?

## Frozen state
Use:
- exact recall cap 16;
- first chunk of 12 stored entities;
- deferred levels D=4,5,6;
- identity rotations 0 and 7;
- value shifts 0,1,2,3.

After the first-chunk early REPORTs and before any second-chunk STORE:
- identify the earliest recall entry that has never yet been REPORTed;
- compute free slots = 16 - current recall entries;
- compute predicted unique-store distance = free slots + FIFO index of that entry + 1.

Then apply unique unseen STORE operations until that target is actually evicted, up to 13 writes.

## Preregistered expectation
- D=4 -> distance 13
- D=5 -> distance 12
- D=6 -> distance 11

## Measurements
Predicted and actual unique-store distance for each arm; exact prediction rate.

## Bounds
Shadow only. No intervention, no future report outcome, no model training, no adaptive rule, no added capacity, no result-informed retry, live activation, or production authority.
