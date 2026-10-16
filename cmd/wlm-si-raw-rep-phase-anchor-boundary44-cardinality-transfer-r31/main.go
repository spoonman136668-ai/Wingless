package main

import (
	"encoding/json"
	"github.com/spoonman136668-ai/Wingless/unitary"
	"os"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmSiRawRepPhaseAnchorBoundary44CardinalityTransferR31()); err != nil {
		panic(err)
	}
}
