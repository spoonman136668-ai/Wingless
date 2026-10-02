package unitary

type poolOp struct{ kind byte; slot int }

type poolState interface {
	Len() int
	Alloc() (int,bool)
	Free(int) bool
	Snapshot() uint32
}

type freeListPool struct{
	free []int
	used [32]bool
	n int
}
func newFreeListPool()*freeListPool{
	p:=&freeListPool{free:make([]int,32)}
	for i:=0;i<32;i++{p.free[i]=i}
	return p
}
func(p *freeListPool)Len()int{return p.n}
func(p *freeListPool)Alloc()(int,bool){
	if len(p.free)==0{return -1,false}
	s:=p.free[0];p.free=p.free[1:];p.used[s]=true;p.n++;return s,true
}
func(p *freeListPool)Free(s int)bool{
	if s<0||s>=32||!p.used[s]{return false}
	p.used[s]=false;p.n--
	i:=0;for i<len(p.free)&&p.free[i]<s{i++}
	p.free=append(p.free,0);copy(p.free[i+1:],p.free[i:]);p.free[i]=s
	return true
}
func(p *freeListPool)Snapshot()uint32{
	var bits uint32
	for i:=0;i<32;i++{if p.used[i]{bits|=uint32(1)<<uint(i)}}
	return bits
}

type bitmapPool struct{ used uint32; n int }
func(p *bitmapPool)Len()int{return p.n}
func(p *bitmapPool)Alloc()(int,bool){
	if p.n>=32{return -1,false}
	for s:=0;s<32;s++{
		m:=uint32(1)<<uint(s)
		if p.used&m==0{p.used|=m;p.n++;return s,true}
	}
	return -1,false
}
func(p *bitmapPool)Free(s int)bool{
	if s<0||s>=32{return false}
	m:=uint32(1)<<uint(s)
	if p.used&m==0{return false}
	p.used&^=m;p.n--;return true
}
func(p *bitmapPool)Snapshot()uint32{return p.used}

func RunWlmSiObjectPoolR1()interface{}{
	seeds:=[...]uint32{1663,1693,1723,1753,1787,1811,1847,1877}
	families:=[...]string{"fill-drain","fragment","reuse-wave","lcg32"}
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
		a:=newFreeListPool();b:=&bitmapPool{}
		valid:=true
		for i:=0;i<initial;i++{
			sa,oka:=a.Alloc();sb,okb:=b.Alloc()
			if !oka||!okb||sa!=i||sb!=i{valid=false}
		}
		state:=seed
		caseMismatch:=false
		cardMismatch:=false
		caseLaw:=false
		for step:=0;step<64;step++{
			ops,ok,next:=poolOps(seed,state,family,step,budget);state=next
			if !ok{m["invalid_operation_schedules"]++;valid=false;break}
			if len(ops)>budget{m["budget_overrun_rows"]++;valid=false}
			beforeA,beforeB:=a.Len(),b.Len()
			allocA,allocB,freeA,freeB:=0,0,0,0
			for _,op:=range ops{
				switch op.kind{
				case 'A':
					sa,oka:=a.Alloc();sb,okb:=b.Alloc()
					if oka{allocA++};if okb{allocB++}
					if oka!=okb||sa!=sb{caseMismatch=true}
				case 'F':
					oka:=a.Free(op.slot);okb:=b.Free(op.slot)
					if oka{freeA++};if okb{freeB++}
					if oka!=okb{caseMismatch=true}
				default:
					m["invalid_operation_schedules"]++;valid=false
				}
			}
			afterA,afterB:=a.Len(),b.Len()
			ca:=absPool(afterA-(beforeA+allocA-freeA))
			cb:=absPool(afterB-(beforeB+allocB-freeB))
			if float64(ca)>m["max_conservation_error"]{m["max_conservation_error"]=float64(ca)}
			if float64(cb)>m["max_conservation_error"]{m["max_conservation_error"]=float64(cb)}
			ra:=absPool(popcount32(a.Snapshot())-afterA)
			rb:=absPool(popcount32(b.Snapshot())-afterB)
			if float64(ra)>m["max_residual_law_error"]{m["max_residual_law_error"]=float64(ra)}
			if float64(rb)>m["max_residual_law_error"]{m["max_residual_law_error"]=float64(rb)}
			if ca!=0||cb!=0||ra!=0||rb!=0||afterA<0||afterA>32||afterB<0||afterB>32{caseLaw=true}
			if afterA!=afterB{cardMismatch=true}
			if a.Snapshot()!=b.Snapshot(){caseMismatch=true}
			m["total_emitted_rows"]+=2
		}
		m["completed_substrate_case_runs"]+=2
		if !valid||caseLaw{m["law_violation_cases"]++}
		if caseMismatch{m["paired_observable_mismatch_cases"]++}
		if cardMismatch{m["state_cardinality_mismatch_cases"]++}
	}}}}
	return struct{
		Metrics map[string]float64 `json:"Metrics"`
	}{Metrics:m}
}

func poolOps(seed,state uint32,family string,step,budget int)([]poolOp,bool,uint32){
	ops:=make([]poolOp,0,budget)
	for i:=0;i<budget;i++{
		op:=poolOp{}
		switch family{
		case "fill-drain":
			if (step/8)%2==0{op.kind='A'}else{op.kind='F';op.slot=(step*budget+i)&31}
		case "fragment":
			if (step+i)%3==0{op.kind='F';op.slot=(int(seed)+step*5+i*7)&31}else{op.kind='A'}
		case "reuse-wave":
			if step%2==0{op.kind='F';op.slot=(step+i*3)&15}else{op.kind='A'}
		case "lcg32":
			state=state*1664525+1013904223
			if state&1==0{op.kind='A'}else{op.kind='F';state=state*1664525+1013904223;op.slot=int(state&31)}
		default:return ops,false,state
		}
		if op.slot<0||op.slot>=32{return ops,false,state}
		ops=append(ops,op)
	}
	return ops,true,state
}
func absPool(v int)int{if v<0{return -v};return v}
func popcount32(x uint32)int{n:=0;for x!=0{x&=x-1;n++};return n}
