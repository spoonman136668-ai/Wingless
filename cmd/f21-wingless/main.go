package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){r,e:=unitary.RunWFoundationF21();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 x:=json.NewEncoder(os.Stdout);x.SetIndent("","  ");if e=x.Encode(r);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}