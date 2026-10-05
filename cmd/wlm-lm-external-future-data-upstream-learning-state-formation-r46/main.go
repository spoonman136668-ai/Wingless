package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r46read(root,name string)[]byte{
	b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b
}
func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	twelfth:=flag.String("twelfth-root","","");thirteenth:=flag.String("thirteenth-root","","");fourteenth:=flag.String("fourteenth-root","","");flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,twelfth,thirteenth,fourteenth}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}
	r:=unitary.RunWlmLmExternalFutureDataUpstreamLearningStateFormationR46(
		r46read(*transfer,"code.bin"),r46read(*transfer,"structured.bin"),r46read(*transfer,"technical-prose.bin"),
		r46read(*third,"code.bin"),r46read(*third,"structured.bin"),r46read(*third,"technical-prose.bin"),
		r46read(*fourth,"code.bin"),r46read(*fourth,"structured.bin"),r46read(*fourth,"technical-prose.bin"),
		r46read(*fifth,"code.bin"),r46read(*fifth,"structured.bin"),r46read(*fifth,"technical-prose.bin"),
		r46read(*twelfth,"code.bin"),r46read(*twelfth,"structured.bin"),r46read(*twelfth,"technical-prose.bin"),
		r46read(*thirteenth,"code.bin"),r46read(*thirteenth,"structured.bin"),r46read(*thirteenth,"technical-prose.bin"),
		r46read(*fourteenth,"code.bin"),r46read(*fourteenth,"structured.bin"),r46read(*fourteenth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
