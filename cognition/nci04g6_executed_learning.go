package cognition

// NCI-04G6: connect actual executed consequences to the existing bounded
// LatentState, with no caller-supplied outcome or hidden state label.
// This is a fixed-generated-source *engineering* seam, not causal transfer.
import (
 "context"
 "encoding/json"
 "errors"
 "unsafe"
)

type ExecutedLearningSession struct {
 Probe *NativeProbe
 State LatentState
}

type nativeLearningWire struct {
 Schema string `json:"schema"`
 Probe json.RawMessage `json:"probe"`
 State json.RawMessage `json:"state"`
}

// Both embeddings come only from hard-pinned generated source/actuator
// descriptions, not from a user-supplied identity or evaluator target.
func (s *ExecutedLearningSession) descriptors(action bool)([FeatureWidth]float32,[FeatureWidth]float32) {
 moduleHash:=probeHash([]byte(probeModule))
 fixtureHash:=probeHash([]byte(probeFalse))
 src:=vectorFromDigest(moduleHash)
 dst:=vectorFromDigest(fixtureHash)
 if action {
  dst=vectorFromDigest(probeHash([]byte(probeTrue)))
 }
 return src,dst
}

func vectorFromDigest(hexSource string)[FeatureWidth]float32{
 var v [FeatureWidth]float32
 for i:=range v {
  // Hex ASCII characters are bounded in 0..127. This is a stable source
  // descriptor, not proof of an invariant across arbitrary source edits.
  v[i]=float32(int(hexSource[2*i]))/64.-1.
 }
 return v
}
func (s *ExecutedLearningSession) ActAndLearn(ctx context.Context,action bool)(NativeProbeObservation,error){
 if s==nil||s.Probe==nil{return NativeProbeObservation{},errors.New("nil native learning session")}
 observed,err:=s.Probe.Intervene(ctx,action)
 if err!=nil{return NativeProbeObservation{},err}
 src,dst:=s.descriptors(action)
 if err:=s.State.Observe(src,dst,int(observed.Action),int(observed.Consequence));err!=nil{
  return NativeProbeObservation{},err
 }
 return observed,nil
}
func (s *ExecutedLearningSession) PredictExecutable(action bool)(float64,error){
 if s==nil{return 0,errors.New("nil session")}
 src,dst:=s.descriptors(action)
 idx:=0;if action{idx=1}
 return s.State.Predict(src,dst,idx)
}
func (s *ExecutedLearningSession) Snapshot()([]byte,error){
 if s==nil||s.Probe==nil{return nil,errors.New("nil session")}
 p,err:=s.Probe.Snapshot();if err!=nil{return nil,err}
 latent,err:=s.State.Snapshot();if err!=nil{return nil,err}
 return json.Marshal(nativeLearningWire{"wingless.nci04g6.executed-learning.v1",p,latent})
}
func (s *ExecutedLearningSession) Restore(blob []byte)error{
 if s==nil||s.Probe==nil||len(blob)>65536{return errors.New("invalid session snapshot")}
 var wire nativeLearningWire
 if err:=json.Unmarshal(blob,&wire);err!=nil{return err}
 if wire.Schema!="wingless.nci04g6.executed-learning.v1"{return errors.New("invalid native session identity")}
 var candidate LatentState
 if err:=candidate.Restore(wire.State);err!=nil{return err}
 // Source, action counter and durable sequence are validated independently.
 if err:=s.Probe.Restore(wire.Probe);err!=nil{return err}
 s.State=candidate
 return nil
}
func (s *ExecutedLearningSession) ResidentBytes()uintptr{
 if s==nil{return 0}
 return unsafe.Sizeof(s.State)+unsafe.Sizeof(*s.Probe)
}
