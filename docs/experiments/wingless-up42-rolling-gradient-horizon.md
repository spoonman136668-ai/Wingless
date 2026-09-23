# Wingless UP-42: rolling full-rank gradient horizon calibration

Status: preregistered calibration branch only; not activated, promoted, or connected to ckb-plane.

Parent calibration: UP-41 sealed at `d5c622da08b3a13b7df3c4d3fd8ade45d2cc961a`.

Independent fixed target: accepted UP-35 capacity-18 geometry from seal `09284aa70203adf8881be42f67084c49b338dabe`.

## Why this experiment exists

UP-41 established three independent facts:

1. the current sequential+sticky optimizer did not enter the known-good fusion basin;
2. the exact-gradient/no-sticky reference also failed to enter that basin in 12 updates, so the horizon is too short;
3. the sequential estimator traveled materially less far than the exact-gradient reference, while sticky projection made cross-target irreversible fusions.

UP-42 removes sticky projection entirely and asks whether a computationally efficient gradient-memory repair can close the estimator gap when given a longer, preregistered horizon.

## Question

**Can a rolling full-rank reconstruction built from the same one plus/minus directional probe pair per update reach the known-good UP-35 fusion basin at the same tested horizon as an exact-gradient reference, while the original sequential rank-one estimator cannot?**

## Calibration-only boundary

UP-42 uses no task labels, tables, held-out data, observer fitting, architecture selection, or final capability evaluation.

The known UP-35 target is used only as a fixed geometry calibration target. No UP-42 arm may become a selected architecture.

## Frozen target and constraints

Target groups remain:

`[0,1,2,5]`, `[3]`, `[4]`

Target capacity remains 18.

All arms start from the same singleton offsets and preserve:

- zero mean;
- RMS 0.05;
- perturbation 0.002 where finite differences are used;
- learning rate 0.002;
- maximum coordinate update 0.004.

Sticky projection is disabled in every arm because UP-41 directly demonstrated cross-target irreversible fusion.

The existing 0.0025 fusion tolerance is used only to define whether an unfused trajectory has entered the target fusion basin.

## Frozen horizon sweep

Independent runs are performed at:

`12, 24, 48, 96`

This doubling schedule is fixed before Windows evidence is observed. No intermediate horizon is added after seeing results.

## Estimators

### Sequential projected central difference

The unchanged UP-38/39/40 one-direction estimator. One deterministic direction is measured and its rank-one projected gradient is immediately used.

### Rolling full-rank directional memory

Still uses exactly one deterministic plus/minus directional probe pair per update.

It stores at most the seven most recent directional derivatives. Once those observations span the five-dimensional zero-mean tangent space, it reconstructs the least-squares full tangent gradient. Before full rank is available, it falls back to the current rank-one update.

This adds directional memory, not more task evaluations per update.

### Analytic cosine reference

Uses the exact gradient of the geometry-only cosine objective. It remains a calibration reference, not a proposed Wingless optimizer.

## Measurements

For every estimator/horizon pair:

- initial, best, and final target cosine;
- final spread of the accepted fused target group;
- whether the 0.0025 target fusion basin is entered;
- first update at which it is entered;
- for the rolling estimator, how many updates use a full-rank reconstruction.

The key result is the **minimum tested horizon** at which each estimator enters the target basin; zero means no tested horizon succeeds.

## Frozen interpretation

- If rolling reconstruction enters the target basin at the same minimum tested horizon as the analytic reference, it is an efficient candidate repair for the directional-estimation defect.
- If rolling enters but needs a longer tested horizon than the analytic reference, directional memory helps but still leaves estimator loss.
- If sequential does not enter by the maximum tested horizon while rolling does, the old one-direction update is inadequate even when merely given more time.
- If neither rolling nor analytic enters by 96, the update rule/clamp/normalization geometry remains the primary limitation and estimator repair alone is insufficient.
- If all three enter at the same horizon, horizon—not estimator—is the dominant limitation.

No numerical criterion beyond the already-frozen target fusion tolerance is introduced.

## Plain speak

UP-41 showed the steering system has too little time, throws away most gradient information, and can glue the wrong coordinates together.

UP-42 removes the glue and gives the steering system memory.

Instead of forgetting each directional measurement after one move, it remembers the last seven measurements and reconstructs the full five-dimensional direction of travel once enough independent information exists. It still spends only one plus/minus probe pair per update.

The horizon sweep then tells us whether that cheap memory is enough to steer like the perfect-gradient reference, and how many updates the geometry actually needs.

## Authority boundary

UP-42 is geometry-only mathematical calibration. It does not invoke a language model, activate Wingless, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
