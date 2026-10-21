package cognition

// P6O is a source-independent candidate-only competition between two
// explicitly defined hypotheses on already executed native motor effects:
// H0: the probability of native-state change is motor-independent.
// H1: each supplied motor has a distinct learned change probability.
// The hypothesis CLASSES are conventional and hand-specified; ONLY their
// evidence and predictions are learned from existing native experience.
// This is neither self-organized hypothesis invention nor scientific transfer.
import (
    "errors"
    "math"
)

type NativeEffectHypothesis struct {
    Motors [2]float32
    Changed [2]uint32
    Unchanged [2]uint32
    SharedProbability float64
    ConditionalProbability [2]float64
    PosteriorConditioned float64
    LogEvidenceRatio float64
    ExpectedInformationNats [2]float64
    ChosenMotor float32
    Informative bool
}

func p6oLogBeta(a,b float64)float64 {
    x,_:=math.Lgamma(a)
    y,_:=math.Lgamma(b)
    z,_:=math.Lgamma(a+b)
    return x+y-z
}
func p6oEntropy(p float64)float64 {
    if p<=0||p>=1{return 0}
    return -p*math.Log(p)-(1-p)*math.Log(1-p)
}
func p6oLogistic(v float64)float64 {
    if v>=0 {return 1/(1+math.Exp(-v))}
    q:=math.Exp(v)
    return q/(1+q)
}

// p6oNativeEffects reads ONLY genuine learned native transitions in the
// original fixed LatentState relation graph. Effects come from an observed
// post-execution state-delta vector. No source/function names or future result.
func (s *LatentState) p6oNativeEffects(before [FeatureWidth]float32,motor float32)(changed,unchanged uint32,err error) {
    if s==nil{return 0,0,errors.New("nil native effect memory")}
    conditioned,err:=NativeMotorPreactionConditions(before,motor)
    if err!=nil{return 0,0,err}
    idx,d:=nearest(s,conditioned)
    if idx<0||d>0.0625||s.ambiguousNeighbor(conditioned,idx,d){
        return 0,0,nil
    }
    for j:=0;j<int(s.Count);j++ {
        n:=uint32(s.Relations[idx][j][0].Updates)
        if n==0{continue}
        post:=s.Entities[j].Mean
        if post[2]!=1 || post[3]!=0 {return 0,0,errors.New("unverified native postaction sensory modality")}
        if post[0]>0.875 {changed+=n
        } else if post[0]< -0.875 {unchanged+=n
        } else {return 0,0,errors.New("native action consequence not identifiable")}
    }
    return changed,unchanged,nil
}

// NativeEffectHypothesisExperiment estimates Bayesian evidence for H0/H1 from
// actual fixed-budget native learned history, then chooses the motor with the
// greatest one-step expected mutual information between hypothesis identity
// and the NEXT unknown original native result. This is a standard parametric
// Beta-Bernoulli computation, with the original motor set externally given.
// It abstains if no genuine evidence supports an informative distinction.
func (s *LatentState) NativeEffectHypothesisExperiment(before [FeatureWidth]float32,motors [2]float32)(NativeEffectHypothesis,error){
    var h NativeEffectHypothesis
    if motors[0]==motors[1]{return h,errors.New("identical native action commands")}
    h.Motors=motors
    var changed,unchanged uint32
    for i,motor:=range motors {
        y,n,err:=s.p6oNativeEffects(before,motor)
        if err!=nil{return NativeEffectHypothesis{},err}
        h.Changed[i],h.Unchanged[i]=y,n
        changed+=y
        unchanged+=n
        h.ConditionalProbability[i]=float64(y+1)/float64(y+n+2)
    }
    total:=changed+unchanged
    if total==0 {return h,nil}
    h.SharedProbability=float64(changed+1)/float64(total+2)
    h.LogEvidenceRatio= -p6oLogBeta(float64(changed+1),float64(unchanged+1))
    for i:=0;i<2;i++ {
        h.LogEvidenceRatio+=p6oLogBeta(float64(h.Changed[i]+1),float64(h.Unchanged[i]+1))
    }
    h.PosteriorConditioned=p6oLogistic(h.LogEvidenceRatio)
    w:=h.PosteriorConditioned
    best:=-1.0
    for i,p1:=range h.ConditionalProbability {
        p0:=h.SharedProbability
        mix:=w*p1+(1-w)*p0
        score:=p6oEntropy(mix)-w*p6oEntropy(p1)-(1-w)*p6oEntropy(p0)
        if score<0 && score> -1e-12 {score=0}
        if math.IsNaN(score)||math.IsInf(score,0)||score<0 {return NativeEffectHypothesis{},errors.New("invalid native hypothesis information score")}
        h.ExpectedInformationNats[i]=score
        if score>best {
            h.ChosenMotor=motors[i]
            best=score
        }
    }
    if best<=1e-12 {
        h.Informative=false
        return h,nil
    }
    h.Informative=true
    return h,nil
}
