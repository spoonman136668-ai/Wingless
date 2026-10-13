package main

import (
	"encoding/json"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmSiRawRepPhaseAnchorFinalContextCausalityR11()); err != nil {
		panic(err)
	}
}
