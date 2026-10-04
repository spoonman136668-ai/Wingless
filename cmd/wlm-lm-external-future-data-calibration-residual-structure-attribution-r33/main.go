package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r33read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCalibrationResidualStructureAttributionR33(
		r33read(*transfer,"code.bin"),r33read(*transfer,"structured.bin"),r33read(*transfer,"technical-prose.bin"),
		r33read(*third,"code.bin"),r33read(*third,"structured.bin"),r33read(*third,"technical-prose.bin"),
		r33read(*fourth,"code.bin"),r33read(*fourth,"structured.bin"),r33read(*fourth,"technical-prose.bin"),
		r33read(*fifth,"code.bin"),r33read(*fifth,"structured.bin"),r33read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
