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
