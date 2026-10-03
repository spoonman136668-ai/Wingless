package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	root:=flag.String("root","","directory containing verified external source bytes")
	flag.Parse()
	if *root=="" {
		panic("EXTERNAL_ROOT_REQUIRED")
	}
	if err:=json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmExternalRawByteGeneralizationR1(*root));err!=nil {
		panic(err)
	}
}
