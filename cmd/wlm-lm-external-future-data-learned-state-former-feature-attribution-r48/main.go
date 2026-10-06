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

func r48read(root,name string)[]byte{b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b}
func r48err(err error)string{if err==nil{return ""};return err.Error()}

func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	eighteenth:=flag.String("eighteenth-root","","");nineteenth:=flag.String("nineteenth-root","","");twentieth:=flag.String("twentieth-root","","")
	resourceOut:=flag.String("resource-out","","")
	flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,eighteenth,nineteenth,twentieth}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}
	host:=&resources.Host{Path:"."};beforeHost,beforeHostErr:=host.Snapshot()
	probe,probeErr:=resources.NewProcessProbe(os.Getpid());var beforeProc resources.ProcessMetrics;var beforeProcErr error
	if probeErr==nil { beforeProc,beforeProcErr=probe.Snapshot() }
	started:=time.Now()
	result:=unitary.RunWlmLmExternalFutureDataLearnedStateFormerFeatureAttributionR48(
		r48read(*transfer,"code.bin"),r48read(*transfer,"structured.bin"),r48read(*transfer,"technical-prose.bin"),
		r48read(*third,"code.bin"),r48read(*third,"structured.bin"),r48read(*third,"technical-prose.bin"),
		r48read(*fourth,"code.bin"),r48read(*fourth,"structured.bin"),r48read(*fourth,"technical-prose.bin"),
		r48read(*fifth,"code.bin"),r48read(*fifth,"structured.bin"),r48read(*fifth,"technical-prose.bin"),
		r48read(*eighteenth,"code.bin"),r48read(*eighteenth,"structured.bin"),r48read(*eighteenth,"technical-prose.bin"),
		r48read(*nineteenth,"code.bin"),r48read(*nineteenth,"structured.bin"),r48read(*nineteenth,"technical-prose.bin"),
		r48read(*twentieth,"code.bin"),r48read(*twentieth,"structured.bin"),r48read(*twentieth,"technical-prose.bin"),
	)
	afterHost,afterHostErr:=host.Snapshot();var afterProc resources.ProcessMetrics;var afterProcErr error
	if probeErr==nil { afterProc,afterProcErr=probe.Snapshot() }
	var governorErr error;if afterHostErr==nil { governorErr=resources.Check(resources.Policy{MaxCPU:100},afterHost) }
	if *resourceOut!="" {
		env:=map[string]interface{}{
			"schema":"wingless.rsi-resource-telemetry.v1","experiment":"WLM-LM-EXTERNAL-FUTURE-DATA-LEARNED-STATE-FORMER-FEATURE-ATTRIBUTION-R48",
			"hardware_identity_assumed":false,"consumer_hardware_policy":"no capacity growth; GPU not required; exact user hardware not invented",
			"algorithmic_envelope":map[string]interface{}{"attribution_arms":4,"baseline_anchors":1,"state_dimension":6,"readout_parameters_including_intercept":7,"total_adaptation_budget":1744,"external_model_calls":0},
			"elapsed_ms":time.Since(started).Milliseconds(),"host_before":beforeHost,"host_after":afterHost,
			"host_before_error":r48err(beforeHostErr),"host_after_error":r48err(afterHostErr),"process_before":beforeProc,"process_after":afterProc,
			"process_probe_error":r48err(probeErr),"process_before_error":r48err(beforeProcErr),"process_after_error":r48err(afterProcErr),"governor_error":r48err(governorErr),
		}
		b,e:=json.Marshal(env);if e!=nil{panic(e)};if e=os.WriteFile(*resourceOut,b,0600);e!=nil{panic(e)}
	}
	if e:=json.NewEncoder(os.Stdout).Encode(result);e!=nil{panic(e)}
}
