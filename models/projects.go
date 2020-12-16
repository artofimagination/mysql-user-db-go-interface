package models

import (
	"github.com/google/uuid"
)

type ProjectData struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Assets    *Asset
	Details   *Asset
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
type Viewers struct {
	OwnerID  uuid.UUID
	UserList []int
}

// ProjectDataList holds the list of viewers a user owns. The ID-s in the list allow the user
// to view the project data belonging to those ID-s.
type ProjectDataList []int

func (*RepoInterface) NewProject(productID *uuid.UUID, detailsID *uuid.UUID, assetsID *uuid.UUID) (*Project, error) {
	var p Project

	newID, err := UUIDImpl.NewUUID()
	if err != nil {
		return nil, err
	}

	p.ID = newID
	p.ProductID = *productID
	p.DetailsID = *detailsID
	p.AssetsID = *assetsID

	return &p, nil
}
