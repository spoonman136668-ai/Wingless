package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r36read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataCrossManifestRepresentationLocalityAttributionR36(
		r36read(*transfer,"code.bin"),r36read(*transfer,"structured.bin"),r36read(*transfer,"technical-prose.bin"),
		r36read(*third,"code.bin"),r36read(*third,"structured.bin"),r36read(*third,"technical-prose.bin"),
		r36read(*fourth,"code.bin"),r36read(*fourth,"structured.bin"),r36read(*fourth,"technical-prose.bin"),
		r36read(*fifth,"code.bin"),r36read(*fifth,"structured.bin"),r36read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
