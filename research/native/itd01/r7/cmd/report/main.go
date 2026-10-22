// Native Wingless source-free scientific-development evidence receipt.
package main

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "os"
 "time"
 r7 "github.com/spoonman136668-ai/Wingless/research/native/itd01/r7"
)
type receipt struct {
 Schema string `json:"schema"`
 EvidenceClass string `json:"evidence_class"`
 ScientificFirstUse bool `json:"scientific_first_use"`
 OriginalHeldSourceUsed bool `json:"original_held_source_used"`
 Architecture string `json:"architecture"`
 Results []r7.Scores `json:"results"`
 PrimaryGate string `json:"primary_gate"`
 TrialCPUApproxWallMillis int64 `json:"trial_wall_millis_includes_generator"`
 Digest string `json:"results_sha256"`
}
func main(){
 start:=time.Now()
 results:=make([]r7.Scores,0,3)
 valid:=true
 for _,seed:=range []int{17,31,47}{
  s,e:=r7.Study(seed)
  if e!=nil{panic(e)}
  results=append(results,s)
  if s.Updates!=640||s.Changes==0||!s.ColdRestoreExact ||
  s.RetainedA<.75||s.B<.75||s.PreA-s.RetainedA>.05+1e-9 ||
  s.RetainedA-s.NoStateA<.15||s.B-s.NoStateB<.15 ||
  s.RetainedA-s.MajorityA<.15||s.B-s.MajorityB<.15||s.SerializedBytes>32768{valid=false}
 }
 raw,_:=json.Marshal(results);sha:=sha256.Sum256(raw)
 status:="SOURCE_FREE_NATIVE_LEARNING_PASS"
 if !valid{status="SOURCE_FREE_NATIVE_LEARNING_NO_GO"}
 out:=receipt{
  Schema:"itd01.r7.wingless.packed-native.v1",
  EvidenceClass:"SYNTHETIC_PUBLIC_FEEDBACK_PILOT_ONLY",
  ScientificFirstUse:false,OriginalHeldSourceUsed:false,
  Architecture:"WINGLESS_NATIVE_PACKED_SIGNED_CAUSAL_STATE_NOT_LLM_NOT_GRU",
  Results:results,PrimaryGate:status,
  TrialCPUApproxWallMillis:time.Since(start).Milliseconds(),
  Digest:hex.EncodeToString(sha[:]),
 }
 encoded,e:=json.MarshalIndent(out,"","  ");if e!=nil{panic(e)}
 fmt.Println(string(encoded))
 if !valid{os.Exit(2)}
}
