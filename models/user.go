package models

import (
	"github.com/google/uuid"
)

// User defines the user structures. Each user must have an associated settings entry.
type User struct {
	ID         uuid.UUID `json:"id" validation:"required"`
	Name       string    `json:"name" validation:"required"`
	Email      string    `json:"email" validation:"required"`
	Password   string    `json:"password" validation:"required"`
	SettingsID uuid.UUID `json:"user_settings_id" validation:"required"`
	AssetsID   uuid.UUID `json:"user_assets_id" validation:"required"`
}
