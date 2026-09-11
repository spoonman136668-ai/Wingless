package inference

import (
	"encoding/json"
	"fmt"
)

// OutputConstraint is a backend-neutral seam. No current adapter claims support.
// Model output cannot request a stronger capability than trusted request configuration.
type OutputConstraint struct {
	Kind   string          `json:"kind"`
	Schema json.RawMessage `json:"schema,omitempty"`
}

func (c *OutputConstraint) Validate() error {
	if c == nil {
		return nil
	}
	if c.Kind != "json_schema" || len(c.Schema) == 0 || len(c.Schema) > 16384 || !json.Valid(c.Schema) {
		return fmt.Errorf("invalid bounded output constraint")
	}
	var obj map[string]any
	if json.Unmarshal(c.Schema, &obj) != nil || obj == nil {
		return fmt.Errorf("schema must be object")
	}
	return nil
}
func ConstraintCapability(c *OutputConstraint) string {
	if c == nil {
		return ""
	}
	return "output:" + c.Kind
}
