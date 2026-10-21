package cognition

import (
 "bytes"
 "context"
 "encoding/json"
 "os"
 "os/exec"
 "path/filepath"
 "testing"
 "unsafe"
)

func TestNCI04P4ActualNativeObservationsConsumedAndColdRestored(t *testing.T){
 // Generate a canonical ICE view only in the detached source checkout.
 // Nothing touches original repo source, hidden labels or first-use custody.
 original,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 temp:=t.TempDir()
 cloned:=filepath.Join(temp,"source")
 clone:=exec.Command("git","clone","--quiet","--no-hardlinks","--no-checkout",original,cloned)
 if out,e:=clone.CombinedOutput();e!=nil{t.Fatalf("native clone: %v %s",e,out)}
 rev:=exec.Command("git","-C",original,"rev-parse","HEAD")
 sha,err:=rev.Output();if err!=nil{t.Fatal(err)}
 checkout:=exec.Command("git","-C",cloned,"checkout","--quiet","--detach",string(bytes.TrimSpace(sha)))
 if out,e:=checkout.CombinedOutput();e!=nil{t.Fatalf("source pin: %v %s",e,out)}
 build:=exec.Command("go","run","./cmd/ice","build")
 build.Dir=cloned
 build.Env=append(os.Environ(),"GOTOOLCHAIN=local","GOPROXY=off","GOSUMDB=off","GOWORK=off")
 if out,e:=build.CombinedOutput();e!=nil{t.Fatalf("scratch canonical ICE generation: %v %s",e,out)}
 r,err:=OpenSourceReasoner(cloned);if err!=nil{t.Fatal(err)}
 if unsafe.Sizeof(r.State)!=6472{t.Fatalf("native memory shape changed: %d",unsafe.Sizeof(r.State))}
 src,e:=r.symbol("TestSourceWins");if e!=nil{t.Fatal(e)}
 dest,e:=r.symbol("ObserveExistingCall");if e!=nil{t.Fatal(e)}
 q0,e:=r.State.Predict(vectorFromSource(src),vectorFromSource(dest),0)
 if e!=nil||q0!=.5{t.Fatal("untrained native prior was not uncertain",q0,e)}
 receipt,e:=r.ObserveNativeICEIntervention(context.Background())
 if e!=nil{t.Fatal(e)}
 if receipt.NativeOutcomes!=[3]uint8{1,0,1}||receipt.LearnedObservations!=3{
  t.Fatalf("native execution was not admitted transactionally: %#v",receipt)
 }
 if !(receipt.BaselinePosterior>receipt.BaselinePrior&&receipt.IntervenedPosterior<receipt.IntervenedPrior){
  t.Fatalf("learned native relation did not move with independently observed consequences: %#v",receipt)
 }
 if receipt.BaselinePrior!=.5||receipt.IntervenedPrior!=.5||r.State.Observations!=3{
  t.Fatal("candidate learned from an invented prior or incorrect budget",receipt)
 }
 for _,v:=range receipt.NativeWallNS{if v<=0{t.Fatal("missing real native execution timing")}}
 frozen,e:=r.State.Snapshot();if e!=nil{t.Fatal(e)}
 outside:=filepath.Join(temp,"native-learned-state.json")
 if err:=r.Save(outside);err!=nil{t.Fatal(err)}
 opened,e:=OpenSourceReasoner(cloned);if e!=nil{t.Fatal(e)}
 if e=opened.Load(outside);e!=nil{t.Fatal(e)}
 recovered,e:=opened.State.Snapshot();if e!=nil||!bytes.Equal(recovered,frozen){t.Fatal("original native state drifted over cold reopen",e)}
 if _,e=opened.ObserveNativeICEIntervention(context.Background());e==nil{
  t.Fatal("unbounded replay of the same intervention cycle accepted")
 }
 if x,_:=opened.State.Snapshot();!bytes.Equal(x,frozen){t.Fatal("invalid repeat mutated native memory")}
 // The native source test failures are used ONLY as native assertions,
 // while a fully fitted count/conditional model remains a strong control.
 // Three samples make no unfamiliar-task transfer or superiority claim.
 measured:=map[string]any{
  "gate":"NCI-04P4-W",
  "real_native_outcomes":receipt.NativeOutcomes,
  "candidate_prior":[2]float64{receipt.BaselinePrior,receipt.IntervenedPrior},
  "candidate_posterior":[2]float64{receipt.BaselinePosterior,receipt.IntervenedPosterior},
  "equal_information_count_posterior":[2]float64{.75,1.0/3},
  "candidate_state_bytes":unsafe.Sizeof(r.State),
  "native_restart_snap_bytes":len(frozen),
  "native_execution_wall_ns":receipt.NativeWallNS,
  "actual_source_commit":receipt.SourceCommit,
  "native_intervention_isolation":"detached isolated code edit, exact restoration",
  "new_trainable_parameters":0,
  "independent_unfamiliar_task_transfer":"NOT_TESTED",
  "resource_parity_all_in":"NOT_QUALIFIED",
  "science_first_use":false,
  "admission":"ENGINEERING_REAL_CAUSAL_LEARNING_PATH_PASS__SCIENCE_NO_GO",
 }
 b,e:=json.Marshal(measured);if e!=nil{t.Fatal(e)}
 t.Logf("NCI04P4_NATIVE_LEARNING_REPORT=%s",b)
}

func TestNCI04P4RejectsUnverifiedNativeLearningBeforeExecution(t *testing.T){
 r:=&SourceReasoner{}
 if _,err:=r.ObserveNativeICEIntervention(context.Background());err==nil{
  t.Fatal("nil/invalid source accepted for native causal learning")
 }
 r.State.Observations=1
 if _,err:=r.ObserveNativeICEIntervention(context.Background());err==nil{
  t.Fatal("non-authenticated live source identity accepted")
 }
}
