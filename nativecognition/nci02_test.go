package nativecognition

// NCI02 is a source-free developmental falsification screen, not official science.
// Edge slots and sign composition are supplied by the harness, NOT learned.
import (
    "encoding/json"
    "fmt"
    "os"
    "reflect"
    "testing"
)

type nci02Obs struct { Channel int; Value float64 }
type nci02Episode struct { Signs [4]int; Events [12]nci02Obs }

type nci02RNG struct { state uint32 }
func (r *nci02RNG) next() uint32 {
    x:=r.state
    x^=x<<13;x^=x>>17;x^=x<<5
    r.state=x
    return x
}
func nci02Generate(seed uint32) nci02Episode {
    r:=nci02RNG{state:seed+0x9e3779b9}
    var e nci02Episode
    for i:=range e.Signs {
        if r.next()&1==1 {e.Signs[i]=1}else{e.Signs[i]=-1}
    }
    for pass:=0;pass<3;pass++ {
        for ch:=0;ch<4;ch++ {
            v:=float64(e.Signs[ch])*.8
            if r.next()%10==0 {v=-v}
            e.Events[pass*4+ch]=nci02Obs{ch,v}
        }
    }
    return e
}
func nci02Truth(e nci02Episode,length int) int {
    p:=1
    for i:=0;i<length;i++ {p*=e.Signs[i]}
    return p
}
func nci02Pred(v float64) int {
    if v>0 {return 1}
    if v<0 {return -1}
    return 0
}
func nci02Native(rate float64,episode nci02Episode,length int,clear,interfere bool) (int,error) {
    m,e:=New(rate);if e!=nil{return 0,e}
    for _,o:=range episode.Events {if _,err:=m.Observe(o.Channel,o.Value);err!=nil {return 0,err}}
    if clear {m,_=New(rate)}
    if interfere {
        other:=nci02Generate(uint32(0x80000000)+uint32(episode.Signs[0]+2)*97+uint32(episode.Signs[1]+2)*59)
        for _,o:=range other.Events {if _,err:=m.Observe(o.Channel,o.Value);err!=nil{return 0,err}}
    }
    p:=1
    for i:=0;i<length;i++{
        v,err:=m.Predict(i);if err!=nil{return 0,err}
        q:=nci02Pred(v);if q==0{return 0,nil}
        p*=q
    }
    return p,nil
}
type nci02EMA struct {
    alpha float64
    state [4]float64
}
func (m *nci02EMA) observe(o nci02Obs) {m.state[o.Channel]+=(o.Value-m.state[o.Channel])*m.alpha}
func (m *nci02EMA) predict(length int) int {
    p:=1
    for i:=0;i<length;i++ {
        x:=nci02Pred(m.state[i]);if x==0{return 0};p*=x
    }
    return p
}
func nci02Baseline(alpha float64,episode nci02Episode,length int,interfere bool) int {
    m:=nci02EMA{alpha:alpha}
    for _,o:=range episode.Events {m.observe(o)}
    if interfere {
        other:=nci02Generate(uint32(0x80000000)+uint32(episode.Signs[0]+2)*97+uint32(episode.Signs[1]+2)*59)
        for _,o:=range other.Events {m.observe(o)}
    }
    return m.predict(length)
}
func nci02TrainSelect(candidate bool) (float64,error) {
    rates:=[]float64{.05,.1,.2}
    winner:=rates[0];scoreBest:=-1
    for _,rate:=range rates {
        correct:=0
        for seed:=uint32(6000);seed<=6031;seed++ {
            e:=nci02Generate(seed)
            var got int
            if candidate {var err error;got,err=nci02Native(rate,e,2,false,false);if err!=nil{return 0,err}}
            if !candidate {got=nci02Baseline(rate,e,2,false)}
            if got==nci02Truth(e,2){correct++}
        }
        if correct>scoreBest {scoreBest=correct;winner=rate}
    }
    return winner,nil
}
func nci02Score(rate,alpha float64,start uint32,length int) (map[string]float64,error) {
    var learned,baseline,frozen,ablated,interfered int
    for seed:=start;seed<start+32;seed++ {
        ep:=nci02Generate(seed)
        answer:=nci02Truth(ep,length)
        p,e:=nci02Native(rate,ep,length,false,false);if e!=nil{return nil,e}
        if p==answer{learned++}
        if nci02Baseline(alpha,ep,length,false)==answer{baseline++}
        f,e:=nci02Native(0,ep,length,false,false);if e!=nil{return nil,e}
        if f==answer{frozen++}
        a,e:=nci02Native(rate,ep,length,true,false);if e!=nil{return nil,e}
        if a==answer{ablated++}
        x,e:=nci02Native(rate,ep,length,false,true);if e!=nil{return nil,e}
        if x==answer{interfered++}
    }
    return map[string]float64{
       "candidate":float64(learned)/32,
       "trained_ema":float64(baseline)/32,
       "frozen_native":float64(frozen)/32,
       "memory_reset":float64(ablated)/32,
       "after_interference":float64(interfered)/32,
       "exact_symbolic_upper_bound":1,
    },nil
}
func nci02Report() (map[string]any,error) {
    r,e:=nci02TrainSelect(true);if e!=nil{return nil,e}
    a,e:=nci02TrainSelect(false);if e!=nil{return nil,e}
    development,e:=nci02Score(r,a,6100,2);if e!=nil{return nil,e}
    stress,e:=nci02Score(r,a,6200,3);if e!=nil{return nil,e}
    gain:=stress["candidate"]-stress["trained_ema"]
    lesionLoss:=stress["candidate"]-stress["memory_reset"]
    interferenceDrop:=stress["candidate"]-stress["after_interference"]
    passed:=gain>=.05 && stress["candidate"]>=.75 && lesionLoss>=.10 && interferenceDrop<=.10
    disposition:="NEGATIVE"
    if passed {disposition="PASS_PROXY_ONLY"}
    return map[string]any{
       "schema":"nci02.wingless.sourcefree-development-result.v1",
       "result":disposition,"candidate_rate_train_selected":r,"ema_rate_train_selected":a,
       "development_two_hop":development,"unseen_three_hop_proxy":stress,
       "gain_vs_trained_ema":gain,"memory_reset_loss":lesionLoss,
       "interference_drop":interferenceDrop,"train_seed_range":[]int{6000,6031},
       "development_seed_range":[]int{6100,6131},"stress_seed_range":[]int{6200,6231},
       "candidate_state_channels":4,"symbolic_composition_supplied":true,
       "entity_binding_learned":false,"native_composition_learned":false,
       "protected_scientific_first_use":0,"learned_transfer_accepted":false,
       "training_episodes_per_arm":32,"observations_per_episode":12,
    },nil
}
func TestNCI02CausalRetentionAndMatchedTrainedControl(t *testing.T) {
    report,err:=nci02Report();if err!=nil {t.Fatal(err)}
    repeat,err:=nci02Report();if err!=nil {t.Fatal(err)}
    if !reflect.DeepEqual(report,repeat){t.Fatal("nondeterministic NCI02 development report")}
    out,err:=json.MarshalIndent(report,"","  ");if err!=nil{t.Fatal(err)}
    fmt.Printf("NCI02_WINGLESS_RESULT=%s\n",string(out))
    if path:=os.Getenv("NCI02_RESULT_PATH");path!="" {
        if err:=os.WriteFile(path,append(out,'\n'),0600);err!=nil{t.Fatal(err)}
    }
    result:=report["result"]
    if result!="PASS_PROXY_ONLY" && result!="NEGATIVE" {t.Fatalf("unexpected classification %v",result)}
}
func TestNCI02GeneratorIsolation(t *testing.T) {
    if reflect.DeepEqual(nci02Generate(6000),nci02Generate(6200)){t.Fatal("training/eval source collision")}
    ep:=nci02Generate(6200)
    if len(ep.Events)!=12 {t.Fatal("observation budget drift")}
    for _,o:=range ep.Events {
        if o.Channel<0||o.Channel>=4||o.Value!=.8 && o.Value!=-.8{t.Fatal("invalid event")}
    }
}
