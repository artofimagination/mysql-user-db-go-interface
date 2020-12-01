package dbcontrollers

import (
	"github.com/google/uuid"
)

type ProjectDBCommon interface {
	DeleteProjects(productID uuid.UUID) error
}

// ProjectDBDummy is a dummy implementation of ProjectDB interface.
// Since project handling can be completely decoupled from user/product management
// It is up to the user what implementation he/she invokes.
type ProjectDBDummy struct {
}

func (ProjectDBDummy) DeleteProjects(productID uuid.UUID) error {
	return nil
}

var projectdb ProjectDBCommon
