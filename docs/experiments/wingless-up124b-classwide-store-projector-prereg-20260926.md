# Wingless UP-124B — class-wide STORE competitor-nullspace projector

Status: preregistered scientific representation-generalization experiment.

Scientific parent: sealed UP-123B b3f33eae8b7898c296cde750ac42553eac2bd834.

## Question

UP-120B and UP-123B showed that hand-selected per-verb competitor-orthogonal corrections can independently fix stores and archives. Can one frozen class-wide linear geometry rule replace those per-verb corrections while preserving all STORE surfaces?

## Frozen class-wide transform

Build one fixed projector before training from the frozen encoder only:
- OBSERVE competitor direction = centroid of observes, sees, notes, watches, notices, spots;
- REPORT competitor direction = centroid of reports, recalls, tells, remembers, recounts, retells;
- form the two-dimensional competitor span from those two directions;
- for every STORE surface representation, project out that competitor span;
- rescale the residual to the original representation norm.

The same transform is applied to every STORE verb at both training and inference.

No per-verb correction coefficient is computed.

## Frozen comparison arms

1. local_corrections
   - exact UP-123B stores + archives frozen local corrections.

2. classwide_store_projector
   - no stores-specific or archives-specific delta;
   - apply the single class-wide projector to every STORE surface:
     stores, keeps, keepsafe, archives.

OBSERVE and REPORT representations remain unchanged.

## Frozen training/evaluation

- 64-D encoder;
- learning rate 0.08;
- 20 epochs per acquisition batch;
- fixed replay budget 6;
- alternate six-verb acquisition lexicon;
- all six class orders;
- current_class_excluded and stage3_anchor2 policies;
- train subjects ada, ben, cy, dee;
- evaluate original heldout eli/fay and unseen heldout kia/leo.

## Metrics

Across every stage:
- stores / keeps / keepsafe / archives accuracy on original and unseen subjects;
- base-lexicon accuracy original/unseen;
- acquired aggregate accuracy;
- STORE-vs-strongest-competitor margins for every STORE verb.

## Interpretation

If the class-wide projector matches the local-correction control, the B lane can replace manually identified verb repairs with one reusable class geometry operator. Degradation would show that the useful correction remains surface-specific.

## Bounds

No adaptive geometry, no result-informed recomputation, no replay/training change, no threshold tuning, no state expansion, no live activation, no production authority.
