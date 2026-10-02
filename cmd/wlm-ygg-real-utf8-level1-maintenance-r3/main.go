package main

import (
	"encoding/json"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func main() {
	if err := json.NewEncoder(os.Stdout).Encode(unitary.RunWlmYggRealUTF8Level1MaintenanceR3()); err != nil {
		panic(err)
	}
}
