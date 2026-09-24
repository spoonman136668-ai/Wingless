# Wingless UP-64A — five-role zero-shot state-space scaling

Status: preregistered scientific cognition experiment.

Scientific parent: UP-63A seal `5f05e9300bcf8ce6319fd01d3f5a0499759061ae`.

## Question

UP-63A showed that the learned reversible program operators generalize across held-out program distributions through length 256 in the three-role substrate.

UP-64A asks a different question: are the learned reversible swap parameters tied to the 27-state training substrate, or can they be instantiated without retraining in a substantially larger compositional state space?

## Frozen design

- learn the existing two reversible angles only from the established three-role UP-56A single-step training task using the already-qualified 1,440-step optimizer;
- perform **no five-role training**;
- instantiate those learned angles in a five-role, three-value substrate with joint dimension 243;
- evaluate direct role-0 value swaps, all four adjacent role swaps, derived mutations of roles 1 through 4, and mixed programs at lengths 16 and 48;
- preserve the existing 0.98 category/aggregate cognition gate and numerical gates;
- evaluate a matched non-unitary control on the single-step five-role cases.

No five-role example can influence the learned parameters.

## Interpretation

A pass would show that the learned reversible parameter values transfer zero-shot across a 9x increase in joint state dimension and new role topology. A failure localizes the first state-space/topology scaling boundary without permitting retraining or threshold changes.
