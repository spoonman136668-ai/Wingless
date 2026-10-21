package cognition

// P6O exact source/custody. Original tool-bound/native action consequences are
// observed only AFTER execution. Competing hypothesis CLASSES are explicitly
// authored by the experimenter, NOT autonomously invented by the learner.
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

func p6oConventionalBeta(a,b float64)float64{
    x,_:=math.Lgamma(a)
    y,_:=math.Lgamma(b)
    z,_:=math.Lgamma(a+b)
    return x+y-z
}
func p6oConventionalEntropy(p float64)float64{
    if p<=0||p>=1{return 0}
    return -p*math.Log(p)-(1-p)*math.Log(1-p)
}
// Independent fully trained Bayesian conventional comparator. Its own
// counter state has exactly the same actually executed outcome history,
// independent of the native learner's representation and predictor code.
func p6oIndependentBayesianController(c p6mCountControl,motors [2]float32)(float32,float64,[2]float64){
    changed:=float64(c.success[0]+c.success[1])
    unchanged:=float64(c.n[0]+c.n[1]-c.success[0]-c.success[1])
    baselineP:=(changed+1)/(changed+unchanged+2)
    logBF:=-p6oConventionalBeta(changed+1,unchanged+1)
    for i:=0;i<2;i++{
        yes:=float64(c.success[i])
        no:=float64(c.n[i]-c.success[i])
        logBF+=p6oConventionalBeta(yes+1,no+1)
    }
    w:=1/(1+math.Exp(-logBF))
    var scores [2]float64
    chosen:=motors[0]
    best:=-1.
    for i,motor:=range motors{
        j:=p6mMotorClass(motor)
        conditional:=(float64(c.success[j])+1)/(float64(c.n[j])+2)
        mixture:=(1-w)*baselineP+w*conditional
        info:=p6oConventionalEntropy(mixture)-
            (1-w)*p6oConventionalEntropy(baselineP)-
            w*p6oConventionalEntropy(conditional)
        scores[i]=info
        if info>best+1e-12 {chosen=motor;best=info}
    }
    return chosen,logBF,scores
}

func TestNCI04P4P6OOriginalNativeHypothesisIntervention(t *testing.T){
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    motors:=[2]float32{p6mWriteMotor,p6mReadMotor}
    var mem0,mem1 runtime.MemStats
    runtime.ReadMemStats(&mem0)
    state,independentHistory,originalWall,learningWall:=p6mLearnNativeMotorActions(t)
    frozen,err:=state.Snapshot()
    if err!=nil{t.Fatal(err)}
    start:=time.Now()
    hypotheses,err:=state.NativeEffectHypothesisExperiment(pre,motors)
    candidateDecisionWall:=time.Since(start)
    if err!=nil || !hypotheses.Informative {t.Fatal("native competing-hypothesis evidence unavailable",err)}
    if hypotheses.Changed!=[2]uint32{2,0} || hypotheses.Unchanged!=[2]uint32{0,2}{
        t.Fatal("learned hypothesis evidence does not match real original native actions")
    }
    if !(hypotheses.PosteriorConditioned>0.5) ||
       !(hypotheses.ExpectedInformationNats[0]>0) ||
       !(hypotheses.ExpectedInformationNats[1]>0){
        t.Fatal("original native experience gave no legitimate competing-hypothesis information")
    }
    start=time.Now()
    controlMotor,controlBF,controlScores:=p6oIndependentBayesianController(independentHistory,motors)
    conventionalDecisionWall:=time.Since(start)
    if math.Abs(hypotheses.LogEvidenceRatio-controlBF)>1e-11 ||
       math.Abs(hypotheses.ExpectedInformationNats[0]-controlScores[0])>1e-11 ||
       math.Abs(hypotheses.ExpectedInformationNats[1]-controlScores[1])>1e-11 ||
       hypotheses.ChosenMotor!=controlMotor{
        t.Fatal("independently trained equal-information conventional Bayesian controller disagrees")
    }
    // A separate conventional least-visited native motor policy also ties
    // on 2-versus-2 native observations; compare decisions explicitly.
    countControl:=p6nConventionalActiveCounts{
        Motors:motors,
        Visits:[2]uint32{independentHistory.n[1],independentHistory.n[0]},
    }
    if countControl.Pick()!=hypotheses.ChosenMotor {
        t.Fatal("native hypothesis choice differs from equally informed visit-count tie control")
    }
    afterPrediction,err:=state.Snapshot()
    if err!=nil||!bytes.Equal(afterPrediction,frozen){
        t.Fatal("counterfactual native planning consumed or mutated future outcome")
    }
    // Native action selected before its outcome; a REAL original call produces
    // the next event. The learned system may only update after that call.
    post,changed,nativeNextWall:=p6mActualMotorExecution(t,hypotheses.ChosenMotor)
    if changed!=(hypotheses.ChosenMotor==p6mWriteMotor) {
        t.Fatal("source-authenticated original actor contradicted recorded command")
    }
    start=time.Now()
    if err:=state.ObserveNativeMotorTransition(pre,hypotheses.ChosenMotor,post);err!=nil{
        t.Fatal("failed to retain actual prospectively chosen native result",err)
    }
    nextLearningWall:=time.Since(start)
    independentHistory.observe(hypotheses.ChosenMotor,changed)
    next,err:=state.NativeEffectHypothesisExperiment(pre,motors)
    if err!=nil || next.LogEvidenceRatio<=hypotheses.LogEvidenceRatio {
        t.Fatal("genuine new native evidence did not revise causal hypothesis support")
    }
    nextControl,_,_:=p6oIndependentBayesianController(independentHistory,motors)
    if nextControl!=next.ChosenMotor {t.Fatal("native posterior ceased to match trained Bayesian control")}
    runtime.ReadMemStats(&mem1)
    snap,err:=state.Snapshot();if err!=nil{t.Fatal(err)}
    if unsafe.Sizeof(state)!=6472 {t.Fatal("native learner capacity expanded")}
    receipt:=map[string]any{
        "gate":"NCI-04P6O",
        "source_native_head_commit":p6iGit(t,"..","rev-parse","HEAD"),
        "hypothesis_classes":[]string{"H0_motors_share_one_consequence_rate","H1_effect_rate_is_motor_conditioned"},
        "hypothesis_classes_manually_specified":true,
        "prior_pseudocounts":"Beta_1_1_each",
        "source_authenticated_native_training_observations":4,
        "candidate_posterior_H1_before_next":hypotheses.PosteriorConditioned,
        "candidate_log_bayes_factor_before_next":hypotheses.LogEvidenceRatio,
        "candidate_expected_information_nats":hypotheses.ExpectedInformationNats,
        "candidate_motor_preaction":hypotheses.ChosenMotor,
        "independent_conventional_bayesian_motor":controlMotor,
        "independent_conventional_bayesian_scores":controlScores,
        "candidate_action_advantage_over_bayesian":false,
        "candidate_action_advantage_over_visit_count":false,
        "prospective_original_native_effect_changed_state":changed,
        "posterior_log_bayes_factor_after_new_native_experience":next.LogEvidenceRatio,
        "native_action_effect_learned_after_execution_only":true,
        "native_fixed_state_bytes":unsafe.Sizeof(state),
        "candidate_native_learned_snapshot_bytes":len(snap),
        "native_initial_training_wall_ns":originalWall.Nanoseconds(),
        "native_initial_learning_wall_ns":learningWall.Nanoseconds(),
        "candidate_hypothesis_decision_wall_ns":candidateDecisionWall.Nanoseconds(),
        "conventional_hypothesis_decision_wall_ns":conventionalDecisionWall.Nanoseconds(),
        "selected_native_original_action_wall_ns":nativeNextWall.Nanoseconds(),
        "selected_native_learning_wall_ns":nextLearningWall.Nanoseconds(),
        "sampled_parent_heap_initial_bytes":mem0.Alloc,
        "sampled_parent_heap_final_bytes":mem1.Alloc,
        "all_in_cpu_rss_io":"NOT_QUALIFIED",
        "new_independent_unfamiliar_source_families":0,
        "scientific_first_use":false,
        "classification":"ORIGINAL_NATIVE_BAYES_HYPOTHESIS_COMPARISON_PASS__CONVENTIONAL_EQUIVALENT",
    }
    b,err:=json.Marshal(receipt);if err!=nil{t.Fatal(err)}
    t.Log("P6O_NATIVE_HYPOTHESIS_RECEIPT="+string(b))
}

type p6oColdReceipt struct {
    Schema string `json:"schema"`
    SourceHead string `json:"source_head"`
    ActuatorContract string `json:"actuator_contract"`
    State json.RawMessage `json:"state"`
}
func TestNCI04P4P6OIndependentProcessHypothesisRestart(t *testing.T){
    p6fPinOriginalSource(t)
    root,err:=filepath.Abs("..");if err!=nil{t.Fatal(err)}
    head:=p6iGit(t,root,"rev-parse","HEAD")
    pre:=p6lOriginalCognitionPreaction(t)
    motors:=[2]float32{p6mWriteMotor,p6mReadMotor}
    if os.Getenv("P6O_COLD_CHILD")=="1"{
        blob,err:=os.ReadFile(os.Getenv("P6O_SEALED_STATE"))
        if err!=nil{t.Fatal(err)}
        var sealed p6oColdReceipt
        if err:=json.Unmarshal(blob,&sealed);err!=nil{t.Fatal(err)}
        if sealed.Schema!="wingless.nci04p6o.hypotheses.v1" ||
           sealed.SourceHead!=head || sealed.ActuatorContract!=p6mContractHash() {
            t.Fatal("cold-process original actuator or source identity drift")
        }
        var restored LatentState
        if err:=restored.Restore(sealed.State);err!=nil{t.Fatal(err)}
        h,err:=restored.NativeEffectHypothesisExperiment(pre,motors)
        if err!=nil || !h.Informative || h.Changed!=[2]uint32{2,0} ||
           h.Unchanged!=[2]uint32{0,2} ||
           !(h.PosteriorConditioned>0.5) {
            t.Fatal("native hypothesis evidence was lost across genuine process restart",err)
        }
        t.Log("P6O_INDEPENDENT_OS_PROCESS_HYPOTHESIS_RESTORE_PASS")
        return
    }
    state,_,_,_:=p6mLearnNativeMotorActions(t)
    snap,err:=state.Snapshot();if err!=nil{t.Fatal(err)}
    sealed:=p6oColdReceipt{
        Schema:"wingless.nci04p6o.hypotheses.v1",
        SourceHead:head,ActuatorContract:p6mContractHash(),State:snap,
    }
    b,err:=json.Marshal(sealed);if err!=nil{t.Fatal(err)}
    path:=filepath.Join(t.TempDir(),"p6o-hypothesis-sealed.json")
    if err:=os.WriteFile(path,b,0600);err!=nil{t.Fatal(err)}
    cmd:=exec.Command(os.Args[0],"-test.run=^TestNCI04P4P6OIndependentProcessHypothesisRestart$","-test.v")
    cmd.Env=append(os.Environ(),"P6O_COLD_CHILD=1","P6O_SEALED_STATE="+path)
    output,err:=cmd.CombinedOutput()
    if err!=nil || !bytes.Contains(output,[]byte("P6O_INDEPENDENT_OS_PROCESS_HYPOTHESIS_RESTORE_PASS")){
        t.Fatalf("actual independent native process lost hypothesis evidence: %v %s",err,output)
    }
    unchanged,err:=os.ReadFile(path)
    if err!=nil||!bytes.Equal(unchanged,b){t.Fatal("native cold restart modified original evidence",err)}
    t.Log("P6O_NATIVE_HYPOTHESIS_COLD_OS_RESTART_PASS",len(b))
}

func TestNCI04P4P6OInsufficientNativeEvidenceAbstains(t *testing.T){
    p6fPinOriginalSource(t)
    pre:=p6lOriginalCognitionPreaction(t)
    var state LatentState
    noEvidence,err:=state.NativeEffectHypothesisExperiment(pre,[2]float32{p6mWriteMotor,p6mReadMotor})
    if err!=nil||noEvidence.Informative {t.Fatal("untrained native effect hypotheses falsely claimed information",err)}
    if _,err:=state.NativeEffectHypothesisExperiment(pre,[2]float32{p6mWriteMotor,p6mWriteMotor});err==nil {
        t.Fatal("duplicate motor hypotheses received an experiment choice")
    }
    if _,err:=state.NativeEffectHypothesisExperiment(pre,[2]float32{p6mWriteMotor,float32(math.NaN())});err==nil {
        t.Fatal("nonfinite motor affordance was admitted")
    }
    original,err:=state.Snapshot();if err!=nil{t.Fatal(err)}
    after,err:=state.Snapshot()
    if err!=nil||!bytes.Equal(after,original){t.Fatal("invalid causal experiments changed learned state",err)}
}
