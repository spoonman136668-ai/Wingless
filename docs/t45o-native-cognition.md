# T45-O: Opt-in native cognitive candidate (not scientific acceptance)

The maintained `cmd/wingless` executable now exposes `wingless cognition` as an offline JSONL event interface. It runs a bounded `cognition.State` logistic prediction-error updater on binary sensor/action and observed consequences, and marginalizes two-step predictions without seeing future outcomes.

Example on stdin:

```jsonl
{"op":"observe","action":1,"prior":0,"outcome":1,"rate":0.05}
{"op":"predict","action":1,"prior":0}
{"op":"forecast_two","prior":0,"action":1,"next_action":0}
{"op":"snapshot"}
```

The caller owns persistent snapshot custody and may send a separate `restore` operation. This command has no file writes, credentials, model APIs, worker dispatch, tool access, accepted-ref mutation, or acceptance rights. It does not modify the default worker/backend path.

**Evidence limitation:** This is a conventional four-weight learned update intended for integration and falsification. It has not been shown to outperform an equally informed symbolic inference or count-table competitor, and is not a novel cognition claim. Any use inside the real worker reasoning path requires separate affected-interface qualification and FORGE-native scientific controls. CKB-plane is not a research execution or acceptance dependency.
