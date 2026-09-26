# Wingless UP-LM1G — symmetric three-family gradient scale

Status: preregistered scientific optimization diagnostic.

Scientific parent: sealed UP-LM1F a24fcfb761d68e6ce191bdb21d994efde68fddb9.

## Question

LM1F showed that averaging matched base/paraphrase/third-family gradients underperforms cyclic sequential updates, especially on the third family. Is that negative caused by a smaller effective optimization step, or by genuine interference/capacity limits?

## Frozen starting point

Use the exact LM1F extended warm20 starting model:
- recurrent transition frozen;
- grounded router frozen;
- exact recall cap 16;
- four adaptation epochs;
- matched base/paraphrase/third-family training corpora;
- no attention.

For every matched triplet, compute all three readout gradients from the same pre-update parameters.

## Frozen arms

1. cyclic_by_index
   - exact LM1F sequential control.

2. symmetric_scale_1over3
   - apply (g_base + g_para + g_third) * (0.08 / 3).
   - exact LM1F symmetric control.

3. symmetric_scale_1oversqrt3
   - apply the same gradient sum at 0.08 / sqrt(3).

4. symmetric_scale_1
   - apply the same gradient sum at 0.08.
   - this approximates the first-order total magnitude of three sequential lr=0.08 updates while removing within-triplet recency.

All scales are preregistered; none is selected post hoc.

## Evaluation

For every arm:
- base held-out block;
- paraphrase held-out block;
- third held-out block;
- all third-family structural order families at stream1/stream4;
- byte top-1/perplexity plus full routing/memory metrics.

## Interpretation

Recovery at larger symmetric scale would attribute LM1F mainly to under-stepping. Persistent inferiority across scales would support a genuine benefit from sequential/nonlinear optimization or a readout-representation limitation.

## Bounds

No recurrent/router training, no extra examples, no state/recall expansion, no adaptive scaling, no result-informed retry, no live activation, no production authority.
