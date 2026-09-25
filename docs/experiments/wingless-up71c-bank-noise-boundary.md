# Wingless UP-71C — golden-tag bank/noise boundary localization

Status: preregistered scientific state-capacity experiment.

Scientific parent: sealed UP-70C Windows evidence `52af265b43b4074e62bdc3836ab68c7d90e58452`.

## Question

UP-70C passed all 6/7/8-bank conditions through noise 0.008, while 10 banks passed at zero noise and collapsed to exact-scenario accuracy 0.5 at noise 0.004 and 0.008 on both schedules. UP-71C localizes the bank-count/noise interaction between the robust 8-bank regime and the fragile 10-bank regime.

## Frozen design

Unchanged: state dimension 16; 256 deterministic scenarios; four entities per bank; existing golden-rotation phase tags; coherence decoder and magnitude control; reversible transport depth 64; value gate 0.99; exact-scenario gate 0.95.

Schedules: 113M, 114M.
Bank counts: 9, 10.
Memory noise: 0, 0.001, 0.002, 0.003, 0.004, 0.008.

No tag optimization, state-dimension increase, decoder change, threshold change, adaptive stopping, schedule search, bank selection, promotion, or activation is permitted.

## Interpretation

The experiment may identify a 9-bank capacity boundary, a low-noise 10-bank transition, schedule sensitivity, or continued success/failure. All points are sealed exactly as observed.
