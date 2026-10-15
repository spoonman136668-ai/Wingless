package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){v,e:=unitary.RunWFoundationF19();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 enc:=json.NewEncoder(os.Stdout);enc.SetIndent("","  ");if e=enc.Encode(v);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}