package cognition

// P5: source-held DEVELOPMENT falsifier on another existing native task
// family. Train ONLY on project ICE Query consequences, cold-restart, then
// predict before observing original toolboundary authorization outcomes.
// Neither holdout result nor damaged-source test text enters LatentState.
// This does not consume protected scientific first-use or qualify transfer.
import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "math"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "testing"
 "time"
)

func runNativeToolBoundaryTest(ctx context.Context,checkout string)(uint8,int64,[]byte,error){
 task:=nativeReasoningTest{Name:"TestBoundary",Package:"./toolboundary",Source:"toolboundary/tools_test.go"}
 deadline,cancel:=context.WithTimeout(ctx,35*time.Second);defer cancel()
 cmd:=exec.CommandContext(deadline,"go","test","-json","./toolboundary","-run","^TestBoundary$","-count=1")
 cmd.Dir=checkout
 cmd.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 start:=time.Now()
 output,runErr:=cmd.CombinedOutput()
 elapsed:=time.Since(start).Nanoseconds()
 if deadline.Err()!=nil{return 0,elapsed,nil,deadline.Err()}
 outcome,err:=nativeGoTestVerdict(output,task,runErr)
 if err!=nil{return 0,elapsed,nil,fmt.Errorf("native holdout execution was not a completed test assertion: %w",err)}
 return uint8(outcome),elapsed,output,nil
}

func p5Brier(predictions [3]float64,outcomes [3]uint8)float64{
 var sum float64
 for i,p:=range predictions{
  delta:=p-float64(outcomes[i])
  sum+=delta*delta
 }
 return sum/3
}
func TestNCI04P5TransferToOriginalToolAuthorization(t *testing.T){
 // Native source family and action are frozen independently of any learned
 // state: the ORIGINAL ICE stale-index test is training; existing original
 // toolboundary TestBoundary is evaluation. No new target test or answer.
 repo,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 detached:=filepath.Join(t.TempDir(),"original-native")
 cloned:=exec.Command("git","clone","--quiet","--no-hardlinks","--no-checkout",repo,detached)
 if out,e:=cloned.CombinedOutput();e!=nil{t.Fatalf("native project checkout unavailable: %v %s",e,out)}
 h:=exec.Command("git","-C",repo,"rev-parse","HEAD")
 head,e:=h.Output();if e!=nil{t.Fatal(e)}
 sha:=strings.TrimSpace(string(head))
 if len(sha)!=40{t.Fatal("native Git source identity absent")}
 detachedHead:=exec.Command("git","-C",detached,"checkout","--quiet","--detach",sha)
 if out,e:=detachedHead.CombinedOutput();e!=nil{t.Fatalf("cannot establish frozen holdout source: %v %s",e,out)}
 index:=exec.Command("go","run","./cmd/ice","build")
 index.Dir=detached
 index.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 if out,e:=index.CombinedOutput();e!=nil{t.Fatalf("ICE source index unavailable: %v %s",e,out)}
 trainer,e:=OpenSourceReasoner(detached);if e!=nil{t.Fatal(e)}
 receipt,e:=trainer.ObserveNativeICEIntervention(context.Background())
 if e!=nil{t.Fatal(e)}
 if receipt.NativeOutcomes!=[3]uint8{1,0,1}||trainer.State.Observations!=3{
  t.Fatalf("training results invalid: %#v",receipt)
 }
 frozen,e:=trainer.State.Snapshot();if e!=nil{t.Fatal(e)}
 outside:=filepath.Join(t.TempDir(),"native-state.json")
 if e=trainer.Save(outside);e!=nil{t.Fatal(e)}
 cold,e:=OpenSourceReasoner(detached);if e!=nil{t.Fatal(e)}
 if e=cold.Load(outside);e!=nil{t.Fatal(e)}
 after,e:=cold.State.Snapshot()
 if e!=nil||!bytes.Equal(frozen,after){t.Fatal("trained state not durable across cold reopen",e)}

 // Crucial: Native query features are source derived from the ORIGINAL
 // different toolboundary package. Do not reuse the training task identity,
 // expected answer, target-function label or future test result.
 source,e:=cold.symbol("TestBoundary");if e!=nil{t.Fatal(e)}
 target,e:=cold.symbol("Execute");if e!=nil{t.Fatal(e)}
 if source.File!="toolboundary/tools_test.go"||target.File!="toolboundary/tools.go"{
  t.Fatal("native holdout source identity drift")
 }
 featFrom,featTo:=vectorFromSource(source),vectorFromSource(target)
 p0,e:=cold.State.Predict(featFrom,featTo,0);if e!=nil{t.Fatal(e)}
 p1,e:=cold.State.Predict(featFrom,featTo,1);if e!=nil{t.Fatal(e)}
 priorPredictions:=[3]float64{p0,p1,p0}
 // Fully fitted conventional action-conditional count sees the same P4
 // three observations, with unit Beta smoothing. It receives no holdout
 // target or additional training data.
 countPredictions:=[3]float64{.75,1.0/3,.75}
 if p0<0||p0>1||p1<0||p1>1{t.Fatal("invalid heldout probability")}

 // Observe *after* predictions are recorded, in a separate real native
 // authorization package, with exactly one reversible source intervention.
 originalPath:=filepath.Join(detached,"toolboundary","tools.go")
 original,e:=os.ReadFile(originalPath);if e!=nil{t.Fatal(e)}
 trusted:=exec.Command("git","-C",detached,"show","HEAD:toolboundary/tools.go")
 blob,e:=trusted.Output();if e!=nil||!bytes.Equal(blob,original){t.Fatal("native holdout source differs from Git")}
 const policyGuard = `if a.Tool != "read_file" || !has(p.AllowedTools, a.Tool) || !has(p.AllowedPaths, a.Path) {`
 const weakenedGuard = `if a.Tool != "read_file" && !has(p.AllowedTools, a.Tool) && !has(p.AllowedPaths, a.Path) {`
 if bytes.Count(original,[]byte(policyGuard))!=1{t.Fatal("native authorization guard identity drift")}
 var actual [3]uint8
 var wall [3]int64
 actual[0],wall[0],_,e=runNativeToolBoundaryTest(context.Background(),detached)
 if e!=nil||actual[0]!=1{t.Fatalf("toolboundary original native task failed: %d %v",actual[0],e)}
 changed:=bytes.Replace(original,[]byte(policyGuard),[]byte(weakenedGuard),1)
 if bytes.Equal(changed,original){t.Fatal("native authorization intervention absent")}
 if e=os.WriteFile(originalPath,changed,0600);e!=nil{t.Fatal(e)}
 scope:=exec.Command("git","-C",detached,"diff","--name-only")
 files,e:=scope.Output()
 if e!=nil||strings.TrimSpace(string(files))!="toolboundary/tools.go"{t.Fatalf("unbounded authorization intervention: %q %v",files,e)}
 var negative []byte
 actual[1],wall[1],negative,e=runNativeToolBoundaryTest(context.Background(),detached)
 if e!=nil{t.Fatalf("holdout failure not a completed native assertion: %v",e)}
 if actual[1]!=0||!bytes.Contains(negative,[]byte("shell"))||!bytes.Contains(negative,[]byte("ok.txt")){
  t.Fatalf("holdout did not fail specifically on unauthorized tool admission: %d %s",actual[1],negative)
 }
 if e=os.WriteFile(originalPath,original,0600);e!=nil{t.Fatal(e)}
 actual[2],wall[2],_,e=runNativeToolBoundaryTest(context.Background(),detached)
 if e!=nil||actual[2]!=1{t.Fatalf("holdout native source reversal failed: %d %v",actual[2],e)}
 restored,e:=os.ReadFile(originalPath)
 if e!=nil||!bytes.Equal(restored,original){t.Fatal("original holdout source not restored")}
 status:=exec.Command("git","-C",detached,"status","--porcelain")
 // Canonical ICE emitted locally by cmd/ice build should match source-bound
 // tracked artifacts except source-index metadata from the disposable head.
 dirty,e:=status.Output();if e!=nil{t.Fatal(e)}
 if bytes.Contains(dirty,[]byte("toolboundary/tools.go")){t.Fatalf("holdout source remained dirty: %q",dirty)}
 if _,e=cold.State.Snapshot();e!=nil{t.Fatal(e)}
 if cold.State.Observations!=3{t.Fatal("held-out outcomes leaked into learner")}
 candidateBrier:=p5Brier(priorPredictions,actual)
 countBrier:=p5Brier(countPredictions,actual)
 if math.IsNaN(candidateBrier)||math.IsNaN(countBrier){t.Fatal("invalid held-out score")}
 result:=map[string]any{
  "gate":"NCI-04P5-W",
  "train_family":"ice stale-source query validation",
  "heldout_family":"toolboundary read-only authorization enforcement",
  "test_source_commit":sha,
  "training_observations":receipt.NativeOutcomes,
  "native_holdout_outcomes":actual,
  "candidate_prior_heldout":priorPredictions,
  "equal_information_trained_count_prior":countPredictions,
  "candidate_heldout_brier":candidateBrier,
  "trained_count_heldout_brier":countBrier,
  "candidate_better_than_count":candidateBrier<countBrier,
  "native_unseen_family_observed_after_prediction":true,
  "holdout_not_used_in_training":cold.State.Observations==3,
  "native_holdout_process_wall_ns":wall,
  "source_inspection_control":"NOT_INDEPENDENTLY_QUALIFIED",
  "full_paired_cpu_rss_snapshot_io":"NOT_QUALIFIED",
  "independent_scientific_first_use":false,
  "classification":"NATIVE_SECOND_FAMILY_TRANSFER_FALSIFIER__SCIENTIFIC_NO_GO",
 }
 b,e:=json.Marshal(result);if e!=nil{t.Fatal(e)}
 t.Logf("NCI04P5_NATIVE_HELDOUT_TRANSFER=%s",b)
}
