package unitary

import (
	"crypto/sha256"
	"encoding/hex"
)

type upLm11dDeficitLawR1Row struct {
	Seed                 int  `json:"seed"`
	WriteBudget          int  `json:"write_budget"`
	DemandOrdinal        int  `json:"demand_ordinal"`
	State                int  `json:"state"`
	Attempted            bool `json:"attempted"`
	Accepted             bool `json:"accepted"`
	AcceptedWritesSoFar  int  `json:"accepted_writes_so_far"`
	DeficitWritesSoFar   int  `json:"deficit_writes_so_far"`
}

type upLm11dDeficitLawR1Metrics struct {
	EmittedRows               float64 `json:"emitted_rows"`
	BudgetCases               float64 `json:"budget_cases"`
	InvalidPermutations       float64 `json:"invalid_permutations"`
	UniqueStatesPerSeed       float64 `json:"unique_states_per_seed"`
	AggregateDemandWrites     float64 `json:"aggregate_demand_writes"`
	AggregateAcceptedWrites   float64 `json:"aggregate_accepted_writes"`
	AggregateDeficit          float64 `json:"aggregate_deficit"`
	LawViolationCases         float64 `json:"law_violation_cases"`
	BudgetOverrunRows         float64 `json:"budget_overrun_rows"`
	MaxConservationError      float64 `json:"max_conservation_error"`
	PrefixDigestComparisons   float64 `json:"prefix_digest_comparisons"`
	PrefixDigestMismatchCases float64 `json:"prefix_digest_mismatch_cases"`
}

// RunUpLm11dDeficitLawR1 evaluates bounded admissions against same-seed full-run prefixes.
func RunUpLm11dDeficitLawR1() interface{} {
	seeds := [...]int{104729, 130363, 155921, 181081, 206369, 231701, 257053, 282407}
	budgets := [...]int{0, 64, 128, 256, 512}
	rows := make([]upLm11dDeficitLawR1Row, 0, len(seeds)*len(budgets)*512)

	invalidPermutations := 0
	minimumUniqueStates := 512
	aggregateDemand := 0
	aggregateAccepted := 0
	aggregateDeficit := 0
	lawViolationCases := 0
	budgetOverrunRows := 0
	maxConservationError := 0
	prefixDigestComparisons := 0
	prefixDigestMismatchCases := 0
	budgetCases := 0

	for _, seed := range seeds {
		order := upLm11dDeficitLawR1Permutation(seed)
		if !upLm11dDeficitLawR1ValidPermutation(order) {
			invalidPermutations++
			continue
		}
		unique := make(map[int]struct{}, 512)
		for _, state := range order {
			unique[state] = struct{}{}
		}
		if len(unique) < minimumUniqueStates {
			minimumUniqueStates = len(unique)
		}

		for _, budget := range budgets {
			budgetCases++
			var boundedState [512]bool
			acceptedWrites := 0
			deficitWrites := 0
			overrunInCase := 0

			for ordinal, state := range order {
				attempted := true
				accepted := false
				if attempted && acceptedWrites < budget {
					boundedState[state] = true
					acceptedWrites++
					accepted = true
				} else if attempted {
					deficitWrites++
				}
				if accepted && acceptedWrites > budget {
					overrunInCase++
					budgetOverrunRows++
				}
				rows = append(rows, upLm11dDeficitLawR1Row{
					Seed: seed, WriteBudget: budget, DemandOrdinal: ordinal, State: state,
					Attempted: attempted, Accepted: accepted,
					AcceptedWritesSoFar: acceptedWrites, DeficitWritesSoFar: deficitWrites,
				})
				aggregateDemand++
			}

			expectedAccepted := budget
			if expectedAccepted > len(order) {
				expectedAccepted = len(order)
			}
			expectedDeficit := len(order) - expectedAccepted
			conservationError := upLm11dDeficitLawR1Abs(len(order) - (acceptedWrites + deficitWrites))
			if conservationError > maxConservationError {
				maxConservationError = conservationError
			}
			aggregateAccepted += acceptedWrites
			aggregateDeficit += deficitWrites

			expectedDigest := upLm11dDeficitLawR1PrefixDigest(order, expectedAccepted)
			boundedDigest := upLm11dDeficitLawR1CanonicalDigest(boundedState)
			prefixDigestComparisons++
			digestMismatch := boundedDigest != expectedDigest
			if digestMismatch {
				prefixDigestMismatchCases++
			}
			if acceptedWrites != expectedAccepted || deficitWrites != expectedDeficit || conservationError != 0 || overrunInCase != 0 || digestMismatch {
				lawViolationCases++
			}
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": upLm11dDeficitLawR1Metrics{
			EmittedRows:               float64(len(rows)),
			BudgetCases:               float64(budgetCases),
			InvalidPermutations:       float64(invalidPermutations),
			UniqueStatesPerSeed:       float64(minimumUniqueStates),
			AggregateDemandWrites:     float64(aggregateDemand),
			AggregateAcceptedWrites:   float64(aggregateAccepted),
			AggregateDeficit:          float64(aggregateDeficit),
			LawViolationCases:         float64(lawViolationCases),
			BudgetOverrunRows:         float64(budgetOverrunRows),
			MaxConservationError:      float64(maxConservationError),
			PrefixDigestComparisons:   float64(prefixDigestComparisons),
			PrefixDigestMismatchCases: float64(prefixDigestMismatchCases),
		},
	}
}

func upLm11dDeficitLawR1Permutation(seed int) [512]int {
	var order [512]int
	for i := range order {
		order[i] = i
	}
	stream := uint64(seed) ^ 0xD1B54A32D192ED03
	for i := len(order) - 1; i > 0; i-- {
		stream = upLm11dDeficitLawR1SplitMix64(stream)
		j := int(stream % uint64(i+1))
		order[i], order[j] = order[j], order[i]
	}
	return order
}

func upLm11dDeficitLawR1ValidPermutation(order [512]int) bool {
	var seen [512]bool
	for _, state := range order {
		if state < 0 || state >= len(seen) || seen[state] {
			return false
		}
		seen[state] = true
	}
	return true
}

func upLm11dDeficitLawR1PrefixDigest(order [512]int, accepted int) string {
	var state [512]bool
	for i := 0; i < accepted; i++ {
		state[order[i]] = true
	}
	return upLm11dDeficitLawR1CanonicalDigest(state)
}

func upLm11dDeficitLawR1CanonicalDigest(state [512]bool) string {
	var canonical [64]byte
	for index, present := range state {
		if present {
			canonical[index/8] |= 1 << uint(index%8)
		}
	}
	digest := sha256.Sum256(canonical[:])
	return hex.EncodeToString(digest[:])
}

func upLm11dDeficitLawR1SplitMix64(state uint64) uint64 {
	state += 0x9E3779B97F4A7C15
	z := state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

func upLm11dDeficitLawR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
