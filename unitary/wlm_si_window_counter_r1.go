package unitary

type wlmWindowRow struct {
	Seed uint32
	ScheduleFamily string
	InitialWindowTotal int
	AdmissionBudget int
	Substrate string
	Tick int
	RequestedAdmissions int
	AcceptedAdmissions int
	RejectedAdmissions int
	ExpiredCount int
	WindowTotalBefore int
	WindowTotalAfter int
	CumulativeAccepted int
	CumulativeRejected int
	CumulativeExpired int
	CanonicalBuckets []int
}

type wlmWindowResult struct {
	Rows []wlmWindowRow
	Metrics map[string]float64
}

type wlmWindowState interface {
	Total() int
	Expire(int) int
	Admit(int,int)
	Canonical(int) []int
}

type wlmWindowDeque struct {
	events []int
}

func (s *wlmWindowDeque) Total() int { return len(s.events) }
func (s *wlmWindowDeque) Expire(tick int) int {
	cutoff:=tick-16
	n:=0
	for n<len(s.events)&&s.events[n]<=cutoff{n++}
	if n>0{s.events=append([]int(nil),s.events[n:]...)}
	return n
}
func (s *wlmWindowDeque) Admit(tick,count int) {
	for i:=0;i<count;i++{s.events=append(s.events,tick)}
}
func (s *wlmWindowDeque) Canonical(tick int) []int {
	out:=make([]int,16)
	start:=tick-15
	for _,t:=range s.events{
		if t>=start&&t<=tick{out[t-start]++}
	}
	return out
}

type wlmWindowBuckets struct {
	counts [16]int
	stamps [16]int
	total int
}

func newWlmWindowBuckets() *wlmWindowBuckets {
	s:=&wlmWindowBuckets{}
	for i:=range s.stamps{s.stamps[i]=-1<<30}
	return s
}
func (s *wlmWindowBuckets) Total() int { return s.total }
func (s *wlmWindowBuckets) Expire(tick int) int {
	expiredTick:=tick-16
	slot:=expiredTick&15
	if s.stamps[slot]!=expiredTick{return 0}
	n:=s.counts[slot]
	s.counts[slot]=0
	s.stamps[slot]=-1<<30
	s.total-=n
	return n
}
func (s *wlmWindowBuckets) Admit(tick,count int) {
	slot:=tick&15
	if s.stamps[slot]!=tick{
		if s.counts[slot]!=0{panic("bucket overwrite before expiry")}
		s.stamps[slot]=tick
	}
	s.counts[slot]+=count
	s.total+=count
}
func (s *wlmWindowBuckets) Canonical(tick int) []int {
	out:=make([]int,16)
	start:=tick-15
	for i:=0;i<16;i++{
		t:=start+i
		slot:=t&15
		if s.stamps[slot]==t{out[i]=s.counts[slot]}
	}
	return out
}

func RunWlmSiWindowCounterR1() interface{} {
	seeds:=[...]uint32{2311,2339,2371,2393,2411,2441,2473,2503}
	families:=[...]string{"steady","burst","alternating","lcg32"}
	initials:=[...]int{0,16,32,48}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{
		"paired_case_count":0,
		"completed_substrate_case_runs":0,
		"total_emitted_rows":0,
		"paired_observable_mismatch_cases":0,
		"row_count_mismatch_cases":0,
		"state_cardinality_mismatch_cases":0,
		"invalid_operation_schedules":0,
		"budget_overrun_rows":0,
		"law_violation_cases":0,
		"max_conservation_error":0,
		"max_residual_law_error":0,
	}
	rows:=make([]wlmWindowRow,0,65536)
	for si,seed:=range seeds{
		for _,family:=range families{
			requests,valid:=wlmWindowRequests(seed,family)
			if !valid{metrics["invalid_operation_schedules"]++;continue}
			for _,initial:=range initials{
				for _,budget:=range budgets{
					metrics["paired_case_count"]++
					var a,b []wlmWindowRow
					var badA,badB bool
					if si%2==0{
						a,badA=wlmWindowRun(seed,family,initial,budget,requests,"event-deque",true,metrics)
						b,badB=wlmWindowRun(seed,family,initial,budget,requests,"circular-buckets",false,metrics)
					}else{
						b,badB=wlmWindowRun(seed,family,initial,budget,requests,"circular-buckets",false,metrics)
						a,badA=wlmWindowRun(seed,family,initial,budget,requests,"event-deque",true,metrics)
					}
					rows=append(rows,a...)
					rows=append(rows,b...)
					metrics["completed_substrate_case_runs"]+=2
					if badA||badB{metrics["law_violation_cases"]++}
					if len(a)!=len(b){metrics["row_count_mismatch_cases"]++}
					if !wlmWindowRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++}
					if !wlmWindowTotalsEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
				}
			}
		}
	}
	metrics["total_emitted_rows"]=float64(len(rows))
	return wlmWindowResult{Rows:rows,Metrics:metrics}
}

func wlmWindowRequests(seed uint32,family string)([]int,bool){
	out:=make([]int,64)
	state:=seed
	for tick:=0;tick<64;tick++{
		switch family{
		case "steady":
			out[tick]=1+int(seed%8)
		case "burst":
			if (tick+int(seed))%4==0{out[tick]=8}else{out[tick]=0}
		case "alternating":
			if (tick+int(seed))%2==0{out[tick]=8}else{out[tick]=1}
		case "lcg32":
			state=state*1664525+1013904223
			out[tick]=int(state%9)
		default:
			return out,false
		}
		if out[tick]<0||out[tick]>8{return out,false}
	}
	return out,true
}

func wlmWindowInitial(seed uint32,total int) []int {
	counts:=make([]int,16)
	for i:=0;i<total;i++{
		slot:=(i*5+int(seed))&15
		counts[slot]++
	}
	return counts
}

func wlmWindowRun(seed uint32,family string,initial,budget int,requests []int,substrate string,useDeque bool,metrics map[string]float64)([]wlmWindowRow,bool){
	var state wlmWindowState
	if useDeque{state=&wlmWindowDeque{}}else{state=newWlmWindowBuckets()}
	initialCounts:=wlmWindowInitial(seed,initial)
	for i,count:=range initialCounts{
		tick:=-15+i
		state.Admit(tick,count)
	}
	if state.Total()!=initial{panic("initial window total mismatch")}
	rows:=make([]wlmWindowRow,0,64)
	ca,cr,ce:=0,0,0
	bad:=false
	for tick:=0;tick<64;tick++{
		before:=state.Total()
		expired:=state.Expire(tick)
		requested:=requests[tick]
		accepted:=requested
		if accepted>budget{accepted=budget}
		rejected:=requested-accepted
		if accepted>budget{metrics["budget_overrun_rows"]++;bad=true}
		state.Admit(tick,accepted)
		after:=state.Total()
		expected:=before-expired+accepted
		conservation:=wlmWindowAbs(after-expected)
		canonical:=state.Canonical(tick)
		sum:=0;for _,v:=range canonical{sum+=v}
		residual:=wlmWindowAbs(sum-after)
		wlmWindowMaxMetric(metrics,"max_conservation_error",conservation)
		wlmWindowMaxMetric(metrics,"max_residual_law_error",residual)
		if conservation!=0||residual!=0||after<0{bad=true}
		ca+=accepted;cr+=rejected;ce+=expired
		rows=append(rows,wlmWindowRow{
			Seed:seed,ScheduleFamily:family,InitialWindowTotal:initial,AdmissionBudget:budget,Substrate:substrate,Tick:tick,
			RequestedAdmissions:requested,AcceptedAdmissions:accepted,RejectedAdmissions:rejected,ExpiredCount:expired,
			WindowTotalBefore:before,WindowTotalAfter:after,CumulativeAccepted:ca,CumulativeRejected:cr,CumulativeExpired:ce,
			CanonicalBuckets:canonical,
		})
	}
	return rows,bad
}

func wlmWindowRowsEqual(a,b []wlmWindowRow) bool {
	if len(a)!=len(b){return false}
	for i:=range a{
		x,y:=a[i],b[i]
		if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialWindowTotal!=y.InitialWindowTotal||
			x.AdmissionBudget!=y.AdmissionBudget||x.Tick!=y.Tick||x.RequestedAdmissions!=y.RequestedAdmissions||
			x.AcceptedAdmissions!=y.AcceptedAdmissions||x.RejectedAdmissions!=y.RejectedAdmissions||
			x.ExpiredCount!=y.ExpiredCount||x.WindowTotalBefore!=y.WindowTotalBefore||x.WindowTotalAfter!=y.WindowTotalAfter||
			x.CumulativeAccepted!=y.CumulativeAccepted||x.CumulativeRejected!=y.CumulativeRejected||x.CumulativeExpired!=y.CumulativeExpired||
			!wlmWindowIntsEqual(x.CanonicalBuckets,y.CanonicalBuckets){return false}
	}
	return true
}
func wlmWindowTotalsEqual(a,b []wlmWindowRow)bool{
	if len(a)!=len(b){return false}
	for i:=range a{if a[i].WindowTotalBefore!=b[i].WindowTotalBefore||a[i].WindowTotalAfter!=b[i].WindowTotalAfter{return false}}
	return true
}
func wlmWindowIntsEqual(a,b []int)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmWindowAbs(v int)int{if v<0{return -v};return v}
func wlmWindowMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
