package mg37

import (
 "errors"
 "math"
 parent "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg28"
)

// SourceAgent adds an opt-in, non-invasive reward guard to the existing MG28
// selected-action feedback flow. It never synthesizes environment observations.
type SourceAgent struct {
 Parent *parent.Agent
 Rewards [7]Evidence
 Pending int
 HasPending bool
}
func NewSourceAgent(enabled bool)*SourceAgent{
 return &SourceAgent{Parent:parent.New(enabled)}
}
func(a *SourceAgent) BeginEpisode()error{
 if a==nil||a.Parent==nil{return errors.New("MG38_BAD_PARENT")}
 if err:=a.Parent.BeginEpisode();err!=nil{return err}
 a.HasPending=false
 return nil
}
func(a *SourceAgent) Choose(before parent.Public)(*int,string,error){
 if a==nil||a.Parent==nil||a.HasPending{return nil,"",errors.New("MG38_PENDING_OR_BAD_PARENT")}
 proposed,status,err:=a.Parent.Choose(before)
 if err!=nil||proposed==nil{return nil,"",errors.New("MG38_PARENT_CHOOSE_FAILED")}
 d,err:=Choose(*proposed,a.Rewards)
 if err!=nil{return nil,"",err}
 result:=d.Action
 if d.Override{return &result,"MG38_SELECTED_REWARD_PROGRESS_OVERRIDE",nil}
 return proposed,status,nil
}
func(a *SourceAgent) Commit(before parent.Public,selected int)error{
 if a==nil||a.Parent==nil||a.HasPending||selected<0||selected>6{return errors.New("MG38_BAD_COMMIT")}
 if err:=a.Parent.Commit(before,selected);err!=nil{return err}
 a.Pending=selected;a.HasPending=true
 return nil
}
func(a *SourceAgent) Feedback(before parent.Public,selected int,after parent.Public,reward float64,terminal bool)error{
 if a==nil||a.Parent==nil||!a.HasPending||a.Pending!=selected{return errors.New("MG38_UNSELECTED_FEEDBACK")}
 if math.IsNaN(reward)||math.IsInf(reward,0)||reward< -1||reward>1{return errors.New("MG38_INVALID_REWARD")}
 if err:=a.Parent.Feedback(before,selected,after,reward,terminal);err!=nil{return err}
 a.Rewards[selected].Count++
 a.Rewards[selected].RewardSum+=reward
 a.HasPending=false
 return nil
}
