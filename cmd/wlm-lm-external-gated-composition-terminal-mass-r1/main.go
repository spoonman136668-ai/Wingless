package main
import("encoding/json";"flag";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func readTerminalMass(path string)[]byte{b,e:=os.ReadFile(path);if e!=nil{panic(e)};return b}
func main(){c:=flag.String("code","","");s:=flag.String("structured","","");p:=flag.String("technical-prose","","");flag.Parse();if *c==""||*s==""||*p==""{panic("SOURCE_PATH_REQUIRED")};if e:=json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmExternalGatedCompositionTerminalMassR1(readTerminalMass(*c),readTerminalMass(*s),readTerminalMass(*p)));e!=nil{panic(e)}}
