package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r43read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataStateConditionalResponseAttributionR43(
		r43read(*transfer,"code.bin"),r43read(*transfer,"structured.bin"),r43read(*transfer,"technical-prose.bin"),
		r43read(*third,"code.bin"),r43read(*third,"structured.bin"),r43read(*third,"technical-prose.bin"),
		r43read(*fourth,"code.bin"),r43read(*fourth,"structured.bin"),r43read(*fourth,"technical-prose.bin"),
		r43read(*fifth,"code.bin"),r43read(*fifth,"structured.bin"),r43read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
