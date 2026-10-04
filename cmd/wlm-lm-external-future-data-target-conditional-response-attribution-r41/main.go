package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r41read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataTargetConditionalResponseAttributionR41(
		r41read(*transfer,"code.bin"),r41read(*transfer,"structured.bin"),r41read(*transfer,"technical-prose.bin"),
		r41read(*third,"code.bin"),r41read(*third,"structured.bin"),r41read(*third,"technical-prose.bin"),
		r41read(*fourth,"code.bin"),r41read(*fourth,"structured.bin"),r41read(*fourth,"technical-prose.bin"),
		r41read(*fifth,"code.bin"),r41read(*fifth,"structured.bin"),r41read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
