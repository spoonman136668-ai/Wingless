package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r40read(root,name string)[]byte{
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
	r:=unitary.RunWlmLmExternalFutureDataProspectiveApplicabilityRegionAttributionR40(
		r40read(*transfer,"code.bin"),r40read(*transfer,"structured.bin"),r40read(*transfer,"technical-prose.bin"),
		r40read(*third,"code.bin"),r40read(*third,"structured.bin"),r40read(*third,"technical-prose.bin"),
		r40read(*fourth,"code.bin"),r40read(*fourth,"structured.bin"),r40read(*fourth,"technical-prose.bin"),
		r40read(*fifth,"code.bin"),r40read(*fifth,"structured.bin"),r40read(*fifth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
