package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewProductDetails_ValidInit(t *testing.T) {
	// Create test data
	details := make(Details)
	newID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create details uuid %s", err)
		return
	}

	expected := ProductDetails{
		ID:      newID,
		Details: details,
	}

	Interface = RepoInterface{}

	UUIDInterface = UUIDInterfaceMock{
		uuidMock: newID,
	}

	// Execute test
	productDetails, err := Interface.NewProductDetails(details)
	if err != nil {
		t.Errorf("Failed to create new product details %s", err)
		return
	}

	if !cmp.Equal(*productDetails, expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", *productDetails, expected)
		return
	}
}

func TestNewProductDetails_NilDetails(t *testing.T) {
	// Create test data
	var details Details
	Interface = RepoInterface{}

	// Execute test
	_, err := Interface.NewProductDetails(details)
	if err == nil || err.Error() != ErrProductDetailsNotInitialised {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrProductDetailsNotInitialised)
		return
	}
}
