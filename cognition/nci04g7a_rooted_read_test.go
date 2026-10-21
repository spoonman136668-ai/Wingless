package cognition

import (
 "context"
 "math"
 "os"
 "path/filepath"
 "testing"
)

func TestNCI04G7ASourceIndependentNativeReadAndRestart(t *testing.T) {
 var predictions []float64
 var consequences []uint8
 for _,external := range []string{"0","1"} {
  root:=t.TempDir()
  if e:=os.WriteFile(filepath.Join(root,rootedReadSignal),[]byte(external),0600);e!=nil{t.Fatal(e)}
  learner,e:=NewRootedReadLearner(root);if e!=nil{t.Fatal(e)}
  unknown,e:=learner.Predict();if e!=nil||unknown!=.5{t.Fatalf("expected abstaining prior, got %v %v",unknown,e)}
  event,e:=learner.Observe(context.Background());if e!=nil{t.Fatal(e)}
  if event.Sequence!=1||event.Action!="rooted-read:native-state.bit"{t.Fatalf("wrong native event: %#v",event)}
  got,e:=learner.Predict();if e!=nil{t.Fatal(e)}
  receipt,e:=learner.Snapshot();if e!=nil{t.Fatal(e)}
  reopened,e:=NewRootedReadLearner(root);if e!=nil{t.Fatal(e)}
  if e=reopened.Restore(receipt);e!=nil{t.Fatal(e)}
  after,e:=reopened.Predict();if e!=nil||math.Abs(after-got)>1e-12{t.Fatal("native learning changed after restart",e)}
  predictions=append(predictions,got)
  consequences=append(consequences,event.Consequence)
  if e=os.WriteFile(filepath.Join(root,rootedReadSignal),[]byte("x"),0600);e!=nil{t.Fatal(e)}
  if _,e=reopened.Observe(context.Background());e==nil{t.Fatal("invalid sensor observation was learned")}
  if reopened.sequence!=1{t.Fatal("rejected observation consumed capacity")}
  other,e:=NewRootedReadLearner(t.TempDir());if e!=nil{t.Fatal(e)}
  if e=other.Restore(receipt);e==nil{t.Fatal("receipt replayed across workspaces")}
 }
 if consequences[0]!=0||consequences[1]!=1||!(predictions[0]<.5&&predictions[1]>.5){
  t.Fatal("source-independent two-world outcome failed",predictions,consequences)
 }
 t.Log("NCI04G7A_ROOTED_NATIVE_READ_PASS: distinct external outcomes, source-stable descriptors and restart; transfer/resource gate NOT QUALIFIED")
}
func TestNCI04G7AFailClosedOnInvalidToolAndBudget(t *testing.T) {
 root:=t.TempDir()
 learner,e:=NewRootedReadLearner(root);if e!=nil{t.Fatal(e)}
 if _,e=learner.Observe(context.Background());e==nil{t.Fatal("missing source-independent state admitted")}
 file:=filepath.Join(root,rootedReadSignal)
 if e=os.WriteFile(file,[]byte("10"),0600);e!=nil{t.Fatal(e)}
 if _,e=learner.Observe(context.Background());e==nil{t.Fatal("oversized input admitted")}
 if e=os.WriteFile(file,[]byte("1"),0600);e!=nil{t.Fatal(e)}
 for i:=0;i<rootedReadBudget;i++ {
  if _,e=learner.Observe(context.Background());e!=nil{t.Fatal(e)}
 }
 if _,e=learner.Observe(context.Background());e==nil{t.Fatal("extra observation admitted")}
 receipt,_:=learner.Snapshot()
 if e=learner.Restore([]byte("{invalid"));e==nil{t.Fatal("invalid receipt accepted")}
 now,_:=learner.Snapshot();if string(now)!=string(receipt){t.Fatal("failed restore changed learned state")}
 if _,e=NewRootedReadLearner("relative");e==nil{t.Fatal("relative workspace accepted")}
}
