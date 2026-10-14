package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){out,err:=unitary.RunWFoundationF17();if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
 enc:=json.NewEncoder(os.Stdout);enc.SetIndent("","  ");if err=enc.Encode(out);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}}
