package nvidiaqual

// CandidateTask is a bounded work item used only by the NVIDIA provider qualification.
type CandidateTask struct {
	ID        string
	Priority  int
	DependsOn []string
}

// ReadyTaskOrder returns the IDs of unfinished tasks whose dependencies are all
// completed, ordered by descending priority and then ascending ID.
//
// Qualification requirements are intentionally enforced by hidden tests.
func ReadyTaskOrder(items []CandidateTask, completed map[string]bool) ([]string, error) {
	panic("not implemented")
}
