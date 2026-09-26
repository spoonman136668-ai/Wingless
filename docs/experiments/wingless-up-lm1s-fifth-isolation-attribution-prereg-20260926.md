# Wingless UP-LM1S — fifth-family isolation attribution

Status: preregistered scientific five-family byte-readout attribution experiment.

Scientific parent: sealed UP-LM1R d11cd78832db0685d7e5485103863266b4792da9.

## Question

UP-LM1R found that the fifth family's output-readout gradient is not systematically more anti-aligned with the prior four families, despite a much larger gradient norm and much worse frozen-start perplexity. With the recurrent state frozen, does the fifth family remain weak when given the exact same fifth-family update dose in isolation, or is its weakness created by interleaving prior-family updates?

## Frozen starting point

Use the exact five-family pre-adaptation starting point from UP-LM1Q/UP-LM1R:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family rotating-palindromic training history;
- output alphabet extended for fifth-family bytes;
- matched base / paraphrase / third / fourth / fifth corpora.

## Arms

All adaptation arms use learning rate 0.08 and 4 epochs.

1. no_adaptation
   - no readout updates.

2. fifth_only
   - one fifth-family update per matched example per epoch;
   - exactly the same fifth-family update count as the cyclic arm;
   - no prior-family updates.

3. five_family_cyclic
   - exact LM1P cyclic control;
   - one update per family per matched example in canonical family order.

No recurrent updates.

## Evaluation

For every arm:
- held-out top-1 accuracy and perplexity for all five families;
- minimum, mean, and spread across families;
- fifth-family held-out accuracy/perplexity;
- mean prior-family held-out accuracy.

Derived comparisons:
- fifth_only minus five_family_cyclic fifth-family accuracy;
- fifth_only minus five_family_cyclic prior-family mean accuracy;
- fifth_only fifth-family gain over no_adaptation.

## Interpretation

If fifth-only adaptation substantially outperforms matched-dose cyclic adaptation on fifth-family held-out accuracy, competition from prior-family updates is causally suppressing fifth-family learning even though broad aggregate gradient anti-alignment was not observed. If fifth-only remains similarly weak, the limiting factor is more consistent with intrinsic representation/readout difficulty or introduction history.

No numeric success threshold is introduced; exact measured differences are the result.

## Bounds

No recurrent training, no router, no exact recall, no adaptive schedule, no learning-rate search, no extra fifth-family updates beyond the matched cyclic fifth dose, no sixth family, no result-informed retry, no live activation, no production authority.
