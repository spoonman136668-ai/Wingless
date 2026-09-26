package unitary

const UP96CWidth1024HorizonSchema = "wingless.up96c-width1024-horizon.v1"

type UP96CPoint struct {
	ChurnWrites              int     `json:"churn_writes"`
	HotQueryAccuracy         float64 `json:"hot_query_accuracy"`
	HotSetExactAccuracy      float64 `json:"hot_set_exact_accuracy"`
	ColdQueryAccuracy        float64 `json:"cold_query_accuracy"`
	RecallEntriesUsed        int     `json:"recall_entries_used"`
	RejectedOneShotWrites    int     `json:"rejected_one_shot_writes"`
	FalsePositiveAdmissions  int     `json:"false_positive_churn_admissions"`
	FilterResets             int     `json:"filter_resets"`
	PolicyMetadataBytes      int     `json:"policy_metadata_bytes"`
	TotalBoundedMemoryBytes  int     `json:"total_bounded_memory_bytes"`
}

type UP96CWidth1024HorizonResult struct {
	Schema                     string       `json:"schema"`
	Experiment                 string       `json:"experiment"`
	SourceUP95CSeal            string       `json:"source_up95c_seal"`
	ActiveKeys                 int          `json:"active_keys"`
	HotKeys                    int          `json:"hot_keys"`
	FilterBits                 int          `json:"filter_bits"`
	HashCount                  int          `json:"hash_count"`
	GenerationInterval         int          `json:"generation_interval"`
	ExactRecallCap             int          `json:"exact_recall_cap"`
	EntryPayloadBytes          int          `json:"entry_payload_bytes"`
	FutureOracleUsed           bool         `json:"future_oracle_used"`
	QueryLabelsUsedForAdmission bool        `json:"query_labels_used_for_admission"`
	Points                     []UP96CPoint `json:"points"`
}

func up96cRun(churn int,seedBases []int) UP96CPoint {
	const active=32
	const hotCount=16
	const width=1024
	hotHits,hotTotal,hotExactHits:=0,0,0
	coldHits,coldTotal:=0,0
	episodes,maxEntries,maxBytes:=0,0,0
	rejectedTotal,falsePositiveAdmissions,resets:=0,0,0

	for _,base:=range seedBases {
		for ep:=0;ep<64;ep++ {
			seed:=sq0Seed(base,3201+churn*61,ep)
			rng:=newSQ0RNG(seed)
			m:=&up95cMemory{width:width,seen:make([]uint64,width/64)}
			truth:=make([]int,active)
			for k:=0;k<active;k++ {
				v:=rng.intn(32);truth[k]=v;m.write(k,v)
			}
			hotWritten:=0
			for k:=0;k<active;k++ {
				if !up85cIsHot(k,hotCount){continue}
				m.write(k,truth[k]);hotWritten++
				if hotWritten%4==0 {
					for h:=0;h<active;h++ {if up85cIsHot(h,hotCount){m.query(h)}}
				}
			}
			for j:=0;j<churn;j++ {
				key:=active+j
				admitted,rejected:=m.write(key,rng.intn(32))
				if rejected{rejectedTotal++}
				if admitted{falsePositiveAdmissions++}
				if (j+1)%4==0 {
					for h:=0;h<active;h++ {if up85cIsHot(h,hotCount){m.query(h)}}
				}
			}
			hotExact:=true
			for k:=0;k<active;k++ {
				got,ok:=m.query(k)
				if up85cIsHot(k,hotCount){
					hotTotal++;if ok&&got==truth[k]{hotHits++}else{hotExact=false}
				}else{
					coldTotal++;if ok&&got==truth[k]{coldHits++}
				}
			}
			if hotExact{hotExactHits++}
			episodes++
			if m.mem.count>maxEntries{maxEntries=m.mem.count}
			if m.totalBytes()>maxBytes{maxBytes=m.totalBytes()}
			resets+=m.resets
		}
	}
	return UP96CPoint{
		ChurnWrites:churn,
		HotQueryAccuracy:float64(hotHits)/float64(hotTotal),
		HotSetExactAccuracy:float64(hotExactHits)/float64(episodes),
		ColdQueryAccuracy:float64(coldHits)/float64(coldTotal),
		RecallEntriesUsed:maxEntries,
		RejectedOneShotWrites:rejectedTotal,
		FalsePositiveAdmissions:falsePositiveAdmissions,
		FilterResets:resets,
		PolicyMetadataBytes:134,
		TotalBoundedMemoryBytes:maxBytes,
	}
}

func RunUP96C()(UP96CWidth1024HorizonResult,error){
	result:=UP96CWidth1024HorizonResult{
		Schema:UP96CWidth1024HorizonSchema,
		Experiment:"UP-96C-width1024-horizon",
		SourceUP95CSeal:"7f4192c8341908834d727ad97678c1080c865c8d",
		ActiveKeys:32,HotKeys:16,FilterBits:1024,HashCount:2,GenerationInterval:32,
		ExactRecallCap:16,EntryPayloadBytes:16,FutureOracleUsed:false,QueryLabelsUsedForAdmission:false,
	}
	seedBases:=[]int{181000000,182000000}
	for _,churn:=range []int{1536,3072,6144,12288} {
		result.Points=append(result.Points,up96cRun(churn,seedBases))
	}
	return result,nil
}
