package unitary

import "math"

const UP74CPartitionedCarrierSchema = "wingless.up74c-partitioned-carrier.v1"

type UP74CPartitionedCarrierPoint struct {
	Arm                     string  `json:"arm"`
	Family                  string  `json:"family"`
	Setting                 string  `json:"setting"`
	PrimaryAccuracy         float64 `json:"primary_accuracy"`
	ExactEpisodeAccuracy    float64 `json:"exact_episode_accuracy"`
	CrossBankCorruptionRate float64 `json:"cross_bank_corruption_rate"`
	RecurrentStateBytes     int     `json:"recurrent_state_bytes"`
}

type UP74CPartitionedCarrierResult struct {
	Schema              string                         `json:"schema"`
	Experiment          string                         `json:"experiment"`
	SourceUPSQ0Seal     string                         `json:"source_up_sq0_seal"`
	SeedBases           []int                          `json:"seed_bases"`
	RecurrentStateBytes int                            `json:"recurrent_state_bytes"`
	ExactRecallUsed     bool                           `json:"exact_recall_used"`
	Points              []UP74CPartitionedCarrierPoint `json:"points"`
}

type up74cMachine struct {
	arm   string
	seed  uint64
	state [64]float64
}

func up74cComponent(arm string, seed uint64, key, value, index int) float64 {
	if arm == "flat64" {
		return sq0Component(seed, key, value, index)
	}
	bank := key & 3
	start := bank * 16
	if index < start || index >= start+16 {
		return 0
	}
	local := index - start
	x := seed
	x ^= uint64(key+1) * 0x9e3779b97f4a7c15
	x ^= uint64(value+1) * 0xbf58476d1ce4e5b9
	x ^= uint64(local+1) * 0x94d049bb133111eb
	x = sq0Mix64(x)
	scale := 1.0 / math.Sqrt(16)
	if x&1 == 0 {
		return -scale
	}
	return scale
}

func (m *up74cMachine) score(key, value int) float64 {
	sum := 0.0
	for i := 0; i < 64; i++ {
		sum += m.state[i] * up74cComponent(m.arm, m.seed, key, value, i)
	}
	return sum
}

func (m *up74cMachine) decode(key, vocab int) (best int, bestScore, secondScore float64, present bool) {
	best = 0
	bestScore = math.Inf(-1)
	secondScore = math.Inf(-1)
	maxAbs := 0.0
	for v := 0; v < vocab; v++ {
		s := m.score(key, v)
		if math.Abs(s) > maxAbs {
			maxAbs = math.Abs(s)
		}
		if s > bestScore {
			secondScore = bestScore
			bestScore = s
			best = v
		} else if s > secondScore {
			secondScore = s
		}
	}
	if math.IsInf(secondScore, -1) {
		secondScore = 0
	}
	present = maxAbs >= sq0PresenceThreshold
	return
}

func (m *up74cMachine) add(key, value int, scale float64) {
	for i := 0; i < 64; i++ {
		m.state[i] += scale * up74cComponent(m.arm, m.seed, key, value, i)
	}
}

func (m *up74cMachine) write(key, value, vocab int) {
	old, _, _, present := m.decode(key, vocab)
	if present {
		m.add(key, old, -1)
	}
	m.add(key, value, 1)
	intended := m.score(key, value)
	competing := math.Inf(-1)
	for v := 0; v < vocab; v++ {
		if v == value {
			continue
		}
		s := m.score(key, v)
		if s > competing {
			competing = s
		}
	}
	if math.IsInf(competing, -1) {
		competing = 0
	}
	den := math.Max(1, math.Abs(intended)+math.Abs(competing))
	margin := (intended - competing) / den
	if margin < sq0CorrectionThreshold {
		m.add(key, value, 1)
	}
}

func (m *up74cMachine) query(key, vocab int) int {
	best, _, _, _ := m.decode(key, vocab)
	return best
}

func up74cMultiBank(arm string, banks int, load float64, seedBases []int) UP74CPartitionedCarrierPoint {
	hits,total,exactHits,episodes := 0,0,0,0
	crossErr,crossTotal := 0,0
	totalSlots := banks*4
	occupied := int(math.Round(float64(totalSlots)*load))
	if occupied < 1 { occupied=1 }
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,301+banks*11+occupied,ep)
			rng:=newSQ0RNG(seed)
			m:=&up74cMachine{arm:arm,seed:seed}
			truth:=map[int]int{}
			order:=make([]int,0,occupied)
			for i:=0;i<occupied;i++ {
				key:=(ep+i*5)%totalSlots
				for {
					if _,ok:=truth[key];!ok {break}
					key=(key+1)%totalSlots
				}
				v:=rng.intn(16)
				truth[key]=v
				order=append(order,key)
				m.write(key,v,16)
				currentBank:=key/4
				for k,want:=range truth {
					if k/4==currentBank {continue}
					crossTotal++
					if m.query(k,16)!=want {crossErr++}
				}
			}
			exact:=true
			for _,key:=range order {
				total++
				if m.query(key,16)==truth[key] {hits++} else {exact=false}
			}
			episodes++
			if exact {exactHits++}
		}
	}
	rate:=0.0
	if crossTotal>0 {rate=float64(crossErr)/float64(crossTotal)}
	return UP74CPartitionedCarrierPoint{
		Arm:arm,Family:"multi_bank_interference",
		Setting:"banks="+itoa(banks)+",load="+formatLoad(load),
		PrimaryAccuracy:float64(hits)/float64(total),
		ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),
		CrossBankCorruptionRate:rate,RecurrentStateBytes:512,
	}
}

func formatLoad(v float64) string {
	if v==0.5 {return "0.50"}
	return "1.00"
}

func up74cAssociative(arm string, load int, seedBases []int) UP74CPartitionedCarrierPoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,337+load,ep)
			rng:=newSQ0RNG(seed)
			m:=&up74cMachine{arm:arm,seed:seed}
			truth:=make([]int,load)
			for key:=0;key<load;key++ {
				v:=rng.intn(32)
				truth[key]=v
				m.write(key,v,32)
			}
			exact:=true
			for q:=0;q<4;q++ {
				key:=(ep*3+q*7)%load
				total++
				if m.query(key,32)==truth[key] {hits++} else {exact=false}
			}
			episodes++
			if exact {exactHits++}
		}
	}
	return UP74CPartitionedCarrierPoint{
		Arm:arm,Family:"associative_recall",Setting:"load="+itoa(load),
		PrimaryAccuracy:float64(hits)/float64(total),
		ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),
		RecurrentStateBytes:512,
	}
}

func up74cOverwrite(arm string, writes int, seedBases []int) UP74CPartitionedCarrierPoint {
	hits,total,exactHits,episodes:=0,0,0,0
	for _,base:=range seedBases {
		for ep:=0;ep<96;ep++ {
			seed:=sq0Seed(base,359+writes,ep)
			rng:=newSQ0RNG(seed)
			m:=&up74cMachine{arm:arm,seed:seed}
			var truth [8]int
			for step:=0;step<writes;step++ {
				slot:=(step+ep)%8
				v:=(truth[slot]+1+rng.intn(15))%16
				truth[slot]=v
				m.write(slot,v,16)
			}
			exact:=true
			for slot:=0;slot<8;slot++ {
				total++
				if m.query(slot,16)==truth[slot] {hits++} else {exact=false}
			}
			episodes++
			if exact {exactHits++}
		}
	}
	return UP74CPartitionedCarrierPoint{
		Arm:arm,Family:"overwrite_latest_value_wins",Setting:"writes="+itoa(writes),
		PrimaryAccuracy:float64(hits)/float64(total),
		ExactEpisodeAccuracy:float64(exactHits)/float64(episodes),
		RecurrentStateBytes:512,
	}
}

func RunUP74C()(UP74CPartitionedCarrierResult,error){
	seedBases:=[]int{125000000,126000000}
	result:=UP74CPartitionedCarrierResult{
		Schema:UP74CPartitionedCarrierSchema,
		Experiment:"UP-74C-partitioned-carrier",
		SourceUPSQ0Seal:"8de9e11f26cc44c75f51d26cb606f97e35816479",
		SeedBases:append([]int(nil),seedBases...),
		RecurrentStateBytes:512,ExactRecallUsed:false,
	}
	for _,arm:=range []string{"flat64","partition4x16"} {
		for _,banks:=range []int{4,6,8} {
			for _,load:=range []float64{0.5,1.0} {
				result.Points=append(result.Points,up74cMultiBank(arm,banks,load,seedBases))
			}
		}
		for _,load:=range []int{8,16,24,32} {
			result.Points=append(result.Points,up74cAssociative(arm,load,seedBases))
		}
		for _,writes:=range []int{16,32,64} {
			result.Points=append(result.Points,up74cOverwrite(arm,writes,seedBases))
		}
	}
	return result,nil
}
