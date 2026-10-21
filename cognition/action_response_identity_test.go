package cognition

import ("fmt";"testing";"unsafe")
func runNCI04ESourceTrial(t *testing.T,m *ActionResponseIdentity,latent int,rounds int){
 t.Helper()
 if err:=m.BeginTrial();err!=nil{t.Fatal(err)}
 for i:=0;i<rounds;i++{
  a:=i%2; outcome:=a^latent
  if err:=m.Observe(a,outcome);err!=nil{t.Fatal(err)}
 }
}
func TestNCI04EActionConditionalSeparation(t *testing.T){
 m:=NewActionResponseIdentity()
 for session:=0;session<8;session++{
  // Hidden context alternates, but learner receives no identity argument.
  runNCI04ESourceTrial(t,m,session%2,8)
 }
 if m.ModeCount()!=2{t.Fatalf("learned mode count %d, require two",m.ModeCount())}
 delta,ok:=m.Separation(0)
 if !ok||delta<.6{t.Fatalf("no action-conditioned separation %.4f %v",delta,ok)}
 correct:=0;pooledOnes,pooledZeros:=0,0
 for session:=0;session<8;session++{
  latent:=session%2
  if e:=m.BeginTrial();e!=nil{t.Fatal(e)}
  if _,known,e:=m.Predict(1);e!=nil||known{t.Fatal("identity guessed before observation")}
  // One controlled action-0 outcome reveals a response, NOT an identity tag.
  if e:=m.Observe(0,latent);e!=nil{t.Fatal(e)}
  p,known,e:=m.Predict(1);if e!=nil||!known{t.Fatalf("no query after evidence: %v %v",known,e)}
  truth:=1-latent
  if (p>=.5)==(truth==1){correct++}
  if truth==1{pooledOnes++}else{pooledZeros++}
 }
 if correct!=8{t.Fatalf("did not predict held-out action responses %d/8",correct)}
 pooled:=float64(pooledOnes+1)/float64(pooledOnes+pooledZeros+2)
 fmt.Printf("NCI04E_W_DISTINGUISHABLE=8/8 POOLED_ACTION1_BASELINE=%.4f MODES=%d BYTES=%d\n",pooled,m.ModeCount(),unsafe.Sizeof(*m))
 // One pooled trained model is uninformative because hidden contexts conflict.
 if pooled!=.5{t.Fatal("positive fixture not balanced")}
}
func TestNCI04EIndistinguishableMustNotClaimIdentity(t *testing.T){
 m:=NewActionResponseIdentity()
 for session:=0;session<8;session++{
  // Two nominal sources; no distinguishing action-history evidence.
  runNCI04ESourceTrial(t,m,0,8)
 }
 if m.ModeCount()!=1{t.Fatalf("HARD_FALSE_POSITIVE: split indistinguishable sources into %d",m.ModeCount())}
 if err:=m.BeginTrial();err!=nil{t.Fatal(err)}
 if _,known,e:=m.Predict(1);e!=nil||known{t.Fatal("HARD_FALSE_POSITIVE: predicted identity without evidence")}
 fmt.Printf("NCI04E_W_INDISTINGUISHABLE_SPLITS=0\n")
}
func TestNCI04EBoundedSourceAndFailClosed(t *testing.T){
 m:=NewActionResponseIdentity()
 if e:=m.Observe(9,1);e==nil{t.Fatal("invalid action")}
 if e:=m.Observe(1,-1);e==nil{t.Fatal("invalid outcome")}
 if m.ModeCount()!=0{t.Fatal("invalid observations mutated model")}
 if _,_,e:=m.Predict(2);e==nil{t.Fatal("invalid query accepted")}
 if err:=m.BeginTrial();err!=nil{t.Fatal(err)}
 if len(m.Modes)!=NCI04EMaxModes{t.Fatal("capacity mismatch")}
}
