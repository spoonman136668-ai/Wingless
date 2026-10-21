package cognition

// P6M is an opt-in, candidate-only motor-efference learning interface.
// The motor value is the actuator's actually selected command signal,
// captured BEFORE it executes. It is NOT an effect/answer label, source ID,
// claim of novelty or grant of action execution authority.
//
// Only existing 6472-byte LatentState storage is used; there is no action
// registry, new trainable parameter, predictor oracle or production dispatch.
import (
    "errors"
    "math"
)

// NativeMotorPreactionConditions augments the actual authorized source
// observation with the efference copy of the command already selected by
// its original caller. Caller must independently source-authenticate and
// bind its actuator contract for cold restart; this helper cannot prove
// which operation a controller will actually execute.
func NativeMotorPreactionConditions(before [FeatureWidth]float32, motor float32) ([FeatureWidth]float32,error) {
    if !validVector(before) || before[2]!=-1 || before[3]!=0 ||
       math.IsNaN(float64(motor)) || math.IsInf(float64(motor),0) ||
       motor < -1 || motor > 1 {
        return [FeatureWidth]float32{},errors.New("unverified motor preaction signal")
    }
    before[3]=motor
    return before,nil
}

// ObserveNativeMotorTransition learns a motor-conditioned, perceptual
// state transition ONLY AFTER its actual consequence is measured. Distinct
// motor signals dynamically occupy existing learned source prototypes;
// all relations use one fixed existing channel, not an action-ID lookup.
func (s *LatentState) ObserveNativeMotorTransition(before [FeatureWidth]float32,motor float32,after [FeatureWidth]float32)error{
    if s==nil{return errors.New("nil native motor learner")}
    conditioned,err:=NativeMotorPreactionConditions(before,motor)
    if err!=nil{return err}
    return s.ObserveNativeTransition(conditioned,after,0)
}

// PredictNativeMotorTransition never receives the post-action outcome.
// If no matching experience exists or distinct observed consequences
// conflict, the existing P6L rule abstains.
func (s *LatentState) PredictNativeMotorTransition(before [FeatureWidth]float32,motor float32)([FeatureWidth]float32,bool,error){
    var unknown [FeatureWidth]float32
    if s==nil{return unknown,false,errors.New("nil native motor learner")}
    conditioned,err:=NativeMotorPreactionConditions(before,motor)
    if err!=nil{return unknown,false,err}
    return s.PredictNativeTransition(conditioned,0)
}
