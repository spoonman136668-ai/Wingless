// Package r5 adds a cost-aware, goal-directed frontier candidate over the
// unmodified native Wingless R4 episode memory and fitted R3 transitions.
// Its route calculations do not observe a grid, privileged map, or reward.
package r5

import (
 "errors"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r4 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r4"
)

const Schema="wingless.itd02.r5.cost-aware-public-frontier.v1"
const MaxNodes=2048
const MaxDepth=36
const StallTrigger=4

// Frozen source-only objective weights, committed to HYPOTHESIS.json before
// candidate implementation. Never optimize these on consumed DEV worlds.
const (
 travelStepWeight=18
 goalDistanceWeight=6
 cellVisitWeight=32
 turnStepWeight=6
 unvisitedBonus=24
 blockedForwardWeight=40
)
type Policy struct{Core *r4.Policy}
func New()*Policy{return &Policy{Core:r4.New()}}
func legal(s r1.PublicGrid)bool{
 return s.X>=0&&s.X<=64&&s.Y>=0&&s.Y<=64&&
  s.GoalX>=0&&s.GoalX<=64&&s.GoalY>=0&&s.GoalY<=64&&
  s.Direction>=0&&s.Direction<4
}
func cellKey(s r1.PublicGrid)uint32{return uint32(s.X*65+s.Y)}
func poseKey(s r1.PublicGrid)uint32{return cellKey(s)*4+uint32(s.Direction)}
func md(s r1.PublicGrid)int{
 x:=s.X-s.GoalX;if x<0{x = -x}
 y:=s.Y-s.GoalY;if y<0{y = -y}
 return x+y
}
func(p *Policy)ObserveFit(s r1.PublicGrid,action int,after r1.PublicGrid,reward float64,done bool)error{
 if p==nil||p.Core==nil{return errors.New("R5_NIL_NATIVE_LEARNER")}
 return p.Core.ObserveFit(s,action,after,reward,done)
}
func(p *Policy)BeginEpisode()error{
 if p==nil||p.Core==nil{return errors.New("R5_NIL_EPISODE_RESET")}
 return p.Core.BeginEpisode()
}
func(p *Policy)Feedback(s r1.PublicGrid,action int,after r1.PublicGrid)error{
 if p==nil||p.Core==nil{return errors.New("R5_NIL_EPISODE_FEEDBACK")}
 return p.Core.Feedback(s,action,after)
}
type route struct{
 s r1.PublicGrid
 first int
 depth int
 travel int
}
// predict applies only past actual public pose evidence, never a source
// simulator rule. Unknown future public clearance remains a hypothesis.
func(p *Policy)predict(s r1.PublicGrid,action int)(r1.PublicGrid,bool,error){
 if cell,ok:=p.Core.Route.Poses[poseKey(s)];ok{s.FrontBlocked=cell.Blocked}
 out,ok,e:=p.Core.Core.Predict(s,action)
 if e!=nil||!ok{return r1.PublicGrid{},false,e}
 if cell,seen:=p.Core.Route.Poses[poseKey(out)];seen{
  out.FrontBlocked=cell.Blocked
 }else{out.FrontBlocked=false}
 return out,true,nil
}
// chooseFrontier performs bounded cheapest-model-route expansion, giving goal
// proximity and action time tangible weight alongside episode visits. It only
// returns a candidate action, not an unobserved future source consequence.
func(p *Policy)chooseFrontier(s r1.PublicGrid)(int,bool,error){
 q:=[]route{{s:s,first:-1,depth:0,travel:0}}
 bestRoute:=map[uint32]int{poseKey(s):0}
 bestFirst,bestScore:=0,1<<30
 for head:=0;head<len(q)&&head<MaxNodes;head++{
  n:=q[head]
  if n.depth>=MaxDepth{continue}
  for action:=0;action<3;action++{
   next,ok,e:=p.predict(n.s,action)
   if e!=nil{return 0,false,e}
   if !ok{continue}
   // A blocked movement or a turn that the learned model predicts does not
   // change pose is never manufactured into exploratory progress.
   if next.X==n.s.X&&next.Y==n.s.Y&&next.Direction==n.s.Direction{continue}
   first:=n.first;if first<0{first=action}
   cost:=travelStepWeight
   if action!=2{cost+=turnStepWeight}
   if action==2 && n.s.FrontBlocked{cost+=blockedForwardWeight}
   cost+=int(p.Core.Route.Visits[cellKey(next)])*cellVisitWeight
   total:=n.travel+cost
   k:=poseKey(next)
   if old,seen:=bestRoute[k];seen && old<=total{continue}
   bestRoute[k]=total
   if len(bestRoute)>MaxNodes{continue}
   if next.X!=s.X||next.Y!=s.Y{
    score:=total+md(next)*goalDistanceWeight
    if p.Core.Route.Visits[cellKey(next)]==0{score-=unvisitedBonus}
    if score<bestScore{bestScore=score;bestFirst=first}
   }
   if len(q)<MaxNodes{q=append(q,route{s:next,first:first,depth:n.depth+1,travel:total})}
  }
 }
 return bestFirst,bestScore<1<<30,nil
}
// Choose remains deterministic, with a legally matched Feedback required
// before any second action. Baseline R4 chooses when not stalled.
func(p *Policy)Choose(s r1.PublicGrid)(int,error){
 if p==nil||p.Core==nil||p.Core.Core==nil||!legal(s){
  return 0,errors.New("R5_INVALID_PUBLIC_OBSERVATION")
 }
 if p.Core.Route.Pending{return 0,errors.New("R5_PENDING_FEEDBACK")}
 if p.Core.Route.Stationary<StallTrigger &&
  p.Core.Route.Visits[cellKey(s)]<StallTrigger{
  return p.Core.Choose(s)
 }
 // Validate all original R3/R4 model and session constraints before changing
 // the session's pending-action custody.
 if _,e:=p.Core.Snapshot();e!=nil{return 0,e}
 action,ok,e:=p.chooseFrontier(s)
 if e!=nil{return 0,e}
 if !ok{return p.Core.Choose(s)}
 if action<0||action>2{return 0,errors.New("R5_ILLEGAL_FRONTIER_ACTION")}
 p.Core.Route.Pending=true
 p.Core.Route.Before=s
 p.Core.Route.Action=action
 return action,nil
}
func(p *Policy)Snapshot()([]byte,error){
 if p==nil||p.Core==nil{return nil,errors.New("R5_NIL_SNAPSHOT")}
 return p.Core.Snapshot()
}
func Restore(raw []byte)(*Policy,error){
 r,e:=r4.Restore(raw)
 if e!=nil{return nil,e}
 return &Policy{Core:r},nil
}
