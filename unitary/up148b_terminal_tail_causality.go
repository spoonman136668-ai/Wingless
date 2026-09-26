package unitary

const UP148BTerminalSchema="wingless.up148b-terminal-tail-causality.v1"

type UP148BArm struct {
	Arm string `json:"arm"`
	TerminalShift int `json:"terminal_shift"`
	PrefixStateOldRetention float64 `json:"prefix_state_old_retention"`
	MeanTerminalPrefixDamage float64 `json:"mean_terminal_prefix_damage"`
	MeanTerminalAnchorRecovery float64 `json:"mean_terminal_anchor_recovery"`
	MeanTerminalTailDamage float64 `json:"mean_terminal_tail_damage"`
	MeanTerminalNetEpochChange float64 `json:"mean_terminal_net_epoch_change"`
	FinalOldHeldout UP129BSplitMetric `json:"final_old_heldout"`
	FinalOldUnseen UP129BSplitMetric `json:"final_old_unseen"`
	FinalMeanOldRetention float64 `json:"final_mean_old_retention"`
	FinalPrimaryNew UP130BSplitMetric `json:"final_primary_new"`
	FinalSecondaryNew UP130BSplitMetric `json:"final_secondary_new"`
	FinalMeanNewAccuracy float64 `json:"final_mean_new_accuracy"`
}

type UP148BResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP147BSeal string `json:"source_up147b_seal"`
	CommonPrefixEpochs int `json:"common_prefix_epochs"`
	TerminalEpochs int `json:"terminal_epochs"`
	LearningRate float64 `json:"learning_rate"`
	NewUpdatesPerEpoch int `json:"new_updates_per_epoch"`
	OldUpdatesPerEpoch int `json:"old_updates_per_epoch"`
	FinalNewUpdates int `json:"final_new_updates"`
	AdaptiveTailSelection bool `json:"adaptive_tail_selection"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	Arms []UP148BArm `json:"arms"`
	FinalOldRetentionDeltaSafeMinusDamaging float64 `json:"final_old_retention_delta_safe_minus_damaging"`
	FinalNewAccuracyDeltaSafeMinusDamaging float64 `json:"final_new_accuracy_delta_safe_minus_damaging"`
}

func up148bTrainEpoch(g *up129bGate,shift,epoch int,o,r [64]float64)(prefixDamage,anchorRecovery,tailDamage,net float64){
	examples:=up135bNewExamples()
	before:=up143bOldMean(g,o,r)
	full:=up139bRotate(examples,shift)
	for i:=0;i<20;i++{up135bStep(g,full[i],o,r)}
	prefix:=up143bOldMean(g,o,r)
	up133bTrainOldBlock(g,epoch,o,r)
	anchor:=up143bOldMean(g,o,r)
	for i:=20;i<24;i++{up135bStep(g,full[i],o,r)}
	tail:=up143bOldMean(g,o,r)
	return prefix-before,anchor-prefix,tail-anchor,tail-before
}

func up148bRun(arm string,terminalShift int,o,r [64]float64) UP148BArm {
	g:=up129bTrainGate(o,r)
	epoch:=0
	for _,shift:=range []int{0,6,12}{
		for rep:=0;rep<5;rep++{
			up148bTrainEpoch(g,shift,epoch,o,r)
			epoch++
		}
	}
	prefixState:=up143bOldMean(g,o,r)
	sp,sa,st,sn:=0.0,0.0,0.0,0.0
	for rep:=0;rep<5;rep++{
		p,a,t,n:=up148bTrainEpoch(g,terminalShift,epoch,o,r)
		sp+=p;sa+=a;st+=t;sn+=n;epoch++
	}
	oldH:=up129bSplit(g,"old_heldout",up121bOriginalNames[4:6],o,r)
	oldU:=up129bSplit(g,"old_unseen",up121bUnseenNames,o,r)
	p,_:=up130bSplit(g,"primary_new",up130bNewNames[4:6],o,r)
	s,_:=up130bSplit(g,"secondary_new",up131bSecondaryNames(),o,r)
	return UP148BArm{
		Arm:arm,TerminalShift:terminalShift,PrefixStateOldRetention:prefixState,
		MeanTerminalPrefixDamage:sp/5,MeanTerminalAnchorRecovery:sa/5,
		MeanTerminalTailDamage:st/5,MeanTerminalNetEpochChange:sn/5,
		FinalOldHeldout:oldH,FinalOldUnseen:oldU,FinalMeanOldRetention:(oldH.ClassAccuracy+oldU.ClassAccuracy)/2,
		FinalPrimaryNew:p,FinalSecondaryNew:s,FinalMeanNewAccuracy:(p.OverallAccuracy+s.OverallAccuracy)/2,
	}
}

func RunUP148B()(UP148BResult,error){
	o,r:=up124bCompetitorDirections()
	safe:=up148bRun("safe_terminal",18,o,r)
	dmg:=up148bRun("damaging_terminal",12,o,r)
	return UP148BResult{
		Schema:UP148BTerminalSchema,Experiment:"UP-148B-terminal-tail-causality",
		SourceUP147BSeal:"8a6c31cc4f3c3c439d8d297d953b847cfe8b53ae",
		CommonPrefixEpochs:15,TerminalEpochs:5,LearningRate:0.08,
		NewUpdatesPerEpoch:24,OldUpdatesPerEpoch:15,FinalNewUpdates:4,
		AdaptiveTailSelection:false,ExtraUpdatesUsed:false,
		Arms:[]UP148BArm{safe,dmg},
		FinalOldRetentionDeltaSafeMinusDamaging:safe.FinalMeanOldRetention-dmg.FinalMeanOldRetention,
		FinalNewAccuracyDeltaSafeMinusDamaging:safe.FinalMeanNewAccuracy-dmg.FinalMeanNewAccuracy,
	},nil
}
