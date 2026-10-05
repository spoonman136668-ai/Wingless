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

func r47read(root,name string)[]byte{
	b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b
}

func errText(err error)string{if err==nil{return ""};return err.Error()}

func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	fifteenth:=flag.String("fifteenth-root","","");sixteenth:=flag.String("sixteenth-root","","");seventeenth:=flag.String("seventeenth-root","","")
	resourceOut:=flag.String("resource-out","","")
	flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,fifteenth,sixteenth,seventeenth}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}

	host:=&resources.Host{Path:"."}
	beforeHost,beforeHostErr:=host.Snapshot()
	probe,probeErr:=resources.NewProcessProbe(os.Getpid())
	var beforeProc resources.ProcessMetrics
	var beforeProcErr error
	if probeErr==nil { beforeProc,beforeProcErr=probe.Snapshot() }
	started:=time.Now()

	result:=unitary.RunWlmLmExternalFutureDataLearnedRawInputFormationR47(
		r47read(*transfer,"code.bin"),r47read(*transfer,"structured.bin"),r47read(*transfer,"technical-prose.bin"),
		r47read(*third,"code.bin"),r47read(*third,"structured.bin"),r47read(*third,"technical-prose.bin"),
		r47read(*fourth,"code.bin"),r47read(*fourth,"structured.bin"),r47read(*fourth,"technical-prose.bin"),
		r47read(*fifth,"code.bin"),r47read(*fifth,"structured.bin"),r47read(*fifth,"technical-prose.bin"),
		r47read(*fifteenth,"code.bin"),r47read(*fifteenth,"structured.bin"),r47read(*fifteenth,"technical-prose.bin"),
		r47read(*sixteenth,"code.bin"),r47read(*sixteenth,"structured.bin"),r47read(*sixteenth,"technical-prose.bin"),
		r47read(*seventeenth,"code.bin"),r47read(*seventeenth,"structured.bin"),r47read(*seventeenth,"technical-prose.bin"),
	)

	afterHost,afterHostErr:=host.Snapshot()
	var afterProc resources.ProcessMetrics
	var afterProcErr error
	if probeErr==nil { afterProc,afterProcErr=probe.Snapshot() }
	var governorErr error
	if afterHostErr==nil { governorErr=resources.Check(resources.Policy{MaxCPU:100},afterHost) }

	if *resourceOut!="" {
		env:=map[string]interface{}{
			"schema":"wingless.rsi-resource-telemetry.v1",
			"experiment":"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-RAW-INPUT-FORMATION-R47",
			"hardware_identity_assumed":false,
			"consumer_hardware_policy":"no capacity growth; GPU not required; exact user hardware not invented",
			"algorithmic_envelope":map[string]interface{}{
				"candidate_mechanisms":4,
				"state_dimension":6,
				"readout_parameters_including_intercept":7,
				"max_candidate_learned_scalars":6,
				"baseline_effective_state_former_scalars":1,
				"total_adaptation_budget":1744,
				"external_model_calls":0,
			},
			"elapsed_ms":time.Since(started).Milliseconds(),
			"host_before":beforeHost,"host_after":afterHost,
			"host_before_error":errText(beforeHostErr),"host_after_error":errText(afterHostErr),
			"process_before":beforeProc,"process_after":afterProc,
			"process_probe_error":errText(probeErr),"process_before_error":errText(beforeProcErr),"process_after_error":errText(afterProcErr),
			"governor_error":errText(governorErr),
		}
		b,e:=json.Marshal(env);if e!=nil{panic(e)}
		if e=os.WriteFile(*resourceOut,b,0600);e!=nil{panic(e)}
	}
	if e:=json.NewEncoder(os.Stdout).Encode(result);e!=nil{panic(e)}
}
