// Package custody enforces pre-original W-F29 data separation and costs.
// It does NOT implement or qualify a learned identity model.
package custody

import (
    "errors"
    "fmt"
    "math"
    "unsafe"
)

const NativeStateCeiling = 576
const ConsensusThreshold = 0.075

// Evidence deliberately contains neither entity identifiers nor truth labels.
// EncounterID is local to an observation episode, not a canonical entity ID.
type Evidence struct {
    EncounterID uint64
    Role uint8
    Left [3]float64
    Right [3]float64
}

type PublicInput struct { Evidence []Evidence }
type PrivateCase struct { Input PublicInput; Truth bool }
type Meter struct {
    NativeBytes int
    HistoryBytes int
    WitnessComparisons int
    Encounters int
}
type Manifest struct {
    Seeds []uint64
    EntityIDs []uint64
    GraphFamilies []uint64
}

func disjointUnique(a, b []uint64) bool {
    have:=make(map[uint64]bool,len(a)+len(b))
    for _,x:=range a { if have[x] { return false };have[x]=true }
    for _,x:=range b { if have[x] { return false };have[x]=true }
    return len(a)>0 && len(b)>0
}

// ValidateSplit fails closed if *any* generator seed, true entity or structural
// family appears in both train and test. This is independent of the predictor.
func ValidateSplit(train,test Manifest) error {
    if !disjointUnique(train.Seeds,test.Seeds) {return errors.New("seed split duplicated or empty")}
    if !disjointUnique(train.EntityIDs,test.EntityIDs) {return errors.New("entity identity split duplicated or empty")}
    if !disjointUnique(train.GraphFamilies,test.GraphFamilies) {return errors.New("graph-family split duplicated or empty")}
    return nil
}

// Prepare returns a copy of the accessible observations; the private truth
// bit never crosses the scorer interface. The outside-core history is charged
// separately from the caller-declared native state. Candidate generation and
// any subsequent pair-search costs require additional meters.
func Prepare(c PrivateCase, nativeBytes int) (PublicInput, Meter, error) {
    if nativeBytes<0 || nativeBytes>NativeStateCeiling {return PublicInput{},Meter{},errors.New("native resident state over budget")}
    if len(c.Input.Evidence)<2 {return PublicInput{},Meter{},errors.New("need independent encounters")}
    seen:=map[uint64]bool{}
    for _,e:=range c.Input.Evidence {
        if seen[e.EncounterID] {return PublicInput{},Meter{},errors.New("duplicate encounter ID")}
        seen[e.EncounterID]=true
        for _,v:=range append(e.Left[:],e.Right[:]...) {
            if math.IsNaN(v)||math.IsInf(v,0) {return PublicInput{},Meter{},errors.New("nonfinite witness")}
        }
    }
    exposed:=append([]Evidence(nil),c.Input.Evidence...)
    historyBytes:=len(exposed)*int(unsafe.Sizeof(Evidence{}))
    return PublicInput{Evidence:exposed},Meter{NativeBytes:nativeBytes,HistoryBytes:historyBytes,Encounters:len(exposed)},nil
}

// Consensus is an intentionally programmed control, NOT the learned model.
func Consensus(input PublicInput, meter *Meter) (bool,error) {
    if meter==nil {return false,errors.New("missing cost meter")}
    if len(input.Evidence)<2 {return false,errors.New("insufficient independent encounters")}
    yes:=0
    for _,e:=range input.Evidence {
        matches:=0
        for j:=0;j<3;j++ {
            x,y:=e.Left[j],e.Right[j]
            if math.IsNaN(x)||math.IsNaN(y)||math.IsInf(x,0)||math.IsInf(y,0) {return false,fmt.Errorf("nonfinite witness %d",j)}
            if math.Abs(x-y)<ConsensusThreshold {matches++}
            meter.WitnessComparisons++
        }
        if matches>=2 {yes++}
    }
    return yes*2>len(input.Evidence),nil
}

// ScoreCase keeps truth inaccessible until after prediction. A future learned
// scorer must receive only PublicInput (not PrivateCase) and fully meter any
// additional search/observations before original execution.
func ScoreCase(c PrivateCase,nativeBytes int,score func(PublicInput,*Meter)(bool,error))(bool,Meter,error) {
    public, meter, err:=Prepare(c,nativeBytes)
    if err!=nil {return false,meter,err}
    answer,err:=score(public,&meter)
    if err!=nil {return false,meter,err}
    return answer==c.Truth,meter,nil
}
