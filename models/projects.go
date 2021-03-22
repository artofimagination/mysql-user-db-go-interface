package models

import (
	"github.com/google/uuid"
)

type ProjectData struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	ProductID uuid.UUID `json:"product-id" validate:"required"`
	Assets    *Asset    `json:"assets" validate:"required"`
	Details   *Asset    `json:"details" validate:"required"`
}

type UserProject struct {
	ProjectData *ProjectData
	Privilege   int
}

type Project struct {
	ID        uuid.UUID `validation:"required"`
	ProductID uuid.UUID `validation:"required"`
	AssetsID  uuid.UUID `validation:"required"`
	DetailsID uuid.UUID `validation:"required"`
}

type UserProjectIDs struct {
	ProjectMap     map[uuid.UUID]int
	ProjectIDArray []uuid.UUID
}

type ProjectUserIDs struct {
	UserMap     map[uuid.UUID]int
	UserIDArray []uuid.UUID
}

// Viewers describes a set of project data. Each set has a single viewer ID.
// Whichever user possese this ID can view the data with this ID. The set of project data has a single owner user.
// This structure contains the owner and the list of users who can view the data.
type ViewersList []Viewer

type Viewer struct {
	IsOwner   bool
	ViewerID  uuid.UUID
	ProjectID uuid.UUID
}

type ViewUsers struct {
	OwnerID   uuid.UUID
	UsersList []uuid.UUID
}

func (f *RepoFunctions) NewProject(productID *uuid.UUID, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Project, error) {
	var p Project

	newID, err := f.UUIDImpl.NewUUID()
	if err != nil {
		return nil, err
	}

	p.ID = newID
	p.ProductID = *productID
	p.DetailsID = *detailsID
	p.AssetsID = *assetsID

	return &p, nil
}
