package inference

import (
	"net/http"
	"testing"
	"time"
)

func TestLocalHTTPInferenceHasNoFixedResponseHeaderTimeout(t *testing.T) {
	backend, err := NewLocalHTTP("local", "model", "http://127.0.0.1:1", []string{"code"})
	if err != nil {
		t.Fatal(err)
	}

	transport, ok := backend.client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected transport type %T", backend.client.Transport)
	}
	if transport.ResponseHeaderTimeout != 0 {
		t.Fatalf("local inference must not impose a fixed response-header timeout; got %s", transport.ResponseHeaderTimeout)
	}
	if backend.client.Timeout != 15*time.Minute {
		t.Fatalf("unexpected client timeout: %s", backend.client.Timeout)
	}
}
