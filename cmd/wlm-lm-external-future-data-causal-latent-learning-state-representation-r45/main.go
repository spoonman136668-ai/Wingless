package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r45read(root,name string)[]byte{
	b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b
}
func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	ninth:=flag.String("ninth-root","","");tenth:=flag.String("tenth-root","","");eleventh:=flag.String("eleventh-root","","");flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,ninth,tenth,eleventh}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}
	r:=unitary.RunWlmLmExternalFutureDataCausalLatentLearningStateRepresentationR45(
		r45read(*transfer,"code.bin"),r45read(*transfer,"structured.bin"),r45read(*transfer,"technical-prose.bin"),
		r45read(*third,"code.bin"),r45read(*third,"structured.bin"),r45read(*third,"technical-prose.bin"),
		r45read(*fourth,"code.bin"),r45read(*fourth,"structured.bin"),r45read(*fourth,"technical-prose.bin"),
		r45read(*fifth,"code.bin"),r45read(*fifth,"structured.bin"),r45read(*fifth,"technical-prose.bin"),
		r45read(*ninth,"code.bin"),r45read(*ninth,"structured.bin"),r45read(*ninth,"technical-prose.bin"),
		r45read(*tenth,"code.bin"),r45read(*tenth,"structured.bin"),r45read(*tenth,"technical-prose.bin"),
		r45read(*eleventh,"code.bin"),r45read(*eleventh,"structured.bin"),r45read(*eleventh,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
