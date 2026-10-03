package main
import("encoding/json";"flag";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func readDist(p string)[]byte{b,e:=os.ReadFile(p);if e!=nil{panic(e)};return b}
func main(){c:=flag.String("code","","");s:=flag.String("structured","","");p:=flag.String("technical-prose","","");flag.Parse();if *c==""||*s==""||*p==""{panic("SOURCE_PATH_REQUIRED")};r:=unitary.RunWlmLmExternalMotifRelationDistributionalCompositionR1(readDist(*c),readDist(*s),readDist(*p));if e:=json.NewEncoder(os.Stdout).Encode(r);e!=nil{panic(e)}}
