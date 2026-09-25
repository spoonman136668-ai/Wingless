# Wingless UP-LM0C — syntax-routed bounded recall language bridge

Status: preregistered scientific language/memory integration experiment.

Scientific parents:
- sealed UP-LM0B `413285fd458b87db03dde2667c59d73148c65b4a`;
- sealed UP-91B `9c2b5c8195bba9589afb488236c5ce6ec9271939`.

## Question

UP-LM0B showed strong next-byte language modeling but weak long-dependency value recall. UP-91B showed that a fixed 16-entry exact store can preserve sparse relevant information when routing is precise. UP-LM0C asks whether a purely local syntactic router can close the long-dependency byte prediction without attention or recurrent-state expansion.

## Frozen corpus and recurrent model

Exactly the UP-LM0B corpus, structural split, 64-float recurrent carrier, deterministic signed-permutation transport, deterministic byte embedding, linear softmax head, training order, 20 epochs, and learning rate 0.08.

No retraining difference is allowed between arms.

## Frozen arms

1. `recurrent64`
   - identical UP-LM0B inference.

2. `recurrent64_syntax_recall16`
   - same recurrent model;
   - exact recall capacity exactly 16 name->value bindings;
   - parser reacts only to already-observed local byte patterns:
     - `<name> stores <value>.` writes the observed binding;
     - `<name> reports ` queries that name immediately before the dependent value begins;
   - memory may affect only the first dependent value byte;
   - all other next-byte probabilities remain recurrent-model probabilities;
   - FIFO replacement only if the 16-entry cap is exceeded;
   - no access to future bytes or held-out labels.

## Evaluation

For train, held-out recombination, and held-out four-sentence streams report:
- overall next-byte top-1 accuracy;
- cross-entropy/perplexity;
- dependent first-byte accuracy;
- dependent first-byte cross-entropy/perplexity;
- maximum exact-recall entries used.

## Bounds

No attention, no recurrent-state expansion, no learned router, no future-query oracle, no memory beyond 16 bindings, no hyperparameter search, no result-informed retry, no live activation, no production authority.
