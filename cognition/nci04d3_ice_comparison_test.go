package cognition

import (
 "encoding/json"
 "fmt"
 "os"
 "path/filepath"
 "runtime"
 "testing"
 "time"
 "unsafe"
 "github.com/spoonman136668-ai/Wingless/ice"
)

// NCI-04D3: existing real ICE symbols, no artificial positive examples.
// This evaluates syntax-call evidence, NOT underlying runtime call truth.
func TestNCI04D3NativeICEAgainstExactSyntax(t *testing.T) {
 root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 if err=ice.Validate(root);err!=nil{t.Fatal("source index not validated:",err)}
 r,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 tasks:=[][][2]string{
  {{"Select","Entries"},{"Select","Check"}},
  {{"RunLocal","summarize"},{"runICESourceCognition","PredictCall"}},
 }
 before:=r.State.Observations
 start:=time.Now()
 for _,task:=range tasks {for _,pair:=range task {
  for n:=0;n<6;n++ {if err=r.ObserveExistingCall(pair[0],pair[1]);err!=nil {t.Fatalf("source observation %v: %v",pair,err)}}
 }}
 trainElapsed:=time.Since(start)
 if r.State.Observations-before!=24 {t.Fatalf("not cumulative learning: %d",r.State.Observations)}
 snapFile:=filepath.Join(t.TempDir(),"state.json")
 if err=r.Save(snapFile);err!=nil{t.Fatal(err)}
 reloaded,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 if err=reloaded.Load(snapFile);err!=nil{t.Fatal(err)}
 // Both earlier and later source tasks remain accessible after reload.
 probes:=[][2]string{
  {"Select","Entries"},{"Select","Check"},
  {"RunLocal","summarize"},{"runICESourceCognition","PredictCall"},
  {"Select","summarize"},{"RunLocal","Entries"},
 }
 score:=0.0;count:=0;exactCorrect:=0;learnedCorrect:=0
 queryStart:=time.Now()
 for _,pair:=range probes {
  p,err:=reloaded.PredictCall(pair[0],pair[1]);if err!=nil{t.Fatalf("source probe %v: %v",pair,err)}
  y:=0.0;if p.VerifiedSyntaxCall {y=1;exactCorrect++}
  delta:=p.LearnedProbability-y;score+=delta*delta;count++
  classification:=0.0;if p.LearnedProbability>=.5 {classification=1}
  if classification==y {learnedCorrect++}
 }
 queryElapsed:=time.Since(queryStart)
 if count!=6||score<=0{t.Fatalf("invalid comparison count %d risk %g",count,score)}
 var mem runtime.MemStats;runtime.ReadMemStats(&mem)
 snap,_:=os.ReadFile(snapFile)
 report:=map[string]interface{}{
 "gate":"NCI-04D3","basis":"source-indexed-syntax-witness-not-verified-runtime-semantics",
 "task_count":len(tasks),"observations":24,"probes":count,
 "exact_syntax_brier":0.0,"candidate_brier":score/float64(count),
 "exact_positive_witnesses":exactCorrect,"candidate_threshold_correct":learnedCorrect,
 "reloaded_snapshot":true,
 "nominal_learned_state_bytes":unsafe.Sizeof(LatentState{}),
 "serialized_persistent_snapshot_bytes":len(snap),
 "train_wall_ns":trainElapsed.Nanoseconds(),"query_wall_ns":queryElapsed.Nanoseconds(),
 "go_runtime_sys_bytes_at_probe":mem.Sys,
 "scientific_superiority":"NO_GO_exact_ground_truth_remains_authoritative",
 }
 encoded,err:=json.Marshal(report);if err!=nil{t.Fatal(err)}
 fmt.Printf("NCI04D3_WINGLESS_COMPARISON=%s\n",encoded)
}
