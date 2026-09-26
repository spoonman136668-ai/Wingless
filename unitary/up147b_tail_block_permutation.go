package unitary

const UP147BPermutationSchema="wingless.up147b-tail-block-permutation.v1"

type UP147BEpochMetric struct {
	Schedule string `json:"schedule"`
	BlockPosition int `json:"block_position"`
	Shift int `json:"shift"`
	Repetition int `json:"repetition"`
	PrefixDamage float64 `json:"prefix_damage"`
	AnchorRecovery float64 `json:"anchor_recovery"`
	TailDamage float64 `json:"tail_damage"`
	NetEpochChange float64 `json:"net_epoch_change"`
}
type UP147BCell struct {
	Schedule string `json:"schedule"`
	BlockPosition int `json:"block_position"`
	Shift int `json:"shift"`
	Samples int `json:"samples"`
	MeanPrefixDamage float64 `json:"mean_prefix_damage"`
	MeanAnchorRecovery float64 `json:"mean_anchor_recovery"`
	MeanTailDamage float64 `json:"mean_tail_damage"`
	MeanNetEpochChange float64 `json:"mean_net_epoch_change"`
}
type UP147BScheduleSummary struct {
	Schedule string `json:"schedule"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}
type UP147BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP146BSeal string `json:"source_up146b_seal"`
	Schedules int `json:"schedules"`
	BlocksPerSchedule int `json:"blocks_per_schedule"`
	EpochsPerBlock int `json:"epochs_per_block"`
	DistinctTailSets int `json:"distinct_tail_sets"`
	SwitchesPerSchedule int `json:"switches_per_schedule"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	EpochMetrics []UP147BEpochMetric `json:"epoch_metrics"`
	Cells []UP147BCell `json:"cells"`
	ScheduleSummaries []UP147BScheduleSummary `json:"schedule_summaries"`
}

func up147bBlocks(rotation int) []int {
	base:=[]int{0,6,12,18}
	out:=make([]int,4)
	for i:=0;i<4;i++ { out[i]=base[(i+rotation)%4] }
	return out
}

func up147bRun(rotation int,o,r [64]float64)([]UP147BEpochMetric,[]UP147BCell,UP147BScheduleSummary){
	name:="rotation_"+itoa(rotation)
	g:=up129bTrainGate(o,r)
	examples:=up135bNewExamples()
	blocks:=up147bBlocks(rotation)
	rows:=make([]UP147BEpochMetric,0,20)
	cells:=make([]UP147BCell,0,4)
	epoch:=0
	for bp,shift:=range blocks {
		sp,sa,st,sn:=0.0,0.0,0.0,0.0
		for rep:=0;rep<5;rep++ {
			before:=up143bOldMean(g,o,r)
			full:=up139bRotate(examples,shift)
			for i:=0;i<20;i++ { up135bStep(g,full[i],o,r) }
			prefix:=up143bOldMean(g,o,r)
			up133bTrainOldBlock(g,epoch,o,r)
			anchor:=up143bOldMean(g,o,r)
			for i:=20;i<24;i++ { up135bStep(g,full[i],o,r) }
			tail:=up143bOldMean(g,o,r)
			m:=UP147BEpochMetric{
				Schedule:name,BlockPosition:bp+1,Shift:shift,Repetition:rep+1,
				PrefixDamage:prefix-before,AnchorRecovery:anchor-prefix,
				TailDamage:tail-anchor,NetEpochChange:tail-before,
			}
			rows=append(rows,m)
			sp+=m.PrefixDamage;sa+=m.AnchorRecovery;st+=m.TailDamage;sn+=m.NetEpochChange
			epoch++
		}
		cells=append(cells,UP147BCell{
			Schedule:name,BlockPosition:bp+1,Shift:shift,Samples:5,
			MeanPrefixDamage:sp/5,MeanAnchorRecovery:sa/5,
			MeanTailDamage:st/5,MeanNetEpochChange:sn/5,
		})
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	return rows,cells,UP147BScheduleSummary{
		Schedule:name,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalMeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
	}
}

func RunUP147B()(UP147BResult,error){
	o,r:=up124bCompetitorDirections()
	res:=UP147BResult{
		Schema:UP147BPermutationSchema,Experiment:"UP-147B-tail-block-permutation",
		SourceUP146BSeal:"fe75f2ca45f4566b71afa14e18b439812514d9e3",
		Schedules:4,BlocksPerSchedule:4,EpochsPerBlock:5,DistinctTailSets:4,SwitchesPerSchedule:3,
		AdaptiveOrderingUsed:false,ExtraUpdatesUsed:false,
	}
	for rotation:=0;rotation<4;rotation++ {
		rows,cells,sum:=up147bRun(rotation,o,r)
		res.EpochMetrics=append(res.EpochMetrics,rows...)
		res.Cells=append(res.Cells,cells...)
		res.ScheduleSummaries=append(res.ScheduleSummaries,sum)
	}
	return res,nil
}
