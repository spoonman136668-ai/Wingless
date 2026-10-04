package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r39read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCalibrationInteractionNonlinearityAttributionR39(
		r39read(*transfer,"code.bin"),r39read(*transfer,"structured.bin"),r39read(*transfer,"technical-prose.bin"),
		r39read(*third,"code.bin"),r39read(*third,"structured.bin"),r39read(*third,"technical-prose.bin"),
		r39read(*fourth,"code.bin"),r39read(*fourth,"structured.bin"),r39read(*fourth,"technical-prose.bin"),
		r39read(*fifth,"code.bin"),r39read(*fifth,"structured.bin"),r39read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
