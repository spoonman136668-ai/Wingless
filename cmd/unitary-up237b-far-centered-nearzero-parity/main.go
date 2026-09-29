package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){r,e:=unitary.RunUP237B();if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)};j:=json.NewEncoder(os.Stdout);j.SetIndent("","  ");if e=j.Encode(r);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}}
