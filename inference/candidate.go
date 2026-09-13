package inference

import "strings"

// VerificationCandidate derives a bounded candidate for trusted deterministic
// verification without changing raw model output. It recognizes only one
// complete outer Markdown fence for the coding payload forms Wingless uses.
// Anything ambiguous, nested, prefixed, suffixed, or differently tagged is
// returned unchanged and is not considered normalized.
func VerificationCandidate(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	lines := strings.Split(trimmed, "\n")
	if len(lines) < 3 {
		return raw, false
	}
	switch lines[0] {
	case "```", "```go", "```json":
	default:
		return raw, false
	}
	if lines[len(lines)-1] != "```" {
		return raw, false
	}
	body := strings.Join(lines[1:len(lines)-1], "\n")
	if strings.Contains(body, "```") || strings.TrimSpace(body) == "" {
		return raw, false
	}
	return body, true
}
