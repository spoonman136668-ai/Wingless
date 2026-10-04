package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r37read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCalibrationTargetDefinitionAttributionR37(
		r37read(*transfer,"code.bin"),r37read(*transfer,"structured.bin"),r37read(*transfer,"technical-prose.bin"),
		r37read(*third,"code.bin"),r37read(*third,"structured.bin"),r37read(*third,"technical-prose.bin"),
		r37read(*fourth,"code.bin"),r37read(*fourth,"structured.bin"),r37read(*fourth,"technical-prose.bin"),
		r37read(*fifth,"code.bin"),r37read(*fifth,"structured.bin"),r37read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
