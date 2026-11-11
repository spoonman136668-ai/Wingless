package mg40
import ("bytes";"testing")

func frame(mission string,heading int,front [3]int)Public {
 p:=Public{Image:make([][][]int,7),Direction:heading,Mission:mission}
 for i:=0;i<7;i++ {p.Image[i]=make([][]int,7);for j:=0;j<7;j++{p.Image[i][j]=[]int{1,0,0}}}
 p.Image[3][5]=[]int{front[0],front[1],front[2]}
 return p
}
func TestNativeSelectedTaskEffectRejectsImageChurn(t *testing.T){
 mission:="pick up the yellow ball"
 base:=frame(mission,0,[3]int{4,1,1})
 opened:=frame(mission,0,[3]int{4,1,0})
 closed:=frame(mission,0,[3]int{4,1,1})
 good,e:=TaskEffect(base,5,opened)
 if e!=nil||!good {t.Fatal("missed source-public actual door opening",e)}
 no,e:=TaskEffect(opened,5,closed)
 if e!=nil||no {t.Fatal("door closing given positive task credit",e)}
 turn:=frame(mission,1,[3]int{4,1,0})
 no,e=TaskEffect(base,0,turn)
 if e!=nil||no {t.Fatal("turn changed view but was credited",e)}
 pickup:=frame(mission,0,[3]int{6,4,0})
 after:=frame(mission,0,[3]int{1,0,0})
 good,e=TaskEffect(pickup,3,after)
 if e!=nil||!good{t.Fatal("exact original mission target pickup missing",e)}
 falseTarget:=frame(mission,0,[3]int{6,0,0})
 no,e=TaskEffect(falseTarget,3,after)
 if e!=nil||no{t.Fatal("nonmission ball credited",e)}
}
func TestTaskPolicyUsesOnlyPriorSelectedPositiveEffect(t *testing.T){
 var row Row
 row[0]=Stat{Count:3,Effect:0}
 row[5]=Stat{Count:1,Effect:1}
 a,o,e:=Propose(0,row)
 if e!=nil||o||a!=0 {t.Fatal("invented prior support",a,o,e)}
 row[5]=Stat{Count:3,Effect:2}
 a,o,e=Propose(0,row)
 if e!=nil||!o||a!=5 {t.Fatal("real selected door evidence not considered",a,o,e)}
 row[0]=Stat{Count:1,Effect:0}
 a,o,e=Propose(0,row)
 if e!=nil||o||a!=0 {t.Fatal("unequal parent evidence leaked",a,o,e)}
 if _,_,e=Propose(7,row);e==nil{t.Fatal("bad original action accepted")}
}
func TestNativeSourceLedgerAndColdTaskCheckpoint(t *testing.T){
 a:=New(true)
 before:=frame("pick up the red key",0,[3]int{4,2,2})
 after:=frame(before.Mission,0,[3]int{4,2,0})
 for i:=0;i<2;i++ {
  if e:=a.BeginEpisode();e!=nil{t.Fatal(e)}
  if e:=a.Feedback(before,5,after,0,false);e==nil{t.Fatal("uncommitted reward admitted")}
  if e:=a.Commit(before,5);e!=nil{t.Fatal(e)}
  if _,e:=a.Snapshot();e==nil{t.Fatal("pending source event checkpointed")}
  if e:=a.Feedback(before,4,after,0,false);e==nil{t.Fatal("wrong selected action admitted")}
  if e:=a.Feedback(before,5,after,0,false);e!=nil{t.Fatal(e)}
  if e:=a.Feedback(before,5,after,0,false);e==nil{t.Fatal("duplicate selected event counted")}
 }
 k,e:=publicKey(before);if e!=nil{t.Fatal(e)}
 if a.Task[k][5]!=(Stat{Count:2,Effect:2}){t.Fatal("original source effect not retained",a.Task[k][5])}
 raw,e:=a.Snapshot();if e!=nil{t.Fatal(e)}
 copy,e:=Restore(raw);if e!=nil{t.Fatal(e)}
 same,e:=copy.Snapshot()
 if e!=nil||!bytes.Equal(raw,same){t.Fatal("cold task checkpoint drift",e)}
 if _,e=Restore(append(append([]byte(nil),raw...),'x'));e==nil {t.Fatal("trailing forged source accepted")}
}
