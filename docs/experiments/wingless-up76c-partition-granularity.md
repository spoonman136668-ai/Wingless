# Wingless UP-76C — fixed partition granularity closure

Status: preregistered scientific carrier-geometry closure experiment.

Scientific parents:
- sealed UP-74C `9b00b8faac8d348bd6a6aa96b5f0437bfad0b505`;
- sealed UP-75C `478c60bef8d4104600043e875faeb8de1ea3e409`.

## Question

UP-74C showed fixed partitioning helps at a constant 512-byte state budget. UP-75C showed dynamic two-choice routing is worse. Before closing this line, which fixed partition granularity is best supported: fewer larger partitions or more smaller partitions?

## Frozen state budget

All arms use exactly 64 float64 recurrent values = 512 bytes, gated correction threshold 0.15, no exact recall, no routing table.

## Frozen arms

1. `partition2x32` — key mod 2, 32 dimensions per partition.
2. `partition4x16` — key mod 4, 16 dimensions per partition.
3. `partition8x8` — key mod 8, 8 dimensions per partition.

Within each partition the associative code is deterministic bipolar with norm 1.

## Frozen tasks

High-information subset only:

Multi-bank interference:
- logical banks 6 and 8;
- load 50% and 100%;
- four entities per bank;
- value vocabulary 16.

Associative recall:
- loads 16, 24, 32;
- four queries per episode;
- value vocabulary 32.

Overwrite/latest-value-wins:
- writes 32 and 64;
- eight logical slots;
- value vocabulary 16.

96 episodes per setting.
Seeds: 133M and 134M.

## Metrics

- primary value/query accuracy;
- exact episode accuracy;
- cross-bank corruption rate;
- recurrent-state bytes.

## Interpretation

This is a closure experiment. After it, fixed partition granularity is frozen for downstream use; no additional partition-count ladder follows unless a later integrated task proves a new failure.

## Bounds

No state expansion, no exact recall, no learned router, no dynamic partition selection, no threshold tuning, no result-informed retry, no production authority, no live activation.
