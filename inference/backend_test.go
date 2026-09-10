package inference

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func request() Request {
	return Request{ID: "test", ParentWorkID: "work", Role: "code", Context: "x", MaxContextBytes: 20, MaxOutputTokens: 16, Deadline: time.Now().Add(time.Second), Workspace: "fixture", Capabilities: []string{"code"}}
}
func TestMockCancellation(t *testing.T) {
	for _, kind := range []string{"parent", "deadline", "explicit"} {
		t.Run(kind, func(t *testing.T) {
			m := &Mock{Name: "fast", Delay: time.Second}
			ctx, c := context.WithCancel(context.Background())
			defer c()
			r := request()
			if kind == "parent" {
				c()
			}
			if kind == "deadline" {
				r.Deadline = time.Now().Add(10 * time.Millisecond)
			}
			if kind == "explicit" {
				done := make(chan error, 1)
				go func() { _, e := m.Invoke(ctx, r); done <- e }()
				until := time.Now().Add(time.Second)
				for !m.Cancel(r.ID) {
					if time.Now().After(until) {
						t.Fatal("never active")
					}
					time.Sleep(time.Millisecond)
				}
				if e := <-done; !errors.Is(e, context.Canceled) {
					t.Fatal(e)
				}
				return
			}
			out, e := m.Invoke(ctx, r)
			if e == nil || out.Status == "completed" {
				t.Fatal(out, e)
			}
			want := "canceled"
			if kind == "deadline" {
				want = "worker_timeout"
			}
			if out.ErrorClass != want {
				t.Fatal(out)
			}
		})
	}
}
func TestBounds(t *testing.T) {
	r := request()
	r.Context = "too much context exceeds twenty bytes"
	if _, e := (&Mock{Name: "m"}).Invoke(context.Background(), r); e == nil {
		t.Fatal("context accepted")
	}
}
func TestLifecycle(t *testing.T) {
	l := NewLifecycle()
	if e := l.Move(Busy); e == nil {
		t.Fatal("started busy")
	}
	for _, s := range []State{Starting, Ready, Busy, Idle, Draining, Stopped, Starting, Failed, Stopped, Unavailable} {
		if e := l.Move(s); e != nil {
			t.Fatal(e)
		}
	}
	if l.State() != Unavailable {
		t.Fatal(l.State())
	}
}
func TestLocalTransport(t *testing.T) {
	for _, mode := range []string{"ok", "429", "401", "503", "wrong-model", "truncated", "oversize", "redirect", "cancel"} {
		t.Run(mode, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/v1/models" {
					fmt.Fprint(w, `{"data":[{"id":"fixture"}]}`)
					return
				}
				switch mode {
				case "429":
					w.WriteHeader(429)
					return
				case "401":
					w.WriteHeader(401)
					return
				case "503":
					w.WriteHeader(503)
					return
				case "redirect":
					http.Redirect(w, r, "http://192.0.2.1", 302)
					return
				case "cancel":
					io.Copy(io.Discard, r.Body)
					select {
					case <-r.Context().Done():
					case <-time.After(time.Second):
					}
					return
				case "oversize":
					for i := 0; i < 70000; i++ {
						fmt.Fprint(w, "x")
					}
					return
				}
				model := "fixture"
				finish := "stop"
				if mode == "wrong-model" {
					model = "other"
				}
				if mode == "truncated" {
					finish = "length"
				}
				fmt.Fprintf(w, `{"model":%q,"choices":[{"message":{"content":"ok"},"finish_reason":%q}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`, model, finish)
			}))
			defer server.Close()
			b, e := NewLocalHTTP("local", "fixture", server.URL, []string{"code"})
			if e != nil {
				t.Fatal(e)
			}
			if e = b.Health(context.Background()); e != nil {
				t.Fatal(e)
			}
			req := request()
			if mode == "cancel" {
				req.Deadline = time.Now().Add(20 * time.Millisecond)
			}
			out, e := b.Invoke(context.Background(), req)
			if mode == "ok" {
				if e != nil || out.Text != "ok" || out.Usage.OutputTokens == nil || out.Status != "completed" {
					t.Fatal(out, e)
				}
			} else if e == nil {
				t.Fatal("bad response accepted", out)
			}
		})
	}
}
func TestEndpointBoundary(t *testing.T) {
	for _, endpoint := range []string{"https://example.com", "http://localhost:8080", "http://127.0.0.1@evil.test", "http://127.0.0.1?token=secret", "http://127.0.0.1/api"} {
		if _, e := NewLocalHTTP("local", "m", endpoint, nil); e == nil {
			t.Fatal(endpoint)
		}
	}
}
