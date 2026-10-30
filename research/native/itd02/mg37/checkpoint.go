package mg37

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "math"
 "strings"
 parent "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg28"
)
const CheckpointSchema = "wingless.itd02.mg50.context-reward-source.v1"
const MaxRewardCheckpoint = 9000000

type rewardState struct {
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 Rewards map[string][7]Evidence `json:"rewards"`
}
type rewardEnvelope struct {
 State rewardState `json:"state"`
 SHA string `json:"sha256"`
}
func rewardHash(raw []byte)string{h:=sha256.Sum256(raw);return hex.EncodeToString(h[:])}
func(a *SourceAgent) rewardValid()bool {
 if a==nil||a.Parent==nil||a.Rewards==nil||len(a.Rewards)>MaxRewardContexts||a.HasPending{return false}
 for key,row:=range a.Rewards{
  if len(key)!=64||strings.Trim(key,"0123456789abcdef")!=""{return false}
  for _,cell:=range row{
   if cell.Count>parent.MaxSamples||math.IsNaN(cell.RewardSum)||math.IsInf(cell.RewardSum,0)||math.Abs(cell.RewardSum)>float64(cell.Count){return false}
  }
 }
 return true
}
// Snapshot refuses in-flight actions; parent state is serialized unchanged.
func(a *SourceAgent) Snapshot()([]byte,error){
 if !a.rewardValid(){return nil,errors.New("MG50_INVALID_SOURCE_STATE")}
 parentBytes,err:=a.Parent.Snapshot();if err!=nil{return nil,err}
 s:=rewardState{CheckpointSchema,parentBytes,a.Rewards}
 body,err:=json.Marshal(s);if err!=nil{return nil,err}
 wire,err:=json.Marshal(rewardEnvelope{s,rewardHash(body)})
 if err!=nil||len(wire)>MaxRewardCheckpoint{return nil,errors.New("MG50_CHECKPOINT_CAP")}
 return wire,nil
}
func RestoreSourceAgent(raw []byte)(*SourceAgent,error){
 if len(raw)==0||len(raw)>MaxRewardCheckpoint{return nil,errors.New("MG50_BAD_SIZE")}
 decoder:=json.NewDecoder(bytes.NewReader(raw))
 decoder.DisallowUnknownFields()
 var wire rewardEnvelope
 if err:=decoder.Decode(&wire);err!=nil{return nil,err}
 if decoder.Decode(new(any))!=io.EOF{return nil,errors.New("MG50_TRAILING_CHECKPOINT")}
 data,err:=json.Marshal(wire.State)
 if err!=nil||wire.State.Schema!=CheckpointSchema||wire.SHA!=rewardHash(data){return nil,errors.New("MG50_BAD_DIGEST")}
 p,err:=parent.Restore(wire.State.Parent);if err!=nil{return nil,err}
 a:=&SourceAgent{Parent:p,Rewards:wire.State.Rewards}
 if !a.rewardValid(){return nil,errors.New("MG50_UNTRUSTED_REWARD_CELLS")}
 rebuilt,err:=a.Snapshot()
 if err!=nil||!bytes.Equal(raw,rebuilt){return nil,errors.New("MG50_NONCANONICAL")}
 return a,nil
}
