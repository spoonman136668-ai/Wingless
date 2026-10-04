package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r28read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCalibrationConfidenceAttributionR28(
		r28read(*transfer,"code.bin"),r28read(*transfer,"structured.bin"),r28read(*transfer,"technical-prose.bin"),
		r28read(*third,"code.bin"),r28read(*third,"structured.bin"),r28read(*third,"technical-prose.bin"),
		r28read(*fourth,"code.bin"),r28read(*fourth,"structured.bin"),r28read(*fourth,"technical-prose.bin"),
		r28read(*fifth,"code.bin"),r28read(*fifth,"structured.bin"),r28read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
