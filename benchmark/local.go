package benchmark

import (
	"context"

	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"go/ast"
	"go/parser"
	"go/token"
	"runtime"
	"time"
)

type LocalRow struct {
	ProcessError          *string                   `json:"process_error"`
	ResourceSamplingError *string                   `json:"resource_sampling_error"`
	Repetition            int                       `json:"repetition"`
	SemanticCorrect       bool                      `json:"semantic_correct"`
	ProtocolCompliant     bool                      `json:"protocol_compliant"`
	StrictCorrect         bool                      `json:"strict_correct"`
	ProcessBefore         *resources.ProcessMetrics `json:"process_before"`
	ProcessAfter          *resources.ProcessMetrics `json:"process_after"`

	Task                  string           `json:"task"`
	Correct               bool             `json:"correct"`
	Error                 string           `json:"error,omitempty"`
	Result                inference.Result `json:"result"`
	OutputTokensPerSecond *float64         `json:"output_tokens_per_second"`
}
type LocalReport struct {
	Provenance Provenance          `json:"provenance"`
	Variance   map[string]Variance `json:"variance"`
	Profile    string              `json:"profile"`

	Version       int        `json:"version"`
	OS            string     `json:"os"`
	Arch          string     `json:"arch"`
	CPUs          int        `json:"logical_cpus"`
	ResourceScope string     `json:"resource_scope"`
	StartedAt     time.Time  `json:"started_at"`
	Rows          []LocalRow `json:"rows"`
	Acceptance    string     `json:"acceptance"`
}

func codeShape(text string) bool {
	f, e := parser.ParseFile(token.NewFileSet(), "candidate.go", text, 0)
	if e != nil || f.Name.Name != "candidate" || len(f.Decls) != 1 {
		return false
	}
	fn, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok || fn.Name.Name != "Add" || fn.Body == nil || len(fn.Body.List) != 1 || fn.Recv != nil || fn.Type.TypeParams != nil {
		return false
	}
	ret, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	bin, ok := ret.Results[0].(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return false
	}
	a, aok := bin.X.(*ast.Ident)
	b, bok := bin.Y.(*ast.Ident)
	if !aok || !bok || a.Name != "a" || b.Name != "b" {
		return false
	}
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}
	typ, ok := fn.Type.Results.List[0].Type.(*ast.Ident)
	if !ok || typ.Name != "int" {
		return false
	}
	names := []string{}
	for _, field := range fn.Type.Params.List {
		typ, ok := field.Type.(*ast.Ident)
		if !ok || typ.Name != "int" {
			return false
		}
		for _, n := range field.Names {
			names = append(names, n.Name)
		}
	}
	return len(names) == 2 && names[0] == "a" && names[1] == "b"
}

// RunLocal uses fixed public microfixtures. It executes no generated code or tools.
func RunLocal(ctx context.Context, b inference.InferenceBackend, h resources.Provider, repeats int) (LocalReport, error) {
	out := LocalReport{Version: 3, OS: runtime.GOOS, Arch: runtime.GOARCH, CPUs: runtime.NumCPU(), ResourceScope: "OS host, not per-process/container allocation", StartedAt: time.Now().UTC(), Acceptance: "external_required"}
	if repeats < 1 || repeats > 3 {
		return out, fmt.Errorf("repeats must be 1..3")
	}
	if e := b.Health(ctx); e != nil {
		return out, e
	}
	fixtures := codingFixtures()

	for n := 0; n < repeats; n++ {
		plan := ""
		for _, f := range fixtures {
			if f.id == "plan-implementation" {
				f.prompt += "\nProposed plan (untrusted):\n" + plan
			}
			if e := ctx.Err(); e != nil {
				return out, e
			}
			before, e := h.Snapshot()
			if e != nil {
				return out, e
			}
			if e = resources.Check(resources.Policy{MinRAM: 256 << 20}, before); e != nil {
				return out, e
			}
			r := inference.Request{ID: fmt.Sprintf("local-%d-%s", n, f.id), ParentWorkID: "local-public-microfixtures", Role: "code", Context: f.prompt, MaxContextBytes: 4096, MaxOutputTokens: 128, Deadline: time.Now().Add(60 * time.Second), Workspace: "no-workspace-access", Capabilities: []string{"code"}}
			var processBefore *resources.ProcessMetrics
			var processErr error
			if observer, ok := h.(processObserver); ok {
				processBefore, processErr = observer.ProcessSnapshot()
			}
			sampler := startResourceSampler(ctx, h, before)
			result, e := b.Invoke(ctx, r)
			peak, samplingErr := sampler.finish(nil)
			after, sampleErr := h.Snapshot()
			if sampleErr == nil {
				merged := mergePeak(*peak, after)
				peak = &merged
			}
			row := LocalRow{Task: f.id, Result: result, Repetition: n, ProcessBefore: processBefore}
			if samplingErr != nil {
				message := samplingErr.Error()
				row.ResourceSamplingError = &message
			}
			if observer, ok := h.(processObserver); ok {
				row.ProcessAfter, processErr = observer.ProcessSnapshot()
			}
			if processErr != nil {
				message := processErr.Error()
				row.ProcessError = &message
			}
			if f.id == "plan" {
				plan = result.Text
			}
			row.Result.Telemetry.Before = &before
			row.Result.Telemetry.Peak = peak
			if sampleErr == nil {
				row.Result.Telemetry.After = &after
			}
			if e != nil {
				row.Error = e.Error()
			} else if sampleErr != nil {
				row.Error = sampleErr.Error()
			} else {
				row.SemanticCorrect = f.check(semanticCandidate(result.Text))
				row.ProtocolCompliant = f.protocol(result.Text)
				row.StrictCorrect = row.ProtocolCompliant && row.SemanticCorrect
				row.Correct = row.StrictCorrect
			}
			if result.Usage.OutputTokens != nil && result.Telemetry.GenerationMS != nil && *result.Telemetry.GenerationMS > 0 {
				v := 1000 * float64(*result.Usage.OutputTokens) / float64(*result.Telemetry.GenerationMS)
				row.OutputTokensPerSecond = &v
			}
			out.Rows = append(out.Rows, row)
		}
	}
	out.Variance = summarize(out.Rows)
	return out, nil
}
