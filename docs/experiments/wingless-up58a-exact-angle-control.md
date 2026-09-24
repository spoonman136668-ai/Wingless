# Wingless UP-58A — exact-angle causal control

Status: preregistered scientific cognition diagnostic.

Scientific parent: UP-57A seal `cba417b76999151cc7e70a2da1d5afc577b26dfc`.

UP-57A located the doubly-controlled composition boundary between 16 and 24 operations while norm preservation and reversibility remained essentially exact. UP-58A asks whether the logical-depth failure is caused by the learned rotation angles being slightly short of the exact permutation angle.

Training is unchanged and performed once. The learned parameters are frozen. The same long programs at lengths 16, 24, 32, 48, 64, 96, and 128 are evaluated twice: once with the learned angles and once with the exact reversible angles [π/2, π/2].

This is a causal control, not a parameter search. If exact angles restore long-program accuracy, optimizer precision is implicated; if not, the composition structure itself remains limiting.


## Authoritative Windows result

Workflow run: `36033570518`

Runner: `WINGLESS-LINKDEADKB`

Source head: `ffa672a5b073d2369a8586405222055410112cdf`

Artifact: `10822813563`

Artifact digest: `sha256:c03f8e7a2f08937b54ad26f37a0bd608a57788f35533eae03f4482fb8e730f7c`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

The learned parameters were approximately `[1.568726759098322, 1.439276537385563]`; the exact control used `[π/2, π/2]`.

At every frozen length 16, 24, 32, 48, 64, 96, and 128, the exact-angle unitary achieved accuracy `1.0`. The learned-angle path reproduced the UP-57A degradation.

## Scientific classification

The long-composition failure in UP-56A/UP-57A is caused by accumulated gate-angle error, not by the controlled-operation topology or the unitary substrate. Exact reversible permutation angles restore perfect behavior through the longest tested program.

The next A-lane experiment maps the tolerance to frozen angle perturbations around π/2. This quantifies how accurately learned reversible gates must be identified for long-horizon composition.
