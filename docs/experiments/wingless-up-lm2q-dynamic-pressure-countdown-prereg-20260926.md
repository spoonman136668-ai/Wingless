# Wingless UP-LM2Q — dynamic-pressure countdown

Status: preregistered scientific fixed-capacity shadow-forecast experiment.

Scientific parent: sealed UP-LM2P bcd8b5ea2a5c6c8c6ed9e5bcaad081e561227a03.

## Question

UP-LM2P predicted time-to-eviction exactly for a static pending set. Does the same native-state countdown remain exact after prior STORE pressure and REPORT closures dynamically change which pending dependency is earliest?

## Frozen arms

Use:
- exact recall cap 16;
- deferred levels D=4,5,6;
- identity rotations 0 and 7;
- value shifts 0,1,2,3.

Initialize the 12 first-chunk memories exactly as in LM2P. Mark the first 12-D identities as already REPORTed.

Then execute eight unique pressure STOREs. After STORE 4, mark the initially earliest pending identity REPORTed. After STORE 8, mark the initially second-earliest pending identity REPORTed. These closure identities are frozen from D before execution; there is no adaptive target selection.

## Frozen snapshots

Measure at:
- initial state;
- after STORE 4 plus closure 1;
- after STORE 8 plus closure 2.

At each snapshot, predict unique STOREs until the next unreported eviction using only:
- current recall occupancy;
- current FIFO order;
- identities already REPORTed in past history.

Prediction formula:
free slots + FIFO index of earliest currently-unreported resident + 1.

Ground truth is computed from a cloned current recall state by applying unique STOREs only until the first unreported resident is evicted. The real state is not mutated by this ground-truth simulation.

## Measurements

Per D × rotation × value shift × snapshot:
- free slots;
- target identity and FIFO index;
- predicted unique STORE countdown;
- cloned-ground-truth countdown;
- exact match.

## Bounds

Shadow only. No intervention, no future report schedule inspection, no extra memory, no semantic priority, no adaptive rule, no result-informed retry, no live activation, or production authority.
