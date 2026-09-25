package unitary

import "math"

const UP76CPartitionGranularitySchema = "wingless.up76c-partition-granularity.v1"

type UP76CPartitionPoint struct {
	Arm                     string  `json:"arm"`
	Family                  string  `json:"family"`
	Setting                 string  `json:"setting"`
	PrimaryAccuracy         float64 `json:"primary_accuracy"`
	ExactEpisodeAccuracy    float64 `json:"exact_episode_accuracy"`
	CrossBankCorruptionRate float64 `json:"cross_bank_corruption_rate"`
	RecurrentStateBytes     int     `json:"recurrent_state_bytes"`
}

type UP76CPartitionGranularityResult struct {
	Schema              string                 `json:"schema"`
	Experiment          string                 `json:"experiment"`
	SourceUP74CSeal     string                 `json:"source_up74c_seal"`
	SourceUP75CSeal     string                 `json:"source_up75c_seal"`
	SeedBases           []int                  `json:"seed_bases"`
	RecurrentStateBytes int                    `json:"recurrent_state_bytes"`
	ExactRecallUsed     bool                   `json:"exact_recall_used"`
	DynamicRoutingUsed  bool                   `json:"dynamic_routing_used"`
	Points              []UP76CPartitionPoint  `json:"points"`
}

type up76cMachine struct {
	partitions int
	seed uint64
	state [64]float64
}

func (m *up76cMachine) dim() int { return 64/m.partitions }

func (m *up76cMachine) part(key int) int {
	p:=key%m.partitions
	if p<0 {p+=m.partitions}
	return p
}

func (m *up76cMachine) component(key,value,part,index int) float64 {
	d:=m.dim();start:=part*d
	if index<start || index>=start+d {return 0}
	local:=index-start
	x:=m.seed
	x^=uint64(key+1)*0x9e3779b97f4a7c15
	x^=uint64(value+1)*0xbf58476d1ce4e5b9
	x^=uint64(local+1)*0x94d049bb133111eb
	x=sq0Mix64(x)
	scale:=1/math.Sqrt(float64(d))
	if x&1==0 {return -scale}
	return scale
}

func (m *up76cMachine) score(key,value int) float64 {
	part:=m.part(key);d:=m.dim();start:=part*d;sum:=0.0
	for i:=start;i<start+d;i++ {sum+=m.state[i]*m.component(key,value,part,i)}
	return sum
}

func (m *up76cMachine) decode(key,vocab int)(best int,bestScore,second float64,present bool){
	best=0;bestScore=math.Inf(-1);second=math.Inf(-1);maxAbs:=0.0
	for v:=0;v<vocab;v++ {
		s:=m.score(key,v)
		if math.Abs(s)>maxAbs {maxAbs=math.Abs(s)}
		if s>bestScore {second=bestScore;bestScore=s;best=v} else if s>second {second=s}
	}
	if math.IsInf(second,-1){second=0}
	present=maxAbs>=sq0PresenceThreshold
	return
}

func (m *up76cMachine) add(key,value int,scale float64){
	part:=m.part(key);d:=m.dim();start:=part*d
	for i:=start;i<start+d;i++ {m.state[i]+=scale*m.component(key,value,part,i)}
}

func (m *up76cMachine) write(key,value,vocab int){
	old,_,_,present:=m.decode(key,vocab)
	if present {m.add(key,old,-1)}
	m.add(key,value,1)
	intended:=m.score(key,value);competing:=math.Inf(-1)
	for v:=0;v<vocab;v++ {if v==value{continue};s:=m.score(key,v);if s>competing{competing=s}}
	if math.IsInf(competing,-1){competing=0}
	den:=math.Max(1,math.Abs(intended)+math.Abs(competing))
	if (intended-competing)/den<sq0CorrectionThreshold {m.add(key,value,1)}
}

func (m *up76cMachine) query(key,vocab int) int {best,_,_,_:=m.decode(key,vocab);return best}

func up76cArmName(parts int) string {
	switch parts {case 2:return "partition2x32";case 4:return "partition4x16";default:return "partition8x8"}
}

func up76cMultiBank(parts,banks int,load float64,seedBases []int) UP76CPartitionPoint {
	hits,total,exactHits,episodes:=0,0,0,0;crossErr,crossTotal:=0,0
	totalSlots:=banks*4;occupied:=int(math.Round(float64(totalSlots)*load));if occupied<1{occupied=1}
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,701+parts*17+banks*13+occupied,ep);rng:=newSQ0RNG(seed);m:=&up76cMachine{partitions:parts,seed:seed}
			truth:=map[int]int{};order:=make([]int,0,occupied)
			for i:=0;i<occupied;i++ {
				key:=(ep+i*5)%totalSlots
				for {if _,ok:=truth[key];!ok{break};key=(key+1)%totalSlots}
				v:=rng.intn(16);truth[key]=v;order=append(order,key);m.write(key,v,16)
				currentBank:=key/4
				for k,want:=range truth {if k/4==currentBank{continue};crossTotal++;if m.query(k,16)!=want{crossErr++}}
			}
			exact:=true
			for _,k:=range order {total++;if m.query(k,16)==truth[k]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	rate:=0.0;if crossTotal>0{rate=float64(crossErr)/float64(crossTotal)}
	return UP76CPartitionPoint{Arm:up76cArmName(parts),Family:"multi_bank_interference",Setting:"banks="+up74cItoa(banks)+",load="+formatLoad(load),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),CrossBankCorruptionRate:rate,RecurrentStateBytes:512}
}

func up76cAssociative(parts,load int,seedBases []int) UP76CPartitionPoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,743+parts*19+load,ep);rng:=newSQ0RNG(seed);m:=&up76cMachine{partitions:parts,seed:seed};truth:=make([]int,load)
			for k:=0;k<load;k++ {v:=rng.intn(32);truth[k]=v;m.write(k,v,32)}
			exact:=true
			for q:=0;q<4;q++ {k:=(ep*3+q*7)%load;total++;if m.query(k,32)==truth[k]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	return UP76CPartitionPoint{Arm:up76cArmName(parts),Family:"associative_recall",Setting:"load="+up74cItoa(load),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),RecurrentStateBytes:512}
}

func up76cOverwrite(parts,writes int,seedBases []int) UP76CPartitionPoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,769+parts*23+writes,ep);rng:=newSQ0RNG(seed);m:=&up76cMachine{partitions:parts,seed:seed};var truth [8]int
			for step:=0;step<writes;step++ {slot:=(step+ep)%8;v:=(truth[slot]+1+rng.intn(15))%16;truth[slot]=v;m.write(slot,v,16)}
			exact:=true
			for slot:=0;slot<8;slot++ {total++;if m.query(slot,16)==truth[slot]{hits++}else{exact=false}}
			episodes++;if exact{exactHits++}
		}
	}
	return UP76CPartitionPoint{Arm:up76cArmName(parts),Family:"overwrite_latest_value_wins",Setting:"writes="+up74cItoa(writes),PrimaryAccuracy:float64(hits)/float64(total),ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),RecurrentStateBytes:512}
}

func RunUP76C()(UP76CPartitionGranularityResult,error){
	seedBases:=[]int{133000000,134000000}
	result:=UP76CPartitionGranularityResult{
		Schema:UP76CPartitionGranularitySchema,Experiment:"UP-76C-partition-granularity",
		SourceUP74CSeal:"9b00b8faac8d348bd6a6aa96b5f0437bfad0b505",
		SourceUP75CSeal:"478c60bef8d4104600043e875faeb8de1ea3e409",
		SeedBases:append([]int(nil),seedBases...),RecurrentStateBytes:512,ExactRecallUsed:false,DynamicRoutingUsed:false,
	}
	for _,parts:=range []int{2,4,8} {
		for _,banks:=range []int{6,8} {for _,load:=range []float64{0.5,1.0} {result.Points=append(result.Points,up76cMultiBank(parts,banks,load,seedBases))}}
		for _,load:=range []int{16,24,32} {result.Points=append(result.Points,up76cAssociative(parts,load,seedBases))}
		for _,writes:=range []int{32,64} {result.Points=append(result.Points,up76cOverwrite(parts,writes,seedBases))}
	}
	return result,nil
}
