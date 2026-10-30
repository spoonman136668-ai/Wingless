package mg37
import ("bytes";"testing")
func TestMG50ContextRewardColdSnapshot(t *testing.T){
 a:=NewSourceAgent(true)
 red:=mg38Public()
 blue:=mg38Public();blue.Mission="pick up the blue ball"
 kr,err:=publicContext(red);if err!=nil{t.Fatal(err)}
 kb,err:=publicContext(blue);if err!=nil{t.Fatal(err)}
 var row [7]Evidence
 row[0]=Evidence{Count:6,RewardSum:0}
 row[1]=Evidence{Count:6,RewardSum:6}
 a.Rewards[kr]=row
 wire,err:=a.Snapshot();if err!=nil{t.Fatal(err)}
 restored,err:=RestoreSourceAgent(wire);if err!=nil{t.Fatal(err)}
 again,err:=restored.Snapshot();if err!=nil||!bytes.Equal(wire,again){t.Fatal("cold checkpoint changed",err)}
 d,err:=Choose(0,restored.Rewards[kr]);if err!=nil||!d.Override||d.Action!=1{t.Fatal("reward evidence lost",d,err)}
 d,err=Choose(0,restored.Rewards[kb]);if err!=nil||d.Override{t.Fatal("cross-context leakage",d,err)}
 if _,err=RestoreSourceAgent(append(append([]byte{},wire...),byte('x')));err==nil{t.Fatal("accepted tamper")}
 a.HasPending=true
 if _,err=a.Snapshot();err==nil{t.Fatal("accepted in-flight snapshot")}
}
