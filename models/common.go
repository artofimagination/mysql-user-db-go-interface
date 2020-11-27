package models

import (
	"github.com/google/uuid"
)

var NullUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

type InterfaceCommon interface {
	NewAsset(references DataMap, generatePath func(assetID *uuid.UUID) string) (*Asset, error)
	NewUser(
		name string,
		email string,
		password []byte,
		settingsID uuid.UUID,
		assetsID uuid.UUID) (*User, error)
	NewProduct(name string, public bool, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Product, error)
}

type UUIDInterfaceCommon interface {
	NewUUID() (uuid.UUID, error)
}

type RepoUUIDInterface struct {
}

type RepoInterface struct {
}

// NewUUID is a wrapper to allow mocking
func (RepoUUIDInterface) NewUUID() (uuid.UUID, error) {
	var newID uuid.UUID
	newID, err := uuid.NewUUID()
	if err != nil {
		return newID, err
	}
	return newID, nil
}

var Interface InterfaceCommon
var UUIDImpl UUIDInterfaceCommon
