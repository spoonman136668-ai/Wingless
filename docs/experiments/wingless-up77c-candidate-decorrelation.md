# Wingless UP-77C — candidate decorrelation on frozen 8×8 carrier geometry

Status: preregistered scientific routing experiment.

Scientific parents:
- sealed UP-76C `194f014c9f3b2e7769905b9237229c2112d613af`;
- sealed UP-75C `478c60bef8d4104600043e875faeb8de1ea3e409`.

## Question

UP-76C froze 8×8 as the best tested fixed partition geometry at a constant 512-byte carrier budget. UP-75C showed that the earlier two-choice router was worse than fixed assignment. One plausible cause is candidate structure: the two candidates were confined to correlated local pairs.

Does decorrelating the second candidate rescue norm-based two-choice routing on the frozen 8×8 geometry?

## Frozen carrier budget

All arms:
- 64 float64 recurrent values = 512 bytes;
- eight 8-dimensional partitions;
- deterministic norm-1 bipolar associative vectors;
- gated correction threshold 0.15;
- no exact recall;
- no routing table;
- no persistent metadata.

## Frozen arms

1. `fixed_mod8`
   - partition = key mod 8.

2. `adjacent_pair`
   - candidate A = key mod 8;
   - candidate B = A xor 1;
   - new insert chooses the lower-L2-norm candidate;
   - overwrite/query inspect both candidates and use the stronger decoded match.

3. `antipodal_pair`
   - candidate A = key mod 8;
   - candidate B = (A + 4) mod 8;
   - same norm-based insertion and dual-candidate overwrite/query.

4. `mixed_second_hash`
   - candidate A = key mod 8;
   - candidate B = high three bits of `sq0Mix64(key * frozen_constant)` mod 8;
   - if B == A, use (B + 3) mod 8;
   - same norm-based insertion and dual-candidate overwrite/query.

No result-dependent candidate selection is allowed.

## Frozen tasks

High-load subset:

Multi-bank interference:
- logical banks 6 and 8;
- load 50% and 100%;
- four entities per logical bank;
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
Seeds: 137M and 138M.

## Metrics

- primary value/query accuracy;
- exact-episode accuracy;
- cross-bank corruption rate;
- recurrent-state bytes.

## Interpretation

If antipodal or independently mixed candidates beat fixed_mod8, candidate correlation is supported as the C75 failure mechanism. If all two-choice arms remain worse than fixed_mod8, the routing instability itself—not just candidate correlation—is the stronger diagnosis.

## Bounds

No state expansion, exact recall, routing table, learned router, threshold tuning, result-informed retry, production authority, or live activation.
