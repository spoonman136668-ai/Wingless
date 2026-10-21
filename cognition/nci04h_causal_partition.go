package cognition

// Opt-in NCI-04H native learner: repartition an aliased percept only when
// a real observed consequence contradicts substantial previous evidence.
// Uses the original LatentState capacity; no extra persistent parameters.
import "errors"

type PartitionUpdate struct {
 Split bool
 Abstained bool
 Count uint8
}

func (s *LatentState) ObservePartitioned(src,dst [FeatureWidth]float32,action,outcome int)(PartitionUpdate,error){
 var zero PartitionUpdate
 if s==nil||!validVector(src)||!validVector(dst)||action<0||action>=MaxActions||(outcome!=0&&outcome!=1){
  return zero,errors.New("invalid causal observation")
 }
 i,di:=nearest(s,src);j,dj:=nearest(s,dst)
 if i>=0&&j>=0&&i!=j&&di<=0.0625&&dj<=0.0625{
  edge:=s.Relations[i][j][action]
  p:=probability(edge.LogOdds)
  if edge.Updates>=12&&((outcome==0&&p>=0.75)||(outcome==1&&p<=0.25)){
   if di<=0.0004{
    return PartitionUpdate{false,true,s.Count},nil
   }
   if int(s.Count)>=MaxEntities{return zero,errors.New("fixed entity budget exhausted")}
   next:=*s
   child:=int(next.Count)
   next.Count++
   next.Entities[child]=EntityPrototype{Mean:src,Seen:1}
   next.Relations[child][j][action]=LatentRelation{LogOdds:float32(0.15*(float64(outcome)-0.5)),Updates:1}
   next.Observations++
   *s=next
   return PartitionUpdate{true,false,next.Count},nil
  }
 }
 if err:=s.Observe(src,dst,action,outcome);err!=nil{return zero,err}
 return PartitionUpdate{false,false,s.Count},nil
}

// SourceReasoner integration admits only positive relations verified by ICE.
// An absent lexical call is not proof of a negative causal consequence.
func (r *SourceReasoner) ObserveVerifiedCallPartition(fromName,toName string)(PartitionUpdate,error){
 var zero PartitionUpdate
 if r==nil{return zero,errors.New("nil source reasoner")}
 if err:=r.validate();err!=nil{return zero,err}
 from,err:=r.symbol(fromName);if err!=nil{return zero,err}
 to,err:=r.symbol(toName);if err!=nil{return zero,err}
 if !lexicalCall(from,to){return zero,errors.New("unverified source relation")}
 return r.State.ObservePartitioned(vectorFromSource(from),vectorFromSource(to),0,1)
}
