// Package r53prefreeze implements the frozen Wingless R53 readout slope-block
// mechanism only. It does not acquire input, train a model, evaluate heldout
// outcomes, mutate accepted refs, or dispatch experiments. The independent CKB
// science authority must qualify external source birth, freeze the full
// preregistration and admit any subsequent integrated experiment.
package r53prefreeze

import (
	"errors"
	"math"
)

type ArmID string

const (
	A0 ArmID = "A0" // untied six-slope baseline, one shared intercept
	A1 ArmID = "A1" // tie the early first two slopes only
	A2 ArmID = "A2" // tie the late last four slopes only
	A3 ArmID = "A3" // original R52 two-group negative/mixed parent comparator

	StateDimension      = 6
	MaxReadoutParameters = 7
	AdaptationPerArm     = 436
	TotalAdaptationUnits = 1744
	ExternalModelCalls   = 0
	TokenizerCalls       = 0
)

var ErrUnfrozenArm = errors.New("R53_UNDECLARED_ARM")
var ErrNonfiniteState = errors.New("R53_NONFINITE_STATE")
var ErrNonfiniteCoefficient = errors.New("R53_NONFINITE_COEFFICIENT")
var ErrHiddenParameter = errors.New("R53_UNUSED_COEFFICIENT_NONZERO")

type definition struct {
	id ArmID
	groups [StateDimension]int
	groupCount int
}

var defs = [...]definition{
	{id: A0, groups: [6]int{0, 1, 2, 3, 4, 5}, groupCount: 6},
	{id: A1, groups: [6]int{0, 0, 1, 2, 3, 4}, groupCount: 5},
	{id: A2, groups: [6]int{0, 1, 2, 2, 2, 2}, groupCount: 3},
	{id: A3, groups: [6]int{0, 0, 1, 1, 1, 1}, groupCount: 2},
}

func lookup(arm ArmID) (definition, error) {
	for _, def := range defs {
		if arm == def.id {
			return def, nil
		}
	}
	return definition{}, ErrUnfrozenArm
}

// Arms returns the fixed preregistered A0/A1/A2/A3 order, not a mutable registry.
func Arms() [4]ArmID { return [4]ArmID{A0, A1, A2, A3} }

// SelectedAlpha is the immutable six-dimensional state-former coefficient
// setting already frozen for R53 by R189, not an additional intervention.
func SelectedAlpha() [StateDimension]float64 {
	return [StateDimension]float64{0.75, 0.75, 0.125, 0.125, 1, 1}
}

// Groups and parameter count are copies. The sole intercept is group 0's
// design-row bias position; there is never an additional group intercept.
func Groups(arm ArmID) ([StateDimension]int, int, error) {
	def, err := lookup(arm)
	if err != nil {
		return [StateDimension]int{}, 0, err
	}
	return def.groups, def.groupCount + 1, nil
}

// DesignRow reproduces the immutable R50 fitting transform:
// z[0]=1; for i in 0..5: z[groups[i]+1]+=state[i].
// Only the first groupCount+1 entries can be populated.
func DesignRow(arm ArmID, state [StateDimension]float64) ([MaxReadoutParameters]float64, error) {
	var row [MaxReadoutParameters]float64
	def, err := lookup(arm)
	if err != nil {
		return row, err
	}
	row[0] = 1
	for i, v := range state {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return [MaxReadoutParameters]float64{}, ErrNonfiniteState
		}
		row[def.groups[i]+1] += v
	}
	for _, v := range row {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return [MaxReadoutParameters]float64{}, ErrNonfiniteState
		}
	}
	return row, nil
}

// ExpandReadout reproduces R50's original readout parameter mapping:
// beta[0]=theta[0]; beta[i+1]=theta[groups[i]+1].
// theta has fixed maximum length 7; inactive slots MUST be exactly zero.
func ExpandReadout(arm ArmID, theta [MaxReadoutParameters]float64) ([MaxReadoutParameters]float64, error) {
	var expanded [MaxReadoutParameters]float64
	def, err := lookup(arm)
	if err != nil {
		return expanded, err
	}
	for i, v := range theta {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return expanded, ErrNonfiniteCoefficient
		}
		if i > def.groupCount && v != 0 {
			return expanded, ErrHiddenParameter
		}
	}
	expanded[0] = theta[0]
	for i, g := range def.groups {
		expanded[i+1] = theta[g+1]
	}
	return expanded, nil
}

// Predict computes a readout with the same single intercept and tied slopes as
// the original R50 design. No source input or evaluation outcome is accessed.
func Predict(arm ArmID, state [StateDimension]float64, theta [MaxReadoutParameters]float64) (float64, error) {
	row, err := DesignRow(arm, state)
	if err != nil {
		return 0, err
	}
	beta, err := ExpandReadout(arm, theta)
	if err != nil {
		return 0, err
	}
	sum := beta[0]
	for i, x := range state {
		sum += beta[i+1] * x
	}
	if math.IsNaN(sum) || math.IsInf(sum, 0) {
		return 0, ErrNonfiniteCoefficient
	}
	// Keep the algebraic form under test: group-specific design rows produce
	// exactly the same result within floating-point rounding tolerance.
	projected := 0.0
	for i, x := range row {
		projected += x * theta[i]
	}
	if math.IsNaN(projected) || math.IsInf(projected, 0) {
		return 0, ErrNonfiniteCoefficient
	}
	return sum, nil
}
