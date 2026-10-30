package mg37
import "testing"
func TestMG43ContextEvidenceDoesNotTransfer(t *testing.T) {
 a:=NewSourceAgent(true)
 red:=mg38Public()
 blue:=mg38Public();blue.Mission="pick up the blue ball"
 cr,err:=publicContext(red);if err!=nil{t.Fatal(err)}
 cb,err:=publicContext(blue);if err!=nil{t.Fatal(err)}
 if cr==cb{t.Fatal("mission contexts collided")}
 var row [7]Evidence
 row[1]=Evidence{Count:5,RewardSum:5}
 row[0]=Evidence{Count:5,RewardSum:0}
 a.Rewards[cr]=row
 if len(a.Rewards[cb])!=0{t.Fatal("unexpected foreign evidence")}
 d,err:=Choose(0,a.Rewards[cr]);if err!=nil||!d.Override||d.Action!=1{t.Fatal("rewardful original context ignored",d,err)}
 d,err=Choose(0,a.Rewards[cb]);if err!=nil||d.Override{t.Fatal("cross-context reward leakage",d,err)}
}
