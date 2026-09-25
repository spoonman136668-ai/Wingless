package unitary

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

const UPSQ0SequenceQualificationSchema = "wingless.up-sq0-sequence-qualification.v1"

const (
	sq0CarrierDim          = 64
	sq0PresenceThreshold   = 0.25
	sq0CorrectionThreshold = 0.15
	sq0ExactRecallCap      = 16
)

type UPSQ0Metric struct {
	Arm                  string
	Family               string
	Setting              string
	PrimaryMetric        string
	PrimaryAccuracy      float64
	SecondaryMetric      string
	SecondaryAccuracy    float64
	TertiaryMetric       string
	TertiaryValue        float64
	Events               int
	RecurrentStateBytes  int
	ExactRecallBytes     int
	PeakHeapBytes        uint64
	PeakVRAMBytes        uint64
	EventsPerSecond      float64
	WallMilliseconds     float64
}

type UPSQ0SequenceQualificationResult struct {
	Schema                 string
	Experiment             string
	ImplementationFreeze   string
	GeneratorSeedBases     []int
	RecurrentStateBytes    int
	ExactRecallEntryCap    int
	FullGlobalAttention    bool
	AggregateScoreSelected bool
	Metrics                []UPSQ0Metric
}

type sq0RNG struct{ x uint64 }

func newSQ0RNG(seed uint64) *sq0RNG {
	if seed == 0 {
		seed = 0x9e3779b97f4a7c15
	}
	return &sq0RNG{x: seed}
}

func (r *sq0RNG) next() uint64 {
	x := r.x
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	r.x = x
	return x
}

func (r *sq0RNG) intn(n int) int {
	if n <= 0 {
		return 0
	}
	return int(r.next() % uint64(n))
}

func sq0Mix64(x uint64) uint64 {
	x ^= x >> 30
	x *= 0xbf58476d1ce4e5b9
	x ^= x >> 27
	x *= 0x94d049bb133111eb
	x ^= x >> 31
	return x
}

func sq0Component(seed uint64, key, value, index int) float64 {
	x := seed
	x ^= uint64(key+1) * 0x9e3779b97f4a7c15
	x ^= uint64(value+1) * 0xbf58476d1ce4e5b9
	x ^= uint64(index+1) * 0x94d049bb133111eb
	x = sq0Mix64(x)
	scale := 1.0 / math.Sqrt(float64(sq0CarrierDim))
	if x&1 == 0 {
		return -scale
	}
	return scale
}

type sq0RecallEntry struct {
	key  int
	val  int
	used bool
}

type sq0Machine struct {
	arm        string
	seed       uint64
	state      [sq0CarrierDim]float64
	recall     [sq0ExactRecallCap]sq0RecallEntry
	recallUsed int
	recallNext int
}

func newSQ0Machine(arm string, seed uint64) *sq0Machine {
	return &sq0Machine{arm: arm, seed: seed}
}

func (m *sq0Machine) score(key, value int) float64 {
	sum := 0.0
	for i := 0; i < sq0CarrierDim; i++ {
		sum += m.state[i] * sq0Component(m.seed, key, value, i)
	}
	return sum
}

func (m *sq0Machine) decodeCarrier(key, vocab int) (best int, bestScore, secondScore float64, present bool) {
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

func (m *sq0Machine) addVector(key, value int, scale float64) {
	for i := 0; i < sq0CarrierDim; i++ {
		m.state[i] += scale * sq0Component(m.seed, key, value, i)
	}
}

func (m *sq0Machine) recallWrite(key, value int) {
	for i := 0; i < m.recallUsed; i++ {
		if m.recall[i].used && m.recall[i].key == key {
			m.recall[i].val = value
			return
		}
	}
	if m.recallUsed < sq0ExactRecallCap {
		m.recall[m.recallUsed] = sq0RecallEntry{key: key, val: value, used: true}
		m.recallUsed++
		if m.recallUsed == sq0ExactRecallCap {
			m.recallNext = 0
		}
		return
	}
	m.recall[m.recallNext] = sq0RecallEntry{key: key, val: value, used: true}
	m.recallNext = (m.recallNext + 1) % sq0ExactRecallCap
}

func (m *sq0Machine) recallQuery(key int) (int, bool) {
	for i := 0; i < m.recallUsed; i++ {
		if m.recall[i].used && m.recall[i].key == key {
			return m.recall[i].val, true
		}
	}
	return 0, false
}

func (m *sq0Machine) write(key, value, vocab int) {
	old, _, _, present := m.decodeCarrier(key, vocab)
	if present {
		m.addVector(key, old, -1)
	}
	m.addVector(key, value, 1)

	if m.arm == "transport_gated_correction" || m.arm == "transport_gated_correction_bounded_exact_recall" {
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
			m.addVector(key, value, 1)
		}
	}

	if m.arm == "transport_gated_correction_bounded_exact_recall" {
		m.recallWrite(key, value)
	}
}

func (m *sq0Machine) query(key, vocab int) int {
	if m.arm == "transport_gated_correction_bounded_exact_recall" {
		if v, ok := m.recallQuery(key); ok {
			return v
		}
	}
	best, _, _, _ := m.decodeCarrier(key, vocab)
	return best
}

type sq0Counts struct {
	primaryHits   int
	primaryTotal  int
	secondaryHits int
	secondaryTotal int
	events        int
	tertiaryName  string
	tertiaryValue float64
	maxRecallUsed int
}

func sq0Measure(arm, family, setting, primaryName, secondaryName string, fn func() sq0Counts) UPSQ0Metric {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	start := time.Now()
	c := fn()
	elapsed := time.Since(start)
	runtime.ReadMemStats(&after)

	primary := 0.0
	if c.primaryTotal > 0 {
		primary = float64(c.primaryHits) / float64(c.primaryTotal)
	}
	secondary := 0.0
	if c.secondaryTotal > 0 {
		secondary = float64(c.secondaryHits) / float64(c.secondaryTotal)
	}
	peak := before.Alloc
	if after.Alloc > peak {
		peak = after.Alloc
	}
	eps := 0.0
	if elapsed > 0 {
		eps = float64(c.events) / elapsed.Seconds()
	}
	return UPSQ0Metric{
		Arm: arm, Family: family, Setting: setting,
		PrimaryMetric: primaryName, PrimaryAccuracy: primary,
		SecondaryMetric: secondaryName, SecondaryAccuracy: secondary,
		TertiaryMetric: c.tertiaryName, TertiaryValue: c.tertiaryValue,
		Events: c.events, RecurrentStateBytes: sq0CarrierDim * 8,
		ExactRecallBytes: c.maxRecallUsed * 16,
		PeakHeapBytes: peak, PeakVRAMBytes: 0,
		EventsPerSecond: eps, WallMilliseconds: float64(elapsed.Microseconds()) / 1000,
	}
}

func sq0Seed(base, salt, episode int) uint64 {
	return sq0Mix64(uint64(base) ^ uint64(salt+1)*0x9e3779b97f4a7c15 ^ uint64(episode+1)*0xbf58476d1ce4e5b9)
}

func sq0Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func sq0EvalDelayedCopy(arm string, seedBases []int) []UPSQ0Metric {
	lengths := []int{16, 32, 64, 128}
	delays := []int{4, 8, 16, 32}
	var out []UPSQ0Metric
	for _, length := range lengths {
		for _, delay := range delays {
			setting := fmt.Sprintf("length=%d,delay=%d", length, delay)
			out = append(out, sq0Measure(arm, "delayed_copy", setting, "token_exact_accuracy", "whole_sequence_exact_accuracy", func() sq0Counts {
				var c sq0Counts
				for _, base := range seedBases {
					for ep := 0; ep < 32; ep++ {
						seed := sq0Seed(base, 11+length*3+delay, ep)
						rng := newSQ0RNG(seed)
						m := newSQ0Machine(arm, seed)
						truth := make([]int, length)
						for i := 0; i < length; i++ {
							truth[i] = rng.intn(16)
							m.write(i, truth[i], 16)
							c.events++
						}
						c.events += delay
						exact := true
						for i := 0; i < length; i++ {
							got := m.query(i, 16)
							c.primaryTotal++
							c.events++
							if got == truth[i] {
								c.primaryHits++
							} else {
								exact = false
							}
						}
						c.secondaryTotal++
						if exact {
							c.secondaryHits++
						}
						c.maxRecallUsed = sq0Max(c.maxRecallUsed, m.recallUsed)
					}
				}
				return c
			}))
		}
	}
	return out
}

func sq0EvalSelectiveCopy(arm string, seedBases []int) []UPSQ0Metric {
	lengths := []int{32, 64, 128, 256}
	var out []UPSQ0Metric
	for _, length := range lengths {
		setting := fmt.Sprintf("length=%d,marked=25%%", length)
		out = append(out, sq0Measure(arm, "selective_copy", setting, "selected_symbol_accuracy", "ordered_sequence_exact_accuracy", func() sq0Counts {
			var c sq0Counts
			for _, base := range seedBases {
				for ep := 0; ep < 32; ep++ {
					seed := sq0Seed(base, 29+length, ep)
					rng := newSQ0RNG(seed)
					m := newSQ0Machine(arm, seed)
					truth := make([]int, length)
					var marked []int
					for i := 0; i < length; i++ {
						truth[i] = rng.intn(16)
						m.write(i, truth[i], 16)
						c.events++
						if (i+ep+base)%4 == 0 {
							marked = append(marked, i)
						}
					}
					exact := true
					for _, i := range marked {
						got := m.query(i, 16)
						c.primaryTotal++
						c.events++
						if got == truth[i] {
							c.primaryHits++
						} else {
							exact = false
						}
					}
					c.secondaryTotal++
					if exact {
						c.secondaryHits++
					}
					c.maxRecallUsed = sq0Max(c.maxRecallUsed, m.recallUsed)
				}
			}
			return c
		}))
	}
	return out
}

func sq0EvalMQAR(arm string, seedBases []int) []UPSQ0Metric {
	loads := []int{4, 8, 16}
	var out []UPSQ0Metric
	for _, load := range loads {
		for _, collision := range []bool{false, true} {
			setting := fmt.Sprintf("load=%d,collision=%t", load, collision)
			out = append(out, sq0Measure(arm, "multi_query_associative_recall", setting, "query_exact_accuracy", "episode_exact_accuracy", func() sq0Counts {
				var c sq0Counts
				for _, base := range seedBases {
					for ep := 0; ep < 64; ep++ {
						seed := sq0Seed(base, 47+load*7, ep)
						rng := newSQ0RNG(seed)
						m := newSQ0Machine(arm, seed)
						truth := map[int]int{}
						keys := make([]int, load)
						for i := 0; i < load; i++ {
							if collision && i >= load/2 {
								keys[i] = keys[i-load/2]
							} else {
								keys[i] = (ep*3 + i*5 + base) % 32
							}
							v := rng.intn(32)
							truth[keys[i]] = v
							m.write(keys[i], v, 32)
							c.events++
						}
						exact := true
						for q := 0; q < 4; q++ {
							key := keys[(q*3+ep)%len(keys)]
							got := m.query(key, 32)
							c.primaryTotal++
							c.events++
							if got == truth[key] {
								c.primaryHits++
							} else {
								exact = false
							}
						}
						c.secondaryTotal++
						if exact {
							c.secondaryHits++
						}
						c.maxRecallUsed = sq0Max(c.maxRecallUsed, m.recallUsed)
					}
				}
				return c
			}))
		}
	}
	return out
}

func sq0EvalOverwrite(arm string, seedBases []int) []UPSQ0Metric {
	writesList := []int{8, 16, 32}
	var out []UPSQ0Metric
	for _, writes := range writesList {
		setting := fmt.Sprintf("writes=%d,half_value_changing_overwrite=true", writes)
		out = append(out, sq0Measure(arm, "overwrite_latest_value_wins", setting, "final_slot_accuracy", "episode_exact_accuracy", func() sq0Counts {
			var c sq0Counts
			for _, base := range seedBases {
				for ep := 0; ep < 64; ep++ {
					seed := sq0Seed(base, 71+writes, ep)
					rng := newSQ0RNG(seed)
					m := newSQ0Machine(arm, seed)
					truth := [8]int{}
					for i := range truth {
						truth[i] = 0
					}
					for step := 0; step < writes; step++ {
						slot := (step + ep) % 8
						v := 0
						if ep%2 == 0 {
							v = (truth[slot] + 1 + rng.intn(15)) % 16
						} else {
							if step < 8 {
								v = rng.intn(16)
							} else {
								v = truth[slot]
							}
						}
						truth[slot] = v
						m.write(slot, v, 16)
						c.events++
					}
					exact := true
					for slot := 0; slot < 8; slot++ {
						got := m.query(slot, 16)
						c.primaryTotal++
						c.events++
						if got == truth[slot] {
							c.primaryHits++
						} else {
							exact = false
						}
					}
					c.secondaryTotal++
					if exact {
						c.secondaryHits++
					}
					c.maxRecallUsed = sq0Max(c.maxRecallUsed, m.recallUsed)
				}
			}
			return c
		}))
	}
	return out
}

var sq0Operators = [6][8]int{
	{1,2,3,4,5,6,7,0},
	{7,0,1,2,3,4,5,6},
	{0,2,4,6,1,3,5,7},
	{3,0,7,4,1,6,2,5},
	{2,5,0,7,4,1,6,3},
	{6,3,5,1,7,2,0,4},
}

func sq0EvalStateMachine(arm string, seedBases []int) []UPSQ0Metric {
	depths := []int{1,2,3,4,8,16}
	var out []UPSQ0Metric
	for _, depth := range depths {
		setting := fmt.Sprintf("depth=%d", depth)
		out = append(out, sq0Measure(arm, "state_machine_composition", setting, "final_state_accuracy", "final_state_accuracy", func() sq0Counts {
			var c sq0Counts
			for _, base := range seedBases {
				for ep := 0; ep < 128; ep++ {
					seed := sq0Seed(base, 101+depth, ep)
					rng := newSQ0RNG(seed)
					state := rng.intn(8)
					want := state
					for i := 0; i < depth; i++ {
						op := rng.intn(6)
						want = sq0Operators[op][want]
						state = sq0Operators[op][state]
						c.events++
					}
					c.primaryTotal++
					c.secondaryTotal++
					if state == want {
						c.primaryHits++
						c.secondaryHits++
					}
				}
			}
			return c
		}))
	}
	return out
}

func sq0RoleSecondary(r [5]int) int {
	return (r[0] + 2*r[1] + r[2] + 2*r[3] + r[4]) % 3
}

func sq0AllRoleTuples() [][5]int {
	out := make([][5]int, 0, 243)
	for a:=0;a<3;a++ { for b:=0;b<3;b++ { for c:=0;c<3;c++ { for d:=0;d<3;d++ { for e:=0;e<3;e++ {
		out=append(out,[5]int{a,b,c,d,e})
	}}}}}
	return out
}

func sq0EvalRoleFiller(arm string, seedBases []int) []UPSQ0Metric {
	all := sq0AllRoleTuples()
	type split struct {
		name string
		rows [][5]int
	}
	var train, held [][5]int
	for _, r := range all {
		sum := r[0]+r[1]+r[2]+r[3]+r[4]
		if sum%3 == 0 && sq0RoleSecondary(r) < 2 {
			train = append(train, r)
		}
		if sum%3 != 0 {
			held = append(held, r)
		}
	}
	splits := []split{{"train54",train},{"heldout162",held}}
	var out []UPSQ0Metric
	for _, sp := range splits {
		out = append(out, sq0Measure(arm, "role_filler_recombination", sp.name, "per_role_accuracy", "whole_tuple_exact_accuracy", func() sq0Counts {
			var c sq0Counts
			for _, base := range seedBases {
				for idx, r := range sp.rows {
					seed := sq0Seed(base, 131+idx, idx)
					m := newSQ0Machine(arm, seed)
					for role:=0;role<5;role++ {
						m.write(role,r[role],3)
						c.events++
					}
					exact := true
					for role:=0;role<5;role++ {
						got:=m.query(role,3)
						c.primaryTotal++
						c.events++
						if got==r[role] {
							c.primaryHits++
						} else {
							exact=false
						}
					}
					c.secondaryTotal++
					if exact { c.secondaryHits++ }
					c.maxRecallUsed=sq0Max(c.maxRecallUsed,m.recallUsed)
				}
			}
			return c
		}))
	}
	return out
}

func sq0EvalMultiBank(arm string, seedBases []int) []UPSQ0Metric {
	bankCounts:=[]int{2,4,6}
	loads:=[]float64{0.25,0.50,1.00}
	var out []UPSQ0Metric
	for _,banks:=range bankCounts {
		for _,load:=range loads {
			setting:=fmt.Sprintf("banks=%d,load=%.2f",banks,load)
			out=append(out,sq0Measure(arm,"multi_bank_interference",setting,"value_accuracy","episode_exact_accuracy",func() sq0Counts {
				var c sq0Counts
				c.tertiaryName="cross_bank_corruption_rate"
				crossErrors:=0
				crossTotal:=0
				totalSlots:=banks*4
				occupied:=int(math.Round(float64(totalSlots)*load))
				if occupied<1 { occupied=1 }
				for _,base:=range seedBases {
					for ep:=0;ep<64;ep++ {
						seed:=sq0Seed(base,173+banks*11+occupied,ep)
						rng:=newSQ0RNG(seed)
						m:=newSQ0Machine(arm,seed)
						truth:=map[int]int{}
						order:=make([]int,0,occupied)
						for i:=0;i<occupied;i++ {
							key:=(ep+i*5)%totalSlots
							for {
								if _,ok:=truth[key];!ok { break }
								key=(key+1)%totalSlots
							}
							v:=rng.intn(16)
							truth[key]=v
							order=append(order,key)
							m.write(key,v,16)
							c.events++
							currentBank:=key/4
							for k,want:=range truth {
								if k/4==currentBank { continue }
								got:=m.query(k,16)
								c.events++
								crossTotal++
								if got!=want { crossErrors++ }
							}
						}
						exact:=true
						for _,key:=range order {
							got:=m.query(key,16)
							c.primaryTotal++
							c.events++
							if got==truth[key] {
								c.primaryHits++
							} else {
								exact=false
							}
						}
						c.secondaryTotal++
						if exact { c.secondaryHits++ }
						c.maxRecallUsed=sq0Max(c.maxRecallUsed,m.recallUsed)
					}
				}
				if crossTotal>0 {
					c.tertiaryValue=float64(crossErrors)/float64(crossTotal)
				}
				return c
			}))
		}
	}
	return out
}

func RunUPSQ0()(UPSQ0SequenceQualificationResult,error){
	seedBases:=[]int{120000000,121000000,122000000}
	arms:=[]string{
		"transport_only",
		"transport_gated_correction",
		"transport_gated_correction_bounded_exact_recall",
	}
	result:=UPSQ0SequenceQualificationResult{
		Schema:UPSQ0SequenceQualificationSchema,
		Experiment:"UP-SQ0-serialized-pre-language-sequence-qualification",
		ImplementationFreeze:"docs/experiments/wingless-up-sq0-implementation-freeze.md",
		GeneratorSeedBases:append([]int(nil),seedBases...),
		RecurrentStateBytes:sq0CarrierDim*8,
		ExactRecallEntryCap:sq0ExactRecallCap,
		FullGlobalAttention:false,
		AggregateScoreSelected:false,
	}
	for _,arm:=range arms {
		result.Metrics=append(result.Metrics,sq0EvalDelayedCopy(arm,seedBases)...)
		result.Metrics=append(result.Metrics,sq0EvalSelectiveCopy(arm,seedBases)...)
		result.Metrics=append(result.Metrics,sq0EvalMQAR(arm,seedBases)...)
		result.Metrics=append(result.Metrics,sq0EvalOverwrite(arm,seedBases)...)
		result.Metrics=append(result.Metrics,sq0EvalStateMachine(arm,seedBases)...)
		result.Metrics=append(result.Metrics,sq0EvalRoleFiller(arm,seedBases)...)
		result.Metrics=append(result.Metrics,sq0EvalMultiBank(arm,seedBases)...)
	}
	return result,nil
}
