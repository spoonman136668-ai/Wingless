# Wingless UP-145C — lexical consolidation integration

Status: preregistered scientific integration experiment.

Scientific parent: sealed UP-144C 7bd71e089b99a4e24e958466f3d06603613c0972.

## Question

The isolated C-line now has a validated two-stage consolidation rule, novelty-pressure clock, finite history lifetime, and age-ranked replacement. Do those mechanisms remain coherent when the admission machine carries actual accepted lexical subject/verb/class associations while durable known facts, temporary recurring evidence, and one-shot novelty are interleaved in the same bounded stream?

## Frozen lexical catalog

Use accepted lexical families from the B-line without training or modifying the B gate.

Class encoding:
- STORE = 0;
- OBSERVE = 1;
- REPORT = 2.

Initial exact-memory lexical facts (16 total):
- quin × {lodges, stashes, caches, files, scans, checks, views, monitors, relays, announces, cites, summarizes};
- rue × {lodges, stashes, caches, files}.

The first 12 `quin` facts are the hot durable set and are queried repeatedly. The four `rue` STORE facts are cold durable controls.

Candidate lexical associations:
- mia caches → STORE;
- noah scans → OBSERVE;
- opal relays → REPORT;
- pax files → STORE.

No semantic priority is given to any class or identity.

## Three 32-candidate generations

Candidate calls advance the frozen consolidation counter. After every four candidate calls, all 12 hot durable facts are queried; these read-only lexical queries do not advance admission time.

Four temporal roles are embedded simultaneously:

1. two_stage
   - pair in G1 and pair in G2;
   - expected to build history then admit.

2. cross_boundary
   - one sighting at the end of G1 and one at the start of G2;
   - expected not to qualify.

3. delayed_two_stage
   - pair in G1, no G2 occurrence, pair in G3;
   - tests use of a still-live aged history trace.

4. isolated
   - one sighting in each generation;
   - expected not to accumulate.

All unscheduled candidate positions are unique one-shot lexical-novelty keys with class values cycling STORE/OBSERVE/REPORT.

## Mirrored role assignments

A:
- mia caches = two_stage;
- noah scans = cross_boundary;
- opal relays = delayed_two_stage;
- pax files = isolated.

B:
- opal relays = two_stage;
- pax files = cross_boundary;
- mia caches = delayed_two_stage;
- noah scans = isolated.

## Measurements

Per assignment and candidate:
- lexical subject/verb/class;
- temporal role and positions;
- admission count;
- final exact-memory retention.

Also report:
- hot durable lexical accuracy;
- cold durable lexical accuracy;
- exact recall entries used;
- current/history maxima;
- false one-shot admissions;
- panic rate.

## Interpretation

Role-dependent outcomes that follow timing rather than lexical identity/class demonstrate that the validated bounded consolidation rule survives integration into a mixed lexical association stream. Loss of hot durable facts or class-specific deviations define integration boundaries.

No threshold is tuned post-result.

## Bounds

No memory increase, no semantic/query priority, no B-gate training or mutation, no future oracle, no adaptive schedule, no result-informed retry, no live activation, no production authority.
