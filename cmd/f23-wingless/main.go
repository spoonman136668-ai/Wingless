package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){v,e:=unitary.RunWFoundationF23();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 x:=json.NewEncoder(os.Stdout);x.SetIndent("","  ");if e=x.Encode(v);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}