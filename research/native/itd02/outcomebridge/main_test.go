package main

import (
 "bytes"
 "encoding/base64"
 "encoding/json"
 "strings"
 "testing"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 outcome "github.com/spoonman136668-ai/Wingless/research/native/itd02/outcome"
)

func publicGrid()r1.PublicGrid{
 return r1.PublicGrid{X:20,Y:20,Direction:0,GoalX:21,GoalY:22,FrontBlocked:true}
}
func TestContextBridgeSelectedFeedbackAndColdRestart(t *testing.T){
 p:=outcome.NewContextSelected()
 before:=publicGrid()
 chosen,_:=process(&p,Request{Op:"choose",Before:before,MinCount:2})
 if !chosen.OK||chosen.Action==nil{t.Fatalf("choose: %+v",chosen)}
 a:=*chosen.Action
 after:=before
 if a==0{after.Direction=3}else if a==1{after.Direction=1}
 rejected,_:=process(&p,Request{Op:"feedback",Before:before,After:after,Action:(a+1)%3,Reward:9,EventID:1})
 if rejected.OK||p.Base.Policy.LiveSeen!=0||!p.Base.Policy.Armed{
  t.Fatal("unselected feedback altered learner")
 }
 // A pending selected action cannot be replaced by loading a checkpoint.
 initial:=outcome.NewContextSelected()
 raw,err:=initial.Snapshot();if err!=nil{t.Fatal(err)}
 load,_:=process(&p,Request{Op:"load",State:base64.StdEncoding.EncodeToString(raw)})
 if load.OK||!p.Base.Policy.Armed{t.Fatal("load bypassed pending action")}
 applied,_:=process(&p,Request{Op:"feedback",Before:before,After:after,Action:a,Reward:9,EventID:1,Terminal:true})
 if !applied.OK||applied.Updates!=1{t.Fatalf("selected feedback: %+v",applied)}
 saved,_:=process(&p,Request{Op:"snapshot"})
 if !saved.OK||saved.State==""{t.Fatal("checkpoint not exposed")}
 restarted:=outcome.NewContextSelected()
 loaded,_:=process(&restarted,Request{Op:"load",State:saved.State})
 if !loaded.OK||loaded.Updates!=1{t.Fatalf("cold load: %+v",loaded)}
 if restarted.Bins[1][a].Count!=1||restarted.Bins[1][a].Sum!=9{
  t.Fatal("contextual selected reward lost on cold restart")
 }
}
func TestContextBridgeStrictProtocol(t *testing.T){
 if _,err:=decode([]byte(`{"op":"quit","private_map":true}`));err==nil{
  t.Fatal("accepted privileged field")
 }
 if _,err:=decode([]byte(`{"op":"quit"} {"op":"choose"}`));err==nil{
  t.Fatal("accepted multiple operations")
 }
 var input bytes.Buffer
 for _,request:=range []Request{{Op:"snapshot"},{Op:"quit"}}{
  raw,err:=json.Marshal(request);if err!=nil{t.Fatal(err)}
  input.Write(raw);input.WriteByte('\n')
 }
 var output bytes.Buffer
 if err:=serve(&input,&output);err!=nil{t.Fatal(err)}
 if n:=strings.Count(strings.TrimSpace(output.String()),"\n");n!=1{
  t.Fatalf("expected exactly two response lines; got %d",n+1)
 }
}

func TestContextBridgeLineProtocolRoundTripAcrossFreshServer(t *testing.T) {
 before:=publicGrid()
 source:=outcome.NewContextSelected()
 chosen,err:=source.Choose(before,2)
 if err!=nil{t.Fatal(err)}
 after:=before
 if chosen==0{after.Direction=3}else if chosen==1{after.Direction=1}
 requests:=[]Request{
  {Op:"choose",Before:before,MinCount:2},
  {Op:"feedback",Before:before,After:after,Action:chosen,EventID:1,Reward:13,Terminal:true},
  {Op:"snapshot"},
  {Op:"quit"},
 }
 var input bytes.Buffer
 for _,req:=range requests {
  raw,err:=json.Marshal(req);if err!=nil{t.Fatal(err)}
  input.Write(raw);input.WriteByte('\n')
 }
 var output bytes.Buffer
 if err:=serve(&input,&output);err!=nil{t.Fatal(err)}
 var replies []Reply
 for _,line:=range bytes.Split(bytes.TrimSpace(output.Bytes()),[]byte{10}){
  var response Reply
  if err:=json.Unmarshal(line,&response);err!=nil{t.Fatal(err)}
  replies=append(replies,response)
 }
 if len(replies)!=4||!replies[0].OK||replies[0].Action==nil||
  *replies[0].Action!=chosen||!replies[1].OK||
  replies[1].Updates!=1||!replies[2].OK||replies[2].State==""||
  !replies[3].OK {
  t.Fatalf("invalid native wire roundtrip: %+v",replies)
 }
 // A distinct in-memory server must accept the authentic combined checkpoint.
 var freshInput bytes.Buffer
 for _,req:=range []Request{{Op:"load",State:replies[2].State},{Op:"snapshot"},{Op:"quit"}}{
  raw,err:=json.Marshal(req);if err!=nil{t.Fatal(err)}
  freshInput.Write(raw);freshInput.WriteByte('\n')
 }
 var freshOutput bytes.Buffer
 if err:=serve(&freshInput,&freshOutput);err!=nil{t.Fatal(err)}
 lines:=bytes.Split(bytes.TrimSpace(freshOutput.Bytes()),[]byte{10})
 if len(lines)!=3{t.Fatalf("fresh server replied %d times",len(lines))}
 var loaded,restored Reply
 if err:=json.Unmarshal(lines[0],&loaded);err!=nil{t.Fatal(err)}
 if err:=json.Unmarshal(lines[1],&restored);err!=nil{t.Fatal(err)}
 if !loaded.OK||loaded.Updates!=1||!restored.OK||
  restored.State!=replies[2].State {
  t.Fatalf("fresh server failed exact checkpoint retention: %+v %+v",loaded,restored)
 }
}

func TestContextBridgeMalformedWireCannotConsumePendingSelection(t *testing.T){
 before:=publicGrid()
 native:=outcome.NewContextSelected()
 action,err:=native.Choose(before,2)
 if err!=nil{t.Fatal(err)}
 after:=before
 if action==0{after.Direction=3}else if action==1{after.Direction=1}
 requests:=[]string{}
 appendRequest:=func(req Request){
  raw,err:=json.Marshal(req);if err!=nil{t.Fatal(err)}
  requests=append(requests,string(raw))
 }
 appendRequest(Request{Op:"choose",Before:before,MinCount:2})
 requests=append(requests,`{"op":"feedback","private_map":true}`)
 requests=append(requests,`{"op":"feedback"} {"op":"quit"}`)
 appendRequest(Request{Op:"feedback",Before:before,After:after,Action:(action+1)%3,EventID:1,Reward:99})
 appendRequest(Request{Op:"feedback",Before:before,After:after,Action:action,EventID:1,Reward:7,Terminal:true})
 appendRequest(Request{Op:"snapshot"})
 appendRequest(Request{Op:"quit"})
 var result bytes.Buffer
 if err:=serve(strings.NewReader(strings.Join(requests,"\n")+"\n"),&result);err!=nil{t.Fatal(err)}
 lines:=bytes.Split(bytes.TrimSpace(result.Bytes()),[]byte{10})
 if len(lines)!=len(requests){t.Fatalf("expected %d responses, got %d",len(requests),len(lines))}
 replies:=make([]Reply,len(lines))
 for i,line:=range lines{
  if err:=json.Unmarshal(line,&replies[i]);err!=nil{t.Fatalf("response %d: %v",i,err)}
 }
 if !replies[0].OK||replies[0].Action==nil||*replies[0].Action!=action{
  t.Fatal("initial chosen action mismatch")
 }
 for i:=1;i<=3;i++{
  if replies[i].OK||replies[i].Error==""{t.Fatalf("rejected wire %d accepted: %+v",i,replies[i])}
 }
 if !replies[4].OK||replies[4].Updates!=1||!replies[5].OK{
  t.Fatal("valid feedback after rejected requests was not conserved")
 }
 raw,err:=base64.StdEncoding.DecodeString(replies[5].State);if err!=nil{t.Fatal(err)}
 restored,err:=outcome.RestoreContextSelected(raw);if err!=nil{t.Fatal(err)}
 if restored.Bins[1][action].Count!=1||restored.Bins[1][action].Sum!=7{
  t.Fatal("rejected wire altered legitimate selected reward evidence")
 }
}
