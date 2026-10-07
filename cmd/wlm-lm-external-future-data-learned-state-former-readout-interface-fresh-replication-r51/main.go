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

func r51read(root,name string)[]byte{b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b}
func r51err(err error)string{if err==nil{return ""};return err.Error()}

func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	a:=flag.String("twenty-fourth-root","","");b:=flag.String("twenty-fifth-root","","");c:=flag.String("twenty-sixth-root","","");resourceOut:=flag.String("resource-out","","");flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,a,b,c}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}
	host:=&resources.Host{Path:"."};beforeHost,beforeHostErr:=host.Snapshot();probe,probeErr:=resources.NewProcessProbe(os.Getpid());var beforeProc resources.ProcessMetrics;var beforeProcErr error;if probeErr==nil{beforeProc,beforeProcErr=probe.Snapshot()}
	started:=time.Now()
	result:=unitary.RunWlmLmExternalFutureDataLearnedStateFormerReadoutInterfaceFreshReplicationR51(
		r51read(*transfer,"code.bin"),r51read(*transfer,"structured.bin"),r51read(*transfer,"technical-prose.bin"),
		r51read(*third,"code.bin"),r51read(*third,"structured.bin"),r51read(*third,"technical-prose.bin"),
		r51read(*fourth,"code.bin"),r51read(*fourth,"structured.bin"),r51read(*fourth,"technical-prose.bin"),
		r51read(*fifth,"code.bin"),r51read(*fifth,"structured.bin"),r51read(*fifth,"technical-prose.bin"),
		r51read(*a,"code.bin"),r51read(*a,"structured.bin"),r51read(*a,"technical-prose.bin"),
		r51read(*b,"code.bin"),r51read(*b,"structured.bin"),r51read(*b,"technical-prose.bin"),
		r51read(*c,"code.bin"),r51read(*c,"structured.bin"),r51read(*c,"technical-prose.bin"))
	afterHost,afterHostErr:=host.Snapshot();var afterProc resources.ProcessMetrics;var afterProcErr error;if probeErr==nil{afterProc,afterProcErr=probe.Snapshot()};var governorErr error;if afterHostErr==nil{governorErr=resources.Check(resources.Policy{MaxCPU:100},afterHost)}
	if *resourceOut!=""{env:=map[string]interface{}{"schema":"wingless.rsi-resource-telemetry.v1","experiment":"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-READOUT-INTERFACE-FRESH-REPLICATION-R51","hardware_identity_assumed":false,"algorithmic_envelope":map[string]interface{}{"readout_arms":2,"state_dimension":6,"max_readout_parameters_including_intercept":7,"total_adaptation_budget":1744,"external_model_calls":0},"elapsed_ms":time.Since(started).Milliseconds(),"host_before":beforeHost,"host_after":afterHost,"host_before_error":r51err(beforeHostErr),"host_after_error":r51err(afterHostErr),"process_before":beforeProc,"process_after":afterProc,"process_probe_error":r51err(probeErr),"process_before_error":r51err(beforeProcErr),"process_after_error":r51err(afterProcErr),"governor_error":r51err(governorErr)};raw,e:=json.Marshal(env);if e!=nil{panic(e)};if e=os.WriteFile(*resourceOut,raw,0600);e!=nil{panic(e)}}
	if e:=json.NewEncoder(os.Stdout).Encode(result);e!=nil{panic(e)}
}
