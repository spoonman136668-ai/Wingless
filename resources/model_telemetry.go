package resources

import (
	"encoding/json"
	"errors"
	"strings"
)

const ModelTelemetrySchema = "wingless.model-resource-telemetry.v1"

type TelemetryOrigin string

const (
	TelemetryOriginDirect  TelemetryOrigin = "direct_measurement"
	TelemetryOriginBackend TelemetryOrigin = "backend_report"
	TelemetryOriginRuntime TelemetryOrigin = "runtime_report"
)

func (o TelemetryOrigin) Valid() bool {
	switch o {
	case TelemetryOriginDirect, TelemetryOriginBackend, TelemetryOriginRuntime:
		return true
	default:
		return false
	}
}

type TelemetryProvenance struct {
	Origin     TelemetryOrigin `json:"origin"`
	ProviderID string          `json:"provider_id"`
	Note       string          `json:"note,omitempty"`
}

func (p TelemetryProvenance) Validate() error {
	if !p.Origin.Valid() {
		return errors.New("invalid telemetry origin")
	}
	if strings.TrimSpace(p.ProviderID) == "" {
		return errors.New("telemetry provider id is required")
	}
	return nil
}

type ParameterEstimate struct {
	Count      uint64               `json:"count"`
	Basis      string               `json:"basis"`
	Provenance *TelemetryProvenance `json:"provenance"`
}

func (e ParameterEstimate) Validate() error {
	if strings.TrimSpace(e.Basis) == "" {
		return errors.New("parameter estimate basis is required")
	}
	if e.Provenance != nil {
		return e.Provenance.Validate()
	}
	return nil
}

func (e ParameterEstimate) MarshalJSON() ([]byte, error) {
	if err := e.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Count      uint64               `json:"count"`
		Estimated  bool                 `json:"estimated"`
		Basis      string               `json:"basis"`
		Provenance *TelemetryProvenance `json:"provenance"`
	}{
		Count: e.Count, Estimated: true, Basis: e.Basis, Provenance: e.Provenance,
	})
}

type ModelCapacityTelemetry struct {
	LogicalModelBytes  *uint64              `json:"logical_model_bytes"`
	ResidentModelBytes *uint64              `json:"resident_model_bytes"`
	ActiveModelBytes   *uint64              `json:"active_model_bytes"`
	Provenance         *TelemetryProvenance `json:"provenance"`
}

type ModelMemoryTierTelemetry struct {
	ModelVRAMBytes    *uint64              `json:"model_vram_bytes"`
	ModelHostRAMBytes *uint64              `json:"model_host_ram_bytes"`
	ModelStorageBytes *uint64              `json:"model_storage_bytes"`
	Provenance        *TelemetryProvenance `json:"provenance"`
}

type ModelDataMovementTelemetry struct {
	HostToDeviceBytes   *uint64              `json:"host_to_device_bytes"`
	DeviceToHostBytes   *uint64              `json:"device_to_host_bytes"`
	StorageToHostBytes  *uint64              `json:"storage_to_host_bytes"`
	StorageToDeviceBytes *uint64             `json:"storage_to_device_bytes"`
	Provenance          *TelemetryProvenance `json:"provenance"`
}

type ModelWorkingSetTelemetry struct {
	WorkingSetBytesCurrent   *uint64              `json:"working_set_bytes_current"`
	WorkingSetBytesPeak      *uint64              `json:"working_set_bytes_peak"`
	ActiveParameterEstimate  *ParameterEstimate   `json:"active_parameter_estimate"`
	Provenance               *TelemetryProvenance `json:"provenance"`
}

type ModelCacheTelemetry struct {
	ModelCacheHits   *uint64              `json:"model_cache_hits"`
	ModelCacheMisses *uint64              `json:"model_cache_misses"`
	ModelPageFaults  *uint64              `json:"model_page_faults"`
	ModelPagesLoaded *uint64              `json:"model_pages_loaded"`
	ModelBytesLoaded *uint64              `json:"model_bytes_loaded"`
	Provenance       *TelemetryProvenance `json:"provenance"`
}

type ModelTelemetry struct {
	Schema       string                      `json:"schema"`
	Capacity     *ModelCapacityTelemetry     `json:"capacity"`
	MemoryTiers  *ModelMemoryTierTelemetry   `json:"memory_tiers"`
	DataMovement *ModelDataMovementTelemetry `json:"data_movement"`
	WorkingSet   *ModelWorkingSetTelemetry   `json:"working_set"`
	Cache        *ModelCacheTelemetry        `json:"cache"`
}

func NewModelTelemetry() ModelTelemetry {
	return ModelTelemetry{Schema: ModelTelemetrySchema}
}

func (m ModelTelemetry) Validate() error {
	if m.Schema != ModelTelemetrySchema {
		return errors.New("invalid model telemetry schema")
	}
	provenance := []*TelemetryProvenance{}
	if m.Capacity != nil {
		provenance = append(provenance, m.Capacity.Provenance)
	}
	if m.MemoryTiers != nil {
		provenance = append(provenance, m.MemoryTiers.Provenance)
	}
	if m.DataMovement != nil {
		provenance = append(provenance, m.DataMovement.Provenance)
	}
	if m.WorkingSet != nil {
		provenance = append(provenance, m.WorkingSet.Provenance)
		if m.WorkingSet.ActiveParameterEstimate != nil {
			if err := m.WorkingSet.ActiveParameterEstimate.Validate(); err != nil {
				return err
			}
		}
	}
	if m.Cache != nil {
		provenance = append(provenance, m.Cache.Provenance)
	}
	for _, p := range provenance {
		if p != nil {
			if err := p.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}
