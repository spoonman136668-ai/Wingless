package outcome

import (
 "math"
 "testing"
)

func TestRewardDependentLearningAndRevision(t *testing.T) {
 a,b:=New(),New()
 for i:=uint64(1);i<=3;i++ {
  if err:=a.Accept(Event{ID:i,Action:0,Reward:2});err!=nil{t.Fatal(err)}
  if err:=b.Accept(Event{ID:i,Action:0,Reward:-2});err!=nil{t.Fatal(err)}
 }
 av,ok:=a.Value(0);if !ok||av!=2{t.Fatalf("positive reward: %v %v",av,ok)}
 bv,ok:=b.Value(0);if !ok||bv!=-2{t.Fatalf("negative reward: %v %v",bv,ok)}
 if err:=a.Accept(Event{ID:4,Action:0,Reward:-10});err!=nil{t.Fatal(err)}
 revised,_:=a.Value(0);if revised>=0{t.Fatalf("contradictory reward did not revise value: %v",revised)}
}

func TestDuplicateAndInvalidEventsFailClosed(t *testing.T) {
 l:=New()
 if err:=l.Accept(Event{ID:1,Action:1,Reward:3});err!=nil{t.Fatal(err)}
 before,_:=l.Snapshot()
 cases:=[]Event{
  {ID:1,Action:1,Reward:99},
  {ID:2,Action:3,Reward:1},
  {ID:2,Action:1,Reward:math.NaN()},
  {ID:2,Action:1,Reward:math.Inf(1)},
  {ID:2,Action:1,Reward:RewardLimit+1},
  {ID:2,Action:1,Terminal:true,Truncated:true},
 }
 for _,e:=range cases {
  if err:=l.Accept(e);err==nil{t.Fatalf("accepted invalid event: %+v",e)}
  after,_:=l.Snapshot();if string(after)!=string(before){t.Fatal("rejected event mutated state")}
 }
}

func TestTerminalResetAndColdRestore(t *testing.T) {
 l:=New()
 if err:=l.Accept(Event{ID:1,Action:2,Reward:7,Terminal:true});err!=nil{t.Fatal(err)}
 if err:=l.Accept(Event{ID:2,Action:2,Reward:9});err==nil{t.Fatal("accepted after terminal")}
 raw,err:=l.Snapshot();if err!=nil{t.Fatal(err)}
 restored,err:=Restore(raw);if err!=nil{t.Fatal(err)}
 if !restored.Closed{t.Fatal("terminal state lost")}
 v,ok:=restored.Value(2);if !ok||v!=7{t.Fatalf("restored value: %v %v",v,ok)}
 if err:=restored.Reset();err!=nil{t.Fatal(err)}
 if err:=restored.Accept(Event{ID:2,Action:2,Reward:-3});err!=nil{t.Fatal(err)}
 v,_=restored.Value(2);if v!=2{t.Fatalf("unexpected mean %v",v)}
 raw[10]^=1
 if _,err:=Restore(raw);err==nil{t.Fatal("accepted tampered checkpoint")}
}
