package ckbplane

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

type serviceEngineFunc func(context.Context, TransportRequest) (CandidateEnvelope, error)

func (f serviceEngineFunc) Run(ctx context.Context, req TransportRequest) (CandidateEnvelope, error) {
	return f(ctx, req)
}

type serviceFixture struct {
	endpoint string
	cancel   context.CancelFunc
	done     chan error
	req      TransportRequest
	session  string
}

func startServiceFixture(t *testing.T) serviceFixture {
	t.Helper()
	session, req, candidate := listenerFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	reader, writer := io.Pipe()
	done := make(chan error, 1)
	go func() {
		done <- ServeLoopbackService(ctx, session, serviceEngineFunc(func(context.Context, TransportRequest) (CandidateEnvelope, error) {
			return candidate, nil
		}), writer)
		_ = writer.Close()
	}()
	line, err := bufio.NewReader(reader).ReadBytes('\n')
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	_ = reader.Close()
	var ready LoopbackServiceReady
	if err := json.Unmarshal(bytes.TrimSpace(line), &ready); err != nil {
		cancel()
		t.Fatal(err)
	}
	if ready.Schema != LoopbackServiceReadySchema || ready.PID <= 0 || ready.Endpoint == "" {
		cancel()
		t.Fatalf("bad ready announcement: %+v", ready)
	}
	return serviceFixture{endpoint: ready.Endpoint, cancel: cancel, done: done, req: req, session: session}
}

func stopServiceFixture(t *testing.T, f serviceFixture) {
	t.Helper()
	f.cancel()
	select {
	case err := <-f.done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("service stop: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("service stop timed out")
	}
}

func TestWinglessLoopbackServiceReadyContract(t *testing.T) {
	f := startServiceFixture(t)
	defer stopServiceFixture(t, f)
	raw, err := json.Marshal(f.req)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, f.endpoint, bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CKB-Wingless-Session", f.session)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.Header.Get("X-CKB-Wingless-Session") != f.session {
		t.Fatalf("status=%d session=%q", resp.StatusCode, resp.Header.Get("X-CKB-Wingless-Session"))
	}
}

func TestWinglessLoopbackServiceNumericIPv4Only(t *testing.T) {
	f := startServiceFixture(t)
	defer stopServiceFixture(t, f)
	u, err := url.Parse(f.endpoint)
	if err != nil {
		t.Fatal(err)
	}
	host := u.Hostname()
	ip := net.ParseIP(host)
	if u.Scheme != "http" || u.Path != TransportPath || ip == nil || !ip.IsLoopback() || ip.To4() == nil || host != "127.0.0.1" || u.Port() == "" {
		t.Fatalf("nonconforming endpoint %q", f.endpoint)
	}
}

func TestWinglessLoopbackServiceCancellationStopsListener(t *testing.T) {
	f := startServiceFixture(t)
	u, err := url.Parse(f.endpoint)
	if err != nil {
		t.Fatal(err)
	}
	stopServiceFixture(t, f)
	conn, err := net.DialTimeout("tcp4", u.Host, 300*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		t.Fatal("listener still reachable after service cancellation")
	}
}

func TestWinglessLoopbackServiceNoDefaultActivationSurface(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	mainPath := filepath.Join(root, "cmd", "wingless", "main.go")
	raw, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	for _, forbidden := range []string{"ServeLoopbackService", "plane-service", "proofservice", "CKB_WINGLESS_SESSION"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("default Wingless command activates service surface: %s", forbidden)
		}
	}
}
