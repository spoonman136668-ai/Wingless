package r7
import("errors"; r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1";r6 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r6")
const MaxNodes=2048
type Policy struct{Core *r6.Policy}
func New()*Policy{return &Policy{Core:r6.New()}}
func cell(s r1.PublicGrid)uint32{return uint32(s.X*65+s.Y)}
func pose(s r1.PublicGrid)uint32{return cell(s)*4+uint32(s.Direction)}
func valid(s r1.PublicGrid)bool{return s.X>=0&&s.X<=64&&s.Y>=0&&s.Y<=64&&s.GoalX>=0&&s.GoalX<=64&&s.GoalY>=0&&s.GoalY<=64&&s.Direction>=0&&s.Direction<=3}
type node struct{k uint32;first int}
func(p *Policy)physical(s r1.PublicGrid)(int,bool){
 root:=pose(s);q:=[]node{{k:root,first:-1}};seen:=map[uint32]bool{root:true}
 candidate,found:=0,false
 for i:=0;i<len(q)&&i<MaxNodes;i++{
  n:=q[i];pos:=int(n.k/4);x,y:=pos/65,pos%65
  if n.k!=root&&(x!=s.X||y!=s.Y){
   if x==s.GoalX&&y==s.GoalY{return n.first,true}
   proof:=p.Core.Core.Core.Route.Poses[n.k]
   _,tried:=p.Core.Edges[n.k*3+2]
   if !found&&proof.Seen>0&&!proof.Blocked&&!tried{candidate=n.first;found=true}
  }
  for a:=0;a<3;a++{
   e,ok:=p.Core.Edges[n.k*3+uint32(a)]
   if !ok||e.Count==0||e.Contradictions!=0||e.To==n.k||seen[e.To]{continue}
   seen[e.To]=true;first:=n.first;if first<0{first=a}
   if len(q)<MaxNodes{q=append(q,node{k:e.To,first:first})}
  }
 }
 return candidate,found
}
func(p *Policy)Choose(s r1.PublicGrid)(int,error){
 if p==nil||p.Core==nil||p.Core.Core==nil||!valid(s){return 0,errors.New("R7_INVALID_PUBLIC_OBSERVATION")}
 session:=&p.Core.Core.Core.Route
 if session.Pending{return 0,errors.New("R7_PENDING_REAL_ACTION")}
 if session.Stationary<4&&session.Visits[cell(s)]<4{return p.Core.Core.Choose(s)}
 if _,e:=p.Core.Snapshot();e!=nil{return 0,e}
 action,ok:=p.physical(s)
 if !ok{return p.Core.Core.Choose(s)}
 if action<0||action>2{return 0,errors.New("R7_INVALID_SOURCE_ACTION")}
 session.Pending=true;session.Before=s;session.Action=action
 return action,nil
}
func(p *Policy)Feedback(b r1.PublicGrid,a int,n r1.PublicGrid)error{if p==nil||p.Core==nil{return errors.New("R7_NIL_FEEDBACK")};return p.Core.Feedback(b,a,n)}
func(p *Policy)ObserveFit(b r1.PublicGrid,a int,n r1.PublicGrid,r float64,d bool)error{if p==nil||p.Core==nil{return errors.New("R7_NIL_FIT")};return p.Core.ObserveFit(b,a,n,r,d)}
func(p *Policy)BeginEpisode()error{if p==nil||p.Core==nil{return errors.New("R7_NIL_EPISODE")};return p.Core.BeginEpisode()}
func(p *Policy)Snapshot()([]byte,error){if p==nil||p.Core==nil{return nil,errors.New("R7_NIL_SNAPSHOT")};return p.Core.Snapshot()}
func Restore(raw []byte)(*Policy,error){c,e:=r6.Restore(raw);if e!=nil{return nil,e};return &Policy{Core:c},nil}
