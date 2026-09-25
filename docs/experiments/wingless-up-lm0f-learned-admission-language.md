# Wingless UP-LM0F — learned admission inside mutable language memory

Status: preregistered scientific language/memory integration experiment.

Scientific parents:
- sealed UP-LM0E `5a4b8d05abe54e1b1c0d8107c4d4f9434cf65801`;
- sealed UP-95B `f15ff278753d3af1786b38f9958e41ded3c00fc5`.

## Question

UP-LM0E showed perfect mutable latest-value language recall when store admission was hand-coded from syntax. UP-95B showed a learned local event classifier can perfectly distinguish `stores` from `observes` on held-out keys. Can that learned local classifier replace explicit store admission inside the language-memory loop without losing mutable recall?

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

Each paragraph contains:
1. four initial `<name> stores <value>.` clauses;
2. four `<name> observes <value>.` clauses using distractor values that must not alter memory;
3. a frozen number of `stores` update clauses: 0, 1, 2, or 4;
4. four `<name> reports <value>.` clauses whose targets are latest stored values.

The four names and initial values are cyclic exactly as UP-LM0E. Update count is determined by paragraph index modulo 4. Observation values are deterministically offset by three positions modulo 6. Report order is a deterministic rotation.

There are 144 paragraphs.

Structural split:
- train iff `(name_offset + 2*value_offset + pattern_index) mod 3 != 2`;
- held out otherwise.

Every name, value, update count, and event verb appears in training.

## Frozen recurrent language model

Identical to UP-LM0E:
- 64 float64 recurrent state values;
- same deterministic signed-permutation transport;
- same deterministic bipolar byte embedding;
- same recurrent update;
- same linear softmax head;
- 20 epochs;
- learning rate 0.08;
- no attention;
- no pretrained weights.

## Frozen event classifier

Identical to UP-95B:
- same 64-D local event encoder;
- binary linear logistic head;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- threshold 0.5;
- training keys k00..k23, held-out keys k24..k31;
- no classifier threshold tuning.

At language inference, the classifier receives only the completed local clause bytes. It does not receive an explicit event-type bit.

## Frozen arms

1. `explicit_store_admission`
   - upper-bound control;
   - admit only clauses with explicit `stores` syntax.

2. `learned_local_admission`
   - completed `stores` and `observes` clauses are classified by the frozen local event classifier;
   - admit to exact recall only when classified as store.

Both arms:
- exact recall cap exactly 16;
- overwrite existing name binding in place;
- reports query exact recall;
- memory may override only the first reported value byte.

## Evaluation

Train, held-out, and held-out stream4:
- overall next-byte top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-paragraph query-set exact accuracy;
- exact accuracy by update count 0/1/2/4;
- classifier admission precision/recall inside language;
- maximum recall entries used.

## Bounds

No explicit event-type bit in learned arm, no attention, no future oracle, no recurrent-state expansion, no memory beyond 16 entries, no threshold search, no result-informed retry, no live activation, no production authority.
