# Wingless UP-1: trained unitary propagator

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-0 Windows proof at source head `2167d5cf255a0f139a17959c2cfeaa888fb53484`.

Branch baseline: UP-0 research seal `9cd058dd43bd042d0dc03df4a35b69ecbd5f36e8`.

## Question

Can the interference-routing behavior demonstrated by UP-0 be learned from examples rather than supplied as a hand-selected coupling angle, while retaining norm preservation and reversibility?

## Construction

UP-1 keeps the same bounded two-coordinate complex state and unitary Givens rotation used by UP-0, but initializes the trainable rotation angle at zero.

Training uses:

- four deterministic training states;
- one trainable angle;
- central-difference numerical gradients;
- fixed learning rate 0.15;
- gradient epsilon 1e-6;
- exactly 20 optimization steps;
- no external ML or autodiff dependency.

The target angle is not supplied to the optimizer. It is retained only as a post-training diagnostic reference.

Training examples use global phases 0 and pi/2. Held-out examples use unseen global phases 0.37 and -1.11. The task remains binary relative-phase routing: in-phase states route to coordinate 1 and opposed-phase states route to coordinate 0.

## Pass conditions

- initial held-out accuracy is exactly 0.5;
- train accuracy reaches 1.0;
- held-out accuracy reaches 1.0;
- final mean negative log likelihood is <= 1e-5;
- learned angle is within 1e-3 radians of pi/4;
- training loss is monotonic non-increasing;
- maximum norm drift is <= 1e-12;
- maximum forward/inverse round-trip error is <= 1e-12;
- repeated complete training results are deterministic;
- existing repository tests remain green after rebuilding advisory ICE.

## Interpretation boundary

A pass would establish only that a tiny norm-preserving interference primitive can be optimized from data and generalize over a deliberately controlled global-phase nuisance variable.

It would not establish language reasoning, quantum advantage, superior parameter efficiency, or frontier-model capability.

A positive UP-1 result permits the next research step: multi-parameter learned propagation over a higher-dimensional latent state with a matched non-unitary baseline.

## Authority boundary

UP-1 does not register an inference backend, invoke a language model, activate a worker/listener, execute model-selected tools, or alter ckb-plane authority. Queue, retry, workspace, acceptance, and promotion authority remain outside Wingless.
