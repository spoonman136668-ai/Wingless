package main

import (
	"encoding/json"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	result, _ := unitary.RunUPLM3D()
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(result)
}
