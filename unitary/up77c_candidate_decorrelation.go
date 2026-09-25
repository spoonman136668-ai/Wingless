package unitary

import "math"

const UP77CCandidateDecorrelationSchema = "wingless.up77c-candidate-decorrelation.v1"

type UP77CPoint struct {
	Arm                     string  `json:"arm"`
	Family                  string  `json:"family"`
	Setting                 string  `json:"setting"`
	PrimaryAccuracy         float64 `json:"primary_accuracy"`
	ExactEpisodeAccuracy    float64 `json:"exact_episode_accuracy"`
	CrossBankCorruptionRate float64 `json:"cross_bank_corruption_rate"`
	RecurrentStateBytes     int     `json:"recurrent_state_bytes"`
}

type UP77CCandidateDecorrelationResult struct {
	Schema               string       `json:"schema"`
	Experiment           string       `json:"experiment"`
	SourceUP76CSeal      string       `json:"source_up76c_seal"`
	SourceUP75CSeal      string       `json:"source_up75c_seal"`
	SeedBases            []int        `json:"seed_bases"`
	RecurrentStateBytes  int          `json:"recurrent_state_bytes"`
	ExactRecallUsed      bool         `json:"exact_recall_used"`
	PersistentRouteTable bool         `json:"persistent_route_table"`
	Points               []UP77CPoint `json:"points"`
}

type up77cMachine struct {
	arm string
	seed uint64
	state [64]float64
}

func up77cCandidates(arm string,key int)(int,int){
	a:=key&7
	switch arm{
	case "adjacent_pair":
		return a,a^1
	case "antipodal_pair":
		return a,(a+4)&7
	case "mixed_second_hash":
		h:=sq0Mix64(uint64(key+1)*0x9e3779b97f4a7c15)
		b:=int((h>>61)&7)
		if b==a{b=(b+3)&7}
		return a,b
	default:
		return a,a
	}
}

func up77cComponent(seed uint64,key,value,part,index int) float64 {
	start:=part*8
	if index<start||index>=start+8{return 0}
	local:=index-start
	x:=seed^uint64(key+1)*0x9e3779b97f4a7c15^uint64(value+1)*0xbf58476d1ce4e5b9^uint64(local+1)*0x94d049bb133111eb
	x=sq0Mix64(x)
	scale:=1/math.Sqrt(8)
	if x&1==0{return -scale}
	return scale
}

func (m *up77cMachine) norm(part int) float64 {
	sum:=0.0
	for i:=part*8;i<part*8+8;i++{sum+=m.state[i]*m.state[i]}
	return sum
}

func (m *up77cMachine) score(part,key,value int) float64 {
	sum:=0.0
	for i:=part*8;i<part*8+8;i++{sum+=m.state[i]*up77cComponent(m.seed,key,value,part,i)}
	return sum
}

func (m *up77cMachine) decode(part,key,vocab int)(best int,bestScore,second float64,present bool){
	best=0;bestScore=math.Inf(-1);second=math.Inf(-1);maxAbs:=0.0
	for v:=0;v<vocab;v++{
		s:=m.score(part,key,v)
		if math.Abs(s)>maxAbs{maxAbs=math.Abs(s)}
		if s>bestScore{second=bestScore;bestScore=s;best=v}else if s>second{second=s}
	}
	if math.IsInf(second,-1){second=0}
	present=maxAbs>=sq0PresenceThreshold
	return
}

func (m *up77cMachine) add(part,key,value int,scale float64){
	for i:=part*8;i<part*8+8;i++{m.state[i]+=scale*up77cComponent(m.seed,key,value,part,i)}
}

func (m *up77cMachine) chooseWritePart(key,vocab int)(part,old int,present bool){
	a,b:=up77cCandidates(m.arm,key)
	if m.arm=="fixed_mod8"{
		old,_,_,present=m.decode(a,key,vocab);return a,old,present
	}
	oldA,scoreA,_,presentA:=m.decode(a,key,vocab)
	oldB,scoreB,_,presentB:=m.decode(b,key,vocab)
	if presentA||presentB{
		if presentA&&!presentB{return a,oldA,true}
		if presentB&&!presentA{return b,oldB,true}
		if scoreA>scoreB||(scoreA==scoreB&&a<b){return a,oldA,true}
		return b,oldB,true
	}
	nA,nB:=m.norm(a),m.norm(b)
	if nA<nB||(nA==nB&&a<b){return a,oldA,false}
	return b,oldB,false
}

func (m *up77cMachine) write(key,value,vocab int){
	part,old,present:=m.chooseWritePart(key,vocab)
	if present{m.add(part,key,old,-1)}
	m.add(part,key,value,1)
	intended:=m.score(part,key,value);competing:=math.Inf(-1)
	for v:=0;v<vocab;v++{if v==value{continue};s:=m.score(part,key,v);if s>competing{competing=s}}
	if math.IsInf(competing,-1){competing=0}
	den:=math.Max(1,math.Abs(intended)+math.Abs(competing))
	if (intended-competing)/den<sq0CorrectionThreshold{m.add(part,key,value,1)}
}

func (m *up77cMachine) query(key,vocab int) int {
	a,b:=up77cCandidates(m.arm,key)
	va,sa,_,_:=m.decode(a,key,vocab)
	if m.arm=="fixed_mod8"{return va}
	vb,sb,_,_:=m.decode(b,key,vocab)
	if sa>sb||(sa==sb&&a<b){return va}
	return vb
}

func up77cMultiBank(arm string,banks int,load float64,seedBases []int) UP77CPoint {
	hits,total,exactHits,episodes:=0,0,0,0;crossErr,crossTotal:=0,0
	totalSlots:=banks*4;occupied:=int(math.Round(float64(totalSlots)*load));if occupied<1{occupied=1}
	for _,base:=range seedBases{
		for ep:=0;ep<96;ep++{
			seed:=sq0Seed(base,901+banks*13+occupied,ep);rng:=newSQ0RNG(seed);m:=&up77cMachine{arm:arm,seed:seed}
			truth:=map[int]int{};order:=make([]int,0,occupied)
			for i:=0;i<occupied;i++{
				key:=(ep+i*5)%totalSlots
				for{if _,ok:=truth[key];!ok{break};key=(key+1)%totalSlots}
				v:=rng.intn(16);truth[key]=v;order=append(order,key);m.write(key,v,16)
				currentBank:=key/4
				for k,want:=range truth{if k/4==currentBank{continue};crossTotal++;if m.query(k,16)!=want{crossErr++}}
			}
			exact:=true
			for _,k:=range order{total++;if m.query(k,16)==truth[k]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	rate:=0.0;if crossTotal>0{rate=float64(crossErr)/float64(crossTotal)}
	return UP77CPoint{Arm:arm,Family:"multi_bank_interference",Setting:"banks="+up74cItoa(banks)+",load="+formatLoad(load),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),CrossBankCorruptionRate:rate,RecurrentStateBytes:512}
}

func up77cAssociative(arm string,load int,seedBases []int) UP77CPoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases{
		for ep:=0;ep<96;ep++{
			seed:=sq0Seed(base,941+load,ep);rng:=newSQ0RNG(seed);m:=&up77cMachine{arm:arm,seed:seed};truth:=make([]int,load)
			for k:=0;k<load;k++{v:=rng.intn(32);truth[k]=v;m.write(k,v,32)}
			exact:=true
			for q:=0;q<4;q++{k:=(ep*3+q*7)%load;total++;if m.query(k,32)==truth[k]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	return UP77CPoint{Arm:arm,Family:"associative_recall",Setting:"load="+up74cItoa(load),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),RecurrentStateBytes:512}
}

func up77cOverwrite(arm string,writes int,seedBases []int) UP77CPoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases{
		for ep:=0;ep<96;ep++{
			seed:=sq0Seed(base,977+writes,ep);rng:=newSQ0RNG(seed);m:=&up77cMachine{arm:arm,seed:seed};var truth [8]int
			for step:=0;step<writes;step++{slot:=(step+ep)%8;v:=(truth[slot]+1+rng.intn(15))%16;truth[slot]=v;m.write(slot,v,16)}
			exact:=true
			for slot:=0;slot<8;slot++{total++;if m.query(slot,16)==truth[slot]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	return UP77CPoint{Arm:arm,Family:"overwrite_latest_value_wins",Setting:"writes="+up74cItoa(writes),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),RecurrentStateBytes:512}
}

func RunUP77C()(UP77CCandidateDecorrelationResult,error){
	seedBases:=[]int{137000000,138000000}
	result:=UP77CCandidateDecorrelationResult{
		Schema:UP77CCandidateDecorrelationSchema,Experiment:"UP-77C-candidate-decorrelation",
		SourceUP76CSeal:"194f014c9f3b2e7769905b9237229c2112d613af",
		SourceUP75CSeal:"478c60bef8d4104600043e875faeb8de1ea3e409",
		SeedBases:append([]int(nil),seedBases...),RecurrentStateBytes:512,ExactRecallUsed:false,PersistentRouteTable:false,
	}
	for _,arm:=range []string{"fixed_mod8","adjacent_pair","antipodal_pair","mixed_second_hash"}{
		for _,banks:=range []int{6,8}{for _,load:=range []float64{0.5,1.0}{result.Points=append(result.Points,up77cMultiBank(arm,banks,load,seedBases))}}
		for _,load:=range []int{16,24,32}{result.Points=append(result.Points,up77cAssociative(arm,load,seedBases))}
		for _,writes:=range []int{32,64}{result.Points=append(result.Points,up77cOverwrite(arm,writes,seedBases))}
	}
	return result,nil
}
