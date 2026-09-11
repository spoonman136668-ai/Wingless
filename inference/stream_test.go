package inference

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStreamingValidation(t *testing.T) {
	for _, mode := range []string{"ok", "truncated", "wrong_model", "length", "budget", "after_finish"} {
		t.Run(mode, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				model := "fixture"
				if mode == "wrong_model" {
					model = "other"
				}
				fmt.Fprintf(w, "data: {\"model\":%q,\"choices\":[{\"index\":0,\"delta\":{\"content\":\"hi\"}}]}\n\n", model)
				w.(http.Flusher).Flush()
				if mode == "truncated" {
					return
				}
				finish := "stop"
				if mode == "length" {
					finish = "length"
				}
				fmt.Fprintf(w, "data: {\"model\":\"fixture\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":%q}]}\n\n", finish)
				if mode == "budget" {
					fmt.Fprint(w, "data: {\"model\":\"fixture\",\"choices\":[],\"usage\":{\"completion_tokens\":999}}\n\n")
				}
				if mode == "after_finish" {
					fmt.Fprint(w, "data: {\"model\":\"fixture\",\"choices\":[{\"delta\":{\"content\":\"bad\"}}]}\n\n")
				}
				fmt.Fprint(w, "data: [DONE]\n\n")
			}))
			defer srv.Close()
			b, e := NewStreamingLocalHTTP("local", "fixture", srv.URL, []string{"code"})
			if e != nil {
				t.Fatal(e)
			}
			defer b.Close()
			out, e := b.Invoke(context.Background(), request())
			if mode == "ok" {
				if e != nil || out.Text != "hi" || out.Telemetry.TimeToFirstTokenMS == nil {
					t.Fatal(out, e)
				}
			} else if e == nil || out.Status == "completed" {
				t.Fatal("invalid stream accepted", out, e)
			}
		})
	}
}

func TestStreamingTimeoutAfterContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "data: {\"model\":\"fixture\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"partial\"}}]}\n\n")
		w.(http.Flusher).Flush()
		select {
		case <-r.Context().Done():
		case <-time.After(time.Second):
		}
	}))
	defer server.Close()
	b, e := NewStreamingLocalHTTP("local", "fixture", server.URL, []string{"code"})
	if e != nil {
		t.Fatal(e)
	}
	defer b.Close()
	req := request()
	req.Deadline = time.Now().Add(30 * time.Millisecond)
	out, e := b.Invoke(context.Background(), req)
	if e == nil || out.Status != "failed" || out.ErrorClass != "worker_timeout" {
		t.Fatal(out, e)
	}
}
