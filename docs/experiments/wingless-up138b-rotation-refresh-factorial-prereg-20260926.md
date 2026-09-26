# Wingless UP-138B — rotation × final-refresh factorial

Status: preregistered scientific adaptation-order interaction experiment.

Scientific parent: sealed UP-137B 34eb0e5cf6432290e5cbf0a01c91db8e3f48c41c.

## Question

UP-136B showed that a fixed four-update final refresh is harmful under fixed new-example order. UP-137B showed that epoch rotation alone is not beneficial. UP-135B nevertheless found a strong internal arm when rotation and a four-update refresh occurred together. Is that result a genuine interaction between rotation and refresh placement?

## Frozen mechanism and budget

- exact dual-view 128-D lexical gate and frozen STORE projector;
- new grounding subjects mia/noah;
- balanced old rehearsal subjects ada/ben;
- 20 epochs at lr 0.08;
- exactly 24 new-family updates per epoch;
- exactly 15 old-family updates per epoch;
- every example used exactly once;
- no extra updates.

## 2×2 factorial

Factor A — new-example order:
1. fixed: exact UP-136B fixed 24-example order.
2. rotated: rotate the same 24-example list left by epoch modulo 24.

Factor B — final new refresh:
1. refresh0: all 24 new updates before all 15 old updates.
2. refresh4: first 20 new, all 15 old, final 4 new.

Arms:
- fixed_refresh0
- fixed_refresh4
- rotated_refresh0
- rotated_refresh4

The same ordered list for a given factor-A arm is split at position 20 for refresh4; no duplication occurs.

## Evaluation

- primary and secondary new-family accuracy;
- worst new-surface accuracy;
- old held-out and old unseen class accuracy;
- old STORE precision/recall.

Also report factorial contrasts for each aggregate accuracy metric:
- rotation main effect;
- refresh main effect;
- rotation×refresh interaction contrast.

## Interpretation

A positive interaction contrast with weak/negative individual main effects would directly explain the UP-135B/136B/137B pattern. Absence of the interaction would mean the earlier UP-135B result was due to another order detail not yet isolated.

## Bounds

No extra updates, no adaptive ordering, no threshold search, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
