package models

import (
	"errors"

	"github.com/google/uuid"
)

type Product struct {
	ID       uuid.UUID `validation:"required"`
	Name     string    `validation:"required"`
	Public   bool      `validation:"required"`
	AssetsID uuid.UUID `validation:"required"`
	Details  Details   `validation:"required"`
}

const (
	SupportClients = "support_clients"
	ClientUI       = "client_ui"
	ProjectUI      = "project_ui"
	Requires3D     = "requires_3d"
	HasTrial       = "has_trial"
	IsFree         = "is_free"
)

type Details map[string]interface{}

// Errors called in multiple places (for example in unittests).

var ErrProductDetailsNotInitialised = "Details map not initialised"

type Privilege struct {
	ID          int    `validation:"required"`
	Name        string `validation:"required"`
	Description string `validation:"required"`
}

type Privileges []Privilege
type UserProducts map[uuid.UUID]int
type ProductUsers map[uuid.UUID]int

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

// NewProduct creates a new product instance where details describe the configuration of the product
// and references contain all asset references.
func (RepoInterface) NewProduct(name string, public bool, details Details, assetsID *uuid.UUID) (*Product, error) {
	var p Product

	if details == nil {
		return nil, errors.New(ErrProductDetailsNotInitialised)
	}

	newID, err := Interface.NewUUID()
	if err != nil {
		return nil, err
	}

	p.ID = newID
	p.Name = name
	p.Public = public
	p.Details = details
	p.AssetsID = *assetsID

	return &p, nil
}
