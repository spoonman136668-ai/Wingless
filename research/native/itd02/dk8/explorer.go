// Package dk8 is an opt-in native action-proposal experiment.
// It reuses the unmodified DK7 public selected-feedback memory. No DoorKey
// answer map or external fallback is embedded in the chooser.
package dk8

import (
    "crypto/sha256"
    "errors"

    dk7 "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7"
)

const Schema = "wingless.itd02.dk8.native-uncertainty-sourceonly.v1"

type Explorer struct {
    Native *dk7.Policy
}

func New() *Explorer {
    return &Explorer{Native: dk7.New()}
}

// Choose is a native decision even with zero previous positive reward.
// It never uses evaluator-owned action suggestions: it proposes the
// least-sampled action in the current public context, with a stable
// context-dependent tie-break. Only actual selected feedback updates memory.
func (e *Explorer) Choose(before dk7.Public) (*int, string, error) {
    if e == nil || e.Native == nil {
        return nil, "", errors.New("DK8_MISSING_NATIVE_MEMORY")
    }
    selected, status, err := e.Native.Choose(before)
    if err != nil {
        return nil, "", err
    }
    if selected != nil {
        return selected, "SELECT_NATIVE_OBSERVED_POSITIVE_REWARD", nil
    }
    if status != "ABSTAIN_NO_OBSERVED_POSITIVE_ACTION" {
        return nil, "", errors.New("DK8_UNRECOGNIZED_BASE_STATUS")
    }
    context, err := dk7.Context(before)
    if err != nil {
        return nil, "", err
    }
    cells := e.Native.Counts[context]
    least := cells[0].Count
    for action := 1; action < len(cells); action++ {
        if cells[action].Count < least {
            least = cells[action].Count
        }
    }
    digest := sha256.Sum256([]byte(context))
    offset := int(digest[0]) % len(cells)
    for i := 0; i < len(cells); i++ {
        action := (offset + i) % len(cells)
        if cells[action].Count == least {
            return &action, "SELECT_NATIVE_LEAST_OBSERVED_ACTION", nil
        }
    }
    return nil, "", errors.New("DK8_INTERNAL_COUNT_INVARIANT")
}

func (e *Explorer) BeginEpisode() error {
    if e == nil || e.Native == nil {
        return errors.New("DK8_MISSING_NATIVE_MEMORY")
    }
    return e.Native.BeginEpisode()
}

func (e *Explorer) Commit(before dk7.Public, action int) error {
    if e == nil || e.Native == nil {
        return errors.New("DK8_MISSING_NATIVE_MEMORY")
    }
    return e.Native.Commit(before, action)
}

func (e *Explorer) Feedback(before dk7.Public, action int, after dk7.Public, reward float64, terminal bool) error {
    if e == nil || e.Native == nil {
        return errors.New("DK8_MISSING_NATIVE_MEMORY")
    }
    return e.Native.Feedback(before, action, after, reward, terminal)
}

func (e *Explorer) ObserveFit(before dk7.Public, action int, after dk7.Public, reward float64, terminal bool) error {
    if e == nil || e.Native == nil {
        return errors.New("DK8_MISSING_NATIVE_MEMORY")
    }
    return e.Native.ObserveFit(before, action, after, reward, terminal)
}

func (e *Explorer) Snapshot() ([]byte, error) {
    if e == nil || e.Native == nil {
        return nil, errors.New("DK8_MISSING_NATIVE_MEMORY")
    }
    return e.Native.Snapshot()
}

func Restore(raw []byte) (*Explorer, error) {
    native, err := dk7.Restore(raw)
    if err != nil {
        return nil, err
    }
    return &Explorer{Native: native}, nil
}
