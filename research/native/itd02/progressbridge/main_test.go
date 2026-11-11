package main

import (
 "testing"
 outcome "github.com/spoonman136668-ai/Wingless/research/native/itd02/outcome"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)

func TestProgressBridgeSelectedCustodyAndColdRestart(t *testing.T){
 p:=outcome.NewProgressSelected()
 before:=r1.PublicGrid{X:20,Y:20,GoalX:21,GoalY:22,Direction:0,FrontBlocked:false}
 a,_:=process(&p,Request{Op:"choose",Before:before,MinCount:2})
 if !a.OK||a.Action==nil{t.Fatalf("native progress choose rejected: %+v",a)}
 action:=*a.Action
 after:=before
 switch action {
 case 0:after.Direction=3
 case 1:after.Direction=1
 case 2:after.X++
 default:t.Fatal("illegal source action")
 }
 wrong,_:=process(&p,Request{Op:"feedback",Before:before,After:after,Action:(action+1)%3,EventID:1})
 if wrong.OK {t.Fatal("progress bridge accepted unselected selected-action feedback")}
 snap,_:=process(&p,Request{Op:"snapshot"})
 if !snap.OK||snap.State=="" {t.Fatalf("failed selected checkpoint: %+v",snap)}
 good,_:=process(&p,Request{Op:"feedback",Before:before,After:after,Action:action,EventID:1,Terminal:true})
 if !good.OK||good.Updates!=1{t.Fatalf("selected source feedback failed: %+v",good)}
 reset,_:=process(&p,Request{Op:"episode_start"})
 if !reset.OK {t.Fatalf("terminated episode cannot reset: %+v",reset)}
 saved,_:=process(&p,Request{Op:"snapshot"})
 if !saved.OK||saved.State==""||saved.State==snap.State{t.Fatal("genuine selected outcome not preserved")}
 fresh:=outcome.NewProgressSelected()
 loaded,_:=process(&fresh,Request{Op:"load",State:saved.State})
 if !loaded.OK||loaded.Updates!=1{t.Fatalf("cold restore rejected: %+v",loaded)}
 recovered,_:=process(&fresh,Request{Op:"snapshot"})
 if !recovered.OK||recovered.State!=saved.State{t.Fatal("cold restarts not byte-identical")}
 peek,_:=process(&fresh,Request{Op:"peek_frozen",Before:before,MinCount:2})
 if !peek.OK||peek.Action==nil{t.Fatalf("frozen query refused: %+v",peek)}
 again,_:=process(&fresh,Request{Op:"snapshot"})
 if again.State!=saved.State {t.Fatal("frozen query changed persistent learner")}
}
