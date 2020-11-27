package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewProduct_ValidInit(t *testing.T) {
	// Create test data
	assetsID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create asset uuid %s", err)
		return
	}
	details := make(Details)
	Expected := Product{
		Name:     "TestProduct",
		AssetsID: assetsID,
		Details:  details,
		Public:   true,
	}
	Interface = RepoInterface{}

	// Execute test
	product, err := Interface.NewProduct(
		"TestProduct",
		true,
		details,
		&assetsID,
	)
	if err != nil {
		t.Errorf("Failed to create new asset %s", err)
		return
	}

	Expected.ID = product.ID
	if !cmp.Equal(*product, Expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", *product, Expected)
		return
	}
}

func TestNewProduct_NilDetails(t *testing.T) {
	// Create test data
	assetsID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create asset uuid %s", err)
		return
	}
	var details Details
	Interface = RepoInterface{}

	// Execute test
	_, err = Interface.NewProduct(
		"TestProduct",
		true,
		details,
		&assetsID,
	)
	if err == nil || err.Error() != ErrProductDetailsNotInitialised {
		t.Errorf("Failed to create new product %s", err)
		return
	}
}
