package models

type Privilege struct {
	ID          int    `validation:"required"`
	Name        string `validation:"required"`
	Description string `validation:"required"`
}

type Privileges []*Privilege

func (l Privileges) IsValidPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege {
			return true
		}
	}
	return false
}

func (l Privileges) IsOwnerPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege && value.Name == "Owner" {
			return true
		}
	}
	return false
}

func (l Privileges) IsPartnerPrivilege(privilege int) bool {
	for _, value := range l {
		if value.ID == privilege && value.Name == "Partner" {
			return true
		}
	}
	return false
}
