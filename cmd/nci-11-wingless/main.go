package main
import("encoding/json";"flag";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){
 phase:=flag.String("phase","development","development or sealed")
 flag.Parse()
 result,err:=unitary.RunNCI11Wingless(*phase)
 if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
 o:=json.NewEncoder(os.Stdout);o.SetIndent("","  ")
 if err=o.Encode(result);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
