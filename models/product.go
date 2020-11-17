package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID      uuid.UUID `json:"id" validation:"required"`
	Name    string    `json:"name" validation:"required"`
	Public  bool      `json:"public" validation:"required"`
	Details Details   `json:"details" validation:"required"`
}

type Details struct {
	SupportClients bool `json:"support_clients" validation:"required"`
	ClientUI       bool `json:"client_ui" validation:"required"`
	ProjectUI      bool `json:"project_ui" validation:"required"`
	Requires3D     bool `json:"requires_3d" validation:"required"`
}

type Privilege struct {
	ID          int    `json:"id" validation:"required"`
	Name        string `json:"name" validation:"required"`
	Description string `json:"description" validation:"required"`
}

type Privileges []Privilege
type OwnershipMap map[uuid.UUID]int

func (l Privileges) IsValidPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege {
			return true
		}
	}
	return false
}
