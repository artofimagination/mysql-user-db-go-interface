package models

import (
	"github.com/google/uuid"
)

// UUIDInterfaceMock overwrites the default github uuid library implementation.
type UUIDInterfaceMock struct {
	uuidMock uuid.UUID
	err      error
}

func (i UUIDInterfaceMock) NewUUID() (uuid.UUID, error) {
	return i.uuidMock, i.err
}
