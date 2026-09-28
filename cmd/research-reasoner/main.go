package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spoonman136668-ai/Wingless/reasoner"
)

func main() {
	requestPath := flag.String("request", "", "path to a wingless.experiment-reasoning-request.v1 JSON file")
	outputPath := flag.String("out", "", "path for the sealed reasoning result JSON")
	flag.Parse()

	if os.Getenv("WINGLESS_REMOTE_REASONER_ENABLE") != "1" {
		die(errors.New("remote reasoner disabled; explicit WINGLESS_REMOTE_REASONER_ENABLE=1 required"))
	}
	if *requestPath == "" || *outputPath == "" {
		die(errors.New("both -request and -out are required"))
	}
	key := os.Getenv("WINGLESS_REASONER_OPENROUTER_API_KEY")
	if key == "" {
		die(errors.New("WINGLESS_REASONER_OPENROUTER_API_KEY is not set"))
	}

	raw, err := os.ReadFile(*requestPath)
	if err != nil {
		die(err)
	}
	var request reasoner.Request
	if err := json.Unmarshal(raw, &request); err != nil {
		die(err)
	}
	if err := request.Validate(); err != nil {
		die(err)
	}

	cfg := reasoner.DefaultConfig()
	if model := os.Getenv("WINGLESS_REASONER_MODEL"); model != "" {
		cfg.Model = model
	}
	client, err := reasoner.NewClient(key, cfg)
	if err != nil {
		die(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(request.TimeoutSeconds)*time.Second)
	defer cancel()
	result, err := client.Invoke(ctx, request)
	if err != nil {
		die(err)
	}
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		die(err)
	}
	encoded = append(encoded, '\n')
	if err := writeAtomic(*outputPath, encoded); err != nil {
		die(err)
	}
}

func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".reasoner-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "research-reasoner:", err)
	os.Exit(1)
}
