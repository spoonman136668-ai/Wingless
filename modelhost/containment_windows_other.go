//go:build windows && !amd64

package modelhost

import (
	"fmt"
	"os"
)

type containment struct{}

func newContainment() (*containment, error) {
	return nil, fmt.Errorf("Windows containment requires amd64")
}
func (c *containment) Attach(*os.Process) error { return fmt.Errorf("unavailable") }
func (c *containment) Close()                   {}
func containmentName() string                   { return "unavailable" }
