package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func readSource(path string) []byte {
	b,err:=os.ReadFile(path)
	if err!=nil { panic(err) }
	return b
}

func main() {
	c:=flag.String("code","","")
	s:=flag.String("structured","","")
	p:=flag.String("technical-prose","","")
	flag.Parse()
	if *c==""||*s==""||*p=="" { panic("SOURCE_PATH_REQUIRED") }
	if err:=json.NewEncoder(os.Stdout).Encode(
		unitary.RunWlmLmExternalReasoningMarginSourceAttributionR3(
			readSource(*c),readSource(*s),readSource(*p),
		),
	);err!=nil{panic(err)}
}
