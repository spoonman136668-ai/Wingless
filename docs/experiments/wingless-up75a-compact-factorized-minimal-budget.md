# Wingless UP-75A — compact factorized minimal-budget boundary

Status: preregistered scientific efficiency-boundary experiment.

Scientific parent: sealed UP-74A `04cd3f9b2e77f706e624790902fcd8bc022d104b`.

UP-74A retained perfect held-out accuracy at 25 optimization steps. UP-75A freezes a lower ladder of 1, 2, 5, 10, and 25 steps.

Everything else is unchanged: the same 54 training states, 162 held-out states, raw factorized one-hot and role-factorized simplex representations, five independent 3-class linear-softmax heads, zero initialization, learning rate 1.0, and 0.98 per-role held-out gate.

No adaptive stopping, learning-rate change, threshold change, added training state, representation change, selection, promotion, or activation is permitted. Passing, failure, or non-monotonicity is valid evidence and must be sealed as observed.
