package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r32read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataResidualMagnitudeCalibrationAttributionR32(
		r32read(*transfer,"code.bin"),r32read(*transfer,"structured.bin"),r32read(*transfer,"technical-prose.bin"),
		r32read(*third,"code.bin"),r32read(*third,"structured.bin"),r32read(*third,"technical-prose.bin"),
		r32read(*fourth,"code.bin"),r32read(*fourth,"structured.bin"),r32read(*fourth,"technical-prose.bin"),
		r32read(*fifth,"code.bin"),r32read(*fifth,"structured.bin"),r32read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
