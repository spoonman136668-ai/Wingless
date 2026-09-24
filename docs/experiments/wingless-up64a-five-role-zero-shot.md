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


## Authoritative Windows result

Workflow run: `36048902496`

Runner: `WINGLESS-LINKDEADKB`

Source head: `d3b402102dc67983831f2a96bc5c6fdcd5e2f3eb`

Artifact: `10829423559`

Artifact digest: `sha256:5cf911f16409a3fd6c7b4ebdb6a68408c844665f0a02c23f4662abb1a9eccc4a`

Results:

- aggregate unitary accuracy: `1.0`;
- role-0 value swaps: `1.0`;
- adjacent swaps 0/1, 1/2, 2/3, 3/4: `1.0` each;
- derived mutations of roles 1–4: `1.0` each;
- length-16 programs: `1.0`;
- length-48 programs: `1.0`;
- matched non-unitary single-step control: `0.5705003248862898`;
- frozen gate: `PASS`.

## Scientific classification

The learned reversible swap angles are not tied to the 27-state training substrate. With no five-role training, they transfer across a 9x increase in joint state dimension and new role topology.

The next A-lane experiment moves from operator/state-space scaling to representation structure: train a readout on sparse logical states from distributed encodings and test unseen combinations, with a factorized encoding control against an entangled encoding.
