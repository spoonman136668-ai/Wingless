# UP-LM9D-permutation-invariance

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM9D-permutation-invariance",
  "candidate_id": "cand-perm-reverse",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-perm-reverse-test.ps1",
    "docs/experiments/UP-LM9D-permutation-invariance.md"
  ],
  "controls": [
    "fixed random_fixed seed [3,0,5,1,4,2]",
    "fixed profiles [deferred_only, layout_only, hybrid_max]",
    "fixed rotations [8,21]",
    "derived rows 144 = 2 rotations x 3 profiles x 4 permutations x 6 frozen arms",
    "same harness version as UP-LM9C",
    "same baseline commit c1d35fa339c245b382450a7cfb222c0e92f0154f"
  ],
  "fixed_parameters": {
    "permutations": "[identity, rotate1, random_fixed, reverse]",
    "profiles": "[deferred_only, layout_only, hybrid_max]",
    "random_fixed": "[3,0,5,1,4,2]",
    "rotations": "[8, 21]",
    "rows": "144"
  },
  "metrics": [
    {
      "name": "mismatch_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_post_target_error",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_write_error",
      "comparator": "==",
      "threshold": 0
    }
  ],
  "seeds": [
    42,
    123,
    456,
    789,
    101112,
    131415,
    161718,
    192021
  ],
  "compute_seconds": 3600,
  "stop_conditions": [
    "any metric exceeds threshold",
    "compute_seconds exceeded",
    "harness error"
  ],
  "positive_meaning": "The exact deficit law transfers to the reverse permutation with 0 mismatch rows, 0 max post-target error, and 0 max write error across all 144 derived Cartesian rows, confirming invariance for this substrate dimension.",
  "negative_meaning": "The exact deficit law does not transfer to the reverse permutation; at least one mismatch row or non-zero error appears, falsifying permutation invariance for this dimension.",
  "mixed_meaning": "Some metrics meet the exact deficit law thresholds while others do not, indicating partial invariance that requires dimension-specific analysis."
}
```
