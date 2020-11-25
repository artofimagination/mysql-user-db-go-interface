package models

import (
	"errors"

	"github.com/google/uuid"
)

type UserSetting struct {
	ID       uuid.UUID `validation:"required"`
	Settings Settings  `validation:"required"`
}

// Errors called in multiple places (for example in unittests).

var ErrSettingNotInitialised = "Settings not initialised"

const (
	TwoStepsVerif = "two_steps_verif"
)

type Settings map[string]interface{}

func (RepoInterface) NewUserSettings(settings Settings) (*UserSetting, error) {
	var s UserSetting

	if settings == nil {
		return nil, errors.New(ErrSettingNotInitialised)
	}

	newID, err := Interface.NewUUID()
	if err != nil {
		return nil, err
	}

	s.ID = newID
	s.Settings = settings

	return &s, nil
}
