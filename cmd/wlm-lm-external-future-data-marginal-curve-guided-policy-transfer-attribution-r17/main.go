package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func readR17(p string) []byte {
	b, e := os.ReadFile(p)
	if e != nil {
		panic(e)
	}
	return b
}

func main() {
	c := flag.String("code", "", "")
	s := flag.String("structured", "", "")
	p := flag.String("technical-prose", "", "")
	flag.Parse()
	if *c == "" || *s == "" || *p == "" {
		panic("SOURCE_PATH_REQUIRED")
	}
	if e := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmExternalFutureDataMarginalCurveGuidedPolicyTransferAttributionR17(readR17(*c), readR17(*s), readR17(*p))); e != nil {
		panic(e)
	}
}
