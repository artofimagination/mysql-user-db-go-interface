package models

import (
	"github.com/google/uuid"
)

type UserSetting struct {
	ID       uuid.UUID `json:"id" validation:"required"`
	Settings Settings  `json:"settings" validation:"required"`
}

// Assets structure contains the identification of all user related documents images.
type Settings struct {
	TwoStepsVerif bool `json:"two_steps_verif,omitempty"`
}
