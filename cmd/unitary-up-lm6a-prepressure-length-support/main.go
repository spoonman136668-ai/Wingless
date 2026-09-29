package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){r,err:=unitary.RunUPLM6A();if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)};e:=json.NewEncoder(os.Stdout);e.SetIndent("","  ");if err:=e.Encode(r);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}}
