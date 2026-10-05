package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func read(root,name string)[]byte{
	b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b
}
func main(){
	mode:=flag.String("mode","r44","")
	tr:=flag.String("transfer-root","","");th:=flag.String("third-root","","");fo:=flag.String("fourth-root","","");fi:=flag.String("fifth-root","","")
	e1:=flag.String("eval1-root","","");e2:=flag.String("eval2-root","","");e3:=flag.String("eval3-root","","")
	flag.Parse()
	for _,v:=range []*string{tr,th,fo,fi,e1,e2,e3}{if *v==""{panic("MANIFEST_ROOT_REQUIRED")}}
	r:=unitary.RunWlmLmExternalFutureDataPairwiseActionValueR44(*mode,
		read(*tr,"code.bin"),read(*tr,"structured.bin"),read(*tr,"technical-prose.bin"),
		read(*th,"code.bin"),read(*th,"structured.bin"),read(*th,"technical-prose.bin"),
		read(*fo,"code.bin"),read(*fo,"structured.bin"),read(*fo,"technical-prose.bin"),
		read(*fi,"code.bin"),read(*fi,"structured.bin"),read(*fi,"technical-prose.bin"),
		read(*e1,"code.bin"),read(*e1,"structured.bin"),read(*e1,"technical-prose.bin"),
		read(*e2,"code.bin"),read(*e2,"structured.bin"),read(*e2,"technical-prose.bin"),
		read(*e3,"code.bin"),read(*e3,"structured.bin"),read(*e3,"technical-prose.bin"))
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
