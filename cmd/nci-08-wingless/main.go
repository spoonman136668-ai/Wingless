package main
import("encoding/json";"flag";"fmt";"os";"github.com/spoonman136668-ai/Wingless/unitary")
func main(){
 phase:=flag.String("phase","development","development or sealed");flag.Parse()
 r,err:=unitary.RunNCI08Wingless(*phase)
 if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
 enc:=json.NewEncoder(os.Stdout);enc.SetIndent("","  ")
 if err:=enc.Encode(r);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}