# Wingless UP-7: direct transported-state observer

Status: Windows-qualified negative research result; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-6 Windows qualification sealed at `c431bf27e021e9aa83f99155bbb597d1a5e24434`.

## Question

Can a learned observer read mutable memory directly from the transported latent state without an explicit inverse, a runtime prototype table, or knowledge of transport depth?

This is intentionally a falsifiable question. A scientifically negative result is valid evidence and must not be reported as an acceptance-artifact failure.

## Methodological change

UP-7 separates two concepts:

1. **harness qualification** — code compiles, results are finite and deterministic, frozen Wingless regression stays green;
2. **scientific hypothesis gate** — direct observation reaches predefined generalization thresholds.

The Windows acceptance script fails only for harness/invariant defects. It prints the scientific gate independently.

## Observer input

The observer receives:

- the 16 normalized transported coordinate probabilities;
- a four-way one-hot entity selector.

It does **not** receive:

- the transport depth;
- an inverse-transformed state;
- a memory-table prototype;
- a table ID.

The unitary and matched non-unitary paths receive identical observer architecture and optimization budget, but each observer is trained on its own transport outputs.

## Data split

All 256 possible four-entity memory tables are partitioned deterministically:

- 128 training tables;
- 128 disjoint held-out tables.

Training depths:

`8, 24, 72, 216, 432, 648`

Held-out depths:

`32, 128, 512, 1024`

Held-out depths never appear in observer training.

Both train and held-out states receive deterministic bounded complex perturbation with amplitude 0.05.

## Learned direct observer

A shared four-class linear-softmax observer uses 20 features:

- 16 transported probabilities;
- 4 entity-selector features.

The same observer parameters are reused for all four entities.

The observer is trained for 500 deterministic full-batch gradient steps with learning rate 1.0.

A shared learned relation head from the UP-6 design maps the two predicted value distributions to one of four relation classes.

## Integration workload

After static held-out evaluation, the direct observer is placed inside 32 unseen mutable-memory programs with 12 writes each.

Every transport gap uses only held-out depths.

The explicit irreversible overwrite/re-encode boundary remains, but the state is decoded directly from its forward transported representation.

## Scientific hypothesis gate

The direct-observation hypothesis is considered positive for a path only when all of the following hold:

- static held-out value accuracy >= 0.95;
- integration commit decode accuracy >= 0.95;
- exact final-table accuracy >= 0.90;
- relational-query accuracy >= 0.90.

The matched control is evaluated by the identical thresholds. It is not required to fail.

## Harness qualification conditions

- no runtime prototype lookup;
- no explicit inverse readout;
- no depth feature;
- 128/128 disjoint table split;
- train and held-out depth sets exact and disjoint;
- relation head trains successfully;
- every reported metric is finite;
- complete output is deterministic;
- full existing Wingless regression remains green after advisory ICE rebuild.

## Interpretation boundary

If the unitary hypothesis passes, UP-7 establishes that a small learned observer can directly decode unseen memories at unseen propagation depths without being handed the inverse or the depth.

If it fails, that is still useful: it indicates the transported state needs an explicit clock/frame signal, a more expressive observer, or a different invariant representation before direct readout is viable.

Either outcome informs UP-8. No language-reasoning claim follows from UP-7.

## Authority boundary

UP-7 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-22 against source head `8f201ecca5d680020dba340464464aa5b2ec6a49` passed the harness and complete Wingless regression, while both scientific hypotheses failed.

The unitary direct observer reached 0.4217122395833333 training accuracy and 0.205078125 static held-out accuracy. In mutable-memory integration it reached 0.005208333333333333 commit accuracy, 0 exact final-table accuracy, and 0.75 relational-query accuracy. The unitary transport itself remained bounded with maximum norm drift 7.571721027943568e-14.

The matched control showed essentially the same observer failure while its transport norm drift reached 4714800782374.951.

The primary next diagnosis is representation loss: UP-7 exposed only coordinate probabilities, discarding relative phase/coherence. UP-8 therefore tests phase-aware, global-phase-invariant observation before introducing an explicit clock/depth signal or a larger nonlinear observer.


## Post-qualification confound discovery

The original 128/128 split used `memoryTableIndex(table) % 2`. Because the table index is base-4 and its parity is determined by the final entity value, this split unintentionally gave entity 3 only values 0/2 during observer training and only values 1/3 during held-out evaluation.

Therefore UP-7 remains a valid negative result for the exact tested harness, but it is **not a clean causal test** of phase loss or depth ambiguity. The held-out set also contains unseen marginal labels for one entity.

UP-8 corrects this before drawing architectural conclusions: it uses a balanced combinatorial split in which every entity/value marginal appears in both partitions, verifies that property explicitly, and then compares magnitude-only and phase/coherence-aware direct observers under identical data and transport conditions.
