package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"

	ProductDescription = "product_description"
)

var DefaultImagePath = ""
var DefaultURL = ""

// Errors called in multiple places (for example in unittests).

var ErrAssetRefNotInitialised = errors.New("Asset references not initialised")

type Asset struct {
	ID         uuid.UUID  `validation:"required"`
	References References `validation:"required"`
	Path       string
}

// Assets structure contains the identification of all user related documents images.
type References map[string]string

func (RepoInterface) NewAsset(references References, generatePath func(assetID *uuid.UUID) string) (*Asset, error) {
	var a Asset

	if references == nil {
		return nil, ErrAssetRefNotInitialised
	}

	newID, err := UUIDInterface.NewUUID()
	if err != nil {
		return nil, err
	}

	a.ID = newID
	a.References = references
	a.Path = generatePath(&a.ID)

	return &a, nil
}

func (r *Asset) GetImagePath(typeString string) string {
	path, ok := r.References[typeString]
	if !ok {
		return DefaultImagePath
	}

	return path
}

func (r *Asset) SetImagePath(typeString string) error {
	if _, ok := r.References[typeString]; ok {
		return nil
	}

	newID, err := UUIDInterface.NewUUID()
	if err != nil {
		return err
	}

	r.References[typeString] = fmt.Sprintf("%s/%s.jpg", r.Path, newID.String())

	return nil
}

func (r *Asset) SetURL(typeString string, url string) {
	r.References[typeString] = url
}

func (r *Asset) GetURL(typeString string) string {
	path, ok := r.References[typeString]
	if !ok {
		return DefaultURL
	}
	return path
}

func (r *Asset) ClearAsset(typeString string) error {
	if _, ok := r.References[typeString]; !ok {
		return errors.New("Unknown asset reference type")
	}
	delete(r.References, typeString)
	return nil
}
