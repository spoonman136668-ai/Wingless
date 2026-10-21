package cognition

// NCI-04: direct native LatentState revision test, NOT an independent
// unfamiliar-task evaluation. The only training evidence is a sensor pair,
// executed action and newly observed consequence. No task/world/entity ID or
// correct answer is supplied to LatentState.
import (
 "bytes"
 "math"
 "testing"
 "unsafe"
)

func TestNCI04ContradictionDrivenStateRevision(t *testing.T){
 var s LatentState
 a:=[FeatureWidth]float32{.15,0,0,0}
 b:=[FeatureWidth]float32{.32,0,0,0} // initially within the original merge radius
 dst:=[FeatureWidth]float32{0,1,0,0}
 for i:=0;i<100;i++{
  if err:=s.Observe(a,dst,0,1);err!=nil{t.Fatal(err)}
 }
 count:=s.Count
 first,_:=s.Predict(a,dst,0)
 if first<.90||count!=2{t.Fatalf("trained baseline not stable: %f count=%d",first,count)}
 // Same action/destination, different sensory state. The consequence
 // contradicts established causal relation strongly enough to split.
 if err:=s.Observe(b,dst,0,0);err!=nil{t.Fatal(err)}
 if s.Count!=count+1{t.Fatalf("observed contradiction did not revise identity: %d",s.Count)}
 initially,_:=s.Predict(b,dst,0)
 if !(initially>0.48&&initially<0.5){t.Fatalf("new hypothesis was already confident: %f",initially)}
 for i:=0;i<100;i++{
  if err:=s.Observe(b,dst,0,0);err!=nil{t.Fatal(err)}
 }
 positive,_:=s.Predict(a,dst,0)
 negative,_:=s.Predict(b,dst,0)
 if positive<.90||negative>.10||s.Count!=3{
  t.Fatalf("native learned state not functionally differentiated: A=%f B=%f entities=%d",positive,negative,s.Count)
 }
 // An input equidistant from contradictory prototypes has no identified
 // causal state. The runtime must abstain rather than make a hard prediction.
 ambiguous:=[FeatureWidth]float32{.235,0,0,0}
 undecided,_:=s.Predict(ambiguous,dst,0)
 if undecided!=.5{t.Fatalf("boundary had unjustified certainty: %f",undecided)}
 prior,_:=s.Snapshot()
 var restored LatentState
 if e:=restored.Restore(prior);e!=nil{t.Fatal(e)}
 after,_:=restored.Snapshot()
 if !bytes.Equal(prior,after){t.Fatal("learned partition lost over cold-state restore")}
 a2,_:=restored.Predict(a,dst,0)
 b2,_:=restored.Predict(b,dst,0)
 if math.Abs(a2-positive)>1e-8||math.Abs(b2-negative)>1e-8{t.Fatal("restart changed native causal response")}
 if unsafe.Sizeof(s)!=unsafe.Sizeof(LatentState{}){t.Fatal("native learned state grew")}
 t.Logf("NCI04_CAUSAL_REVISION_DEVELOPMENT_PASS: native entity slots=%d bytes=%d; controls, real native consequences, full resource parity and scientific transfer NOT QUALIFIED",s.Count,unsafe.Sizeof(s))
}

func TestNCI04RejectsUndistinguishableOrUntrainedCausalSplit(t *testing.T){
 var s LatentState
 a:=[FeatureWidth]float32{.15,0,0,0}
 b:=[FeatureWidth]float32{.32,0,0,0}
 dst:=[FeatureWidth]float32{0,1,0,0}
 // A few observations cannot justify a new causal state.
 for i:=0;i<2;i++{if e:=s.Observe(a,dst,0,1);e!=nil{t.Fatal(e)}}
 if e:=s.Observe(b,dst,0,0);e!=nil{t.Fatal(e)}
 if s.Count!=2{t.Fatal("weak evidence created unsupported causal split")}
 // Identical observations cannot be distinguished even after strong
 // training. Do not fabricate an entity identity from consequence alone.
 var trained LatentState
 for i:=0;i<100;i++{if e:=trained.Observe(a,dst,0,1);e!=nil{t.Fatal(e)}}
 existing:=trained.Count
 if e:=trained.Observe(a,dst,0,0);e!=nil{t.Fatal(e)}
 if trained.Count!=existing{t.Fatal("indistinguishable sensor falsely split")}
}

func TestNCI04RevisionFailClosedWhenFixedSlotsExhausted(t *testing.T){
 var s LatentState
 src:=[FeatureWidth]float32{0,0,0,0}
 dst:=[FeatureWidth]float32{0,1,0,0}
 for i:=0;i<100;i++{if e:=s.Observe(src,dst,0,1);e!=nil{t.Fatal(e)}}
 // Only an already-valid, fully occupied native memory should refuse
 // further causal splits rather than overwrite learned evidence.
 for s.Count<MaxEntities{
  var newv [FeatureWidth]float32
  v:=float32(s.Count)
  newv=[FeatureWidth]float32{-.8+v*.1,-.8,.8,0}
  if e:=s.Observe(newv,newv,1,1);e!=nil{t.Fatal(e)}
 }
 original,_:=s.Snapshot()
 contrast:=[FeatureWidth]float32{.2,0,0,0}
 if e:=s.Observe(contrast,dst,0,0);e==nil{t.Fatal("capacity bypassed")}
 after,_:=s.Snapshot()
 if !bytes.Equal(original,after){t.Fatal("capacity rejection modified learned evidence")}
}
