package mg37
import "testing"
func TestObservedRewardBeatsChurn(t *testing.T){
 var e [7]Evidence
 e[0]=Evidence{6,0}; e[1]=Evidence{6,6}
 d,err:=Choose(0,e); if err!=nil||!d.Override||d.Action!=1 {t.Fatalf("rewardful alternative not selected: %+v %v",d,err)}
 d,err=Choose(1,e);if err!=nil||d.Override||d.Action!=1 {t.Fatalf("churn replaced productive action: %+v %v",d,err)}
}
func TestUnknownAndUnproductiveAbstain(t *testing.T){
 var e [7]Evidence
 e[0]=Evidence{6,0};e[1]=Evidence{6,0}
 d,err:=Choose(0,e);if err!=nil||d.Override {t.Fatal(d,err)}
 e[1]=Evidence{2,2};d,err=Choose(0,e);if err!=nil||d.Override {t.Fatal(d,err)}
}
func TestMalformedEvidenceReject(t *testing.T){
 var e [7]Evidence;e[1]=Evidence{3,4}
 if _,err:=Choose(0,e);err==nil{t.Fatal("accepted impossible reward sum")}
 if _,err:=Choose(7,[7]Evidence{});err==nil{t.Fatal("accepted invalid action")}
}
