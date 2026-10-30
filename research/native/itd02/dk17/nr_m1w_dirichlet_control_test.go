package dk17

// NR-M1W: independent conventional Dirichlet control. This file is only a
// source-free Go test, not original-world policy integration.
import (
    "math/big"
    "testing"
)

type referencePrior struct {
    changed int64
    noEffect int64
    other int64
}

// Dirichlet(1,1,62) models the same three completed selected-outcome classes
// as source DK17. The exact posterior EXPECTED changed-minus-no-effect rate
// is (Changed-NoEffect)/(N+64), before native 4096-point quantization.
func dirichletScore(s Stat) (int64, bool) {
    if s.Changed > s.N || s.NoEffect > s.N || s.Changed+s.NoEffect > s.N {
        return 0, false
    }
    n := int64(s.N)
    if n > MaxSamples {
        return 0, false
    }
    prior := referencePrior{changed:1,noEffect:1,other:62}
    changed := big.NewRat(int64(s.Changed)+prior.changed,n+prior.changed+prior.noEffect+prior.other)
    noEffect := big.NewRat(int64(s.NoEffect)+prior.noEffect,n+prior.changed+prior.noEffect+prior.other)
    contrast := new(big.Rat).Sub(changed,noEffect)
    contrast.Mul(contrast,big.NewRat(4096,1))
    // big.Int Quo and Go's built-in signed integer division both truncate
    // toward zero. Rounding negative outcomes to floor would be an error.
    return new(big.Int).Quo(contrast.Num(),contrast.Denom()).Int64(), true
}

// Independent conventional control receives exactly the original parent's
// proposed action and the same seven per-action selected history counters.
func dirichletPick(parent int, rows Actions)(int,bool){
    if parent<0 || parent>=7 {return parent,false}
    baseline,ok:=dirichletScore(rows[parent])
    if !ok {return parent,false}
    best:=baseline
    selected:=parent
    found:=false
    for a,s:=range rows {
        if s.Changed<2 {continue}
        signal,valid:=dirichletScore(s)
        if !valid {return parent,false}
        if signal>best{
            selected=a
            best=signal
            found=true
        }
    }
    return selected,found
}

// Native scoring is the frozen source implementation. This separate helper
// projects its selection rule without altering that source or its parent.
func sourceScoringPick(parent int, rows Actions)(int,bool){
    if parent<0||parent>=7{return parent,false}
    best:=score(rows[parent])
    selected:=parent
    found:=false
    for action,s:=range rows{
        if s.Changed<2{continue}
        v:=score(s)
        if v>best{selected=action;best=v;found=true}
    }
    return selected,found
}

func validStates(limit int) []Stat {
    var states []Stat
    for n:=0;n<=limit;n++{
        for changed:=0;changed<=n;changed++{
            for noEffect:=0;noEffect<=n-changed;noEffect++{
                states=append(states,Stat{
                    N:uint32(n),Changed:uint32(changed),NoEffect:uint32(noEffect),
                })
            }
        }
    }
    return states
}

func TestNRM1WExactDirichletScoreEquivalence(t *testing.T){
    states:=validStates(8)
    if len(states)!=165{t.Fatalf("unexpected preregistered state count: %d",len(states))}
    for _,s:=range states{
        expected,valid:=dirichletScore(s)
        if !valid || score(s)!=expected{
            t.Fatalf("true native score is not Dirichlet posterior: %+v native=%d baseline=%d valid=%t",s,score(s),expected,valid)
        }
    }
    for _,s:=range []Stat{
        {N:4096},{N:4096,Changed:4096},{N:4096,NoEffect:4096},
        {N:4096,Changed:1,NoEffect:4095},
        {N:4096,Changed:2048,NoEffect:2048},
        {N:4095,Changed:2048,NoEffect:2047},
        {N:4095,Changed:0,NoEffect:1},
    }{
        expected,valid:=dirichletScore(s)
        if !valid||expected!=score(s){t.Fatalf("cap-score discrepancy at %+v",s)}
    }
    for _,bad:=range []Stat{
        {N:1,Changed:2},{N:2,Changed:2,NoEffect:2},
        {N:4097},{N:4096,NoEffect:4097},
    } {
        _,valid:=dirichletScore(bad)
        if valid{t.Fatalf("Dirichlet control accepted invalid selected source: %+v",bad)}
    }
}

// Exhaustive supported pair states at N<=6, with all seven parent proposals
// and each of six distinct competing actions. N=0..6 has 84 valid three-class
// states; 7*6*84*84 = 296352 source-only policy comparisons.
func TestNRM1WMatchedConventionalPolicyPairExhaustive(t *testing.T){
    small:=validStates(6)
    if len(small)!=84{t.Fatalf("incorrect source-state set: %d",len(small))}
    var cases,overrides int
    for parent:=0;parent<7;parent++ {
        for alternative:=0;alternative<7;alternative++ {
            if alternative==parent{continue}
            for _,p:=range small{
                for _,q:=range small{
                    var rows Actions
                    rows[parent]=p
                    rows[alternative]=q
                    actual,actualChanged:=sourceScoringPick(parent,rows)
                    baseline,baseChanged:=dirichletPick(parent,rows)
                    if actual!=baseline||actualChanged!=baseChanged{
                        t.Fatalf("conventional decision disagrees: p=%d alt=%d pstat=%+v qstat=%+v native=%d,%t control=%d,%t",
                            parent,alternative,p,q,actual,actualChanged,baseline,baseChanged)
                    }
                    cases++
                    if actualChanged{overrides++}
                }
            }
        }
    }
    if cases!=296352||overrides==0{
        t.Fatalf("wrong bounded coverage or no positive proposed changes: cases=%d overrides=%d",cases,overrides)
    }
    t.Logf("NRM1W_SOURCE_ONLY_MATCHED_DIRICHLET_PAIR_PASS cases=%d overrides=%d actual_world_steps=0",cases,overrides)
}

// Exercise the REAL native Choose with its unchanged selected-source policy,
// including the shared local context on an unseen full public observation.
func TestNRM1WActualNativeChoiceMatchesMatchedDirichletControl(t *testing.T){
    m:=New()
    mission:="use the key to open the door and then get to the goal"
    before:=scene(mission,0,1,1)
    after:=scene(mission,0,2,1)
    other:=scene(mission,0,1,4)
    otherAfter:=scene(mission,0,2,4)
    for i:=0;i<12;i++ {
        if e:=m.ObserveFit(other,0,otherAfter,0,false);e!=nil{t.Fatal(e)}
    }
    for i:=0;i<2;i++ {
        if e:=m.ObserveFit(before,5,after,0,false);e!=nil{t.Fatal(e)}
    }
    for _,probe:=range []Public{
        before, scene(mission,2,7,1),scene(mission,0,3,4),
        scene(mission,0,1,2),
    } {
        native,_,err:=m.Choose(probe)
        parent,_,perr:=m.Policy.Choose(probe)
        if err!=nil||perr!=nil||native==nil||parent==nil {
            t.Fatalf("actual parent/native decisions unavailable: %v / %v",err,perr)
        }
        key,e:=context(probe)
        if e!=nil{t.Fatal(e)}
        expected,_:=dirichletPick(*parent,m.Local[key])
        if *native!=expected {
            t.Fatalf("actual native choice differs from equal-information Dirichlet, parent=%d native=%d control=%d",*parent,*native,expected)
        }
    }
}
