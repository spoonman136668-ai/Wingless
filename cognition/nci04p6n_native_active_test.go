package cognition

// P6N exact-source, real-native closed-loop intervention test. The experimental
// actions are the ORIGINAL LatentState.Observe and LatentState.Predict. Only the
// source-authenticated actuator dispatch maps motor signals to these methods;
// no expected effect, function name or future verdict enters native selection.
import (
    "bytes"
    "encoding/json"
    "math"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "testing"
    "time"
    "unsafe"
)

type p6nConventionalActiveCounts struct {
    Motors [2]float32
    Visits [2]uint32
}
func (c *p6nConventionalActiveCounts)Pick()float32 {
    if c.Visits[0]<=c.Visits[1] {return c.Motors[0]}
    return c.Motors[1]
}
func (c *p6nConventionalActiveCounts)Observe(motor float32) {
    if c.Motors[0]==motor {c.Visits[0]++;return}
    if c.Motors[1]==motor {c.Visits[1]++;return}
    panic("unregistered original actuator command in control")
}

func TestNCI04P4P6NAuthenticActiveNativeExperimentVersusTrainedCounts(t *testing.T) {
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    candidates:=[]float32{p6mWriteMotor,p6mReadMotor}
    control:=p6nConventionalActiveCounts{Motors:[2]float32{p6mWriteMotor,p6mReadMotor}}
    var active LatentState
    var nativeWall,selectWall,learnWall,controlWall time.Duration
    var before,after runtime.MemStats
    runtime.ReadMemStats(&before)
    var selected []float32
    var actual []bool
    var scores []float64
    for round:=0;round<6;round++ {
        selectionStart:=time.Now()
        proposed,err:=active.ChooseNativeMotorExperiment(pre,candidates)
        selectWall+=time.Since(selectionStart)
        if err!=nil{t.Fatal("native candidate action selection failed",err)}
        if proposed.ExplorationScore<=0||proposed.ExplorationScore>1 {
            t.Fatal("invalid epistemic exploration proxy")
        }
        conventionalStart:=time.Now()
        conventional:=control.Pick()
        controlWall+=time.Since(conventionalStart)
        if conventional!=proposed.Motor {
            t.Fatal("same-information independently executed conventional count policy differs")
        }
        // The selected motor is passed to the original actuator, and the
        // native state delta is measured only after executing that method.
        observed,changed,wall:=p6mActualMotorExecution(t,proposed.Motor)
        nativeWall+=wall
        learnStart:=time.Now()
        if err:=active.ObserveNativeMotorTransition(pre,proposed.Motor,observed);err!=nil {
            t.Fatal("failed to learn original native action consequence",err)
        }
        learnWall+=time.Since(learnStart)
        control.Observe(conventional)
        selected=append(selected,proposed.Motor)
        actual=append(actual,changed)
        scores=append(scores,proposed.ExplorationScore)
    }
    runtime.ReadMemStats(&after)
    if active.Observations!=6 || control.Visits!=[2]uint32{3,3} {
        t.Fatal("active native experiments did not preserve equal information/coverage")
    }
    if unsafe.Sizeof(active)!=6472 {t.Fatal("fixed native learned state changed capacity")}
    for i,motor:=range selected {
        if (motor==p6mWriteMotor)!=actual[i] {
            t.Fatal("source-bound native actuator effect mismatch")
        }
        if motor!=candidates[i%2] {t.Fatal("learner did not reduce visit-count uncertainty as expected")}
    }
    snapshot,err:=active.Snapshot();if err!=nil{t.Fatal(err)}
    frozen:=append([]byte(nil),snapshot...)
    for _,motor:=range candidates {
        p,known,err:=active.PredictNativeMotorTransition(pre,motor)
        if err!=nil||!known||((p[0]>0)!=(motor==p6mWriteMotor)){
            t.Fatal("closed loop failed to retain actual motor-conditioned native consequences",err)
        }
    }
    afterPrediction,err:=active.Snapshot()
    if err!=nil||!bytes.Equal(frozen,afterPrediction){
        t.Fatal("native active experiment predicted by mutating learned state",err)
    }
    receipt:=map[string]any{
        "gate":"NCI-04P6N",
        "native_original_methods":[]string{"LatentState.Observe","LatentState.Predict"},
        "source_authenticated_native_actuator_contract_sha256":p6mContractHash(),
        "authentic_adaptively_chosen_executions":6,
        "candidate_action_selection_before_outcome":true,
        "each_native_action_chosen":3,
        "original_native_effect_measured_only_after_execution":true,
        "candidate_motor_sequence":selected,
        "observed_native_changes":actual,
        "exploration_scores":scores,
        "equal_information_count_policy_sequence_identical":true,
        "candidate_advantage_over_count_policy":false,
        "information_metric":"inverse_existing_relation_visits_NOT_SHANNON_INFORMATION_GAIN",
        "learned_native_effect_memory_after_experiment":true,
        "learner_fixed_state_bytes":unsafe.Sizeof(active),
        "learner_snapshot_bytes":len(frozen),
        "sampled_native_action_wall_ns":nativeWall.Nanoseconds(),
        "sampled_candidate_selection_wall_ns":selectWall.Nanoseconds(),
        "sampled_candidate_update_wall_ns":learnWall.Nanoseconds(),
        "sampled_control_selection_wall_ns":controlWall.Nanoseconds(),
        "parent_heap_before_bytes":before.Alloc,
        "parent_heap_after_bytes":after.Alloc,
        "full_parent_child_cpu_rss_io":"NOT_QUALIFIED",
        "unseen_action_source_families":0,
        "scientific_first_use":false,
        "classification":"AUTHENTIC_ACTIVE_NATIVE_ACTIONS_BUT_COUNT_BASELINE_EQUIVALENT",
    }
    b,err:=json.Marshal(receipt)
    if err!=nil{t.Fatal(err)}
    t.Log("P6N_NATIVE_ACTIVE_EXPERIMENT_RECEIPT="+string(b))
}

type p6nColdWire struct {
    Schema string `json:"schema"`
    Head string `json:"head"`
    Contract string `json:"contract_sha256"`
    State json.RawMessage `json:"state"`
}
func TestNCI04P4P6NTrueColdProcessActionSelection(t *testing.T) {
    p6fPinOriginalSource(t)
    root,err:=filepath.Abs("..")
    if err!=nil{t.Fatal(err)}
    head:=p6iGit(t,root,"rev-parse","HEAD")
    pre:=p6lOriginalCognitionPreaction(t)
    candidates:=[]float32{p6mWriteMotor,p6mReadMotor}
    if os.Getenv("P6N_NEW_PROCESS")=="1" {
        b,err:=os.ReadFile(os.Getenv("P6N_SEAL"))
        if err!=nil{t.Fatal(err)}
        var wire p6nColdWire
        if err:=json.Unmarshal(b,&wire);err!=nil{t.Fatal(err)}
        if wire.Schema!="wingless.nci04p6n.active.v1" ||
           wire.Head!=head||wire.Contract!=p6mContractHash() {
            t.Fatal("source-bound actuator authority drift on cold restart")
        }
        var restored LatentState
        if err:=restored.Restore(wire.State);err!=nil{t.Fatal(err)}
        next,err:=restored.ChooseNativeMotorExperiment(pre,candidates)
        if err!=nil||next.Motor!=p6mWriteMotor||next.Observations!=2 {
            t.Fatal("separate OS process lost native intervention experience",err)
        }
        t.Log("P6N_COLD_PROCESS_ACTION_SELECTION_PASS")
        return
    }
    // Persist exact finite native experience, with motor-pair visit counts
    // 2 and 2, so the next original intervention is source-bound deterministic.
    state,_,_,_:=p6mLearnNativeMotorActions(t)
    snap,err:=state.Snapshot();if err!=nil{t.Fatal(err)}
    wire:=p6nColdWire{
        Schema:"wingless.nci04p6n.active.v1",Head:head,
        Contract:p6mContractHash(),State:snap,
    }
    b,err:=json.Marshal(wire);if err!=nil{t.Fatal(err)}
    path:=filepath.Join(t.TempDir(),"sealed-native-active.json")
    if err:=os.WriteFile(path,b,0600);err!=nil{t.Fatal(err)}
    cmd:=exec.Command(os.Args[0],"-test.run=^TestNCI04P4P6NTrueColdProcessActionSelection$","-test.v")
    cmd.Env=append(os.Environ(),"P6N_NEW_PROCESS=1","P6N_SEAL="+path)
    out,err:=cmd.CombinedOutput()
    if err!=nil||!bytes.Contains(out,[]byte("P6N_COLD_PROCESS_ACTION_SELECTION_PASS")){
        t.Fatalf("real separate-process active-memory restoration failed: %v %s",err,out)
    }
    copied,err:=os.ReadFile(path)
    if err!=nil||!bytes.Equal(copied,b){t.Fatal("native active cold replay modified sealed evidence",err)}
    t.Log("P6N_SEPARATE_OS_PROCESS_COLD_RECALL_PASS")
}

func TestNCI04P4P6NInvalidInterventionsFailClosed(t *testing.T) {
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    var state LatentState
    initial,err:=state.Snapshot()
    if err!=nil{t.Fatal(err)}
    invalid:=[][]float32{
        nil,{p6mWriteMotor},{p6mWriteMotor,p6mWriteMotor},
        {p6mWriteMotor,float32(math.NaN())},
        {p6mWriteMotor,1.5},
    }
    for _,bad:=range invalid {
        if _,err:=state.ChooseNativeMotorExperiment(pre,bad);err==nil {
            t.Fatalf("invalid motor experiment admitted: %v",bad)
        }
    }
    after,err:=state.Snapshot()
    if err!=nil||!bytes.Equal(initial,after){t.Fatal("invalid preaction motor mutated learned state",err)}
}
