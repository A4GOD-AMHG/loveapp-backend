package repository

import (
	"testing"

	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/models"
	"github.com/A4GOD-AMHG/LoveApp-Backend/internal/testhelpers"
)

func TestDevicePushTokenUpsertAndList(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	userID := testhelpers.InsertTestUser(t, "anyel", "Anyel", "hash")
	repo := NewDevicePushTokenRepository()

	err := repo.Upsert(&models.DevicePushToken{
		UserID:     userID,
		Platform:   "android",
		PushToken:  "token-1",
		DeviceName: "Pixel 8",
	})
	if err != nil {
		t.Fatalf("error al crear token push: %v", err)
	}

	err = repo.Upsert(&models.DevicePushToken{
		UserID:     userID,
		Platform:   "android",
		PushToken:  "token-2",
		DeviceName: "Pixel 8",
	})
	if err != nil {
		t.Fatalf("error al actualizar token push: %v", err)
	}

	tokens, err := repo.FindByUserID(userID)
	if err != nil {
		t.Fatalf("error al consultar tokens push: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("se esperaba 1 token push, se obtuvieron %d", len(tokens))
	}

	if tokens[0].PushToken != "token-2" {
		t.Fatalf("se esperaba token actualizado 'token-2', se obtuvo '%s'", tokens[0].PushToken)
	}
}

func TestDevicePushTokenDeleteByToken(t *testing.T) {
	cleanup := testhelpers.SetupTestDB(t)
	defer cleanup()

	userID := testhelpers.InsertTestUser(t, "alexis", "Alexis", "hash")
	repo := NewDevicePushTokenRepository()

	if err := repo.Upsert(&models.DevicePushToken{
		UserID:     userID,
		Platform:   "ios",
		PushToken:  "token-ios",
		DeviceName: "iPhone",
	}); err != nil {
		t.Fatalf("error al crear token push: %v", err)
	}

	if err := repo.DeleteByToken(userID, "token-ios"); err != nil {
		t.Fatalf("error al eliminar token push: %v", err)
	}

	tokens, err := repo.FindByUserID(userID)
	if err != nil {
		t.Fatalf("error al consultar tokens push: %v", err)
	}

	if len(tokens) != 0 {
		t.Fatalf("se esperaban 0 tokens push, se obtuvieron %d", len(tokens))
	}
}
