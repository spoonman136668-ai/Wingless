package cognition

// P6M original-native actor/motor memory test. Every effect is measured from
// the original State snapshots AFTER native execution; command selection is
// logged beforehand. No user-facing action, broker or protected science.
import (
    "bytes"
    "crypto/sha256"
    "encoding/hex"
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

const (
    p6mReadMotor float32 = -0.75
    p6mWriteMotor float32 = 0.75
    // The actual motor/actuator contract is a source-authenticated custody
    // boundary. It is never embedded in feature vectors or learned as a label.
    p6mActorContract = "p6m.v1: -0.75=original-LatentState.Predict; +0.75=original-LatentState.Observe"
)
func p6mContractHash()string{
    h:=sha256.Sum256([]byte(p6mActorContract))
    return hex.EncodeToString(h[:])
}
func p6mActualMotorExecution(t *testing.T,motor float32)([FeatureWidth]float32,bool,time.Duration){
    t.Helper()
    var original LatentState
    before,err:=original.Snapshot();if err!=nil{t.Fatal(err)}
    src:=[FeatureWidth]float32{-0.6,-0.3,0,0}
    dst:=[FeatureWidth]float32{0.6,0.3,0,0}
    start:=time.Now()
    switch motor {
    case p6mWriteMotor:
        err=original.Observe(src,dst,0,1)
    case p6mReadMotor:
        _,err=original.Predict(src,dst,0)
    default:
        t.Fatal("actuator refused unregistered motor command before original execution")
    }
    elapsed:=time.Since(start)
    if err!=nil{t.Fatal("original native motor execution failed",err)}
    after,err:=original.Snapshot();if err!=nil{t.Fatal(err)}
    changed:=!bytes.Equal(before,after)
    output:=[FeatureWidth]float32{-1,0,1,0}
    if changed {output[0]=1}
    return output,changed,elapsed
}

type p6mCountControl struct {
    n [2]uint32
    positives [2]uint32
}
func p6mMotorClass(m float32)int{
    if m<0{return 0}
    return 1
}
func (c *p6mCountControl)observe(m float32,changed bool){
    i:=p6mMotorClass(m)
    c.n[i]++
    if changed{c.positives[i]++}
}
func (c *p6mCountControl)predict(m float32)float64{
    i:=p6mMotorClass(m)
    return float64(c.positives[i]+1)/float64(c.n[i]+2)
}

func p6mLearnNativeMotorActions(t *testing.T)(LatentState,p6mCountControl,time.Duration,time.Duration){
    t.Helper()
    pre:=p6lOriginalCognitionPreaction(t)
    var state LatentState
    var conventional p6mCountControl
    var nativeWall,learnerWall time.Duration
    for _,motor:=range []float32{p6mReadMotor,p6mWriteMotor,p6mReadMotor,p6mWriteMotor}{
        // Efference copy exists BEFORE the action. No future outcome is passed.
        if _,known,err:=state.PredictNativeMotorTransition(pre,motor);err!=nil{
            t.Fatal("preaction motor interface rejected valid command",err)
        }else if state.Observations==0 && known{t.Fatal("untrained motor state was confident")}
        observed,changed,wall:=p6mActualMotorExecution(t,motor)
        nativeWall+=wall
        start:=time.Now()
        if err:=state.ObserveNativeMotorTransition(pre,motor,observed);err!=nil{
            t.Fatal("real native motor outcome not learned",err)
        }
        learnerWall+=time.Since(start)
        conventional.observe(motor,changed)
    }
    if state.Observations!=4 || unsafe.Sizeof(state)!=6472{
        t.Fatal("not exactly four native experiences or fixed capacity violated")
    }
    return state,conventional,nativeWall,learnerWall
}

func TestNCI04P4P6MNativeMotorEfferenceAndParity(t *testing.T){
    p6fPinOriginalSource(t)
    var memBefore,memAfter runtime.MemStats
    runtime.ReadMemStats(&memBefore)
    trained,control,originalWall,updateWall:=p6mLearnNativeMotorActions(t)
    pre:=p6lOriginalCognitionPreaction(t)
    var learnerLoss,countLoss,sourceInspectorLoss float64
    var predictWall,controlWall time.Duration
    for _,motor:=range []float32{p6mWriteMotor,p6mReadMotor}{
        start:=time.Now()
        prediction,known,err:=trained.PredictNativeMotorTransition(pre,motor)
        predictWall+=time.Since(start)
        if err!=nil || !known{t.Fatal("trained native motor lost commanded-effect evidence",err)}
        start=time.Now()
        countForecast:=control.predict(motor)
        controlWall+=time.Since(start)
        actualPost,changed,_:=p6mActualMotorExecution(t,motor)
        if prediction!=actualPost{t.Fatal("learned native motor transition mismatched actual execution")}
        expected:=0.0
        if changed{expected=1}
        learnedProbability:=0.0
        if prediction[0]>0{learnedProbability=1}
        learnerLoss+=(learnedProbability-expected)*(learnedProbability-expected)
        countLoss+=(countForecast-expected)*(countForecast-expected)
        // Separate conventional actuator-contract inspector sees exactly
        // the same command BEFORE execution and knows the original call
        // chosen by the dispatch. It therefore predicts the same effect.
        sourceForecast:=0.0
        if motor==p6mWriteMotor{sourceForecast=1}
        sourceInspectorLoss+=(sourceForecast-expected)*(sourceForecast-expected)
    }
    runtime.ReadMemStats(&memAfter)
    snapshot,err:=trained.Snapshot();if err!=nil{t.Fatal(err)}
    report:=map[string]any{
        "gate":"NCI-04P6M",
        "original_native_training_actions":4,
        "original_native_dev_validation_actions":2,
        "action_input":"PRE_EXECUTION_ACTUATOR_EFFERENCE_COPY_NOT_RESULT_LABEL",
        "motor_effect_sign_learned_from_original_native_state_deltas":true,
        "source_or_method_name_in_learner_vector":false,
        "genuine_unfamiliar_task_source_family":false,
        "fixed_latent_state_bytes":unsafe.Sizeof(trained),
        "snapshot_bytes":len(snapshot),
        "learner_brier_sum":learnerLoss,
        "fully_trained_conditional_count_brier_sum":countLoss,
        "fully_informed_original_actuator_inspector_brier_sum":sourceInspectorLoss,
        "learner_advantage_over_best_control":learnerLoss<sourceInspectorLoss,
        "native_action_wall_ns":originalWall.Nanoseconds(),
        "learner_update_wall_ns":updateWall.Nanoseconds(),
        "learner_predict_wall_ns":predictWall.Nanoseconds(),
        "conditional_count_predict_wall_ns":controlWall.Nanoseconds(),
        "parent_heap_before_bytes":memBefore.Alloc,
        "parent_heap_after_bytes":memAfter.Alloc,
        "matched_all_in_cpu_rss_io":"NOT_QUALIFIED",
        "scientific_first_use":false,
        "classification":"REAL_NATIVE_MOTOR_ASSOCIATION_LEARNED_NO_SCIENTIFIC_CONTROL_ADVANTAGE",
    }
    b,err:=json.Marshal(report);if err!=nil{t.Fatal(err)}
    t.Log("P6M_NATIVE_MOTOR_EVIDENCE="+string(b))
    if learnerLoss!=0 || sourceInspectorLoss!=0 || learnerLoss<sourceInspectorLoss{
        t.Fatal("source-native motor result falsely claimed conventional superiority")
    }
}

type p6mSealedState struct {
    Schema string `json:"schema"`
    Head string `json:"head"`
    ActuatorContractSHA string `json:"actuator_contract_sha256"`
    State json.RawMessage `json:"state"`
}
func p6mCheckSeal(sealed p6mSealedState,head string)error {
    if sealed.Schema!="wingless.nci04p6m.motor.v1" ||
       sealed.Head!=head || sealed.ActuatorContractSHA!=p6mContractHash() {
        return os.ErrPermission
    }
    return nil
}
func TestNCI04P4P6MTrueColdProcessMotorRestore(t *testing.T){
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
    head:=p6iGit(t,root,"rev-parse","HEAD")
    if os.Getenv("P6M_CHILD")=="1"{
        b,err:=os.ReadFile(os.Getenv("P6M_SEAL"))
        if err!=nil{t.Fatal(err)}
        var seal p6mSealedState
        if err:=json.Unmarshal(b,&seal);err!=nil{t.Fatal(err)}
        if err:=p6mCheckSeal(seal,head);err!=nil{t.Fatal("actuator/source drift on cold restart",err)}
        var trained LatentState
        if err:=trained.Restore(seal.State);err!=nil{t.Fatal(err)}
        for _,caseMotor:=range []struct{motor float32;change bool}{
            {p6mReadMotor,false},{p6mWriteMotor,true},
        } {
            result,known,err:=trained.PredictNativeMotorTransition(pre,caseMotor.motor)
            if err!=nil||!known||(result[0]>0)!=caseMotor.change{
                t.Fatal("cold-process original motor semantic recall failed",err)
            }
        }
        t.Log("P6M_FRESH_PROCESS_NATIVE_MOTOR_RECALL_PASS")
        return
    }
    state,_,_,_:=p6mLearnNativeMotorActions(t)
    snapshot,err:=state.Snapshot();if err!=nil{t.Fatal(err)}
    seal:=p6mSealedState{Schema:"wingless.nci04p6m.motor.v1",Head:head,ActuatorContractSHA:p6mContractHash(),State:snapshot}
    if err:=p6mCheckSeal(seal,head);err!=nil{t.Fatal(err)}
    tampered:=seal
    tampered.ActuatorContractSHA="not-the-original-actuator-map"
    if p6mCheckSeal(tampered,head)==nil{t.Fatal("changed native motor actuator map was silently admitted")}
    b,err:=json.Marshal(seal);if err!=nil{t.Fatal(err)}
    path:=filepath.Join(t.TempDir(),"original-native-motor-seal.json")
    if err:=os.WriteFile(path,b,0600);err!=nil{t.Fatal(err)}
    cmd:=exec.Command(os.Args[0],"-test.run=^TestNCI04P4P6MTrueColdProcessMotorRestore$","-test.v")
    cmd.Env=append(os.Environ(),"P6M_CHILD=1","P6M_SEAL="+path)
    out,err:=cmd.CombinedOutput()
    if err!=nil || !bytes.Contains(out,[]byte("P6M_FRESH_PROCESS_NATIVE_MOTOR_RECALL_PASS")){
        t.Fatalf("separate native process cannot restore source-bound motor state: %v %s",err,out)
    }
    unchanged,err:=os.ReadFile(path)
    if err!=nil || !bytes.Equal(unchanged,b){t.Fatal("native motor cold restart rewrote durable custody",err)}
    t.Log("P6M_REAL_COLD_PROCESS_PASS sealed_bytes=",len(b))
}
func TestNCI04P4P6MAmbiguousActuationAndInvalidMotorFailClosed(t *testing.T){
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    var trained LatentState
    zero,err:=trained.Snapshot();if err!=nil{t.Fatal(err)}
    positive,_:=p6lExecuteOriginalMemoryAction(t,true)
    negative,_:=p6lExecuteOriginalMemoryAction(t,false)
    if err:=trained.ObserveNativeMotorTransition(pre,p6mWriteMotor,positive);err!=nil{t.Fatal(err)}
    if err:=trained.ObserveNativeMotorTransition(pre,p6mWriteMotor,negative);err!=nil{t.Fatal(err)}
    answer,known,err:=trained.PredictNativeMotorTransition(pre,p6mWriteMotor)
    if err!=nil||known||answer!=([FeatureWidth]float32{}){
        t.Fatal("conflicting original motor effects did not trigger abstention",err)
    }
    stable,err:=trained.Snapshot();if err!=nil{t.Fatal(err)}
    if err:=trained.ObserveNativeMotorTransition(pre,float32(math.NaN()),positive);err==nil{
        t.Fatal("nonfinite motor command was admitted")
    }
    if _,_,err:=trained.PredictNativeMotorTransition(pre,1.5);err==nil{
        t.Fatal("out-of-budget motor command predicted")
    }
    after,err:=trained.Snapshot()
    if err!=nil||!bytes.Equal(stable,after){t.Fatal("invalid native motor command mutated learned state",err)}
    if bytes.Equal(zero,stable){t.Fatal("genuine native motor effects were not retained")}
}
