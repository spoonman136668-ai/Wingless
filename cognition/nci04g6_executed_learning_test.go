package cognition

import (
 "context"
 "encoding/json"
 "fmt"
 "math"
 "path/filepath"
 "testing"
 "time"
 "unsafe"
)

// Strong direct conditional-count control: it uses every action/consequence
// observation available to the candidate, with no additional information.
type exactExecutedControl struct{
 Count [2][2]uint8
}
func (c *exactExecutedControl) Learn(a,b uint8){c.Count[a][b]++}
func (c *exactExecutedControl) Forecast(a uint8)float64 {
 n0,n1:=c.Count[a][0],c.Count[a][1]
 return (float64(n1)+1)/(float64(n0)+float64(n1)+2)
}

func TestNCI04G6SourceBoundLearningAgainstTrainedCountControl(t *testing.T){
 p,err:=NewNativeProbe(filepath.Join(t.TempDir(),"private"))
 if err!=nil{t.Fatal(err)}
 s:=&ExecutedLearningSession{Probe:p}
 c:=&exactExecutedControl{}
 if pred,err:=s.PredictExecutable(true);err!=nil||pred!=.5{t.Fatalf("prior was not unknown: %v %v",pred,err)}
 var nativeDur,timeDur time.Duration
 for _,a:=range []bool{false,true,true}{
  start:=time.Now()
  observed,err:=s.ActAndLearn(context.Background(),a)
  nativeDur+=time.Since(start)
  if err!=nil{t.Fatal(err)}
  start=time.Now()
  c.Learn(observed.Action,observed.Consequence)
  timeDur+=time.Since(start)
 }
 after,err:=s.PredictExecutable(true)
 if err!=nil||! (after>.5&&after<1){t.Fatalf("no learned causal update: %v %v",after,err)}
 strong:=c.Forecast(1)
 if strong<=after{t.Fatal("expected strong count control not trained with equal information")}
 snap,err:=s.Snapshot();if err!=nil{t.Fatal(err)}
 reopened,err:=OpenNativeProbe(p.root);if err!=nil{t.Fatal(err)}
 replay:=&ExecutedLearningSession{Probe:reopened}
 if err:=replay.Restore(snap);err!=nil{t.Fatal(err)}
 check,err:=replay.PredictExecutable(true)
 if err!=nil||math.Abs(after-check)>1e-12{t.Fatal("learned state lost over restart",err)}
 if _,err:=replay.ActAndLearn(context.Background(),false);err==nil{t.Fatal("restart reset action counter")}
 baseline:=unsafe.Sizeof(LatentState{})
 report:=map[string]any{
  "gate":"NCI-04G6-W","engineering_native_path":"PASS",
  "exact_source_execution":true,"source_custody":"generated hard-pinned fixture, not general source tasks",
  "learning_updates":s.State.Observations,"native_interventions":p.sequence,
  "learned_probability":after,"trained_conventional_probability":strong,
  "candidate_resident_bytes":s.ResidentBytes(),"native_baseline_resident_bytes":baseline,
  "total_native_execution_plus_learning_wall_ns":nativeDur.Nanoseconds(),
  "control_incremental_training_wall_ns":timeDur.Nanoseconds(),
  "restart_snapshot_bytes":len(snap),
  "scientific_first_use":false,
  "gate_disposition":"NOT_QUALIFIED: conventional control outperforms learned model and all-in model resident size exceeds baseline; full RSS/CPU equivalence and unseen native task families unavailable"}
 data,_:=json.Marshal(report)
 fmt.Println("NCI04G6_WINGLESS_COMPARISON="+string(data))
 if s.ResidentBytes()<=baseline{t.Fatal("resource gate changed; reassess budget claim")}
}
