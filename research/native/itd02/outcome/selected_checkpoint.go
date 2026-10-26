package outcome

import (
 "bytes"
 "crypto/sha256"
 "encoding/json"
 "errors"
 "io"

 r8 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r8"
)

const SelectedSchema = "wingless.itd02.selected-r8-outcome.v1"

type selectedState struct {
 Native []byte `json:"native_r8"`
 Rewards []byte `json:"outcome_ledger"`
}
type selectedEnvelope struct {
 Schema string `json:"schema"`
 State selectedState `json:"state"`
 Digest [32]byte `json:"digest"`
}

// SnapshotSelected binds both independently validated native checkpoints.
func(s *Selected) SnapshotSelected()([]byte,error){
 if s==nil||s.Policy==nil||s.Rewards==nil{return nil,errors.New("OUTCOME_SELECTED_INVALID")}
 native,err:=s.Policy.Snapshot();if err!=nil{return nil,err}
 rewards,err:=s.Rewards.Snapshot();if err!=nil{return nil,err}
 state:=selectedState{Native:native,Rewards:rewards}
 body,err:=json.Marshal(state);if err!=nil{return nil,err}
 return json.Marshal(selectedEnvelope{Schema:SelectedSchema,State:state,Digest:sha256.Sum256(body)})
}

// RestoreSelected refuses partial, noncanonical, corrupted, or mismatched
// checkpoints. Neither constituent can be independently substituted.
func RestoreSelected(raw []byte)(*Selected,error){
 if len(raw)==0||len(raw)>4<<20{return nil,errors.New("OUTCOME_SELECTED_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw))
 dec.DisallowUnknownFields()
 var envelope selectedEnvelope
 if err:=dec.Decode(&envelope);err!=nil{return nil,err}
 if err:=dec.Decode(new(struct{}));err!=io.EOF{return nil,errors.New("OUTCOME_SELECTED_TRAILING_DATA")}
 canonical,err:=json.Marshal(envelope)
 if err!=nil||!bytes.Equal(raw,canonical)||envelope.Schema!=SelectedSchema{return nil,errors.New("OUTCOME_SELECTED_NONCANONICAL")}
 body,err:=json.Marshal(envelope.State);if err!=nil{return nil,err}
 if envelope.Digest!=sha256.Sum256(body){return nil,errors.New("OUTCOME_SELECTED_DIGEST")}
 native,err:=r8.Restore(envelope.State.Native);if err!=nil{return nil,err}
 rewards,err:=Restore(envelope.State.Rewards);if err!=nil{return nil,err}
 if rewards.NextID!=native.LiveSeen+1{return nil,errors.New("OUTCOME_SELECTED_EVENT_MISMATCH")}
 if rewards.Closed&&native.Armed{return nil,errors.New("OUTCOME_SELECTED_TERMINAL_PENDING")}
 return &Selected{Policy:native,Rewards:rewards},nil
}
