package unitary

import "math"

const UP75CTwoChoicePartitionRoutingSchema = "wingless.up75c-two-choice-partition-routing.v1"

type UP75CTwoChoicePoint struct {
	Arm                     string  `json:"arm"`
	Family                  string  `json:"family"`
	Setting                 string  `json:"setting"`
	PrimaryAccuracy         float64 `json:"primary_accuracy"`
	ExactEpisodeAccuracy    float64 `json:"exact_episode_accuracy"`
	CrossBankCorruptionRate float64 `json:"cross_bank_corruption_rate"`
	RecurrentStateBytes     int     `json:"recurrent_state_bytes"`
}

type UP75CTwoChoicePartitionRoutingResult struct {
	Schema              string                 `json:"schema"`
	Experiment          string                 `json:"experiment"`
	SourceUP74CSeal     string                 `json:"source_up74c_seal"`
	SeedBases           []int                  `json:"seed_bases"`
	RecurrentStateBytes int                    `json:"recurrent_state_bytes"`
	ExactRecallUsed     bool                   `json:"exact_recall_used"`
	PersistentRouteTable bool                  `json:"persistent_route_table"`
	Points              []UP75CTwoChoicePoint  `json:"points"`
}

type up75cMachine struct {
	arm string
	seed uint64
	state [64]float64
}

func up75cCandidates(key int)(int,int){
	a:=key&3
	b:=(3*key+1)&3
	if b==a { b=(a+1)&3 }
	return a,b
}

func up75cComponent(seed uint64,key,value,part,index int) float64 {
	start:=part*16
	if index<start || index>=start+16 { return 0 }
	local:=index-start
	x:=seed
	x^=uint64(key+1)*0x9e3779b97f4a7c15
	x^=uint64(value+1)*0xbf58476d1ce4e5b9
	x^=uint64(local+1)*0x94d049bb133111eb
	x=sq0Mix64(x)
	if x&1==0 { return -0.25 }
	return 0.25
}

func (m *up75cMachine) partNorm(part int) float64 {
	sum:=0.0
	start:=part*16
	for i:=start;i<start+16;i++ { sum+=m.state[i]*m.state[i] }
	return sum
}

func (m *up75cMachine) scorePart(part,key,value int) float64 {
	sum:=0.0
	for i:=part*16;i<part*16+16;i++ { sum+=m.state[i]*up75cComponent(m.seed,key,value,part,i) }
	return sum
}

func (m *up75cMachine) decodePart(part,key,vocab int)(best int,bestScore,second float64,present bool){
	best=0;bestScore=math.Inf(-1);second=math.Inf(-1);maxAbs:=0.0
	for v:=0;v<vocab;v++ {
		s:=m.scorePart(part,key,v)
		if math.Abs(s)>maxAbs { maxAbs=math.Abs(s) }
		if s>bestScore { second=bestScore;bestScore=s;best=v } else if s>second { second=s }
	}
	if math.IsInf(second,-1){second=0}
	present=maxAbs>=sq0PresenceThreshold
	return
}

func (m *up75cMachine) addPart(part,key,value int,scale float64){
	for i:=part*16;i<part*16+16;i++ { m.state[i]+=scale*up75cComponent(m.seed,key,value,part,i) }
}

func (m *up75cMachine) choosePartForWrite(key,vocab int)(part,old int,present bool){
	a,b:=up75cCandidates(key)
	if m.arm=="fixed_mod4" {
		old,_,_,present=m.decodePart(a,key,vocab)
		return a,old,present
	}
	oldA,scoreA,_,presentA:=m.decodePart(a,key,vocab)
	oldB,scoreB,_,presentB:=m.decodePart(b,key,vocab)
	if presentA || presentB {
		if presentA && !presentB { return a,oldA,true }
		if presentB && !presentA { return b,oldB,true }
		if scoreA>scoreB || (scoreA==scoreB && a<b) { return a,oldA,true }
		return b,oldB,true
	}
	nA,nB:=m.partNorm(a),m.partNorm(b)
	if nA<nB || (nA==nB && a<b) { return a,oldA,false }
	return b,oldB,false
}

func (m *up75cMachine) write(key,value,vocab int){
	part,old,present:=m.choosePartForWrite(key,vocab)
	if present { m.addPart(part,key,old,-1) }
	m.addPart(part,key,value,1)
	intended:=m.scorePart(part,key,value)
	competing:=math.Inf(-1)
	for v:=0;v<vocab;v++ {
		if v==value {continue}
		s:=m.scorePart(part,key,v)
		if s>competing {competing=s}
	}
	if math.IsInf(competing,-1){competing=0}
	den:=math.Max(1,math.Abs(intended)+math.Abs(competing))
	if (intended-competing)/den < sq0CorrectionThreshold { m.addPart(part,key,value,1) }
}

func (m *up75cMachine) query(key,vocab int) int {
	a,b:=up75cCandidates(key)
	va,sa,_,_:=m.decodePart(a,key,vocab)
	if m.arm=="fixed_mod4" { return va }
	vb,sb,_,_:=m.decodePart(b,key,vocab)
	if sa>sb || (sa==sb && a<b) { return va }
	return vb
}

func up75cMultiBank(arm string,banks int,load float64,seedBases []int) UP75CTwoChoicePoint {
	hits,total,exactHits,episodes:=0,0,0,0
	crossErr,crossTotal:=0,0
	totalSlots:=banks*4
	occupied:=int(math.Round(float64(totalSlots)*load))
	if occupied<1 {occupied=1}
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,501+banks*13+occupied,ep)
			rng:=newSQ0RNG(seed)
			m:=&up75cMachine{arm:arm,seed:seed}
			truth:=map[int]int{}
			order:=make([]int,0,occupied)
			for i:=0;i<occupied;i++ {
				key:=(ep+i*5)%totalSlots
				for { if _,ok:=truth[key];!ok {break}; key=(key+1)%totalSlots }
				v:=rng.intn(16);truth[key]=v;order=append(order,key);m.write(key,v,16)
				currentBank:=key/4
				for k,want:=range truth {
					if k/4==currentBank {continue}
					crossTotal++
					if m.query(k,16)!=want {crossErr++}
				}
			}
			exact:=true
			for _,k:=range order { total++;if m.query(k,16)==truth[k]{hits++}else{exact=false} }
			episodes++;if exact{exactHits++}
		}
	}
	rate:=0.0;if crossTotal>0{rate=float64(crossErr)/float64(crossTotal)}
	return UP75CTwoChoicePoint{Arm:arm,Family:"multi_bank_interference",Setting:"banks="+up74cItoa(banks)+",load="+formatLoad(load),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),CrossBankCorruptionRate:rate,RecurrentStateBytes:512}
}

func up75cAssociative(arm string,load int,seedBases []int) UP75CTwoChoicePoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,541+load,ep);rng:=newSQ0RNG(seed);m:=&up75cMachine{arm:arm,seed:seed}
			truth:=make([]int,load)
			for k:=0;k<load;k++ {v:=rng.intn(32);truth[k]=v;m.write(k,v,32)}
			exact:=true
			for q:=0;q<4;q++ {k:=(ep*3+q*7)%load;total++;if m.query(k,32)==truth[k]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	return UP75CTwoChoicePoint{Arm:arm,Family:"associative_recall",Setting:"load="+up74cItoa(load),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),RecurrentStateBytes:512}
}

func up75cOverwrite(arm string,writes int,seedBases []int) UP75CTwoChoicePoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,577+writes,ep);rng:=newSQ0RNG(seed);m:=&up75cMachine{arm:arm,seed:seed}
			var truth [8]int
			for step:=0;step<writes;step++ {slot:=(step+ep)%8;v:=(truth[slot]+1+rng.intn(15))%16;truth[slot]=v;m.write(slot,v,16)}
			exact:=true
			for slot:=0;slot<8;slot++ {total++;if m.query(slot,16)==truth[slot]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	return UP75CTwoChoicePoint{Arm:arm,Family:"overwrite_latest_value_wins",Setting:"writes="+up74cItoa(writes),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),RecurrentStateBytes:512}
}

func RunUP75C()(UP75CTwoChoicePartitionRoutingResult,error){
	seedBases:=[]int{129000000,130000000}
	result:=UP75CTwoChoicePartitionRoutingResult{
		Schema:UP75CTwoChoicePartitionRoutingSchema,
		Experiment:"UP-75C-two-choice-partition-routing",
		SourceUP74CSeal:"9b00b8faac8d348bd6a6aa96b5f0437bfad0b505",
		SeedBases:append([]int(nil),seedBases...),
		RecurrentStateBytes:512,ExactRecallUsed:false,PersistentRouteTable:false,
	}
	for _,arm:=range []string{"fixed_mod4","two_choice_norm"} {
		for _,banks:=range []int{6,8} {for _,load:=range []float64{0.5,1.0} {result.Points=append(result.Points,up75cMultiBank(arm,banks,load,seedBases))}}
		for _,load:=range []int{16,24,32} {result.Points=append(result.Points,up75cAssociative(arm,load,seedBases))}
		for _,writes:=range []int{32,64} {result.Points=append(result.Points,up75cOverwrite(arm,writes,seedBases))}
	}
	return result,nil
}
