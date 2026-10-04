package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r34read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCalibrationFeatureRepresentationAttributionR34(
		r34read(*transfer,"code.bin"),r34read(*transfer,"structured.bin"),r34read(*transfer,"technical-prose.bin"),
		r34read(*third,"code.bin"),r34read(*third,"structured.bin"),r34read(*third,"technical-prose.bin"),
		r34read(*fourth,"code.bin"),r34read(*fourth,"structured.bin"),r34read(*fourth,"technical-prose.bin"),
		r34read(*fifth,"code.bin"),r34read(*fifth,"structured.bin"),r34read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
