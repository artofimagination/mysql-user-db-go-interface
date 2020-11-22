package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type LeaveType string

const (
	Avatar            = "avatar"
	ProfileBackground = "profile-background"
)

var NullUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

type Asset struct {
	ID         uuid.UUID  `validation:"required"`
	References References `validation:"required"`
	Path       string
}

// Assets structure contains the identification of all user related documents images.
type References struct {
	AvatarID          uuid.UUID `json:"avatar_id,omitempty"`
	ProfileBackground uuid.UUID `json:"profile_backgr,omitempty"`
}

func (r *Asset) GetPath(typeString string) (string, error) {
	defaultPath := "/assets/images/avatar.jpg"
	var ID uuid.UUID

	switch typeString {
	case Avatar:
		ID = r.References.AvatarID
	case ProfileBackground:
		ID = r.References.ProfileBackground
	default:
		return defaultPath, errors.New("Unknown asset reference type")
	}

	if ID == NullUUID {
		return defaultPath, nil
	}

	return fmt.Sprintf("%s/%s.jpg", r.Path, ID), nil
}

func (r *Asset) SetID(typeString string) error {
	if r.References.AvatarID != NullUUID {
		return nil
	}

	newID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	switch typeString {
	case "avatar":
		r.References.AvatarID = newID
	default:
		return errors.New("Unknown asset reference type")
	}

	return nil
}

func (r *Asset) ClearID(typeString string) error {
	switch typeString {
	case "avatar":
		r.References.AvatarID = NullUUID
	default:
		return errors.New("Unknown asset reference type")
	}
	return nil
}
