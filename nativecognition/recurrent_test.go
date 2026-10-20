package nativecognition

import (
    "bytes"
    "math"
    "testing"
)

func TestNativeOnlineAdaptationAndFrozenControl(t *testing.T) {
    learned,_:=New(.2)
    frozen,_:=New(0)
    var early,late,frozenLoss float64
    for i:=0;i<120;i++ {
        a,err:=learned.Observe(0,.8);if err!=nil { t.Fatal(err) }
        b,err:=frozen.Observe(0,.8);if err!=nil { t.Fatal(err) }
        e:=(.8-a)*(.8-a)
        if i<12 {early+=e}
        if i>=108 {late+=e}
        frozenLoss+=(.8-b)*(.8-b)
    }
    if !(late<early && late<frozenLoss/10) {
        t.Fatalf("no bounded online adaptation: early=%v late=%v frozen=%v",early,late,frozenLoss)
    }
    if learned.Cells[0].Readout==0 { t.Fatal("weights did not learn") }
    if learned.Cells[1].State!=0 { t.Fatal("unused channel mutated") }
}

func TestQueryBlindReadAndDurableSnapshotIsolation(t *testing.T) {
    m,_:=New(.15)
    for i:=0;i<18;i++ {if _,err:=m.Observe(2,.6);err!=nil {t.Fatal(err)}}
    snap,err:=m.Checkpoint();if err!=nil {t.Fatal(err)}
    for i:=0;i<20;i++ {_,err:=m.Predict(2);if err!=nil {t.Fatal(err)}}
    again,_:=m.Checkpoint()
    if !bytes.Equal(snap,again) { t.Fatal("query mutated memory") }
    restored,err:=FromCheckpoint(snap);if err!=nil {t.Fatal(err)}
    if _,err:=restored.Observe(2,-.5);err!=nil {t.Fatal(err)}
    if bytes.Equal(snap,mustCheckpoint(t,restored)) { t.Fatal("restored copy did not evolve") }
    original,_:=m.Checkpoint();if !bytes.Equal(original,snap) {t.Fatal("snapshot aliased original")}
}

func mustCheckpoint(t *testing.T,m *Memory) []byte {
    t.Helper();b,e:=m.Checkpoint();if e!=nil {t.Fatal(e)};return b
}

func TestInvalidInputsFailWithoutMutation(t *testing.T) {
    m,_:=New(.1);before:=mustCheckpoint(t,m)
    bad:=[]struct{ch int;v float64}{{-1,.1},{4,.1},{0,math.NaN()},{1,math.Inf(1)},{1,1.001}}
    for _,x:=range bad {if _,err:=m.Observe(x.ch,x.v);err==nil {t.Fatalf("accepted invalid %+v",x)}}
    if !bytes.Equal(before,mustCheckpoint(t,m)) {t.Fatal("invalid inputs changed state")}
    for _,raw:=range [][]byte{[]byte("{}"),[]byte(`{"version":2,"memory":{}}`),append(append([]byte{},before...),[]byte("{}")...)} {
        if _,err:=FromCheckpoint(raw);err==nil {t.Fatal("accepted invalid checkpoint")}
    }
    if _,err:=New(-.1);err==nil {t.Fatal("invalid rate")}
    if _,err:=New(math.Inf(1));err==nil {t.Fatal("nonfinite rate")}
}
