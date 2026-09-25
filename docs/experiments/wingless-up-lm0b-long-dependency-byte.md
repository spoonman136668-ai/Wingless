# Wingless UP-LM0B — long-dependency byte micro-language

Status: preregistered scientific language experiment.

Scientific parent: sealed UP-LM0A evidence `ddc471404518ba61ca0c56bd37142d6b29164e35`.

## Question

UP-LM0A showed that the frozen 64-state recurrent substrate can learn a genuine next-byte objective and generalize across held-out lexical recombinations. UP-LM0B asks whether that same substrate can carry a value across a longer linguistic dependency where the decisive next byte cannot be inferred from the local bigram context.

## Frozen corpus

Sentence form:

`<name> stores <value>. <distractor> <name> reports <value>.\n`

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

Distractor phrases:
- birds sing.
- rain falls slowly.
- lamps glow at dusk.
- quiet winds cross hills.

All 144 name/value/distractor combinations exist.

Structural split:
- train iff `(name_index + 2*value_index + distractor_index) mod 3 != 2`;
- held out otherwise.

Every name, value, and distractor appears in training.

## Frozen model

Identical to UP-LM0A:
- 64 float64 recurrent state values;
- same signed-permutation transport;
- same deterministic bipolar byte embedding;
- same update `tanh(0.90*transport(h)+0.35*embedding(byte))`;
- linear softmax output only;
- no attention;
- no exact recall;
- no pretrained weights.

Training:
- 20 epochs;
- learning rate 0.08;
- zero output weights/biases;
- online SGD;
- deterministic lexical order;
- state reset at sentence boundaries;
- no validation stopping or parameter search.

## Evaluation

Report for train and held-out:
- overall next-byte top-1 accuracy;
- cross-entropy and perplexity;
- dependent-value first-byte accuracy at the second occurrence of `<value>`;
- dependent-value first-byte cross-entropy/perplexity.

Also report a four-sentence held-out stream with no recurrent reset between sentences.

Control:
- unsmoothed bigram MLE on the identical training corpus;
- zero-probability diagnostic floor 1e-12 only for cross-entropy.

## Bounds

No new trainable recurrent parameters, tokenizer, external corpus, attention, exact-recall channel, hyperparameter search, result-informed retry, hidden-state expansion, live activation, or production authority.
