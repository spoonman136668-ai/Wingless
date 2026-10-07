package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/spoonman136668-ai/Wingless/resources"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r52read(root,name string)[]byte{b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b}
func r52err(err error)string{if err==nil{return ""};return err.Error()}

func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	a:=flag.String("twenty-fourth-root","","");b:=flag.String("twenty-fifth-root","","");c:=flag.String("twenty-sixth-root","","");resourceOut:=flag.String("resource-out","","");flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,a,b,c}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}
	host:=&resources.Host{Path:"."};beforeHost,beforeHostErr:=host.Snapshot();probe,probeErr:=resources.NewProcessProbe(os.Getpid());var beforeProc resources.ProcessMetrics;var beforeProcErr error;if probeErr==nil{beforeProc,beforeProcErr=probe.Snapshot()}
	started:=time.Now()
	result:=unitary.RunWlmLmExternalFutureDataLearnedStateFormerRepresentationReadoutCouplingDiagnosisR52(
		r52read(*transfer,"code.bin"),r52read(*transfer,"structured.bin"),r52read(*transfer,"technical-prose.bin"),
		r52read(*third,"code.bin"),r52read(*third,"structured.bin"),r52read(*third,"technical-prose.bin"),
		r52read(*fourth,"code.bin"),r52read(*fourth,"structured.bin"),r52read(*fourth,"technical-prose.bin"),
		r52read(*fifth,"code.bin"),r52read(*fifth,"structured.bin"),r52read(*fifth,"technical-prose.bin"),
		r52read(*a,"code.bin"),r52read(*a,"structured.bin"),r52read(*a,"technical-prose.bin"),
		r52read(*b,"code.bin"),r52read(*b,"structured.bin"),r52read(*b,"technical-prose.bin"),
		r52read(*c,"code.bin"),r52read(*c,"structured.bin"),r52read(*c,"technical-prose.bin"))
	afterHost,afterHostErr:=host.Snapshot();var afterProc resources.ProcessMetrics;var afterProcErr error;if probeErr==nil{afterProc,afterProcErr=probe.Snapshot()};var governorErr error;if afterHostErr==nil{governorErr=resources.Check(resources.Policy{MaxCPU:100},afterHost)}
	if *resourceOut!=""{env:=map[string]interface{}{"schema":"wingless.rsi-resource-telemetry.v1","experiment":"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-REPRESENTATION-READOUT-COUPLING-DIAGNOSIS-R52","hardware_identity_assumed":false,"algorithmic_envelope":map[string]interface{}{"coupling_arms":4,"state_dimension":6,"max_readout_parameters_including_intercept":7,"total_adaptation_budget":1744,"external_model_calls":0},"elapsed_ms":time.Since(started).Milliseconds(),"host_before":beforeHost,"host_after":afterHost,"host_before_error":r52err(beforeHostErr),"host_after_error":r52err(afterHostErr),"process_before":beforeProc,"process_after":afterProc,"process_probe_error":r52err(probeErr),"process_before_error":r52err(beforeProcErr),"process_after_error":r52err(afterProcErr),"governor_error":r52err(governorErr)};raw,e:=json.Marshal(env);if e!=nil{panic(e)};if e=os.WriteFile(*resourceOut,raw,0600);e!=nil{panic(e)}}
	if e:=json.NewEncoder(os.Stdout).Encode(result);e!=nil{panic(e)}
}
