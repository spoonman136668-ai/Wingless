package main

import (
	"encoding/json"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/unitary"
	"os"
)

func main() {
	b, err := json.Marshal(unitary.RunWlmSiRawRepPhaseAnchorFinalRecordLocalityR12())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
