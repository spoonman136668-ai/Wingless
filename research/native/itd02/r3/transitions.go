// Package r3 is an opt-in, source-only action-conditioned transition learner.
// It receives ONLY completed public action/consequence observations. Neither
// the MiniGrid dynamics, a known turn law, nor an obstacle map is hardcoded.
// A bounded graph search over learned transitions is a conventional planning
// mechanism, not proof of novel native cognition or transfer.
package r3

import (
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 "math"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r2 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r2"
)

const Schema = "wingless.itd02.r3.local-action-conditioned-transition-state.v1"
const MaxCells=2048
const MaxNodes=4096
const MinEffectObservations=4

// Effect is learned from (public before, selected action, actual public after).
// No unobserved action gets a retrospective label.
type Effect struct{
 DX int8 `json:"dx"`
 DY int8 `json:"dy"`
 Heading uint8 `json:"heading"`
 Count uint32 `json:"count"`
}
type Consequences struct{
 Samples uint32 `json:"samples"`
 Outcomes [4]Effect `json:"outcomes"`
}
// A pose records the latest sensed occupancy, rather than assuming that a
// previously blocked location stays blocked in an unfamiliar/changed world.
type CellEvidence struct{
 Blocked bool `json:"blocked"`
 Seen uint32 `json:"seen"`
 Contradictions uint32 `json:"contradictions"`
}
type Policy struct{
 Core *r2.Policy
 Learned [24]Consequences
 Cells map[uint32]CellEvidence
 Events uint64
}
func New()*Policy{return &Policy{Core:r2.New(),Cells:make(map[uint32]CellEvidence)}}
func legal(g r1.PublicGrid)bool{
 return g.X>=0&&g.X<=64&&g.Y>=0&&g.Y<=64&&
  g.GoalX>=0&&g.GoalX<=64&&g.GoalY>=0&&g.GoalY<=64&&
  g.Direction>=0&&g.Direction<4
}
func valid(p *Policy)bool{
 if p==nil||p.Core==nil||p.Core.Core==nil||len(p.Cells)>MaxCells||p.Events>=1<<52{return false}
 for _,v:=range p.Learned{
  if v.Samples>=1<<30{return false}
  var sum uint64
  for _,e:=range v.Outcomes{
   if e.Count>0&&(e.DX < -1||e.DX>1||e.DY < -1||e.DY>1||e.Heading>3){return false}
   sum+=uint64(e.Count)
  }
  if sum!=uint64(v.Samples){return false}
 }
 for k,v:=range p.Cells{
  if k>(64<<9|64<<2|3)||v.Seen==0||v.Seen>=1<<30||v.Contradictions>v.Seen{return false}
  x,y,d:=int(k>>9),int((k>>2)&127),int(k&3)
  if x>64||y>64||d>3{return false}
 }
 return true
}
func poseKey(x,y,d int)uint32{return uint32((x<<9)|(y<<2)|d)}
func effectKey(s r1.PublicGrid,a int)int{
 b:=0;if s.FrontBlocked{b=1}
 return (s.Direction*2+b)*3+a
}
func same(a,b Effect)bool{return a.DX==b.DX&&a.DY==b.DY&&a.Heading==b.Heading}
func updateConsequences(c *Consequences,e Effect)error{
 if c.Samples>=1<<30-1{return errors.New("R3_EXHAUSTED_CONSEQUENCE_BUDGET")}
 slot:=-1
 for i:=range c.Outcomes {
  if c.Outcomes[i].Count>0 && same(c.Outcomes[i],e){slot=i;break}
 }
 if slot<0 {for i:=range c.Outcomes {if c.Outcomes[i].Count==0{slot=i;break}}}
 if slot<0{return errors.New("R3_CONSEQUENCE_VARIANTS_EXHAUSTED")}
 if c.Outcomes[slot].Count==0{c.Outcomes[slot]=e}
 c.Outcomes[slot].Count++
 c.Samples++
 return nil
}
// Observe is transactional: if any public observation, count, capacity, or
// downstream native R2 update is invalid, neither R2 nor R3 is changed.
func(p *Policy)Observe(before r1.PublicGrid,action int,after r1.PublicGrid,reward float64,terminal bool)error{
 if !valid(p)||!legal(before)||!legal(after)||action<0||action>2||
  math.IsNaN(reward)||math.IsInf(reward,0)||math.Abs(reward)>100{
  return errors.New("R3_INVALID_PUBLIC_CONSEQUENCE")
 }
 dx,dy:=after.X-before.X,after.Y-before.Y
 if dx < -1||dx>1||dy < -1||dy>1{return errors.New("R3_NONLOCAL_TRANSITION_REJECTED")}
 key:=poseKey(before.X,before.Y,before.Direction)
 beforeCell,seen:=p.Cells[key]
 if !seen&&len(p.Cells)>=MaxCells{return errors.New("R3_CELL_BUDGET_EXHAUSTED")}
 if seen&&beforeCell.Seen>=1<<30-1{return errors.New("R3_CELL_COUNTER_EXHAUSTED")}
 nextEvidence:=p.Learned[effectKey(before,action)]
 if e:=updateConsequences(&nextEvidence,Effect{DX:int8(dx),DY:int8(dy),Heading:uint8(after.Direction)});e!=nil{return e}
 nextCell:=beforeCell
 if seen&&nextCell.Blocked!=before.FrontBlocked{nextCell.Contradictions++}
 nextCell.Blocked=before.FrontBlocked
 nextCell.Seen++
 nextCore:=*p.Core
 if e:=nextCore.Observe(before,action,after,reward,terminal);e!=nil{return e}
 p.Core=&nextCore
 p.Learned[effectKey(before,action)]=nextEvidence
 if p.Cells==nil{p.Cells=make(map[uint32]CellEvidence)}
 p.Cells[key]=nextCell
 p.Events++
 return nil
}
// Predict is an empirical, action-conditioned next-pose distribution's modal
// outcome. No prediction is returned until enough real evidence exists.
func(p *Policy)Predict(s r1.PublicGrid,action int)(r1.PublicGrid,bool,error){
 if !valid(p)||!legal(s)||action<0||action>2{return r1.PublicGrid{},false,errors.New("R3_INVALID_PREDICT_INPUT")}
 c:=p.Learned[effectKey(s,action)]
 if c.Samples<MinEffectObservations{return r1.PublicGrid{},false,nil}
 best:=Effect{}
 for _,v:=range c.Outcomes{if v.Count>best.Count{best=v}}
 if uint64(best.Count)*4<uint64(c.Samples)*3{return r1.PublicGrid{},false,nil}
 nx,ny:=s.X+int(best.DX),s.Y+int(best.DY)
 if nx<0||nx>64||ny<0||ny>64{return r1.PublicGrid{},false,nil}
 out:=s
 out.X=nx;out.Y=ny;out.Direction=int(best.Heading)
 if c,ok:=p.Cells[poseKey(nx,ny,out.Direction)];ok{out.FrontBlocked=c.Blocked}else{out.FrontBlocked=false}
 return out,true,nil
}
type routeNode struct{
 State r1.PublicGrid
 First int
 Depth int
}
// Plan searches ONLY the learned transition model, including the latest
// observed local affordance. For unobserved future poses it assumes public
// front clearance provisionally; future real feedback can contradict this.
// If no reliable route is available it returns false (native R2 fallback).
func(p *Policy)Plan(s r1.PublicGrid)(int,bool,error){
 if !valid(p)||!legal(s){return 0,false,errors.New("R3_INVALID_PLAN_INPUT")}
 if s.X==s.GoalX&&s.Y==s.GoalY{return 0,false,nil}
 queue:=[]routeNode{{State:s,First:-1,Depth:0}}
 visited:=map[uint32]bool{poseKey(s.X,s.Y,s.Direction):true}
 for head:=0;head<len(queue)&&head<MaxNodes;head++{
  n:=queue[head]
  if n.Depth>=64{continue}
  for a:=0;a<3;a++{
   nxt,ok,e:=p.Predict(n.State,a)
   if e!=nil{return 0,false,e}
   if !ok{continue}
   first:=n.First;if first<0{first=a}
   if nxt.X==s.GoalX&&nxt.Y==s.GoalY{return first,true,nil}
   key:=poseKey(nxt.X,nxt.Y,nxt.Direction)
   if visited[key]{continue}
   visited[key]=true
   if len(queue)>=MaxNodes{return 0,false,nil}
   queue=append(queue,routeNode{State:nxt,First:first,Depth:n.Depth+1})
  }
 }
 return 0,false,nil
}
func(p *Policy)Choose(s r1.PublicGrid)(int,error){
 a,ok,e:=p.Plan(s);if e!=nil{return 0,e}
 if ok{return a,nil}
 return p.Core.Choose(s)
}
type payload struct{
 R2 []byte `json:"original_r2_checkpoint"`
 Learned [24]Consequences `json:"learned_effects"`
 Cells map[uint32]CellEvidence `json:"public_cell_evidence"`
 Events uint64 `json:"completed_events"`
}
type signed struct{
 Schema string `json:"schema"`
 Data payload `json:"data"`
 SHA256 string `json:"sha256"`
}
func(p *Policy)Snapshot()([]byte,error){
 if !valid(p){return nil,errors.New("R3_INVALID_CHECKPOINT_MODEL")}
 core,e:=p.Core.Snapshot();if e!=nil{return nil,e}
 data:=payload{R2:core,Learned:p.Learned,Cells:p.Cells,Events:p.Events}
 raw,e:=json.Marshal(data);if e!=nil{return nil,e}
 sum:=sha256.Sum256(raw)
 return json.Marshal(signed{Schema:Schema,Data:data,SHA256:hex.EncodeToString(sum[:])})
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>2<<20{return nil,errors.New("R3_CHECKPOINT_SIZE")}
 d:=json.NewDecoder(bytes.NewReader(raw))
 d.DisallowUnknownFields()
 var s signed
 if e:=d.Decode(&s);e!=nil{return nil,e}
 if e:=d.Decode(new(struct{}));e!=io.EOF{return nil,errors.New("R3_CHECKPOINT_TRAILING_CONTENT")}
 if s.Schema!=Schema{return nil,errors.New("R3_CHECKPOINT_SCHEMA")}
 body,e:=json.Marshal(s.Data);if e!=nil{return nil,e}
 h:=sha256.Sum256(body)
 if s.SHA256!=hex.EncodeToString(h[:]){return nil,errors.New("R3_CHECKPOINT_AUTHENTICATION")}
 core,e:=r2.Restore(s.Data.R2);if e!=nil{return nil,e}
 cells:=s.Data.Cells;if cells==nil{cells=make(map[uint32]CellEvidence)}
 p:=&Policy{Core:core,Learned:s.Data.Learned,Cells:cells,Events:s.Data.Events}
 if !valid(p)||p.Events!=p.Core.Core.Updates{return nil,errors.New("R3_CHECKPOINT_NATIVE_CUSTODY")}
 return p,nil
}
