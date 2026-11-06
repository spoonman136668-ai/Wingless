package outcome

import (
 "bytes"
 "testing"
)

// Synthetic public-shaped feedback isolates the control arbitration. It is
// NOT authentic source evidence and does not establish navigation success.
func TestProgressV4SelectedStagnationPermitsTurnDespiteGlobalForwardReward(t *testing.T) {
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=true // surrounded source-shaped cell; turns also face walls
 before.GoalX=23
 before.GoalY=25 // nonzero lateral goal in every heading
 forwardRewards:=0
 for i:=0;i<8;i++ {
  action,err:=p.Choose(before,2)
  if err!=nil {t.Fatalf("choose %d: %v",i,err)}
  after:=before
  switch action {
  case 0: after.Direction=(before.Direction+3)%4
  case 1: after.Direction=(before.Direction+1)%4
  case 2: forwardRewards++
  default:t.Fatal("illegal selected action")
  }
  reward:=0.0
  if action==2 {reward=10} // deliberately synthetic reward-priority stress
  if err=p.Feedback(before,action,after,Event{ID:uint64(i+1),Action:action,Reward:reward});err!=nil {
   t.Fatalf("authenticated synthetic selected feedback %d: %v",i,err)
  }
  before=after
 }
 if forwardRewards==0||p.Base.Base.Base.Rewards.Bins[2].Sum==0{
  t.Fatalf("fixture never established positive global forward reward: %+v",p.Base.Base.Base.Rewards.Bins)
 }
 if p.Stalled!=8||!p.stalledProbe(before,8){
  t.Fatalf("authentic selected stagnation not observed: %d",p.Stalled)
 }
 bins:=p.Base.Base.Base.Rewards.Bins
 expected:=0
 if bins[1].Count<bins[0].Count||bins[0].Count==bins[1].Count{
  expected=1 // completed/8 = 1, deterministic tie rotation
 }
 action,err:=p.Choose(before,2)
 if err!=nil||action!=expected{
  t.Fatalf("stalled evidence probe %d expected %d despite positive forward reward: %v",action,expected,err)
 }
 raw,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreProgressSelected(raw)
 if err!=nil||cold.Stalled!=8||!cold.Base.Base.Base.Policy.Armed{t.Fatal("pending evidence or cold stagnation lost")}
 roundTrip,err:=cold.Snapshot()
 if err!=nil||!bytes.Equal(raw,roundTrip){t.Fatal("cold checkpoint was not canonical")}
 // Original R8 still owns the ONLY pending selected action.
 if cold.Base.Base.Base.Policy.Core.Core.Route.Action!=action{
  t.Fatal("stagnation probe bypassed native selected-action custody")
 }
 wrong:=(action+1)%3
 after:=before
 if action==0{after.Direction=(before.Direction+3)%4}else if action==1{after.Direction=(before.Direction+1)%4}
 if err=p.Feedback(before,wrong,after,Event{ID:9,Action:wrong});err==nil{
  t.Fatal("accepted unselected source feedback")
 }
 unchanged,err:=p.Snapshot()
 if err!=nil||!bytes.Equal(unchanged,raw){t.Fatal("invalid feedback mutated checkpoint or stall ledger")}
}

func TestProgressV4StagnationNeedsObservedFailureAndPublicTurnRelevance(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=false
 before.Direction=0
 before.GoalX=before.X+5
 before.GoalY=before.Y
 p.Stalled=4
 if p.stalledProbe(before,8){t.Fatal("forced turn for directly aligned unblocked goal")}
 before.GoalY=before.Y+2
 if !p.stalledProbe(before,8){t.Fatal("missing lateral-goal stagnant probe")}
 if p.stalledProbe(before,7){t.Fatal("unbounded probe cadence")}
 p.Stalled=3
 if p.stalledProbe(before,8){t.Fatal("probed before four actual stalled transitions")}
 before.GoalY=before.Y
 before.FrontBlocked=true
 p.Stalled=4
 if !p.stalledProbe(before,8){t.Fatal("blocked source-public front cannot trigger evidence probe")}
}

func TestProgressV4NoStagnationAfterEpisodeReset(t *testing.T){
 p:=NewProgressSelected()
 before:=selectedFixture()
 before.FrontBlocked=true
 action,err:=p.Choose(before,2);if err!=nil{t.Fatal(err)}
 after:=before
 if action==0{after.Direction=(before.Direction+3)%4}else if action==1{after.Direction=(before.Direction+1)%4}
 if err=p.Feedback(before,action,after,Event{ID:1,Action:action,Terminal:true});err!=nil{t.Fatal(err)}
 if p.Stalled!=1{t.Fatalf("selected nonprogress event did not count: %d",p.Stalled)}
 raw,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreProgressSelected(raw)
 if err!=nil||cold.Stalled!=1{t.Fatalf("terminal checkpoint lost observed stall: %v",err)}
 if err=p.BeginEpisode();err!=nil{t.Fatal(err)}
 if p.Stalled!=0||p.stalledProbe(before,8){t.Fatal("stale test-episode progress leaked into next episode")}
}

func TestProgressV4CheckpointRejectsPriorSchemaAndCorruption(t *testing.T) {
 p:=NewProgressSelected()
 raw,err:=p.Snapshot();if err!=nil{t.Fatal(err)}
 historic:=bytes.Replace(raw,[]byte("selected-goal-progress.v4"),[]byte("selected-goal-progress.v3"),1)
 if bytes.Equal(raw,historic){t.Fatal("source schema marker absent")}
 if _,err=RestoreProgressSelected(historic);err==nil{t.Fatal("v3 checkpoint silently promoted to v4")}
 tampered:=append([]byte(nil),raw...)
 tampered[len(tampered)/2]^=1
 if _,err=RestoreProgressSelected(tampered);err==nil{t.Fatal("tampered v4 evidence accepted")}
}
