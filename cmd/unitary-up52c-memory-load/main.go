package main

import(
 "encoding/json"
 "fmt"
 "os"
 "github.com/spoonman136668-ai/Wingless/unitary"
)

func main(){
 r,e:=unitary.RunUP52C()
 if e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
 enc:=json.NewEncoder(os.Stdout);enc.SetIndent("","  ")
 if e:=enc.Encode(r);e!=nil{fmt.Fprintln(os.Stderr,e);os.Exit(1)}
}
