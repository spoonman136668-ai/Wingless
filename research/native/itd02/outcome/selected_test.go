package outcome

import (
 "bytes"
 "testing"

 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)

func selectedFixture()r1.PublicGrid{
 return r1.PublicGrid{X:20,Y:20,Direction:0,GoalX:21,GoalY:22,FrontBlocked:true}
}
func TestSelectedFeedbackAtomicity(t *testing.T){
 s:=NewSelected()
 before:=selectedFixture()
 a,err:=s.Policy.Choose(before);if err!=nil{t.Fatal(err)}
 after:=before
 if a==0{after.Direction=3}else if a==1{after.Direction=1}
 pre,err:=s.Policy.Snapshot();if err!=nil{t.Fatal(err)}
 lpre,err:=s.Rewards.Snapshot();if err!=nil{t.Fatal(err)}
 for _,e:=range []Event{
  {ID:1,Action:(a+1)%3,Reward:2},
  {ID:1,Action:a,Reward:RewardLimit+1},
  {ID:2,Action:a,Reward:2},
 }{
  if err:=s.Feedback(before,a,after,e);err==nil{t.Fatalf("accepted invalid selected outcome: %+v",e)}
  now,_:=s.Policy.Snapshot();if !bytes.Equal(now,pre){t.Fatal("invalid event changed pending native learner")}
  now,_=s.Rewards.Snapshot();if !bytes.Equal(now,lpre){t.Fatal("invalid event changed reward ledger")}
 }
 if err:=s.Feedback(before,a,after,Event{ID:1,Action:a,Reward:5,Terminal:true});err!=nil{t.Fatal(err)}
 v,ok:=s.Rewards.Value(a);if !ok||v!=5{t.Fatalf("missing selected reward: %v %v",v,ok)}
 if s.Policy.Armed||s.Policy.LiveSeen!=1{t.Fatal("native selected feedback not accepted exactly once")}
 if err:=s.Feedback(before,a,after,Event{ID:1,Action:a,Reward:5});err==nil{t.Fatal("duplicate accepted")}
 if err:=s.BeginEpisode();err!=nil{t.Fatal(err)}
 if s.Rewards.Closed{t.Fatal("terminal reset not applied")}
}

func TestUnselectedFeedbackDoesNotConsumeReward(t *testing.T){
 s:=NewSelected();before:=selectedFixture()
 a,err:=s.Policy.Choose(before);if err!=nil{t.Fatal(err)}
 after:=before
 if err:=s.Feedback(before,(a+1)%3,after,Event{ID:1,Action:(a+1)%3,Reward:9});err==nil{t.Fatal("unselected action accepted")}
 if s.Rewards.NextID!=1||!s.Policy.Armed{t.Fatal("rejected action consumed state")}
}
