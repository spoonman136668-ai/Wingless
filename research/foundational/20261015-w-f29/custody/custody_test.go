package custody
import("math";"reflect";"testing";"unsafe")
func specimen(truth bool) PrivateCase {return PrivateCase{Truth:truth,Input:PublicInput{Evidence:[]Evidence{
 {EncounterID:11,Role:0,Left:[3]float64{.1,.2,.3},Right:[3]float64{.11,.21,.6}},
 {EncounterID:12,Role:2,Left:[3]float64{.3,.4,.5},Right:[3]float64{.31,.41,.9}},
 {EncounterID:13,Role:1,Left:[3]float64{.2,.7,.8},Right:[3]float64{.22,.73,.2}},
}}}}
func TestSplitIsolation(t *testing.T){
 train:=Manifest{[]uint64{1,2},[]uint64{11,12},[]uint64{112,113}}
 test:=Manifest{[]uint64{7},[]uint64{20},[]uint64{116}}
 if err:=ValidateSplit(train,test);err!=nil{t.Fatal(err)}
 test.EntityIDs=[]uint64{12};if ValidateSplit(train,test)==nil{t.Fatal("identity leak")}
 test.EntityIDs=[]uint64{20};test.GraphFamilies=[]uint64{112};if ValidateSplit(train,test)==nil{t.Fatal("family reuse")}
 test.GraphFamilies=[]uint64{116};test.Seeds=[]uint64{2};if ValidateSplit(train,test)==nil{t.Fatal("seed reuse")}
}
func TestPrivateTruthUnreachableAndCostMetered(t *testing.T){
 if _,ok:=reflect.TypeOf(PublicInput{}).FieldByName("Truth");ok {t.Fatal("truth leaked")}
 a,ma,e:=Prepare(specimen(true),576);if e!=nil{t.Fatal(e)}
 b,mb,e:=Prepare(specimen(false),576);if e!=nil{t.Fatal(e)}
 if !reflect.DeepEqual(a,b)||!reflect.DeepEqual(ma,mb){t.Fatal("public model input depends on truth")}
 ok,m,e:=ScoreCase(specimen(true),576,Consensus);if e!=nil||!ok{t.Fatalf("consensus %v %v",ok,e)}
 if m.WitnessComparisons!=9||m.HistoryBytes!=3*int(unsafe.Sizeof(Evidence{}))||m.NativeBytes!=576{t.Fatalf("cost drift %+v",m)}
}
func TestFailClosedOnBudgetDataAndEncounter(t *testing.T){
 if _,_,err:=Prepare(specimen(true),577);err==nil{t.Fatal("resident overflow accepted")}
 c:=specimen(true);c.Input.Evidence[1].EncounterID=11
 if _,_,err:=Prepare(c,128);err==nil{t.Fatal("duplicate encounter accepted")}
 c=specimen(true);c.Input.Evidence[0].Left[0]=math.NaN()
 if _,_,err:=Prepare(c,128);err==nil{t.Fatal("NaN accepted")}
}
