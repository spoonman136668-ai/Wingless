package outcome

import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "testing"
)

func TestSelectedCheckpointColdRestartAndPendingAction(t *testing.T){
 s:=NewSelected()
 before:=selectedFixture()
 a,err:=s.Policy.Choose(before);if err!=nil{t.Fatal(err)}
 raw,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 cold,err:=RestoreSelected(raw);if err!=nil{t.Fatal(err)}
 if !cold.Policy.Armed||cold.Rewards.NextID!=1{t.Fatal("pending selected action lost across restart")}
 after:=before
 if a==0{after.Direction=3}else if a==1{after.Direction=1}
 if err:=cold.Feedback(before,a,after,Event{ID:1,Action:a,Reward:4,Terminal:true});err!=nil{t.Fatal(err)}
 saved,err:=cold.SnapshotSelected();if err!=nil{t.Fatal(err)}
 again,err:=RestoreSelected(saved);if err!=nil{t.Fatal(err)}
 if again.Policy.LiveSeen!=1||!again.Rewards.Closed{t.Fatal("reward and native evidence not conserved")}
 v,ok:=again.Rewards.Value(a);if !ok||v!=4{t.Fatal("reward lost on cold restart")}
 if err:=again.BeginEpisode();err!=nil{t.Fatal(err)}
 if again.Rewards.Closed{t.Fatal("terminal state did not reset")}
}

func TestSelectedCheckpointRejectsCorruptionAndMismatchedCounts(t *testing.T){
 s:=NewSelected()
 raw,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 corrupted:=append([]byte(nil),raw...)
 corrupted[len(corrupted)/2]^=1
 if _,err:=RestoreSelected(corrupted);err==nil{t.Fatal("corrupt checkpoint accepted")}
 if _,err:=RestoreSelected(append(append([]byte(nil),raw...),byte(' ')));err==nil{t.Fatal("noncanonical checkpoint accepted")}
 var wire selectedEnvelope
 if err:=json.Unmarshal(raw,&wire);err!=nil{t.Fatal(err)}
 l,err:=Restore(wire.State.Rewards);if err!=nil{t.Fatal(err)}
 if err:=l.Accept(Event{ID:1,Action:0,Reward:1});err!=nil{t.Fatal(err)}
 wire.State.Rewards,err=l.Snapshot();if err!=nil{t.Fatal(err)}
 // A valid inner ledger cannot be paired with an unchanged native R8.
 // Even when an attacker recomputes the outer digest, custody must reject.
 body,err:=json.Marshal(wire.State);if err!=nil{t.Fatal(err)}
 wire.Digest=sha256Sum(body)
 mismatched,err:=json.Marshal(wire);if err!=nil{t.Fatal(err)}
 if bytes.Equal(raw,mismatched){t.Fatal("mismatch fixture unchanged")}
 if _,err:=RestoreSelected(mismatched);err==nil{t.Fatal("native/outcome count mismatch accepted")}
}

func sha256Sum(raw []byte)[32]byte{return sha256.Sum256(raw)}

func TestSelectedObservedRewardEvidenceSurvivesColdRestart(t *testing.T) {
 s:=NewSelected()
 before:=selectedFixture()
 action,err:=s.Policy.Choose(before);if err!=nil{t.Fatal(err)}
 after:=before
 if action==0 {after.Direction=3} else if action==1 {after.Direction=1}
 if err:=s.Feedback(before,action,after,Event{ID:1,Action:action,Reward:7,Terminal:true});err!=nil{t.Fatal(err)}
 original,ok:=s.RewardEvidence(action)
 if !ok||original!=7{t.Fatalf("unexpected original reward evidence %v %v",original,ok)}
 raw,err:=s.SnapshotSelected();if err!=nil{t.Fatal(err)}
 restored,err:=RestoreSelected(raw);if err!=nil{t.Fatal(err)}
 got,ok:=restored.RewardEvidence(action)
 if !ok||got!=original{t.Fatalf("cold restart changed selected reward evidence %v %v",got,ok)}
 if _,err:=restored.RecommendRewardAction(2);err==nil {
  t.Fatal("accepted unsupported reward recommendation after restart")
 }
}
