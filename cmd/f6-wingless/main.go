package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){v,e:=unitary.RunWFoundationF6();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 p:=json.NewEncoder(os.Stdout);p.SetIndent("","  ");if e=p.Encode(v);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}
