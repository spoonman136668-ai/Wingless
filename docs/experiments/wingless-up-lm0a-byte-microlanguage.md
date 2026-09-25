# Wingless UP-LM0A — byte-level micro-language objective

Status: preregistered scientific language-bridge experiment.

Scientific parent: sealed UP-SQ0 evidence `8de9e11f26cc44c75f51d26cb606f97e35816479`.

## Question

Can a compact recurrent Wingless substrate learn a genuine next-byte prediction objective and generalize to structurally held-out recombinations without attention, a KV cache, or exact-recall assistance?

This is the first language-objective bridge. It is not a claim of real-language capability.

## Frozen corpus

ASCII micro-language sentence form:

`<name> <verb> <object>.\n`

Names:
- ada
- ben
- cy
- dee
- eli
- fay

Verbs:
- sees
- likes
- finds
- moves

Objects:
- red fox
- blue cat
- green owl
- small dog
- bright sun
- calm sea

All 144 name/verb/object combinations exist.

Structural split:
- train iff `(name_index + 2*verb_index + object_index) mod 3 != 2`
- held out iff the expression equals 2.

Every lexical item therefore appears in training; held-out evaluation tests recombination, not unseen spelling.

## Frozen recurrent substrate

- recurrent state dimension: 64 float64 values;
- recurrent transport: deterministic signed permutation, therefore norm-preserving before input injection;
- byte input embedding: deterministic bipolar 64-vector derived only from the byte value;
- update: `h = tanh(0.90 * transport(h) + 0.35 * embedding(byte))`;
- no attention;
- no exact recall;
- no external model;
- no pretrained weights;
- output: learned linear softmax over the frozen observed-byte alphabet;
- training: teacher-forced next-byte prediction.

## Frozen training

- epochs: 20;
- learning rate: 0.08;
- zero-initialized output weights and biases;
- sentence order: deterministic lexical order each epoch;
- recurrent state resets at sentence boundaries during training;
- no validation-driven stopping;
- no hyperparameter search.

## Evaluation

Report:
- train next-byte top-1 accuracy;
- train cross-entropy and perplexity;
- held-out recombination next-byte top-1 accuracy;
- held-out cross-entropy and perplexity;
- 4-sentence held-out stream top-1 accuracy/perplexity without recurrent reset between sentences;
- bigram maximum-likelihood baseline on the exact same train/held-out bytes.

No pass threshold is used to tune the system after observing results.

## Bounds

No tokenizer, external corpus, transformer, attention, exact-recall channel, result-informed retry, parameter search, hidden-state expansion, live activation, or production authority.
