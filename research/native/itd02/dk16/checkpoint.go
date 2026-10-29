package dk16
import (
 "bytes"
 "encoding/json"
 "errors"
 "io"
 prior "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk14"
)
const MaxCheckpoint=3000000
type state struct{
 Schema string `json:"schema"`
 Parent json.RawMessage `json:"parent"`
 States map[string]Actions `json:"states"`
}
type envelope struct{State state `json:"state"`;SHA256 string `json:"sha256"`}
func (m *Policy) Snapshot()([]byte,error){
 if !m.valid()||m.Policy.Pending!=nil{return nil,errors.New("DK16_PENDING_OR_DRIFT")}
 source,e:=m.Policy.Snapshot();if e!=nil{return nil,e}
 record:=state{Schema:Schema,Parent:source,States:m.States}
 data,e:=json.Marshal(record);if e!=nil{return nil,e}
 raw,e:=json.Marshal(envelope{State:record,SHA256:hash(data)})
 if e!=nil||len(raw)>MaxCheckpoint{return nil,errors.New("DK16_CHECKPOINT_BUDGET")}
 return raw,nil
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>MaxCheckpoint{return nil,errors.New("DK16_CHECKPOINT_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var wire envelope
 if e:=d.Decode(&wire);e!=nil{return nil,e}
 if e:=d.Decode(new(any));e!=io.EOF{return nil,errors.New("DK16_TRAILING_SOURCE")}
 data,e:=json.Marshal(wire.State)
 if e!=nil||wire.State.Schema!=Schema||hash(data)!=wire.SHA256{
  return nil,errors.New("DK16_CHECKPOINT_SOURCE_HASH")
 }
 parent,e:=prior.Restore(wire.State.Parent);if e!=nil{return nil,e}
 result:=&Policy{Policy:parent,States:wire.State.States}
 if !result.valid(){return nil,errors.New("DK16_CHECKPOINT_COUNTS")}
 same,e:=result.Snapshot()
 if e!=nil||!bytes.Equal(same,raw){return nil,errors.New("DK16_NONCANONICAL")}
 return result,nil
}
