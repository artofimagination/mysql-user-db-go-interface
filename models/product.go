package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID      uuid.UUID `validation:"required"`
	Name    string    `validation:"required"`
	Public  bool      `validation:"required"`
	Details Details   `validation:"required"`
}

type Details struct {
	SupportClients bool `json:"support_clients" validation:"required"`
	ClientUI       bool `json:"client_ui" validation:"required"`
	ProjectUI      bool `json:"project_ui" validation:"required"`
	Requires3D     bool `json:"requires_3d" validation:"required"`
}

type Privilege struct {
	ID          int    `validation:"required"`
	Name        string `validation:"required"`
	Description string `validation:"required"`
}

type Privileges []Privilege
type UserProducts map[uuid.UUID]int
type ProductUsers map[uuid.UUID]int

func (l Privileges) IsValidPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege {
			return true
		}
	}
	return false
}

func NewProduct(name string, public bool, details *Details) (*Product, error) {
	var p Product
	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	p.ID = newID
	p.Name = name
	p.Public = public
	p.Details = *details

	return &p, nil
}
