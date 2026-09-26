# Wingless UP-LM2K — dependency closure boundary

Status: preregistered scientific fixed-capacity language experiment.

Scientific parent: sealed UP-LM2J f7d186efd394989b9cc8b98a06af0cc394b0e49b.

## Question

UP-LM2J showed that 24 entities can be processed exactly with cap16 when dependencies close inside working sets of at most 16. How many unfinished first-chunk dependencies may cross a 12+12 chunk boundary before exact recall fails?

## Frozen model and memory

Reuse UP-LM2J unchanged:
- same trained five-family language model;
- same router;
- exact recall cap = 16;
- 24 fixed entity names;
- same STORE, OBSERVE, REPORT event multiset;
- no extra training;
- no memory increase.

## Frozen two-chunk organization

Chunk A = entities 0..11.
Chunk B = entities 12..23.

For every arm:
1. STORE all 12 A entities.
2. OBSERVE all 12 A entities.
3. REPORT the preregistered early-closure subset of A.
4. STORE all 12 B entities.
5. OBSERVE all 12 B entities.
6. REPORT the deferred A entities.
7. REPORT all 12 B entities.

The event multiset is identical across arms; only REPORT placement changes.

## Deferred dependency arms

Number of A reports deferred across the B STORE boundary:
- defer0;
- defer2;
- defer4;
- defer6;
- defer8;
- defer12.

Deferred A reports are always the latest A entities, so under FIFO16 they are the newest A dependencies.

Before B is stored, A contributes 12 live memory entries. Storing B adds 12 more. With cap16, exactly eight oldest A entries are evicted. Therefore:
- defer <= 4 should remain exact;
- defer > 4 should lose exactly defer-4 unresolved A dependencies.

Frozen expected recall-hit rate:
- defer0 = 1;
- defer2 = 1;
- defer4 = 1;
- defer6 = 22/24;
- defer8 = 20/24;
- defer12 = 16/24.

These expectations are preregistered and do not control harness acceptance.

## Measurements

Per allocation × training schedule × family × deferral arm:
- top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- report-set exact accuracy;
- recall hit rate;
- preregistered expected recall hit rate;
- event/STORE/OBSERVE/REPORT routing accuracy;
- maximum recall entries.

## Interpretation

A sharp exactness boundary at four unfinished cross-chunk dependencies would show that working-set scheduling succeeds only when unresolved dependencies plus the next chunk fit within cap16. Deviations define a more complex closure boundary.

## Bounds

No capacity change, no adaptive chunking, no extra training, no router change, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
