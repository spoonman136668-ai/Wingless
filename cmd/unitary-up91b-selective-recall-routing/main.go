package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	r, err := unitary.RunUP91B()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
