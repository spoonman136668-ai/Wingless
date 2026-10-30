package mg23

import "testing"

// NR-M1X source-free falsification. No original MiniGrid environment is run.
// We compare the unchanged native MG23 source-selected API with an independent
// tabular 32-step eligibility trace given the EXACT same selected events.
func TestNRM1XTerminalCreditEqualsTabularEligibilityAndIncludesDistractors(t *testing.T) {
    model := New()
    before := frame("pick up the yellow ball", 0, [3]int{1, 0, 0})
    key, err := context(before)
    if err != nil { t.Fatal(err) }
    if err = model.BeginEpisode(); err != nil { t.Fatal(err) }
    // First selected action is outside the terminal-credit window.
    // The next 31 actions are zero-effect turns, then a terminal feedback.
    // This is deliberately an abstract selected-feedback source fixture, NOT
    // a claim about physically achievable original-world trajectories.
    selected := []int{2}
    for i:=0; i<31; i++ {selected=append(selected,0)}
    selected=append(selected,6)
    conventional:=map[int]uint32{}
    trace:=[]int{}
    for i,action:=range selected {
        positive:=i==len(selected)-1
        reward:=0.0
        if positive { reward=1.0 }
        if err=model.ObserveFit(before,action,before,reward,positive);err!=nil{
            t.Fatalf("native selected source event %d failed: %v",i,err)
        }
        trace=append(trace,action)
        if len(trace)>Window { trace=trace[1:] }
        if positive {
            for age:=0;age<len(trace);age++{
                a:=trace[len(trace)-1-age]
                conventional[a]+=uint32(Window-age)
            }
        }
    }
    if len(model.Trace)!=Window {t.Fatalf("wrong strict terminal window length %d",len(model.Trace))}
    row:=model.Cells[key]
    for action:=0;action<7;action++{
        if row[action].Delayed!=conventional[action]{
            t.Fatalf("native MG23 differs from conventional same-info trace action %d native=%d expected=%d",
                action,row[action].Delayed,conventional[action])
        }
        if row[action].Target!=0 {t.Fatalf("nonexistent source-public target observation credited to %d",action)}
    }
    if row[2].Count!=1||row[2].Delayed!=0{
        t.Fatalf("out-of-window selected action received terminal credit %+v",row[2])
    }
    if row[0].Count!=31||row[0].Delayed!=496{
        t.Fatalf("31 unchanged public-view distractor turns failed to receive full terminal credit %+v",row[0])
    }
    if row[6].Count!=1||row[6].Delayed!=32 {t.Fatalf("actual terminal selected feedback missing %+v",row[6])}
    observed,reason,e:=model.Choose(before)
    if e!=nil||observed==nil||*observed!=0||
        reason!="SELECT_NATIVE_MG23_OWN_OBSERVED_TARGET_OR_TERMINAL_CREDIT"{
        t.Fatalf("high terminal credit failed to favor arbitrary zero-effect distractor: %v %s %v",observed,reason,e)
    }
    cold,e:=model.Snapshot();if e!=nil{t.Fatal(e)}
    restored,e:=Restore(cold);if e!=nil{t.Fatal(e)}
    if restored.Cells[key][0].Delayed!=496{t.Fatal("terminal distractor credit lost after cold restore")}
    noGoal:=New()
    if e=noGoal.BeginEpisode();e!=nil{t.Fatal(e)}
    for _,action:=range selected{
        if e=noGoal.ObserveFit(before,action,before,0,false);e!=nil{t.Fatal(e)}
    }
    for _,s:=range noGoal.Cells[key]{
        if s.Delayed!=0 {t.Fatal("reward-free source hallucinated delayed credit")}
    }
    t.Log("NRM1X_WINGLESS_SELECTED_SOURCE_CONVENTIONAL_TRACE_PARITY_PASS 33 events, 31 zero-effect distractors, zero original-world steps")
}
