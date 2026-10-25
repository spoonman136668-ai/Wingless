// Package r6 layers authenticated episode-local observed transition edges over
// unchanged native R5. Unlike predicted R3 topology, these edges only enter
// memory after an actually chosen public action and its observed consequence.
package r6

import(
 "bytes"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "io"
 r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
 r5 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r5"
)
const Schema="wingless.itd02.r6.observed-edge-routing.v1"
const MaxEdges=4096
const MaxEvents=4096
const MaxGraphNodes=2048
type Edge struct{
 To uint32 `json:"to"`
 Count uint16 `json:"count"`
 Contradictions uint16 `json:"contradictions"`
}
type Policy struct{
 Core *r5.Policy
 Edges map[uint32]Edge
 Seen uint32
}
func New()*Policy{return &Policy{Core:r5.New(),Edges:make(map[uint32]Edge)}}
func validPublic(s r1.PublicGrid)bool{
 return s.X>=0&&s.X<=64&&s.Y>=0&&s.Y<=64&&
 s.GoalX>=0&&s.GoalX<=64&&s.GoalY>=0&&s.GoalY<=64&&
 s.Direction>=0&&s.Direction<=3
}
func posKey(s r1.PublicGrid)uint32{return uint32((s.X*65+s.Y)*4+s.Direction)}
func edgeKey(s r1.PublicGrid,a int)uint32{return posKey(s)*3+uint32(a)}
func decodePos(k uint32)(x,y,dir int){
 dir=int(k%4)
 n:=int(k/4)
 return n/65,n%65,dir
}
func(p *Policy)valid()bool{
 if p==nil||p.Core==nil||p.Core.Core==nil||p.Core.Core.Core==nil||
  len(p.Edges)>MaxEdges||p.Seen>MaxEvents||
  p.Seen!=p.Core.Core.Route.Steps{return false}
 var count uint64
 for k,v:=range p.Edges{
  if k>=65*65*4*3||v.Count==0||v.Contradictions>v.Count||
   v.To>=65*65*4{return false}
  count+=uint64(v.Count)
 }
 return count==uint64(p.Seen)
}
func(p *Policy)BeginEpisode()error{
 if !p.valid(){return errors.New("R6_INVALID_EPISODE_RESET")}
 if e:=p.Core.BeginEpisode();e!=nil{return e}
 p.Edges=make(map[uint32]Edge)
 p.Seen=0
 return nil
}
func(p *Policy)ObserveFit(s r1.PublicGrid,a int,n r1.PublicGrid,reward float64,done bool)error{
 if !p.valid()||p.Seen!=0||p.Core.Core.Route.Pending{
  return errors.New("R6_FIT_INTERLEAVED_WITH_REAL_EPISODE")
 }
 return p.Core.ObserveFit(s,a,n,reward,done)
}
func(p *Policy)Feedback(s r1.PublicGrid,a int,n r1.PublicGrid)error{
 if !p.valid()||!validPublic(s)||!validPublic(n)||a<0||a>2||
  !p.Core.Core.Route.Pending||p.Core.Core.Route.Before!=s||
  p.Core.Core.Route.Action!=a||p.Seen>=MaxEvents{
  return errors.New("R6_UNSELECTED_OR_INVALID_SOURCE_CONSEQUENCE")
 }
 k:=edgeKey(s,a)
 old,ok:=p.Edges[k]
 if !ok&&len(p.Edges)>=MaxEdges{return errors.New("R6_EDGE_CAPACITY")}
 if old.Count>=65535{return errors.New("R6_EDGE_COUNT_CAPACITY")}
 // Preflight the original R4 native source transaction: no graph learns
 // anything until original public-action custody has accepted feedback.
 if e:=p.Core.Feedback(s,a,n);e!=nil{return e}
 if !ok{old=Edge{To:posKey(n)}}
 if old.To!=posKey(n){old.Contradictions++}
 old.Count++
 p.Edges[k]=old
 p.Seen++
 return nil
}
type path struct{k uint32;first int}
func(p *Policy)observedShortest(s r1.PublicGrid)(int,bool){
 if !p.valid(){return 0,false}
 start:=posKey(s)
 // A previously verified goal coordinate is usable; no unobserved route
 // is materialized from the public goal coordinates alone.
 adj:=make(map[uint32][]path)
 for key,e:=range p.Edges{
  if e.Contradictions!=0{continue}
  source:=key/3
  action:=int(key%3)
  if source==e.To{continue}
  adj[source]=append(adj[source],path{k:e.To,first:action})
 }
 // Queue in deterministic action order; shortest edges by actual observed
 // action count, rather than task-specific or learned-model imagined steps.
 q:=[]path{{k:start,first:-1}}
 seen:=map[uint32]bool{start:true}
 bestFrontier,frontierDepth:=0,-1
 for head:=0;head<len(q)&&head<MaxGraphNodes;head++{
  node:=q[head]
  x,y,_:=decodePos(node.k)
  if node.k!=start&&x==s.GoalX&&y==s.GoalY{return node.first,true}
  // When goal connectivity is not known, prefer the nearest recorded
  // reachable public pose with missing observed outgoing actions.
  if node.k!=start && frontierDepth<0{
   tried:=0
   for action:=0;action<3;action++{
    if _,ok:=p.Edges[node.k*3+uint32(action)];ok{tried++}
   }
   if tried<3{bestFrontier=node.first;frontierDepth=head}
  }
  for action:=0;action<3;action++{
   k:=node.k*3+uint32(action)
   e,ok:=p.Edges[k]
   if !ok||e.Contradictions!=0||e.To==node.k||seen[e.To]{continue}
   seen[e.To]=true
   first:=node.first;if first<0{first=action}
   q=append(q,path{k:e.To,first:first})
   if len(q)>=MaxGraphNodes{break}
  }
 }
 if frontierDepth>=0{return bestFrontier,true}
 return 0,false
}
func(p *Policy)Choose(s r1.PublicGrid)(int,error){
 if !p.valid()||!validPublic(s)||p.Core.Core.Route.Pending{
  return 0,errors.New("R6_INVALID_SOURCE_ACTION_REQUEST")
 }
 action,ok:=p.observedShortest(s)
 if !ok{return p.Core.Choose(s)}
 if action<0||action>2{return 0,errors.New("R6_ILLEGAL_VERIFIED_EDGE")}
 p.Core.Core.Route.Pending=true
 p.Core.Core.Route.Before=s
 p.Core.Core.Route.Action=action
 return action,nil
}
type data struct{
 OriginalR5 []byte `json:"original_r5_model"`
 Edges map[uint32]Edge `json:"observed_public_edges"`
 Seen uint32 `json:"actual_public_steps"`
}
type envelope struct{
 Schema string `json:"schema"`
 Data data `json:"data"`
 SHA string `json:"sha256"`
}
func(p *Policy)Snapshot()([]byte,error){
 if !p.valid(){return nil,errors.New("R6_INVALID_MODEL_SNAPSHOT")}
 original,e:=p.Core.Snapshot();if e!=nil{return nil,e}
 d:=data{OriginalR5:original,Edges:p.Edges,Seen:p.Seen}
 b,e:=json.Marshal(d);if e!=nil{return nil,e}
 h:=sha256.Sum256(b)
 return json.Marshal(envelope{Schema:Schema,Data:d,SHA:hex.EncodeToString(h[:])})
}
func Restore(raw []byte)(*Policy,error){
 if len(raw)==0||len(raw)>3<<20{return nil,errors.New("R6_CHECKPOINT_SIZE")}
 dec:=json.NewDecoder(bytes.NewReader(raw))
 dec.DisallowUnknownFields()
 var env envelope
 if e:=dec.Decode(&env);e!=nil{return nil,e}
 if e:=dec.Decode(new(struct{}));e!=io.EOF{return nil,errors.New("R6_CHECKPOINT_TRAILING_BYTES")}
 if env.Schema!=Schema{return nil,errors.New("R6_CHECKPOINT_SCHEMA")}
 body,e:=json.Marshal(env.Data);if e!=nil{return nil,e}
 checksum:=sha256.Sum256(body)
 if env.SHA!=hex.EncodeToString(checksum[:]){return nil,errors.New("R6_CHECKPOINT_HASH")}
 core,e:=r5.Restore(env.Data.OriginalR5);if e!=nil{return nil,e}
 p:=&Policy{Core:core,Edges:env.Data.Edges,Seen:env.Data.Seen}
 if p.Edges==nil{p.Edges=make(map[uint32]Edge)}
 if !p.valid(){return nil,errors.New("R6_CHECKPOINT_SOURCE_CUSTODY")}
 return p,nil
}
