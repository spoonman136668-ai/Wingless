// Package r4 is an opt-in source-only transient frontier experiment on top
// of unchanged R3. Source outcomes remain real selected-action public feedback.
// Episode routing state is distinct from native learning and cleared on reset.
package r4

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r3 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r3"
)

const Schema="wingless.itd02.r4.episode-public-frontier.v1"
const MaxEpisodeSteps=4096
const MaxEpisodeCells=2048
const MaxEpisodePoses=4096
const FrontierDepth=36
const FrontierNodes=2048
const StallTrigger=4

type PoseEvidence struct{
 Blocked bool `json:"blocked"`
 Seen uint16 `json:"seen"`
}
type Episode struct{
 Steps uint32 `json:"steps"`
 Stationary uint16 `json:"stationary"`
 Visits map[uint32]uint16 `json:"visits"`
 Poses map[uint32]PoseEvidence `json:"poses"`
 Pending bool `json:"pending"`
 Before r1.PublicGrid `json:"pending_before"`
 Action int `json:"pending_action"`
}
type Policy struct{
 Core *r3.Policy
 Route Episode
}
func New()*Policy{return &Policy{Core:r3.New(),Route:newEpisode()}}
func newEpisode()Episode{return Episode{
 Visits:make(map[uint32]uint16),
 Poses:make(map[uint32]PoseEvidence),
}}
func cellKey(s r1.PublicGrid)uint32{return uint32(s.X*65+s.Y)}
func poseKey(s r1.PublicGrid)uint32{return cellKey(s)*4+uint32(s.Direction)}
func legal(s r1.PublicGrid)bool{
 return s.X>=0&&s.X<=64&&s.Y>=0&&s.Y<=64&&
 s.GoalX>=0&&s.GoalX<=64&&s.GoalY>=0&&s.GoalY<=64&&
 s.Direction>=0&&s.Direction<=3
}
func valid(p *Policy)bool{
 if p==nil||p.Core==nil||p.Route.Steps>MaxEpisodeSteps||
 len(p.Route.Visits)>MaxEpisodeCells||len(p.Route.Poses)>MaxEpisodePoses||
 p.Route.Stationary>uint16(MaxEpisodeSteps){return false}
 var sumCells uint64
 for key,n:=range p.Route.Visits{
  if key>=65*65||n==0{return false}
  sumCells+=uint64(n)
 }
 var sumPoses uint64
 for key,n:=range p.Route.Poses{
  if key>=65*65*4||n.Seen==0{return false}
  sumPoses+=uint64(n.Seen)
 }
 if sumCells!=uint64(p.Route.Steps)||sumPoses!=uint64(p.Route.Steps){
  return false
 }
 if p.Route.Pending&&(!legal(p.Route.Before)||
  p.Route.Action<0||p.Route.Action>2){return false}
 if !p.Route.Pending&&p.Route.Action!=0{return false}
 return true
}
// BeginEpisode resets source-local search history, never model weights.
func(p *Policy)BeginEpisode()error{
 if !valid(p)||p.Route.Pending{return errors.New("R4_INVALID_EPISODE_RESET")}
 p.Route=newEpisode()
 return nil
}
func seenBlocked(e Episode,s r1.PublicGrid)(bool,bool){
 v,ok:=e.Poses[poseKey(s)]
 return v.Blocked,ok
}
func(p *Policy)predictEpisode(e Episode,s r1.PublicGrid,action int)(r1.PublicGrid,bool,error){
 if b,ok:=seenBlocked(e,s);ok{s.FrontBlocked=b}
 nxt,ok,err:=p.Core.Predict(s,action)
 if err!=nil||!ok{return r1.PublicGrid{},false,err}
 if b,known:=seenBlocked(e,nxt);known{nxt.FrontBlocked=b}else{
  // Unknown future walls are a provisional hypothesis, never remembered
  // as observed clearance; actual feedback can contradict the source model.
  nxt.FrontBlocked=false
 }
 return nxt,true,nil
}
type node struct{state r1.PublicGrid;first,depth int}
func manhattan(s r1.PublicGrid)int{
 x:=s.X-s.GoalX;if x<0{x= -x}
 y:=s.Y-s.GoalY;if y<0{y= -y}
 return x+y
}
// frontier searches for a state with less source-local visitation, never
// using a privileged environment map or learned heldout consequences.
func(p *Policy)frontier(start r1.PublicGrid)(int,bool,error){
 q:=[]node{{state:start,first:-1}}
 reached:=map[uint32]bool{poseKey(start):true}
 bestFirst,bestScore:=0,1<<30
 for head:=0;head<len(q)&&head<FrontierNodes;head++{
  n:=q[head]
  if n.depth>=FrontierDepth{continue}
  for a:=0;a<3;a++{
   next,ok,err:=p.predictEpisode(p.Route,n.state,a)
   if err!=nil{return 0,false,err}
   if !ok{continue}
   k:=poseKey(next)
   if reached[k]{continue}
   reached[k]=true
   first:=n.first;if first<0{first=a}
   // Do not reward phantom progress from a stationary action.
   if next.X!=start.X||next.Y!=start.Y{
    v:=int(p.Route.Visits[cellKey(next)])
    // Novel cell gets priority over a cheap revisited detour. Ties
    // prefer proximity to public goal, then shorter learned route.
    score:=v*1000+(n.depth+1)*7+manhattan(next)
    if score<bestScore{bestScore=score;bestFirst=first}
   }
   if len(q)<FrontierNodes{
    q=append(q,node{state:next,first:first,depth:n.depth+1})
   }
  }
 }
 return bestFirst,bestScore<1<<30,nil
}
// Choose has no training effects. A second call before actual source
// Feedback is illegal, protecting one action => one observed consequence.
func(p *Policy)Choose(s r1.PublicGrid)(int,error){
 if !valid(p)||p.Route.Pending||!legal(s){
  return 0,errors.New("R4_INVALID_PUBLIC_ACTION_REQUEST")
 }
 chosen,err:=p.Core.Choose(s)
 if err!=nil{return 0,err}
 if p.Route.Stationary>=StallTrigger||
  p.Route.Visits[cellKey(s)]>=StallTrigger{
  first,found,e:=p.frontier(s)
  if e!=nil{return 0,e}
  if found{chosen=first}
 }
 if chosen<0||chosen>2{return 0,errors.New("R4_UNLAWFUL_CHOICE")}
 p.Route.Pending=true
 p.Route.Before=s
 p.Route.Action=chosen
 return chosen,nil
}
// Feedback is an OBSERVATION of a fresh rollout, not a model-training call.
// It must follow the exact chosen action and actual before/after state.
func(p *Policy)Feedback(before r1.PublicGrid,action int,after r1.PublicGrid)error{
 if !valid(p)||!p.Route.Pending||!legal(before)||!legal(after)||
  before!=p.Route.Before||action!=p.Route.Action||
  p.Route.Steps>=MaxEpisodeSteps{return errors.New("R4_INVALID_OBSERVED_CONSEQUENCE")}
 dx,dy:=after.X-before.X,after.Y-before.Y
 if dx < -1||dx>1||dy < -1||dy>1{
  return errors.New("R4_NONLOCAL_SOURCE_CONSEQUENCE")
 }
 pos:=cellKey(before);pose:=poseKey(before)
 cells:=p.Route.Visits[pos]
 observed:=p.Route.Poses[pose]
 if cells>=uint16(MaxEpisodeSteps)||observed.Seen>=uint16(MaxEpisodeSteps)||
  (cells==0&&len(p.Route.Visits)>=MaxEpisodeCells)||
  (observed.Seen==0&&len(p.Route.Poses)>=MaxEpisodePoses){
  return errors.New("R4_SESSION_BUDGET_EXHAUSTED")
 }
 stall:=p.Route.Stationary
 if before.X==after.X&&before.Y==after.Y{
  if stall>=uint16(MaxEpisodeSteps){return errors.New("R4_STALL_COUNTER_EXHAUSTED")}
  stall++
 }else{stall=0}
 if p.Route.Visits==nil{p.Route.Visits=make(map[uint32]uint16)}
 if p.Route.Poses==nil{p.Route.Poses=make(map[uint32]PoseEvidence)}
 p.Route.Visits[pos]=cells+1
 p.Route.Poses[pose]=PoseEvidence{Blocked:before.FrontBlocked,Seen:observed.Seen+1}
 p.Route.Steps++
 p.Route.Stationary=stall
 p.Route.Pending=false
 p.Route.Before=r1.PublicGrid{}
 p.Route.Action=0
 return nil
}
// ObserveFit is exclusively authorized for a genuine training stream, and
// not called by Feedback or any CHECK/DEV decision.
func(p *Policy)ObserveFit(before r1.PublicGrid,action int,after r1.PublicGrid,reward float64,terminal bool)error{
 if !valid(p)||p.Route.Steps!=0||p.Route.Pending{
  return errors.New("R4_FIT_INTERLEAVED_WITH_PUBLIC_EVALUATION")
 }
 return p.Core.Observe(before,action,after,reward,terminal)
}
type snapshotData struct{
 R3 []byte `json:"r3_original_checkpoint"`
 Episode Episode `json:"episode_session"`
}
type envelope struct{
 Schema string `json:"schema"`
 Data snapshotData `json:"data"`
 SHA256 string `json:"sha256"`
}
func(p *Policy)Snapshot()([]byte,error){
 if !valid(p){return nil,errors.New("R4_INVALID_SOURCE_MODEL")}
 raw,e:=p.Core.Snapshot();if e!=nil{return nil,e}
 data:=snapshotData{R3:raw,Episode:p.Route}
 body,e:=json.Marshal(data);if e!=nil{return nil,e}
 hash:=sha256.Sum256(body)
 return json.Marshal(envelope{Schema:Schema,Data:data,SHA256:hex.EncodeToString(hash[:])})
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>2<<20{return nil,errors.New("R4_CHECKPOINT_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw));d.DisallowUnknownFields()
 var enc envelope
 if e:=d.Decode(&enc);e!=nil{return nil,e}
 if e:=d.Decode(new(struct{}));e!=io.EOF{return nil,errors.New("R4_TRAILING_CHECKPOINT_BYTES")}
 if enc.Schema!=Schema{return nil,errors.New("R4_CHECKPOINT_SCHEMA")}
 body,e:=json.Marshal(enc.Data);if e!=nil{return nil,e}
 h:=sha256.Sum256(body)
 if enc.SHA256!=hex.EncodeToString(h[:]){return nil,errors.New("R4_CHECKPOINT_HASH")}
 core,e:=r3.Restore(enc.Data.R3);if e!=nil{return nil,e}
 ep:=enc.Data.Episode
 if ep.Visits==nil{ep.Visits=make(map[uint32]uint16)}
 if ep.Poses==nil{ep.Poses=make(map[uint32]PoseEvidence)}
 p:=&Policy{Core:core,Route:ep}
 if !valid(p){return nil,errors.New("R4_INVALID_RESTORED_EPISODE")}
 return p,nil
}
