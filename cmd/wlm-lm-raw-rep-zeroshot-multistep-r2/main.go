package main

import (
	"encoding/json"
	"github.com/spoonman136668-ai/Wingless/unitary"
	"os"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmLmRawRepZeroshotMultistepR2()); err != nil {
		panic(err)
	}
}