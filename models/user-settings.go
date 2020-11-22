package models

import (
	"github.com/google/uuid"
)

type UserSetting struct {
	ID       uuid.UUID `validation:"required"`
	Settings Settings  `validation:"required"`
}

// Assets structure contains the identification of all user related documents images.
type Settings struct {
	TwoStepsVerif bool `json:"two_steps_verif" validation:"required"`
}
