// Package toolboundary deliberately exposes only a read-only, rooted file tool.
package toolboundary

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Proposal struct {
	Tool string `json:"tool"`
	Path string `json:"path"`
}
type Policy struct {
	Workspace      string
	AllowedTools   []string
	AllowedPaths   []string
	Timeout        time.Duration
	MaxOutputBytes int
}
type Observation struct {
	Tool      string `json:"tool"`
	Path      string `json:"path"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitCode  int    `json:"exit_code"`
	ElapsedMS int64  `json:"elapsed_ms"`
	Reason    string `json:"reason"`
}

func has(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}
func Execute(parent context.Context, p Policy, a Proposal) (o Observation, err error) {
	o = Observation{Tool: a.Tool, Path: a.Path, ExitCode: -1}
	start := time.Now()
	defer func() {
		o.ElapsedMS = time.Since(start).Milliseconds()
		if err != nil {
			o.Reason = err.Error()
		}
	}()
	if a.Tool != "read_file" || !has(p.AllowedTools, a.Tool) || !has(p.AllowedPaths, a.Path) {
		err = fmt.Errorf("tool/path unauthorized")
		return
	}
	if p.Timeout <= 0 || p.Timeout > time.Minute || p.MaxOutputBytes <= 0 || p.MaxOutputBytes > 1048576 || !filepath.IsAbs(p.Workspace) {
		err = fmt.Errorf("unbounded tool policy")
		return
	}
	if !filepath.IsLocal(a.Path) || strings.ContainsAny(a.Path, "\\:") {
		err = fmt.Errorf("workspace escape denied")
		return
	}
	ctx, c := context.WithTimeout(parent, p.Timeout)
	defer c()
	if err = ctx.Err(); err != nil {
		return
	}
	root, e := os.OpenRoot(p.Workspace)
	if e != nil {
		err = e
		return
	}
	defer root.Close()
	file, e := root.Open(a.Path)
	if e != nil {
		err = fmt.Errorf("rooted read denied: %w", e)
		return
	}
	defer file.Close()
	info, e := file.Stat()
	if e != nil {
		err = e
		return
	}
	if !info.Mode().IsRegular() {
		err = fmt.Errorf("nonregular file denied")
		return
	}
	b, e := io.ReadAll(io.LimitReader(file, int64(p.MaxOutputBytes)+1))
	if e != nil {
		err = e
		return
	}
	if err = ctx.Err(); err != nil {
		return
	}
	if len(b) > p.MaxOutputBytes {
		err = fmt.Errorf("output bound exceeded")
		return
	}
	o.Stdout = string(b)
	o.ExitCode = 0
	o.Reason = "authorized_read"
	return
}
