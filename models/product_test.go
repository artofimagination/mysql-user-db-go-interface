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

	detailsID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create asset uuid %s", err)
		return
	}

	Expected := Product{
		Name:      "TestProduct",
		AssetsID:  assetsID,
		DetailsID: detailsID,
		Public:    true,
	}
	Interface = RepoInterface{}

	// Execute test
	product, err := Interface.NewProduct(
		"TestProduct",
		true,
		&detailsID,
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
