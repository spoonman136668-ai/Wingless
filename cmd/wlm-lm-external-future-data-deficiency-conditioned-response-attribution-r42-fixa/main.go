package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r42fixaread(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataDeficiencyConditionalResponseAttributionR42FixA(
		r42fixaread(*transfer,"code.bin"),r42fixaread(*transfer,"structured.bin"),r42fixaread(*transfer,"technical-prose.bin"),
		r42fixaread(*third,"code.bin"),r42fixaread(*third,"structured.bin"),r42fixaread(*third,"technical-prose.bin"),
		r42fixaread(*fourth,"code.bin"),r42fixaread(*fourth,"structured.bin"),r42fixaread(*fourth,"technical-prose.bin"),
		r42fixaread(*fifth,"code.bin"),r42fixaread(*fifth,"structured.bin"),r42fixaread(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
