package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r35read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataRawRepresentationAttributionR35(
		r35read(*transfer,"code.bin"),r35read(*transfer,"structured.bin"),r35read(*transfer,"technical-prose.bin"),
		r35read(*third,"code.bin"),r35read(*third,"structured.bin"),r35read(*third,"technical-prose.bin"),
		r35read(*fourth,"code.bin"),r35read(*fourth,"structured.bin"),r35read(*fourth,"technical-prose.bin"),
		r35read(*fifth,"code.bin"),r35read(*fifth,"structured.bin"),r35read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
