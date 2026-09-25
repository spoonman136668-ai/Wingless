# Wingless UP-75C — fixed-budget two-choice partition routing

Status: preregistered scientific carrier-routing experiment.

Scientific parent: sealed UP-74C evidence `9b00b8faac8d348bd6a6aa96b5f0437bfad0b505`.

## Question

UP-74C showed that deterministic 4×16 partitioning improves interference behavior at the same 512-byte state budget, but high-load collapse remains. UP-75C asks whether a local two-choice routing rule can distribute keys more evenly without adding recurrent state or persistent routing metadata.

## Frozen state budget

Both arms use exactly:
- 64 float64 recurrent values;
- four 16-dimensional partitions;
- 512 recurrent-state bytes;
- gated-correction threshold 0.15;
- no exact recall.

## Frozen arms

1. `fixed_mod4`
   - partition = `key mod 4`.

2. `two_choice_norm`
   - candidate A = `key mod 4`;
   - candidate B = `(3*key + 1) mod 4`; if identical, use `(candidate A + 1) mod 4`;
   - first insertion chooses the candidate partition with lower current L2 norm; ties choose the lower partition id;
   - overwrite first searches both candidate partitions and updates the candidate with the stronger decoded match;
   - query decodes both candidate partitions and returns the value from the candidate with the stronger best-score match;
   - no persistent routing table.

## Frozen tasks

Multi-bank interference:
- logical banks: 6 and 8;
- four entities per bank;
- load: 50% and 100%;
- value vocabulary: 16;
- 96 episodes per setting.

Associative recall:
- loads: 16, 24, 32;
- four queries per episode;
- value vocabulary: 32;
- 96 episodes per load.

Overwrite/latest-value-wins:
- eight slots;
- writes: 32 and 64;
- value vocabulary: 16;
- 96 episodes per setting.

Seeds: 129M and 130M.

## Metrics

- primary value/query accuracy;
- exact-episode accuracy;
- cross-bank corruption rate;
- recurrent-state bytes.

## Bounds

No state expansion, exact recall, learned router, routing-table memory, threshold tuning, result-informed retry, production authority, or live activation.
