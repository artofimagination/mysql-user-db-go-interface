package models

import (
	"errors"

	"github.com/google/uuid"
)

type ProductDetails struct {
	ID      uuid.UUID `validation:"required"`
	Details Details   `validation:"required"`
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

func (RepoInterface) NewProductDetails(details Details) (*ProductDetails, error) {
	var d ProductDetails

	if details == nil {
		return nil, errors.New(ErrProductDetailsNotInitialised)
	}

	newID, err := UUIDInterface.NewUUID()
	if err != nil {
		return nil, err
	}

	d.ID = newID
	d.Details = details

	return &d, nil
}
