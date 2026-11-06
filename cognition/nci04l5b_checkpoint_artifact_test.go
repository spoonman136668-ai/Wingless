package cognition

import (
    "bytes"
    "context"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

// Independent evidence export using the previously qualified L5A native fixture.
// No modification to the frozen L5A tests, accepted learner, or original source.
func TestNCI04L5BCheckpointArtifactHelper(t *testing.T) {
    mode := os.Getenv("NCI04L5B_MODE")
    if mode == "" { return }
    root := os.Getenv("NCI04L5B_PROBE_ROOT")
    snapshotPath := os.Getenv("NCI04L5B_PREFIX_CHECKPOINT")
    if root == "" || snapshotPath == "" { t.Fatal("missing source-free fixture path") }
    var probe *NativeProbe
    var err error
    if mode == "full" || mode == "prefix" {
        probe, err = NewNativeProbe(root)
    } else {
        probe, err = OpenNativeProbe(root)
    }
    if err != nil { t.Fatal(err) }
    session := &ExecutedLearningSession{Probe: probe}
    if mode == "resume" || mode == "stale" {
        raw, err := os.ReadFile(snapshotPath)
        if err != nil { t.Fatal(err) }
        err = session.Restore(raw)
        if mode == "stale" {
            if err == nil || probe.sequence != 2 { t.Fatal("stale checkpoint accepted or source cursor corrupted") }
            report, _ := json.Marshal(coldReport{RejectedStale: true, Sequence: probe.sequence})
            fmt.Println(coldReportMarker + string(report))
            return
        }
        if err != nil { t.Fatal(err) }
    }
    if mode == "full" || mode == "prefix" {
        if _, err := session.ActAndLearn(context.Background(), false); err != nil { t.Fatal(err) }
    }
    if mode == "full" || mode == "resume" {
        if _, err := session.ActAndLearn(context.Background(), true); err != nil { t.Fatal(err) }
    }
    if mode != "full" && mode != "prefix" && mode != "resume" { t.Fatal("invalid fixture mode") }
    saved, err := session.Snapshot()
    if err != nil { t.Fatal(err) }
    if mode == "prefix" {
        if err := os.WriteFile(snapshotPath, saved, 0600); err != nil { t.Fatal(err) }
    } else {
        if err := os.WriteFile(filepath.Join(root, "final.snapshot"), saved, 0600); err != nil { t.Fatal(err) }
    }
    prediction, err := session.PredictExecutable(true)
    if err != nil { t.Fatal(err) }
    rawSource, err := probeRead(root)
    if err != nil { t.Fatal(err) }
    sum := sha256.Sum256(saved)
    r := coldReport{SnapshotSHA256: hex.EncodeToString(sum[:]), SourceSHA256: probeHash(rawSource),
        Prediction: prediction, Observations: session.State.Observations, Sequence: probe.sequence}
    payload, err := json.Marshal(r)
    if err != nil { t.Fatal(err) }
    fmt.Println(coldReportMarker + string(payload))
}

func TestNCI04L5BCheckpointArtifactRetention(t *testing.T) {
    if os.Getenv("NCI04L5B_MODE") != "" { return }
    temp := t.TempDir()
    export := os.Getenv("NCI04L5B_ARTIFACT_ROOT")
    if export != "" {
        if !filepath.IsAbs(export) { t.Fatal("artifact root must be absolute") }
        if err := os.Mkdir(export, 0700); err != nil { t.Fatal("artifact root not newly created", err) }
        temp = export
    }
    exe, err := os.Executable()
    if err != nil { t.Fatal(err) }
    parent := filepath.Join(temp, "prefix.snapshot")
    fullRoot := filepath.Join(temp, "full")
    resumedRoot := filepath.Join(temp, "resumed")
    run := func(mode, root string) coldReport {
        t.Helper()
        child := exec.Command(exe, "-test.run=^TestNCI04L5BCheckpointArtifactHelper$", "-test.v")
        child.Env = append(os.Environ(), "NCI04L5B_MODE="+mode,
            "NCI04L5B_PROBE_ROOT="+root, "NCI04L5B_PREFIX_CHECKPOINT="+parent)
        output, err := child.CombinedOutput()
        if err != nil { t.Fatalf("cold-process artifact %s: %v\n%s", mode, err, output) }
        for _, line := range strings.Split(string(output), "\n") {
            if ix := strings.Index(line, coldReportMarker); ix >= 0 {
                var r coldReport
                if err := json.Unmarshal([]byte(strings.TrimSpace(line[ix+len(coldReportMarker):])), &r); err != nil { t.Fatal(err) }
                return r
            }
        }
        t.Fatalf("missing native process report for %s", mode)
        return coldReport{}
    }
    full := run("full", fullRoot)
    prefix := run("prefix", resumedRoot)
    resumed := run("resume", resumedRoot)
    stale := run("stale", resumedRoot)
    if prefix.Sequence != 1 || prefix.Observations != 1 ||
        full.Sequence != 2 || resumed.Sequence != 2 ||
        full.Observations != 2 || resumed.Observations != 2 ||
        full != resumed || !stale.RejectedStale || stale.Sequence != 2 {
        t.Fatal("native cold-start equivalence or stale replay rejection failed")
    }
    prefBytes, err := os.ReadFile(parent)
    if err != nil { t.Fatal(err) }
    fullBytes, err := os.ReadFile(filepath.Join(fullRoot, "final.snapshot"))
    if err != nil { t.Fatal(err) }
    resumedBytes, err := os.ReadFile(filepath.Join(resumedRoot, "final.snapshot"))
    if err != nil { t.Fatal(err) }
    before := sha256.Sum256(prefBytes)
    after := sha256.Sum256(resumedBytes)
    if hex.EncodeToString(before[:]) != prefix.SnapshotSHA256 ||
        hex.EncodeToString(after[:]) != resumed.SnapshotSHA256 ||
        !bytes.Equal(fullBytes, resumedBytes) { t.Fatal("exact raw checkpoint byte custody failed") }
    if export != "" {
        sha := os.Getenv("NCI04L5B_SOURCE_SHA")
        if len(sha) != 40 || strings.Trim(sha, "0123456789abcdef") != "" { t.Fatal("source SHA missing") }
        report := map[string]any{
            "schema": "wingless.nci04l5b.native-cold-state-sourcefree.dev.v1",
            "project": "Wingless",
            "source_commit_sha": sha,
            "prefix_checkpoint_sha256": prefix.SnapshotSHA256,
            "continued_checkpoint_sha256": resumed.SnapshotSHA256,
            "uninterrupted_checkpoint_sha256": full.SnapshotSHA256,
            "source_probe_sha256": resumed.SourceSHA256,
            "source_free_generated_fixture": true,
            "independent_process_continuation": true,
            "stale_checkpoint_rejected": stale.RejectedStale,
            "learned_observations": resumed.Observations,
            "independent_source_first_use": false,
            "scientific_accepted": false,
            "accepted_model_or_portable_state_former": false,
            "matched_baseline_superiority": "UNVERIFIED",
            "unfamiliar_transfer": "UNVERIFIED",
            "promotion_authority": false,
        }
        data, err := json.MarshalIndent(report, "", "  ")
        if err != nil { t.Fatal(err) }
        if err := os.WriteFile(filepath.Join(temp,"evidence.json"), append(data, '\n'),0600); err != nil { t.Fatal(err) }
    }
}
