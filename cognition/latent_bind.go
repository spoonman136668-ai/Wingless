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
// Observe receives a sensory pair, an action and an observed consequence.
// No task label, symbolic entity ID or future query is supplied.
func (s *LatentState) Observe(src,dst [FeatureWidth]float32,action,outcome int)error{
    if !validVector(src)||!validVector(dst)||action<0||action>=MaxActions||(outcome!=0&&outcome!=1){return errors.New("invalid latent observation")}
    next:=*s
    i,e:=next.bind(src);if e!=nil{return e}
    j,e:=next.bind(dst);if e!=nil{return e}
    edge:= &next.Relations[i][j][action]
    p:=probability(edge.LogOdds)
    v:=float64(edge.LogOdds)+0.15*(float64(outcome)-p)
    if v>12 {v=12};if v< -12 {v= -12}
    edge.LogOdds=float32(v)
    if edge.Updates<65535 {edge.Updates++}
    next.Observations++
    *s=next
    return nil
}
// Queries are strictly read-only. Unknown entities or relations use 0.5 prior.
func (s *LatentState) Predict(src,dst [FeatureWidth]float32,action int)(float64,error){
    if !validVector(src)||!validVector(dst)||action<0||action>=MaxActions{return 0,errors.New("invalid latent query")}
    i,a:=nearest(s,src);j,b:=nearest(s,dst)
    if i<0||j<0||a>0.0625||b>0.0625{return 0.5,nil}
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
