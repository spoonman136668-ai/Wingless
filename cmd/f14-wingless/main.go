package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){r,e:=unitary.RunWFoundationF14();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 v:=json.NewEncoder(os.Stdout);v.SetIndent("","  ");if e=v.Encode(r);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}
