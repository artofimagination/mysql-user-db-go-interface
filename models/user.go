package models

import (
	"errors"

	"github.com/google/uuid"
)

// Errors called in multiple places (for example in unittests).

var ErrInvalidSettingsID = "Invalid settings uuid"
var ErrInvalidAssetsID = "Invalid assets uuid"

type UserData struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Name     string    `json:"username" validate:"required"`
	Email    string    `json:"email" validate:"required"`
	Password []byte    `json:"password" validate:"required"`
	Settings *Asset    `json:"settings" validate:"required"`
	Assets   *Asset    `json:"assets" validate:"required"`
}

type ProductUser struct {
	UserData  UserData
	Privilege int
}

// User defines the user structures. Each user must have an associated settings entry.
type User struct {
	ID         uuid.UUID
	Name       string
	Email      string
	Password   []byte
	SettingsID uuid.UUID
	AssetsID   uuid.UUID
}

func (f *RepoFunctions) NewUser(
	name string,
	email string,
	password []byte,
	settingsID uuid.UUID,
	assetsID uuid.UUID) (*User, error) {
	var u User

	emptyUUID := uuid.UUID{}
	if settingsID == emptyUUID {
		return nil, errors.New(ErrInvalidSettingsID)
	}

	if assetsID == emptyUUID {
		return nil, errors.New(ErrProductDetailsNotInitialised)
	}

	newID, err := f.UUIDImpl.NewUUID()
	if err != nil {
		return nil, err
	}

	u.ID = newID
	u.Name = name
	u.Email = email
	u.Password = password
	u.SettingsID = settingsID
	u.AssetsID = assetsID

	return &u, nil
}
