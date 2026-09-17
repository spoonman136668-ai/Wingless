package cognitive

func (r Route) ExecutableInCR1A() bool {
	switch r {
	case RouteMemory, RouteSkill, RouteInferenceSingle, RouteInferenceMulti:
		return true
	default:
		return false
	}
}
