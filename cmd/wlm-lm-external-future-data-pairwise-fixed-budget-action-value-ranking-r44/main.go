package main

import(
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"github.com/spoonman136668-ai/Wingless/unitary"
)

func r44read(root,name string)[]byte{
	b,e:=os.ReadFile(filepath.Join(root,name));if e!=nil{panic(e)};return b
}
func main(){
	transfer:=flag.String("transfer-root","","");third:=flag.String("third-root","","");fourth:=flag.String("fourth-root","","");fifth:=flag.String("fifth-root","","")
	sixth:=flag.String("sixth-root","","");seventh:=flag.String("seventh-root","","");eighth:=flag.String("eighth-root","","");flag.Parse()
	for _,p:=range []*string{transfer,third,fourth,fifth,sixth,seventh,eighth}{if *p==""{panic("MANIFEST_ROOT_REQUIRED")}}
	r:=unitary.RunWlmLmExternalFutureDataPairwiseFixedBudgetActionValueRankingR44(
		r44read(*transfer,"code.bin"),r44read(*transfer,"structured.bin"),r44read(*transfer,"technical-prose.bin"),
		r44read(*third,"code.bin"),r44read(*third,"structured.bin"),r44read(*third,"technical-prose.bin"),
		r44read(*fourth,"code.bin"),r44read(*fourth,"structured.bin"),r44read(*fourth,"technical-prose.bin"),
		r44read(*fifth,"code.bin"),r44read(*fifth,"structured.bin"),r44read(*fifth,"technical-prose.bin"),
		r44read(*sixth,"code.bin"),r44read(*sixth,"structured.bin"),r44read(*sixth,"technical-prose.bin"),
		r44read(*seventh,"code.bin"),r44read(*seventh,"structured.bin"),r44read(*seventh,"technical-prose.bin"),
		r44read(*eighth,"code.bin"),r44read(*eighth,"structured.bin"),r44read(*eighth,"technical-prose.bin"),
	)
	if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}
}
