package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const (
	BaseAssetPath  = "base_asset_path"
	UserAvatar     = "user_avatar"
	UserBackground = "user_background"

	ProductDescription = "product_description"
)

var DefaultImagePath = ""
var DefaultURL = ""

// Errors called in multiple places (for example in unittests).

var ErrAssetRefNotInitialised = errors.New("Asset references not initialised")

type Asset struct {
	ID      uuid.UUID `validation:"required"`
	DataMap DataMap   `validation:"required"`
}

// Assets structure contains the identification of all user related documents images.
type DataMap map[string]interface{}

func (RepoInterface) NewAsset(references DataMap, generatePath func(assetID *uuid.UUID) string) (*Asset, error) {
	var a Asset

	if references == nil {
		return nil, ErrAssetRefNotInitialised
	}

	newID, err := UUIDImpl.NewUUID()
	if err != nil {
		return nil, err
	}

	a.ID = newID
	a.DataMap = references
	a.DataMap[BaseAssetPath] = generatePath(&a.ID)

	return &a, nil
}

func (r *Asset) GetImagePath(typeString string) string {
	path, ok := r.DataMap[typeString].(string)
	if !ok {
		return DefaultImagePath
	}

	return path
}

func (r *Asset) SetImagePath(typeString string) error {
	if _, ok := r.DataMap[typeString]; ok {
		return nil
	}

	newID, err := UUIDImpl.NewUUID()
	if err != nil {
		return err
	}

	r.DataMap[typeString] = fmt.Sprintf("%s/%s.jpg", r.DataMap[BaseAssetPath], newID.String())

	return nil
}

func (r *Asset) SetURL(typeString string, url string) {
	r.DataMap[typeString] = url
}

func (r *Asset) GetURL(typeString string) string {
	path, ok := r.DataMap[typeString].(string)
	if !ok {
		return DefaultURL
	}
	return path
}

func (r *Asset) ClearAsset(typeString string) error {
	if _, ok := r.DataMap[typeString]; !ok {
		return errors.New("Unknown asset reference type")
	}
	delete(r.DataMap, typeString)
	return nil
}
