# NR-M1T — independent review packet for changed inherited CLI `main()`

**State:** REVIEW REQUIRED / NOT APPROVED. This is a narrow source-review packet, not proof of independent approval, native model promotion, bypassed ICE validation, scientific first use, or authorization to merge.

## Exact prior review and current source

- Historical source-reviewed atlas was rebased *only* by already approved Wingless [#344](https://github.com/spoonman136668-ai/Wingless/pull/344). Its exact reviewed source commit `fd00165e92b4fc21d93d33292c32038e6e96bb22` is a **real ancestor** of the current native research branch (GitHub comparison verified).
- Frozen reviewed declaration: `cmd/wingless/main.go`, `main()`, normalized declaration SHA-256 `1b79588a2867fa9951569289c4bd24a43eb3a7a8957c1b7293aa9ddcb5916743`.
- Current inherited declaration: `cmd/wingless/main.go`, `main()`, normalized declaration SHA-256 `803dd3a99f4234b9667c9ca44770d0ada3b7ed2e59eb8b678d0f1ea43c347f02`; current file Git blob `5f1ddb8e1c5fbe28005451c7c1d76f217cdf3519`. **These are unequal**. No claim of byte-identical source reuse is permissible.
- The change was originally introduced by independent prior engineering candidate commit `1131ae7d000bedf8c9c24c8f3f7dea1a60a01c58` in draft [T45-O PR #340](https://github.com/spoonman136668-ai/Wingless/pull/340), not by new MG40 task-memory files.
- Maintainer audit [run #41029835192](https://github.com/spoonman136668-ai/Wingless/actions/runs/41029835192) reported **957** canonical current Go declarations: 352 historically reviewed identities retained, 605 *new lexical-provisional/unreviewed* declarations, **exactly one changed historically reviewed declaration** (`main`). The unchanged G26 source-aware validator correctly returned `G26_REVIEW_REQUIRED`, not a false PASS.
- Existing G28 provisional pathway is already merged ([#345](https://github.com/spoonman136668-ai/Wingless/pull/345)): it may register new source-bound lexical symbols without claiming independent review. It **cannot** promote or replace a modified reviewed declaration with unsupported review, and must never be used to bypass this one changed `main`.

## Auditable behavior delta (original commit #340)

The original entrypoint was restricted to the historical `gpu`, `local-benchmark`, `host`, `demo`, and `benchmark` paths. The new source adds exactly this dispatch branch at the top of `main()`:

```go
if len(os.Args) == 2 && os.Args[1] == "cognition" {
    if e := cognitionCLI(); e != nil {
        fmt.Fprintln(os.Stderr, e)
        os.Exit(1)
    }
    return
}
```

The usage string also adds `cognition`. No historical dispatch cases were removed in this specific originating diff. The new handler in `cmd/wingless/cognition.go` is an offline JSON-lines stateful cognition adapter:

- Passes only `os.Stdin` and `os.Stdout` to `runCognition`. No source environment or hardware command is invoked in `main`.
- `runCognition` accepts explicitly implemented `observe`, `predict`, `forecast_two`, `snapshot`, and `restore` operations; rejects default/unknown operation including `execute`.
- The input has a 1 MiB reader cap, scanner max token 64 KiB and up to 4,096 processed events per call. Decoder rejects unknown JSON fields and extra JSON text. One concern for review: rejected or malformed input exits nonzero, but a scanner/read error should also remain nonzero.
- Existing source test `cmd/wingless/cognition_test.go::TestOfflineCognition` checks an observe/forecast/snapshot path and rejects unauthorized `execute`.
- This candidate's NR-M1T CI runs the **actual** inherited `go test ./cmd/wingless -run '^TestOfflineCognition$' -count=1 -v` and `go vet ./cmd/wingless` on Linux and Windows. These prove an engineering guard only; the review decision remains independent.

## Exact acceptance checklist for independent reviewer

1. Independently compare the reviewed historical `main()` with the candidate's Go source bytes; verify the verified source SHA / line spans and **the originating commit**. Do not assert unchanged declaration identity.
2. Confirm only `len(os.Args)==2 && os.Args[1]=="cognition"` enters the offline adapter, failed cognition commands exit nonzero, and all historical dispatch cases and their usage semantics are preserved.
3. Confirm the handler's bounded input/allowed operations and failure semantics; no original-world stepping, side-effecting OS commands, latent source oracle, network, unapproved paid model, or production executor is reachable from the newly introduced path.
4. Re-run focused CLI+MG40 tests and Go vet on Linux/Windows, then independently review any related source-controlled failure. Check all existing Wingless ICE syntax, reviewed 352-label and relationship, browser/tamper, source-index source custody and native acceptance. **Do not waive a failure.**
5. If and only if reviewer approval is genuine, submit the approved evidence through the repository's established independent source-review process, with immutable original/current declaration SHAs and reviewed content. Self-asserted automated receipts, this packet, the connected GitHub accounts, an AI summary, and a green CI job **do not** authenticate an independent reviewer.
6. Regenerate canonical syntax index and the G28 atlas from the **exact approved** commit; only the 605 genuinely new declarations may be labeled `lexical-provisional` until independently reviewed. Preserve all 352 old reviews, correct any changed main review only through accepted migration, and pass full final-head ICE before native promotion.

## Current stop decision

**Blocked on the altered inherited `cmd/wingless/main.go::main` reviewed label.** The lexical source proves the previous source and origin, but cannot confer independent review. Keep native [PR #438](https://github.com/spoonman136668-ai/Wingless/pull/438) draft, research-only, not merged or activated. The independently source-only and full G29-qualified [Yggdrasil #700](https://github.com/spoonman136668-ai/Yggdrasil/pull/700) remains unaffected. This packet does not require delivering local receipts to the owner; source history stays in GitHub.
