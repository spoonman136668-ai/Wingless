package cognition

// P6N is a candidate-only, bounded preaction native motor-selection policy.
// It acts on actual prior motor-conditioned learned relation evidence. It does
// not discover action primitives, authorize a tool, infer new actuator meaning,
// or prove causal/scientific advantage over a matched count policy.
import (
    "errors"
    "math"
)

// NativeMotorChoice describes the experimenter's supplied candidate motor,
// selected before its original native actuator executes. An external native
// actor must bind the actual command and uphold authorization and source custody.
type NativeMotorChoice struct {
    Motor float32
    ExplorationScore float64
    Observations uint32
}

// ChooseNativeMotorExperiment chooses the LEAST observed of a bounded list of
// existing motor commands, based on relation evidence in fixed LatentState.
// Source and motor may not disclose native future effects; no outcome labels.
// It is structurally equivalent to conditional action-visit counts and should
// not be characterized as an information-theoretically novel selection policy.
func (s *LatentState) ChooseNativeMotorExperiment(before [FeatureWidth]float32,motors []float32)(NativeMotorChoice,error) {
    var choice NativeMotorChoice
    if s==nil || len(motors)<2 || len(motors)>MaxActions {
        return choice,errors.New("invalid bounded native motor choice")
    }
    var seen [MaxActions]float32
    for k,motor:=range motors {
        if math.IsNaN(float64(motor)) || math.IsInf(float64(motor),0) {
            return choice,errors.New("nonfinite motor command")
        }
        for j:=0;j<k;j++ {
            if seen[j]==motor{return choice,errors.New("duplicated motor command")}
        }
        seen[k]=motor
        if _,err:=NativeMotorPreactionConditions(before,motor);err!=nil {
            return choice,err
        }
    }
    best:=uint32(^uint32(0))
    for _,motor:=range motors {
        src,err:=NativeMotorPreactionConditions(before,motor)
        if err!=nil{return NativeMotorChoice{},err}
        idx,d:=nearest(s,src)
        visits:=uint32(0)
        if idx>=0 && d<=0.0625 && !s.ambiguousNeighbor(src,idx,d) {
            for j:=0;j<int(s.Count);j++ {
                if s.Entities[j].Mean[2]!=1 {continue}
                visits+=uint32(s.Relations[idx][j][0].Updates)
            }
        }
        if visits<best {
            best=visits
            choice.Motor=motor
            choice.Observations=visits
            choice.ExplorationScore=1.0/float64(1+visits)
        }
    }
    return choice,nil
}
