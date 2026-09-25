# Wingless UP-78C — marginal-impact routing closure on frozen 8×8 geometry

Status: preregistered scientific routing closure experiment.

Scientific parent: sealed UP-77C `e1fdd4071748f9e833cdfb293cb7c9897103485b`.

## Question

UP-77C showed that decorrelating two-choice candidates does not rescue norm-based routing; every tested two-choice arm remained worse than fixed 8×8 assignment. One distinct hypothesis remains: absolute load norm may be the wrong local objective. Does routing a new item to the candidate with smaller predicted post-write norm improve performance?

## Frozen carrier

All arms:
- 64 float64 recurrent values = 512 bytes;
- eight 8-dimensional fixed partitions;
- deterministic norm-1 bipolar associative vectors;
- gated correction threshold 0.15;
- no exact recall;
- no routing table;
- no persistent routing metadata.

## Frozen arms

1. `fixed_mod8`
   - partition = key mod 8.

2. `mixed_norm`
   - candidate A = key mod 8;
   - candidate B = independently mixed second hash from UP-77C;
   - first insertion selects the lower current L2 norm candidate;
   - overwrite/query inspect both candidates and use the stronger decoded match.

3. `mixed_marginal_impact`
   - same two candidates;
   - for a first insertion, compute the exact post-write squared norm for each candidate after adding that key/value vector once;
   - choose the candidate with lower predicted post-write squared norm; ties choose the lower partition id;
   - overwrite/query behavior is identical to `mixed_norm`.

No result-dependent arm selection.

## Frozen tasks

Same high-load closure subset:
- multi-bank interference: banks 6 and 8, load 50% and 100%, value vocabulary 16;
- associative recall: loads 16, 24, 32, four queries per episode, value vocabulary 32;
- overwrite/latest-value-wins: writes 32 and 64 over eight logical slots, value vocabulary 16;
- 96 episodes per setting;
- seeds 141M and 142M.

## Metrics

- primary value/query accuracy;
- exact-episode accuracy;
- cross-bank corruption rate;
- recurrent-state bytes.

## Interpretation

This is the final dynamic two-choice closure experiment. If `mixed_marginal_impact` does not materially exceed fixed_mod8 across the integrated task family, dynamic two-choice routing is frozen out of the downstream Wingless architecture unless a later integrated language task proves a new need.

## Bounds

No state expansion, exact recall, routing table, learned router, threshold tuning, result-informed retry, production authority, or live activation.
