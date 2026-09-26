package unitary

const UP170BMonotonicSchema="wingless.up170b-margin-rank-monotonic-risk.v1"

type UP170BBinMetric struct{
	Class string `json:"class"`
	RankStart int `json:"rank_start"`
	RankEnd int `json:"rank_end"`
	ProbeStates int `json:"probe_states"`
	ExampleSlots int `json:"example_slots"`
	Crossings int `json:"crossings"`
	CrossingDensity float64 `json:"crossing_density"`
}
type UP170BClassTrend struct{
	Class string `json:"class"`
	NonIncreasingAcrossBins bool `json:"non_increasing_across_bins"`
}
type UP170BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP169BSeal string `json:"source_up169b_seal"`
	ContextPhase int `json:"context_phase"`
	NewPrefixRotation int `json:"new_prefix_rotation"`
	PrefixExamples int `json:"prefix_examples"`
	Subjects int `json:"subjects"`
	Paths int `json:"paths"`
	ProbeSurfaces int `json:"probe_surfaces"`
	ExamplesPerState int `json:"examples_per_state"`
	RankBins [][2]int `json:"rank_bins"`
	AdaptiveRotationUsed bool `json:"adaptive_rotation_used"`
	ThresholdFitted bool `json:"threshold_fitted"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	BinMetrics []UP170BBinMetric `json:"bin_metrics"`
	ClassTrends []UP170BClassTrend `json:"class_trends"`
}
func up170bPostReport(common *up129bGate,subject string,epoch,rotation int,o,r [64]float64)*up129bGate{
	g:=*common;pair:=up159bPair(subject);reportBlock:=up154bReportBlock(subject)
	prefix:=[]up135bExample{}
	for _,e:=range up159bExamples(pair){if !up153bInBlock(e,reportBlock){prefix=append(prefix,e)}}
	if len(prefix)!=20{panic("UP170B_PREFIX_COUNT")}
	rot:=make([]up135bExample,20);for i:=0;i<20;i++{rot[i]=prefix[(rotation+i)%20]}
	for _,e:=range rot{up135bStep(&g,e,o,r)}
	pre:=up158bItems(epoch,[]int{10,11,12});up156bApplyOld(&g,pre,0,len(pre),o,r)
	for _,e:=range reportBlock{up135bStep(&g,e,o,r)}
	return &g
}
func up170bBin(rank int,bins [][2]int)int{for i,b:=range bins{if rank>=b[0]&&rank<=b[1]{return i}};return -1}
func RunUP170B()(UP170BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch,rotation:=17,7
	bins:=[][2]int{{1,4},{5,10},{11,20},{21,40},{41,120}}
	classes:=[]string{"STORE","OBSERVE","REPORT"}
	type acc struct{states int;cross []int}
	accs:=map[string]*acc{};for _,c:=range classes{accs[c]=&acc{cross:make([]int,len(bins))}}
	for _,subject:=range up130bNewNames{
		base:=up170bPostReport(common,subject,epoch,rotation,o,r);pa,pb:=up169bPaths(base,epoch,o,r)
		for idx:=0;idx<15;idx++{
			class:=up169bClass(idx)
			for _,p0:=range []*up129bGate{pa,pb}{
				p:=*p0;before:=up165bMargins(&p,o,r);items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&p,items,0,1,o,r);after:=up165bMargins(&p,o,r)
				x:=accs[class];x.states++
				for i:=range before{if before[i].correct!=after[i].correct{rank:=up166bRank(before,i);bi:=up170bBin(rank,bins);if bi>=0{x.cross[bi]++}}}
			}
		}
	}
	res:=UP170BResult{Schema:UP170BMonotonicSchema,Experiment:"UP-170B-margin-rank-monotonic-risk",SourceUP169BSeal:"938afd01c7da1ff384e1a3198c268df5e042271f",ContextPhase:epoch,NewPrefixRotation:rotation,PrefixExamples:20,Subjects:6,Paths:2,ProbeSurfaces:15,ExamplesPerState:120,RankBins:bins,AdaptiveRotationUsed:false,ThresholdFitted:false,MaintenanceTriggered:false}
	for _,class:=range classes{
		x:=accs[class];dens:=[]float64{}
		for i,b:=range bins{
			slots:=x.states*(b[1]-b[0]+1);d:=0.0;if slots>0{d=float64(x.cross[i])/float64(slots)}
			dens=append(dens,d);res.BinMetrics=append(res.BinMetrics,UP170BBinMetric{Class:class,RankStart:b[0],RankEnd:b[1],ProbeStates:x.states,ExampleSlots:slots,Crossings:x.cross[i],CrossingDensity:d})
		}
		nonInc:=true;for i:=1;i<len(dens);i++{if dens[i]>dens[i-1]{nonInc=false}}
		res.ClassTrends=append(res.ClassTrends,UP170BClassTrend{Class:class,NonIncreasingAcrossBins:nonInc})
	}
	return res,nil
}
