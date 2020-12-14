package models

import (
	"github.com/google/uuid"
)

// This contants describe the visibility of proucts and projects
const (
	Public    = iota // Public is visible for everyone including nont registered visitors
	Protected        // Content is available for registered users
	Private          // Content is available for the owner and users the product or project is shared with.
)

type InterfaceCommon interface {
	NewAsset(references DataMap, generatePath func(assetID *uuid.UUID) (string, error)) (*Asset, error)
	NewUser(
		name string,
		email string,
		password []byte,
		settingsID uuid.UUID,
		assetsID uuid.UUID) (*User, error)
	NewProduct(name string, public bool, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Product, error)
	NewProject(detailsID *uuid.UUID, assetsID *uuid.UUID) (*Project, error)
}

type UUIDInterfaceCommon interface {
	NewUUID() (uuid.UUID, error)
}

type RepoUUIDInterface struct {
}

type RepoInterface struct {
}

type Privilege struct {
	ID          int    `validation:"required"`
	Name        string `validation:"required"`
	Description string `validation:"required"`
}

type Privileges []Privilege

func (l Privileges) IsValidPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege {
			return true
		}
	}
	return false
}

func (l Privileges) IsOwnerPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege && value.Name == "Owner" {
			return true
		}
	}
	return false
}

func (l Privileges) IsPartnerPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege && value.Name == "Partner" {
			return true
		}
	}
	return false
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
