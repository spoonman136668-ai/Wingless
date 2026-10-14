package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){
 r,err:=unitary.RunWFoundationF1B()
 if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
 d:=json.NewEncoder(os.Stdout);d.SetIndent("","  ")
 if err=d.Encode(r);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
