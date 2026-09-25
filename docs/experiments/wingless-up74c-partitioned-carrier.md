# Wingless UP-74C — fixed-budget partitioned carrier

Status: preregistered scientific carrier-geometry experiment.

Scientific parent: sealed UP-SQ0 evidence `8de9e11f26cc44c75f51d26cb606f97e35816479`.

## Question

SQ0 showed strong interference in a flat 64-dimensional associative carrier as working-set load rises. Can deterministic key partitioning reduce interference while holding total recurrent state fixed at exactly 64 float64 values (512 bytes)?

## Frozen arms

1. `flat64`
   - existing SQ0 64-dimensional associative carrier;
   - gated correction enabled;
   - no exact recall.

2. `partition4x16`
   - same total 64 float64 state values;
   - four disjoint 16-dimensional partitions;
   - route key to partition `key mod 4`;
   - deterministic bipolar vector inside the selected partition;
   - gated correction threshold remains 0.15;
   - no exact recall.

## Frozen tasks

### Multi-bank interference
- logical banks: 4, 6, 8;
- four entities per logical bank;
- load: 50% and 100%;
- value vocabulary: 16;
- 96 episodes per setting.

### Associative recall
- loads: 8, 16, 24, 32 unique key/value pairs;
- four queries per episode;
- value vocabulary: 32;
- 96 episodes per load.

### Overwrite
- eight logical slots;
- writes: 16, 32, 64;
- latest-value-wins;
- value vocabulary: 16;
- 96 episodes per setting.

Seeds: 125M and 126M.

## Metrics

- value/query accuracy;
- exact-episode accuracy;
- cross-bank corruption rate where applicable;
- recurrent state bytes.

## Bounds

Total recurrent state must remain 512 bytes in both arms. No exact recall, no state expansion, no routing search, no learned router, no threshold tuning, no result-informed retry, no production authority, no live activation.
