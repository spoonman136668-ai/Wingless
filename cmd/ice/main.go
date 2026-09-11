package main

import (
	"encoding/json"
	"fmt"
	"github.com/spoonman136668-ai/Wingless/ice"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: ice build|validate|query [symbol]")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "relationships":
		var report ice.TypeReport
		report, err = ice.Relationships(".")
		if err == nil {
			err = json.NewEncoder(os.Stdout).Encode(report)
		}
	case "build":
		err = ice.Save(".")
	case "validate":
		err = ice.Validate(".")
	case "seams", "decisions":
		err = ice.Validate(".")
		if err == nil {
			name := "decisions.json"
			if os.Args[1] == "seams" {
				name = "integration-seams.json"
			}
			var b []byte
			b, err = os.ReadFile(filepath.Join(".ice", name))
			if err == nil {
				_, err = os.Stdout.Write(b)
			}
		}
	case "query":
		if len(os.Args) != 3 {
			err = fmt.Errorf("query requires term")
		} else {
			var x []ice.Symbol
			x, err = ice.Query(".", os.Args[2])
			if err == nil {
				err = json.NewEncoder(os.Stdout).Encode(x)
			}
		}
	default:
		err = fmt.Errorf("unknown command")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
