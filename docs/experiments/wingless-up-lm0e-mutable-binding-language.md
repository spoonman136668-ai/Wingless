# Wingless UP-LM0E — mutable language bindings

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0D `c44b28d8feb60320f0ed5c15548b26b2da727c1a`.

## Question

UP-LM0D showed perfect four-binding delayed language recall with the fixed 16-entry syntax-routed store. Can the same architecture preserve latest-value-wins semantics when language facts are explicitly updated before later reports?

## Frozen corpus

Names:
- ada
- ben
- cy
- dee
- eli
- fay

Values:
- amber
- cobalt
- ivory
- jade
- mauve
- silver

Distractors:
- birds sing.
- rain falls slowly.
- lamps glow at dusk.
- quiet winds cross hills.

Each paragraph contains four initial store clauses, one distractor, a frozen number of update clauses, another distractor, then four report clauses.

For each `name_offset n`, `value_offset v`, and distractor index `d`:

- bindings use names `n,n+1,n+2,n+3 mod 6`;
- initial values use `v,v+1,v+2,v+3 mod 6`;
- update values are each initial value advanced by two positions mod 6;
- update count is determined by `d`: 0, 1, 2, or 4;
- the updated names are always the first `update_count` names in binding order;
- report order is rotated by `d`;
- each report target is the latest stored value for that name.

There are 144 paragraphs.

Structural split:
- train iff `(n + 2*v + d) mod 3 != 2`;
- held out otherwise.

Every name, value, update count, and distractor occurs in training.

## Frozen recurrent model

Identical to UP-LM0D:
- 64 float64 recurrent-state values;
- same deterministic signed-permutation transport;
- same deterministic bipolar byte embedding;
- same recurrent update;
- same linear softmax output;
- 20 epochs;
- learning rate 0.08;
- no attention;
- no pretrained weights.

## Frozen arms

1. `recurrent64`
2. `recurrent64_syntax_recall16`
   - same local `stores` / `reports` parser;
   - an update uses the same `<name> stores <value>.` syntax and overwrites the existing binding;
   - exact recall capacity remains exactly 16 entries;
   - memory may override only the first byte of a reported value;
   - no future-byte access.

## Evaluation

Train, held-out, and held-out stream4:
- overall next-byte top-1 accuracy;
- perplexity;
- dependent first-byte accuracy across all four report queries;
- whole-paragraph four-query exact accuracy;
- accuracy partitioned by update count 0/1/2/4;
- maximum exact-recall entries used.

## Bounds

No attention, state expansion, learned router, future oracle, memory beyond 16 entries, hyperparameter search, result-informed retry, live activation, or production authority.
