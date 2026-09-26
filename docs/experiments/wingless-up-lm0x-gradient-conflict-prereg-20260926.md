# Wingless UP-LM0X — matched gradient-conflict diagnostic

Status: preregistered scientific language-model diagnostic experiment.

Scientific parent: sealed UP-LM0W 205791da63c3ce4a214d44006ab77564e93c305e.

## Question

LM0W removed between-example update order but produced an intermediate retention/adaptation profile rather than a joint improvement. Are matched base and paraphrase examples genuinely pushing the byte readout in conflicting directions?

## Frozen model and corpus

- train the exact base byte model for 20 epochs as in LM0W;
- do not perform paraphrase adaptation;
- use the exact matched base/paraphrase training pairs from the LM0W corpus;
- deterministic 64-D recurrent states;
- readout weights and biases only are analyzed;
- no router training or evaluation is required for this diagnostic.

## Frozen gradient measurement

For every matched training pair at the same frozen base-model parameters:
1. compute the complete readout gradient of the base example without applying it;
2. compute the complete readout gradient of the matched paraphrase example without applying it;
3. include all readout weights and biases;
4. compute pair cosine similarity, dot product, and gradient norms.

Aggregate:
- pair count;
- mean, minimum, and maximum pair cosine;
- fraction of pairs with cosine < 0;
- fraction with cosine <= 0;
- concatenated global cosine over all matched pair gradients.

Also report a per-output-byte row conflict table:
- byte value;
- row cosine across all matched pairs and row parameters;
- accumulated dot product;
- base gradient norm;
- paraphrase gradient norm.

Row metrics are reported for every byte in the frozen alphabet in deterministic alphabet order.

## Interpretation

Substantial negative pair or row cosine would directly support competing-gradient interference. Predominantly aligned gradients would reject that diagnosis and redirect the next experiment toward optimization trajectory or capacity effects.

## Bounds

No parameter updates after the base model is trained, no adaptation, no anchoring, no replay, no gradient surgery, no threshold tuning, no result-informed retry, no live activation, no production authority.
