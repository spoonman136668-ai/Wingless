// Package contextbroker retrieves source excerpts, never model authority.
package contextbroker

import (
	"fmt"
	"github.com/spoonman136668-ai/Wingless/ice"
	"os"
	"path/filepath"
	"strings"
)

type Snippet struct {
	File  string `json:"file"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Text  string `json:"text"`
}
type Provider interface {
	Retrieve(query string, maxBytes int) ([]Snippet, error)
}
type ICE struct{ Root string }

func (p ICE) Retrieve(q string, budget int) ([]Snippet, error) {
	if budget <= 0 || budget > 1048576 {
		return nil, fmt.Errorf("invalid context budget")
	}
	hits, e := ice.Query(p.Root, q)
	if e != nil {
		return nil, e
	}
	var out []Snippet
	used := 0
	seen := map[string]bool{}
	for _, s := range hits {
		key := fmt.Sprintf("%s:%d", s.File, s.Start)
		if seen[key] {
			continue
		}
		seen[key] = true
		b, e := os.ReadFile(filepath.Join(p.Root, filepath.FromSlash(s.File)))
		if e != nil {
			return nil, e
		}
		lines := strings.Split(string(b), "\n")
		if s.Start < 1 || s.End > len(lines) {
			return nil, fmt.Errorf("invalid ICE range")
		}
		text := strings.Join(lines[s.Start-1:s.End], "\n")
		if used+len(text) > budget {
			continue
		}
		used += len(text)
		out = append(out, Snippet{s.File, s.Start, s.End, text})
	}
	if e = ice.Validate(p.Root); e != nil {
		return nil, e
	}
	return out, nil
}
func Collect(p Provider, q string, budget int) ([]Snippet, error) {
	if budget <= 0 || budget > 1048576 {
		return nil, fmt.Errorf("invalid context budget")
	}
	out, e := p.Retrieve(q, budget)
	if e != nil {
		return nil, e
	}
	total := 0
	for _, s := range out {
		total += len(s.Text)
		if total > budget {
			return nil, fmt.Errorf("provider exceeded context budget")
		}
	}
	return out, nil
}
