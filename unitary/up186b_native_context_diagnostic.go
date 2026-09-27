package unitary

import "math"

const UP186BNativeContextSchema="wingless.up186b-native-context-diagnostic.v1"

type UP186BPhase struct{
	Phase int `json:"phase"`
	ActualCrossings int `json:"actual_crossings"`
	MeanMargin float64 `json:"mean_margin"`
	MeanAbsMargin float64 `json:"mean_abs_margin"`
	MinAbsMargin float64 `json:"min_abs_margin"`
	NearBoundaryCount int `json:"near_boundary_count"`
	CorrectCount int `json:"correct_count"`
	GateWeightL2 float64 `json:"gate_weight_l2"`
	MeanBias0 float64 `json:"mean_bias_0"`
	MeanBias1 float64 `json:"mean_bias_1"`
	MeanBias2 float64 `json:"mean_bias_2"`
	Examples int `json:"examples"`
}
type UP186BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP185BSeal string `json:"source_up185b_seal"`
	Phases []int `json:"phases"`
	FittingUsed bool `json:"fitting_used"`
	MaintenanceTriggered bool `json:"maintenance_triggered"`
	Points []UP186BPhase `json:"points"`
}
func up186bGateNorm(g *up129bGate)float64{
	s:=0.0
	for k:=0;k<3;k++{for i:=0;i<128;i++{s+=g.w[k][i]*g.w[k][i]}}
	return math.Sqrt(s)
}
func up186bCrossings(base *up129bGate,phase int,o,r [64]float64)int{
	paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
	total:=0
	for _,path:=range paths{
		g:=*base
		for _,idx:=range path{
			before:=up165bMargins(&g,o,r)
			items:=up158bItems(phase,[]int{idx});up156bApplyOld(&g,items,0,1,o,r)
			after:=up165bMargins(&g,o,r)
			for i:=range before{if before[i].correct!=after[i].correct{total++}}
		}
	}
	return total
}
func RunUP186B()(UP186BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	phases:=[]int{26,27,28,29,30,31,32,33}
	res:=UP186BResult{Schema:UP186BNativeContextSchema,Experiment:"UP-186B-native-context-diagnostic",SourceUP185BSeal:"0316209ffd00c60d12a89066935e7bf1e26e2234",Phases:phases,FittingUsed:false,MaintenanceTriggered:false}
	for _,phase:=range phases{
		p:=UP186BPhase{Phase:phase,MinAbsMargin:math.Inf(1)}
		subjects:=0
		for _,subject:=range up130bNewNames{
			g:=up169bPostReport(common,subject,phase,o,r);subjects++
			ev:=up165bMargins(g,o,r)
			for _,e:=range ev{
				p.MeanMargin+=e.margin;a:=math.Abs(e.margin);p.MeanAbsMargin+=a
				if a<p.MinAbsMargin{p.MinAbsMargin=a};if a<=0.05{p.NearBoundaryCount++};if e.correct{p.CorrectCount++};p.Examples++
			}
			p.GateWeightL2+=up186bGateNorm(g);p.MeanBias0+=g.b[0];p.MeanBias1+=g.b[1];p.MeanBias2+=g.b[2]
			p.ActualCrossings+=up186bCrossings(g,phase,o,r)
		}
		n:=float64(p.Examples);p.MeanMargin/=n;p.MeanAbsMargin/=n
		s:=float64(subjects);p.GateWeightL2/=s;p.MeanBias0/=s;p.MeanBias1/=s;p.MeanBias2/=s
		res.Points=append(res.Points,p)
	}
	return res,nil
}
