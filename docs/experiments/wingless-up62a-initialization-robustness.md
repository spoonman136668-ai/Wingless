# Wingless UP-62A — initialization robustness at confirmed training budget

Status: preregistered scientific cognition confirmation.

Scientific parent: UP-61A seal `e3b319607b677dd7e9537d9651767c013e50d07e`.

UP-61A showed that the unchanged learner reaches the long-horizon precision band by 1440 steps from initialization [0.1,0.1]. UP-62A tests whether that conclusion depends on that initialization.

The exact same 1440-step optimizer, learning rate 0.10, epsilon 1e-6, task, and gates are run from five frozen initializations: [0.05,0.05], [0.1,0.1], [0.25,0.25], [0.5,0.5], and [-0.1,0.1]. Every learned model is evaluated at mixed-program length 128.

No adaptive restart, learning-rate change, or result-informed initialization is permitted.
