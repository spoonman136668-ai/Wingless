package main
import("encoding/json";"flag";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func readR19(p string)[]byte{b,e:=os.ReadFile(p);if e!=nil{panic(e)};return b}
func main(){c:=flag.String("code","","");s:=flag.String("structured","","");p:=flag.String("technical-prose","","");flag.Parse();if *c==""||*s==""||*p==""{panic("SOURCE_PATH_REQUIRED")};if e:=json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmExternalFutureDataMarginalSignalRankCalibrationR19(readR19(*c),readR19(*s),readR19(*p)));e!=nil{panic(e)}}
