package models

import (
	"errors"
	"path"
	"strings"

	"github.com/google/uuid"
)

const (
	BaseAssetPath = "base_asset_path"
)

// Errors called in multiple places (for example in unittests).

var ErrAssetRefNotInitialised = errors.New("Asset references not initialised")

type Asset struct {
	ID      uuid.UUID `json:"id" validate:"required"`
	DataMap DataMap   `json:"datamap" validate:"required"`
}

// Assets structure contains the identification of all user related documents images.
type DataMap map[string]interface{}

func (f *RepoFunctions) NewAsset(references DataMap) (*Asset, error) {
	var a Asset

	if references == nil {
		return nil, ErrAssetRefNotInitialised
	}

	newID, err := f.UUIDImpl.NewUUID()
	if err != nil {
		return nil, err
	}

	a.ID = newID
	a.DataMap = references
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (f *RepoFunctions) GetFilePath(asset *Asset, typeString string, defaultPath string) string {
	path, ok := asset.DataMap[typeString].(string)
	if !ok {
		return defaultPath
	}

	return path
}

func (f *RepoFunctions) SetFilePath(asset *Asset, typeString string, extension string) error {
	if _, ok := asset.DataMap[typeString]; ok {
		return nil
	}

	newID, err := f.UUIDImpl.NewUUID()
	if err != nil {
		return err
	}

	asset.DataMap[typeString] = strings.Join([]string{path.Join(asset.DataMap[BaseAssetPath].(string), newID.String()), extension}, "")

	return nil
}

func (f *RepoFunctions) SetField(asset *Asset, typeString string, field interface{}) {
	asset.DataMap[typeString] = field
}

func (f *RepoFunctions) GetField(asset *Asset, typeString string, defaultValue string) interface{} {
	field, ok := asset.DataMap[typeString]
	if !ok {
		return defaultValue
	}
	return field
}

func (f *RepoFunctions) ClearAsset(asset *Asset, typeString string) error {
	if _, ok := asset.DataMap[typeString]; !ok {
		return errors.New("Unknown asset reference type")
	}
	delete(asset.DataMap, typeString)
	return nil
}
