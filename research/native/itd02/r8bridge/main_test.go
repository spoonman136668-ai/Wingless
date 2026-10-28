package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "strings"
    "testing"

    r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
    r5 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r5"
    r8 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r8"
)

func wireTestPublic(blocked bool)r1.PublicGrid{
    return r1.PublicGrid{X:20,Y:20,Direction:0,GoalX:21,GoalY:22,FrontBlocked:blocked}
}
func wireTurn(before r1.PublicGrid,arm int,blocked bool)r1.PublicGrid{
    after:=before
    if arm==0{after.Direction=(before.Direction+3)%4}
    if arm==1{after.Direction=(before.Direction+1)%4}
    after.FrontBlocked=blocked
    return after
}
func TestWireMalformedQuitDoesNotTerminatePendingState(t *testing.T){
    // A rejected quit must not discard the live source-custodied learner.
    input:=strings.Join([]string{
        `{"op":"quit","source_private_graph":"must reject"}`,
        `{"op":"snapshot"}`,
        `{"op":"quit"}`,
        "",
    },"\n")
    var output bytes.Buffer
    if err:=serve(strings.NewReader(input),&output);err!=nil{t.Fatal(err)}
    lines:=strings.Split(strings.TrimSpace(output.String()),"\n")
    if len(lines)!=3{t.Fatalf("wrong reply count: %d",len(lines))}
    var replies [3]Reply
    for i,line:=range lines{
        if err:=json.Unmarshal([]byte(line),&replies[i]);err!=nil{t.Fatal(err)}
    }
    if replies[0].OK||replies[0].Error=="" {t.Fatal("private quit was accepted")}
    if !replies[1].OK||replies[1].State==""{t.Fatal("rejected quit stopped bridge")}
    if !replies[2].OK{t.Fatal("valid quit not acknowledged")}
}
func TestNativeR8FitChangesActionAndRetainsAcrossRestart(t *testing.T){
    p:=r8.New()
    s:=wireTestPublic(true)
    // Authored public-shaped engineering fixture: never label as an
    // original MiniGrid experiment or protected scientific evidence.
    for i:=0;i<8;i++{
        for arm:=0;arm<2;arm++{
            reply,quit:=process(&p,Request{Op:"learn",Before:s,After:wireTurn(s,arm,arm==0),Action:arm,Terminal:true})
            if quit||!reply.OK{t.Fatalf("native FIT rejected %d %d: %s",i,arm,reply.Error)}
        }
    }
    if !p.Healthy()||p.FitSeen!=16{t.Fatal("original native source-FIT custody not conserved")}
    original,err:=p.Core.Snapshot()
    if err!=nil{t.Fatal(err)}
    old,err:=r5.Restore(original)
    if err!=nil{t.Fatal(err)}
    oldAction,err:=old.Choose(s)
    if err!=nil{t.Fatal(err)}
    chosen,quit:=process(&p,Request{Op:"choose",Before:s})
    if quit||!chosen.OK||chosen.Action==nil||*chosen.Action!=1{
        t.Fatalf("native R8 did not choose lower learned risk: %+v",chosen)
    }
    if oldAction==*chosen.Action{t.Fatal("R8 decision did not differ from frozen native R5")}
    pre,err:=p.Snapshot()
    if err!=nil{t.Fatal(err)}
    // Failed load or duplicate choice must never lose the pending original action.
    rejected,_:=process(&p,Request{Op:"load",State:base64.StdEncoding.EncodeToString(pre)})
    if rejected.OK||!p.Armed{t.Fatal("load erased selected pending source action")}
    bad,_:=process(&p,Request{Op:"feedback",Before:s,Action:0,After:wireTurn(s,0,true)})
    if bad.OK||!p.Armed{t.Fatal("unselected feedback modified native custody")}
    feedback,_:=process(&p,Request{Op:"feedback",Before:s,Action:1,After:wireTurn(s,1,false)})
    if !feedback.OK||p.LiveSeen!=1{t.Fatalf("genuine selected feedback refused: %+v",feedback)}
    if restart,_:=process(&p,Request{Op:"episode_start"});!restart.OK{t.Fatal(restart.Error)}
    snap,_:=process(&p,Request{Op:"snapshot"})
    if !snap.OK||snap.State==""{t.Fatal("native checkpoint failed")}
    after,_:=p.Snapshot()
    if !p.Healthy(){t.Fatal("native learner invalid after selected feedback")}
    invalid,_:=process(&p,Request{Op:"load",State:"not-valid-base64"})
    if invalid.OK{t.Fatal("malformed checkpoint was accepted")}
    afterInvalid,err:=p.Snapshot()
    if err!=nil||!bytes.Equal(after,afterInvalid){t.Fatal("failed load replaced live model")}
    cold:=r8.New()
    loaded,_:=process(&cold,Request{Op:"load",State:snap.State})
    if !loaded.OK||!cold.Healthy(){t.Fatalf("cold native restore failed: %+v",loaded)}
    choice,_:=process(&cold,Request{Op:"choose",Before:s})
    if !choice.OK||choice.Action==nil||*choice.Action!=1{t.Fatal("cold retained native policy did not choose learned lower-risk action")}
}
func TestPublicWireStrictness(t *testing.T){
    bad:=[]string{
        `{"op":"choose","info":{"privileged_graph":true}}`,
        `{"op":"choose"}{"op":"quit"}`,
        `{"op":"choose",`,
    }
    for _,raw:=range bad{
        if _,err:=decodePublic([]byte(raw));err==nil{t.Fatalf("bad JSON accepted: %s",raw)}
    }
}
