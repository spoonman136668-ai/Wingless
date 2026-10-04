package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r29read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCalibrationOutcomeMechanismAttributionR29(
		r29read(*transfer,"code.bin"),r29read(*transfer,"structured.bin"),r29read(*transfer,"technical-prose.bin"),
		r29read(*third,"code.bin"),r29read(*third,"structured.bin"),r29read(*third,"technical-prose.bin"),
		r29read(*fourth,"code.bin"),r29read(*fourth,"structured.bin"),r29read(*fourth,"technical-prose.bin"),
		r29read(*fifth,"code.bin"),r29read(*fifth,"structured.bin"),r29read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
