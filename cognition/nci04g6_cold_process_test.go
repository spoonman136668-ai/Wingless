package cognition

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const coldModeEnv = "NCI04L5A_COLD_MODE"
const coldRootEnv = "NCI04L5A_COLD_ROOT"
const coldSnapshotEnv = "NCI04L5A_COLD_SNAPSHOT"
const coldReportMarker = "NCI04L5A_COLD_REPORT="
const coldArtifactRootEnv = "NCI04L5A_ARTIFACT_ROOT"
const coldSourceSHAEnv = "NCI04L5A_SOURCE_SHA"

type coldReport struct {
	SnapshotSHA256 string  `json:"snapshot_sha256"`
	SourceSHA256   string  `json:"source_sha256"`
	Prediction     float64 `json:"prediction"`
	Observations   uint32  `json:"observations"`
	Sequence       uint8   `json:"sequence"`
	RejectedStale  bool    `json:"rejected_stale"`
}

// This helper runs ONLY in a separately launched test process. The parent
// process never holds the resumed learner. Fixed-generated-source actions are
// development evidence, NOT independently held-out scientific first-use.
func TestNCI04G6ColdProcessHelper(t *testing.T) {
	mode := os.Getenv(coldModeEnv)
	if mode == "" {
		return
	}
	root, stateFile := os.Getenv(coldRootEnv), os.Getenv(coldSnapshotEnv)
	if root == "" || stateFile == "" {
		t.Fatal("missing cold-process custody")
	}

	var probe *NativeProbe
	var err error
	if mode == "full" || mode == "prefix" {
		probe, err = NewNativeProbe(root)
	} else {
		probe, err = OpenNativeProbe(root)
	}
	if err != nil {
		t.Fatal(err)
	}
	session := &ExecutedLearningSession{Probe: probe}
	if mode == "resume" || mode == "stale" {
		raw, err := os.ReadFile(stateFile)
		if err != nil {
			t.Fatal(err)
		}
		err = session.Restore(raw)
		if mode == "stale" {
			if err == nil {
				t.Fatal("stale checkpoint was accepted after continued execution")
			}
			if probe.sequence != 2 {
				t.Fatalf("failed restore mutated durable cursor: %d", probe.sequence)
			}
			b, err := probe.Snapshot()
			if err != nil || len(b) == 0 {
				t.Fatal("rejected stale restore damaged custody", err)
			}
			bReport, _ := json.Marshal(coldReport{RejectedStale: true, Sequence: probe.sequence})
			fmt.Println(coldReportMarker + string(bReport))
			return
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	if mode == "full" || mode == "prefix" {
		if _, err := session.ActAndLearn(context.Background(), false); err != nil {
			t.Fatal(err)
		}
	}
	if mode == "full" || mode == "resume" {
		if _, err := session.ActAndLearn(context.Background(), true); err != nil {
			t.Fatal(err)
		}
	}
	if mode != "full" && mode != "prefix" && mode != "resume" {
		t.Fatal(errors.New("unauthorized cold-process mode"))
	}
	snapshot, err := session.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if mode == "prefix" {
		if err := os.WriteFile(stateFile, snapshot, 0600); err != nil {
			t.Fatal(err)
		}
	}
	p, err := session.PredictExecutable(true)
	if err != nil {
		t.Fatal(err)
	}
	source := session.Probe
	b, err := probeRead(root)
	if err != nil {
		t.Fatal(err)
	}
	if mode == "full" || mode == "resume" {
		if err := os.WriteFile(filepath.Join(root, "final.snapshot"), snapshot, 0600); err != nil {
			t.Fatal(err)
		}
	}
	digest := sha256.Sum256(snapshot)
	report := coldReport{SnapshotSHA256: hex.EncodeToString(digest[:]), SourceSHA256: probeHash(b), Prediction: p, Observations: session.State.Observations, Sequence: source.sequence}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(coldReportMarker + string(encoded))
}

// L5A engineering acceptance: the same fixed-source sequence is executed in
// one uninterrupted process, then in separate prefix and continuation
// processes using a durable, independently reopened checkpoint. This proves
// process restart for this exact generated fixture only.
func TestNCI04G6ColdProcessRestart(t *testing.T) {
	if os.Getenv(coldModeEnv) != "" {
		return
	}
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	if supplied := os.Getenv(coldArtifactRootEnv); supplied != "" {
		if !filepath.IsAbs(supplied) { t.Fatal("artifact root must be absolute") }
		if err := os.Mkdir(supplied, 0700); err != nil { t.Fatal("artifact root must be newly created", err) }
		tmp = supplied
	}
	snap := filepath.Join(tmp, "prefix.snapshot")
	fullRoot := filepath.Join(tmp, "full")
	resumedRoot := filepath.Join(tmp, "resumed")
	run := func(mode, root string) coldReport {
		t.Helper()
		command := exec.Command(exe, "-test.run=^TestNCI04G6ColdProcessHelper$", "-test.v")
		command.Env = append(os.Environ(), coldModeEnv+"="+mode, coldRootEnv+"="+root, coldSnapshotEnv+"="+snap)
		result, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("separate process %s failed: %v\n%s", mode, err, result)
		}
		for _, line := range strings.Split(string(result), "\n") {
			if index := strings.Index(line, coldReportMarker); index >= 0 {
				var report coldReport
				if err := json.Unmarshal([]byte(strings.TrimSpace(line[index+len(coldReportMarker):])), &report); err != nil {
					t.Fatal(err)
				}
				return report
			}
		}
		t.Fatalf("missing process report for mode %s: %s", mode, result)
		return coldReport{}
	}
	uninterrupted := run("full", fullRoot)
	prefix := run("prefix", resumedRoot)
	if prefix.Sequence != 1 || prefix.Observations != 1 {
		t.Fatal("prefix custody invalid")
	}
	resumed := run("resume", resumedRoot)
	if uninterrupted.Sequence != 2 || uninterrupted.Observations != 2 || resumed.Sequence != 2 || resumed.Observations != 2 {
		t.Fatal("continuation lost action or training cursor")
	}
	if uninterrupted != resumed {
		t.Fatalf("cold-process continuation drift: full=%+v restarted=%+v", uninterrupted, resumed)
	}
	stale := run("stale", resumedRoot)
	if !stale.RejectedStale || stale.Sequence != 2 {
		t.Fatal("stale checkpoint replay was not rejected")
	}
	if os.Getenv(coldArtifactRootEnv) != "" {
		prefixBytes, err := os.ReadFile(snap)
		if err != nil { t.Fatal(err) }
		fullBytes, err := os.ReadFile(filepath.Join(fullRoot, "final.snapshot"))
		if err != nil { t.Fatal(err) }
		resumedBytes, err := os.ReadFile(filepath.Join(resumedRoot, "final.snapshot"))
		if err != nil { t.Fatal(err) }
		parentHash := sha256.Sum256(prefixBytes)
		childHash := sha256.Sum256(resumedBytes)
		if hex.EncodeToString(parentHash[:]) != prefix.SnapshotSHA256 || hex.EncodeToString(childHash[:]) != resumed.SnapshotSHA256 || !bytes.Equal(fullBytes, resumedBytes) {
			t.Fatal("durable checkpoint artifact or uninterrupted native state mismatched")
		}
		sourceSHA := os.Getenv(coldSourceSHAEnv)
		if len(sourceSHA) != 40 || strings.Trim(sourceSHA, "0123456789abcdef") != "" {
			t.Fatal("exact checked-out source SHA missing")
		}
		manifest := map[string]any{
			"schema": "wingless.nci04l5a.native-cold-state-sourcefree.dev.v1",
			"source_commit_sha": sourceSHA,
			"project": "Wingless",
			"parent_snapshot_sha256": prefix.SnapshotSHA256,
			"child_snapshot_sha256": resumed.SnapshotSHA256,
			"uninterrupted_snapshot_sha256": uninterrupted.SnapshotSHA256,
			"source_probe_sha256": resumed.SourceSHA256,
			"observations": resumed.Observations,
			"action_sequence": resumed.Sequence,
			"checkpoint_bytes_equal_after_cold_process": true,
			"rejected_stale_checkpoint": stale.RejectedStale,
			"independent_source_first_use": false,
			"source_free_generated_fixture": true,
			"matched_control_superiority": "UNVERIFIED",
			"unfamiliar_transfer": "UNVERIFIED",
			"scientific_accepted": false,
			"promotion_authority": false,
		}
		raw, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil { t.Fatal(err) }
		if err := os.WriteFile(filepath.Join(tmp, "evidence.json"), append(raw, '\n'), 0600); err != nil { t.Fatal(err) }
	}
	t.Logf("NCI04L5A_WING_FIXTURE_COLD_PROCESS_PASS snapshot_sha256=%s sequence=%d observations=%d scientific_first_use=false", resumed.SnapshotSHA256, resumed.Sequence, resumed.Observations)
}