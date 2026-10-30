package mg37
import (
 "testing"
 original "github.com/spoonman136668-ai/Wingless/research/native/itd02/mg28"
)
func mg38Public() original.Public {
 p:=original.Public{Image:make([][][]int,7),Direction:0,Mission:"pick up the red ball"}
 for i:=0;i<7;i++ {p.Image[i]=make([][]int,7);for j:=0;j<7;j++ {p.Image[i][j]=[]int{1,0,0}}}
 return p
}
func TestMG38SelectedFeedbackAndNoReplay(t *testing.T){
 a:=NewSourceAgent(true)
 before:=mg38Public()
 after:=mg38Public()
 after.Direction=1
 if err:=a.BeginEpisode();err!=nil {t.Fatal(err)}
 if err:=a.Feedback(before,1,after,1,false);err==nil {t.Fatal("unselected feedback accepted")}
 proposed,_,err:=a.Choose(before);if err!=nil||proposed==nil {t.Fatal(err)}
 if err:=a.Commit(before,*proposed);err!=nil {t.Fatal(err)}
 if err:=a.Feedback(before,*proposed,after,1,false);err!=nil {t.Fatal(err)}
 if a.Rewards[*proposed].Count!=1 || a.Rewards[*proposed].RewardSum!=1 {t.Fatal("selected reward missing")}
 if err:=a.Feedback(before,*proposed,after,1,false);err==nil {t.Fatal("duplicate feedback accepted")}
}
