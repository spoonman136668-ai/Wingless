package cognition

// P6H is an opt-in native engineering learner for perceptual transitions.
// It changes no existing LatentState representation or production executor.
// Preaction and observed consequence are deliberately different modalities.
// This mechanism is not scientific transfer or an accepted source semantic model.
import (
    "errors"
    "math"

    tb "github.com/spoonman136668-ai/Wingless/toolboundary"
)

// NativePreactionFeatures uses available permitted workspace measurements;
// no path, source hash, target verdict, or post-action result enters prediction.
// Numeric size is bounded. The third coordinate marks the observable modality,
// not an identity, causality label, or expected answer.
func NativePreactionFeatures(v tb.PreactionObservation) ([FeatureWidth]float32,error) {
    var x [FeatureWidth]float32
    x[2]=-1
    switch v.Resource {
    case "missing":
        if v.Size!=0 || v.SHA256!="" {return x,errors.New("malformed missing preaction resource")}
        x[0]=-1; x[1]=0
    case "regular":
        if v.Size<0 || v.Size>1048576 || len(v.SHA256)!=64 {return x,errors.New("malformed regular preaction resource")}
        x[0]=1
        x[1]=float32(math.Min(1,float64(v.Size)/64))
    default:
        return x,errors.New("unverified preaction resource")
    }
    return x,nil
}

// NativeExecutedConsequenceFeatures is derived strictly AFTER native Execute.
// The success/failure and returned output size are measured consequences,
// not target labels available to the pre-action prediction.
func NativeExecutedConsequenceFeatures(v tb.Observation,err error) [FeatureWidth]float32 {
    var x [FeatureWidth]float32
    x[2]=1
    if err!=nil || v.ExitCode!=0 {
        x[0]=-1
        return x
    }
    x[0]=1
    x[1]=float32(math.Min(1,float64(len(v.Stdout))/64))
    return x
}

// ObserveNativeTransition learns a previously experienced preaction-to-result
// relation. The constant '1' denotes that an observed transition occurred;
// native success/failure is represented ONLY in its observed postaction vector.
// No new trainable fields, event registry, or expanded fixed state.
func (s *LatentState) ObserveNativeTransition(before,after [FeatureWidth]float32,action int) error {
    if s==nil{return errors.New("nil latent native transition state")}
    if before[2]!=-1 || after[2]!=1{return errors.New("native pre/post modality violation")}
    return s.Observe(before,after,action,1)
}

// PredictNativeTransition infers the most supported next-state prototype from
// learned transitions without receiving the destination or actual outcome.
// Unseen/ambiguous sources, missing links, or tied evidence must abstain.
func (s *LatentState) PredictNativeTransition(before [FeatureWidth]float32,action int)([FeatureWidth]float32,bool,error) {
    var empty [FeatureWidth]float32
    if s==nil{return empty,false,errors.New("nil latent native transition state")}
    if !validVector(before) || before[2]!=-1 || action<0 || action>=MaxActions {
        return empty,false,errors.New("invalid preaction prediction")
    }
    src,d:=nearest(s,before)
    if src<0 || d>0.0625 || s.ambiguousNeighbor(before,src,d){return empty,false,nil}
    best:=-1
    // Distinct supported native consequences from indistinguishable public
    // preaction information are not resolvable by choosing a majority.
    // A plurality here can confidently predict the wrong real native action.
    // Reject without adding any state fields, thresholds, or action identities.
    for j:=0;j<int(s.Count);j++ {
        if s.Entities[j].Mean[2]!=1 || s.Relations[src][j][action].Updates==0 {continue}
        if best>=0{return empty,false,nil}
        best=j
    }
    if best<0{return empty,false,nil}
    return s.Entities[best].Mean,true,nil
}
