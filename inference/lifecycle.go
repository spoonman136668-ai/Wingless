package inference

import (
	"fmt"
	"sync"
)

type State string

const (
	Unavailable State = "UNAVAILABLE"
	Stopped     State = "STOPPED"
	Starting    State = "STARTING"
	Ready       State = "READY"
	Busy        State = "BUSY"
	Idle        State = "IDLE"
	Draining    State = "DRAINING"
	Failed      State = "FAILED"
)

type Lifecycle struct {
	mu    sync.Mutex
	state State
}

func NewLifecycle() *Lifecycle    { return &Lifecycle{state: Stopped} }
func (l *Lifecycle) State() State { l.mu.Lock(); defer l.mu.Unlock(); return l.state }
func (l *Lifecycle) Move(next State) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	allowed := map[State][]State{Stopped: {Starting, Unavailable}, Unavailable: {Stopped}, Starting: {Ready, Failed}, Ready: {Busy, Draining, Failed}, Busy: {Idle, Failed}, Idle: {Busy, Draining, Failed}, Draining: {Stopped, Failed}, Failed: {Stopped}}
	for _, s := range allowed[l.state] {
		if s == next {
			l.state = next
			return nil
		}
	}
	return fmt.Errorf("invalid lifecycle transition %s -> %s", l.state, next)
}
