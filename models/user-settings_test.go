package models

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestNewUserSettings_ValidInit(t *testing.T) {
	// Create test data
	settings := make(Settings)
	newID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to create settings uuid %s", err)
		return
	}

	Expected := UserSettings{
		ID:       newID,
		Settings: settings,
	}

	Interface = RepoInterface{}

	UUIDInterface = UUIDInterfaceMock{
		uuidMock: newID,
	}

	// Execute test
	userSettings, err := Interface.NewUserSettings(settings)
	if err != nil {
		t.Errorf("Failed to create new user settings %s", err)
		return
	}

	if !cmp.Equal(*userSettings, Expected) {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", *userSettings, Expected)
		return
	}
}

func TestNewUserSettings_NilSettings(t *testing.T) {
	// Create test data
	var settings Settings
	Interface = RepoInterface{}

	// Execute test
	_, err := Interface.NewUserSettings(settings)
	if err == nil || err.Error() != ErrSettingNotInitialised {
		t.Errorf("\nTest returned:\n %+v\nExpected:\n %+v", err, ErrSettingNotInitialised)
		return
	}
}
