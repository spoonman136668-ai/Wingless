package inference

import "testing"

func TestVerificationCandidateSingleFence(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  string
		want string
	}{
		{"json", "```json\n{\"sum\":5}\n```", "{\"sum\":5}"},
		{"go", "```go\npackage candidate\n\nfunc Add(a,b int) int { return a+b }\n```", "package candidate\n\nfunc Add(a,b int) int { return a+b }"},
		{"bare", "```\nplain payload\n```", "plain payload"},
		{"outer-whitespace", " \r\n```json\n{\"ok\":true}\n```\r\n ", "{\"ok\":true}"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, changed := VerificationCandidate(tc.raw)
			if !changed || got != tc.want {
				t.Fatalf("got=%q changed=%v want=%q", got, changed, tc.want)
			}
		})
	}
}

func TestVerificationCandidateFailsClosed(t *testing.T) {
	for _, raw := range []string{
		"raw payload",
		"prefix\n```json\n{\"sum\":5}\n```",
		"```json\n{\"sum\":5}\n```\nsuffix",
		"```python\nprint('x')\n```",
		"```json\n{\"x\":\"```\"}\n```",
		"```json\n{\"sum\":5}",
		"```json\n\n```",
	} {
		got, changed := VerificationCandidate(raw)
		if changed || got != raw {
			t.Fatalf("ambiguous payload changed: raw=%q got=%q changed=%v", raw, got, changed)
		}
	}
}
