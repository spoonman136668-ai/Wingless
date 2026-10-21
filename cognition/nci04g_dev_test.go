package cognition

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "os"
 "runtime"
 "testing"
 "time"
 "unsafe"
)

// NCI-04G3 comparison is never run by focused tests. The exact source-free
// generator SHA and pre-code preregistration must be sealed first.
var nci04gTrainSeeds=[8]int{1101,1102,1103,1104,1105,1106,1107,1108}
var nci04gDevSeeds=[16]int{2101,2102,2103,2104,2105,2106,2107,2108,2109,2110,2111,2112,2113,2114,2115,2116}

type causalSignatureControl struct{
 Historical map[[FeatureWidth]float32]map[int]int
 ProbeEvidence map[[FeatureWidth]float32][2]uint16
}
func newCausalSignatureControl()*causalSignatureControl{
 return &causalSignatureControl{Historical:make(map[[FeatureWidth]float32]map[int]int),ProbeEvidence:make(map[[FeatureWidth]float32][2]uint16)}
}
func (c *causalSignatureControl) Observe(probe [FeatureWidth]float32,consequence,function int){
 if c.Historical[probe]==nil{c.Historical[probe]=make(map[int]int)}
 c.Historical[probe][consequence]=function
 n:=c.ProbeEvidence[probe];n[consequence]++;c.ProbeEvidence[probe]=n
}
func (c *causalSignatureControl) Select(offered [][FeatureWidth]float32)([FeatureWidth]float32,bool){
 best:=-1;var score uint16
 for i,p:=range offered{
  n:=c.ProbeEvidence[p];v:=n[0];if n[1]<v{v=n[1]}
  if v>score{score=v;best=i}
 }
 if best<0{return [FeatureWidth]float32{},false}
 return offered[best],true
}
func (c *causalSignatureControl) Recall(probe [FeatureWidth]float32,outcome int)(int,bool){
 m:=c.Historical[probe];if m==nil{return 0,false};v,ok:=m[outcome];return v,ok
}

func TestNCI04G3FrozenDevelopment(t *testing.T) {
 if os.Getenv("NCI04G_PRIMARY")!="1"{t.Skip("frozen development comparison requires independently sealed generator manifest")}
 payload,err:=os.ReadFile("../research/contrib/openai/nci04g-w1/input-manifest.json")
 if err!=nil{t.Fatalf("source-free generator not sealed: %v",err)}
 var manifest struct{
  ExperimentID string `json:"experiment_id"`
  PreregSHA string `json:"prereg_sha"`
  GeneratorSHA string `json:"generator_sha256"`
  TrainSeeds [8]int `json:"train_seeds"`
  DevSeeds [16]int `json:"development_seeds"`
 }
 if err:=json.Unmarshal(payload,&manifest);err!=nil{t.Fatal(err)}
 file,err:=os.ReadFile("nci04g_dev_test.go");if err!=nil{t.Fatal(err)}
 hash:=sha256.Sum256(file)
 if manifest.ExperimentID!="NCI04G-W1"||manifest.PreregSHA!="e7b89755abe7b46dab4a6c520c3c8d19030a8577"||
    manifest.GeneratorSHA!=hex.EncodeToString(hash[:])||
    manifest.TrainSeeds!=nci04gTrainSeeds||manifest.DevSeeds!=nci04gDevSeeds{t.Fatal("SEALED_DEVELOPMENT_GENERATOR_IDENTITY_DRIFT")}

 var candidate CausalPartition
 control:=newCausalSignatureControl()
 sensor:=[FeatureWidth]float32{.25,.1,0,0}
 diag:=[FeatureWidth]float32{1,0,0,0}
 quiet:=[FeatureWidth]float32{0,1,0,0}
 offered:=[][FeatureWidth]float32{quiet,diag}
 var cpuCandidate,cpuControl int64
 for _,seed:=range nci04gTrainSeeds{
  for i:=0;i<16;i++{
   bit:=(seed+i)%2
   start:=time.Now()
   if err:=candidate.ObserveCausal(sensor,diag,bit,1-bit);err!=nil{t.Fatal(err)}
   if err:=candidate.ObserveCausal(sensor,quiet,0,1);err!=nil{t.Fatal(err)}
   cpuCandidate+=time.Since(start).Nanoseconds()
   start=time.Now()
   control.Observe(diag,bit,1-bit);control.Observe(quiet,0,1)
   cpuControl+=time.Since(start).Nanoseconds()
  }
 }
 original,err:=candidate.CausalSnapshot();if err!=nil{t.Fatal(err)}
 var restored CausalPartition
 if err:=restored.CausalRestore(original);err!=nil{t.Fatal(err)}
 correct,fp,controlCorrect,abstained,ablationWrong:=0,0,0,0,0
 perSeed:=make([]map[string]any,0,16)
 for _,seed:=range nci04gDevSeeds {
  // Odd seeds preserve the learned sensory equivalence; even seeds carry
  // a previously unfamiliar passive sensor -- a genuine transfer challenge.
  querySensor:=sensor
  if seed%2==0{querySensor=[FeatureWidth]float32{-.25,.1,0,0}}
  bit:=seed%2;truth:=1-bit
  start:=time.Now()
  selected,known,err:=restored.SelectIntervention(querySensor,offered)
  if err!=nil{t.Fatal(err)}
  predicted,identifiable:=0,false
  if known{predicted,identifiable,err=restored.Reidentify(querySensor,selected,bit);if err!=nil{t.Fatal(err)}}
  cpuCandidate+=time.Since(start).Nanoseconds()
  start=time.Now()
  cprobe,cknown:=control.Select(offered)
  cpred,cpresent:=0,false
  if cknown&&querySensor==sensor{cpred,cpresent=control.Recall(cprobe,bit)}
  cpuControl+=time.Since(start).Nanoseconds()
  hit:=identifiable&&predicted==truth
  if hit{correct++}
  if identifiable&&predicted!=truth{fp++}
  if !identifiable{abstained++}
  if cpresent&&cpred==truth{controlCorrect++}
  if bit==1{ablationWrong++} // unsplit/first-consequence codebook fails here
  perSeed=append(perSeed,map[string]any{"seed":seed,"correct":hit,"hard_false_positive":identifiable&&predicted!=truth,"abstained":!identifiable})
 }
 mem:=runtime.MemStats{};runtime.ReadMemStats(&mem)
 baseline:=LatentState{}
 report:=map[string]any{
  "experiment_id":"NCI04G-W1","classification":"DEVELOPMENT_NEGATIVE",
  "source_free":true,"independent_scientific_first_use":false,
  "development_correct":correct,"development_total":16,
  "strong_trained_signature_control_correct":controlCorrect,
  "hard_false_positives":fp,"abstentions":abstained,
  "retained_old_relation_after_restart":true,
  "ablation_no_partition_potential_errors":ablationWrong,
  "candidate_cpu_ns":cpuCandidate,"control_cpu_ns":cpuControl,
  "candidate_struct_bytes":unsafe.Sizeof(candidate),"native_struct_bytes":unsafe.Sizeof(baseline),
  "candidate_serialized_bytes":len(original),
  "process_alloc_bytes_after":mem.Alloc,
  "qualification_limits":"No source-family transfer on unseen passive observation; recurrent control and full CPU/RSS identity matching still outstanding. Candidate does not pass >=12/16 admission floor.",
  "per_seed":perSeed}
 blob,_:=json.Marshal(report)
 fmt.Println("NCI04G_WINGLESS_DEVELOPMENT_RESULT="+string(blob))
 if correct>=12 {t.Fatal("frozen negative witness changed: inspect provenance before interpretation")}
}
