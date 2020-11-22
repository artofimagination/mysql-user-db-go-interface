package models

import (
	"github.com/google/uuid"
)

// User defines the user structures. Each user must have an associated settings entry.
type User struct {
	ID         uuid.UUID `validation:"required"`
	Name       string    `validation:"required"`
	Email      string    `validation:"required"`
	Password   string    `validation:"required"`
	SettingsID uuid.UUID `validation:"required"`
	AssetsID   uuid.UUID `validation:"required"`
}
