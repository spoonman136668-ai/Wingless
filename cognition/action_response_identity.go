package cognition

// NCI04E is an isolated, source-free necessary-condition probe.
// A trial supplies only ordered (action,outcome) observations and a terminal
// boundary. No latent identity, entity index, future answer or context label
// is passed to the learner. A mode is formed only on mature predictive surprise.
import ("errors"; "math")

const NCI04EMaxModes = 2

type responseMode struct {
    Used bool
    Zero [2]uint32
    One [2]uint32
}
type ActionResponseIdentity struct {
    Modes [NCI04EMaxModes]responseMode
    Active int
}
func NewActionResponseIdentity() *ActionResponseIdentity {
    return &ActionResponseIdentity{Active:-1}
}
func (m *ActionResponseIdentity) BeginTrial() error {
    if m==nil {return errors.New("nil identity model")}
    m.Active=-1
    return nil
}
func (m *ActionResponseIdentity) ModeCount() int {
    if m==nil{return 0}
    n:=0;for _,mode:=range m.Modes{if mode.Used{n++}}
    return n
}
func modeEstimate(v responseMode,a int) float64 {
    return float64(v.One[a]+1)/float64(v.Zero[a]+v.One[a]+2)
}
func (m *ActionResponseIdentity) Observe(a,outcome int)error {
    if m==nil||a<0||a>1||outcome<0||outcome>1{return errors.New("invalid action/outcome")}
    if m.Active < -1||m.Active>=NCI04EMaxModes{return errors.New("invalid active context")}
    chosen:=m.Active
    if chosen<0 {
        best:=-1
        bestLikelihood:=-1.0
        allMature:=true
        for i,v:=range m.Modes {
            if !v.Used {continue}
            n:=v.Zero[a]+v.One[a]
            if n<3 {allMature=false}
            p:=modeEstimate(v,a)
            likelihood:=p
            if outcome==0 {likelihood=1-p}
            if likelihood>bestLikelihood{bestLikelihood=likelihood;best=i}
        }
        if best<0 {chosen=0;m.Modes[0].Used=true
        }else if bestLikelihood<.25 && allMature && m.ModeCount()<NCI04EMaxModes {
            chosen=-1
            for i,v:=range m.Modes{if !v.Used{chosen=i;break}}
            if chosen<0{return errors.New("mode allocation failed")}
            m.Modes[chosen]=responseMode{Used:true}
        }else{chosen=best}
        m.Active=chosen
    }
    v:=&m.Modes[chosen]
    if v.Zero[a]+v.One[a]>=1000000{return errors.New("observation budget exhausted")}
    if outcome==0{v.Zero[a]++}else{v.One[a]++}
    return nil
}
// Predict cannot infer an unobserved current identity. It returns unknown
// until this trial supplies distinguishing evidence and the action is trained.
func (m *ActionResponseIdentity) Predict(a int)(float64,bool,error){
    if m==nil||a<0||a>1{return 0,false,errors.New("invalid query")}
    if m.Active==-1{return 0,false,nil}
    v:=m.Modes[m.Active]
    if !v.Used ||v.One[a]+v.Zero[a]<3{return 0,false,nil}
    return modeEstimate(v,a),true,nil
}
func (m *ActionResponseIdentity) CurrentMode() (int,bool) {
    if m==nil ||m.Active<0{return -1,false}
    return m.Active,true
}
// WeightedModeConfusion is deliberately not an intervention or identity oracle.
// It reports only conditional response separation learned from experience.
func (m *ActionResponseIdentity) Separation(a int)(float64,bool){
    if m==nil||a<0||a>1||m.ModeCount()!=2{return 0,false}
    for _,v:=range m.Modes{if v.Zero[a]+v.One[a]<3{return 0,false}}
    return math.Abs(modeEstimate(m.Modes[0],a)-modeEstimate(m.Modes[1],a)),true
}
