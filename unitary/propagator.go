package unitary

import (
	"fmt"
	"math"
	"math/cmplx"
)

// State is a bounded complex-valued latent state. UP-0 keeps this deliberately
// small and local; it is not a quantum state running on quantum hardware.
type State []complex128

// Coupling applies a two-coordinate unitary rotation. Theta controls mixing and
// Phi controls relative phase. The operation is a complex Givens rotation.
type Coupling struct {
	A     int     `json:"a"`
	B     int     `json:"b"`
	Theta float64 `json:"theta"`
	Phi   float64 `json:"phi"`
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func validateState(state State) error {
	if len(state) < 2 {
		return fmt.Errorf("unitary state must contain at least two amplitudes")
	}
	for i, v := range state {
		if !finite(real(v)) || !finite(imag(v)) {
			return fmt.Errorf("unitary state amplitude %d is not finite", i)
		}
	}
	return nil
}

// NormSquared returns the total probability mass represented by the state.
func NormSquared(state State) (float64, error) {
	if err := validateState(state); err != nil {
		return 0, err
	}
	var total float64
	for _, v := range state {
		a := cmplx.Abs(v)
		total += a * a
	}
	if !finite(total) {
		return 0, fmt.Errorf("unitary state norm is not finite")
	}
	return total, nil
}

// Normalize returns a normalized copy of state.
func Normalize(state State) (State, error) {
	norm2, err := NormSquared(state)
	if err != nil {
		return nil, err
	}
	if norm2 <= 0 {
		return nil, fmt.Errorf("unitary state norm must be positive")
	}
	scale := complex(1/math.Sqrt(norm2), 0)
	out := append(State(nil), state...)
	for i := range out {
		out[i] *= scale
	}
	return out, nil
}

func validateCoupling(c Coupling, dimension int) error {
	if c.A < 0 || c.B < 0 || c.A >= dimension || c.B >= dimension {
		return fmt.Errorf("unitary coupling index out of range")
	}
	if c.A == c.B {
		return fmt.Errorf("unitary coupling requires distinct coordinates")
	}
	if !finite(c.Theta) || !finite(c.Phi) {
		return fmt.Errorf("unitary coupling parameters must be finite")
	}
	return nil
}

// Propagate applies a product of local unitary rotations to a copy of initial.
// Every step has the matrix:
//
//	[ cos(theta)              -exp(-i*phi) sin(theta) ]
//	[ exp(+i*phi) sin(theta)   cos(theta)             ]
//
// so propagation preserves the L2 norm up to floating-point error.
func Propagate(initial State, program []Coupling) (State, error) {
	if err := validateState(initial); err != nil {
		return nil, err
	}
	out := append(State(nil), initial...)
	for _, step := range program {
		if err := validateCoupling(step, len(out)); err != nil {
			return nil, err
		}

		a := out[step.A]
		b := out[step.B]
		c := complex(math.Cos(step.Theta), 0)
		s := complex(math.Sin(step.Theta), 0)
		phase := cmplx.Rect(1, step.Phi)

		out[step.A] = c*a - cmplx.Conj(phase)*s*b
		out[step.B] = phase*s*a + c*b
	}
	return out, nil
}

// Invert returns the exact inverse program in reverse application order.
func Invert(program []Coupling) []Coupling {
	out := make([]Coupling, len(program))
	for i := range program {
		step := program[len(program)-1-i]
		step.Theta = -step.Theta
		out[i] = step
	}
	return out
}

// Probabilities converts amplitudes into normalized measurement probabilities.
func Probabilities(state State) ([]float64, error) {
	norm2, err := NormSquared(state)
	if err != nil {
		return nil, err
	}
	if norm2 <= 0 {
		return nil, fmt.Errorf("unitary state norm must be positive")
	}
	out := make([]float64, len(state))
	for i, v := range state {
		a := cmplx.Abs(v)
		out[i] = (a * a) / norm2
	}
	return out, nil
}

// ArgMax returns the first coordinate with maximum measurement probability.
func ArgMax(state State) (int, error) {
	probabilities, err := Probabilities(state)
	if err != nil {
		return 0, err
	}
	best := 0
	for i := 1; i < len(probabilities); i++ {
		if probabilities[i] > probabilities[best] {
			best = i
		}
	}
	return best, nil
}

// L2Distance returns the Euclidean distance between two complex states.
func L2Distance(a, b State) (float64, error) {
	if err := validateState(a); err != nil {
		return 0, err
	}
	if err := validateState(b); err != nil {
		return 0, err
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("unitary states must have equal dimensions")
	}
	var total float64
	for i := range a {
		d := cmplx.Abs(a[i] - b[i])
		total += d * d
	}
	return math.Sqrt(total), nil
}
