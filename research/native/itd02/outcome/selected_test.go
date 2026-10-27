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

func TestRewardRecommendationCustodyAndEvidence(t *testing.T){
 s:=NewSelected()
 if _,err:=s.RecommendRewardAction(2);err==nil{t.Fatal("accepted absent evidence")}
 s.Rewards.Bins=[3]Bin{{Count:2,Sum:2},{Count:2,Sum:8},{Count:2,Sum:4}}
 s.Rewards.NextID=7
 a,err:=s.RecommendRewardAction(2)
 if err!=nil||a!=1{t.Fatalf("expected action 1, got %d %v",a,err)}
 s.Rewards.Bins[2].Sum=8
 if _,err:=s.RecommendRewardAction(2);err==nil{t.Fatal("accepted reward tie")}
 s.Rewards.Bins[2].Sum=4
 s.Policy.Armed=true
 if _,err:=s.RecommendRewardAction(2);err==nil{t.Fatal("bypassed pending selected action")}
 s.Policy.Armed=false
 s.Rewards.Closed=true
 if _,err:=s.RecommendRewardAction(2);err==nil{t.Fatal("bypassed terminal boundary")}
}

func TestSelectedSnapshotRejectsUnownedRewardEvidence(t *testing.T) {
 s:=NewSelected()
 s.Rewards.Bins[1]=Bin{Count:1,Sum:3}
 s.Rewards.NextID=2
 if _,err:=s.SnapshotSelected();err==nil {
  t.Fatal("accepted reward evidence without matching native selected feedback")
 }
}

func TestRewardAwareUntrainedChoiceUsesNativeCustody(t *testing.T) {
 before:=selectedFixture()
 native:=NewSelected()
 expected,err:=native.Policy.Choose(before)
 if err!=nil{t.Fatal(err)}
 s:=NewSelected()
 chosen,err:=s.ChooseRewardAware(before,2)
 if err!=nil||chosen!=expected{t.Fatalf("untrained native fallback %d vs %d: %v",chosen,expected,err)}
 if !s.Policy.Armed||!s.Policy.Core.Core.Route.Pending||
  s.Policy.Core.Core.Route.Action!=chosen {t.Fatal("native selected-action route not armed")}
 saved,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreSelected(saved);if err!=nil{t.Fatal(err)}
 if !cold.Policy.Armed||cold.Policy.Core.Core.Route.Action!=chosen {
  t.Fatal("native selected-action route lost after restart")
 }
}

func TestRewardAwareRejectsUnownedEvidenceWithoutArming(t *testing.T) {
 s:=NewSelected()
 s.Rewards.Bins=[3]Bin{{Count:2,Sum:2},{Count:2,Sum:8},{Count:2,Sum:4}}
 s.Rewards.NextID=7
 if _,err:=s.ChooseRewardAware(selectedFixture(),2);err==nil {
  t.Fatal("accepted outcome evidence without matching native selected feedback")
 }
 if s.Policy.Armed {t.Fatal("rejected action armed native custodian")}
}

func TestRewardAwareInvalidObservationDoesNotMutateSelectedCustody(t *testing.T) {
 s:=NewSelected()
 before,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 bad:=selectedFixture()
 bad.Direction=9
 if _,err:=s.ChooseRewardAware(bad,2);err==nil {
  t.Fatal("invalid public observation accepted")
 }
 after,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 if !bytes.Equal(before,after) {
  t.Fatal("invalid reward-aware request mutated native or reward state")
 }
}

func TestRewardAwarePendingActionDoesNotMutateState(t *testing.T) {
 s:=NewSelected()
 before:=selectedFixture()
 if _,err:=s.ChooseRewardAware(before,2);err!=nil{t.Fatal(err)}
 saved,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 if _,err:=s.ChooseRewardAware(before,2);err==nil {
  t.Fatal("second choice bypassed pending selected action")
 }
 again,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 if !bytes.Equal(saved,again){
  t.Fatal("rejected pending choice changed checkpoint")
 }
}

func TestRewardAwareSelectedFeedbackAndColdRestart(t *testing.T) {
 s:=NewSelected()
 before:=selectedFixture()
 action,err:=s.ChooseRewardAware(before,2)
 if err!=nil{t.Fatal(err)}
 after:=before
 if action==0 {after.Direction=3} else if action==1 {after.Direction=1}
 if err:=s.Feedback(before,action,after,Event{ID:1,Action:action,Reward:11,Terminal:true});err!=nil{
  t.Fatalf("reward-aware selected feedback rejected: %v",err)
 }
 if s.Policy.Armed||s.Policy.LiveSeen!=1||s.Rewards.NextID!=2||!s.Rewards.Closed {
  t.Fatal("reward-aware selected feedback was not jointly conserved")
 }
 saved,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreSelected(saved);if err!=nil{t.Fatal(err)}
 value,ok:=cold.RewardEvidence(action)
 if !ok||value!=11||cold.Policy.LiveSeen!=1||cold.Policy.Armed {
  t.Fatal("cold restart lost selected outcome or native custody")
 }
 if err:=cold.BeginEpisode();err!=nil{t.Fatal(err)}
 if _,err:=cold.ChooseRewardAware(before,2);err!=nil{
  t.Fatalf("reward-aware next episode rejected: %v",err)
 }
}

func TestBalancedSelectedExplorationChangesRewardChoiceAfterEvidence(t *testing.T) {
 before:=selectedFixture()
 baseline:=NewSelected()
 nativeAction,err:=baseline.Policy.Choose(before);if err!=nil{t.Fatal(err)}
 preferred:=(nativeAction+1)%3
 s:=NewSelected()
 for i:=0;i<6;i++ {
  action,err:=s.ChooseRewardExplore(before,2)
  if err!=nil{t.Fatalf("selected exploration at event %d: %v",i+1,err)}
  after:=before
  if action==0{after.Direction=(before.Direction+3)%4
  }else if action==1{after.Direction=(before.Direction+1)%4}
  reward:=float64(1)
  if action==preferred{reward=9}
  if err:=s.Feedback(before,action,after,Event{
   ID:uint64(i+1),Action:action,Reward:reward,Terminal:true,
  });err!=nil{t.Fatalf("selected feedback at event %d: %v",i+1,err)}
  if err:=s.BeginEpisode();err!=nil{t.Fatalf("episode reset %d: %v",i+1,err)}
 }
 for action,b:=range s.Rewards.Bins{
  if b.Count!=2{t.Fatalf("action %d has %d rather than 2 selected samples",action,b.Count)}
 }
 recommended,err:=s.RecommendRewardAction(2)
 if err!=nil||recommended!=preferred{
  t.Fatalf("learned reward preference %d vs %d: %v",recommended,preferred,err)
 }
 raw,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreSelected(raw);if err!=nil{t.Fatal(err)}
 learned,err:=cold.ChooseRewardExplore(before,2)
 if err!=nil||learned!=preferred||learned==nativeAction{
  t.Fatalf("learned native choice after cold restart: %d baseline %d preferred %d error %v",learned,nativeAction,preferred,err)
 }
 if !cold.Policy.Armed||cold.Policy.Core.Core.Route.Action!=learned{
  t.Fatal("learned native choice failed to arm selected source custody")
 }
}
