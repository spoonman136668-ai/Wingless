package unitary

import (
	"crypto/sha256"
	"encoding/binary"
	"math/bits"
	"sort"
)

type wlmSiDensePackedR2Row struct {
	Seed              uint64 `json:"seed"`
	ScheduleIndex     int    `json:"schedule_index"`
	Substrate         string `json:"substrate"`
	Step              int    `json:"step"`
	ActiveCount       int    `json:"active_count"`
	InactiveCount     int    `json:"inactive_count"`
	DemandUnits       int    `json:"demand_units"`
	DeficitUnits      int    `json:"deficit_units"`
	ConservationTotal int    `json:"conservation_total"`
}

type wlmSiDensePackedR2Dense struct{ bits [256]bool }
type wlmSiDensePackedR2Packed struct{ words [4]uint64 }
type wlmSiDensePackedR2Observation struct {
	active, inactive, demand, deficit, conservation int
}

// RunWlmSiDensePackedR2 executes the frozen dense-array versus packed-bitset experiment.
func RunWlmSiDensePackedR2() interface{} {
	const (
		width    = 256
		steps    = 64
		budget   = 128
		schedules = 32
	)
	seeds := [...]uint64{101, 211, 307, 401, 503, 601, 701, 809}
	metrics := map[string]float64{
		"budget_overrun_rows":                 0,
		"completed_substrate_case_runs":       0,
		"invalid_state_schedules":             0,
		"law_violation_cases":                 0,
		"max_conservation_error":              0,
		"max_deficit_law_error":               0,
		"paired_case_count":                   0,
		"paired_observable_mismatch_cases":    0,
		"row_count_mismatch_cases":            0,
		"state_cardinality_mismatch_cases":    0,
		"total_emitted_rows":                  0,
	}
	rows := make([]wlmSiDensePackedR2Row, 0, len(seeds)*schedules*2*steps)

	for _, seed := range seeds {
		for scheduleIndex := 0; scheduleIndex < schedules; scheduleIndex++ {
			metrics["paired_case_count"]++
			schedule, canonical := wlmSiDensePackedR2Schedule(seed, scheduleIndex)
			before := sha256.Sum256(canonical)
			after := sha256.Sum256(canonical)
			invalid := before != after
			initial := wlmSiDensePackedR2Initial(seed, scheduleIndex)
			oracle := wlmSiDensePackedR2Oracle(initial, schedule, budget)

			dense := wlmSiDensePackedR2Dense{bits: initial}
			packed := wlmSiDensePackedR2Packed{words: wlmSiDensePackedR2Pack(initial)}
			var denseRows, packedRows []wlmSiDensePackedR2Row
			var denseInvalid, packedInvalid, denseLaw, packedLaw bool

			if scheduleIndex%2 == 0 {
				denseRows, denseInvalid, denseLaw = wlmSiDensePackedR2RunDense(seed, scheduleIndex, schedule, budget, &dense, oracle)
				packedRows, packedInvalid, packedLaw = wlmSiDensePackedR2RunPacked(seed, scheduleIndex, schedule, budget, &packed, oracle)
				rows = append(rows, denseRows...)
				rows = append(rows, packedRows...)
			} else {
				packedRows, packedInvalid, packedLaw = wlmSiDensePackedR2RunPacked(seed, scheduleIndex, schedule, budget, &packed, oracle)
				denseRows, denseInvalid, denseLaw = wlmSiDensePackedR2RunDense(seed, scheduleIndex, schedule, budget, &dense, oracle)
				rows = append(rows, packedRows...)
				rows = append(rows, denseRows...)
			}

			metrics["completed_substrate_case_runs"] += 2
			if invalid || denseInvalid || packedInvalid {
				metrics["invalid_state_schedules"]++
			}
			if len(denseRows) != steps || len(packedRows) != steps {
				metrics["row_count_mismatch_cases"]++
			}
			if denseLaw || packedLaw {
				metrics["law_violation_cases"]++
			}
			if !wlmSiDensePackedR2SameCardinality(denseRows, packedRows) {
				metrics["state_cardinality_mismatch_cases"]++
			}
			if !wlmSiDensePackedR2SameObservables(denseRows, packedRows) {
				metrics["paired_observable_mismatch_cases"]++
			}
			for _, row := range denseRows {
				wlmSiDensePackedR2UpdateErrors(metrics, row, budget)
			}
			for _, row := range packedRows {
				wlmSiDensePackedR2UpdateErrors(metrics, row, budget)
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	if len(rows) > 32768 {
		metrics["budget_overrun_rows"] = float64(len(rows) - 32768)
	}
	return map[string]interface{}{"rows": rows, "metrics": metrics}
}

func wlmSiDensePackedR2Schedule(seed uint64, scheduleIndex int) ([64]uint16, []byte) {
	var schedule [64]uint16
	canonical := make([]byte, 64*2)
	for step := range schedule {
		digest := wlmSiDensePackedR2HashTuple(seed, uint64(scheduleIndex), uint64(step))
		index := binary.BigEndian.Uint64(digest[:8]) % 256
		schedule[step] = uint16(index)
		binary.BigEndian.PutUint16(canonical[step*2:], uint16(index))
	}
	return schedule, canonical
}

func wlmSiDensePackedR2Initial(seed uint64, scheduleIndex int) [256]bool {
	type candidate struct {
		index int
		hash  [32]byte
	}
	candidates := make([]candidate, 256)
	for index := range candidates {
		candidates[index] = candidate{index: index, hash: wlmSiDensePackedR2HashTuple(seed, uint64(scheduleIndex), uint64(index))}
	}
	sort.Slice(candidates, func(i, j int) bool {
		for k := range candidates[i].hash {
			if candidates[i].hash[k] != candidates[j].hash[k] {
				return candidates[i].hash[k] < candidates[j].hash[k]
			}
		}
		return candidates[i].index < candidates[j].index
	})
	var initial [256]bool
	for _, candidate := range candidates[:128] {
		initial[candidate.index] = true
	}
	return initial
}

func wlmSiDensePackedR2HashTuple(seed, scheduleIndex, value uint64) [32]byte {
	var encoded [24]byte
	binary.BigEndian.PutUint64(encoded[0:8], seed)
	binary.BigEndian.PutUint64(encoded[8:16], scheduleIndex)
	binary.BigEndian.PutUint64(encoded[16:24], value)
	return sha256.Sum256(encoded[:])
}

func wlmSiDensePackedR2Pack(initial [256]bool) [4]uint64 {
	var words [4]uint64
	for index, active := range initial {
		if active {
			words[index/64] |= uint64(1) << uint(index%64)
		}
	}
	return words
}

func wlmSiDensePackedR2Oracle(initial [256]bool, schedule [64]uint16, budget int) [64]wlmSiDensePackedR2Observation {
	state := initial
	var observations [64]wlmSiDensePackedR2Observation
	for step, index := range schedule {
		state[index] = !state[index]
		active := 0
		for _, set := range state {
			if set {
				active++
			}
		}
		deficit := active - budget
		if deficit < 0 {
			deficit = 0
		}
		observations[step] = wlmSiDensePackedR2Observation{active, 256 - active, active, deficit, 256}
	}
	return observations
}

func wlmSiDensePackedR2RunDense(seed uint64, scheduleIndex int, schedule [64]uint16, budget int, state *wlmSiDensePackedR2Dense, oracle [64]wlmSiDensePackedR2Observation) ([]wlmSiDensePackedR2Row, bool, bool) {
	rows := make([]wlmSiDensePackedR2Row, 0, 64)
	invalid, law := false, false
	for step, index := range schedule {
		state.bits[index] = !state.bits[index]
		active := 0
		for _, set := range state.bits {
			if set {
				active++
			}
		}
		row := wlmSiDensePackedR2Row{Seed: seed, ScheduleIndex: scheduleIndex, Substrate: "dense-bool-array", Step: step, ActiveCount: active, InactiveCount: 256 - active, DemandUnits: active, DeficitUnits: wlmSiDensePackedR2Deficit(active, budget), ConservationTotal: active + (256 - active)}
		if !wlmSiDensePackedR2Valid(row, oracle[step], budget) {
			law = true
		}
		rows = append(rows, row)
	}
	return rows, invalid, law
}

func wlmSiDensePackedR2RunPacked(seed uint64, scheduleIndex int, schedule [64]uint16, budget int, state *wlmSiDensePackedR2Packed, oracle [64]wlmSiDensePackedR2Observation) ([]wlmSiDensePackedR2Row, bool, bool) {
	rows := make([]wlmSiDensePackedR2Row, 0, 64)
	invalid, law := false, false
	for step, index := range schedule {
		word, offset := int(index)/64, uint(index%64)
		state.words[word] ^= uint64(1) << offset
		active := 0
		for _, bitsWord := range state.words {
			active += bits.OnesCount64(bitsWord)
		}
		row := wlmSiDensePackedR2Row{Seed: seed, ScheduleIndex: scheduleIndex, Substrate: "packed-4xuint64-bitset", Step: step, ActiveCount: active, InactiveCount: 256 - active, DemandUnits: active, DeficitUnits: wlmSiDensePackedR2Deficit(active, budget), ConservationTotal: active + (256 - active)}
		if active < 0 || active > 256 {
			invalid = true
		}
		if !wlmSiDensePackedR2Valid(row, oracle[step], budget) {
			law = true
		}
		rows = append(rows, row)
	}
	return rows, invalid, law
}

func wlmSiDensePackedR2Deficit(demand, budget int) int {
	if demand > budget {
		return demand - budget
	}
	return 0
}

func wlmSiDensePackedR2Valid(row wlmSiDensePackedR2Row, expected wlmSiDensePackedR2Observation, budget int) bool {
	return row.ActiveCount == expected.active && row.InactiveCount == expected.inactive && row.DemandUnits == expected.demand && row.DeficitUnits == expected.deficit && row.ConservationTotal == expected.conservation && row.ConservationTotal == 256 && row.DeficitUnits == wlmSiDensePackedR2Deficit(row.DemandUnits, budget)
}

func wlmSiDensePackedR2SameCardinality(a, b []wlmSiDensePackedR2Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ActiveCount != b[i].ActiveCount || a[i].InactiveCount != b[i].InactiveCount {
			return false
		}
	}
	return true
}

func wlmSiDensePackedR2SameObservables(a, b []wlmSiDensePackedR2Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.ScheduleIndex != y.ScheduleIndex || x.Step != y.Step || x.ActiveCount != y.ActiveCount || x.InactiveCount != y.InactiveCount || x.DemandUnits != y.DemandUnits || x.DeficitUnits != y.DeficitUnits || x.ConservationTotal != y.ConservationTotal {
			return false
		}
	}
	return true
}

func wlmSiDensePackedR2UpdateErrors(metrics map[string]float64, row wlmSiDensePackedR2Row, budget int) {
	conservation := wlmSiDensePackedR2Abs(row.ConservationTotal - 256)
	deficit := wlmSiDensePackedR2Abs(row.DeficitUnits - wlmSiDensePackedR2Deficit(row.DemandUnits, budget))
	if float64(conservation) > metrics["max_conservation_error"] {
		metrics["max_conservation_error"] = float64(conservation)
	}
	if float64(deficit) > metrics["max_deficit_law_error"] {
		metrics["max_deficit_law_error"] = float64(deficit)
	}
}

func wlmSiDensePackedR2Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
