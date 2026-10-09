# W-SYN01 — synthetic-only causal readout falsification

**Not R53, not external validation, not READY_RESEARCH, no primary result yet.**

Mechanism foundation: Wingless R53 PR #21 exact head `2e3512beeef075d846373cc813c66b774fef8437`.
Frozen standalone preregistration commit: `b9b7944a6fe66247c87605e7db1068fcd0cd6010`.
Guarded implementation candidate: `4304510f0fab3ffebe151a4328e59e1c478a74ac`.

The previous **unexecuted** W-SYN01 draft (`e34ab20c...`) incorrectly counted 4,096 unique cases as 4,096 predictions; it has been superseded, without any primary results, by the immutable R2 freeze above: **4,096 cases × four arms = 16,384 predictions**. Do not run or merge the obsolete branch.

## Scientific question
Test whether the frozen six-dimensional R53 readout groups discriminate early-only versus late-only tied truth with fixed ridge fitting, controlled noise, exactly 436 adaptation examples per arm per regime/seed and independent synthetic heldout seeds. Wrong-block bias is deliberately tested in both directions; negative outcomes must be retained.

## Execution separation
A clean ordinary `go test ./research/contrib/ckb-synthetic/w-syn01` only tests mechanism plumbing. The primary seeded study is behind `CKB_SYN_W01_QUALIFIED_RUN=yes`; **setting that variable is not scientific authority**. Only the already-existing CKB scientific authority may authenticate exact prereg/code/tree/harness identity, issue a genuinely qualified run receipt, invoke the primary once under dedicated research resources, and classify its evidence. Do **not** set the variable in PR CI or ordinary qualification. Primary execution produces synthetic-only diagnostics with `official_ckb_acceptance=false` until separately accepted.

The current R53 external-source origin rule, R52 mixed lineage, R53/Y095 disjoint cohort and 41,453/14,365/1,454-byte budgets are unchanged. Synthetic material may **never** be used to satisfy the external R53 source cohort.

No autonomous scheduler, sidecar reactivation, accepted-ref mutation, broker path, KTRADE change, controller change, or capacity growth.
