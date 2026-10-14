package unitary
import("encoding/json";"fmt";"math")
type nci06Fact struct{Edge nci02WEdge `json:"edge"`; Exposure,Inserted,Last int}
type nci06Trial struct{
 Seed int `json:"seed"`; Family int `json:"family"`
 LearnedTwo bool `json:"learned_two"`; LearnedFour bool `json:"learned_four"`
 LearnedRestored bool `json:"learned_restored"`
 FIFOTwo bool `json:"fifo_two"`; FIFOFour bool `json:"fifo_four"`
 LRUTwo bool `json:"lru_two"`; LRUFour bool `json:"lru_four"`
 OracleTwo bool `json:"oracle_two"`; OracleFour bool `json:"oracle_four"`
 LearnedEdges int `json:"learned_edges"`; FIFOEdges int `json:"fifo_edges"`
 LRUEdges int `json:"lru_edges"`
}
func nci06Fixture(seed,family int)([5]int,[]nci06Fact,[]nci02WEdge){
 positions:=[4][4]int{{0,1,2,3},{0,2,4,6},{1,3,5,7},{2,4,7,9}}
 var nodes [5]int
 for i:=0;i<5;i++{nodes[i]=64+(seed*7+family*3+i*[4]int{1,3,5,7}[family])%16}
 chain:=[]nci02WEdge{}
 for i:=0;i<4;i++{chain=append(chain,nci02WEdge{From:nodes[i],To:nodes[i+1]})}
 tcount,ncount:=0,0;stream:=[]nci06Fact{}
 for time:=0;time<12;time++{
  target:=false;for _,pos:=range positions[family]{if time==pos{target=true}}
  edge:=nci02WEdge{}
  if target{edge=chain[tcount];tcount++}else{
   edge=nci02WEdge{From:96+2*ncount+(seed+family)%2,To:120+2*ncount+(seed+2*family)%2};ncount++
  }
  n:=1
  switch family{
  case 0:if target{n=4}
  case 1:if target{n=3}else{n=1+(seed+time)%2}
  case 2:if !target{n=4}
  default:n=1+(seed*3+time*5)%4
  }
  stream=append(stream,nci06Fact{Edge:edge,Exposure:n,Inserted:time,Last:time})
 }
 if tcount!=4||ncount!=8{panic("bad fixture")}
 return nodes,stream,chain
}
func nci06Feature(f nci06Fact)[3]float64{
 return [3]float64{1,float64(f.Exposure)/4,float64(f.Inserted)/12}
}
func nci06Prob(w [3]float64,x [3]float64)float64{
 z:=0.;for i:=range x{z+=w[i]*x[i]}
 if z>35{z=35};if z< -35{z= -35}
 return 1/(1+math.Exp(-z))
}
func nci06Train()([3]float64,int,float64){
 type example struct{x [3]float64; y float64}
 rows:=[]example{}
 for seed:=25200;seed<25360;seed++{for family:=0;family<4;family++{
  _,stream,chain:=nci06Fixture(seed,family)
  for _,f:=range stream{
   label:=0.
   for _,e:=range chain{if e==f.Edge{label=1}}
   rows=append(rows,example{nci06Feature(f),label})
  }
 }}
 w:=[3]float64{}
 for epoch:=0;epoch<50;epoch++{
  for _,row:=range rows{
   g:=.035*(nci06Prob(w,row.x)-row.y)
   for k:=range w{w[k]-=g*row.x[k]}
  }
 }
 right,total:=0,0
 for seed:=25360;seed<25392;seed++{for family:=0;family<4;family++{
  _,stream,chain:=nci06Fixture(seed,family)
  for _,f:=range stream{
   truth:=false;for _,e:=range chain{if e==f.Edge{truth=true}}
   if (nci06Prob(w,nci06Feature(f))>=.5)==truth{right++}
   total++
  }
 }}
 return w,len(rows),float64(right)/float64(total)
}
func nci06Retain(stream []nci06Fact,mode string,w [3]float64)([]nci02WEdge,int,int){
 kept:=[]nci06Fact{};evict,mults,clock:=0,0,0
 for _,incoming:=range stream{
  clock++
  if len(kept)==7{
   at:=0;best:=math.Inf(1)
   for i,f:=range kept{
    score:=float64(f.Inserted)
    if mode=="lru"{score=float64(f.Last)}
    if mode=="learned"{score=nci06Prob(w,nci06Feature(f));mults+=3}
    if score<best{at=i;best=score}
   }
   kept=append(kept[:at],kept[at+1:]...);evict++
  }
  incoming.Last=clock;kept=append(kept,incoming)
  // One public rehearsal checkpoint, BEFORE capacity overflow, independent
  // of future questions. LRU sees only observed rehearsals, not labels.
  if incoming.Inserted==6{
   for i:=range kept{if kept[i].Exposure>=3{clock++;kept[i].Last=clock}}
  }
 }
 out:=[]nci02WEdge{};for _,f:=range kept{out=append(out,f.Edge)}
 return out,evict,mults
}
func nci06Reload(x []nci02WEdge)([]nci02WEdge,error){
 raw,e:=json.Marshal(x);if e!=nil{return nil,e}
 var y []nci02WEdge
 e=json.Unmarshal(raw,&y);return y,e
}
func RunNCI06Wingless(phase string)(map[string]any,error){
 if phase!="development"&&phase!="sealed"{return nil,fmt.Errorf("bad phase")}
 w,examples,validation:=nci06Train()
 result:=map[string]any{
  "schema":"wingless.nci06.learned-bounded-retention.v1","phase":phase,
  "parent":"a03f7efb53972ba4c811aa01000dacaf85e1cd3f",
  "retention_weights":w,"development_examples":examples,
  "development_validation_accuracy":validation,
  "sealed_seed_range":[]int{25632,25656},
  "memory_cap":7,"learnable_retention_parameters":3,"trained_join_parameters":6,
  "trained_admission_parameters":65,
  "learned_relational_storage":false,"learned_graph_search":false,
  "nci01_full_acceptance":false,"accepted_ref_mutation":false,"external_llm_calls":0,
  "verdict":"DEVELOPMENT_ONLY",
 }
 if phase=="development"{return result,nil}
 join,cut,_:=nci02WTrain()
 admission:=uplm0fTrainClassifier()
 admitted,missed:=0,0
 rows:=[]nci06Trial{}
 work:=map[string]int{"learned":0,"fifo":0,"lru":0}
 successes:=map[string]int{}
 evictions,scoreMult,checkpoints:=0,0,0
 for seed:=25632;seed<25656;seed++{for family:=0;family<4;family++{
  nodes,stream,chain:=nci06Fixture(seed,family)
  observed:=[]nci06Fact{}
  for _,fact:=range stream{
   clause:=nci01WClause(fact.Edge.From,fact.Edge.To,true)
   if admission.prob(uplm0fEncode(clause))>=.5{
    admitted++;observed=append(observed,fact)
   }else{missed++}
  }
  // All candidate policies see the same native-admitted stream.
  stream=observed
  L,eL,cL:=nci06Retain(stream,"learned",w)
  F,eF,_:=nci06Retain(stream,"fifo",w)
  R,eR,_:=nci06Retain(stream,"lru",w)
  evictions+=eL+eF+eR;scoreMult+=cL
  wl,wf,wr:=0,0,0
  r:=nci06Trial{Seed:seed,Family:family}
  r.LearnedTwo=nci02WReach(join,cut,L,nodes[0],nodes[2],2,&wl)
  r.LearnedFour=nci02WReach(join,cut,L,nodes[0],nodes[4],4,&wl)
  r.FIFOTwo=nci02WReach(join,cut,F,nodes[0],nodes[2],2,&wf)
  r.FIFOFour=nci02WReach(join,cut,F,nodes[0],nodes[4],4,&wf)
  r.LRUTwo=nci02WReach(join,cut,R,nodes[0],nodes[2],2,&wr)
  r.LRUFour=nci02WReach(join,cut,R,nodes[0],nodes[4],4,&wr)
  work["learned"]+=wl;work["fifo"]+=wf;work["lru"]+=wr
  restored,err:=nci06Reload(L);if err!=nil{return nil,err};checkpoints++
  recOps:=0;r.LearnedRestored=nci02WReach(join,cut,restored,nodes[0],nodes[4],4,&recOps)
  work["learned"]+=recOps
  r.OracleTwo=nci04WOracle(chain,nodes[0],nodes[2],2)
  r.OracleFour=nci04WOracle(chain,nodes[0],nodes[4],4)
  r.LearnedEdges=len(L);r.FIFOEdges=len(F);r.LRUEdges=len(R)
  if len(L)!=7||len(F)!=7||len(R)!=7{return nil,fmt.Errorf("memory budget failure")}
  for k,b:=range map[string]bool{
   "learned_two":r.LearnedTwo,"learned_four":r.LearnedFour,
   "learned_restored":r.LearnedRestored,"fifo_two":r.FIFOTwo,
   "fifo_four":r.FIFOFour,"lru_two":r.LRUTwo,"lru_four":r.LRUFour,
   "oracle_two":r.OracleTwo,"oracle_four":r.OracleFour}{
    if b{successes[k]++}
  }
  rows=append(rows,r)
 }}
 if len(rows)!=96{return nil,fmt.Errorf("unexpected cohort")}
 accuracy:=map[string]float64{}
 for k,v:=range successes{accuracy[k]=float64(v)/96}
 byFamily:=map[string]map[string]float64{}
 for f:=0;f<4;f++{
  scores:=map[string]float64{};L,F,R:=0,0,0
  for _,r:=range rows{if r.Family==f{
   if r.LearnedFour{L++};if r.FIFOFour{F++};if r.LRUFour{R++}
  }}
  scores["learned_four"]=float64(L)/24
  scores["fifo_four"]=float64(F)/24
  scores["lru_four"]=float64(R)/24
  byFamily[fmt.Sprintf("%d",f)]=scores
 }
 gainFIFO:=accuracy["learned_four"]-accuracy["fifo_four"]
 gainLRU:=accuracy["learned_four"]-accuracy["lru_four"]
 noBigRegression:=true
 for _,g:=range byFamily{if g["learned_four"]<g["lru_four"]-.20{noBigRegression=false}}
 verdict:="NO_TRANSFERABLE_LEARNED_RETENTION_GAIN"
 if gainFIFO>=.15&&gainLRU>=.10&&accuracy["learned_four"]>=.60&&noBigRegression{
  verdict="LEARNED_BOUNDED_RETENTION_TRANSFER"
 }else if gainFIFO>=.10{verdict="PARTIAL_RETENTION_GAIN_WITH_GENERALIZATION_COST"}
 result["native_admission_true_positive"]=admitted
 result["native_admission_false_negative"]=missed
 result["records"]=rows;result["accuracy"]=accuracy;result["by_family"]=byFamily
 result["gain_vs_fifo_four"]=gainFIFO;result["gain_vs_lru_four"]=gainLRU
 result["learned_join_candidate_evaluations"]=work
 result["retention_score_multiplies"]=scoreMult
 result["evictions_all_arms"]=evictions;result["checkpoint_restores"]=checkpoints
 result["verdict"]=verdict
 return result,nil
}
