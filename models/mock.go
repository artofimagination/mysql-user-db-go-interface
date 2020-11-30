package models

import (
	"github.com/google/uuid"
)

// UUIDImplMock overwrites the default github uuid library implementation.
type UUIDImplMock struct {
	uuidMock uuid.UUID
	err      error
}

func (i UUIDImplMock) NewUUID() (uuid.UUID, error) {
	return i.uuidMock, i.err
}
