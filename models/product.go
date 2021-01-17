package models

import (
	"github.com/google/uuid"
)

type ProductData struct {
	ID      uuid.UUID
	Name    string
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
	AssetsID  uuid.UUID `validation:"required"`
	DetailsID uuid.UUID `validation:"required"`
}

type UserProductIDs struct {
	ProductMap     map[uuid.UUID]int
	ProductIDArray []uuid.UUID
}
type ProductUserIDs struct {
	UserMap     map[uuid.UUID]int
	UserIDArray []uuid.UUID
}

// NewProduct creates a new product instance where details describe the configuration of the product
// and references contain all asset references.
func (f *RepoFunctions) NewProduct(name string, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Product, error) {
	var p Product

	newID, err := f.UUIDImpl.NewUUID()
	if err != nil {
		return nil, err
	}

	p.ID = newID
	p.Name = name
	p.DetailsID = *detailsID
	p.AssetsID = *assetsID

	return &p, nil
}
