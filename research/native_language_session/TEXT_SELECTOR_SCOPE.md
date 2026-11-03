# Whole-response native text selector scope

Status: ENGINEERING COMPONENT ONLY; preregistered before implementation. No scientific first use or language-model acceptance.

Objective: accept complete UTF-8 input/response examples, learn through the project's original native learner, and select a whole response from the learned bank. Output selection is closed-set, not novel-string generation or LLM-class capability.

Limits frozen for this component: 2–16 examples; 2–8 distinct supplied responses; each input 1–256 UTF-8 bytes; total input <=1024 bytes; response <=128 UTF-8 bytes; checkpoint <=131072 bytes. Validate all training input before learning. Selection is read-only. Training starts from scratch, not continual adaptation.

Wingless: deterministically fold raw UTF-8 bytes into 13 bounded integer percepts and 12 edges, with no language grammar. Train one original NCI11 learned associative classifier per supplied response via the unchanged native learning routine. Both learned encoding/readout projections and the native associative matrix determine scores. This byte folding is lossy and supplies no semantic parser; familiar wording transfer is unverified.

Yggdrasil: present each raw UTF-8 byte as high/low four-boolean sensory events in the original A/B native local-cell contexts. Each candidate response has an original Colony. Completed example-level response labels train the original local Hebbian/consolidation rule; original bounded-neighbor predictions select the response. This bounded learner bank is an input/output component, not a substitute cognitive foundation or proof of the full genome/phenotype, growth, repair or compact inheritance thesis. Existing developmental lineage remains intact.

No external model, transformer, next-token objective, conventional model replacement, hardcoded response/grammar, old experiment cohort invocation, protected source/data admission, controller, queue, live runtime activation, accepted-ref mutation, FORGE takeover or promotion. Input examples are caller-supplied; schema/hash checks are not scientific provenance.

Tests use invented engineering fixtures only: validate bounds, actual native-state updates and score influence, complete checkpoint and fresh-process parity, read-only prediction, corrupted/foreign state rejection. Report any native learning failure honestly; engineering fixture fit is not a language proof. No performance/transfer threshold is invented or historical negative relabeled.

Before a language proof, separately freeze legitimate development data, equally informed trained controls, no-state/untrained/shuffled-label comparisons, genuine held-out unfamiliar transfer and prospective paired retention measures. Original R7 negative, original Wingless scientific outcomes and frozen thresholds remain unchanged.

Project: Wingless
Exact parent: a400bf5faaf027b886a19f4437ef401a489a7d2f
Original native learner Git blob: f248fd6dec93d04571f7cb4c3ccd721017812bfd
