package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){r,e:=unitary.RunWFoundationF10();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 a:=json.NewEncoder(os.Stdout);a.SetIndent("","  ");if e=a.Encode(r);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}
