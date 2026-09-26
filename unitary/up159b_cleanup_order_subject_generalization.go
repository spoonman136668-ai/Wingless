package unitary

const UP159BSubjectGeneralizationSchema="wingless.up159b-cleanup-order-subject-generalization.v1"

type UP159BResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP158BSeal string `json:"source_up158b_seal"`
	Subjects int `json:"subjects"`
	Orders int `json:"orders"`
	Arms int `json:"arms"`
	RehearsalBefore int `json:"rehearsal_before"`
	RehearsalAfter int `json:"rehearsal_after"`
	CleanupStore int `json:"cleanup_store"`
	CleanupObserve int `json:"cleanup_observe"`
	CleanupReport int `json:"cleanup_report"`
	ExtraUpdatesUsed bool `json:"extra_updates_used"`
	AdaptiveOrderingUsed bool `json:"adaptive_ordering_used"`
	Metrics []UP158BArm `json:"metrics"`
}
func RunUP159B()(UP159BResult,error){
	o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r)
	orders:=[]string{"STORE_OBSERVE_REPORT","STORE_REPORT_OBSERVE","OBSERVE_STORE_REPORT","OBSERVE_REPORT_STORE","REPORT_STORE_OBSERVE","REPORT_OBSERVE_STORE"}
	res:=UP159BResult{
		Schema:UP159BSubjectGeneralizationSchema,Experiment:"UP-159B-cleanup-order-subject-generalization",
		SourceUP158BSeal:"50d4146aa51fca6289a8271ca1e92b9a8ce801bb",
		Subjects:6,Orders:6,Arms:36,RehearsalBefore:3,RehearsalAfter:12,CleanupStore:5,CleanupObserve:5,CleanupReport:2,
		ExtraUpdatesUsed:false,AdaptiveOrderingUsed:false,
	}
	for _,subject:=range up130bNewNames{
		for _,order:=range orders{res.Metrics=append(res.Metrics,up158bRun(common,subject,order,o,r))}
	}
	return res,nil
}
