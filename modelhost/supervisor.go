// Package modelhost supervises one explicitly configured local llama-server.
// It grants no model-requested shell, tool, queue or workspace execution authority.
package modelhost

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/inference"
	"github.com/spoonman136668-ai/Wingless/resources"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Executable       string `json:"executable"`
	ExecutableSHA256 string `json:"executable_sha256"`
	Model            string `json:"model"`
	ModelSHA256      string `json:"model_sha256"`
	ModelID          string `json:"model_id"`
	Port             int    `json:"port"`
	Threads          int    `json:"threads"`
	ContextTokens    int    `json:"context_tokens"`
	StartupSeconds   int    `json:"startup_seconds"`
	RuntimeSeconds   int    `json:"runtime_seconds"`
	MinRAM           uint64 `json:"min_ram"`
}
type Evidence struct {
	State       inference.State `json:"state"`
	PID         int             `json:"pid"`
	StartedAt   time.Time       `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at"`
	Reason      string          `json:"reason"`
	ExitCode    *int            `json:"exit_code"`
	StdoutTail  string          `json:"stdout_tail"`
	StderrTail  string          `json:"stderr_tail"`
}
type tail struct {
	mu sync.Mutex
	b  []byte
}

func (t *tail) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	n := len(p)
	if n >= 8192 {
		t.b = append(t.b[:0], p[n-8192:]...)
	} else {
		t.b = append(t.b, p...)
		if len(t.b) > 8192 {
			t.b = t.b[len(t.b)-8192:]
		}
	}
	return n, nil
}
func (t *tail) text() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return strings.ToValidUTF8(string(t.b), "?")
}

type Supervisor struct {
	started  bool
	mu       sync.Mutex
	cfg      Config
	metrics  resources.Provider
	cancel   context.CancelFunc
	done     chan struct{}
	evidence Evidence
	out, err tail
}

func New(c Config, m resources.Provider) (*Supervisor, error) {
	if m == nil || !filepath.IsAbs(c.Executable) || !filepath.IsAbs(c.Model) || c.ModelID == "" || strings.ContainsAny(c.ModelID, "\r\n") || c.Port < 1024 || c.Port > 65535 || c.Threads < 1 || c.Threads > 16 || c.ContextTokens < 128 || c.ContextTokens > 4096 || c.StartupSeconds < 1 || c.StartupSeconds > 300 || c.RuntimeSeconds < c.StartupSeconds || c.RuntimeSeconds > 900 || c.MinRAM < 1<<30 {
		return nil, fmt.Errorf("invalid bounded model configuration")
	}
	base := strings.ToLower(filepath.Base(c.Executable))
	if base != "llama-server" && base != "llama-server.exe" {
		return nil, fmt.Errorf("only pinned llama-server supported")
	}
	return &Supervisor{cfg: c, metrics: m, evidence: Evidence{State: inference.Stopped}}, nil
}
func pinned(path, want string, max int64) error {
	h, e := hex.DecodeString(want)
	if e != nil || len(h) != 32 {
		return fmt.Errorf("SHA256 pin required")
	}
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil {
		return e
	}
	if !st.Mode().IsRegular() || st.Size() < 1 || st.Size() > max {
		return fmt.Errorf("file size/type denied")
	}
	sum := sha256.New()
	n, e := io.Copy(sum, io.LimitReader(f, max+1))
	if e != nil {
		return e
	}
	if n != st.Size() || !strings.EqualFold(hex.EncodeToString(sum.Sum(nil)), want) {
		return fmt.Errorf("file hash mismatch")
	}
	return nil
}
func (s *Supervisor) Endpoint() string { return "http://127.0.0.1:" + strconv.Itoa(s.cfg.Port) }
func (s *Supervisor) Status() Evidence {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := s.evidence
	e.StdoutTail = s.out.text()
	e.StderrTail = s.err.text()
	if e.ExitCode != nil {
		v := *e.ExitCode
		e.ExitCode = &v
	}
	if e.CompletedAt != nil {
		v := *e.CompletedAt
		e.CompletedAt = &v
	}
	return e
}
func (s *Supervisor) Start(parent context.Context) error {
	s.mu.Lock()
	if s.started || s.evidence.State != inference.Stopped {
		s.mu.Unlock()
		return fmt.Errorf("model already started or failed; create a new supervisor")
	}
	s.started = true
	s.evidence.State = inference.Starting
	s.mu.Unlock()
	fail := func(e error) error {
		s.mu.Lock()
		s.evidence.State = inference.Failed
		s.evidence.Reason = e.Error()
		s.mu.Unlock()
		return e
	}
	if e := parent.Err(); e != nil {
		return fail(e)
	}
	if e := pinned(s.cfg.Executable, s.cfg.ExecutableSHA256, 512<<20); e != nil {
		return fail(e)
	}
	if e := pinned(s.cfg.Model, s.cfg.ModelSHA256, 2<<30); e != nil {
		return fail(e)
	}
	metrics, e := s.metrics.Snapshot()
	if e != nil {
		return fail(e)
	}
	if e = resources.Check(resources.Policy{MinRAM: s.cfg.MinRAM}, metrics); e != nil {
		return fail(e)
	}
	conn, e := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(s.cfg.Port)), 200*time.Millisecond)
	if e == nil {
		conn.Close()
		return fail(fmt.Errorf("configured port already occupied"))
	}
	ctx, cancel := context.WithTimeout(parent, time.Duration(s.cfg.RuntimeSeconds)*time.Second)
	cmd := exec.CommandContext(ctx, s.cfg.Executable, "--model", s.cfg.Model, "--alias", s.cfg.ModelID, "--host", "127.0.0.1", "--port", strconv.Itoa(s.cfg.Port), "--threads", strconv.Itoa(s.cfg.Threads), "--ctx-size", strconv.Itoa(s.cfg.ContextTokens), "--parallel", "1", "--n-gpu-layers", "0")
	// Never inherit API keys, proxies or arbitrary model runtime flags.
	cmd.Env = []string{}
	for _, key := range []string{"SystemRoot", "WINDIR", "TEMP", "TMP"} {
		if v := os.Getenv(key); v != "" {
			cmd.Env = append(cmd.Env, key+"="+v)
		}
	}
	cmd.Dir = filepath.Dir(s.cfg.Executable)
	cmd.Stdout = &s.out
	cmd.Stderr = &s.err
	cmd.WaitDelay = time.Second
	s.mu.Lock()
	if s.evidence.State != inference.Starting {
		s.mu.Unlock()
		cancel()
		return fmt.Errorf("startup stopped")
	}
	if e = cmd.Start(); e != nil {
		s.mu.Unlock()
		cancel()
		return fail(e)
	}
	s.cancel = cancel
	s.done = make(chan struct{})
	done := s.done
	s.evidence.PID = cmd.Process.Pid
	s.evidence.StartedAt = time.Now().UTC()
	s.mu.Unlock()
	go func() {
		e := cmd.Wait()
		s.mu.Lock()
		defer s.mu.Unlock()
		now := time.Now().UTC()
		code := cmd.ProcessState.ExitCode()
		s.evidence.CompletedAt = &now
		s.evidence.ExitCode = &code
		if s.evidence.Reason == "operator_stop" {
			s.evidence.State = inference.Stopped
		} else {
			s.evidence.State = inference.Failed
			switch {
			case s.evidence.Reason != "":
			case ctx.Err() == context.DeadlineExceeded:
				s.evidence.Reason = "runtime_timeout"
			case ctx.Err() != nil:
				s.evidence.Reason = "parent_canceled"
			case e != nil:
				s.evidence.Reason = "process_failed"
			default:
				s.evidence.Reason = "process_exited"
			}
		}
		cancel()
		close(done)
	}()
	ready, stop := context.WithTimeout(ctx, time.Duration(s.cfg.StartupSeconds)*time.Second)
	defer stop()
	backend, e := inference.NewLocalHTTP("supervised", s.cfg.ModelID, s.Endpoint(), []string{"code"})
	if e != nil {
		cancel()
		<-done
		return e
	}
	defer backend.Close()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		if e = backend.Health(ready); e == nil {
			s.mu.Lock()
			if s.evidence.State == inference.Starting {
				s.evidence.State = inference.Ready
				s.mu.Unlock()
				return nil
			}
			s.mu.Unlock()
			return fmt.Errorf("process exited during startup")
		}
		select {
		case <-done:
			return fmt.Errorf("model exited during startup")
		case <-ready.Done():
			s.mu.Lock()
			if ready.Err() == context.DeadlineExceeded && ctx.Err() == nil {
				s.evidence.Reason = "startup_timeout"
			}
			s.mu.Unlock()
			cancel()
			<-done
			return fmt.Errorf("startup timeout/cancellation: %w", ready.Err())
		case <-ticker.C:
		}
	}
}
func (s *Supervisor) Stop(ctx context.Context) error {
	s.mu.Lock()
	if s.evidence.State == inference.Starting && s.cancel == nil {
		s.evidence.State = inference.Stopped
		s.evidence.Reason = "operator_stop"
		s.mu.Unlock()
		return nil
	}
	if s.done == nil {
		s.mu.Unlock()
		return nil
	}
	done := s.done
	select {
	case <-done:
		s.mu.Unlock()
		return nil
	default:
	}
	s.evidence.State = inference.Draining
	s.evidence.Reason = "operator_stop"
	s.cancel()
	s.mu.Unlock()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
