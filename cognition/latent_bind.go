// Package cognition implements opt-in, candidate-only native cognitive mechanisms.
// This program does not establish scientific novelty or grant execution authority.
package cognition

import (
    "encoding/json"
    "errors"
    "math"
)

const MaxEntities=16
const MaxActions=3
const FeatureWidth=4

// LatentRelation is a learned directed association between two dynamically
// discovered sensory identities, not an explicit symbol-labelled rule table.
type LatentRelation struct {
    LogOdds float32 `json:"log_odds"`
    Updates uint16 `json:"updates"`
}
type EntityPrototype struct {
    Mean [FeatureWidth]float32 `json:"mean"`
    Seen uint16 `json:"seen"`
}
type LatentState struct {
    Entities [MaxEntities]EntityPrototype `json:"entities"`
    Relations [MaxEntities][MaxEntities][MaxActions]LatentRelation `json:"relations"`
    Count uint8 `json:"count"`
    Observations uint32 `json:"observations"`
}

func validVector(v [FeatureWidth]float32)bool{
    for _,x:=range v {if math.IsNaN(float64(x))||math.IsInf(float64(x),0)||x < -1||x>1{return false}}
    return true
}
func nearest(s *LatentState, v [FeatureWidth]float32)(int,float32){
    idx:=-1
    best:=float32(100)
    for i:=0;i<int(s.Count);i++ {
        var d float32
        for k,x:=range v{q:=x-s.Entities[i].Mean[k];d+=q*q}
        if d<best{best,idx=d,i}
    }
    return idx,best
}
func (s *LatentState) bind(v [FeatureWidth]float32)(int,error){
    i,d:=nearest(s,v)
    if i<0||d>0.0625 {
        if int(s.Count)>=MaxEntities{return 0,errors.New("entity capacity exceeded")}
        i=int(s.Count)
        s.Count++
        s.Entities[i]=EntityPrototype{Mean:v,Seen:1}
        return i,nil
    }
    e:= &s.Entities[i]
    if e.Seen<65535 {e.Seen++}
    factor:=float32(1.0)/float32(e.Seen)
    for k,x:=range v {e.Mean[k]+=factor*(x-e.Mean[k])}
    return i,nil
}
func probability(v float32)float64{
    z:=float64(v)
    if z>=0 {q:=math.Exp(-z);return 1/(1+q)}
    q:=math.Exp(z);return q/(1+q)
}
// contradictionCandidate creates a distinct causal state only when the
// EXISTING trained action relation strongly contradicts a newly observed
// consequence AND the new observation is perceptually distinguishable.
// It adds no fields: prototypes and relations use the original fixed budget.
// A novel state has no confident answer until it learns from consequences.
// It never splits two truly identical observations; those remain uncertain.
// Candidate-only learning policy; not an independent scientific threshold.
func (s *LatentState) contradictionCandidate(src,dst [FeatureWidth]float32,action,outcome int)(bool,bool){
    i,ds:=nearest(s,src)
    j,dd:=nearest(s,dst)
    if i<0||j<0||i==j||ds>0.0625||dd>0.0625{return false,false}
    edge:=s.Relations[i][j][action]
    if edge.Updates<24{return false,false}
    p:=probability(edge.LogOdds)
    if !((outcome==0&&p>=0.90)||(outcome==1&&p<=0.10)){return false,false}
    // A distant observation is already a new entity under normal binding.
    // Split only the perceptually divergent endpoint, not a task label.
    if ds>=dd&&ds>0.0016{return true,false}
    if dd>0.0016{return false,true}
    return false,false
}
func (s *LatentState) ambiguousNeighbor(v [FeatureWidth]float32,chosen int,chosenDistance float32)bool{
    for i:=0;i<int(s.Count);i++ {
        if i==chosen{continue}
        var d float32
        for k,x:=range v{q:=x-s.Entities[i].Mean[k];d+=q*q}
        if d<=0.0625&&math.Abs(float64(d-chosenDistance))<=0.0008{return true}
    }
    return false
}

// Observe receives a sensory pair, an action and an observed consequence.
// No task label, symbolic entity ID or future query is supplied.
func (s *LatentState) Observe(src,dst [FeatureWidth]float32,action,outcome int)error{
    if !validVector(src)||!validVector(dst)||action<0||action>=MaxActions||(outcome!=0&&outcome!=1){return errors.New("invalid latent observation")}
    next:=*s
    splitSource,splitTarget:=next.contradictionCandidate(src,dst,action,outcome)
    // Transactional capacity handling. Never silently overwrite a learned
    // prototype or convert a contradictory single observation into truth.
    if (splitSource||splitTarget)&&int(next.Count)>=MaxEntities{
        return errors.New("causal state revision exceeds fixed entity capacity")
    }
    var i,j int
    var e error
    if splitSource{
        i=int(next.Count)
        next.Count++
        next.Entities[i]=EntityPrototype{Mean:src,Seen:1}
    }else{
        i,e=next.bind(src);if e!=nil{return e}
    }
    if splitTarget{
        j=int(next.Count)
        next.Count++
        next.Entities[j]=EntityPrototype{Mean:dst,Seen:1}
    }else{
        j,e=next.bind(dst);if e!=nil{return e}
    }
    edge:= &next.Relations[i][j][action]
    p:=probability(edge.LogOdds)
    v:=float64(edge.LogOdds)+0.15*(float64(outcome)-p)
    if v>12 {v=12};if v< -12 {v= -12}
    edge.LogOdds=float32(v)
    if edge.Updates<65535 {edge.Updates++}
    // Preserve the original fixed-size online evidence counter. Saturating
    // instead of wrapping avoids treating exhausted memory as untrained state.
    // The staged copy has not been committed if this guard fails.
    if next.Observations == ^uint32(0) {
        return errors.New("native causal observation counter exhausted")
    }
    next.Observations++
    *s=next
    return nil
}
// Queries are strictly read-only. Unknown entities or relations use 0.5 prior.
func (s *LatentState) Predict(src,dst [FeatureWidth]float32,action int)(float64,error){
    if !validVector(src)||!validVector(dst)||action<0||action>=MaxActions{return 0,errors.New("invalid latent query")}
    i,a:=nearest(s,src);j,b:=nearest(s,dst)
    if i<0||j<0||a>0.0625||b>0.0625{return 0.5,nil}
    // A boundary between learned causal states carries insufficient
    // evidence for confident identity binding. Abstain rather than choose
    // an arbitrary nearest prototype.
    if s.ambiguousNeighbor(src,i,a)||s.ambiguousNeighbor(dst,j,b){return 0.5,nil}
    edge:=s.Relations[i][j][action]
    if edge.Updates==0{return 0.5,nil}
    return probability(edge.LogOdds),nil
}
func (s *LatentState) ForecastVia(src,via,dst [FeatureWidth]float32,first,second int)(float64,error){
    a,e:=s.Predict(src,via,first);if e!=nil{return 0,e}
    b,e:=s.Predict(via,dst,second);if e!=nil{return 0,e}
    return a*b,nil
}
func (s *LatentState) Snapshot()([]byte,error){return json.Marshal(s)}
func (s *LatentState) Restore(data []byte)error{
    if len(data)>65536{return errors.New("oversized latent snapshot")}
    var next LatentState
    if e:=json.Unmarshal(data,&next);e!=nil{return e}
    if next.Count>MaxEntities{return errors.New("invalid entity count")}
    for i:=0;i<int(next.Count);i++{if !validVector(next.Entities[i].Mean)||next.Entities[i].Seen==0{return errors.New("invalid prototype")}}
    for i:=0;i<MaxEntities;i++{for j:=0;j<MaxEntities;j++{for a:=0;a<MaxActions;a++{
        r:=next.Relations[i][j][a]
        if math.IsNaN(float64(r.LogOdds))||math.IsInf(float64(r.LogOdds),0)||math.Abs(float64(r.LogOdds))>12{return errors.New("invalid relation")}
        if (i>=int(next.Count)||j>=int(next.Count))&&r.Updates!=0{return errors.New("orphan relation")}
    }}}
    *s=next;return nil
}
