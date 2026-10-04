package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataProspectiveCalibrationR27(
		read(*transfer,"code.bin"),read(*transfer,"structured.bin"),read(*transfer,"technical-prose.bin"),
		read(*third,"code.bin"),read(*third,"structured.bin"),read(*third,"technical-prose.bin"),
		read(*fourth,"code.bin"),read(*fourth,"structured.bin"),read(*fourth,"technical-prose.bin"),
		read(*fifth,"code.bin"),read(*fifth,"structured.bin"),read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
