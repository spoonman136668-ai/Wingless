package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeInvoker struct {
	choice modelChoice
	request []byte
	out []byte
	err error
}
func (f *fakeInvoker) Invoke(_ context.Context, req []byte, choice modelChoice)([]byte,error){
	f.choice=choice; f.request=append([]byte(nil),req...); return f.out,f.err
}

func TestSelectModelAdaptive(t *testing.T){
	var small implementationRequest
	small.Schema=requestSchema
	small.Spec.ChangedPaths=[]string{"a","b"}
	small.Context=append(small.Context,struct{Content string `json:"content"`}{Content:"x"})
	got,err:=selectModel(small)
	if err!=nil||got.Model!="gpt-6-luna"||got.Effort!="medium"{t.Fatalf("small=%+v err=%v",got,err)}

	var normal implementationRequest
	normal.Schema=requestSchema
	for i:=0;i<5;i++{normal.Spec.ChangedPaths=append(normal.Spec.ChangedPaths,string(rune('a'+i)))}
	normal.Context=append(normal.Context,struct{Content string `json:"content"`}{Content:strings.Repeat("x",300<<10)})
	got,err=selectModel(normal)
	if err!=nil||got.Model!="gpt-6-sol"||got.Effort!="medium"{t.Fatalf("normal=%+v err=%v",got,err)}

	var hard implementationRequest
	hard.Schema=requestSchema
	for i:=0;i<13;i++{hard.Spec.ChangedPaths=append(hard.Spec.ChangedPaths,"p"+string(rune('a'+i)))}
	got,err=selectModel(hard)
	if err!=nil||got.Model!="gpt-6-sol"||got.Effort!="high"{t.Fatalf("hard=%+v err=%v",got,err)}
}

func TestGatewayLoopbackAndJSONContract(t *testing.T){
	bundle:=[]byte(`{"schema":"ckb-plane.research-implementation-bundle.v1","cycle_id":"c","experiment_id":"e","spec_sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","files":[{"path":"docs/x.md","content":"x"}]}`)
	fake:=&fakeInvoker{out:bundle}
	g:=&gateway{invoker:fake,timeout:0}

	impl:=map[string]any{
		"schema":requestSchema,
		"spec":map[string]any{"changed_paths":[]string{"docs/x.md"}},
		"context":[]map[string]any{{"content":"context"}},
	}
	implRaw,_:=json.Marshal(impl)
	chat:=map[string]any{
		"model":gatewayModel,
		"messages":[]map[string]string{{"role":"system","content":"system"},{"role":"user","content":string(implRaw)}},
		"temperature":0,
		"max_tokens":4096,
		"stream":false,
		"response_format":map[string]string{"type":"json_object"},
	}
	raw,_:=json.Marshal(chat)

	req:=httptest.NewRequest(http.MethodPost,"http://127.0.0.1/v1/chat/completions",bytes.NewReader(raw))
	req.RemoteAddr="203.0.113.4:1234"
	w:=httptest.NewRecorder();g.chat(w,req)
	if w.Code!=http.StatusForbidden{t.Fatalf("nonloopback status=%d",w.Code)}

	req=httptest.NewRequest(http.MethodPost,"http://127.0.0.1/v1/chat/completions",bytes.NewReader(raw))
	req.RemoteAddr="127.0.0.1:1234"
	w=httptest.NewRecorder();g.chat(w,req)
	if w.Code!=http.StatusOK{t.Fatalf("status=%d body=%s",w.Code,w.Body.String())}
	if fake.choice.Model!="gpt-6-luna"||fake.choice.Effort!="medium"{t.Fatalf("choice=%+v",fake.choice)}
	if !strings.Contains(w.Body.String(),"ckb-plane.research-implementation-bundle.v1"){t.Fatalf("body=%s",w.Body.String())}
}

func TestGatewayRejectsNonJSONFinal(t *testing.T){
	fake:=&fakeInvoker{out:[]byte("not-json")}
	g:=&gateway{invoker:fake}
	implRaw:=[]byte(`{"schema":"ckb-plane.research-implementation-request.v1","spec":{"changed_paths":["x"]}}`)
	chat:=map[string]any{
		"model":gatewayModel,
		"messages":[]map[string]string{{"role":"system","content":"system"},{"role":"user","content":string(implRaw)}},
		"temperature":0,"max_tokens":4096,"stream":false,
		"response_format":map[string]string{"type":"json_object"},
	}
	raw,_:=json.Marshal(chat)
	req:=httptest.NewRequest(http.MethodPost,"http://127.0.0.1/v1/chat/completions",bytes.NewReader(raw))
	req.RemoteAddr="127.0.0.1:9"
	w:=httptest.NewRecorder();g.chat(w,req)
	if w.Code!=http.StatusBadGateway{t.Fatalf("status=%d body=%s",w.Code,w.Body.String())}
}
