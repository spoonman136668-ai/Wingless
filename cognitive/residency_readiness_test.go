package cognitive

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type sessionCaptureBackend struct {
	scriptedBackend
	sessionIDs []string
}

func (b *sessionCaptureBackend) Invoke(ctx context.Context, req inference.Request) (inference.Result, error) {
	b.sessionIDs = append(b.sessionIDs, req.SessionID)
	return b.scriptedBackend.Invoke(ctx, req)
}

func TestInferenceSessionResidencyTelemetry(t *testing.T) {
	logical := uint64(100)
	zero := uint64(0)
	model := resources.NewModelTelemetry()
	model.Capacity = &resources.ModelCapacityTelemetry{
		LogicalModelBytes: &logical,
		Provenance: &resources.TelemetryProvenance{
			Origin: resources.TelemetryOriginBackend, ProviderID: "scripted",
		},
	}
	model.DataMovement = &resources.ModelDataMovementTelemetry{
		HostToDeviceBytes: &zero,
		Provenance: &resources.TelemetryProvenance{
			Origin: resources.TelemetryOriginBackend, ProviderID: "scripted",
		},
	}
	result := completed("done")
	result.Telemetry.Model = &model
	backend := &sessionCaptureBackend{scriptedBackend: scriptedBackend{results: []inference.Result{result}}}
	rt := Runtime{Backend: backend, Policy: Policy{Version: "session-v1", MaxPasses: 1}}

	res, err := rt.Run(context.Background(), baseRun("session-telemetry"))
	if err != nil {
		t.Fatal(err)
	}
	s := res.Evidence.Session
	if s.SessionID != "session-telemetry:inference-session" || s.ModelCalls != 1 || s.CognitivePasses != 1 {
		t.Fatalf("bad session accounting: %+v", s)
	}
	if len(backend.sessionIDs) != 1 || backend.sessionIDs[0] != s.SessionID {
		t.Fatalf("backend did not receive bounded session identity: %+v", backend.sessionIDs)
	}
	if s.BackendID == nil || *s.BackendID != "scripted" || s.ModelID == nil || *s.ModelID != "fixture" {
		t.Fatalf("bad session identity: %+v", s)
	}
	if s.ResourceTelemetry == nil || s.ResourceTelemetry.Capacity == nil || s.ResourceTelemetry.Capacity.LogicalModelBytes == nil || *s.ResourceTelemetry.Capacity.LogicalModelBytes != logical {
		t.Fatalf("model telemetry missing: %+v", s.ResourceTelemetry)
	}
	b, err := json.Marshal(s.ResourceTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	text := string(b)
	if !strings.Contains(text, `"resident_model_bytes":null`) || !strings.Contains(text, `"active_model_bytes":null`) {
		t.Fatalf("unknown residency was fabricated or omitted: %s", text)
	}
	if !strings.Contains(text, `"host_to_device_bytes":0`) {
		t.Fatalf("measured zero lost: %s", text)
	}
}

func TestMultiPassSessionDoesNotInventTelemetryAggregate(t *testing.T) {
	logical := uint64(100)
	model := resources.NewModelTelemetry()
	model.Capacity = &resources.ModelCapacityTelemetry{LogicalModelBytes: &logical}
	first := completed("draft")
	first.Telemetry.Model = &model
	second := completed("final")
	second.Telemetry.Model = &model
	backend := &sessionCaptureBackend{scriptedBackend: scriptedBackend{results: []inference.Result{first, second}}}
	rt := Runtime{
		Backend: backend,
		Policy:  Policy{Version: "session-v1", MaxPasses: 2, AdditionalPassPrompt: "finish"},
		Evaluator: EvaluatorFunc(func(pass int, _ inference.Result) Evaluation {
			return Evaluation{Sufficient: pass == 2, Reason: "fixture"}
		}),
	}
	res, err := rt.Run(context.Background(), baseRun("multi-telemetry"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Evidence.Session.ModelCalls != 2 || res.Evidence.Session.CognitivePasses != 2 {
		t.Fatalf("bad multi-pass accounting: %+v", res.Evidence.Session)
	}
	if len(backend.sessionIDs) != 2 || backend.sessionIDs[0] != res.Evidence.Session.SessionID || backend.sessionIDs[1] != res.Evidence.Session.SessionID {
		t.Fatalf("multi-pass session identity split: %+v", backend.sessionIDs)
	}
	if res.Evidence.Session.ResourceTelemetry != nil {
		t.Fatal("multi-pass telemetry was unsafely aggregated")
	}
	if len(res.Evidence.Passes) != 2 || res.Evidence.Passes[0].Telemetry.Model == nil || res.Evidence.Passes[1].Telemetry.Model == nil {
		t.Fatal("per-pass telemetry was lost")
	}
}

func TestRouteExtensibilityDoesNotGrantExecution(t *testing.T) {
	if RouteAuthorizedDeeperBackend.ExecutableInCR1A() || RouteToolProposal.ExecutableInCR1A() {
		t.Fatal("representable future/authorized routes became executable")
	}
	if !RouteMemory.ExecutableInCR1A() || !RouteSkill.ExecutableInCR1A() || !RouteInferenceSingle.ExecutableInCR1A() || !RouteInferenceMulti.ExecutableInCR1A() {
		t.Fatal("existing CR-1A route execution set regressed")
	}
	decision, _, _, err := (Router{Policy: Policy{Version: "routing-v1", MaxPasses: 1}}).Decide(context.Background(), baseRun("route-version"))
	if err != nil {
		t.Fatal(err)
	}
	if decision.Version != RoutingDecisionVersion {
		t.Fatalf("routing version=%d want %d", decision.Version, RoutingDecisionVersion)
	}
}
