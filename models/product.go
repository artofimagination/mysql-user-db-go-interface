package models

import (
	"github.com/google/uuid"
)

type ProductData struct {
	ID      uuid.UUID
	Name    string
	Public  bool
	Assets  *Asset
	Details *Asset
}

type UserProduct struct {
	ProductData ProductData
	Privilege   int
}

type Product struct {
	ID        uuid.UUID `validation:"required"`
	Name      string    `validation:"required"`
	Public    bool
	AssetsID  uuid.UUID `validation:"required"`
	DetailsID uuid.UUID `validation:"required"`
}

type Privilege struct {
	ID          int    `validation:"required"`
	Name        string `validation:"required"`
	Description string `validation:"required"`
}

type Privileges []*Privilege
type UserProductIDs struct {
	ProductMap     map[uuid.UUID]int
	ProductIDArray []uuid.UUID
}
type ProductUserIDs struct {
	UserMap     map[uuid.UUID]int
	UserIDArray []uuid.UUID
}

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

// NewProduct creates a new product instance where details describe the configuration of the product
// and references contain all asset references.
func (*RepoInterface) NewProduct(name string, public bool, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Product, error) {
	var p Product

	newID, err := UUIDImpl.NewUUID()
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
