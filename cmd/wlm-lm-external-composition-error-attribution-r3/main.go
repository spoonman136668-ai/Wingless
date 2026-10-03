package main
import("encoding/json";"flag";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func readAttr(path string)[]byte{b,e:=os.ReadFile(path);if e!=nil{panic(e)};return b}
func main(){c:=flag.String("code","","");s:=flag.String("structured","","");p:=flag.String("technical-prose","","");flag.Parse();if *c==""||*s==""||*p==""{panic("SOURCE_PATH_REQUIRED")};if e:=json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmExternalCompositionErrorAttributionR3(readAttr(*c),readAttr(*s),readAttr(*p)));e!=nil{panic(e)}}
