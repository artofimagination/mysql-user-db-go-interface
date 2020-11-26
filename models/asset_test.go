package models

import (
	"fmt"
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

	newID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create details uuid %s", err)
		return
	}

	UUIDInterface = UUIDInterfaceMock{
		uuidMock: newID,
	}

	// Execute test
	err = asset.SetImagePath("test")
	if err != nil {
		t.Errorf("Failed to set path %s", err)
		return
	}

	expected := fmt.Sprintf("%s/%s.jpg", asset.Path, newID.String())
	if !cmp.Equal(asset.References["test"], expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", asset.References["test"], expected)
		return
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
	newID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create details uuid %s", err)
		return
	}
	expected := Asset{
		ID:         newID,
		References: references,
		Path:       "testPath",
	}

	Interface = RepoInterface{}
	UUIDInterface = UUIDInterfaceMock{
		uuidMock: newID,
	}

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
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrAssetRefNotInitialised)
		return
	}
}
