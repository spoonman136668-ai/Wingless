package mg37

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "math"
 parent "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg28"
)

const MaxRewardContexts = 4096

// SourceAgent isolates selected-action reward evidence by original-public context.
// This is an opt-in draft adapter over the unchanged MG28 native learner.
type SourceAgent struct {
 Parent *parent.Agent
 Rewards map[string][7]Evidence
 Pending int
 PendingContext string
 HasPending bool
}
func NewSourceAgent(enabled bool)*SourceAgent {
 return &SourceAgent{Parent:parent.New(enabled),Rewards:make(map[string][7]Evidence)}
}
func publicContext(p parent.Public)(string,error) {
 if len(p.Image)!=7 || p.Direction<0 || p.Direction>3 || len(p.Mission)==0 {return "",errors.New("MG43_BAD_PUBLIC")}
 for _,row:=range p.Image {if len(row)!=7{return "",errors.New("MG43_BAD_IMAGE")};for _,pixel:=range row {if len(pixel)!=3{return "",errors.New("MG43_BAD_PIXEL")}}}
 front:=p.Image[3][5]
 value:=struct { Mission string `json:"mission"`; Direction int `json:"direction"`; Front [3]int `json:"front"` }{p.Mission,p.Direction,[3]int{front[0],front[1],front[2]}}
 raw,err:=json.Marshal(value);if err!=nil{return "",err}
 sum:=sha256.Sum256(raw);return hex.EncodeToString(sum[:]),nil
}
func(a *SourceAgent) BeginEpisode()error{
 if a==nil||a.Parent==nil||a.Rewards==nil{return errors.New("MG43_BAD_PARENT")}
 if err:=a.Parent.BeginEpisode();err!=nil{return err}
 a.HasPending=false;a.PendingContext=""
 return nil
}
func(a *SourceAgent) Choose(before parent.Public)(*int,string,error){
 if a==nil||a.Parent==nil||a.Rewards==nil||a.HasPending{return nil,"",errors.New("MG43_PENDING_OR_BAD_PARENT")}
 context,err:=publicContext(before);if err!=nil{return nil,"",err}
 proposed,status,err:=a.Parent.Choose(before)
 if err!=nil||proposed==nil{return nil,"",errors.New("MG43_PARENT_CHOOSE_FAILED")}
 d,err:=Choose(*proposed,a.Rewards[context]);if err!=nil{return nil,"",err}
 result:=d.Action
 if d.Override{return &result,"MG43_CONTEXT_REWARD_PROGRESS_OVERRIDE",nil}
 return proposed,status,nil
}
func(a *SourceAgent) Commit(before parent.Public,selected int)error{
 if a==nil||a.Parent==nil||a.Rewards==nil||a.HasPending||selected<0||selected>6{return errors.New("MG43_BAD_COMMIT")}
 context,err:=publicContext(before);if err!=nil{return err}
 if err=a.Parent.Commit(before,selected);err!=nil{return err}
 a.Pending=selected;a.PendingContext=context;a.HasPending=true
 return nil
}
func(a *SourceAgent) Feedback(before parent.Public,selected int,after parent.Public,reward float64,terminal bool)error{
 if a==nil||a.Parent==nil||a.Rewards==nil||!a.HasPending||a.Pending!=selected{return errors.New("MG43_UNSELECTED_FEEDBACK")}
 context,err:=publicContext(before);if err!=nil||context!=a.PendingContext{return errors.New("MG43_CONTEXT_MISMATCH")}
 if math.IsNaN(reward)||math.IsInf(reward,0)||reward< -1||reward>1{return errors.New("MG43_INVALID_REWARD")}
 if _,ok:=a.Rewards[context];!ok && len(a.Rewards)>=MaxRewardContexts{return errors.New("MG43_REWARD_CONTEXT_CAP")}
 if err=a.Parent.Feedback(before,selected,after,reward,terminal);err!=nil{return err}
 row:=a.Rewards[context];row[selected].Count++;row[selected].RewardSum+=reward;a.Rewards[context]=row
 a.HasPending=false;a.PendingContext=""
 return nil
}
