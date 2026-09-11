//go:build !windows

package modelhost

import "os"

type containment struct{}

func newContainment() (*containment, error)     { return &containment{}, nil }
func (c *containment) Attach(*os.Process) error { return nil }
func (c *containment) Close()                   {}
func containmentName() string                   { return "direct_child_only" }
