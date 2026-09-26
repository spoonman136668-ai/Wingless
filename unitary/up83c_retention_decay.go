package unitary

const UP83CRetentionDecaySchema = "wingless.up83c-retention-decay.v1"

type UP83CPoint struct {
	Policy                  string  `json:"policy"`
	ActiveKeys              int     `json:"active_keys"`
	ChurnWrites             int     `json:"churn_writes"`
	HotQueryAccuracy        float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy     float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy       float64 `json:"cold_query_accuracy"`
	ColdSetExactAccuracy    float64 `json:"cold_set_exact_accuracy"`
	RecallEntriesUsed       int     `json:"recall_entries_used"`
	PolicyMetadataBytes     int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes int     `json:"total_bounded_memory_bytes"`
}

type UP83CRetentionDecayResult struct {
	Schema            string       `json:"schema"`
	Experiment        string       `json:"experiment"`
	SourceUP82CSeal   string       `json:"source_up82c_seal"`
	ExactRecallCap    int          `json:"exact_recall_cap"`
	EntryPayloadBytes int          `json:"entry_payload_bytes"`
	FutureOracleUsed  bool         `json:"future_oracle_used"`
	Points            []UP83CPoint `json:"points"`
}

func up83cRun(policy string, active, churn int, seedBases []int) UP83CPoint {
	hotHits, hotTotal, hotExactHits := 0, 0, 0
	coldHits, coldTotal, coldExactHits := 0, 0, 0
	episodes, maxEntries, maxBytes := 0, 0, 0
	meta := 16

	for _, base := range seedBases {
		for ep := 0; ep < 64; ep++ {
			seed := sq0Seed(base, 1901+active*17+churn, ep)
			rng := newSQ0RNG(seed)

			var write func(int, int)
			var query func(int) (int, bool)
			var used func() int
			var bytesUsed func() int

			switch policy {
			case "fifo", "lru":
				m := &up79cMemory{policy: policy}
				write = m.write
				query = m.query
				used = func() int { return m.count }
				bytesUsed = m.totalBytes
				if policy == "lru" {
					meta = 136
				} else {
					meta = 16
				}
			case "two_bit_aging":
				m := &up81cAging{}
				write = m.write
				query = m.query
				used = func() int { return m.count }
				bytesUsed = m.totalBytes
				meta = 5
			}

			truth := make([]int, active)
			for k := 0; k < active; k++ {
				v := rng.intn(32)
				truth[k] = v
				write(k, v)
			}

			for k := 0; k < active; k += 2 {
				query(k)
			}

			for j := 0; j < churn; j++ {
				write(active+j, rng.intn(32))
			}

			hotExact, coldExact := true, true
			for k := 0; k < active; k++ {
				got, ok := query(k)
				if k%2 == 0 {
					hotTotal++
					if ok && got == truth[k] {
						hotHits++
					} else {
						hotExact = false
					}
				} else {
					coldTotal++
					if ok && got == truth[k] {
						coldHits++
					} else {
						coldExact = false
					}
				}
			}
			if hotExact {
				hotExactHits++
			}
			if coldExact {
				coldExactHits++
			}
			episodes++
			if used() > maxEntries {
				maxEntries = used()
			}
			if bytesUsed() > maxBytes {
				maxBytes = bytesUsed()
			}
		}
	}

	return UP83CPoint{
		Policy: policy,
		ActiveKeys: active,
		ChurnWrites: churn,
		HotQueryAccuracy: float64(hotHits) / float64(hotTotal),
		HotSetExactAccuracy: float64(hotExactHits) / float64(episodes),
		ColdQueryAccuracy: float64(coldHits) / float64(coldTotal),
		ColdSetExactAccuracy: float64(coldExactHits) / float64(episodes),
		RecallEntriesUsed: maxEntries,
		PolicyMetadataBytes: meta,
		TotalBoundedMemoryBytes: maxBytes,
	}
}

func RunUP83C() (UP83CRetentionDecayResult, error) {
	result := UP83CRetentionDecayResult{
		Schema: UP83CRetentionDecaySchema,
		Experiment: "UP-83C-retention-decay",
		SourceUP82CSeal: "362fc1ef42363c3cc9da74c0638a4973bdc778e4",
		ExactRecallCap: 16,
		EntryPayloadBytes: 16,
		FutureOracleUsed: false,
	}
	seedBases := []int{157000000, 158000000}
	for _, policy := range []string{"fifo", "lru", "two_bit_aging"} {
		for _, active := range []int{24, 32} {
			for _, churn := range []int{0, 4, 8, 12, 16, 20, 24, 32} {
				result.Points = append(result.Points, up83cRun(policy, active, churn, seedBases))
			}
		}
	}
	return result, nil
}
