package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/spoonman136668-ai/Wingless/resources"
	"os"
)

func runGPU(args []string) error {
	f := flag.NewFlagSet("gpu", flag.ContinueOnError)
	path := f.String("nvidia-smi", "", "absolute trusted NVIDIA utility path")
	index := f.Int("index", 0, "device index")
	if e := f.Parse(args); e != nil {
		return e
	}
	g, e := (resources.NVIDIA{Executable: *path, Index: *index}).GPU(context.Background())
	if e != nil {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(g)
}
