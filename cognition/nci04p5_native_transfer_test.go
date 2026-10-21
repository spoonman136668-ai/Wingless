package cognition

// NCI-04P5 is a one-time native cross-family development falsifier.
// It freezes the genuinely trained P4 source-cognition state and asks whether
// that state contains a transferable relation for the existing (unchanged)
// toolboundary security workflow. This is NOT scientific first-use: this
// repository's source is visible, and the action semantics have no learned
// equivalence. The result must not masquerade as unfamiliar-task success.
import (
 "bytes"
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "math"
 "os"
 "os/exec"
 "path/filepath"
 "runtime"
 "strings"
 "testing"
 "time"
 "unsafe"
)

func TestNCI04P5FrozenCausalStateNativeCrossFamilyFalsifier(t *testing.T) {
 root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 checkout:=filepath.Join(t.TempDir(),"source")
 clone:=exec.Command("git","clone","--quiet","--no-hardlinks","--no-checkout",root,checkout)
 if out,e:=clone.CombinedOutput();e!=nil{t.Fatalf("isolated checkout unavailable: %v %s",e,out)}
 rev:=exec.Command("git","-C",root,"rev-parse","HEAD")
 originalSHA,err:=rev.Output();if err!=nil{t.Fatal(err)}
 sha:=strings.TrimSpace(string(originalSHA))
 detached:=exec.Command("git","-C",checkout,"checkout","--quiet","--detach",sha)
 if out,e:=detached.CombinedOutput();e!=nil{t.Fatalf("source identity checkout failed: %v %s",e,out)}
 build:=exec.Command("go","run","./cmd/ice","build")
 build.Dir=checkout
 build.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 if out,e:=build.CombinedOutput();e!=nil{t.Fatalf("scratch canonical ICE unavailable: %v %s",e,out)}
 learner,err:=OpenSourceReasoner(checkout);if err!=nil{t.Fatal(err)}
 train,err:=learner.ObserveNativeICEIntervention(context.Background())
 if err!=nil{t.Fatal(err)}
 if train.NativeOutcomes!=[3]uint8{1,0,1}||learner.State.Observations!=3{
  t.Fatalf("did not train on authenticated native P4 source: %#v",train)
 }
 frozen,err:=learner.State.Snapshot();if err!=nil{t.Fatal(err)}
 // The actual held-out workflow is a different native package and concerns
 // rooted file access, symlink containment, limits, and cancellation.
 task:=nativeReasoningTest{Name:"TestBoundary",Package:"./toolboundary",Source:"toolboundary/tools_test.go"}
 for _,prior:=range nativeReasoningTasks{
  if prior.Name==task.Name||prior.Source==task.Source{t.Fatal("candidate cross-family task appeared in training allowlist")}
 }
 observedSymbol,err:=learner.symbol(task.Name);if err!=nil{t.Fatal(err)}
 if observedSymbol.File!=task.Source||observedSymbol.Kind!="function"{t.Fatal("held-out task is not original source")}
 destination,err:=learner.symbol("Execute");if err!=nil{t.Fatal(err)}
 if destination.File!="toolboundary/tools.go"{t.Fatalf("untrusted held-out native endpoint: %#v",destination)}
 originalFile:=filepath.Join(checkout,task.Source)
 raw,err:=os.ReadFile(originalFile);if err!=nil{t.Fatal(err)}
 committed:=exec.Command("git","-C",checkout,"show","HEAD:"+task.Source)
 immutable,err:=committed.Output()
 if err!=nil||!bytes.Equal(raw,immutable){t.Fatal("held-out native source provenance mismatch")}
 sourceDigest:=sha256.Sum256(raw)

 // There is NO legitimate learned mapping between P4's "disable ICE
 // freshness guard" action and toolboundary's file-access policy actions.
 // Query action zero only as a frozen-state diagnostic, never relabel it
 // as a semantic transfer prediction. A novel relation must be uncertain.
 prior,err:=learner.State.Predict(vectorFromSource(observedSymbol),vectorFromSource(destination),0)
 if err!=nil{t.Fatal(err)}
 if math.Abs(prior-0.5)>1e-12{
  t.Fatalf("unsupported confident cross-domain relation: %g",prior)
 }
 var before,after runtime.MemStats
 runtime.ReadMemStats(&before)
 deadline,cancel:=context.WithTimeout(context.Background(),45*time.Second)
 defer cancel()
 testCmd:=exec.CommandContext(deadline,"go","test","-json",task.Package,"-run","^TestBoundary$","-count=1")
 testCmd.Dir=checkout
 testCmd.Env=append(os.Environ(),"GOFLAGS=","GOWORK=off","GOPROXY=off","GOSUMDB=off","GOTOOLCHAIN=local")
 started:=time.Now()
 output,runErr:=testCmd.CombinedOutput()
 wall:=time.Since(started).Nanoseconds()
 if deadline.Err()!=nil{t.Fatal("held-out native task timed out")}
 nativeOutcome,err:=nativeGoTestVerdict(output,task,runErr)
 if err!=nil{t.Fatalf("held-out native consequence not authenticated: %v",err)}
 if nativeOutcome!=1{t.Fatal("previously passing native security workflow regressed")}
 runtime.ReadMemStats(&after)
 surviving,err:=learner.State.Snapshot()
 if err!=nil||!bytes.Equal(surviving,frozen){t.Fatal("held-out evaluation changed frozen learned state")}
 restart:=filepath.Join(t.TempDir(),"research-state.json")
 if err:=learner.Save(restart);err!=nil{t.Fatal(err)}
 cold,err:=OpenSourceReasoner(checkout);if err!=nil{t.Fatal(err)}
 if err:=cold.Load(restart);err!=nil{t.Fatal(err)}
 unseen,err:=cold.State.Predict(vectorFromSource(observedSymbol),vectorFromSource(destination),0)
 if err!=nil||unseen!=prior{t.Fatal("native cross-family abstention changed after cold restore")}
 audit:=map[string]any{
  "gate":"NCI-04P5",
  "source_commit":sha,
  "training_native_outcomes":train.NativeOutcomes,
  "training_actions":[]int{0,1,0},
  "heldout_native_package":task.Package,
  "heldout_native_test":task.Name,
  "heldout_source_sha256":hex.EncodeToString(sourceDigest[:]),
  "heldout_native_outcome":nativeOutcome,
  "frozen_candidate_probability_without_action_mapping":prior,
  "source_held_count_control_prior":0.5,
  "trained_ICE_action_count_control":0.75,
  "count_comparator_transfer_applicable":false,
  "unfamiliar_action_semantics_equivalence":"NOT_ESTABLISHED",
  "source_family": "separate original toolboundary package; not independently sequestered scientific source",
  "heldout_test_child_wall_ns":wall,
  "go_test_parent_heap_alloc_delta_bytes":int64(after.TotalAlloc)-int64(before.TotalAlloc),
  "candidate_native_struct_bytes":unsafe.Sizeof(learner.State),
  "candidate_restart_snapshot_bytes":len(frozen),
  "all_in_matched_cpu_rss_and_io":"NOT_QUALIFIED",
  "native_state_changed_during_heldout":false,
  "external_truth_oracle":"original native Go security test assertion only",
  "independent_scientific_first_use":false,
  "scientific_transfer_result":"NOT_DEMONSTRATED",
  "disposition":"STOP_P4_FIXTURE_TUNING; REQUIRE_LEARNED_CROSS_TASK_OBSERVATION_AND_ACTION_INVARIANTS",
 }
 encoded,e:=json.Marshal(audit);if e!=nil{t.Fatal(e)}
 t.Logf("NCI04P5_NATIVE_TRANSFER_FALSIFIER=%s",encoded)
}
