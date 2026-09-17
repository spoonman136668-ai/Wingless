package resources

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestModelTelemetryPreservesUnknownVersusMeasuredZero(t *testing.T) {
	zero := uint64(0)
	m := NewModelTelemetry()
	m.DataMovement = &ModelDataMovementTelemetry{
		HostToDeviceBytes: &zero,
		Provenance: &TelemetryProvenance{
			Origin: TelemetryOriginBackend, ProviderID: "fixture-backend",
		},
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, `"host_to_device_bytes":0`) {
		t.Fatalf("measured zero was not preserved: %s", text)
	}
	if !strings.Contains(text, `"device_to_host_bytes":null`) || !strings.Contains(text, `"capacity":null`) {
		t.Fatalf("unknown telemetry was not encoded as null: %s", text)
	}
}

func TestParameterEstimateAlwaysMarksEstimate(t *testing.T) {
	e := ParameterEstimate{Count: 123, Basis: "backend-reported activation estimate"}
	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"estimated":true`) {
		t.Fatalf("parameter estimate lost estimate marker: %s", b)
	}
	if _, err := json.Marshal(ParameterEstimate{Count: 123}); err == nil {
		t.Fatal("parameter estimate without basis was accepted")
	}
}
