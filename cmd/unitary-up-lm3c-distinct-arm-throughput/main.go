package main

import (
	"encoding/json"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	result, err := unitary.RunUPLM3C()
	if err != nil {
		panic(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		panic(err)
	}
}
