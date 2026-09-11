package inference

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSSEFragmentationAndBounds(t *testing.T) {
	good := "data: {\"model\":\"fixture\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"é\"}}]}\n\ndata: {\"model\":\"fixture\",\"choices\":[{\"delta\":{},\"finish_reason\":\"stop\"}]}\n\ndata: [DONE]\n\n"
	b, _ := NewStreamingLocalHTTP("local", "fixture", "http://127.0.0.1:12345", nil)
	defer b.Close()
	r := request()
	out, e := b.readStream(oneByte{strings.NewReader(good)}, r, r.Deadline.Add(-1e9))
	if e != nil || out.Text != "é" || out.Usage.OutputTokens != nil {
		t.Fatal(out, e)
	}
	for _, bad := range []string{strings.Replace(good, "é", string([]byte{255}), 1), strings.Replace(good, "data: [DONE]", "data: {\"model\":\"fixture\",\"choices\":[{\"finish_reason\":\"stop\"}]}\n\ndata: [DONE]", 1), strings.Repeat(":padding\n", 10000), strings.Replace(good, "é", strings.Repeat("x", 1000), 1)} {
		if _, e = b.readStream(strings.NewReader(bad), r, r.Deadline); e == nil {
			t.Fatal("malformed/oversize accepted")
		}
	}
}

type oneByte struct{ io.Reader }

func (r oneByte) Read(p []byte) (int, error) { return r.Reader.Read(p[:1]) }
func TestOversizeHealthAndConstraintUnavailable(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, strings.Repeat(" ", 65537)) }))
	defer s.Close()
	b, _ := NewLocalHTTP("local", "fixture", s.URL, nil)
	defer b.Close()
	if b.Health(context.Background()) == nil {
		t.Fatal("oversize health accepted")
	}
	r := request()
	r.OutputConstraint = &OutputConstraint{Kind: "json_schema", Schema: []byte(`{"type":"object"}`)}
	if _, e := b.Invoke(context.Background(), r); e == nil {
		t.Fatal("unsupported constraint ignored")
	}
}

func TestDuplicateCompletionAndUsage(t *testing.T) {
	b, _ := NewStreamingLocalHTTP("local", "fixture", "http://127.0.0.1:12345", nil)
	defer b.Close()
	finish := "data: {\"model\":\"fixture\",\"choices\":[{\"finish_reason\":\"stop\"}]}\n\n"
	usage := "data: {\"model\":\"fixture\",\"choices\":[],\"usage\":{\"completion_tokens\":1}}\n\n"
	for _, s := range []string{finish + "data: [DONE]\n\ndata: [DONE]\n\n", finish + usage + usage + "data: [DONE]\n\n"} {
		if _, e := b.readStream(strings.NewReader(s), request(), request().Deadline); e == nil {
			t.Fatal("duplicate accepted")
		}
	}
}
