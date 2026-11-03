package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/unitary"
	"io"
	"os"
)

type request struct {
	Op         string                      `json:"op"`
	Examples   []unitary.NativeTextExample `json:"examples"`
	Input      string                      `json:"input"`
	Checkpoint json.RawMessage             `json:"checkpoint"`
}

func handle(raw []byte, state **unitary.NativeTextSelector) (any, error) {
	if unitary.ValidateNativeSessionJSON(raw) != nil {
		return nil, fmt.Errorf("TEXT_JSON_INVALID")
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return nil, fmt.Errorf("TEXT_REQUEST_INVALID")
	}
	var r request
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&r) != nil || d.Decode(new(any)) != io.EOF {
		return nil, fmt.Errorf("TEXT_REQUEST_INVALID")
	}
	allowed := map[string][]string{"train": {"op", "examples"}, "predict": {"op", "input"}, "checkpoint": {"op"}, "restore": {"op", "checkpoint"}}
	keys, ok := allowed[r.Op]
	if !ok || len(fields) != len(keys) {
		return nil, fmt.Errorf("TEXT_REQUEST_FIELDS_INVALID")
	}
	for _, k := range keys {
		if _, ok := fields[k]; !ok {
			return nil, fmt.Errorf("TEXT_REQUEST_FIELDS_INVALID")
		}
	}
	switch r.Op {
	case "train":
		if *state != nil {
			return nil, fmt.Errorf("SILENT_STATE_REPLACEMENT_DENIED")
		}
		m, err := unitary.TrainNativeTextSelector(r.Examples)
		if err != nil {
			return nil, err
		}
		*state = m
		return map[string]any{"ready": true}, nil
	case "restore":
		if *state != nil {
			return nil, fmt.Errorf("SILENT_STATE_REPLACEMENT_DENIED")
		}
		m, err := unitary.RestoreNativeTextSelector(r.Checkpoint)
		if err != nil {
			return nil, err
		}
		*state = m
		return map[string]any{"ready": true}, nil
	}
	if *state == nil {
		return nil, fmt.Errorf("TRAIN_OR_RESTORE_REQUIRED")
	}
	if r.Op == "predict" {
		return (*state).Predict(r.Input)
	}
	rawCheckpoint, err := (*state).Checkpoint()
	if err != nil {
		return nil, err
	}
	return map[string]any{"checkpoint": json.RawMessage(rawCheckpoint)}, nil
}
func main() {
	var state *unitary.NativeTextSelector
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 65536)
	out := json.NewEncoder(os.Stdout)
	for scanner.Scan() {
		value, err := handle(scanner.Bytes(), &state)
		if err != nil {
			out.Encode(map[string]any{"ok": false, "error": err.Error(), "language_proof": false})
			continue
		}
		out.Encode(map[string]any{"ok": true, "result": value, "language_proof": false})
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "TEXT_REQUEST_SIZE_OR_STREAM_ERROR")
		os.Exit(1)
	}
}
