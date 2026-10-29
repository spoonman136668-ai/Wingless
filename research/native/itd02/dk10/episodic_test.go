package dk10

import (
 "bytes"
 "testing"
 dk7 "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7"
)

func sourcePublic(direction int) dk7.Public {
 p:=dk7.Public{Image:make([][][]int,7),Direction:direction,Mission:dk7.Mission}
 for x:=0;x<7;x++{p.Image[x]=make([][]int,7);for y:=0;y<7;y++{p.Image[x][y]=[]int{1,0,0}}}
 p.Image[3][5]=[]int{5,4,0}
 return p
}

func TestSelectedPublicTransitionNoveltyAndColdCheckpoint(t *testing.T){
 m:=New()
 a,b:=sourcePublic(0),sourcePublic(1)
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 selected,status,e:=m.Choose(a)
 if e!=nil||selected==nil||status!="SELECT_NATIVE_BASE_EXPLORATION"{t.Fatalf("no base native proposal: %v %s %v",selected,status,e)}
 if e=m.Commit(a,0);e!=nil{t.Fatal(e)}
 if e=m.Feedback(a,1,a,0,false);e==nil{t.Fatal("wrong selected action accepted")}
 if m.Samples!=0||m.Base.Native.Samples!=0{t.Fatal("mismatched feedback mutated native memory")}
 if e=m.Feedback(a,0,a,0,false);e!=nil{t.Fatal(e)}
 if e=m.Feedback(a,0,a,0,false);e==nil{t.Fatal("duplicate native event accepted")}
 if e=m.Commit(a,1);e!=nil{t.Fatal(e)}
 if e=m.Feedback(a,1,b,0,false);e!=nil{t.Fatal(e)}
 if m.Samples!=2||m.Base.Native.Samples!=2{t.Fatal("original source selected counts differ")}
 proposed,status,e:=m.Choose(a)
 if e!=nil||proposed==nil||*proposed!=1||status!="SELECT_NATIVE_LEARNED_PUBLIC_TRANSITION_NOVELTY"{
  t.Fatalf("selected evidence did not distinguish observed state changes: %v %s %v",proposed,status,e)
 }
 raw,e:=m.Snapshot();if e!=nil{t.Fatal(e)}
 cold,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 raw2,e:=cold.Snapshot()
 if e!=nil||!bytes.Equal(raw,raw2){t.Fatal("native cold custody mismatch")}
 next,nextStatus,e:=cold.Choose(a)
 if e!=nil||next==nil||*next!=*proposed||nextStatus!=status{
  t.Fatalf("cold native choice drift: %v %s %v",next,nextStatus,e)
 }
 if _,e=Restore(append(append([]byte{},raw...),0));e==nil{t.Fatal("trailing checkpoint accepted")}
 if _,e=Restore(raw[:len(raw)-1]);e==nil{t.Fatal("truncated checkpoint accepted")}
}

func TestBadPublicAndPendingCannotMutateNativeMemory(t *testing.T){
 m:=New();a:=sourcePublic(0)
 bad:=sourcePublic(0);bad.Mission="hidden task coordinate"
 if proposal,_,e:=m.Choose(bad);e==nil||proposal!=nil{t.Fatal("hidden mission accepted")}
 if e:=m.ObserveFit(bad,1,a,0,false);e==nil{t.Fatal("bad FIT accepted")}
 if m.Samples!=0||m.Base.Native.Samples!=0{t.Fatal("bad input modified memory")}
 if e:=m.Commit(a,2);e!=nil{t.Fatal(e)}
 if _,_,e:=m.Choose(a);e==nil{t.Fatal("proposal with unreturned selected source")}
 if _,e:=m.Snapshot();e==nil{t.Fatal("pending source checkpoint accepted")}
 if e:=m.Feedback(a,2,a,0,false);e!=nil{t.Fatal(e)}
 if e:=m.BeginEpisode();e!=nil{t.Fatal(e)}
 if len(m.Visits)!=0||m.Samples!=1{t.Fatal("session reset erased learned source or leaked visits")}
}
