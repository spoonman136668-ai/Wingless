package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r31read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataResidualCalibrationAttributionR31(
		r31read(*transfer,"code.bin"),r31read(*transfer,"structured.bin"),r31read(*transfer,"technical-prose.bin"),
		r31read(*third,"code.bin"),r31read(*third,"structured.bin"),r31read(*third,"technical-prose.bin"),
		r31read(*fourth,"code.bin"),r31read(*fourth,"structured.bin"),r31read(*fourth,"technical-prose.bin"),
		r31read(*fifth,"code.bin"),r31read(*fifth,"structured.bin"),r31read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
