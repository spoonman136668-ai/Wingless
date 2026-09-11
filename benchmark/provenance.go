package benchmark

import (
	"fmt"
	"github.com/spoonman136668-ai/Wingless/resources"
)

type Provenance struct {
	RuntimeVersion *string        `json:"runtime_version"`
	RuntimeHash    *string        `json:"runtime_hash"`
	ModelRepo      *string        `json:"model_repo"`
	ModelRevision  *string        `json:"model_revision"`
	ModelHash      *string        `json:"model_hash"`
	Quantization   *string        `json:"quantization"`
	ModelBytes     *int64         `json:"model_bytes"`
	Context        *int           `json:"context_tokens"`
	Threads        *int           `json:"threads"`
	GPULayers      *int           `json:"gpu_layers"`
	Batch          *int           `json:"batch"`
	MicroBatch     *int           `json:"micro_batch"`
	CPUModel       *string        `json:"cpu_model"`
	GPU            *resources.GPU `json:"gpu"`
}

func Profile(name string, gpu bool) error {
	switch name {
	case "CPU_ONLY":
		if gpu {
			return fmt.Errorf("CPU_ONLY forbids GPU layers")
		}
	case "GPU_RESIDENT_SMALL":
		if !gpu {
			return fmt.Errorf("GPU_RESIDENT_SMALL requires GPU configuration")
		}
	default:
		return fmt.Errorf("profile %s unavailable", name)
	}
	return nil
}

type processObserver interface {
	ProcessSnapshot() (*resources.ProcessMetrics, error)
}
