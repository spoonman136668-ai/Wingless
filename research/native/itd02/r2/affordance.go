// Package r2 adds a FIT-learned action-consequence affordance to original
// Wingless R1. It does not contain a planner or an a-priori blocked-action rule.
package r2
import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "math"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)
const Schema="wingless.itd02.r2.public-forward-outcome-counts.v1"
type Policy struct {
 Core *r1.Policy
 BlockedAttempts uint64
 BlockedProgress uint64
}
func New()*Policy{return &Policy{Core:r1.New()}}
func(p *Policy)valid()bool{return p!=nil&&p.Core!=nil&&p.BlockedProgress<=p.BlockedAttempts&&p.BlockedAttempts<1<<52}
func(p *Policy)Scores(s r1.PublicGrid)([3]float64,error){
 if !p.valid(){return [3]float64{},errors.New("R2_INVALID_LEARNED_COUNTS")}
 scores,e:=p.Core.Scores(s);if e!=nil{return scores,e}
 if s.FrontBlocked&&p.BlockedAttempts>=8{
  post:=(float64(p.BlockedProgress)+1)/(float64(p.BlockedAttempts)+2)
  scores[2]-=8*(1-post)
 }
 if math.IsNaN(scores[2])||math.IsInf(scores[2],0){return [3]float64{},errors.New("R2_NONFINITE_SCORE")}
 return scores,nil
}
func(p *Policy)Choose(s r1.PublicGrid)(int,error){
 scores,e:=p.Scores(s);if e!=nil{return 0,e}
 best:=0;for a:=1;a<3;a++{if scores[a]>scores[best]{best=a}}
 return best,nil
}
func(p *Policy)Observe(before r1.PublicGrid,action int,after r1.PublicGrid,reward float64,terminal bool)error{
 if !p.valid(){return errors.New("R2_INVALID_STATE")}
 if p.BlockedAttempts>=1<<52-1{return errors.New("R2_EXHAUSTED_COUNTS")}
 // Call original authentic R1 update ONLY on a local cloned copy. Failed
 // observations must never change either the native R1 state or counts.
 next:=*p.Core
 if e:=next.Observe(before,action,after,reward,terminal);e!=nil{return e}
 attempts,progress:=p.BlockedAttempts,p.BlockedProgress
 if action==2&&before.FrontBlocked{
  attempts++
  if after.X!=before.X||after.Y!=before.Y {progress++}
 }
 p.Core=&next
 p.BlockedAttempts=attempts
 p.BlockedProgress=progress
 return nil
}
type payload struct{
 R1 []byte `json:"original_r1_checkpoint"`
 BlockedAttempts uint64 `json:"blocked_attempts"`
 BlockedProgress uint64 `json:"blocked_progress"`
}
type signed struct{
 Schema string `json:"schema"`
 Data payload `json:"data"`
 SHA256 string `json:"sha256"`
}
func(p *Policy)Snapshot()([]byte,error){
 if !p.valid(){return nil,errors.New("R2_INVALID_CHECKPOINT_SOURCE")}
 r,e:=p.Core.Snapshot();if e!=nil{return nil,e}
 b:=payload{r,p.BlockedAttempts,p.BlockedProgress}
 body,e:=json.Marshal(b);if e!=nil{return nil,e}
 digest:=sha256.Sum256(body)
 return json.Marshal(signed{Schema:Schema,Data:b,SHA256:hex.EncodeToString(digest[:])})
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>32768{return nil,errors.New("R2_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw))
 dec.DisallowUnknownFields()
 var env signed
 if e:=dec.Decode(&env);e!=nil{return nil,e}
 if e:=dec.Decode(new(struct{}));e!=io.EOF{return nil,errors.New("R2_CHECKPOINT_TRAILING_CONTENT")}
 if env.Schema!=Schema{return nil,errors.New("R2_CHECKPOINT_SCHEMA")}
 body,e:=json.Marshal(env.Data);if e!=nil{return nil,e}
 digest:=sha256.Sum256(body)
 if hex.EncodeToString(digest[:])!=env.SHA256{return nil,errors.New("R2_CHECKPOINT_AUTHENTICATION")}
 core,e:=r1.Restore(env.Data.R1);if e!=nil{return nil,e}
 p:=&Policy{Core:core,BlockedAttempts:env.Data.BlockedAttempts,BlockedProgress:env.Data.BlockedProgress}
 if !p.valid(){return nil,errors.New("R2_CHECKPOINT_COUNTS")}
 return p,nil
}
