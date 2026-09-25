package unitary

const UP91BSelectiveRecallRoutingSchema = "wingless.up91b-selective-recall-routing.v1"

type UP91BSelectiveRecallPoint struct {
	Arm                  string  `json:"arm"`
	Family               string  `json:"family"`
	Setting              string  `json:"setting"`
	QueryAccuracy        float64 `json:"query_accuracy"`
	ExactQuerySetAccuracy float64 `json:"exact_query_set_accuracy"`
	RecallEntriesUsed    int     `json:"recall_entries_used"`
	RecurrentStateBytes  int     `json:"recurrent_state_bytes"`
	ExactRecallBytes     int     `json:"exact_recall_bytes"`
}

type UP91BSelectiveRecallRoutingResult struct {
	Schema             string                     `json:"schema"`
	Experiment         string                     `json:"experiment"`
	SourceUPSQ0Seal    string                     `json:"source_up_sq0_seal"`
	SeedBases          []int                      `json:"seed_bases"`
	ExactRecallCap     int                        `json:"exact_recall_cap"`
	RecurrentStateBytes int                       `json:"recurrent_state_bytes"`
	Points             []UP91BSelectiveRecallPoint `json:"points"`
}

func up91bAdmit(m *sq0Machine, policy string, key, value int, salient bool) {
	switch policy {
	case "fifo_all_writes":
		m.recallWrite(key, value)
	case "salience_only":
		if salient {
			m.recallWrite(key, value)
		}
	}
}

func up91bQuery(m *sq0Machine, key, vocab int) int {
	if v, ok := m.recallQuery(key); ok {
		return v
	}
	return m.query(key, vocab)
}

func up91bSelectiveCopy(policy string, length int, seedBases []int) UP91BSelectiveRecallPoint {
	queryHits := 0
	queryTotal := 0
	exactHits := 0
	episodes := 0
	maxRecall := 0
	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 211+length, ep)
			rng := newSQ0RNG(seed)
			m := newSQ0Machine("transport_gated_correction", seed)
			truth := make([]int, length)
			marked := make([]int, 0, length/4)
			for i := 0; i < length; i++ {
				v := rng.intn(16)
				truth[i] = v
				salient := (i+ep+base)%4 == 0
				m.write(i, v, 16)
				up91bAdmit(m, policy, i, v, salient)
				if salient {
					marked = append(marked, i)
				}
			}
			exact := true
			for _, key := range marked {
				got := up91bQuery(m, key, 16)
				queryTotal++
				if got == truth[key] {
					queryHits++
				} else {
					exact = false
				}
			}
			episodes++
			if exact {
				exactHits++
			}
			if m.recallUsed > maxRecall {
				maxRecall = m.recallUsed
			}
		}
	}
	return UP91BSelectiveRecallPoint{
		Arm: policy,
		Family: "selective_copy",
		Setting: "length=" + itoa(length) + ",marked=25%",
		QueryAccuracy: float64(queryHits) / float64(queryTotal),
		ExactQuerySetAccuracy: float64(exactHits) / float64(episodes),
		RecallEntriesUsed: maxRecall,
		RecurrentStateBytes: sq0CarrierDim * 8,
		ExactRecallBytes: maxRecall * 16,
	}
}

func up91bSparseDistractor(policy string, totalWrites, salientCount int, seedBases []int) UP91BSelectiveRecallPoint {
	queryHits := 0
	queryTotal := 0
	exactHits := 0
	episodes := 0
	maxRecall := 0
	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 241+totalWrites*3+salientCount, ep)
			rng := newSQ0RNG(seed)
			m := newSQ0Machine("transport_gated_correction", seed)
			truth := make(map[int]int, salientCount)
			for s := 0; s < salientCount; s++ {
				key := s
				v := rng.intn(32)
				truth[key] = v
				m.write(key, v, 32)
				up91bAdmit(m, policy, key, v, true)
			}
			distractors := totalWrites - 2*salientCount
			pre := distractors / 2
			post := distractors - pre
			for d := 0; d < pre; d++ {
				key := 1000 + ep*1000 + d
				v := rng.intn(32)
				m.write(key, v, 32)
				up91bAdmit(m, policy, key, v, false)
			}
			for s := 0; s < salientCount; s++ {
				key := s
				old := truth[key]
				v := (old + 1 + rng.intn(31)) % 32
				truth[key] = v
				m.write(key, v, 32)
				up91bAdmit(m, policy, key, v, true)
			}
			for d := 0; d < post; d++ {
				key := 200000 + ep*1000 + d
				v := rng.intn(32)
				m.write(key, v, 32)
				up91bAdmit(m, policy, key, v, false)
			}
			exact := true
			for s := 0; s < salientCount; s++ {
				got := up91bQuery(m, s, 32)
				queryTotal++
				if got == truth[s] {
					queryHits++
				} else {
					exact = false
				}
			}
			episodes++
			if exact {
				exactHits++
			}
			if m.recallUsed > maxRecall {
				maxRecall = m.recallUsed
			}
		}
	}
	return UP91BSelectiveRecallPoint{
		Arm: policy,
		Family: "sparse_query_distractor_stream",
		Setting: "writes=" + itoa(totalWrites) + ",salient=" + itoa(salientCount),
		QueryAccuracy: float64(queryHits) / float64(queryTotal),
		ExactQuerySetAccuracy: float64(exactHits) / float64(episodes),
		RecallEntriesUsed: maxRecall,
		RecurrentStateBytes: sq0CarrierDim * 8,
		ExactRecallBytes: maxRecall * 16,
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var buf [24]byte
	i := len(buf)
	n := v
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

func RunUP91B() (UP91BSelectiveRecallRoutingResult, error) {
	seedBases := []int{123000000, 124000000}
	result := UP91BSelectiveRecallRoutingResult{
		Schema: UP91BSelectiveRecallRoutingSchema,
		Experiment: "UP-91B-selective-recall-routing",
		SourceUPSQ0Seal: "8de9e11f26cc44c75f51d26cb606f97e35816479",
		SeedBases: append([]int(nil), seedBases...),
		ExactRecallCap: sq0ExactRecallCap,
		RecurrentStateBytes: sq0CarrierDim * 8,
	}
	for _, policy := range []string{"fifo_all_writes", "salience_only"} {
		for _, length := range []int{32, 64, 128, 256} {
			result.Points = append(result.Points, up91bSelectiveCopy(policy, length, seedBases))
		}
		for _, writes := range []int{32, 64, 128, 256} {
			for _, salient := range []int{4, 8, 16} {
				result.Points = append(result.Points, up91bSparseDistractor(policy, writes, salient, seedBases))
			}
		}
	}
	return result, nil
}
