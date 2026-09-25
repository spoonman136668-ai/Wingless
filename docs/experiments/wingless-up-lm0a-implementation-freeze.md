# Wingless UP-LM0A — implementation freeze

Status: frozen before first execution.

This file only fixes implementation details left open by the preregistration. It does not change the corpus, structural split, recurrent-state size, epochs, learning rate, or evaluation families.

## Byte alphabet

The output alphabet is the sorted unique set of raw ASCII byte values appearing anywhere in the complete frozen 144-sentence corpus.

Held-out examples contain no byte absent from this alphabet.

## Recurrent transport

For output coordinate `i` in 0..63:

- source coordinate = `(13*i + 7) mod 64`;
- sign = +1 when `i` is even and -1 when `i` is odd.

Because 13 is coprime with 64, the source mapping is a permutation. Signed permutation transport is norm-preserving before input injection.

## Byte embedding

For byte value `b` and state coordinate `i`:

- derive a deterministic 64-bit word from `b`, `i`, and a frozen constant using the existing Wingless `sq0Mix64` mixer;
- map the low bit to ±1/sqrt(64).

No embedding parameter is trained.

## State update

Exactly:

`h_new[i] = tanh(0.90 * transported_h[i] + 0.35 * byte_embedding[b][i])`

No layer normalization, learned recurrence, residual branch, dropout, attention, or exact-recall side channel.

## Output training

- linear softmax head only;
- weights and biases initialized to zero;
- standard cross-entropy gradient;
- online SGD;
- learning rate 0.08;
- 20 epochs;
- deterministic lexical sentence order every epoch;
- no momentum, weight decay, batching, shuffling, scheduler, validation stopping, or retry.

## Scoring

Sentence evaluation:
- recurrent state resets at each sentence start;
- score byte pairs `sentence[t] -> sentence[t+1]`.

Four-sentence stream evaluation:
- group held-out sentences in deterministic lexical groups of four;
- recurrent state resets only at the start of each four-sentence group;
- state is preserved across sentence boundaries;
- the artificial cross-boundary pair `newline -> first byte of next sentence` is not included in accuracy or cross-entropy, because that transition is absent from sentence-level training.

## Bigram baseline

- estimate conditional next-byte probabilities by unsmoothed maximum likelihood from the same training byte pairs;
- top-1 ties resolve to the lowest byte value;
- during cross-entropy reporting only, a zero-probability held-out transition is floored to 1e-12 to keep the diagnostic finite;
- the floor does not alter top-1 predictions or training counts.

## Interpretation

This is a first language-objective qualification, not a real-language benchmark. A negative result is valid evidence that the fixed recurrent carrier/readout is insufficient and must be changed before moving to real text.
