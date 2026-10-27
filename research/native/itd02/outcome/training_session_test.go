package outcome

import (
    "bytes"
    "errors"
    "testing"

    r1 "github.com/spoonman136668-ai/Wingless/research/native/itd02/r1"
)

// Deliberately public-shaped engineering source. This is NOT original-world
// evidence or a substitute for an authorized MiniGrid development adapter.
type sessionFixture struct {
    initial r1.PublicGrid
    preferred int
    resets int
    actions []int
    failStep bool
    terminal bool
}
func (f *sessionFixture) ResetPublic()(r1.PublicGrid,error){
    f.resets++
    return f.initial,nil
}
func (f *sessionFixture) StepPublic(action int)(PublicConsequence,error){
    f.actions=append(f.actions,action)
    if f.failStep{return PublicConsequence{},errors.New("SOURCE_STEP_OUTCOME_AMBIGUOUS")}
    after:=f.initial
    switch action {
    case 0: after.Direction=(after.Direction+3)%4
    case 1: after.Direction=(after.Direction+1)%4
    }
    reward:=float64(1)
    if action==f.preferred{reward=9}
    return PublicConsequence{
        After:after,Reward:reward,Terminal:f.terminal,
    },nil
}

func TestNativeSelfTrainingEndToEndColdRestart(t *testing.T){
    before:=selectedFixture()
    native:=NewSelected()
    baseline,err:=native.Policy.Choose(before)
    if err!=nil{t.Fatal(err)}
    preferred:=(baseline+1)%3
    source:=&sessionFixture{initial:before,preferred:preferred,terminal:true}
    model:=NewSelected()
    var persisted []byte
    for i:=0;i<6;i++{
        training:=&TrainingSession{Learner:model,
            MinimumSelectedSamples:2,MaximumEpisodeSteps:1}
        result,state,err:=training.RunEpisode(source)
        if err!=nil{t.Fatalf("episode %d: %v",i,err)}
        if !result.Terminal||result.Steps!=1||
            result.NativeUpdates!=uint64(i+1){t.Fatal("not a complete native episode")}
        if len(state)==0{t.Fatal("no native checkpoint")}
        persisted=state
        // Simulated independent native process: NO Go state carried forward
        // except the exact cumulative checkpoint returned by the learner.
        model,err=RestoreSelected(persisted)
        if err!=nil{t.Fatal(err)}
    }
    if source.resets!=6 || len(source.actions)!=6{
        t.Fatal("not six actually source-stepped decisions")
    }
    for i,b:=range model.Rewards.Bins{
        if b.Count!=2{t.Fatalf("action %d selected %d times",i,b.Count)}
    }
    chosen,err:=model.RecommendRewardAction(2)
    if err!=nil||chosen!=preferred{t.Fatalf("did not learn reward advantage %d, want %d: %v",chosen,preferred,err)}
    last,_,err:=(&TrainingSession{Learner:model,MinimumSelectedSamples:2,
        MaximumEpisodeSteps:1}).RunEpisode(source)
    if err!=nil || last.CumulativeReward!=9{
        t.Fatalf("resumed learner did not improve chosen reward %+v %v",last,err)
    }
    if source.actions[len(source.actions)-1]!=preferred ||
        source.actions[len(source.actions)-1]==baseline{
        t.Fatal("learned action did not change after self training")
    }
}

func TestNativeTrainingRejectsAmbiguousSourceStepWithoutFabrication(t *testing.T){
    s:=NewSelected()
    source:=&sessionFixture{initial:selectedFixture(),failStep:true,terminal:true}
    _,_,err:=(&TrainingSession{Learner:s,MinimumSelectedSamples:2,
        MaximumEpisodeSteps:1}).RunEpisode(source)
    if err==nil{t.Fatal("ambiguous source step passed")}
    if !s.Policy.Armed||s.Policy.LiveSeen!=0||s.Rewards.NextID!=1 {
        t.Fatal("ambiguous outcome consumed or discarded pending action")
    }
    if _,_,err=(&TrainingSession{Learner:s,
        MinimumSelectedSamples:2,MaximumEpisodeSteps:1}).RunEpisode(source);err==nil{
        t.Fatal("replayed ambiguous source action")
    }
    if len(source.actions)!=1{t.Fatal("ambiguous source was stepped twice")}
}

func TestNativeTrainingNeverForgesTerminal(t *testing.T){
    s:=NewSelected()
    source:=&sessionFixture{initial:selectedFixture(),terminal:false}
    _,_,err:=(&TrainingSession{Learner:s,MinimumSelectedSamples:2,
        MaximumEpisodeSteps:1}).RunEpisode(source)
    if err==nil{t.Fatal("fabricated source terminal")}
    if s.Rewards.Closed||s.Rewards.NextID!=2 || s.Policy.LiveSeen!=1 {
        t.Fatal("lost real nonterminal native evidence")
    }
    if _,_,err:=(&TrainingSession{Learner:s,MinimumSelectedSamples:2,
        MaximumEpisodeSteps:1}).RunEpisode(source);err==nil{
        t.Fatal("restarted unfinished source episode")
    }
}

func TestNativeTrainingCheckpointSinkIsBoundedAndFailClosed(t *testing.T){
    source:=&sessionFixture{initial:selectedFixture(),preferred:1,terminal:true}
    s:=&TrainingSession{Learner:NewSelected(),MinimumSelectedSamples:2,
        MaximumEpisodeSteps:1}
    var saved [][]byte
    result,err:=s.TrainBounded(source,3,func(data []byte)error{
        saved=append(saved,append([]byte(nil),data...))
        return nil
    })
    if err!=nil||len(result)!=3||len(saved)!=3{t.Fatalf("%v %+v",err,result)}
    for i,raw:=range saved {
        cold,err:=RestoreSelected(raw)
        if err!=nil || cold.Policy.LiveSeen!=uint64(i+1){
            t.Fatalf("saved cold checkpoint %d invalid %v",i,err)
        }
    }
    if bytes.Equal(saved[0],saved[2]) {t.Fatal("checkpoints not cumulative")}
    n:=source.resets
    _,err=s.TrainBounded(source,2,func(data []byte)error{return errors.New("CHECKPOINT_WRITE_UNKNOWN")})
    if err==nil||source.resets!=n+1 {t.Fatal("continued after ambiguous checkpoint sink")}
    if _,err:=s.TrainBounded(source,129,func([]byte)error{return nil});err==nil{
        t.Fatal("accepted unbounded training")}
}
