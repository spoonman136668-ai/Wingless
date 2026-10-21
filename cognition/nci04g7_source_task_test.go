package cognition

import (
 "context"
 "encoding/json"
 "fmt"
 "math"
 "path/filepath"
 "runtime"
 "testing"
 "time"
 "unsafe"
)

// Exact same 3 native task receipts for every model. The conventional model
// is fully fitted to all received outcomes; the learned model has no advantage.
type nativeOutcomeCount struct{ Outcomes [3][2]uint8 }
func (c *nativeOutcomeCount) Observe(action int, outcome uint8){c.Outcomes[action][outcome]++}
func (c *nativeOutcomeCount) Predict(action int)float64 {
 n:=c.Outcomes[action]
 return (float64(n[1])+1)/(float64(n[0]+n[1])+2)
}

func TestNCI04G7RealSourceTaskAndCostAccounting(t *testing.T){
 root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 r,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 var count nativeOutcomeCount
 var sourceWall,modelWall,controlWall time.Duration
 for i:=0;i<3;i++ {
  before,err:=r.PredictNativeReasoningTask(i);if err!=nil||before!=.5{t.Fatalf("untrained action must be unknown: %d %.6f %v",i,before,err)}
  began:=time.Now()
  observation,err:=r.ObserveNativeReasoningTask(context.Background(),i)
  sourceWall+=time.Since(began)
  if err!=nil{t.Fatal(err)}
  if observation.ObservedPass!=1||observation.Action!=uint8(i)||
     observation.SourceCommit==""||observation.IndexSHA!=r.IndexSHA||
     observation.SourceFileSHA256==""{t.Fatalf("not authenticated by native task: %+v",observation)}
  started:=time.Now()
  count.Observe(i,observation.ObservedPass)
  controlWall+=time.Since(started)
 }
 if _,err:=r.ObserveNativeReasoningTask(context.Background(),0);err==nil{
  t.Fatal("source task exhausted observation capacity")
 }
 began:=time.Now()
 p,err:=r.PredictNativeReasoningTask(0)
 modelWall+=time.Since(began)
 if err!=nil||math.IsNaN(p)||p<=.5{t.Fatal("native observation did not change learned prediction",err)}
 started:=time.Now()
 strong:=count.Predict(0)
 controlWall+=time.Since(started)
 file:=filepath.Join(t.TempDir(),"native-state.json")
 if err:=r.Save(file);err!=nil{t.Fatal(err)}
 restored,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 if err:=restored.Load(file);err!=nil{t.Fatal(err)}
 again,err:=restored.PredictNativeReasoningTask(0)
 if err!=nil||p!=again{t.Fatal("real native task learning failed restart equivalence",err)}
 state,_:=r.State.Snapshot()
 var mem runtime.MemStats
 runtime.ReadMemStats(&mem)
 report:=map[string]any{
  "gate":"NCI-04G7-W", "native_source_tasks_executed":3,
  "source_authentication":"GIT_BLOB_AND_CURRENT_ICE_AND_NATIVE_GO_TEST",
  "candidate_native_state_bytes":unsafe.Sizeof(r.State),
  "existing_native_latent_state_bytes":unsafe.Sizeof(LatentState{}),
  "existing_source_reasoner_struct_bytes":unsafe.Sizeof(*r),
  "trained_count_control_struct_bytes":unsafe.Sizeof(count),
  "native_state_snapshot_bytes":len(state),
  "runtime_process_heap_alloc_bytes":mem.Alloc,
  "common_source_provenance_and_execution_wall_ns":sourceWall.Nanoseconds(),
  "candidate_incremental_prediction_wall_ns":modelWall.Nanoseconds(),
  "control_incremental_update_and_prediction_wall_ns":controlWall.Nanoseconds(),
  "candidate_fitted_prediction":p,"fully_trained_count_control_prediction":strong,
  "all_in_matched_cpu_rss":"NOT_QUALIFIED: separate-process equal-budget controls, snapshot and execution custody accounting still missing",
  "learned_transfer":"NOT_PROVEN: three positive allowlisted authentic tests, no unseen native tasks or causal state changes",
  "scientific_first_use":false,
  "disposition":"SOURCE_NATIVE_TASK_INTEGRATION_ONLY; SCIENTIFIC_NO_GO"}
 b,_:=json.Marshal(report)
 fmt.Println("NCI04G7_WINGLESS_NATIVE_REPORT="+string(b))
 if unsafe.Sizeof(r.State)!=unsafe.Sizeof(LatentState{}){t.Fatal("new resident learned state exceeded original exact LatentState")}
 if strong<p{t.Fatal("unearned candidate superiority claim")}
}

func TestNCI04G7InvalidActionRejectsWithoutExecuting(t *testing.T){
 root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
 r,err:=OpenSourceReasoner(root);if err!=nil{t.Fatal(err)}
 if _,err:=r.ObserveNativeReasoningTask(context.Background(),3);err==nil{t.Fatal("unlisted task was admitted")}
 if _,err:=r.ObserveNativeReasoningTask(context.Background(),-1);err==nil{t.Fatal("negative task action was admitted")}
 if r.State.Observations!=0{t.Fatal("invalid action changed learned source state")}
}
