# Wingless UP-LM0H — whole lexical-family holdout

Status: preregistered scientific language/memory generalization experiment.

Scientific parent: sealed UP-LM0G `e2b6dd94977bfde44c80f91480a308aae5d943ee`.

## Question

UP-LM0G closed the LM0F domain-transfer failure when the admission classifier trained inside the natural lexical domain. Does that learned admission rule generalize to entirely unseen verb families rather than merely held-out name/value recombinations?

## Frozen semantic families

STORE:
- training surfaces: `stores`, `saves`;
- held-out surface: `retains`.

OBSERVE:
- training surfaces: `observes`, `sees`;
- held-out surface: `notes`.

REPORT remains `reports`.

All six names and six values remain unchanged.

## Frozen classifier

- deterministic 64-D recurrent event encoder;
- binary logistic head;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- threshold exactly 0.5;
- training uses only the four training verb surfaces;
- held-out classifier evaluation uses only `retains` and `notes`;
- no explicit event type at inference.

## Frozen language model

- same 64-value recurrent state;
- same transport/update rule;
- same linear softmax output;
- 20 epochs;
- learning rate 0.08;
- no attention;
- exact recall cap 16;
- 512 recurrent-state bytes.

Training paragraphs use only training verb surfaces, deterministically alternating:
- STORE between `stores` and `saves`;
- OBSERVE between `observes` and `sees`.

Held-out paragraphs replace STORE with `retains` and OBSERVE with `notes`.
REPORT clauses are unchanged.

## Frozen arms

1. `explicit_store_admission`
2. `learned_family_holdout_admission`

Both use identical recurrent and exact-memory budgets.

## Evaluation

- classifier train-family and held-out-family accuracy/precision/recall;
- language train-family, held-out-family, and held-out-family stream4;
- top-1 next-byte accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- exact accuracy by update count;
- admission precision/recall;
- maximum recall entries.

## Interpretation

The explicit arm separates language-surface generalization from admission-classifier generalization. If both arms degrade similarly, the language model is the bottleneck. If only the learned arm degrades, semantic admission is the bottleneck.

## Bounds

No vocabulary-dependent capacity increase, no state expansion, no memory above 16 entries, no attention, no threshold tuning, no future oracle, no result-informed retry, no live activation, no production authority.
