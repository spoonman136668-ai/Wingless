package resources

import (
	"errors"
	"strings"
)

const ModelTelemetrySchema = "wingless.model-resource-telemetry.v1"

type TelemetryProvenanceKind string

const (
	TelemetryDirectMeasurement TelemetryProvenanceKind = "direct_measurement"
	TelemetryBackendReported   TelemetryProvenanceKind = "backend_reported"
	TelemetryRuntimeReported   TelemetryProvenanceKind = "runtime_reported"
	TelemetryEstimated         TelemetryProvenanceKind = "estimated"
)

func (k TelemetryProvenanceKind) Valid() bool {
	switch k {
	case TelemetryDirectMeasurement, TelemetryBackendReported, TelemetryRuntimeReported, TelemetryEstimated:
		return true
	default:
		return false
	}
}

type TelemetryProvenance struct {
	Kind   TelemetryProvenanceKind `json:"kind"`
	Source string                  `json:"source"`
	Detail string                  `json:"detail,omitempty"`
}

func (p TelemetryProvenance) Validate() error {
	if !p.Kind.Valid() {
		return errors.New("invalid telemetry provenance kind")
	}
	if strings.TrimSpace(p.Source) == "" {
		return errors.New("telemetry provenance source is required")
	}
	return nil
}

type ActiveParameterEstimate struct {
	Value      uint64              `json:"value"`
	Estimated  bool                `json:"estimated"`
	Method     string              `json:"method"`
	Provenance TelemetryProvenance `json:"provenance"`
}

func (e ActiveParameterEstimate) Validate() error {
	if !e.Estimated {
		return errors.New("active parameter estimate must explicitly be marked estimated")
	}
	if strings.TrimSpace(e.Method) == "" {
		return errors.New("active parameter estimate method is required")
	}
	if err := e.Provenance.Validate(); err != nil {
		return err
	}
	return nil
}

// ModelTelemetry records model-attributable resource evidence only. Host/process
// metrics remain in Metrics. Nil means unavailable; a non-nil pointer to zero means
// the runtime measured or reported zero.
type ModelTelemetry struct {
	Schema string `json:"schema"`

	LogicalModelBytes  *uint64 `json:"logical_model_bytes"`
	ResidentModelBytes *uint64 `json:"resident_model_bytes"`
	ActiveModelBytes   *uint64 `json:"active_model_bytes"`

	ModelVRAMBytes    *uint64 `json:"model_vram_bytes"`
	ModelHostRAMBytes *uint64 `json:"model_host_ram_bytes"`
	ModelStorageBytes *uint64 `json:"model_storage_bytes"`

	HostToDeviceBytes   *uint64 `json:"host_to_device_bytes"`
	DeviceToHostBytes   *uint64 `json:"device_to_host_bytes"`
	StorageToHostBytes  *uint64 `json:"storage_to_host_bytes"`
	StorageToDeviceBytes *uint64 `json:"storage_to_device_bytes"`

	WorkingSetBytesCurrent *uint64                  `json:"working_set_bytes_current"`
	WorkingSetBytesPeak    *uint64                  `json:"working_set_bytes_peak"`
	ActiveParameterEstimate *ActiveParameterEstimate `json:"active_parameter_estimate"`

	ModelCacheHits   *uint64 `json:"model_cache_hits"`
	ModelCacheMisses *uint64 `json:"model_cache_misses"`
	ModelPageFaults  *uint64 `json:"model_page_faults"`
	ModelPagesLoaded *uint64 `json:"model_pages_loaded"`
	ModelBytesLoaded *uint64 `json:"model_bytes_loaded"`

	CapacityProvenance   *TelemetryProvenance `json:"capacity_provenance"`
	ResidencyProvenance  *TelemetryProvenance `json:"residency_provenance"`
	MovementProvenance   *TelemetryProvenance `json:"movement_provenance"`
	WorkingSetProvenance *TelemetryProvenance `json:"working_set_provenance"`
	CacheProvenance      *TelemetryProvenance `json:"cache_provenance"`
}

func NewModelTelemetry() ModelTelemetry {
	return ModelTelemetry{Schema: ModelTelemetrySchema}
}

func (m ModelTelemetry) Validate() error {
	if m.Schema != ModelTelemetrySchema {
		return errors.New("unsupported model telemetry schema")
	}
	for _, p := range []*TelemetryProvenance{
		m.CapacityProvenance,
		m.ResidencyProvenance,
		m.MovementProvenance,
		m.WorkingSetProvenance,
		m.CacheProvenance,
	} {
		if p != nil {
			if err := p.Validate(); err != nil {
				return err
			}
		}
	}
	if m.ActiveParameterEstimate != nil {
		if err := m.ActiveParameterEstimate.Validate(); err != nil {
			return err
		}
	}
	return nil
}
