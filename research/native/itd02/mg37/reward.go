// Package mg37 provides an opt-in, reward-conditioned decision guard for MG28.
// It never reads hidden world state and does not run original environment steps.
package mg37

import ("errors"; "math")

const MinEvidence = 3
const ImprovementMargin = 0.1

type Evidence struct { Count uint32; RewardSum float64 }
type Decision struct { Action int; Override bool }

// Choose requires paired evidence for the parent proposal and any alternative.
// Unknown outcomes remain unknown: no override is allowed on mere image changes.
func Choose(parent int, evidence [7]Evidence) (Decision, error) {
 if parent < 0 || parent >= 7 { return Decision{},errors.New("MG37_INVALID_PARENT_ACTION") }
 for _, e := range evidence {
  if math.IsNaN(e.RewardSum)||math.IsInf(e.RewardSum,0)||math.Abs(e.RewardSum)>float64(e.Count) {
   return Decision{}, errors.New("MG37_INVALID_OBSERVED_REWARD")
  }
 }
 result:=Decision{Action:parent}
 baseline:=evidence[parent]
 if baseline.Count<MinEvidence { return result,nil }
 baselineMean:=baseline.RewardSum/float64(baseline.Count)
 bestMean:=baselineMean
 for action,e:=range evidence {
  if action==parent || e.Count<MinEvidence {continue}
  mean:=e.RewardSum/float64(e.Count)
  if mean>0 && mean>baselineMean+ImprovementMargin && mean>bestMean+ImprovementMargin {
   result=Decision{Action:action,Override:true};bestMean=mean
  }
 }
 return result,nil
}
