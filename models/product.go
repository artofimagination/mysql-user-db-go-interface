package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `validation:"required"`
	Name      string    `validation:"required"`
	Public    bool      `validation:"required"`
	AssetsID  uuid.UUID `validation:"required"`
	DetailsID uuid.UUID `validation:"required"`
}

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
func (RepoInterface) NewProduct(name string, public bool, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Product, error) {
	var p Product

	newID, err := UUIDInterface.NewUUID()
	if err != nil {
		return nil, err
	}

	p.ID = newID
	p.Name = name
	p.Public = public
	p.DetailsID = *detailsID
	p.AssetsID = *assetsID

	return &p, nil
}
