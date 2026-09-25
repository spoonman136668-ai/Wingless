# Wingless UP-LM0D — multi-binding language recall

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0C `6323443c61da2841f9ad6889d13abacd23590e61`.

## Question

UP-LM0C showed that a local syntax router plus the fixed 16-entry exact store closes a single long dependency. Can the same architecture preserve and retrieve several simultaneous name/value bindings inside one paragraph, with multiple delayed queries and no attention?

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

Each paragraph contains four distinct cyclic names, four distinct cyclic values, one distractor, then four report clauses.

For name offset `n`, value offset `v`, and distractor index `d`:

- bindings use names `n,n+1,n+2,n+3 mod 6`;
- values use `v,v+1,v+2,v+3 mod 6`;
- store clauses occur in binding order;
- report order is a deterministic rotation by `d`.

Paragraph form:

`<n0> stores <v0>. <n1> stores <v1>. <n2> stores <v2>. <n3> stores <v3>. <distractor> <q0> reports <vq0>. <q1> reports <vq1>. <q2> reports <vq2>. <q3> reports <vq3>.\n`

There are 144 paragraphs.

Structural split:
- train iff `(n + 2*v + d) mod 3 != 2`;
- held out otherwise.

Every name, value, and distractor occurs in training.

## Frozen recurrent model

Identical to UP-LM0B/UP-LM0C:
- 64 float64 recurrent state values;
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
   - same local `stores` / `reports` parser as UP-LM0C;
   - exact recall cap stays exactly 16 entries;
   - memory may override only the first byte of a reported value;
   - no future-byte access.

## Evaluation

Train, held-out, and held-out stream4:
- overall next-byte top-1 accuracy;
- perplexity;
- dependent first-byte accuracy over all four report queries;
- whole-paragraph four-query exact accuracy;
- maximum exact-recall entries used.

## Bounds

No attention, state expansion, learned router, future oracle, memory beyond 16 entries, hyperparameter search, result-informed retry, live activation, or production authority.
