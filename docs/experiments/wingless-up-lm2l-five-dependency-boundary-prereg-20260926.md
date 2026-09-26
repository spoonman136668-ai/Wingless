# Wingless UP-LM2L — five-dependency boundary

Status: preregistered scientific fixed-capacity language experiment.

Scientific parent: sealed UP-LM2K bc6469d94e4b71d79ceaa85a55aeb6ab905b8464.

## Question

UP-LM2K established exact processing with four deferred first-chunk dependencies and failure with six. Is five deferred dependencies the first failing level, independent of entity identity assignment and deferred-report order?

## Frozen system

Reuse UP-LM2K unchanged:
- 24 entities;
- two chunks of 12;
- exact recall cap 16;
- same five lexical families;
- same five fixed-mass allocations;
- canonical_prior and balanced_prior;
- no recurrent retraining, router change, attention, or future oracle.

## Frozen boundary levels

Test deferred first-chunk reports:
- 4;
- 5;
- 6.

## Frozen identity controls

Use two deterministic cyclic identity rotations over the same 24 accepted entity names:
- rotation 0;
- rotation 7.

Each identity keeps its original value association. The full STORE/OBSERVE/REPORT event multiset is unchanged.

## Frozen deferred-report orders

For the deferred first-chunk set:
- forward;
- reverse.

The same deferred positions are used; only report order changes.

## Preregistered capacity expectation

For every identity rotation and report order:
- D=4 => 24/24 recall hits;
- D=5 => 23/24 recall hits;
- D=6 => 22/24 recall hits.

These expectations do not control harness acceptance.

## Measurements

Per allocation × training schedule × family × identity rotation × deferred-report order × D:
- recall hit rate;
- expected hit rate;
- dependent first-byte accuracy;
- report-set exactness;
- event and class routing;
- max recall entries.

## Interpretation

A consistent first failure at D=5 establishes the direct dependency-closure boundary. Any variation by identity rotation or deferred-report order defines a sequence/identity interaction instead.

## Bounds

No capacity increase, extra training, adaptive chunking, semantic priority, router change, attention, future oracle, result-informed retry, live activation, or production authority.
