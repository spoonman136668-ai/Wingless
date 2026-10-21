package cognition

// G7A development seam: actual source-independent workspace observation.
// Rooted file reading is an information-gathering action, NOT proven causal
// transfer, external sensor attestation or scientific first-use.
import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "path/filepath"
 "time"
 "github.com/spoonman136668-ai/Wingless/toolboundary"
)

const rootedReadSignal = "native-state.bit"
const rootedReadBudget = 3

type RootedReadEvent struct {
 Action string
 Consequence uint8
 Sequence uint8
 EvidenceBasis string
}
type RootedReadLearner struct {
 root string
 state LatentState
 sequence uint8
}
type rootedReadReceipt struct {
 Schema string
 WorkspaceSHA256 string
 Sequence uint8
 Learned json.RawMessage
}
func NewRootedReadLearner(root string)(*RootedReadLearner,error){
 if !filepath.IsAbs(root)||filepath.Clean(root)!=root{return nil,errors.New("absolute canonical workspace required")}
 return &RootedReadLearner{root:root},nil
}
func (s *RootedReadLearner) descriptors()([FeatureWidth]float32,[FeatureWidth]float32){
 // Fixed descriptors do not contain file contents, expected outcomes or task IDs.
 return [FeatureWidth]float32{-.75,.25,-.25,.5},[FeatureWidth]float32{.5,-.5,.25,-.75}
}
func (s *RootedReadLearner) Predict()(float64,error){
 if s==nil{return 0,errors.New("nil native learner")}
 src,dst:=s.descriptors()
 return s.state.Predict(src,dst,0)
}
// Observe reads the actual authorized workspace, never a caller-supplied label.
func (s *RootedReadLearner) Observe(ctx context.Context)(RootedReadEvent,error){
 var zero RootedReadEvent
 if s==nil||s.root==""{return zero,errors.New("nil native learner")}
 if s.sequence>=rootedReadBudget{return zero,errors.New("native observation budget exhausted")}
 policy:=toolboundary.Policy{Workspace:s.root,AllowedTools:[]string{"read_file"},
  AllowedPaths:[]string{rootedReadSignal},Timeout:5*time.Second,MaxOutputBytes:1}
 observed,err:=toolboundary.Execute(ctx,policy,toolboundary.Proposal{Tool:"read_file",Path:rootedReadSignal})
 if err!=nil{return zero,err}
 if observed.ExitCode!=0||len(observed.Stdout)!=1||(observed.Stdout!="0"&&observed.Stdout!="1"){
  return zero,errors.New("invalid source-independent native observation")
 }
 outcome:=int(observed.Stdout[0]-'0')
 src,dst:=s.descriptors()
 next:=s.state
 if err:=next.Observe(src,dst,0,outcome);err!=nil{return zero,err}
 s.state=next
 s.sequence++
 return RootedReadEvent{Action:"rooted-read:native-state.bit",Consequence:uint8(outcome),
  Sequence:s.sequence,EvidenceBasis:"authorized rooted read of workspace-owned data; external attestation absent"},nil
}
func (s *RootedReadLearner) Snapshot()([]byte,error){
 if s==nil{return nil,errors.New("nil native learner")}
 state,err:=s.state.Snapshot();if err!=nil{return nil,err}
 digest:=sha256.Sum256([]byte(s.root))
 return json.Marshal(rootedReadReceipt{"wingless.nci04g7a.rooted-read.v1",hex.EncodeToString(digest[:]),s.sequence,state})
}
func (s *RootedReadLearner) Restore(data []byte)error{
 if s==nil||len(data)>65536{return errors.New("invalid native restart receipt")}
 var wire rootedReadReceipt
 if err:=json.Unmarshal(data,&wire);err!=nil{return err}
 digest:=sha256.Sum256([]byte(s.root))
 if wire.Schema!="wingless.nci04g7a.rooted-read.v1"||wire.WorkspaceSHA256!=hex.EncodeToString(digest[:])||wire.Sequence>rootedReadBudget{
  return errors.New("native workspace identity or observation budget mismatch")
 }
 var state LatentState
 if err:=state.Restore(wire.Learned);err!=nil{return err}
 if state.Observations!=uint32(wire.Sequence){return errors.New("native observation count drift")}
 s.sequence=wire.Sequence;s.state=state
 return nil
}
