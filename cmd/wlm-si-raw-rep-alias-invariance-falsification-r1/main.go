package main

import (
	"encoding/json"
	"github.com/spoonman136668-ai/Wingless/unitary"
	"os"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmSiRawRepAliasInvarianceFalsificationR1()); err != nil {
		panic(err)
	}
}
