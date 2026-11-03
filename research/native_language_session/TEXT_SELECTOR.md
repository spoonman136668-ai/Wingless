# Native whole-response text component

This component accepts complete UTF-8 input/response pairs, trains the project's
original native learner, and selects one complete response from the supplied
bank. It is **closed-response selection**, not novel-string generation, a
general language model, or the language proof point. Every response reports
`language_proof:false`; predictions also report `novel_string_generation:false`.

Limits: 2–16 examples, 2–8 distinct supplied responses, 1–256 UTF-8 bytes per
input, <=1024 input bytes total, <=128 UTF-8 bytes per response. No corpus,
external model, old scientific cohort, parser grammar, or hardcoded response is
embedded. All training examples are validated before native updates. Training
is from scratch. Prediction is read-only. Ties use deterministic bank order.

The retained state includes the response bank and **every native model's full
learned state**, source/schema/encoding pins, training-data digest and integrity
digests. Digests are not authentication, proof of legitimate data, or scientific
custody. Source drift and corrupted/foreign state are rejected.

JSONL operations:
- `{"op":"train","examples":[{"input":"aaaa","response":"amber"},{"input":"zzzz","response":"violet"}]}`
- `{"op":"predict","input":"aaaa"}`
- `{"op":"checkpoint"}`
- `{"op":"restore","checkpoint":<returned checkpoint object>}` in a fresh process.

The example strings are invented engineering fixtures. No linguistic
performance is inferred from them. New tests check native-state updates, score
influence versus untrained state, complete checkpoint/fresh-process identity,
prediction nonmutation and invalid inputs. Existing mandatory regression and
source-index acceptance remain required. Historical negatives remain negative.

Before any language proof: separately freeze legitimate development data and
genuine held-out unfamiliar transfer; train equally informed controls; compare
untrained, no-state and shuffled-label arms; establish prospective paired
retention measures; measure all state, training and inference costs. No such
scientific proof, architectural novelty, efficiency advantage or autonomous
FORGE development is claimed here.

Wingless uses the unchanged NCI11 learned associative implementation (blob
`f248fd6dec93d04571f7cb4c3ccd721017812bfd`). Raw UTF-8 bytes fold into 13
integer percepts and 12 edges. Original native learning trains encoding/readout
projections and calibration for one classifier per response; prediction builds
and reads the original associative matrix. Folding is lossy and can collide:
it is not a semantic parser, and wording/meaning transfer is unverified.

Build `go build -o wingless-native-text-selector ./cmd/wingless-native-text-selector`.
Run the executable for JSONL. Check `go test ./unitary -run '^TestNativeText' -v`
and `python research/native_language_session/check_text_cli.py ./wingless-native-text-selector`.
