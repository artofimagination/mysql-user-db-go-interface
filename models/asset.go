package models

import (
	"errors"
	"path"

	"github.com/google/uuid"
)

const (
	BaseAssetPath = "base_asset_path"
)

// Errors called in multiple places (for example in unittests).

var ErrAssetRefNotInitialised = errors.New("Asset references not initialised")

type Asset struct {
	ID      uuid.UUID `validation:"required"`
	DataMap DataMap   `validation:"required"`
}

// Assets structure contains the identification of all user related documents images.
type DataMap map[string]interface{}

func (*RepoInterface) NewAsset(references DataMap, generatePath func(assetID *uuid.UUID) (string, error)) (*Asset, error) {
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
	a.DataMap[BaseAssetPath], err = generatePath(&a.ID)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *Asset) GetFilePath(typeString string, defaultPath string) string {
	path, ok := r.DataMap[typeString].(string)
	if !ok {
		return defaultPath
	}

	return path
}

func (r *Asset) SetFilePath(typeString string, extension string) error {
	if _, ok := r.DataMap[typeString]; ok {
		return nil
	}

	newID, err := UUIDImpl.NewUUID()
	if err != nil {
		return err
	}

	r.DataMap[typeString] = path.Join(r.DataMap[BaseAssetPath].(string), newID.String()+extension)

	return nil
}

func (r *Asset) SetField(typeString string, url string) {
	r.DataMap[typeString] = url
}

func (r *Asset) GetField(typeString string, defaultURL string) string {
	path, ok := r.DataMap[typeString].(string)
	if !ok {
		return defaultURL
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
