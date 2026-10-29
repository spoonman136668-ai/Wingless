package dk8

import (
    "bytes"
    "testing"

    dk7 "github.com/spoonman136668-ai/Wingless/research/native/itd02/dk7"
)

func public() dk7.Public {
    p := dk7.Public{Image: make([][][]int, 7), Direction: 1, Mission: dk7.Mission}
    for x := 0; x < 7; x++ {
        p.Image[x] = make([][]int, 7)
        for y := 0; y < 7; y++ {
            p.Image[x][y] = []int{1, 0, 0}
        }
    }
    p.Image[3][5] = []int{5, 0, 0}
    return p
}

func TestUnknownNativeActionWithoutExternalFallbackAndColdIdentity(t *testing.T) {
    e, p := New(), public()
    first, status, err := e.Choose(p)
    if err != nil || first == nil || *first < 0 || *first > 6 || status != "SELECT_NATIVE_LEAST_OBSERVED_ACTION" {
        t.Fatalf("native zero-reward action missing: %v %s %v", first, status, err)
    }
    if err := e.Commit(p, *first); err != nil { t.Fatal(err) }
    if _, _, err := e.Choose(p); err == nil { t.Fatal("reproposal during pending action") }
    if _, err := e.Snapshot(); err == nil { t.Fatal("pending checkpoint must reject") }
    if err := e.Feedback(p, (*first+1)%7, p, 0, false); err == nil {
        t.Fatal("wrong-action feedback accepted")
    }
    if e.Native.Samples != 0 { t.Fatal("wrong action mutated native evidence") }
    if err := e.Feedback(p, *first, p, 0, false); err != nil { t.Fatal(err) }
    second, status, err := e.Choose(p)
    if err != nil || second == nil || *second == *first || status != "SELECT_NATIVE_LEAST_OBSERVED_ACTION" {
        t.Fatalf("sample counts did not change native action: %v -> %v %s %v", first, second, status, err)
    }
    raw, err := e.Snapshot()
    if err != nil { t.Fatal(err) }
    restored, err := Restore(raw)
    if err != nil { t.Fatal(err) }
    copyRaw, err := restored.Snapshot()
    if err != nil || !bytes.Equal(copyRaw, raw) { t.Fatal("cold snapshot identity") }
    cold, coldStatus, err := restored.Choose(p)
    if err != nil || cold == nil || *cold != *second || coldStatus != status {
        t.Fatalf("cold native action identity changed: %v -> %v %v", second, cold, err)
    }
    corrupted := append(append([]byte{}, raw...), 0)
    if _, err := Restore(corrupted); err == nil { t.Fatal("tampered source checkpoint accepted") }
}

func TestObservedRewardTakesPrecedenceAndInvalidMissionFailsClosed(t *testing.T) {
    e, p := New(), public()
    if err := e.ObserveFit(p, 5, p, 1, false); err != nil { t.Fatal(err) }
    selected, status, err := e.Choose(p)
    if err != nil || selected == nil || *selected != 5 || status != "SELECT_NATIVE_OBSERVED_POSITIVE_REWARD" {
        t.Fatalf("native learned reward priority broken: %v %s %v", selected, status, err)
    }
    originalSamples := e.Native.Samples
    bad := public()
    bad.Mission = "hidden evaluator goal"
    if a, _, err := e.Choose(bad); err == nil || a != nil { t.Fatal("invalid public mission accepted") }
    if err := e.ObserveFit(bad, 1, p, 0, false); err == nil { t.Fatal("invalid fit mission accepted") }
    if e.Native.Samples != originalSamples { t.Fatal("invalid public input mutated native memory") }
}

func TestBoundedNativeProposalsPreserveSelectedFeedbackOnly(t *testing.T) {
    e, p := New(), public()
    for i := 0; i < 7; i++ {
        a, status, err := e.Choose(p)
        if err != nil || a == nil || status != "SELECT_NATIVE_LEAST_OBSERVED_ACTION" {
            t.Fatalf("missing native action at %d: %v %s %v", i, a, status, err)
        }
        if err := e.Commit(p, *a); err != nil { t.Fatal(err) }
        if err := e.Feedback(p, *a, p, 0, false); err != nil { t.Fatal(err) }
    }
    ctx, err := dk7.Context(p)
    if err != nil { t.Fatal(err) }
    bins := e.Native.Counts[ctx]
    for i, b := range bins {
        if b.Count != 1 || b.Positive != 0 { t.Fatalf("action %d has inconsistent native evidence %+v", i, b) }
    }
    if e.Native.Samples != 7 { t.Fatal("invented or dropped selected feedback") }
}
