package toolboundary

// ObservePreaction captures read-only source facts before Execute has a result.
// The snapshot is not a prediction, a success label, or an atomic future guarantee.
import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type PreactionObservation struct {
    Resource string `json:"resource"` // "missing" or "regular"; other conditions fail closed
    Size int64 `json:"size"`
    SHA256 string `json:"sha256,omitempty"`
}

func ObservePreaction(parent context.Context, p Policy, a Proposal) (PreactionObservation, error) {
    var none PreactionObservation
    if a.Tool != "read_file" || !has(p.AllowedTools, a.Tool) || !has(p.AllowedPaths, a.Path) {
        return none, errors.New("tool/path unauthorized")
    }
    if p.Timeout <= 0 || p.Timeout > time.Minute || p.MaxOutputBytes <= 0 || p.MaxOutputBytes > 1048576 || !filepath.IsAbs(p.Workspace) {
        return none, errors.New("unbounded tool policy")
    }
    if !filepath.IsLocal(a.Path) || strings.Contains(a.Path, ":") || strings.ContainsRune(a.Path, 92) {
        return none, errors.New("workspace escape denied")
    }
    ctx, cancel := context.WithTimeout(parent, p.Timeout)
    defer cancel()
    if err := ctx.Err(); err != nil { return none, err }
    root, err := os.OpenRoot(p.Workspace)
    if err != nil { return none, err }
    defer root.Close()
    file, err := root.Open(a.Path)
    if errors.Is(err, os.ErrNotExist) { return PreactionObservation{Resource:"missing"}, nil }
    if err != nil { return none, fmt.Errorf("rooted preaction inspection denied: %w", err) }
    defer file.Close()
    info, err := file.Stat()
    if err != nil { return none, err }
    if !info.Mode().IsRegular() { return none, errors.New("nonregular file denied") }
    if info.Size() > int64(p.MaxOutputBytes) { return none, errors.New("preaction output bound exceeded") }
    h := sha256.New()
    count, err := io.Copy(h, io.LimitReader(file, int64(p.MaxOutputBytes)+1))
    if err != nil { return none, err }
    if count > int64(p.MaxOutputBytes) { return none, errors.New("preaction output bound exceeded") }
    after, err := file.Stat()
    if err != nil { return none, err }
    if count != info.Size() || after.Size() != info.Size() || !after.ModTime().Equal(info.ModTime()) {
        return none, errors.New("preaction source changed during observation")
    }
    if err := ctx.Err(); err != nil { return none, err }
    return PreactionObservation{Resource:"regular", Size:count, SHA256:hex.EncodeToString(h.Sum(nil))}, nil
}
