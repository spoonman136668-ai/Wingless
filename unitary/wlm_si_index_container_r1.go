package unitary

type wlmIndexEntry struct { Key int; Value int }
type wlmIndexLookup struct { Key int; Hit bool; Value int }
type wlmIndexRow struct {
	Seed uint32; ScheduleFamily string; InitialCardinality int; OperationBudget int; Substrate string; Step int
	RequestedOperations int; ExecutedOperations int; InsertedCount int; UpdatedCount int; RejectedUpdateCount int
	DeleteHitCount int; DeleteMissCount int; LookupHitCount int; LookupMissCount int
	CardinalityBefore int; CardinalityAfter int; CumulativeInserted int; CumulativeUpdated int; CumulativeRejected int; CumulativeDeleted int
	Entries []wlmIndexEntry; Lookups []wlmIndexLookup
}
type wlmIndexResult struct { Rows []wlmIndexRow; Metrics map[string]float64 }
type wlmIndexOp struct { Kind byte; Key int; Value int }
type wlmIndex interface { Len() int; Set(int,int)(bool,bool,bool); Delete(int) bool; Get(int)(int,bool); Snapshot() []wlmIndexEntry }

type wlmDenseIndex struct { present [32]bool; values [32]int; size int }
func(d *wlmDenseIndex)Len()int{return d.size}
func(d *wlmDenseIndex)Set(k,v int)(bool,bool,bool){if d.present[k]{d.values[k]=v;return false,true,true};if d.size>=32{return false,false,false};d.present[k]=true;d.values[k]=v;d.size++;return true,false,true}
func(d *wlmDenseIndex)Delete(k int)bool{if !d.present[k]{return false};d.present[k]=false;d.values[k]=0;d.size--;return true}
func(d *wlmDenseIndex)Get(k int)(int,bool){if !d.present[k]{return 0,false};return d.values[k],true}
func(d *wlmDenseIndex)Snapshot()[]wlmIndexEntry{out:=make([]wlmIndexEntry,0,d.size);for k:=0;k<32;k++{if d.present[k]{out=append(out,wlmIndexEntry{Key:k,Value:d.values[k]})}};return out}

type wlmOpenIndex struct { state [32]uint8; keys [32]int; values [32]int; size int }
func(h *wlmOpenIndex)Len()int{return h.size}
func wlmIndexHash(k int)int{return ((k&7)*7+3)&31}
func(h *wlmOpenIndex)find(k int)(int,bool,int){start:=wlmIndexHash(k);firstTomb:=-1;for i:=0;i<32;i++{p:=(start+i)&31;switch h.state[p]{case 0:if firstTomb>=0{return p,false,firstTomb};return p,false,p;case 1:if h.keys[p]==k{return p,true,p};case 2:if firstTomb<0{firstTomb=p}}};return -1,false,firstTomb}
func(h *wlmOpenIndex)Set(k,v int)(bool,bool,bool){idx,found,slot:=h.find(k);if found{h.values[idx]=v;return false,true,true};if h.size>=32||slot<0{return false,false,false};h.state[slot]=1;h.keys[slot]=k;h.values[slot]=v;h.size++;return true,false,true}
func(h *wlmOpenIndex)Delete(k int)bool{idx,found,_:=h.find(k);if !found{return false};h.state[idx]=2;h.values[idx]=0;h.size--;return true}
func(h *wlmOpenIndex)Get(k int)(int,bool){idx,found,_:=h.find(k);if !found{return 0,false};return h.values[idx],true}
func(h *wlmOpenIndex)Snapshot()[]wlmIndexEntry{out:=make([]wlmIndexEntry,0,h.size);for i:=0;i<32;i++{if h.state[i]==1{out=append(out,wlmIndexEntry{Key:h.keys[i],Value:h.values[i]})}};for i:=1;i<len(out);i++{x:=out[i];j:=i-1;for j>=0&&out[j].Key>x.Key{out[j+1]=out[j];j--};out[j+1]=x};return out}

func RunWlmSiIndexContainerR1()interface{}{
	seeds:=[...]uint32{907,929,953,977,1009,1031,1061,1091};families:=[...]string{"sequential","collision","delete-reinsert","lcg32"};initials:=[...]int{0,8,16,24};budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,"max_conservation_error":0,"max_residual_law_error":0}
	rows:=make([]wlmIndexRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{metrics["paired_case_count"]++;var a,b []wlmIndexRow;var badA,badB bool;if si%2==0{a,badA=wlmIndexRun(seed,family,initial,budget,"dense-fixed-array",true,metrics);b,badB=wlmIndexRun(seed,family,initial,budget,"open-addressed",false,metrics)}else{b,badB=wlmIndexRun(seed,family,initial,budget,"open-addressed",false,metrics);a,badA=wlmIndexRun(seed,family,initial,budget,"dense-fixed-array",true,metrics)};rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2;if badA||badB{metrics["law_violation_cases"]++};if len(a)!=len(b){metrics["row_count_mismatch_cases"]++};if !wlmIndexRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++};if !wlmIndexCardinalitiesEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}}}}}
	metrics["total_emitted_rows"]=float64(len(rows));return wlmIndexResult{Rows:rows,Metrics:metrics}
}

func wlmIndexRun(seed uint32,family string,initial,budget int,substrate string,dense bool,metrics map[string]float64)([]wlmIndexRow,bool){
	var idx wlmIndex;if dense{idx=&wlmDenseIndex{}}else{idx=&wlmOpenIndex{}};for k:=0;k<initial;k++{idx.Set(k,int(seed)+k*13)}
	rows:=make([]wlmIndexRow,0,64);ci,cu,cr,cd:=0,0,0,0;bad:=false;state:=seed
	for step:=0;step<64;step++{ops,ok,next:=wlmIndexOps(seed,state,family,step,budget);state=next;if !ok{metrics["invalid_operation_schedules"]++;return rows,true};before:=idx.Len();inserted,updated,rejected,delHit,delMiss,lookHit,lookMiss:=0,0,0,0,0,0,0;lookups:=make([]wlmIndexLookup,0,budget);if len(ops)>budget{metrics["budget_overrun_rows"]++;bad=true};for _,op:=range ops{switch op.Kind{case 'S':ins,upd,accepted:=idx.Set(op.Key,op.Value);if ins{inserted++};if upd{updated++};if !accepted{rejected++};case 'D':if idx.Delete(op.Key){delHit++}else{delMiss++};case 'G':v,hit:=idx.Get(op.Key);if hit{lookHit++}else{lookMiss++};lookups=append(lookups,wlmIndexLookup{Key:op.Key,Hit:hit,Value:v});default:metrics["invalid_operation_schedules"]++;bad=true}};after:=idx.Len();expected:=before+inserted-delHit;conservation:=wlmIndexAbs(after-expected);residual:=wlmIndexAbs(len(idx.Snapshot())-after);wlmIndexMaxMetric(metrics,"max_conservation_error",conservation);wlmIndexMaxMetric(metrics,"max_residual_law_error",residual);if conservation!=0||residual!=0||after<0||after>32{bad=true};ci+=inserted;cu+=updated;cr+=rejected;cd+=delHit;rows=append(rows,wlmIndexRow{Seed:seed,ScheduleFamily:family,InitialCardinality:initial,OperationBudget:budget,Substrate:substrate,Step:step,RequestedOperations:budget,ExecutedOperations:len(ops),InsertedCount:inserted,UpdatedCount:updated,RejectedUpdateCount:rejected,DeleteHitCount:delHit,DeleteMissCount:delMiss,LookupHitCount:lookHit,LookupMissCount:lookMiss,CardinalityBefore:before,CardinalityAfter:after,CumulativeInserted:ci,CumulativeUpdated:cu,CumulativeRejected:cr,CumulativeDeleted:cd,Entries:idx.Snapshot(),Lookups:lookups})};return rows,bad
}
func wlmIndexOps(seed,state uint32,family string,step,budget int)([]wlmIndexOp,bool,uint32){ops:=make([]wlmIndexOp,0,budget);for ordinal:=0;ordinal<budget;ordinal++{var kind byte;var key,value int;switch family{case "sequential":key=(step*budget+ordinal)&31;switch step%3{case 0:kind='S';case 1:kind='G';default:kind='D'};value=int(seed)+step*97+ordinal;case "collision":key=int(seed&7)+8*((step+ordinal)&3);switch(step+ordinal)%3{case 0:kind='S';case 1:kind='G';default:kind='D'};value=int(seed)+step*71+ordinal*3;case "delete-reinsert":key=(step+ordinal)&31;if step%2==0{kind='D'}else{kind='S'};value=int(seed)+step*53+ordinal*5;case "lcg32":state=state*1664525+1013904223;key=int(state&31);state=state*1664525+1013904223;switch state%3{case 0:kind='S';case 1:kind='D';default:kind='G'};state=state*1664525+1013904223;value=int(state&0x7fffffff);default:return ops,false,state};if key<0||key>=32{return ops,false,state};ops=append(ops,wlmIndexOp{Kind:kind,Key:key,Value:value})};return ops,true,state}
func wlmIndexRowsEqual(a,b []wlmIndexRow)bool{if len(a)!=len(b){return false};for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialCardinality!=y.InitialCardinality||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.RequestedOperations!=y.RequestedOperations||x.ExecutedOperations!=y.ExecutedOperations||x.InsertedCount!=y.InsertedCount||x.UpdatedCount!=y.UpdatedCount||x.RejectedUpdateCount!=y.RejectedUpdateCount||x.DeleteHitCount!=y.DeleteHitCount||x.DeleteMissCount!=y.DeleteMissCount||x.LookupHitCount!=y.LookupHitCount||x.LookupMissCount!=y.LookupMissCount||x.CardinalityBefore!=y.CardinalityBefore||x.CardinalityAfter!=y.CardinalityAfter||x.CumulativeInserted!=y.CumulativeInserted||x.CumulativeUpdated!=y.CumulativeUpdated||x.CumulativeRejected!=y.CumulativeRejected||x.CumulativeDeleted!=y.CumulativeDeleted||!wlmIndexEntriesEqual(x.Entries,y.Entries)||!wlmIndexLookupsEqual(x.Lookups,y.Lookups){return false}};return true}
func wlmIndexCardinalitiesEqual(a,b []wlmIndexRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].CardinalityBefore!=b[i].CardinalityBefore||a[i].CardinalityAfter!=b[i].CardinalityAfter||len(a[i].Entries)!=a[i].CardinalityAfter||len(b[i].Entries)!=b[i].CardinalityAfter{return false}};return true}
func wlmIndexEntriesEqual(a,b []wlmIndexEntry)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmIndexLookupsEqual(a,b []wlmIndexLookup)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmIndexAbs(v int)int{if v<0{return -v};return v}
func wlmIndexMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
