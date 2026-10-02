package unitary

type accOp struct{ bin int; delta int }

type accState interface{
	Get(int) int
	Apply(int,int) bool
	Snapshot() [32]int
	Nonzero() int
	Total() int
}

type denseAcc struct{ bins [32]int; nonzero int; total int }
func(a *denseAcc)Get(i int)int{return a.bins[i]}
func(a *denseAcc)Apply(i,d int)bool{
	next:=a.bins[i]+d
	if next<0{return false}
	if a.bins[i]==0&&next>0{a.nonzero++}
	if a.bins[i]>0&&next==0{a.nonzero--}
	a.bins[i]=next;a.total+=d;return true
}
func(a *denseAcc)Snapshot()[32]int{return a.bins}
func(a *denseAcc)Nonzero()int{return a.nonzero}
func(a *denseAcc)Total()int{return a.total}

type sparseEntry struct{ bin int; count int }
type sparseAcc struct{ e []sparseEntry; total int }
func(a *sparseAcc)find(bin int)(int,bool){
	lo,hi:=0,len(a.e)
	for lo<hi{m:=(lo+hi)/2;if a.e[m].bin<bin{lo=m+1}else{hi=m}}
	return lo,lo<len(a.e)&&a.e[lo].bin==bin
}
func(a *sparseAcc)Get(bin int)int{_,ok:=a.find(bin);if !ok{return 0};i,_:=a.find(bin);return a.e[i].count}
func(a *sparseAcc)Apply(bin,d int)bool{
	i,ok:=a.find(bin);cur:=0;if ok{cur=a.e[i].count}
	next:=cur+d;if next<0{return false}
	if ok{
		if next==0{copy(a.e[i:],a.e[i+1:]);a.e=a.e[:len(a.e)-1]}else{a.e[i].count=next}
	}else if next>0{
		a.e=append(a.e,sparseEntry{});copy(a.e[i+1:],a.e[i:]);a.e[i]=sparseEntry{bin:bin,count:next}
	}
	a.total+=d;return true
}
func(a *sparseAcc)Snapshot()[32]int{var out [32]int;for _,e:=range a.e{out[e.bin]=e.count};return out}
func(a *sparseAcc)Nonzero()int{return len(a.e)}
func(a *sparseAcc)Total()int{return a.total}

func RunWlmSiAccumulatorR1()interface{}{
	seeds:=[...]uint32{1901,1931,1973,1999,2039,2063,2099,2131}
	families:=[...]string{"uniform","hotspot","signed-wave","lcg32"}
	initials:=[...]int{0,8,16,24}
	budgets:=[...]int{1,2,4,8}
	m:=map[string]float64{
		"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,
		"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,
		"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,
		"budget_overrun_rows":0,"law_violation_cases":0,
		"max_conservation_error":0,"max_residual_law_error":0,
	}
	for _,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		m["paired_case_count"]++
		a:=&denseAcc{};b:=&sparseAcc{}
		for bin:=0;bin<initial;bin++{v:=1+int((seed+uint32(bin))%3);a.Apply(bin,v);b.Apply(bin,v)}
		state:=seed;caseMismatch:=false;cardMismatch:=false;caseLaw:=false
		for step:=0;step<64;step++{
			ops,ok,next:=accOps(seed,state,family,step,budget);state=next
			if !ok{m["invalid_operation_schedules"]++;caseLaw=true;break}
			if len(ops)>budget{m["budget_overrun_rows"]++;caseLaw=true}
			beforeA,beforeB:=a.Total(),b.Total()
			appliedDeltaA,appliedDeltaB:=0,0
			for _,op:=range ops{
				oka:=a.Apply(op.bin,op.delta);okb:=b.Apply(op.bin,op.delta)
				if oka{appliedDeltaA+=op.delta};if okb{appliedDeltaB+=op.delta}
				if oka!=okb{caseMismatch=true}
			}
			afterA,afterB:=a.Total(),b.Total()
			ca:=absAcc(afterA-(beforeA+appliedDeltaA));cb:=absAcc(afterB-(beforeB+appliedDeltaB))
			if float64(ca)>m["max_conservation_error"]{m["max_conservation_error"]=float64(ca)}
			if float64(cb)>m["max_conservation_error"]{m["max_conservation_error"]=float64(cb)}
			sa,sb:=a.Snapshot(),b.Snapshot()
			sumA,sumB:=0,0
			prefixA,prefixB:=0,0
			for i:=0;i<32;i++{
				sumA+=sa[i];sumB+=sb[i];prefixA+=sa[i];prefixB+=sb[i]
				if sa[i]!=sb[i]||prefixA!=prefixB{caseMismatch=true}
			}
			ra:=absAcc(sumA-afterA);rb:=absAcc(sumB-afterB)
			if float64(ra)>m["max_residual_law_error"]{m["max_residual_law_error"]=float64(ra)}
			if float64(rb)>m["max_residual_law_error"]{m["max_residual_law_error"]=float64(rb)}
			if ca!=0||cb!=0||ra!=0||rb!=0||afterA<0||afterB<0{caseLaw=true}
			if a.Nonzero()!=b.Nonzero(){cardMismatch=true}
			if afterA!=afterB{caseMismatch=true}
			m["total_emitted_rows"]+=2
		}
		m["completed_substrate_case_runs"]+=2
		if caseMismatch{m["paired_observable_mismatch_cases"]++}
		if cardMismatch{m["state_cardinality_mismatch_cases"]++}
		if caseLaw{m["law_violation_cases"]++}
	}}}}
	return struct{Metrics map[string]float64 `json:"Metrics"`}{Metrics:m}
}

func accOps(seed,state uint32,family string,step,budget int)([]accOp,bool,uint32){
	ops:=make([]accOp,0,budget)
	for i:=0;i<budget;i++{
		op:=accOp{}
		switch family{
		case "uniform":
			op.bin=(step*budget+i)&31;op.delta=1
		case "hotspot":
			op.bin=(int(seed)+step+i)&7;if (step+i)%4==3{op.delta=-1}else{op.delta=1}
		case "signed-wave":
			op.bin=(step+i*5)&31;if ((step/4)+i)%2==0{op.delta=2}else{op.delta=-1}
		case "lcg32":
			state=state*1664525+1013904223;op.bin=int(state&31)
			state=state*1664525+1013904223;if state&1==0{op.delta=1}else{op.delta=-1}
		default:return ops,false,state
		}
		if op.bin<0||op.bin>=32||op.delta==0{return ops,false,state}
		ops=append(ops,op)
	}
	return ops,true,state
}
func absAcc(v int)int{if v<0{return -v};return v}
