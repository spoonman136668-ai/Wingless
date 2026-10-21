package cognition

// NCI04F: source-free falsification of source-identity/response-hypothesis
// ambiguity. The production NCI04E candidate implementation is UNMODIFIED.
import (
  "encoding/json"
  "fmt"
  "os"
  "reflect"
  "testing"
  "unsafe"
)

type nci04fProfile [2]int

type nci04fResponseSignature struct {
  Used bool
  Outcomes nci04fProfile
  Seen uint32
}
type nci04fSymbolicControl struct {
  Modes [2]nci04fResponseSignature
}
func (m *nci04fSymbolicControl) TrainObservedTrial(events [][2]int) {
  var values nci04fProfile
  var seen [2]bool
  for _,event:=range events{
    a,y:=event[0],event[1]
    if a<0||a>1||y<0||y>1{panic("bad fixture")}
    values[a]=y;seen[a]=true
  }
  if !seen[0]||!seen[1]{panic("incomplete observed trial")}
  for i,mode:=range m.Modes{
    if mode.Used && mode.Outcomes==values {m.Modes[i].Seen++;return}
  }
  for i,mode:=range m.Modes{
    if !mode.Used {m.Modes[i]=nci04fResponseSignature{Used:true,Outcomes:values,Seen:1};return}
  }
  panic("frozen control capacity exceeded")
}
func(m *nci04fSymbolicControl) PredictAfterAction0(y int)(int,bool){
  predictions:=map[int]bool{}
  for _,mode:=range m.Modes {
    if mode.Used&&mode.Outcomes[0]==y {predictions[mode.Outcomes[1]]=true}
  }
  if len(predictions)!=1{return 0,false}
  for prediction:=range predictions{return prediction,true}
  return 0,false
}
func(m *nci04fSymbolicControl) PredictAfterBoth(y0,y1 int)(int,bool){
  for _,mode:=range m.Modes{
    if mode.Used&&mode.Outcomes[0]==y0&&mode.Outcomes[1]==y1{return mode.Outcomes[1],true}
  }
  return 0,false
}
func TestNCI04FIndistinguishableFirstProbe(t *testing.T) {
  profiles:=[2]nci04fProfile{{0,0},{0,1}}
  candidate:=NewActionResponseIdentity()
  var conventional nci04fSymbolicControl
  for session:=0;session<12;session++{
    latent:=session%2
    if e:=candidate.BeginTrial();e!=nil{t.Fatal(e)}
    eventStream:=make([][2]int,0,8)
    for step:=0;step<8;step++{
      a:=step%2;outcome:=profiles[latent][a]
      if e:=candidate.Observe(a,outcome);e!=nil{t.Fatal(e)}
      eventStream=append(eventStream,[2]int{a,outcome})
    }
    conventional.TrainObservedTrial(eventStream)
  }
  if !conventional.Modes[0].Used||!conventional.Modes[1].Used{
    t.Fatal("fair trained signature control missed full observed response profiles")
  }
  candidateModeCount:=candidate.ModeCount()
  invalidKnown:=0
  correctlyAmbiguousControl:=0
  controlAfterDiscriminatingProbe:=0
  candidateAfterProbeCorrect:=0
  candidateAfterProbeKnown:=0
  for session:=0;session<8;session++{
    latent:=session%2
    if e:=candidate.BeginTrial();e!=nil{t.Fatal(e)}
    if _,known,e:=candidate.Predict(1);e!=nil||known{t.Fatal("source identity leaked before first action")}
    // First action 0 has SAME result under both physical sources.
    if e:=candidate.Observe(0,profiles[latent][0]);e!=nil{t.Fatal(e)}
    if _,known,e:=candidate.Predict(1);e!=nil{t.Fatal(e)}else if known {invalidKnown++}
    if _,ok:=conventional.PredictAfterAction0(profiles[latent][0]);!ok{correctlyAmbiguousControl++}
    // The second action exposes the previously hidden response difference.
    if e:=candidate.Observe(1,profiles[latent][1]);e!=nil{t.Fatal(e)}
    if p,known,e:=candidate.Predict(1);e!=nil{t.Fatal(e)}else if known{
      candidateAfterProbeKnown++
      if (p>=.5)==(profiles[latent][1]==1){candidateAfterProbeCorrect++}
    }
    if v,ok:=conventional.PredictAfterBoth(profiles[latent][0],profiles[latent][1]);ok&&v==profiles[latent][1]{
      controlAfterDiscriminatingProbe++
    }
  }
  // Defensible comparator does not claim identity before the discriminating
  // probe, and learns both response families without externally supplied IDs.
  if correctlyAmbiguousControl!=8||controlAfterDiscriminatingProbe!=8{
    t.Fatalf("conventional trained control failed %d %d",correctlyAmbiguousControl,controlAfterDiscriminatingProbe)
  }
  negativeModel:=NewActionResponseIdentity()
  for session:=0;session<8;session++{
    if e:=negativeModel.BeginTrial();e!=nil{t.Fatal(e)}
    for step:=0;step<8;step++{
      a:=step%2
      if e:=negativeModel.Observe(a,profiles[0][a]);e!=nil{t.Fatal(e)}
    }
  }
  if negativeModel.ModeCount()!=1{t.Fatal("negative source wrongly split identical profiles")}
  disposition:="NEGATIVE"
  if candidateModeCount==2&&invalidKnown==0&&candidateAfterProbeCorrect==8{disposition="PASS_NECESSARY_CONDITION_ONLY"}
  report:=map[string]any{
    "schema":"nci04f.wingless.sourcefree-v1",
    "result":disposition,
    "observed_training_sessions":12,
    "heldout_sessions":8,
    "latent_identity_exposed_to_candidate":false,
    "source_profiles_share_action0":true,
    "candidate_learned_modes":candidateModeCount,
    "candidate_known_action1_after_ambiguous_action0":invalidKnown,
    "candidate_known_action1_after_discriminating_probe":candidateAfterProbeKnown,
    "candidate_correct_after_discriminating_probe":candidateAfterProbeCorrect,
    "unlabelled_strong_control_distinct_signatures":2,
    "strong_control_abstention_after_ambiguous_probe":correctlyAmbiguousControl,
    "strong_control_correct_after_discriminating_probe":controlAfterDiscriminatingProbe,
    "identical_response_negative_control_modes":negativeModel.ModeCount(),
    "candidate_resident_bytes":unsafe.Sizeof(*candidate),
    "control_resident_bytes":unsafe.Sizeof(conventional),
    "scientific_first_use":0,
    "learned_entity_identity_accepted":false,
  }
  repeated:=report["result"]
  if !reflect.DeepEqual(repeated,disposition){t.Fatal("nondeterministic disposition")}
  raw,e:=json.Marshal(report);if e!=nil{t.Fatal(e)}
  fmt.Println("NCI04F_W_RESULT="+string(raw))
  if path:=os.Getenv("NCI04F_RESULT_PATH");path!=""{if e:=os.WriteFile(path,append(raw,'\n'),0600);e!=nil{t.Fatal(e)}}
}
