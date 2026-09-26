package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main(){
	result,err:=unitary.RunUP50A()
	if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
	enc:=json.NewEncoder(os.Stdout);enc.SetIndent("","  ")
	if err:=enc.Encode(result);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
