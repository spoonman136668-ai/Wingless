package unitary
import("fmt";"math";"time")

// F23: infer which of THREE shuffled, competing identity claims should
// create an internal recurrent relationship. Claims, evidence, and candidate
// endpoints are written into the 576-byte state, not separately supplied.
// The candidate alias proposals are externally given; graph endpoint pair
// matching remains CODED. Frozen W21 4-weight testimony model, F14/F22 core.
type wf23Row struct{
 Seed,Family,Depth,Shift,Variant int
 Truth bool
 Selected,FirstOnly,AllLinks,NoLinks,CodedUpper wf15Choice
 SelectedLinks int
 ProposalPerm [3]int
}
func wf23Infer(state wf1State,q [2]int,query [8]int,k int,
 role wf13Roles,enc wf7Encoder,core wf5Composer,p wf21Weights,mode string)(wf15Choice,int){
 var src,dst,tags [12]int
 for i:=0;i<12;i++{
  bits:=math.Float32bits(state[i])
  src[i],dst[i]=int(bits>>16),int(bits&65535)
  tags[i]=int(state[12+i])
 }
 var old,alias [3]int
 var trust [3]bool
 chosen:=0
 for c:=0;c<3;c++{
  b:=math.Float32bits(state[24+4*c])
  alias[c],old[c]=int(b>>16),int(b&65535)
  x:=[3]float64{float64(state[25+c*4]),float64(state[26+c*4]),float64(state[27+c*4])}
  confidence:=wf21Sig(p.V[0]+p.V[1]*x[0]+p.V[2]*x[1]+p.V[3]*x[2])
  switch mode{
  case "learned":trust[c]=confidence>=.5
  case "first":trust[c]=c==0
  case "all":trust[c]=true
  case "none":trust[c]=false
  case "coded":trust[c]=wf21EvidenceTruth(x)
  }
  if trust[c]{chosen++}
 }
 var active [12]float64
 for h:=0;h<k;h++{
  var next [12]float64
  for j:=0;j<12;j++{
   activation:=0.
   if h==0{activation,_,_=wf7Pair(enc,q[0],src[j])}else{
    msg:=0.
    for i:=0;i<12;i++{
     if i==j{continue}
     sim,_,_:=wf7Pair(enc,dst[i],src[j])
     for c:=0;c<3;c++{
      if trust[c]&&old[c]!=alias[c]&&dst[i]==old[c]&&src[j]==alias[c]{
       sim=math.Max(sim,1.)
      }
     }
     msg+=active[i]*sim
    }
    activation=wf1Prob(core.W[0]*msg+core.W[1]*active[j]+core.W[2])
   }
   next[j]=activation*wf13RoleScore(role,tags[j],query[h])
  }
  active=next
 }
 score:=0.
 for i:=0;i<12;i++{
  z,_,_:=wf7Pair(enc,dst[i],q[1]);score+=active[i]*z
 }
 prob:=wf1Prob(core.W[3]*score+core.W[4])
 return wf15Choice{Answer:prob>=.5,Steps:k,EarlyHalt:false},chosen
}
func RunWFoundationF23()(map[string]any,error){
 start:=time.Now()
 policy:=wf21Train()
 old:=[4]float64{1.1898503414192076e-16,2.8522704808696697,2.85227048086967,-2.85227048086967}
 for i,w:=range policy.V{if math.Abs(w-old[i])>1e-8{return nil,fmt.Errorf("F21 FROZEN POLICY DRIFT")}}
 entity,err:=wf8LoadFrozen();if err!=nil{return nil,err}
 role:=wf14FrozenRoles()
 core:=wf5Composer{W:[5]float64{6.470679800286544,-1.3847365385924115,-3.0464843010034035,7.51083733991036,-3.4005406394991717}}
 rows:=[]wf23Row{}
 for seed:=88400;seed<88406;seed++{for fam:=108;fam<=111;fam++{
  for k:=7;k<=8;k++{for shift:=0;shift<2;shift++{
   w,main,alt,noise,err:=wf19Graph(seed,k,fam,shift)
   if err!=nil{return nil,err}
   fork:=1+fam-108
   original:=wf1Rename(k+35,seed,2)
   aliasA:=wf1Rename(k+90,seed,2)
   aliasB:=wf1Rename(k+91,seed,2)
   unrelated:=wf1Rename(k+92,seed,2)
   second:=-1
   for i,e:=range w.Facts.Edges{
    if e.From==original&&e.To==wf1Rename(fork+2,seed,2){second=i}
   }
   if second<0{return nil,fmt.Errorf("original role-bearing bridge missing")}
   q:=[2]int{w.Src,w.Dst}
   for variant:=0;variant<12;variant++{
    graph:=w
    evidence:=wf21Patterns[variant%8]
    trueCert:=false
    if variant<8{
     graph.Roles[main]=(graph.Roles[main]+1)%4
     graph.Facts.Edges[second].From=aliasA
     trueCert=wf21EvidenceTruth(evidence)
    }else{
     evidence=[3]float64{1,1,-1}
     switch variant{
     case 8:
     case 9:
      graph.Roles[main]=(graph.Roles[main]+1)%4
      graph.Facts.Edges[second].From=aliasB
     case 10:
      graph.Roles[main]=(graph.Roles[main]+1)%4
      graph.Roles[alt]=(graph.Roles[alt]+1)%4
     case 11:graph.Roles[noise]=(graph.Roles[noise]+1)%4
     }
    }
    // Ordered proposals: A correct potentially, B unverified negative,
    // unrelated decoy positive; order permuted deterministically.
    claims:=[3]struct{alias,old int;e [3]float64}{
     {aliasA,original,evidence},
     {aliasB,original,[3]float64{-1,-1,1}},
     {unrelated,original,[3]float64{1,1,-1}},
    }
    if variant==9{claims[1].e=[3]float64{-1,-1,1}}
    if variant==8{claims[1].e=[3]float64{-1,-1,1}}
    expected:=true
    if variant<8{expected=trueCert}
    if variant==9||variant==10{expected=false}
    reference:=graph
    if variant<8&&trueCert{reference.Facts.Edges[second].From=original}
    if got:=wf13Exact(reference,w.Required,k,nil);got!=expected{
     return nil,fmt.Errorf("INVALID_W23_SOURCE %d/%d/%d",seed,fam,variant)
    }
    state,err:=wf14Serialize(graph);if err!=nil{return nil,err}
    var permutation [3]int
    shiftPerm:=(seed+fam+variant+shift)%3
    for j:=0;j<3;j++{
     selected:=(j+shiftPerm)%3
     permutation[j]=selected
     c:=claims[selected]
     state[24+4*j]=math.Float32frombits(uint32(c.alias)<<16|uint32(c.old))
     for z,v:=range c.e{state[25+4*j+z]=float32(v)}
    }
    selected,n:=wf23Infer(state,q,w.Required,k,role,entity,core,policy,"learned")
    first,_:=wf23Infer(state,q,w.Required,k,role,entity,core,policy,"first")
    all,_:=wf23Infer(state,q,w.Required,k,role,entity,core,policy,"all")
    none,_:=wf23Infer(state,q,w.Required,k,role,entity,core,policy,"none")
    coded,_:=wf23Infer(state,q,w.Required,k,role,entity,core,policy,"coded")
    rows=append(rows,wf23Row{Seed:seed,Family:fam,Depth:k,Shift:shift,
      Variant:variant,Truth:expected,Selected:selected,FirstOnly:first,
      AllLinks:all,NoLinks:none,CodedUpper:coded,
      SelectedLinks:n,ProposalPerm:permutation})
   }
  }}
 }}
 if len(rows)!=1152{return nil,fmt.Errorf("W23 count drift")}
 scores:=map[string]map[string]float64{}
 for _,name:=range []string{"learned","first","all","none","coded"}{
  corr,fp,fn:=0,0,0
  worst:=0.
  for fam:=108;fam<=111;fam++{
   n,bad:=0,0
   for _,r:=range rows{if r.Family!=fam{continue}
    var p wf15Choice
    switch name{
    case "learned":p=r.Selected
    case "first":p=r.FirstOnly
    case "all":p=r.AllLinks
    case "none":p=r.NoLinks
    case "coded":p=r.CodedUpper
    }
    n++
    if p.Answer==r.Truth{corr++}else{bad++}
    if p.Answer&&!r.Truth{fp++}
    if !p.Answer&&r.Truth{fn++}
   }
   worst=math.Max(worst,float64(bad)/float64(n))
  }
  scores[name]=map[string]float64{"correct":float64(corr),"accuracy":float64(corr)/1152,
    "fp":float64(fp),"fn":float64(fn),"worst_family_error":worst}
 }
 learned:=scores["learned"];coded:=scores["coded"]
 verdict:="INVALID_COMPETING_CLAIM_REFERENCE"
 if coded["correct"]==1152{
  if learned["accuracy"]>=.95&&learned["worst_family_error"]<=.05&&
   learned["accuracy"]-scores["first"]["accuracy"]>=.05&&
   learned["accuracy"]-scores["all"]["accuracy"]>=.05{
   verdict="DEVELOPMENT_LEARNED_COMPETING_ALIAS_CLAIM_SELECTION"
  }else{verdict="DEVELOPMENT_NO_COMPETING_ALIAS_CLAIM_SELECTION"}
 }
 return map[string]any{
  "schema":"wingless.f23.in-state-competing-identity-claims.dev.v1",
  "stage":"DEVELOPMENT_ONLY","parent_w22":"5a785822bf8fe7aa9878ab374dfd799a9d6f3633",
  "pretrained_trust_weights":policy.V,"retrained":false,
  "candidate_proposals_externally_supplied":true,
  "candidate_alias_pair_matching_coded":true,
  "three_candidates_per_case":true,"candidate_order_changed_per_case":true,
  "native_state_bytes":576,"state_claim_word_offsets":[]int{24,28,32},
  "new_seeds":[]int{88400,88406},"motifs_reused":[]int{108,109,110,111},
  "graph_count":96,"question_count":len(rows),"rows":rows,
  "scores":scores,"registered_verdict":verdict,
  "source_birth_independent":false,"sealed_first_use":false,
  "accepted_ref_mutation":false,"paid_api_calls":0,
  "elapsed_ms":time.Since(start).Milliseconds()},nil
}
