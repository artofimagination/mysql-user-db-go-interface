package models

import (
	"github.com/google/uuid"
)

// This contants describe the visibility of proucts and projects
const (
	Public    = "Public"    // Public is visible for everyone including nont registered visitors
	Protected = "Protected" // Content is available for registered users
	Private   = "Private"   // Content is available for the owner and users the product or project is shared with.
)

type ModelFunctionsCommon interface {
	NewAsset(references DataMap, generatePath func(assetID *uuid.UUID) (string, error)) (*Asset, error)
	NewUser(
		name string,
		email string,
		password []byte,
		settingsID uuid.UUID,
		assetsID uuid.UUID) (*User, error)
	NewProduct(name string, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Product, error)
	NewProject(productID *uuid.UUID, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Project, error)
}

type UUIDCommon interface {
	NewUUID() (uuid.UUID, error)
}

type RepoUUID struct {
}

type RepoFunctions struct {
	UUIDImpl UUIDCommon
}

// NewUUID is a wrapper to allow mocking
func (*RepoUUID) NewUUID() (uuid.UUID, error) {
	var newID uuid.UUID
	newID, err := uuid.NewUUID()
	if err != nil {
		return newID, err
	}
	return newID, nil
}
