package unitary

import "fmt"
const UPLM6BSchema="wingless.up-lm6b-prepressure-extension-map.v1"
type UPLM6BRow struct{Profile string `json:"profile"`;Rotation int `json:"rotation"`;Permutation string `json:"permutation"`;Arm int `json:"arm"`;Target int `json:"target"`;PreLength int `json:"pre_length"`;PreCountdown int `json:"pre_countdown"`;Writes int `json:"writes"`;PostLength int `json:"post_length"`;PostCountdown int `json:"post_countdown"`}
type UPLM6BResult struct{Schema string `json:"schema"`;Experiment string `json:"experiment"`;SourceUPLM6ASeal string `json:"source_up_lm6a_seal"`;ParentClassification string `json:"parent_classification"`;Rows []UPLM6BRow `json:"rows"`;MinWrites int `json:"min_writes"`;MaxWrites int `json:"max_writes"`;DistinctWriteCounts int `json:"distinct_write_counts"`;MinPostLength int `json:"min_post_length"`;MaxPostLength int `json:"max_post_length"`;LiveActivation bool `json:"live_activation"`;PolicyChanged bool `json:"policy_changed"`;Classification string `json:"classification"`}
func RunUPLM6B()(UPLM6BResult,error){
 p,err:=RunUPLM6A();if err!=nil{return UPLM6BResult{},err}
 profiles:=[]string{"deferred_only","layout_only","hybrid_min"};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
 r:=UPLM6BResult{Schema:UPLM6BSchema,Experiment:"UP-LM6B-prepressure-extension-map",SourceUPLM6ASeal:"8c70e27c629159262bbe7c04cfde269757713a77",ParentClassification:p.Classification,MinWrites:1<<30,MinPostLength:1<<30}
 counts:=map[int]bool{}
 for _,profile:=range profiles{for _,rot:=range rots{for _,perm:=range perms{
  arms:=uplm2yPermute(uplm2xArms(rot),perm)
  for i:=range arms{
   target:=uplm3tTarget(arms[i],profile);preLen:=len(arms[i].r.order);preCD:=0;if _,cd,ok:=uplm2xFirstPending(&arms[i]);ok{preCD=cd}
   writes:=0
   for n:=0;n<32;n++{_,cd,ok:=uplm2xFirstPending(&arms[i]);if !ok||cd<=target{break};arms[i].r.write(fmt.Sprintf("3t-pre-%s-%d-%d-%d",profile,rot,i,n),"x");writes++}
   postLen:=len(arms[i].r.order);postCD:=0;if _,cd,ok:=uplm2xFirstPending(&arms[i]);ok{postCD=cd}
   r.Rows=append(r.Rows,UPLM6BRow{Profile:profile,Rotation:rot,Permutation:perm,Arm:i,Target:target,PreLength:preLen,PreCountdown:preCD,Writes:writes,PostLength:postLen,PostCountdown:postCD})
   counts[writes]=true;if writes<r.MinWrites{r.MinWrites=writes};if writes>r.MaxWrites{r.MaxWrites=writes};if postLen<r.MinPostLength{r.MinPostLength=postLen};if postLen>r.MaxPostLength{r.MaxPostLength=postLen}
  }
 }}}
 r.DistinctWriteCounts=len(counts)
 if r.ParentClassification!="PREPRESSURE_CREATES_LENGTH16_FLOOR"{r.Classification="ANCHOR_NOT_REPRODUCED"}else if r.MinWrites>0&&r.MinWrites==r.MaxWrites{r.Classification="UNIFORM_EXTENSION"}else if r.MaxWrites>0&&r.DistinctWriteCounts>1{r.Classification="TARGET_DEPENDENT_EXTENSION"}else if r.MaxWrites==0{r.Classification="NO_PREPRESSURE_EXTENSION"}else{r.Classification="OTHER_VALID_PATTERN"}
 return r,nil
}
