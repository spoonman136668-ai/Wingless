# Wingless UP-59A — gate-angle component attribution

Status: preregistered scientific cognition diagnostic.

Scientific parent: UP-58A seal `0bf110e1938fdde619364e26fb73f98e63e4a6e0`.

UP-58A proved exact π/2 angles restore perfect long-horizon composition. UP-59A isolates which learned angle causes the accumulated logical error.

Training is unchanged and performed once. Four frozen parameter arms are evaluated at program lengths 16, 24, 48, and 128:

- both learned angles;
- exact role-swap angle + learned doubly-controlled angle;
- learned role-swap angle + exact doubly-controlled angle;
- both exact π/2.

The existing 0.98 cognition accuracy and numerical gates remain unchanged. This is component attribution, not optimization or parameter search.


## Authoritative Windows result

Workflow run: `36034045166`

Runner: `WINGLESS-LINKDEADKB`

Source head: `5a7299b86fe250bb28b1e1b6af674322c694ca97`

Artifact: `10823268842`

Artifact digest: `sha256:b1885ba37665955a2bbea126846c748c221dffb47847f6b72e980876059cf330`

Results at length 128:

- both learned: `0.691358024691358`;
- exact role-swap + learned double-control: `0.7037037037037037`;
- learned role-swap + exact double-control: `1.0`;
- both exact: `1.0`.

The learned-role-swap/exact-double-control arm also remains perfect at lengths 16, 24, and 48.

## Scientific classification

Accumulated cognition error is attributable to the learned doubly-controlled rotation angle, not the small role-swap angle error. The reversible topology itself remains exact.

The next A-lane experiment holds every other variable fixed and maps a preregistered perturbation ladder around the exact π/2 double-control angle.
