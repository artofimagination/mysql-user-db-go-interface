package models

import (
	"encoding/json"
)

// Feature describes the available simulation features.
type Feature struct {
	ID     int             `json:"id"`
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}
