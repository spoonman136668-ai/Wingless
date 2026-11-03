# Wingless native learned-state session

This is the first runtime component toward a native language prototype, not a
language model or a language proof point. Inputs are explicitly supplied integer
relations and boolean query labels. There is no text parser, generator, corpus,
external model, controller, queue or automatic promotion.

The runtime calls the unchanged NCI11 associative learner from parent
`b820ea1ff76f17e350a914177ea823f6984a8e09`; its source blob is
`f248fd6dec93d04571f7cb4c3ccd721017812bfd`. The original negative research
results and frozen experiments remain unchanged. The native matrix and its
learned encoding/readout parameters both participate in prediction.

Build `go build -o wingless-native-session ./cmd/wingless-native-session`.
The executable reads one JSON object per line on stdin and writes one reply.
Every reply explicitly reports `language_capability:false`.

Operations (exact fields):

| Operation | Additional fields | Behavior |
|---|---|---|
| `train` | `examples`, `leaky` | From-scratch native training; 1–16 examples |
| `observe` | `edges` | Replaces the current graph memory using learned parameters |
| `predict` | `from`, `to` | Read-only native boolean prediction and logit |
| `checkpoint` | none | Complete parameters, matrix, loaded flag, schema and source pin |
| `restore` | `checkpoint` | Load the complete checkpoint into an empty process |

Each example has exactly 12 `edges` (`from`/`to` integer pairs), two
`questions` (two integer pairs), and two boolean `truth` labels. IDs are
1–1,000,000. `leaky` is boolean. Retraining or restoring over an existing
session is denied. `observe` is intentionally graph replacement, not continual
training. There is no claim of learned language or arbitrary graph reasoning.

Checkpoint SHA256 detects accidental corruption; it is not authentication,
provenance, trusted feedback, or authority to consume scientific data. Source
identity is checked in CI against the original Git blob. A caller can construct
a new digest and must not treat it as evidence of genuine training.

Run `go test ./...`, `go vet ./...`, `go build ./...`, and
`python research/native_language_session/check_cli.py ./wingless-native-session`.
The new tests use invented engineering inputs only. They check learning changes
parameters, native-memory influence, complete state recovery in a fresh process,
and rejection without mutation. They do not run the NCI11 development/sealed
cohorts or establish scientific success. Existing ICE and final acceptance
requirements remain separate and mandatory before acceptance.

Next research work must add genuinely learned native language representations
and output, legitimate development data, frozen transfer evaluation, and fair
equally informed trained controls. FORGE may take over that development only
after independently qualified language evidence and its authorized integration.
