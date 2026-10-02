package unitary

type wlmAllocEvent struct {
	Kind string
	RequestedSlot int
	AllocatedSlot int
	Success bool
}
type wlmAllocRow struct {
	Seed uint32
	ScheduleFamily string
	InitialOccupancy int
	OperationBudget int
	Substrate string
	Step int
	ExecutedOperations int
	AllocateSuccessCount int
	AllocateRejectCount int
	FreeSuccessCount int
	FreeMissingCount int
	OccupancyBefore int
	OccupancyAfter int
	CumulativeAllocated int
	CumulativeRejected int
	CumulativeFreed int
	OccupiedSlots []int
	Events []wlmAllocEvent
}
type wlmAllocResult struct {
	Rows []wlmAllocRow
	Metrics map[string]float64
}
type wlmAllocOp struct {
	Kind byte
	Slot int
}
type wlmAllocator interface {
	Occupancy() int
	Allocate() (int,bool)
	Free(int) bool
	SnapshotOccupied() []int
}

type wlmBitmapAllocator struct {
	used [64]bool
	occupied int
}
func (a *wlmBitmapAllocator) Occupancy() int { return a.occupied }
func (a *wlmBitmapAllocator) Allocate() (int,bool) {
	if a.occupied>=64 { return -1,false }
	for i:=0;i<64;i++ {
		if !a.used[i] {
			a.used[i]=true
			a.occupied++
			return i,true
		}
	}
	return -1,false
}
func (a *wlmBitmapAllocator) Free(slot int) bool {
	if slot<0||slot>=64||!a.used[slot] { return false }
	a.used[slot]=false
	a.occupied--
	return true
}
func (a *wlmBitmapAllocator) SnapshotOccupied() []int {
	out:=make([]int,0,a.occupied)
	for i:=0;i<64;i++ { if a.used[i] { out=append(out,i) } }
	return out
}

type wlmFreeListAllocator struct {
	free []int
}
func newWlmFreeListAllocator() *wlmFreeListAllocator {
	free:=make([]int,64)
	for i:=0;i<64;i++ { free[i]=i }
	return &wlmFreeListAllocator{free:free}
}
func (a *wlmFreeListAllocator) Occupancy() int { return 64-len(a.free) }
func (a *wlmFreeListAllocator) find(slot int) (int,bool) {
	lo,hi:=0,len(a.free)
	for lo<hi {
		m:=(lo+hi)/2
		if a.free[m]<slot { lo=m+1 } else { hi=m }
	}
	return lo,lo<len(a.free)&&a.free[lo]==slot
}
func (a *wlmFreeListAllocator) Allocate() (int,bool) {
	if len(a.free)==0 { return -1,false }
	slot:=a.free[0]
	copy(a.free[0:],a.free[1:])
	a.free=a.free[:len(a.free)-1]
	return slot,true
}
func (a *wlmFreeListAllocator) Free(slot int) bool {
	if slot<0||slot>=64 { return false }
	i,isFree:=a.find(slot)
	if isFree { return false }
	a.free=append(a.free,0)
	copy(a.free[i+1:],a.free[i:])
	a.free[i]=slot
	return true
}
func (a *wlmFreeListAllocator) SnapshotOccupied() []int {
	out:=make([]int,0,64-len(a.free))
	fi:=0
	for slot:=0;slot<64;slot++ {
		for fi<len(a.free)&&a.free[fi]<slot { fi++ }
		if fi<len(a.free)&&a.free[fi]==slot { continue }
		out=append(out,slot)
	}
	return out
}

func RunWlmSiFreelistAllocatorR1() interface{} {
	seeds:=[...]uint32{3023,3049,3083,3119,3163,3191,3221,3253}
	families:=[...]string{"fill-drain","fragment","reuse-hot","lcg32"}
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
	rows:=make([]wlmAllocRow,0,65536)
	for si,seed:=range seeds {
		for _,family:=range families {
			for _,initial:=range initials {
				for _,budget:=range budgets {
					metrics["paired_case_count"]++
					var a,b []wlmAllocRow
					var badA,badB bool
					if si%2==0 {
						a,badA=wlmAllocRun(seed,family,initial,budget,"free-bitmap",true,metrics)
						b,badB=wlmAllocRun(seed,family,initial,budget,"ordered-free-list",false,metrics)
					} else {
						b,badB=wlmAllocRun(seed,family,initial,budget,"ordered-free-list",false,metrics)
						a,badA=wlmAllocRun(seed,family,initial,budget,"free-bitmap",true,metrics)
					}
					rows=append(rows,a...)
					rows=append(rows,b...)
					metrics["completed_substrate_case_runs"]+=2
					if badA||badB { metrics["law_violation_cases"]++ }
					if len(a)!=len(b) { metrics["row_count_mismatch_cases"]++ }
					if !wlmAllocRowsEqual(a,b) { metrics["paired_observable_mismatch_cases"]++ }
					if !wlmAllocCardinalitiesEqual(a,b) { metrics["state_cardinality_mismatch_cases"]++ }
				}
			}
		}
	}
	metrics["total_emitted_rows"]=float64(len(rows))
	return wlmAllocResult{Rows:rows,Metrics:metrics}
}

func wlmAllocRun(seed uint32,family string,initial,budget int,substrate string,bitmap bool,metrics map[string]float64)([]wlmAllocRow,bool) {
	var a wlmAllocator
	if bitmap { a=&wlmBitmapAllocator{} } else { a=newWlmFreeListAllocator() }
	for i:=0;i<initial;i++ {
		slot,ok:=a.Allocate()
		if !ok||slot!=i { panic("initial first-fit allocation failed") }
	}
	rows:=make([]wlmAllocRow,0,64)
	state:=seed
	ca,cr,cf:=0,0,0
	bad:=false
	for step:=0;step<64;step++ {
		ops,ok,next:=wlmAllocOps(seed,state,family,step,budget)
		state=next
		if !ok { metrics["invalid_operation_schedules"]++;return rows,true }
		if len(ops)>budget { metrics["budget_overrun_rows"]++;bad=true }
		before:=a.Occupancy()
		allocOK,allocReject,freeOK,freeMissing:=0,0,0,0
		events:=make([]wlmAllocEvent,0,budget)
		for _,op:=range ops {
			switch op.Kind {
			case 'A':
				slot,success:=a.Allocate()
				if success { allocOK++ } else { allocReject++ }
				events=append(events,wlmAllocEvent{Kind:"allocate",RequestedSlot:-1,AllocatedSlot:slot,Success:success})
			case 'F':
				success:=a.Free(op.Slot)
				if success { freeOK++ } else { freeMissing++ }
				events=append(events,wlmAllocEvent{Kind:"free",RequestedSlot:op.Slot,AllocatedSlot:-1,Success:success})
			default:
				metrics["invalid_operation_schedules"]++
				bad=true
			}
		}
		after:=a.Occupancy()
		snapshot:=a.SnapshotOccupied()
		conservation:=wlmAllocAbs(after-(before+allocOK-freeOK))
		residual:=wlmAllocAbs(len(snapshot)-after)
		wlmAllocMaxMetric(metrics,"max_conservation_error",conservation)
		wlmAllocMaxMetric(metrics,"max_residual_law_error",residual)
		if conservation!=0||residual!=0||after<0||after>64 { bad=true }
		ca+=allocOK
		cr+=allocReject
		cf+=freeOK
		rows=append(rows,wlmAllocRow{
			Seed:seed,ScheduleFamily:family,InitialOccupancy:initial,OperationBudget:budget,
			Substrate:substrate,Step:step,ExecutedOperations:len(ops),
			AllocateSuccessCount:allocOK,AllocateRejectCount:allocReject,
			FreeSuccessCount:freeOK,FreeMissingCount:freeMissing,
			OccupancyBefore:before,OccupancyAfter:after,
			CumulativeAllocated:ca,CumulativeRejected:cr,CumulativeFreed:cf,
			OccupiedSlots:snapshot,Events:events,
		})
	}
	return rows,bad
}

func wlmAllocOps(seed,state uint32,family string,step,budget int)([]wlmAllocOp,bool,uint32) {
	ops:=make([]wlmAllocOp,0,budget)
	for ordinal:=0;ordinal<budget;ordinal++ {
		var op wlmAllocOp
		switch family {
		case "fill-drain":
			if step%4<2 {
				op=wlmAllocOp{Kind:'A',Slot:-1}
			} else {
				op=wlmAllocOp{Kind:'F',Slot:(step*budget+ordinal+int(seed))&63}
			}
		case "fragment":
			if (step+ordinal)%2==0 {
				op=wlmAllocOp{Kind:'F',Slot:(step*3+ordinal*5+int(seed))&63}
			} else {
				op=wlmAllocOp{Kind:'A',Slot:-1}
			}
		case "reuse-hot":
			if (step+ordinal)%3==0 {
				op=wlmAllocOp{Kind:'F',Slot:(int(seed)+step+ordinal)&15}
			} else {
				op=wlmAllocOp{Kind:'A',Slot:-1}
			}
		case "lcg32":
			state=state*1664525+1013904223
			if state%3==0 {
				state=state*1664525+1013904223
				op=wlmAllocOp{Kind:'F',Slot:int(state&63)}
			} else {
				op=wlmAllocOp{Kind:'A',Slot:-1}
			}
		default:
			return ops,false,state
		}
		if op.Kind=='F'&&(op.Slot<0||op.Slot>=64) { return ops,false,state }
		ops=append(ops,op)
	}
	return ops,true,state
}

func wlmAllocRowsEqual(a,b []wlmAllocRow) bool {
	if len(a)!=len(b) { return false }
	for i:=range a {
		x,y:=a[i],b[i]
		if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialOccupancy!=y.InitialOccupancy||
			x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.ExecutedOperations!=y.ExecutedOperations||
			x.AllocateSuccessCount!=y.AllocateSuccessCount||x.AllocateRejectCount!=y.AllocateRejectCount||
			x.FreeSuccessCount!=y.FreeSuccessCount||x.FreeMissingCount!=y.FreeMissingCount||
			x.OccupancyBefore!=y.OccupancyBefore||x.OccupancyAfter!=y.OccupancyAfter||
			x.CumulativeAllocated!=y.CumulativeAllocated||x.CumulativeRejected!=y.CumulativeRejected||
			x.CumulativeFreed!=y.CumulativeFreed||!wlmAllocIntsEqual(x.OccupiedSlots,y.OccupiedSlots)||
			!wlmAllocEventsEqual(x.Events,y.Events) { return false }
	}
	return true
}
func wlmAllocCardinalitiesEqual(a,b []wlmAllocRow) bool {
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i].OccupancyBefore!=b[i].OccupancyBefore||a[i].OccupancyAfter!=b[i].OccupancyAfter{return false} }
	return true
}
func wlmAllocIntsEqual(a,b []int) bool {
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i]!=b[i] { return false } }
	return true
}
func wlmAllocEventsEqual(a,b []wlmAllocEvent) bool {
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i]!=b[i] { return false } }
	return true
}
func wlmAllocAbs(v int) int { if v<0{return -v};return v }
func wlmAllocMaxMetric(m map[string]float64,k string,v int) { if float64(v)>m[k]{m[k]=float64(v)} }
