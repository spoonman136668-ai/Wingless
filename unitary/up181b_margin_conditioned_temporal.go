package unitary

import "math"

const UP181BMarginConditionalSchema="wingless.up181b-margin-conditioned-temporal.v1"

type UP181BMetric struct{
 Scope string `json:"scope"`
 MarginBand string `json:"margin_band"`
 TemporalState string `json:"temporal_state"`
 Slots int `json:"slots"`
 Crossings int `json:"crossings"`
 CrossingDensity float64 `json:"crossing_density"`
}
type UP181BResult struct{
 Schema string `json:"schema"`
 Experiment string `json:"experiment"`
 SourceUP180BSeal string `json:"source_up180b_seal"`
 ContextPhase int `json:"context_phase"`
 MarginBands []string `json:"margin_bands"`
 TemporalStates []string `json:"temporal_states"`
 BandsFrozenBeforeRun bool `json:"bands_frozen_before_run"`
 AdaptiveBandsUsed bool `json:"adaptive_bands_used"`
 MaintenanceTriggered bool `json:"maintenance_triggered"`
 Metrics []UP181BMetric `json:"metrics"`
}
type up181bCount struct{slots,cross int}
func up181bBand(x float64)string{
 if x<0.05{return "m0_005"}
 if x<0.10{return "m005_010"}
 if x<0.20{return "m010_020"}
 if x<0.40{return "m020_040"}
 return "m040_plus"
}
func up181bTemporal(streak int)string{
 if streak==1{return "current_only"}
 if streak==2{return "transition_to_persistent"}
 if streak>=3{return "established_persistent"}
 return ""
}
func up181bRate(a,b int)float64{if b==0{return 0};return float64(a)/float64(b)}
func RunUP181B()(UP181BResult,error){
 o,r:=up124bCompetitorDirections();common:=up150bCommonPrefix(o,r);epoch:=26
 paths:=[][]int{{0,1,2,3,4,5,6,7,8,9,13,14},{5,6,7,8,9,0,1,2,3,4,13,14}}
 scopes:=[]string{"ALL","STORE","OBSERVE","REPORT"}
 mbands:=[]string{"m0_005","m005_010","m010_020","m020_040","m040_plus"}
 states:=[]string{"current_only","transition_to_persistent","established_persistent"}
 acc:=map[string]map[string]map[string]*up181bCount{}
 for _,s:=range scopes{acc[s]=map[string]map[string]*up181bCount{};for _,m:=range mbands{acc[s][m]=map[string]*up181bCount{};for _,t:=range states{acc[s][m][t]=&up181bCount{}}}}
 for _,subject:=range up130bNewNames{
  base:=up169bPostReport(common,subject,epoch,o,r)
  for _,path:=range paths{
   g:=*base;streak:=make([]int,120)
   for _,idx:=range path{
    before:=up165bMargins(&g,o,r);cls:=up171bClass(idx);band:=up171bBand(cls)
    for i:=range before{if up166bRank(before,i)<=band{streak[i]++}else{streak[i]=0}}
    items:=up158bItems(epoch,[]int{idx});up156bApplyOld(&g,items,0,1,o,r);after:=up165bMargins(&g,o,r)
    for i:=range before{
     ts:=up181bTemporal(streak[i]);if ts==""{continue}
     mb:=up181bBand(math.Abs(before[i].margin));changed:=before[i].correct!=after[i].correct
     for _,scope:=range []string{"ALL",cls}{x:=acc[scope][mb][ts];x.slots++;if changed{x.cross++}}
    }
   }
  }
 }
 res:=UP181BResult{Schema:UP181BMarginConditionalSchema,Experiment:"UP-181B-margin-conditioned-temporal",SourceUP180BSeal:"4373af4553fa7a86c03a705f9ee85a1e62d68da5",ContextPhase:26,MarginBands:mbands,TemporalStates:states,BandsFrozenBeforeRun:true,AdaptiveBandsUsed:false,MaintenanceTriggered:false}
 for _,scope:=range scopes{for _,mb:=range mbands{for _,ts:=range states{x:=acc[scope][mb][ts];res.Metrics=append(res.Metrics,UP181BMetric{Scope:scope,MarginBand:mb,TemporalState:ts,Slots:x.slots,Crossings:x.cross,CrossingDensity:up181bRate(x.cross,x.slots)})}}}
 return res,nil
}
