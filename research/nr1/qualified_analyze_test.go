package nr1

import (
	"math"
	"testing"
)

func TestQualifiedTopTrafficUsesLogicalExpertUniverse(t *testing.T) {
	report, err := AnalyzeQualifiedQwen3Coder(qualifiedTrace(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Layers) != QualifiedQwen3CoderLayers {
		t.Fatalf("layers=%d want=%d", len(report.Layers), QualifiedQwen3CoderLayers)
	}
	if len(report.Layers[0].TopTraffic) == 0 {
		t.Fatal("missing per-layer top-traffic metrics")
	}

	// 5%% of the logical 128-expert layer is ceil(6.4)=7 experts. This
	// synthetic token selects eight experts once each, so the top 5%% carries
	// 7/8 of observed routing traffic. Using only the observed-expert count as
	// the denominator would incorrectly choose one expert and report 1/8.
	if got, want := report.Layers[0].TopTraffic[0].Traffic, 0.875; got != want {
		t.Fatalf("layer top-5%% traffic=%v want=%v", got, want)
	}

	if len(report.GlobalTopTraffic) == 0 {
		t.Fatal("missing global top-traffic metrics")
	}
	// Global logical universe is 48*128=6144 layer/expert keys. Five percent
	// rounds up to 308 keys. qualifiedTrace(1) observes 48*8=384 keys once.
	if got, want := report.GlobalTopTraffic[0].Traffic, float64(308)/384.0; math.Abs(got-want) > 1e-12 {
		t.Fatalf("global top-5%% traffic=%v want=%v", got, want)
	}
}
