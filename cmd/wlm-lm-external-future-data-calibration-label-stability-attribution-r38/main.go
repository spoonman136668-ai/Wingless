package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r38read(root,name string)[]byte{
	b,e:=os.ReadFile(filepath.Join(root,name))
	if e!=nil{panic(e)}
	return b
}

func main(){
	transfer:=flag.String("transfer-root","","")
	third:=flag.String("third-root","","")
	fourth:=flag.String("fourth-root","","")
	fifth:=flag.String("fifth-root","","")
	flag.Parse()
	if *transfer==""||*third==""||*fourth==""||*fifth==""{panic("MANIFEST_ROOT_REQUIRED")}
	r:=unitary.RunWlmLmExternalFutureDataCalibrationLabelStabilityAttributionR38(
		r38read(*transfer,"code.bin"),r38read(*transfer,"structured.bin"),r38read(*transfer,"technical-prose.bin"),
		r38read(*third,"code.bin"),r38read(*third,"structured.bin"),r38read(*third,"technical-prose.bin"),
		r38read(*fourth,"code.bin"),r38read(*fourth,"structured.bin"),r38read(*fourth,"technical-prose.bin"),
		r38read(*fifth,"code.bin"),r38read(*fifth,"structured.bin"),r38read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
