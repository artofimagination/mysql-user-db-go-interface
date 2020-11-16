package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Data struct {
	ID       int     `json:"id"`
	DataType int     `json:"type"`
	Speed    float32 `json:"speed"`
}

// Project defines the project structure.
type Project struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	UserID    uuid.UUID       `json:"user_id"`
	FeatureID int             `json:"features_id"`
	Config    json.RawMessage `json:"config"`
}
