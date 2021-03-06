package models

import (
	"github.com/google/uuid"
)

var ModelFunctions *RepoFunctions
var UUIDImpl *UUIDImplMock

// UUIDImplMock overwrites the default github uuid library implementation.
type UUIDImplMock struct {
	uuidMock uuid.UUID
	err      error
}

func (i *UUIDImplMock) NewUUID() (uuid.UUID, error) {
	return i.uuidMock, i.err
}
