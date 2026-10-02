package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	result, err := unitary.RunUP8()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
