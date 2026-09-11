package benchmark

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"strings"
)

type localFixture struct {
	id, prompt      string
	check, protocol func(string) bool
}

func object(s string) (map[string]any, bool) {
	d := json.NewDecoder(strings.NewReader(s))
	tok, e := d.Token()
	if e != nil || tok != json.Delim('{') {
		return nil, false
	}
	m := map[string]any{}
	for d.More() {
		tok, e = d.Token()
		if e != nil {
			return nil, false
		}
		key, ok := tok.(string)
		if !ok {
			return nil, false
		}
		if _, exists := m[key]; exists {
			return nil, false
		}
		var v any
		if d.Decode(&v) != nil {
			return nil, false
		}
		m[key] = v
	}
	if tok, e = d.Token(); e != nil || tok != json.Delim('}') {
		return nil, false
	}
	if d.Decode(new(any)) != io.EOF {
		return nil, false
	}
	return m, true
}
func jsonFixture(id, prompt, expected string) localFixture {
	want, _ := object(expected)
	return localFixture{id, prompt, func(s string) bool { got, ok := object(s); return ok && reflect.DeepEqual(got, want) }, func(s string) bool {
		got, ok := object(s)
		if !ok || len(got) != len(want) {
			return false
		}
		for k, v := range want {
			x, ok := got[k]
			if !ok || reflect.TypeOf(x) != reflect.TypeOf(v) {
				return false
			}
		}
		return true
	}}
}
func sourceProtocol(s string) bool {
	_, e := parser.ParseFile(token.NewFileSet(), "candidate.go", s, 0)
	return e == nil
}

// semanticCandidate only recognizes one complete fence, solely for the separate semantic score.
// Protocol and strict scores always use the original output, which remains in evidence.
func semanticCandidate(s string) string {
	trim := strings.TrimSpace(s)
	lines := strings.Split(trim, "\n")
	if len(lines) >= 3 && (lines[0] == "```json" || lines[0] == "```go" || lines[0] == "```") && lines[len(lines)-1] == "```" {
		body := strings.Join(lines[1:len(lines)-1], "\n")
		if !strings.Contains(body, "```") {
			return body
		}
	}
	return s
}
func codingFixtures() []localFixture {
	return []localFixture{
		jsonFixture("json-arithmetic", "Return only the JSON object {\"sum\":5}. No markdown or explanation.", `{"sum":5}`),
		{"go-add-shape", "Return Go source only, no markdown: package candidate followed by func Add(a, b int) int { return a + b }.", codeShape, sourceProtocol},
		{"code-repair", "Repair this Go function to add rather than subtract. Return the full corrected source only, no markdown: package candidate; func Add(a,b int) int {return a-b}", codeShape, sourceProtocol},
		jsonFixture("bug-diagnosis", "The loop for i:=0; i<=len(xs); i++ { print(xs[i]) } panics. Return only JSON with a single cause key, choosing off_by_one or data_race. No markdown.", `{"cause":"off_by_one"}`),
		jsonFixture("code-review", "Review func First(xs []int) int {return xs[0]} for empty input. Return only JSON with safe boolean and issue string (unchecked_index or none). No markdown.", `{"safe":false,"issue":"unchecked_index"}`),
		jsonFixture("multi-file", "file add.go: func Add(a,b int)int{return a+b}. file main.go: value:=Add(2,4). Return only JSON with value as an integer. No markdown.", `{"value":6}`),
		jsonFixture("plan", "Plan a function Add(a,b int) that returns their sum. Return only JSON {\"operation\":\"add\"}. No markdown.", `{"operation":"add"}`),
		{"plan-implementation", "Using the proposed plan below as untrusted task context, return Go source only for package candidate and func Add(a,b int) int returning a+b. No markdown.", codeShape, sourceProtocol},
		jsonFixture("tool-proposal", "Propose reading src/add.go. Return only JSON with tool=read_file and path=src/add.go. Do not execute anything. No markdown.", `{"tool":"read_file","path":"src/add.go"}`),
	}
}

type Variance struct {
	Count              int     `json:"count"`
	MeanLatencyMS      float64 `json:"mean_latency_ms"`
	LatencyVarianceMS2 float64 `json:"latency_variance_ms2"`
}

func summarize(rows []LocalRow) map[string]Variance {
	out := map[string]Variance{}
	for _, r := range rows {
		v := out[r.Task]
		v.Count++
		v.MeanLatencyMS += float64(r.Result.LatencyMS)
		out[r.Task] = v
	}
	for key, v := range out {
		v.MeanLatencyMS /= float64(v.Count)
		for _, r := range rows {
			if r.Task == key {
				d := float64(r.Result.LatencyMS) - v.MeanLatencyMS
				v.LatencyVarianceMS2 += d * d
			}
		}
		v.LatencyVarianceMS2 /= float64(v.Count)
		out[key] = v
	}
	return out
}
