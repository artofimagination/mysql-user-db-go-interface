package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestSetImagePath(t *testing.T) {
	// Create test data
	asset := Asset{}
	asset.References = make(References)
	asset.Path = "test/path"
	Interface = RepoInterface{}

	// Execute test
	err := asset.SetImagePath("test")
	if err != nil {
		t.Errorf("Failed to set path %s", err)
		return
	}

	result := strings.Split(asset.References["test"], "/")
	for index, slice := range result {
		switch index {
		case 0:
			if slice != "test" {
				t.Errorf("Invalid path returned %s", asset.References["test"])
				return
			}
		case 1:
			if slice != "path" {
				t.Errorf("Invalid path returned %s", asset.References["test"])
				return
			}
		case len(result) - 1:
			uuidString := strings.Split(result[len(result)-1], ".")
			_, err = uuid.Parse(uuidString[0])
			if err != nil {
				t.Errorf("Invalid path returned %s", asset.References["test"])
				return
			}
		}
	}
}

func TestGetPath_ValidKeyUUIDValue(t *testing.T) {
	// Create test data
	asset := Asset{}
	asset.References = make(References)
	testID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to generate test uuid %s", err)
		return
	}

	asset.Path = "test/path"
	expected := fmt.Sprintf("%s/%s", asset.Path, testID.String())
	asset.References["test"] = expected

	// Execute test
	path := asset.GetImagePath("test")

	if !cmp.Equal(path, expected) {
		t.Errorf("Test returned:\n %+v\nExpected:\n %+v", path, expected)
		return
	}
}

func TestGetPath_InvalidKeyUUIDValue(t *testing.T) {
	// Create test data
	asset := Asset{}
	asset.References = make(References)
	DefaultImagePath = "default/path/image.jpg"

	asset.Path = "test/path"

	// Execute test
	path := asset.GetImagePath("test")

	if !cmp.Equal(path, DefaultImagePath) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", path, DefaultImagePath)
		return
	}
}

func TestGetURL(t *testing.T) {
	// Create test data
	asset := Asset{}
	asset.References = make(References)
	DefaultURL = "http://test.com"

	// Execute test
	path := asset.GetURL("test")

	if !cmp.Equal(path, DefaultURL) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", path, DefaultURL)
		return
	}
}

func TestNewAsset_ValidInit(t *testing.T) {
	// Create test data
	references := make(References)
	expected := Asset{
		References: references,
		Path:       "testPath",
	}

	Interface = RepoInterface{}

	// Execute test
	asset, err := Interface.NewAsset(
		references,
		func(assetID *uuid.UUID) string {
			return "testPath"
		})
	if err != nil {
		t.Errorf("Failed to create new asset %s", err)
		return
	}

	expected.ID = asset.ID
	if !cmp.Equal(*asset, expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", *asset, expected)
		return
	}
}

func TestNewAsset_NilReferences(t *testing.T) {
	// Create test data
	var references References

	Interface = RepoInterface{}

	// Execute test
	_, err := Interface.NewAsset(
		references,
		func(assetID *uuid.UUID) string {
			return "testPath"
		})
	if err == nil || err.Error() != ErrAssetRefNotInitialised {
		t.Errorf("Failed to create new asset %s", err)
		return
	}
}
