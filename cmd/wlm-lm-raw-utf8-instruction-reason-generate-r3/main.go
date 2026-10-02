package main

import (
	"encoding/json"
	"github.com/spoonman136668-ai/Wingless/unitary"
	"os"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmRawUTF8InstructionReasonGenerateR3()); err != nil {
		panic(err)
	}
}
