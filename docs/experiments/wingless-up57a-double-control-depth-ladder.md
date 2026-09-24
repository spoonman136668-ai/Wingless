# Wingless UP-57A — doubly-controlled composition depth ladder

Status: preregistered scientific cognition diagnostic.

Scientific parent: UP-56A seal `fd1731dba6b0f8c73ef5173143c52de1ad435a98`.

UP-56A transfers its doubly-controlled primitive perfectly but misses the frozen long-program gate. UP-57A does not change training. It learns the exact same primitive once, freezes the learned parameters, and evaluates mixed-program lengths 2, 4, 8, 12, 16, 24, 32, 48, 64, 96, and 128.

Each point uses the unchanged 0.98 cognition accuracy gate and numerical gates. The experiment maps the composition-depth boundary; it does not tune it.


## Authoritative Windows result

Workflow run: `36033149929`

Runner: `WINGLESS-LINKDEADKB`

Source head: `7bae02b7c5114c05caceb49517db984cca377cbd`

Artifact: `10822798019`

Artifact digest: `sha256:717e8888c6c1e4c77f977825538831d70d4291d3591cd4e4ea21cda2d4398882`

Frozen unitary depth result:

- lengths 2, 4, 8, 12, 16: accuracy `1.0`, PASS;
- length 24: `0.9506172839506173`, FAIL;
- length 32: `0.9506172839506173`;
- length 48: `0.8888888888888888`;
- length 64: `0.7160493827160493`;
- length 96: `0.7654320987654321`;
- length 128: `0.691358024691358`.

Norm drift and round-trip error remain on the order of 1e-15.

## Scientific classification

The doubly-controlled composition boundary is logical/parameteric rather than loss of unitarity: exact reversibility remains intact while repeated imperfect gate rotations accumulate task error. The next A-lane causal control freezes the task and compares the learned angles against exact π/2 rotations without retraining.
