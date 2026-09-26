# Wingless UP-161B — fixed balanced cleanup schedule

Status: preregistered scientific lexical-stability intervention experiment.

Scientific parent: sealed UP-160B d4a12a5ecd2d5f24671d1af3a2a7267732f10fdd.

## Question

UP-159B showed that cleanup-order preference is subject-dependent, and UP-160B proved the chained cleanup updates are noncommutative. Can a fixed, nonadaptive schedule that rotates class position across terminal epochs reduce subject-specific order bias without adding updates?

## Frozen starting point

Reuse UP-159B/160B unchanged:
- exact accepted 15-epoch common prefix;
- three matched two-subject new-family corpora:
  - mia/noah;
  - opal/pax;
  - quin/rue;
- each terminal epoch contains exactly 24 new-family updates and 15 old-rehearsal updates;
- old rehearsal indices 10,11,12 occur before REPORT;
- fixed target-subject four-example REPORT block;
- fixed post-REPORT cleanup composition:
  - 5 STORE;
  - 5 OBSERVE;
  - 2 REPORT;
- learning rate remains 0.08;
- five terminal epochs only.

## Frozen schedule arms

1. canonical_fixed
   - every terminal epoch: STORE_OBSERVE_REPORT.

2. balanced_forward
   - epochs 15..19:
     - STORE_OBSERVE_REPORT
     - OBSERVE_REPORT_STORE
     - REPORT_STORE_OBSERVE
     - STORE_OBSERVE_REPORT
     - OBSERVE_REPORT_STORE

3. balanced_reverse
   - epochs 15..19:
     - STORE_REPORT_OBSERVE
     - REPORT_OBSERVE_STORE
     - OBSERVE_STORE_REPORT
     - STORE_REPORT_OBSERVE
     - REPORT_OBSERVE_STORE

Because five terminal epochs cannot divide three cleanup positions exactly evenly, each balanced arm is the closest deterministic cyclic allocation: every class visits every position, with two positions repeated once. The two balanced arms are complementary and both are preregistered.

No schedule is selected by subject or by observed performance.

## Measurements

Per subject × schedule:
- exact epoch-by-epoch cleanup orders;
- mean new-prefix effect;
- mean pre-REPORT recovery;
- mean REPORT damage;
- mean post-REPORT cleanup recovery;
- mean net epoch change;
- final mean old retention;
- held-out new-name accuracy.

Per schedule:
- mean old retention across six subjects;
- minimum old retention across six subjects;
- subject-to-subject retention spread.

## Interpretation

If either balanced schedule reduces cross-subject spread or raises the minimum subject retention without extra updates, fixed position rotation can average away part of the noncommutative context bias. If not, context sensitivity cannot be neutralized by simple nonadaptive cycling.

No adaptive schedule, post-hoc selection, or numeric threshold is introduced.

## Bounds

No extra updates, no extra epochs, no rehearsal coverage change, no class-composition change, no memory change, no architecture change, no adaptive ordering, no result-informed retry, no live activation, no production authority.
