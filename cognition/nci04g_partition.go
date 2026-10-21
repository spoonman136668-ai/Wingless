// Package cognition contains development-only mechanisms; no acceptance authority.
package cognition

import (
 "encoding/json"
 "errors"
 "math"
)

var ErrUnknownCausalState = errors.New("causal evidence insufficient")

// PartitionMode is a consequence-distinguished hypothesis, not a caller-supplied identity.
type PartitionMode struct {
 Sensor [FeatureWidth]float32 `json:"sensor"`
 Probe [FeatureWidth]float32 `json:"probe"`
 Consequence uint8 `json:"consequence"`
 Function uint8 `json:"function"`
 Samples uint16 `json:"samples"`
}

// PartitionProbe is a learned feature-bearing intervention prototype; the
// implementation has no designer-supplied diagnostic action IDs.
type PartitionProbe struct {
 Sensor [FeatureWidth]float32 `json:"sensor"`
 Descriptor [FeatureWidth]float32 `json:"descriptor"`
 Outcomes [2]uint16 `json:"outcomes"`
}

// CausalPartition is a bounded, opt-in *development* candidate.
// Conflicting consequences create competing modes without rewriting the old mode.
// It cannot identify a new hidden state before observing a distinguishing probe.
type CausalPartition struct {
 Modes [MaxEntities]PartitionMode `json:"modes"`
 Count uint8 `json:"count"`
 Probes [MaxActions]PartitionProbe `json:"probes"`
 ProbeCount uint8 `json:"probe_count"`
}

func equalVec(a,b [FeatureWidth]float32) bool {
 for i:=range a { if a[i]!=b[i] {return false} }
 return true
}

func (s *CausalPartition) ObserveCausal(sensor,probe [FeatureWidth]float32,consequence,function int) error {
 if !validVector(sensor)||!validVector(probe)||(consequence!=0&&consequence!=1)||(function!=0&&function!=1) {
  return errors.New("invalid causal observation")
 }
 next:=*s // commit only after every capacity/integrity check has passed
 pidx:=-1
 for i:=0;i<int(next.ProbeCount);i++ {
  if equalVec(next.Probes[i].Sensor,sensor)&&equalVec(next.Probes[i].Descriptor,probe) {pidx=i;break}
 }
 if pidx<0 {
  if int(next.ProbeCount)>=MaxActions{return errors.New("probe capacity exceeded")}
  pidx=int(next.ProbeCount);next.ProbeCount++
  next.Probes[pidx]=PartitionProbe{Sensor:sensor,Descriptor:probe}
 }
 midx:=-1
 for i:=0;i<int(next.Count);i++ {
  m:=next.Modes[i]
  if equalVec(m.Sensor,sensor)&&equalVec(m.Probe,probe)&&int(m.Consequence)==consequence {
   midx=i;break
  }
 }
 if midx<0 {
  if int(next.Count)>=MaxEntities{return errors.New("causal partition capacity exceeded")}
  midx=int(next.Count);next.Count++
  next.Modes[midx]=PartitionMode{Sensor:sensor,Probe:probe,Consequence:uint8(consequence),Function:uint8(function),Samples:0}
 } else if int(next.Modes[midx].Function)!=function {
  // No hidden authority can silently overwrite an established relation.
  return errors.New("conflicting function evidence: abstain")
 }
 if next.Modes[midx].Samples==math.MaxUint16 || next.Probes[pidx].Outcomes[consequence]==math.MaxUint16 {
  return errors.New("counter capacity exceeded")
 }
 next.Modes[midx].Samples++
 next.Probes[pidx].Outcomes[consequence]++
 *s=next
 return nil
}

// SelectIntervention uses ONLY variation actually observed for each offered
// descriptor. An unfamiliar or nondiscriminating set stays explicitly unknown.
func (s *CausalPartition) SelectIntervention(sensor [FeatureWidth]float32,offered [][FeatureWidth]float32)([FeatureWidth]float32,bool,error) {
 if !validVector(sensor)||len(offered)==0||len(offered)>MaxActions{return [FeatureWidth]float32{},false,errors.New("invalid action affordances")}
 best:=-1
 var strength uint16
 for i,probe:=range offered {
  if !validVector(probe){return [FeatureWidth]float32{},false,errors.New("invalid intervention descriptor")}
  for k:=0;k<int(s.ProbeCount);k++ {
   p:=s.Probes[k]
   if !equalVec(sensor,p.Sensor)||!equalVec(probe,p.Descriptor){continue}
   contrast:=p.Outcomes[0];if p.Outcomes[1]<contrast {contrast=p.Outcomes[1]}
   if contrast>strength {strength=contrast;best=i}
  }
 }
 if best<0{return [FeatureWidth]float32{},false,nil}
 return offered[best],true,nil
}

// Reidentify returns an earlier learned relation only AFTER a response has
// distinguished a causal mode; missing and ambiguous evidence never asserts it.
func (s *CausalPartition) Reidentify(sensor,probe [FeatureWidth]float32,consequence int)(int,bool,error) {
 if !validVector(sensor)||!validVector(probe)||(consequence!=0&&consequence!=1){return 0,false,errors.New("invalid causal query")}
 for i:=0;i<int(s.Count);i++ {
  m:=s.Modes[i]
  if equalVec(m.Sensor,sensor)&&equalVec(m.Probe,probe)&&int(m.Consequence)==consequence&&m.Samples>0 {
   return int(m.Function),true,nil
  }
 }
 return 0,false,nil
}

func (s *CausalPartition) CausalSnapshot()([]byte,error){return json.Marshal(s)}

func (s *CausalPartition) CausalRestore(raw []byte) error {
 if len(raw)>16384{return errors.New("oversized causal snapshot")}
 var next CausalPartition
 if err:=json.Unmarshal(raw,&next);err!=nil{return err}
 if int(next.Count)>MaxEntities||int(next.ProbeCount)>MaxActions{return errors.New("invalid causal capacity")}
 for i:=0;i<int(next.Count);i++{
  m:=next.Modes[i]
  if !validVector(m.Sensor)||!validVector(m.Probe)||m.Consequence>1||m.Function>1||m.Samples==0{return errors.New("invalid causal mode")}
  for j:=0;j<i;j++{old:=next.Modes[j];if equalVec(m.Sensor,old.Sensor)&&equalVec(m.Probe,old.Probe)&&m.Consequence==old.Consequence{return errors.New("duplicate causal mode")}}
 }
 for i:=0;i<int(next.ProbeCount);i++ {
  p:=next.Probes[i]
  if !validVector(p.Sensor)||!validVector(p.Descriptor){return errors.New("invalid causal probe")}
  for j:=0;j<i;j++{old:=next.Probes[j];if equalVec(p.Sensor,old.Sensor)&&equalVec(p.Descriptor,old.Descriptor){return errors.New("duplicate probe")}
  }
  counts:=[2]uint32{}
  for j:=0;j<int(next.Count);j++ {m:=next.Modes[j];if equalVec(m.Sensor,p.Sensor)&&equalVec(m.Probe,p.Descriptor){counts[m.Consequence]+=uint32(m.Samples)}}
  if uint32(p.Outcomes[0])!=counts[0]||uint32(p.Outcomes[1])!=counts[1]{return errors.New("probe evidence mismatch")}
 }
 *s=next
 return nil
}
