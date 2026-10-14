package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){
 d,err:=unitary.RunNCI11R3Wingless()
 if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
 o:=json.NewEncoder(os.Stdout);o.SetIndent("","  ")
 if err=o.Encode(d);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
