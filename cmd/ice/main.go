package main
import("encoding/json";"fmt";"os";"github.com/spoonman136668-ai/Wingless/ice")
func main(){if len(os.Args)<2{fmt.Fprintln(os.Stderr,"usage: ice build|validate|query [symbol]");os.Exit(2)};var err error;switch os.Args[1]{case "build":err=ice.Save(".");case "validate":err=ice.Validate(".");case "query":if len(os.Args)!=3{err=fmt.Errorf("query requires term")}else{var x []ice.Symbol;x,err=ice.Query(".",os.Args[2]);if err==nil{err=json.NewEncoder(os.Stdout).Encode(x)}};default:err=fmt.Errorf("unknown command")};if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}}
