package cognition

// Candidate-only experience memory. Raw sensor impressions have no entity IDs,
// memory-slot selectors, future query, or task labels. Temporal paired views
// estimate stable features; capacity-bounded prototypes are allocated on novelty.
// This is NOT an identified causal model: action/outcome are observed associations.
import (
 "bytes"
 "encoding/json"
 "errors"
 "io"
 "math"
)
const ExperienceCapacity=8
type experienceTrack struct {
 Used bool `json:"used"`
 Centre [4]float64 `json:"centre"`
 Values [2]float64 `json:"values"`
 Counts [2]uint32 `json:"counts"`
 Last uint64 `json:"last"`
}
type ExperienceMemory struct {
 Schema int `json:"schema"`
 Noise [4]float64 `json:"noise"`
 Tracks [ExperienceCapacity]experienceTrack `json:"tracks"`
 Updates uint64 `json:"updates"`
}
func NewExperienceMemory() *ExperienceMemory {
 m:=&ExperienceMemory{Schema:1}; for k:=0;k<4;k++{m.Noise[k]=.25};return m
}
func checkSense(a [4]float64)error{
 for _,v:=range a {if math.IsNaN(v)||math.IsInf(v,0)||v< -1||v>1{return errors.New("invalid bounded sensory value")}}
 return nil
}
func (m *ExperienceMemory) validate()error{
 if m==nil||m.Schema!=1{return errors.New("invalid memory version")}
 for _,n:=range m.Noise{if math.IsNaN(n)||math.IsInf(n,0)||n<0||n>2{return errors.New("invalid noise")}}
 for _,t:=range m.Tracks{
  if !t.Used{continue}
  if err:=checkSense(t.Centre);err!=nil{return err}
  for _,v:=range t.Values{if math.IsNaN(v)||math.IsInf(v,0)||v<0||v>1{return errors.New("invalid value")}}
 }
 return nil
}
func (m *ExperienceMemory) closest(sensor [4]float64)(int,float64){
 index:=-1;best:=math.Inf(1)
 for i,t:=range m.Tracks{
  if !t.Used{continue}
  dist:=0.0
  for k:=0;k<4;k++{d:=sensor[k]-t.Centre[k];dist+=d*d/(.08+m.Noise[k])}
  dist/=4
  if dist<best{index,best=i,dist}
 }
 return index,best
}
func (m *ExperienceMemory) findSlot()int{
 for i,t:=range m.Tracks{if !t.Used{return i}}
 oldest:=0;for i:=1;i<ExperienceCapacity;i++{if m.Tracks[i].Last<m.Tracks[oldest].Last{oldest=i}}
 return oldest
}
// ObserveExperience retains bounded raw experience. "after" is the next
// observation of the same physical sensor stream, not a supplied entity ID.
func (m *ExperienceMemory) ObserveExperience(before,after [4]float64,action int,outcome float64)error{
 if err:=m.validate();err!=nil{return err}
 if err:=checkSense(before);err!=nil{return err}
 if err:=checkSense(after);err!=nil{return err}
 if action!=0&&action!=1{return errors.New("invalid action")}
 if math.IsNaN(outcome)||math.IsInf(outcome,0)||outcome<0||outcome>1{return errors.New("invalid observed outcome")}
 if m.Updates==math.MaxUint64{return errors.New("experience counter exhausted")}
 // Learn which sensory coordinates are stable across observed continuity.
 next:=*m
 for k:=0;k<4;k++{next.Noise[k]=.9*next.Noise[k]+.1*math.Abs(after[k]-before[k])}
 index,dist:=next.closest(before)
 if index<0||dist>.27{index=next.findSlot();next.Tracks[index]=experienceTrack{Used:true, Centre:before, Values:[2]float64{.5,.5}}}
 t:=&next.Tracks[index]
 // Centroid adaptation does not depend on query identity or evaluator hints.
 for k:=0;k<4;k++{t.Centre[k]=.85*t.Centre[k]+.15*(before[k]+after[k])/2}
 t.Values[action]=.8*t.Values[action]+.2*outcome
 if t.Counts[action]<math.MaxUint32{t.Counts[action]++}
 next.Updates++
 t.Last=next.Updates
 *m=next
 return nil
}
// InferExperience chooses a prototype through learned feature stability, never
// by a caller-given slot or entity ID. Unseen observations return unknown.
func (m *ExperienceMemory) InferExperience(sensor [4]float64,action int)(float64,bool,error){
 if err:=m.validate();err!=nil{return 0,false,err}
 if err:=checkSense(sensor);err!=nil{return 0,false,err}
 if action!=0&&action!=1{return 0,false,errors.New("invalid action")}
 i,d:=m.closest(sensor)
 if i<0||d>.27||m.Tracks[i].Counts[action]==0{return 0,false,nil}
 return m.Tracks[i].Values[action],true,nil
}
func (m *ExperienceMemory) ExperienceSnapshot()([]byte,error){
 if err:=m.validate();err!=nil{return nil,err}
 return json.Marshal(m)
}
func (m *ExperienceMemory) ExperienceRestore(data []byte)error{
 if len(data)==0||len(data)>8192{return errors.New("invalid snapshot size")}
 dec:=json.NewDecoder(bytes.NewReader(data));dec.DisallowUnknownFields()
 var candidate ExperienceMemory
 if err:=dec.Decode(&candidate);err!=nil{return err}
 var extra any
 if err:=dec.Decode(&extra);err!=io.EOF{return errors.New("trailing snapshot data")}
 if err:=candidate.validate();err!=nil{return err}
 *m=candidate
 return nil
}
