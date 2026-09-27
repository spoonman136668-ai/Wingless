# Wingless UP-LM3D — compressed intervention window

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3C bb29b01bf50424f94cc40412fd395059d3b0cc48.

## Question

Does per-round throughput become causal for earliest-deadline scheduling when the fixed total action budget cannot be spent serially before the remaining intervention window closes?

## Frozen scenario

Reuse UP-LM3C unchanged:
- exact recall cap 16;
- six arms;
- deadline profiles by_deferred_level and by_layout;
- identity rotations 3 and 11;
- permutations identity and reverse;
- one action per arm per round;
- total action budget 6;
- global pressure rounds 6;
- throughput 1, 2, 3;
- baseline, earliest_deadline, fixed_order.

## Frozen compression

Corrective actions are disabled until:
- global round 1; or
- global round 2.

After onset, the frozen policy may use its normal throughput subject to the same total budget and distinct-arm constraint.

## Measurements

Per deadline profile × onset round × throughput:
- actions used;
- completed and failed originals;
- failures prevented versus baseline;
- earliest-deadline advantage over fixed-order.

## Interpretation

If higher throughput reduces earliest-deadline failures only after delayed onset, throughput is a temporal reachability resource: it matters when the action window is too short to deploy the fixed budget serially. If outcomes remain flat, the remaining bottleneck lies deeper than simple budget deployment rate.

## Bounds

Counterfactual only. No adaptive onset, budget, or throughput; no capacity change; no future schedule oracle; no semantic priority; no extra training; no live activation; no production authority.
