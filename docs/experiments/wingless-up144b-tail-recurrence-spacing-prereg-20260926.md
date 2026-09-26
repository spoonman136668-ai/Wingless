# Wingless UP-144B — tail recurrence spacing at fixed coverage

Status: preregistered scientific lexical sequencing experiment.

Scientific parent: sealed UP-143B e52125cd6b57be5564ec797327ff5819d822cb04.

## Question

UP-143B localized the original-family retention effect to repeated post-anchor tail damage. If tail-identity coverage and frequency are held exactly fixed, does temporal spacing of those tail identities control the unrecovered damage?

## Frozen starting point

Use the exact original lexical family and pre-adaptation gate from UP-143B:
- 64-D state;
- 128-D dual-view gate input;
- frozen STORE competitor-nullspace projector;
- exact old-family training history.

## Frozen adaptation budget

All arms:
- 20 epochs;
- learning rate 0.08;
- exactly 24 original-family updates per epoch;
- exactly 15 balanced old-family rehearsal updates per epoch;
- exactly 4 new-family updates after the old anchor;
- every original-family example exactly once per epoch;
- exactly four rotation shifts: {0, 6, 12, 18};
- each shift used exactly five times across 20 epochs.

## Arms

1. interleaved
   - epoch shifts cycle 0, 6, 12, 18 and repeat five times.

2. blocked
   - shift 0 for epochs 1–5;
   - shift 6 for epochs 6–10;
   - shift 12 for epochs 11–15;
   - shift 18 for epochs 16–20.

The multiset of tail identities and per-identity frequency is identical. Only recurrence spacing differs.

## Measurements

For every epoch:
- old retention before epoch;
- old retention after first 20 new updates;
- old retention after old-family anchor;
- old retention after final four new updates;
- prefix damage;
- anchor recovery;
- tail damage;
- net epoch change.

Final:
- old held-out / unseen retention;
- primary / secondary new-family accuracy;
- mean new-family accuracy.

## Interpretation

Better retention and smaller tail damage under interleaving would show that temporal spacing, not just coverage breadth, is a causal part of the tail-diversity mechanism. Similar outcomes would indicate that coverage/frequency dominates over spacing.

No post-result numeric threshold is introduced; exact measured differences are the result.

## Bounds

No adaptive shift choice, no extra updates, no new memory, no projector recomputation, no architecture change, no threshold search, no result-informed retry, no live activation, no production authority.
